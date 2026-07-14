import grpc
import os
import sys
from loguru import logger
from peewee import DoesNotExist
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from model.models import ShoppingCart,OrderInfo,OrderGoods
from proto import order_pb2,order_pb2_grpc
from settings.settings import DB

# 订单服务 - 购物车相关接口
class ShoppingCartServicer(order_pb2_grpc.ShoppingCartServicer):
    @logger.catch
    def CartItemList(self,request,context):
        #获取用户购物车列表
        items = ShoppingCart.select().where(ShoppingCart.user==request.id)
        rsp = order_pb2.CartItemListResponse(total=items.count())
        for item in items:
            rsp.data.append(order_pb2.ShopCartInfoResponse(
                id=item.id,
                userId=item.user,
                goodsId=item.goods,
                nums=item.nums,
                checked=item.checked
            ))
        return rsp
    @logger.catch
    def CreateCartItem(self,request,context):
        #添加商品到购物车
        #如果商品记录已经存在则合并购物车
        existed_items = ShoppingCart.select().where(ShoppingCart.user==request.userId,ShoppingCart.goods==request.goodsId)
        if existed_items:
            item = existed_items[0] 
            item.nums += request.nums
        else:
            item = ShoppingCart(
                user=request.userId,
                goods=request.goodsId,
                nums=request.nums,
            )
        item.save()
        return order_pb2.ShopCartInfoResponse(id=item.id) #只要返回id即可
    
    @logger.catch
    def UpdateCartItem(self,request,context):
        #修改购物车条目信息 - 数量和选中状态
        try:
            item = ShoppingCart.get(ShoppingCart.user==request.userId,ShoppingCart.goods==request.goodsId)
            item.checked = request.checked
            if request.nums:
                item.nums = request.nums
            item.save()
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("购物车条目不存在")
            return empty_pb2.Empty()
    @logger.catch
    def DeleteCartItem(self,request,context):
        #删除购物车条目
        try:
            item = ShoppingCart.get(ShoppingCart.user==request.userId,ShoppingCart.goods==request.goodsId)
            item.delete_instance()
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("购物车条目不存在")
            return empty_pb2.Empty()