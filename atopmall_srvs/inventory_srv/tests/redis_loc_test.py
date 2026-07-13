#演示基于Redis的分布式锁机制
import os
import sys
import threading
import redis
from redis_lock import Lock
from datetime import datetime
from peewee import *
from playhouse.shortcuts import ReconnectMixin
from playhouse.pool import PooledMySQLDatabase
path = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.append(path)
from settings import settings

class ReconnectMysqlDatebase(ReconnectMixin, PooledMySQLDatabase):
    pass

db = ReconnectMysqlDatebase(database="lock_inventory_test",host="192.168.1.106",port=3306,user="root",password="123456")


#删除：物理删除和逻辑删除(软删除)
#通过save方法做了修改如何确保update_time字段更新时间更新而不是add_time字段
class BaseModel(Model):
    add_time = DateTimeField(default=datetime.now,verbose_name="添加时间")
    update_time = DateTimeField(default=datetime.now,verbose_name="更新时间")
    is_deleted = BooleanField(default=False,verbose_name="是否删除")

    #重写save方法，确保update_time字段更新时间更新而不是add_time字段 拦截添加操作
    def save(self, *args, **kwargs): #*args, **kwargs是可变参数，用于接收任意数量的参数和关键字参数
        if self._pk is not None: #_pk是主键，如果主键存在，说明是更新操作，需要更新update_time字段 主键不一定是id
            self.update_time = datetime.now()
        else:
            self.add_time = datetime.now()
        return super().save(*args, **kwargs)

    #以下两个方法是删除方法，一个类方法，一个实例方法都能拦截删除操作
    @classmethod
    def delete(cls,permanently=False): #permanently表示是否永久删除
        if permanently:
            return super().delete() #表示物理删除，将数据从数据库中删除
        else:
            return super().update(is_deleted=True) #表示逻辑删除，将is_deleted字段设置为True
        
    def delete_instance(self, permanently=False,recursive=False, delete_nullable=False):
        if permanently:
            return self.delete(permanently).where(self._pk_expr()).execute() #表示物理删除，将数据从数据库中删除
        else:
            self.is_deleted = True
            return self.save() #表示逻辑删除，将is_deleted字段设置为True
    #拦截查询操作，只查询未删除的数据
    @classmethod
    def select(cls, *fields):
        return super().select(*fields).where(cls.is_deleted == False)

    class Meta:
        database = db


class Inventory(BaseModel):
    #商品库存表
    # stock = PrimaryKeyField(Stock)
    goods = IntegerField(verbose_name="商品id",unique=True)
    stocks = IntegerField(verbose_name="库存数量",default=0)
    version = IntegerField(verbose_name="版本号",default=0) #分布式锁的乐观锁机制


import uuid
class Lock_Test():
    def __init__(self,name,id=None): #初始化锁
        if id is None:
            self.id = uuid.uuid4()
        else:
            self.id = id
        self.redis_client = redis.Redis(
            host='localhost',
            port=6379,
            protocol=2  # 强制使用 RESP2 协议 兼容版本6.0.0以上
            )
        self.name = name
    
    def acquire(self): #获取锁
    # 问题1-互斥问题:锁操作的原子性->使用setnx确保获取和设置锁的操作是原子操作，不会出现多个线程同时获取到锁的情况
    #    if not self.redis_client.get(self.name): #如果大量高并发，可以出现多个线程同时进入这行代码，实现多个线程同时获取到锁的情况，需要添加一个判断，确保只有一个线程获取到锁
    #        #如果为空，说明没有其他线程占用锁，获取到锁，可以进行一个简单的设置操作
    #        self.redis_client.set(self.name,1)
    #        return True
        # if self.redis_client.setnx(self.name,1): #使用setnx确保获取和设置锁的操作是原子操作，不会出现多个线程同时获取到锁的情况，关键
        #     return True
    # 问题2-死锁问题:加锁后进程崩溃，锁永远不释放->设置过期时间，确保锁在一定时间内被释放
        if self.redis_client.set(self.name,1,nx=True,ex=15): #nx和ex参数
            #问题2引申:进程没有执行完key过期了->续租锁，15s的2/3的时候去续租
            #启动一个线程定时的刷新这个过期(续租)，这个操作最好也是使用lua脚本来完成
            return True
        else:
           #如果不为空，说明有其他线程占用锁，等待锁释放
            while True: #阻塞等待锁释放
                import time
                time.sleep(1)
                # if not self.redis_client.get(self.name):
                #     #如果为空，说明锁已被释放，获取到锁，可以进行一个简单的设置操作
                #     self.redis_client.set(self.name,1) 
                #     return True
                # if self.redis_client.setnx(self.name,1):
                #     return True
                if self.redis_client.set(self.name,1,nx=True,ex=15): #nx和ex参数
                    return True
    def release(self):
        # 问题3-误删问题:线程A删了线程B的锁 无法保证安全性 -> 添加一个唯一标识（如uuid），释放时校验
        #删除的时候先判断取出来的值是否和当前lock的id一致，一致删除，不一致报错
        # 问题引申:这里不安全:get 和 delete 是两个独立的操作->使用lua脚本将get和delete操作合并起来，确保原子性
        lock_id = self.redis_client.get(self.name)
        if lock_id == self.id:
            self.redis_client.delete(self.name)
        else:
            print(f"不能删除不属于当前线程的锁{self.name}")

def sell1():
    #演示普通锁机制(两把锁[pyhton进程锁+数据库行级排他锁])，不适合高并发场景，因为数据库行级排他锁是基于行的锁，只能锁住当前行，不能锁住其他行
    goods_list = [(1,99),(2,20),(3,30)]
    for goods_id,num in goods_list:
        # R = Lock_Test(f"lock:goods_{goods_id}")
        # R.acquire()
        #引入第三方库 redis_lock 实现 也不算三方
        #auto_renewal=True：库内部新开子线程，每隔 expire/3 秒发送一次 Redis 命令刷新锁过期时间（续租）
        R = Lock(settings.REDIS_CLIENT,f"lock:goods_{goods_id}",auto_renewal=True,expire=15)
        R.acquire()
        with db.atomic() as txn:
            try:
                inv = Inventory.get(Inventory.goods == goods_id)
                print(f"商品{goods_id}售出{num}件")
                import time
                from random import randint
                time.sleep(randint(1,3))
            except DoesNotExist:
                txn.rollback()
                print(f"商品{goods_id}不存在库存记录")
                break
            if inv.stocks < num:
                print(f"商品{goods_id}库存不足")
                txn.rollback()
                # R.release()，如果是分布式锁，那么这里就不需要释放锁，因为分布式锁是基于redis的锁，不是基于进程的锁，进程结束时会自动释放锁
                R.release()
                break
            #让数据库自己根据当前的值更新库存
            ok = Inventory.update(stocks=Inventory.stocks-num).where(Inventory.goods == goods_id).execute()
            if ok:
                print(f"商品{goods_id}库存更新成功")
            else:
                print(f"商品{goods_id}库存更新失败")
        R.release()
 # 改用with上下文，无论分支break/异常都会自动release，不会卡死BLPOP
def sell2():
    #演示普通锁机制(两把锁[pyhton进程锁+数据库行级排他锁])，不适合高并发场景，因为数据库行级排他锁是基于行的锁，只能锁住当前行，不能锁住其他行
    goods_list = [(1,99),(2,20),(3,30)]
    for goods_id,num in goods_list:
        # R = Lock_Test(f"lock:goods_{goods_id}")
        # R.acquire()
        #引入第三方库 redis_lock 实现
        R = Lock(settings.REDIS_CLIENT,f"lock:goods_{goods_id}",auto_renewal=True,expire=15)
        # with上下文管理锁，无论分支break/异常都会自动release，不会卡死BLPOP
        with R.acquire(blocking=True,timeout=10):
            with db.atomic() as txn:
                try:
                    inv = Inventory.get(Inventory.goods == goods_id)
                    print(f"商品{goods_id}售出{num}件")
                    import time
                    from random import randint
                    time.sleep(randint(1,3))
                except DoesNotExist:
                    txn.rollback()
                    print(f"商品{goods_id}不存在库存记录")
                    break
                if inv.stocks < num:
                    print(f"商品{goods_id}库存不足")
                    txn.rollback()
                    break
                #让数据库自己根据当前的值更新库存
                ok = Inventory.update(stocks=Inventory.stocks-num).where(Inventory.goods == goods_id).execute()
                if ok:
                    print(f"商品{goods_id}库存更新成功")
                else:
                    print(f"商品{goods_id}库存更新失败")
if __name__ == "__main__":
    # lock = Lock("yyds")
    # # lock.acquire()
    # lock.release()
    t1 = threading.Thread(target=sell1)
    t2 = threading.Thread(target=sell1)
    t1.start()
    t2.start()
    t1.join()
    t2.join()
    print("所有线程执行完成")