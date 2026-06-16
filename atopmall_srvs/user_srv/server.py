import grpc
import logging
from concurrent import futures
from user_srv.proto import user_pb2_grpc
from user_srv.handler.user import UserServicer

def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServicer_to_server(UserServicer(),server) #这行代码的意思是将UserServicer类添加到server中,并将其作为UserServicer服务端，简单来说就是注册UserServicer类
    server.add_insecure_port("[::]:50051")
    print("启动服务:127.0.0.1:50051")
    server.start()
    server.wait_for_termination() #线程阻塞,等待服务端终止

if __name__ == "__main__":
    logging.basicConfig()
    server()
