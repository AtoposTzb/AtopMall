import sys
import os
import grpc
import consul
import json

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0,BASE_DIR)
from order_srv.proto import order_pb2,order_pb2_grpc
from order_srv.settings import settings
from google.protobuf import empty_pb2

class TestOrderServicer:
    def __init__(self):
        #连接grpc服务器

        c = consul.Consul(host=settings.CONSUL_HOST,port=settings.CONSUL_PORT)
        services = c.agent.services() #直接获取所有的服务
        ip = ""
        port = ""
        for key,value in services.items():
            if value["Service"] == settings.SERVICE_NAME:
                ip = value["Address"]
                port = value["Port"]
                break
        if not ip or not port:
            raise Exception("未找到服务")
        
        channel = grpc.insecure_channel(f"{ip}:{port}")
        self.order_stub = order_pb2_grpc.OrderStub(channel)  # 订单服务
        self.cart_stub = order_pb2_grpc.ShoppingCartStub(channel) #购物车服务
    
    def create_cart_item(self):
        rsp = self.cart_stub.CreateCartItem(
            order_pb2.CartItemRequest(goodsId=422,userId=11,nums=15)
        )
        print(rsp)
    
    def create_order(self):
        rsp = self.order_stub.CreateOrder(
            order_pb2.OrderRequest(
                userId=11,
                address="云南菌子地",
                mobile="13800000001",
                name="张三",
                post="快发货"
                )
        )
        print(rsp)
    
    def order_list(self):
        rsp = self.order_stub.OrderList(
            order_pb2.OrderFilterRequest(userId=11)
        )
        print(rsp)

if __name__ == "__main__":
    test = TestOrderServicer()
    # test.set_inv()
    # test.get_inv()
    # test.sell()
    # test.create_cart_item()
    # test.create_order()
    test.order_list()
