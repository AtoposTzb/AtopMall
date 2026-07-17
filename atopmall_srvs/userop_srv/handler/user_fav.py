import os
import sys
import grpc
import datetime
from loguru import logger
from peewee import DoesNotExist
from google.protobuf import empty_pb2

BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from proto import userfav_pb2, userfav_pb2_grpc
from model.models import UserFav

class UserFavServicer(userfav_pb2_grpc.UserFavServicer):
    @logger.catch
    def GetFavList(self, request: userfav_pb2.UserFavRequest, context):
        # 获取收藏列表
        rsp = userfav_pb2.UserFavListResponse()
        user_favs = UserFav.select()
        if request.userId:
            user_favs = user_favs.where(UserFav.user==request.userId)
        if request.goodsId:
            user_favs = user_favs.where(UserFav.goods==request.goodsId)

        rsp.total = user_favs.count()
        for user_fav in user_favs:
            user_fav_rsp = userfav_pb2.UserFavResponse()
            user_fav_rsp.userId = user_fav.user
            user_fav_rsp.goodsId = user_fav.goods

            rsp.data.append(user_fav_rsp)

        return rsp

    @logger.catch
    def AddUserFav(self, request: userfav_pb2.UserFavRequest, context):
        # 添加收藏
        user_fav = UserFav()
        user_fav.user = request.userId
        user_fav.goods = request.goodsId
        user_fav.save(force_insert=True)

        return empty_pb2.Empty()

    @logger.catch
    def DeleteUserFav(self, request: userfav_pb2.UserFavRequest, context):
        # 删除收藏
        try:
            user_fav = UserFav.get(UserFav.user==request.userId, UserFav.goods==request.goodsId)
            user_fav.delete_instance(permanently=True)

            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()

    @logger.catch
    def GetUserFavDetail(self, request: userfav_pb2.UserFavRequest, context):
        # 获取收藏详情
        try:
            UserFav.get(UserFav.user == request.userId, UserFav.goods == request.goodsId)
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()