import grpc
import os
import sys
import time
import json
import threading
from random import Random
from loguru import logger
from datetime import datetime
from peewee import DoesNotExist
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from model.models import ShoppingCart,OrderInfo,OrderGoods
from proto import order_pb2,order_pb2_grpc,goods_pb2,goods_pb2_grpc,inventory_pb2,inventory_pb2_grpc
from settings import settings
from common.register.consul import ConsulRegister 
from rocketmq.client import Producer,TransactionMQProducer,TransactionStatus,Message,SendStatus
from rocketmq.client import ConsumeStatus

#生成订单号
def create_order_sn(user_id):
    return f'{time.strftime("%Y%m%d%H%M%S",time.localtime())}{user_id}{Random().randint(10,99)}'

#监听延时消息(超时订单消息)
def order_timeout(msg):
    msg_body_str = msg.body.decode("utf-8")
    print(f"超时消息接收时间:{datetime.now()},消息内容:{msg_body_str}")
    msg_body = json.loads(msg_body_str)
    order_sn = msg_body["orderSn"]

    #1.查询订单支付状态
    with settings.DB.atomic() as txn:
        try:
            order = OrderInfo.get(OrderInfo.order_sn==order_sn)
            if order.status != "TRADE_SUCCESS":
                order.status = "TRADE_CLOSED"
                order.save()

                #2.给库存服务发送一个归还库存的消息
                msg = Message("order_reback")
                msg.set_keys("atopmall")
                msg.set_tags("reback")
                msg.set_body(json.dumps({"orderSn":order_sn}))

                #3.给订单服务发送一个取消订单的消息
                sync_producer = Producer("order_sender")
                sync_producer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
                sync_producer.start()

                ret = sync_producer.send_sync(msg)
                if ret.status != SendStatus.OK:
                    raise Exception("发送失败")
                sync_producer.shutdown()
        except Exception as e:
            print(e)
            txn.rollback()
            return ConsumeStatus.RECONSUME_LATER #重试
    return ConsumeStatus.CONSUME_SUCCESS
        


local_execute_dict = {} #本地事物字典，用于存储半消息的执行状态
local_execute_lock = threading.Lock() #保护local_execute_dict的线程锁

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
            tmp_rsp.addTime = order.add_time.strftime("%Y-%m-%d %H:%M:%S")

            rsp.data.append(tmp_rsp)
        return rsp
    
    @logger.catch
    def OrderDetail(self,request,context):
        #获取订单详情
        rsp = order_pb2.OrderInfoDetailResponse()
        try:
            if request.userId: #如果有用户id，就根据用户id来查询订单(简单的权限校验)
                order = OrderInfo.get(OrderInfo.id==request.id,OrderInfo.user==request.userId)
            else:
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
    def check_callback(self,msg):
        #消息回查
        msg_body = json.loads(msg.body.decode("utf-8"))
        order_sn = msg_body["orderSn"]
        #查询本地数据库，看一下order_sn的订单是否已经入库了
        orders = OrderInfo.select().where(OrderInfo.order_sn==order_sn)
        if orders :
            return TransactionStatus.ROLLBACK
        else:
            return TransactionStatus.COMMIT  
        #订单不存在，本地事务失败了，库存被扣了但没订单，此时消息事物-确认通知库存服务，赶紧把库存归还了

    @logger.catch
    def local_execute(self,msg,user_args):
        #本地事物回调 用于处理半消息
        #因为涉及两张表，所以需要事务来完成
        msg_body = json.loads(msg.body.decode("utf-8"))
        order_sn = msg_body["orderSn"] 
        with local_execute_lock:
            local_execute_dict[order_sn] = {}
        with settings.DB.atomic() as txn:
        #购物车查询
            checke_goods_ids=[]
            checke_goods_nums={}
            order_allmount = 0.0
            order_goods_list=[]
            goods_sell_info = []
            for item in ShoppingCart.select().where(ShoppingCart.user==msg_body["userId"],ShoppingCart.checked==True):
                checke_goods_ids.append(item.goods)
                checke_goods_nums[item.goods]=item.nums

            if not checke_goods_ids:
                """
                格式：
                {"orderSn":{
                    "code":"",
                    "detail":"",
                }}
                """ 
                with local_execute_lock:
                    local_execute_dict[order_sn]["code"] = grpc.StatusCode.NOT_FOUND
                    local_execute_dict[order_sn]["detail"] = "请选择商品"
                return TransactionStatus.ROLLBACK #回滚事务 取消半消息
            
            #查询商品的信息 需要获取商品的服务(从consul中获取，然后grpc调用)
            goods_consul = ConsulRegister(settings.CONSUL_HOST,settings.CONSUL_PORT)
            goods_srv_address,goods_srv_port = goods_consul.get_host_port(f'Service=="{settings.GOODS_SRV_NAME}"')
            if not goods_srv_address or not goods_srv_port:
                with local_execute_lock:
                    local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                    local_execute_dict[order_sn]["detail"] = "商品服务不可用"
                return TransactionStatus.ROLLBACK  #回滚事务 取消半消息
            
            goods_channel = grpc.insecure_channel(f"{goods_srv_address}:{goods_srv_port}") #创建商品服务的channel
            goods_stub = goods_pb2_grpc.GoodsStub(goods_channel) #创建商品服务的stub 用于调用商品服务的方法

            #批量获取商品的信息
            try:
                goods_info_rsp = goods_stub.BatchGetGoods(goods_pb2.BatchGoodsIdInfo(id=checke_goods_ids))
            except Exception as e:
                with local_execute_lock:
                    local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                    local_execute_dict[order_sn]["detail"] = f"商品服务调用失败：{e}"
                return TransactionStatus.ROLLBACK  #回滚事务 取消半消息
            
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
                with local_execute_lock:
                    local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                    local_execute_dict[order_sn]["detail"] = "库存服务不可用"
                return TransactionStatus.ROLLBACK  #回滚事务 取消半消息
            inv_channel = grpc.insecure_channel(f"{inv_srv_address}:{inv_srv_port}") #创建库存服务的channel
            inv_stub = inventory_pb2_grpc.InventoryStub(inv_channel) #创建库存服务的stub 用于调用库存服务的方法

            try:
                #调用失败问题比较复杂，比如库存不足-网络问题-服务不可达等
                inv_stub.SellInv(inventory_pb2.SellInfo(orderSn = order_sn,goodsInfo = goods_sell_info))
            except Exception as e:
                with local_execute_lock:
                    local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                    local_execute_dict[order_sn]["detail"] = f"扣减库存失败：{e}"
                err_code = e.code()
                if err_code == grpc.StatusCode.UNKNOWN or err_code == grpc.StatusCode.DEADLINE_EXCEEDED:
                    return TransactionStatus.COMMIT
                else:
                    return TransactionStatus.ROLLBACK  
            #创建订单
            try:
                order = OrderInfo()
                order.user = msg_body["userId"]
                order.order_sn = order_sn
                order.order_mount = order_allmount
                order.address = msg_body["address"]
                order.signer_name = msg_body["name"]
                order.singer_mobile = msg_body["mobile"]
                order.post = msg_body["post"]
                order.save()

                #批量插入订单商品表
                for order_goods in order_goods_list:
                    order_goods.order = order.id  #id不能忘
                OrderGoods.bulk_create(order_goods_list)

                #删除购物车中的商品
                ShoppingCart.delete().where(ShoppingCart.user == msg_body["userId"], ShoppingCart.checked == True).execute()
                #返回数据
                with local_execute_lock:
                    local_execute_dict[order_sn] = {
                        "code":grpc.StatusCode.OK,
                        "detail":"下单成功",
                        "order":{
                            "id":order.id,
                            "orderSn":order_sn,
                            "total":order.order_mount
                        }
                    }

                #发送延时消息
                msg = Message("order_timeout")
                msg.set_delay_time_level(5) #设置为超时时间为1min
                msg.set_keys("atopmall")
                msg.set_tags("cancel")
                msg.set_body(json.dumps({"orderSn":order_sn}))
                sync_producer = Producer("cancel") #此处的topic 是cancel，用于取消订单，不能和之前的topic一样
                sync_producer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
                sync_producer.start()

                ret = sync_producer.send_sync(msg)
                if ret.status != SendStatus.OK:
                    raise Exception(f"发送延时消息失败")
                print(f"发送延时消息时间:{datetime.now()}")
                sync_producer.shutdown()


            except Exception as e:
                #调用库存服务的归还库存的接口就行了
                """
                    调用库存归还接口的问题：
                    详见笔记
                """
                txn.rollback() #回滚事务
                with local_execute_lock:
                    local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                    local_execute_dict[order_sn]["detail"] = f"创建订单失败：{e}"
                return TransactionStatus.COMMIT
        return TransactionStatus.ROLLBACK 


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

        #先准备好一个half消息
        producer = TransactionMQProducer("order_srv",checker_callback=self.check_callback)
        producer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
        producer.start()
        msg = Message("order_reback") #这是order_reback的topic，用于库存归还
        msg.set_keys("atopmall")
        msg.set_tags("order_reback")

        order_sn = create_order_sn(request.userId) #生成订单号，这样全局变量和高并发就不会出现数据的冲突，哪个订单对应哪个半消息

       #通过 msg_body 将local_execute回调中需要使用的数据传递过去
        msg_body={ 
            "orderSn":order_sn,
            "userId":request.userId,
            "address":request.address,
            "name":request.name,
            "mobile":request.mobile,
            "post":request.post,
        } 
        msg.set_body(json.dumps(msg_body)) #设置消息体为json字符串

        ret = producer.send_message_in_transaction(msg,self.local_execute,user_args=None) #发送半消息, 并在local_execute回调中处理
        logger.info(f"发送状态:{ret.status},消息id:{ret.msg_id}")
        if ret.status != SendStatus.OK:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details("新建订单失败")
            return order_pb2.OrderInfoResponse()
        #此处直接往下执行，会出现各种问题，因为半消息的发送是异步的，所以需要等待半消息发送完成
        #本地事物应该写在local_execute回调中，但是local_execute回调中不能有事务，否则会导致死锁，因为事务是全局的，不能在回调中开启事务
        #需要local_execute通知主函数，执行状态是什么告诉主函数，主函数根据执行状态来判断是否继续执行
        while True:
            with local_execute_lock:
                if order_sn in local_execute_dict:#监控local_execute回调执行状态
                    context.set_code(local_execute_dict[order_sn]["code"])
                    context.set_details(local_execute_dict[order_sn]["detail"])
                    #关闭事务消息
                    producer.shutdown()
                    if local_execute_dict[order_sn]["code"] == grpc.StatusCode.OK:
                        return order_pb2.OrderInfoResponse(id=local_execute_dict[order_sn]["order"]["id"],
                                                           orderSn=local_execute_dict[order_sn]["order"]["orderSn"],
                                                           total=local_execute_dict[order_sn]["order"]["total"]
                                                           )
                    else:
                        return order_pb2.OrderInfoResponse()
            time.sleep(0.1)