import sys
import os
import grpc
import consul
import json

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0,BASE_DIR)
from inventory_srv.proto import inventory_pb2,inventory_pb2_grpc
from inventory_srv.settings import settings
from google.protobuf import empty_pb2

class TestInventoryServicer:
    def __init__(self):
        #连接grpc服务器

        c = consul.Consul(host="192.168.1.106",port=8500)
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
        self.inventory_stub = inventory_pb2_grpc.InventoryStub(channel)

    def set_inv(self):
        rsp = self.inventory_stub.SetInv(inventory_pb2.GoodsInvInfo(goodsId=10,num=120))
    
    def get_inv(self):
        rsp = self.inventory_stub.InvDetail(inventory_pb2.GoodsInvInfo(goodsId=10))
        print(rsp.num)
    
    def sell(self):
        goods_list = [(1,10),(2,6)]
        request = inventory_pb2.SellInfo()
        for goodsId,num in goods_list:
            request.goodsInfo.append(inventory_pb2.GoodsInvInfo(goodsId=goodsId,num=num))
        rsp = self.inventory_stub.SellInv(request)
    
    def reback(self):
        goods_list = [(1,6),(2,3)]
        request = inventory_pb2.SellInfo()
        for goodsId,num in goods_list:
            request.goodsInfo.append(inventory_pb2.GoodsInvInfo(goodsId=goodsId,num=num))
        rsp = self.inventory_stub.RebackInv(request)
if __name__ == "__main__":
    test = TestInventoryServicer()
    # test.set_inv()
    # test.get_inv()
    # test.sell()
    test.reback()
