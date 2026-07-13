import grpc
import os
import sys
import time
from random import Random
from loguru import logger
from peewee import DoesNotExist
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from model.models import ShoppingCart,OrderInfo,OrderGoods
from proto import order_pb2,order_pb2_grpc,goods_pb2,goods_pb2_grpc,inventory_pb2,inventory_pb2_grpc
from settings import settings
from common.register.consul import ConsulRegister 

#生成订单号
def create_order_sn(user_id):
    return f'{time.strftime("%Y%m%d%H%M%S",time.localtime())}{user_id}{Random().randint(10,99)}'



# 订单服务 - 订单相关接口
class OrderServicer(order_pb2_grpc.OrderServicer):
    @logger.catch
    def OrderList(self,request,context):
        #获取订单列表
        rsp = order_pb2.OrderListResponse()
        orders = OrderInfo.select()
        if request.userId:
            orders = orders.where(OrderInfo.user==request.userId)
        rsp.total = orders.count()

        #分页
        per_page_nums = request.pagePerNums if request.pagePerNums else 10
        start = (request.pages-1)*per_page_nums if request.pages else 0
        orders = orders.limit(per_page_nums).offset(start) #从start开始取per_page_nums条数据

        for order in orders:
            tmp_rsp = order_pb2.OrderInfoResponse()
            tmp_rsp.id = order.id
            tmp_rsp.userId = order.user
            tmp_rsp.orderSn = order.order_sn
            tmp_rsp.payType = order.pay_type
            tmp_rsp.status = order.status
            tmp_rsp.post = order.post
            tmp_rsp.total = order.order_mount
            tmp_rsp.address = order.address
            tmp_rsp.name = order.signer_name
            tmp_rsp.mobile = order.singer_mobile

            rsp.data.append(tmp_rsp)
        return rsp
    
    @logger.catch
    def OrderDetail(self,request,context):
        #获取订单详情
        rsp = order_pb2.OrderInfoDetailResponse()
        try:
            order = OrderInfo.get(OrderInfo.id==request.id)
            rsp.orderInfo.id = order.id
            rsp.orderInfo.userId = order.user
            rsp.orderInfo.orderSn = order.order_sn
            rsp.orderInfo.payType = order.pay_type
            rsp.orderInfo.status = order.status
            rsp.orderInfo.post = order.post
            rsp.orderInfo.total = order.order_mount
            rsp.orderInfo.address = order.address
            rsp.orderInfo.name = order.signer_name
            rsp.orderInfo.mobile = order.singer_mobile

            order_goods = OrderGoods.select().where(OrderGoods.order==order.id)
            for order_goods in order_goods:
                order_goods_rsp = order_pb2.OrderItemResponse()
                order_goods_rsp.goodsId = order_goods.goods
                order_goods_rsp.goodsImage = order_goods.goods_image
                order_goods_rsp.goodsName = order_goods.goods_name
                order_goods_rsp.goodsPrice = float(order_goods.goods_price)
                order_goods_rsp.nums = order_goods.nums

                rsp.data.append(order_goods_rsp)
            return rsp
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("订单不存在")
            return rsp
        
    @logger.catch
    def UpdateOrderStatus(self,request,context):
        #更新订单的支付状态
        OrderInfo.update(status=request.status).where(OrderInfo.order_sn==request.OrderSn).execute()
        return empty_pb2.Empty()
    
    @logger.catch
    def CreateOrder(self,request,context):
        """
        创建订单
            1.价格 --- 访问商品服务(grpc来完成)
            2.库存的扣减 --- 访问库存服务(grpc来完成)
            3.订单的基本信息 --- 订单的商品信息表
            4.从购物车中获取到选中的商品
            5.从购物车中删除已经购买的商品
        """
        #因为涉及两张表，所以需要事务来完成
        with settings.DB.atomic() as txn:
        #购物车查询
            checke_goods_ids=[]
            checke_goods_nums={}
            order_allmount = 0.0
            order_goods_list=[]
            goods_sell_info = []
            for item in ShoppingCart.select().where(ShoppingCart.user==request.userId,ShoppingCart.checked==True):
                checke_goods_ids.append(item.goods)
                checke_goods_nums[item.goods]=item.nums

            if not checke_goods_ids:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("请选择商品")
                return order_pb2.OrderInfoResponse()
            
            #查询商品的信息 需要获取商品的服务(从consul中获取，然后grpc调用)
            goods_consul = ConsulRegister(settings.CONSUL_HOST,settings.CONSUL_PORT)
            goods_srv_address,goods_srv_port = goods_consul.get_host_port(f'Service=="{settings.GOODS_SRV_NAME}"')
            if not goods_srv_address or not goods_srv_port:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details("商品服务不可用")
                return order_pb2.OrderInfoResponse()
            
            goods_channel = grpc.insecure_channel(f"{goods_srv_address}:{goods_srv_port}") #创建商品服务的channel
            goods_stub = goods_pb2_grpc.GoodsStub(goods_channel) #创建商品服务的stub 用于调用商品服务的方法

            #批量获取商品的信息
            try:
                goods_info_rsp = goods_stub.BatchGetGoods(goods_pb2.BatchGoodsIdInfo(id=checke_goods_ids))
            except Exception as e:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"商品服务调用失败：{e}")
                return order_pb2.OrderInfoResponse()
            
            for goods_info in goods_info_rsp.data:
                order_allmount += goods_info.shopPrice * checke_goods_nums[goods_info.id]
                order_goods = OrderGoods(
                    goods=goods_info.id,
                    goods_name=goods_info.name,
                    goods_image=goods_info.goodsFrontImage,
                    goods_price=goods_info.shopPrice,
                    nums=checke_goods_nums[goods_info.id],
                )
                order_goods_list.append(order_goods)
                goods_sell_info.append(inventory_pb2.GoodsInvInfo(
                    goodsId = goods_info.id,
                    num = checke_goods_nums[goods_info.id],
                ))
            
            #扣减库存
            #这里需要负载均衡吗？ 这里已经完成了一个负载均衡里面比较简单的做法(随机挑一个商品服务)
            # - 如果深究的话，qrpc 中的dns的resolver 机制 go语言
            inv_consul = ConsulRegister(settings.CONSUL_HOST,settings.CONSUL_PORT)
            inv_srv_address,inv_srv_port = inv_consul.get_host_port(f'Service=="{settings.INVENTORY_SRV_NAME}"')
            if not inv_srv_address or not inv_srv_port:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details("库存服务不可用")
                return order_pb2.OrderInfoResponse()
            inv_channel = grpc.insecure_channel(f"{inv_srv_address}:{inv_srv_port}") #创建库存服务的channel
            inv_stub = inventory_pb2_grpc.InventoryStub(inv_channel) #创建库存服务的stub 用于调用库存服务的方法

            try:
                inv_stub.SellInv(inventory_pb2.SellInfo(goodsInfo = goods_sell_info))
            except Exception as e:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"扣减库存失败：{e}")
                return order_pb2.OrderInfoResponse()
            #创建订单
            try:
                order = OrderInfo()
                order.order_sn = create_order_sn(request.userId)
                order.order_mount = order_allmount
                order.address = request.address
                order.signer_name = request.name
                order.singer_mobile = request.mobile
                order.post = request.post
                order.save()

                #批量插入订单商品表
                for order_goods in order_goods_list:
                    order_goods.order = order.id  #id不能忘
                OrderGoods.bulk_create(order_goods_list)

                #删除购物车中的商品
                ShoppingCart.delete().where(ShoppingCart.user == request.userId, ShoppingCart.checked == True).execute()
            except Exception as e:
                txn.rollback() #回滚事务
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"创建订单失败：{e}")
                return order_pb2.OrderInfoResponse()
        #返回订单信息
        return order_pb2.OrderInfoResponse(
            id=order.id,
            orderSn=order.order_sn,
            total=order.order_mount,
        )
