import os
import sys
from loguru import logger
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from google.protobuf import empty_pb2
from model.models import LeavingMessages
from proto import message_pb2,message_pb2_grpc

class MessageServicer(message_pb2_grpc.MessageServicer):
    pass