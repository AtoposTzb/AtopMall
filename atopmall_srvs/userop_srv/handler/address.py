import grpc
import os
import sys
from loguru import logger
from peewee import DoesNotExist
from google.protobuf import empty_pb2
BASE_DIR =  os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from proto import address_pb2, address_pb2_grpc
from model.models import Address

class AddressServicer(address_pb2_grpc.AddressServicer):
    pass