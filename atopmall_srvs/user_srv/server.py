import sys
import os
import logging
import signal
import argparse
from concurrent import futures

import grpc
from loguru import logger

#获取当前文件的上一级目录,即user_srv目录，方便vscode|traeIDE终端运行，其实就是找包的路径，不然会报“module not found”错误
BASE_DIR = os.path.dirname(os.path.abspath(os.path.dirname(__file__)))
sys.path.insert(0,BASE_DIR)

from user_srv.proto import user_pb2_grpc
from user_srv.handler.user import UserServicer

def on_exit(sig,frame):
    logger.info("进程中断")
    sys.exit(0)

def server():
    #解析命令行参数,这个的作用是通过--ip和--port参数来指定服务端的IP地址和端口号
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip',
                        nargs='?',
                        type=str,
                        default="127.0.0.1",
                        help="服务端IP地址"
                        )
    parser.add_argument('--port',
                        nargs='?',
                        type=str,
                        default="50051",
                        help="服务端端口号，默认50051"
                        )
    args = parser.parse_args()  #解析命令行参数,返回一个命名空间对象,包含所有解析后的参数
    # args.ip  获取--ip参数的值  # args.port  获取--port参数的值

    logger.add("logs/user_srv_{time}.log") #将日志写入到文件夹logs下
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10)) #创建一个grpc服务器,并指定最大线程数为10
    user_pb2_grpc.add_UserServicer_to_server(UserServicer(),server) #这行代码的意思是将UserServicer类添加到server中,并将其作为UserServicer服务端，简单来说就是注册UserServicer类
    server.add_insecure_port(f"{args.ip}:{args.port}")

    #主进程退出信号监听
    '''
    window下支持的信号有限:
        SIGINT: Ctrl+C终端
        SIGTERM: kill命令
    '''
    #优雅退出服务
    signal.signal(signal.SIGINT,on_exit)
    signal.signal(signal.SIGTERM,on_exit)

    logger.info(f"启动服务:{args.ip}:{args.port}")
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
