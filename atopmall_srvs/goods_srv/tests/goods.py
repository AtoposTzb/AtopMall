import sys
import os
import grpc
import consul

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0,BASE_DIR)
from goods_srv.proto import goods_pb2,goods_pb2_grpc
from goods_srv.settings import settings

class TestGoodsServicer:
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
        self.stub = goods_pb2_grpc.GoodsStub(channel) 
        
    def goods_list(self):
        #获取商品列表
        #测试获取riceMin=50也就是50元以上商品
        # rsp :goods_pb2.GoodsListResponse = self.stub.GoodsList(goods_pb2.GoodsFilterRequest(priceMin=50))
        #测试获取topCategory=135485也就是level3分类的商品
        # rsp :goods_pb2.GoodsListResponse = self.stub.GoodsList(goods_pb2.GoodsFilterRequest(topCategory=135485))
        #测试获取topCategory=135485也就是level1分类的商品
        # rsp :goods_pb2.GoodsListResponse = self.stub.GoodsList(goods_pb2.GoodsFilterRequest(topCategory=130358))
        rsp :goods_pb2.GoodsListResponse = self.stub.GoodsList(goods_pb2.GoodsFilterRequest(keyWords="四川"))
        for goods in rsp.data:
            print(goods.name,goods.shopPrice)
          
        
if __name__ == "__main__":
    test = TestGoodsServicer()
    test.goods_list()
