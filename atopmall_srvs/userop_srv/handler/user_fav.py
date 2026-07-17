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
    pass