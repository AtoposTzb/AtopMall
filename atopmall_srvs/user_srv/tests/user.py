import sys
import os
import grpc

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0,BASE_DIR)
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
    
    def user_by_id(self,id):
        #根据用户id获取用户信息
        rsp:user_pb2.UserInfoResponse = self.stub.GetUserById(user_pb2.IdRequest(id = id))
        print(rsp.id,rsp.mobile)
    def user_by_mobile(self,mobile):
        #根据用户mobile获取用户信息
        rsp:user_pb2.UserInfoResponse = self.stub.GetUserByMobile(user_pb2.MobileRequest(mobile = mobile))
        print(rsp.id,rsp.mobile)
    def create_user(self,nick_name,mobile,password):
        #创建用户
        rsp:user_pb2.UserInfoResponse = self.stub.CreateUser(user_pb2.CreateUserInfo(
            nickName = nick_name,
            password = password,
            mobile = mobile
        ))
        print(rsp.id)        
        
if __name__ == "__main__":
    testUser = TestUserServicer()
    # testUser.user_list()
    # testUser.user_by_id(3)
    # testUser.user_by_mobile("13800000000")
    testUser.create_user("test_user","13800000000","123456")