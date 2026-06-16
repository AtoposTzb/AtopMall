import grpc
from user_srv.proto import user_pb2,user_pb2_grpc
class TestUserServicer:
    def __init__(self):
        #连接grpc服务器
        channel = grpc.insecure_channel("127.0.0.1:50051")
        self.stub = user_pb2_grpc.UserStub(channel) #创建一个UserStub对象,用于调用UserServicer服务端的方法
        
    def user_list(self):
        #获取用户列表
        rsp :user_pb2.UserListResponse = self.stub.GetUserList(user_pb2.PageInfo(pageNum=2,pageSize=2))
        print(rsp.total)
        for user in rsp.data:
            print(user.mobile,user.birthDay)

if __name__ == "__main__":
    testuser = TestUserServicer()
    testuser.user_list()