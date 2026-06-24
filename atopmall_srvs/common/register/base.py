import abc  #导入abc模块,用于创建抽象类和抽象方法
class Register(metaclass=abc.ABCMeta): #abc.ABCMeta是一个元类,用于创建抽象类,抽象类不能被实例化,只能被继承
    #这是一个抽象方法,必须在子类中实现 @abc.abstractmethod 是一个装饰器,用于标记一个方法为抽象方法
    @abc.abstractmethod 
    def register(self,name,id,address,port,tags,check):
        pass

    @abc.abstractmethod 
    def deregister(self,service_id):
        pass

    @abc.abstractmethod 
    def get_all_services(self):
        pass

    @abc.abstractmethod 
    def filter_service(self,filter):
        pass    
    
    
"""
抽象类,用于定义服务注册器的基本接口
    1. register: 注册服务
    2. deregister: 取消注册服务
    3. get_all_services: 获取所有服务
    4. filter_service: 过滤服务
在后续可以根据consul 或者 etcd 等服务发现工具来实现具体的服务注册器类
"""