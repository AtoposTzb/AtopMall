import time
from datetime import date

import grpc
from loguru import logger
from peewee import DoesNotExist
from passlib.hash import pbkdf2_sha256
from google.protobuf import empty_pb2

from user_srv.model.models import User
from user_srv.proto import user_pb2,user_pb2_grpc

class UserServicer(user_pb2_grpc.UserServicer):
    @logger.catch
    def convert_user_to_rsp(self,user):
        #将model.User对象转换为message对象
        user_info_rsp = user_pb2.UserInfoResponse() 
        
        user_info_rsp.id = user.id
        user_info_rsp.password = user.password
        user_info_rsp.mobile = user.mobile
        user_info_rsp.rolo = user.rolo

        if user.nick_name:
            user_info_rsp.nickName = user.nick_name
        if user.gender:
            user_info_rsp.gender = user.gender
        if user.birthday:
            user_info_rsp.birthDay = int(time.mktime(user.birthday.timetuple()))

        return user_info_rsp
    
    @logger.catch
    def GetUserList(self,request:user_pb2.PageInfo,context): #参数：request,context 类型：user_pb2.PageInfo,grpc.ServerContext
        #获取用户的列表
        rsp = user_pb2.UserListResponse()  #创建一个UserListResponse对象,用于存储用户列表，并将数据传递给客户端调用
        users = User.select()
        rsp.total = users.count()  #设置用户总数,将数据给客户端调用
        #分页
        start = 0 #起始索引 页码
        page = 1 #当前页码
        per_page_num = 10 #每页显示的用户数量
        if request.pageSize:
            per_page_num = request.pageSize
        if request.pageNum: 
            page = request.pageNum  # 更新当前页码
        start = per_page_num * (page - 1)  # 计算起始索引

        users = users.limit(per_page_num).offset(start) #分页查询,从start开始,每页显示per_page_num个用户,共rsp.total条,users查询数量为per_page_num

        for user in users:
            user_info_rsp = self.convert_user_to_rsp(user)
            rsp.data.append(user_info_rsp)  #将用户信息添加到用户列表中
        return rsp
    
    @logger.catch
    def GetUserById(self,request:user_pb2.IdRequest,context):#context是grpc.ServerContext对象,用于获取客户端调用的元数据
        #通过id查询用户
        try:
            user = User.get(User.id == request.id)
            return self.convert_user_to_rsp(user)
        except DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"user用户不存在,id={request.id}")
            return user_pb2.UserInfoResponse()
        
    @logger.catch
    def GetUserByMobile(self,request:user_pb2.MobileRequest,context):#context是grpc.ServerContext对象,用于获取客户端调用的元数据
        #通过mobile查询用户
        try:
            user = User.get(User.mobile == request.mobile)
            return self.convert_user_to_rsp(user)
        except DoesNotExist as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"user用户不存在,mobile={request.mobile}")
            return user_pb2.UserInfoResponse()
    
    @logger.catch
    def CreateUser(self,request:user_pb2.CreateUserInfo,context):
        #新建用户
        try:
            User.get(User.mobile == request.mobile)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details("用户已存在")
            return user_pb2.UserInfoResponse()
        except DoesNotExist:
            pass
        
        user = User()
        user.nick_name = request.nickName
        user.mobile = request.mobile
        user.password = pbkdf2_sha256.hash(request.password)
        user.save()
        return self.convert_user_to_rsp(user) #将model.User对象转换为message对象
    
    @logger.catch
    def UpdateUser(self,request:user_pb2.UpdateUserInfo,context):
        try:
            user = User.get(User.id ==request.id)

            user.nick_name = request.nickName
            user.gender = request.gender
            user.birthday = date.fromtimestamp(request.birthDay)
            user.save()
            return empty_pb2.Empty() #不关心返回值,返回空对象 详见proto文件

        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在")
            return empty_pb2.Empty()