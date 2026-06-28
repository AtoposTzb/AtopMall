import json
import os
import sys
import grpc
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)

from goods_srv.proto import goods_pb2,goods_pb2_grpc
from loguru import logger
from google.protobuf import empty_pb2
from goods_srv.model.models import Category
from peewee import DoesNotExist

# 分类服务
class CategoryServicer(goods_pb2_grpc.CategoryServicer):
    def category_model_to_dict(self, category):
        re = {}
        re["id"] = category.id
        re["name"] = category.name
        re["parentCategory"] = category.parent_category_id 
        re["level"] = category.level
        re["isTab"] = category.is_tab
        return re


    @logger.catch
    def GetAllCategorysList(self, request:empty_pb2.Empty, context):
        categorys = Category.select() #先获取所有的分类
        category_list_rsp = goods_pb2.CategoryListResponse() #这个消息体返回两种格式数据
        category_list_rsp.total = categorys.count()
        level1 = []
        level2 = []
        level3 = []
        for category in categorys:
            category_rsp = goods_pb2.CategoryInfoResponse() #这个是需要返回给客户端的分类对象
            
            category_rsp.id = category.id
            category_rsp.name = category.name
            if category.parent_category_id:
                category_rsp.parentCategory = category.parent_category_id 
            category_rsp.level = category.level
            category_rsp.isTab = category.is_tab

            category_list_rsp.data.append(category_rsp)
        
            if category.level == 1:
                level1.append(self.category_model_to_dict(category))
            elif category.level == 2:
                level2.append(self.category_model_to_dict(category))
            elif category.level == 3:
                level3.append(self.category_model_to_dict(category))

        #开始整理数据 JsonData格式，方便客户端解析
        '''
        商品分类分3个级别,格式处理
        [{
            "id":1,
            "name":"一级分类1",
            "sub_category":[{
                "id":2,
                "name":"二级分类1-1",
                "sub_category":[{
                    "id":3,
                    "name":"三级分类1-1-1",
                },{},{},...]
            },{},{},...]
        },{},{},{},...]
        
        '''
        for data3 in level3:
            for data2 in level2:
                if data3["parentCategory"] == data2["id"]:
                    if "sub_category" not in data2:
                        data2["sub_category"] = [data3] #如果二级分类下没有三级分类，直接添加三级分类
                    else:
                        data2["sub_category"].append(data3) #如果有三级分类，直接添加三级分类
        for data2 in level2:
            for data1 in level1:
                if data2["parentCategory"] == data1["id"]:
                    if "sub_category" not in data1:
                        data1["sub_category"] = [data2] #如果一级分类下没有二级分类，直接添加二级分类
                    else:
                        data1["sub_category"].append(data2) #如果有二级分类，直接添加二级分类
        #上面已将将一级分类下的二级分类和三级分类都添加到了一级分类下
        #返回数据
        category_list_rsp.jsonData = json.dumps(level1) #将一级分类转换为json字符串
        return category_list_rsp
    
    @logger.catch
    def GetSubCategory(self, request: goods_pb2.CategoryListRequest, context):
        category_list_rsp = goods_pb2.SubCategoryListResponse()

        try:
            category_info = Category.get(Category.id == request.id)
            category_list_rsp.info.id = category_info.id
            category_list_rsp.info.name = category_info.name
            category_list_rsp.info.level = category_info.level
            category_list_rsp.info.isTab = category_info.is_tab
            if category_info.parent_category:
                category_list_rsp.info.parentCategory = category_info.parent_category_id
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return goods_pb2.SubCategoryListResponse()

        categorys = Category.select().where(Category.parent_category == request.id)
        category_list_rsp.total = categorys.count()
        for category in categorys:
            category_rsp = goods_pb2.CategoryInfoResponse()
            category_rsp.id = category.id
            category_rsp.name = category.name
            if category_info.parent_category:
                category_rsp.parentCategory = category_info.parent_category_id
            category_rsp.level = category.level
            category_rsp.isTab = category.is_tab

            category_list_rsp.subCategorys.append(category_rsp)

        return category_list_rsp

    @logger.catch
    def CreateCategory(self, request: goods_pb2.CategoryInfoRequest, context):
        try:
            category = Category()
            category.name = request.name
            if request.level != 1:
                category.parent_category = request.parentCategory
            category.level = request.level
            category.is_tab = request.isTab
            category.save()

            category_rsp = goods_pb2.CategoryInfoResponse()
            category_rsp.id = category.id
            category_rsp.name = category.name
            if category.parent_category:
                category_rsp.parentCategory = category.parent_category.id
            category_rsp.level = category.level
            category_rsp.isTab = category.is_tab
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details('插入数据失败：'+str(e))
            return goods_pb2.CategoryInfoResponse()

        return category_rsp

    @logger.catch
    def DeleteCategory(self, request: goods_pb2.DeleteCategoryRequest, context):
        try:
            category = Category.get(request.id)
            category.delete_instance()
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()

    @logger.catch
    def UpdateCategory(self, request: goods_pb2.CategoryInfoRequest, context):
        try:
            category = Category.get(request.id)
            if request.name:
                category.name = request.name
            if request.parentCategory:
                category.parent_category = request.parentCategory
            if request.level:
                category.level = request.level
            if request.isTab:
                category.is_tab = request.isTab
            category.save()

            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()
