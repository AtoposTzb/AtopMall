import sys
import os
import logging
import signal
import argparse
import socket
import grpc
import uuid

from concurrent import futures
from loguru import logger

#获取当前文件的上一级目录,即goods_srv目录，方便vscode|traeIDE终端运行，其实就是找包的路径，不然会报“module not found”错误
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)

from goods_srv.proto import goods_pb2_grpc
from goods_srv.handler.goods import GoodsServicer
from goods_srv.handler.category import CategoryServicer
from goods_srv.handler.brands import BrandServicer
from goods_srv.handler.banners import BannerServicer
from goods_srv.handler.category_brand import CategoryBrandServicer
from grpc_health.v1 import health  #使用官方提供的健康检查服务，下载然后直接导入库即可使用
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc
from common.register import consul
from settings import settings
from functools import partial #偏函数,用于固定参数,返回一个新的函数,新的函数可以少传参数

#，注销服务到consul
def on_exit(sig,frame,service_id):
    deregister = consul.ConsulRegister(settings.CONSUL_HOST,settings.CONSUL_PORT)
    logger.info(f"注销{service_id}服务")
    deregister.deregister(service_id=service_id)
    logger.info("注销服务成功")
    sys.exit(0)

#动态获取可用的端口号
def get_free_port():
    s = socket.socket(socket.AF_INET,socket.SOCK_STREAM) #创建一个socket对象,并指定协议为IPv4,套接字类型为TCP
    s.bind(("",0)) #绑定到任意IP地址和端口号0，0表示自动选择一个可用的端口号
    port = s.getsockname()[1] #获取绑定的端口号，[1]表示获取端口号，[0]表示获取IP地址，
    s.close()
    return port

def server():
    #解析命令行参数,这个的作用是通过--ip和--port参数来指定服务端的IP地址和端口号
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip',
                        nargs='?',
                        type=str,
                        default="192.168.1.6",
                        help="服务端IP地址"
                        )
    parser.add_argument('--port',
                        nargs='?',
                        type=int,
                        default=0,
                        help="服务端端口号,一般默认为50051,0表示动态获取端口号"
                        )
    args = parser.parse_args()  #解析命令行参数,返回一个命名空间对象,包含所有解析后的参数
    # args.ip  获取--ip参数的值  # args.port  获取--port参数的值
    # 没有指定端口号，则动态获取端口号
    if args.port == 0:
        args.port = get_free_port()

    logger.add("logs/goods_srv_{time}.log") #将日志写入到文件夹logs下
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10)) #创建一个grpc服务器,并指定最大线程数为10
    #1.注册商品相关的服务(商品服务,分类服务,品牌服务,轮播图服务,品牌分类服务)
    goods_pb2_grpc.add_GoodsServicer_to_server(GoodsServicer(),server) #这行代码的意思是将GoodsServicer类添加到server中,并将其作为GoodsServicer服务端，简单来说就是注册GoodsServicer类
    goods_pb2_grpc.add_CategoryServicer_to_server(CategoryServicer(),server)
    goods_pb2_grpc.add_BrandServicer_to_server(BrandServicer(),server)
    goods_pb2_grpc.add_BannerServicer_to_server(BannerServicer(),server)
    goods_pb2_grpc.add_CategoryBrandServicer_to_server(CategoryBrandServicer(),server)

    #2.注册健康检查consul
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(),server)
    server.add_insecure_port(f"{args.ip}:{args.port}")
    service_id = str(uuid.uuid4())
    #主进程退出信号监听
    '''
    window下支持的信号有限:
        SIGINT: Ctrl+C终端
        SIGTERM: kill命令
    '''
    #优雅退出服务,注销服务到consul
    signal.signal(signal.SIGINT,partial(on_exit,service_id=service_id))
    signal.signal(signal.SIGTERM,partial(on_exit,service_id=service_id))

    logger.info(f"启动服务:{args.ip}:{args.port}")
    server.start()
    #3.注册服务到consul
    logger.info("注册服务开始")
    register = consul.ConsulRegister(settings.CONSUL_HOST,settings.CONSUL_PORT)
    if not register.register(name=settings.SERVICE_NAME,id=service_id,address=args.ip,port=args.port,tags=settings.SERVICE_TAGS,check=None):
        logger.info("注册服务失败")
        sys.exit(1) #退出进程,状态码为1,表示失败
    else:
        logger.info("注册服务成功")
    
    server.wait_for_termination() #线程阻塞,等待服务端终止

if __name__ == "__main__":
    #日志信息
    # logger.debug("调试信息")
    # logger.info("普通信息")
    # logger.warning("警告信息")
    # logger.error("错误信息")
    # logger.critical("严重错误信息")
    # print(get_free_port())
    logging.basicConfig()
    #监听nacos配置修改
    settings.client.add_config_watcher(settings.NACOS["data_id"],settings.NACOS["group"],settings.update_config)
    server()
