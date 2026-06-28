import sys
import os
import grpc
import consul
import json

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0,BASE_DIR)
from goods_srv.proto import goods_pb2,goods_pb2_grpc
from goods_srv.settings import settings
from google.protobuf import empty_pb2

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
        self.goods_stub = goods_pb2_grpc.GoodsStub(channel) 
        self.category_stub = goods_pb2_grpc.CategoryStub(channel)
        
    def goods_list(self):
        #获取商品列表
        #测试获取riceMin=50也就是50元以上商品
        # rsp :goods_pb2.GoodsListResponse = self.goods_stub.GoodsList(goods_pb2.GoodsFilterRequest(priceMin=50))
        #测试获取topCategory=135485也就是level3分类的商品
        # rsp :goods_pb2.GoodsListResponse = self.goods_stub.GoodsList(goods_pb2.GoodsFilterRequest(topCategory=135485))
        #测试获取topCategory=135485也就是level1分类的商品
        # rsp :goods_pb2.GoodsListResponse = self.goods_stub.GoodsList(goods_pb2.GoodsFilterRequest(topCategory=130358))
        rsp :goods_pb2.GoodsListResponse = self.goods_stub.GoodsList(goods_pb2.GoodsFilterRequest(keyWords="四川"))
        for goods in rsp.data:
            print(goods.name,goods.shopPrice)
    
    def batch_get_goods(self):
        ids = [421,422]
        rsp :goods_pb2.GoodsListResponse = self.goods_stub.BatchGetGoods(goods_pb2.BatchGoodsIdInfo(id=ids))
        for goods in rsp.data:
            print(goods.name,goods.shopPrice)
        
    def get_goods_detail(self):
        rsp:goods_pb2.GoodsDetailResponse = self.goods_stub.GetGoodsDetail(goods_pb2.GoodInfoRequest(id=421))
        print(rsp.name)
        
    def get_category_list(self):
        rsp:goods_pb2.CategoryListResponse = self.category_stub.GetAllCategorysList(empty_pb2.Empty())
        data = json.loads(rsp.jsonData)
        print(data)
        
if __name__ == "__main__":
    test = TestGoodsServicer()
    # test.goods_list()
    # test.batch_get_goods()
    # test.get_goods_detail()
    test.get_category_list()
