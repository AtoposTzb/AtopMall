import os
import sys
from datetime import datetime
from peewee import *
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from settings.settings import DB
# from playhouse.pool import PooledMySQLDatabase
# from playhouse.shortcuts import ReconnectMixin

# class ReconnectMysqlDatebase(ReconnectMixin,PooledMySQLDatabase):
#     pass

# db = ReconnectMysqlDatebase(database="atopmall_order_srv",host="192.168.1.106",port=3306,user="root",password="123456")

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

class ShoppingCart(BaseModel):
    """
    购物车
    """
    user = IntegerField(verbose_name="用户id")
    goods = IntegerField(verbose_name="商品id")
    nums = IntegerField(verbose_name="购买数量")
    checked = BooleanField(default=True, verbose_name="是否选中")


class OrderInfo(BaseModel):
    """
    订单
    """
    ORDER_STATUS = (
        ("TRADE_SUCCESS", "成功"),
        ("TRADE_CLOSED", "超时关闭"),
        ("WAIT_BUYER_PAY", "交易创建"),
        ("TRADE_FINISHED", "交易结束"),
    )

    PAY_TYPE = (
        ("alipay", "支付宝"),
    )

    user = IntegerField(verbose_name="用户id")
    order_sn = CharField(max_length=30, null=True, unique=True, verbose_name="订单号")
    pay_type = CharField(choices=PAY_TYPE, default="alipay", max_length=30, verbose_name="支付方式")
    status = CharField(choices=ORDER_STATUS, default="paying", max_length=30, verbose_name="订单状态")
    trade_no = CharField(max_length=100, unique=True, null=True, verbose_name=u"交易号") #支付宝的交易号
    order_mount = FloatField(default=0.0, verbose_name="订单金额")
    pay_time = DateTimeField(null=True, verbose_name="支付时间")

    # 用户信息
    address = CharField(max_length=100, default="", verbose_name="收货地址")
    signer_name = CharField(max_length=20, default="", verbose_name="签收人")
    singer_mobile = CharField(max_length=11, verbose_name="联系电话")
    post = CharField(max_length=200, default="", verbose_name="留言")


class OrderGoods(BaseModel):
    """
    订单的商品详情
    """
    order = IntegerField(verbose_name="订单id")
    goods = IntegerField(verbose_name="商品id")
    goods_name = CharField(max_length=20, default="", verbose_name="商品名称")
    goods_image = CharField(max_length=200, default="", verbose_name="商品图片")
    goods_price = DecimalField(verbose_name="商品价格") #镜像
    nums = IntegerField(default=0, verbose_name="商品数量")


if __name__ == "__main__":
    DB.create_tables([ShoppingCart, OrderInfo, OrderGoods])
    # ShoppingCart.delete().where(ShoppingCart.user == 1, ShoppingCart.checked == True)