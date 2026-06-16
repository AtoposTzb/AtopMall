import time
import datetime
import grpc
from user_srv.model.models import User
from user_srv.proto import user_pb2,user_pb2_grpc

class UserServicer(user_pb2_grpc.UserServicer):
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
            user_info_rsp = user_pb2.UserInfoResponse() #创建一个UserInfoResponse对象,用于存储用户信息，并将数据传递给客户端调用
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

            rsp.data.append(user_info_rsp)  #将用户信息添加到用户列表中
        return rsp
    

        