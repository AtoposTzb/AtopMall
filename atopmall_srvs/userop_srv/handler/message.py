import os
import sys
from loguru import logger
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from model.models import LeavingMessages
from proto import message_pb2,message_pb2_grpc

class MessageServicer(message_pb2_grpc.MessageServicer):
    @logger.catch
    def MessageList(self, request: message_pb2.MessageRequest, context):
        # 获取分类列表
        rsp = message_pb2.MessageListResponse()
        messages = LeavingMessages.select()
        if request.userId:
            messages = messages.where(LeavingMessages.user==request.userId)

        rsp.total = messages.count()
        for message in messages:
            brand_rsp = message_pb2.MessageResponse()

            brand_rsp.id = message.id
            brand_rsp.userId = message.user
            brand_rsp.messageType = message.message_type
            brand_rsp.subject = message.subject
            brand_rsp.message = message.message
            brand_rsp.file = message.file

            rsp.data.append(brand_rsp)

        return rsp

    @logger.catch
    def CreateMessage(self, request: message_pb2.MessageRequest, context):
        # 创建消息
        message = LeavingMessages()

        message.user = request.userId
        message.message_type = request.messageType
        message.subject = request.subject
        message.message = request.message
        message.file = request.file

        message.save()

        rsp = message_pb2.MessageResponse()
        rsp.id = message.id
        rsp.messageType = message.message_type
        rsp.subject = message.subject
        rsp.message = message.message
        rsp.file = message.file

        return rsp
