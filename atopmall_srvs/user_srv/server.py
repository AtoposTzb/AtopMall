import sys
import os
import logging
import signal
from concurrent import futures

import grpc
from loguru import logger

BASE_DIR = os.path.dirname(os.path.abspath(os.path.dirname(__file__)))
sys.path.insert(0,BASE_DIR)

from user_srv.proto import user_pb2_grpc
from user_srv.handler.user import UserServicer

def on_exit(sig,frame):
    logger.info("进程终端")
    sys.exit(0)

def server():
    logger.add("logs/user_srv_{time}.log") #将日志写入到文件夹logs下
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServicer_to_server(UserServicer(),server) #这行代码的意思是将UserServicer类添加到server中,并将其作为UserServicer服务端，简单来说就是注册UserServicer类
    server.add_insecure_port("[::]:50051")

    #主进程退出信号监听
    '''
    window下支持的信号有限:
        SIGINT: Ctrl+C终端
        SIGTERM: kill命令
    '''
    signal.signal(signal.SIGINT,on_exit)
    signal.signal(signal.SIGTERM,on_exit)

    logger.info("启动服务:127.0.0.1:50051")
    server.start()
    server.wait_for_termination() #线程阻塞,等待服务端终止

if __name__ == "__main__":
    #日志信息
    # logger.debug("调试信息")
    # logger.info("普通信息")
    # logger.warning("警告信息")
    # logger.error("错误信息")
    # logger.critical("严重错误信息")
    logging.basicConfig()
    server()
