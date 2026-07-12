import grpc
import os
import sys
from loguru import logger
from peewee import DoesNotExist
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from model.models import Inventory
from proto import inventory_pb2,inventory_pb2_grpc
from settings.settings import DB

# 库存服务
class InventoryServicer(inventory_pb2_grpc.InventoryServicer):
    @logger.catch
    def SetInv(self, request:inventory_pb2.GoodsInvInfo, context):
        #设置库存
        force_insert = False  #默认更新库存
        invs = Inventory.select().where(Inventory.goods == request.goodsId) 
        if not invs:
            inv = Inventory()
            inv.goods = request.goodsId
            force_insert = True #如果不存在库存记录，则插入新记录
        else:
            inv = invs[0]  #为什么取第一个？因为一个商品只能有一个库存记录
        inv.stocks = request.num
        inv.save(force_insert=force_insert) #force_insert=True 表示如果存在则更新，否则插入 所以这个接口可以设置库存也可以更新库存
        return empty_pb2.Empty()

    @logger.catch
    def InvDetail(self, request:inventory_pb2.GoodsInvInfo, context):
        #获取库存详情
        try:
            inv = Inventory.get(Inventory.goods == request.goodsId)
            return inventory_pb2.GoodsInvInfo(goodsId=request.goodsId, num=inv.stocks)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("商品不存在库存记录")
            return inventory_pb2.GoodsInvInfo()
        
    @logger.catch
    def SellInv(self, request:inventory_pb2.SellInfo, context):
        #扣减库存 超卖问题 事务处理:执行多个sql是原子性的
        with DB.atomic() as txn: #peewee使用atomic()方法开启事务
            for item in request.goodsInfo:
                try:
                    inv = Inventory.get(Inventory.goods == item.goodsId)
                except DoesNotExist:
                    txn.rollback() #回滚事务
                    context.set_code(grpc.StatusCode.NOT_FOUND)
                    context.set_details("商品不存在库存记录")
                    return empty_pb2.Empty() 
                if inv.stocks < item.num:
                    context.set_code(grpc.StatusCode.RESOURCE_EXHAUSTED)
                    context.set_details("库存不足")
                    txn.rollback() #回滚事务
                    return empty_pb2.Empty()
                else:
                    #TODO: 超卖问题 可能引起数据不一致-分布式锁解决
                    inv.stocks -= item.num
                    inv.save()
            return empty_pb2.Empty()
        
    @logger.catch
    def RebackInv(self, request:inventory_pb2.SellInfo, context):
        #库存归还 几种种情况：1.订单超时会自动归还 2.订单创建失败，需要归还之前的库存 3.手动归还库存
        with DB.atomic() as txn: #peewee使用atomic()方法开启事务
            for item in request.goodsInfo:
                try:
                    inv = Inventory.get(Inventory.goods == item.goodsId)
                except DoesNotExist:
                    txn.rollback() #回滚事务
                    context.set_code(grpc.StatusCode.NOT_FOUND)
                    context.set_details("商品不存在库存记录")
                    return empty_pb2.Empty() 
                #TODO: 可能引起数据不一致-分布式锁解决
                inv.stocks += item.num
                inv.save()
            return empty_pb2.Empty()