import consul
import requests
from common.register.base import Register

class ConsulRegister(Register):
    def __init__(self,host,port):
        self.host = host
        self.port = port
        #实例化一个consul客户端
        self.client = consul.Consul(host=host,port=port)
    def register(self, name, id, address, port,tags,check)->bool:   
        if check is None:
            check = { #健康检查
                    "GRPC":f"{address}:{port}",
                    "GRPCUseTLS":False,
                    "Timeout":"5s",
                    "Interval":"5s",
                    "DeregisterCriticalServiceAfter":"15s"
                }
        else:
            check = check
        #注册服务
        return  self.client.agent.service.register(name=name,service_id=id,
                            address=address,port=port,tags=tags,check=check)

        
    def deregister(self,service_id):
        return self.client.agent.service.deregister(service_id)
    
    def get_all_services(self):
        return self.client.agent.services()
    
    def filter_service(self,filter):
        url = f"http://{self.host}:{self.port}/v1/agent/services"
        params = {
            "filter":filter
        }
        return requests.get(url,params=params).json()