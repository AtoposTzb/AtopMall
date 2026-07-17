import sys
import os
import grpc
import consul
import json

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from userop_srv.settings import settings

class TestMessageServicer:
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
        #TODO:先暂时留空

if __name__ == "__main__":
    test = TestMessageServicer()
