import os
import sys
from datetime import datetime
from peewee import *
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from settings.settings import DB

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
        database = DB

#简化不写
# class Stock(BaseModel):
#     #仓库表
#     name = CharField(verbose_name="仓库名")
#     address = CharField(verbose_name="仓库地址")

class Inventory(BaseModel):
    #商品库存表
    # stock = PrimaryKeyField(Stock)
    goods = IntegerField(verbose_name="商品id",unique=True)
    stocks = IntegerField(verbose_name="库存数量",default=0)
    version = IntegerField(verbose_name="版本号",default=0) #分布式锁的乐观锁

#根据服务还可以多设计一张表来处理库存的扣减和归还操作
# class InventoryLog(BaseModel):
#     #库存日志表
#     goods = IntegerField(verbose_name="商品id")
#     num = IntegerField(verbose_name="库存数量")
#     version = IntegerField(verbose_name="版本号",default=0) #分布式锁的乐观锁
#     create_time = DateTimeField(default=datetime.now,verbose_name="创建时间")
#     update_time = DateTimeField(default=datetime.now,verbose_name="更新时间")
#     is_deleted = BooleanField(default=False,verbose_name="是否删除")

if __name__ == "__main__":
    DB.create_tables([Inventory])
    # for i in range(421,841):
    #     inv = Inventory(goods=i,stocks = 100)
    #     inv.save()