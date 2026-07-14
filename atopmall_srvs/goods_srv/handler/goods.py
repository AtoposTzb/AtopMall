import grpc
from loguru import logger
from peewee import DoesNotExist
from google.protobuf import empty_pb2
from goods_srv.model.models import Goods,Category,Brands
from goods_srv.proto import goods_pb2,goods_pb2_grpc

# 商品服务
class GoodsServicer(goods_pb2_grpc.GoodsServicer):
        
    def convert_model_to_message(self, goods):
        info_rsp = goods_pb2.GoodsInfoResponse()

        info_rsp.id = goods.id
        info_rsp.categoryId = goods.category_id
        info_rsp.name = goods.name
        info_rsp.goodsSn = goods.goods_sn
        info_rsp.clickNum = goods.click_num
        info_rsp.soldNum = goods.sold_num
        info_rsp.favNum = goods.fav_num
        info_rsp.marketPrice = goods.market_price
        info_rsp.shopPrice = goods.shop_price
        info_rsp.goodsBrief = goods.goods_brief
        info_rsp.shipFree = goods.ship_free
        info_rsp.goodsFrontImage = goods.goods_front_image
        info_rsp.isNew = goods.is_new
        info_rsp.descImages.extend(goods.desc_images)
        info_rsp.images.extend(goods.images)      #使用的extend方法而不是append方法，是因为goods.images是一个列表，而info_rsp.images是一个列表，需要将goods.images中的元素添加到info_rsp.images中，合并两个列表
        info_rsp.isHot = goods.is_hot
        info_rsp.onSale = goods.on_sale

        info_rsp.category.id = goods.category.id
        info_rsp.category.name = goods.category.name

        info_rsp.brand.id = goods.brand.id
        info_rsp.brand.name = goods.brand.name
        info_rsp.brand.logo = goods.brand.logo

        return info_rsp

    @logger.catch # 商品列表页
    def GoodsList(self, request:goods_pb2.GoodsFilterRequest, context):
        #商品列表页
        rsp = goods_pb2.GoodsListResponse() #存储查询到的商品列表信息返回给客户端
        
        #过滤条件
        goods = Goods.select() #查询所有商品
        if request.keyWords: #一个简单的搜索
            goods = goods.filter(Goods.name.contains(request.keyWords))  #filter()方法用于过滤查询结果，返回符合条件的商品列表 where()方法用于指定查询条件，返回符合条件的商品列表
        if request.isHot:
            goods = goods.filter(Goods.is_hot == True)
        if request.isNew:
            goods = goods.filter(Goods.is_new == True)
        if request.priceMin:
            goods = goods.filter(Goods.shop_price >= request.priceMin)
        if request.priceMax:
            goods = goods.filter(Goods.shop_price <= request.priceMax)
        if request.brand:
            goods = goods.filter(Goods.brand_id == request.brand)
        #分类级别过滤 通过category查询商品，可能是一级分类，可能是二级分类，可能是三级分类
        if request.topCategory:
            try:
                ids=[]
                category = Category.get(Category.id == request.topCategory)
                level = category.level
                if level == 1: 
                    #SELECT * FROM category WHERE parent_category_id IN ( SELECT category.id FROM category WHERE parent_category_id=130358)
                    categorys = Category.select().where(Category.parent_category_id.in_(Category.select().where(Category.parent_category_id == request.topCategory)))
                    for category in categorys:
                        ids.append(category.id)
                elif level == 2:
                    categorys = Category.select().where(Category.parent_category_id == request.topCategory)
                    for category in categorys:
                        ids.append(category.id)
                elif level == 3:
                    ids.append(request.topCategory)
                goods = goods.filter(Goods.category_id.in_(ids))
            except DoesNotExist:
                logger.error(f"分类id={request.topCategory}不存在")
                rsp.total = 0
                return rsp
        #分页 limit offset
        start = 0
        page = 1
        pagePerNums = 10
        if request.pagePerNums:
            pagePerNums = request.pagePerNums
        if request.pages:
            page = request.pages
            start = (page-1)*pagePerNums

        rsp.total = goods.count()
        goods = goods.limit(pagePerNums).offset(start)
        for good in goods:
            rsp.data.append(self.convert_model_to_message(good))

        return rsp
    
    @logger.catch # 批量获取商品详情
    def BatchGetGoods(self,request:goods_pb2.BatchGoodsIdInfo,context):
        #批量获取商品详情,订单新建的时候可以使用
        rsp = goods_pb2.GoodsListResponse()
        goods = Goods.select().where(Goods.id.in_(list(request.id))) #根据商品id列表查询所有商品
        rsp.total = goods.count()
        for good in goods:
            rsp.data.append(self.convert_model_to_message(good))
        return rsp
    
    @logger.catch # 删除商品
    def DeleteGoods(self,request:goods_pb2.DeleteGoodsInfo,context):
        #删除商品
        # rows = Goods.delete().where(Goods.id==request.id)
        # if rows == 0:
        #     logger.error(f"商品id={request.id}不存在")
        #     return empty_pb2.Empty()
        try:
            goods =  Goods.get(Goods.id == request.id)
            goods.delete_instance() #删除商品
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"商品id={request.id}不存在")
            return empty_pb2.Empty()
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL) #严重的内部错误
            context.set_details(str(e))
            return empty_pb2.Empty()
    
    @logger.catch # 获取商品详情
    def GetGoodsDetail(self,request:goods_pb2.GoodInfoRequest,context):
        #获取商品的详情
        try:
            goods =  Goods.get(Goods.id == request.id)
            #每次查询商品详情，商品的点击量+1
            goods.click_num += 1
            goods.save() #保存商品信息
            rsp = self.convert_model_to_message(goods)
            return rsp
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"商品id={request.id}不存在")
            return goods_pb2.GoodsInfoResponse()
        
    @logger.catch # 创建商品
    def CreateGoods(self,request:goods_pb2.CreateGoodsInfo,context):
        #创建商品
        #先处理外键 商品分类和品牌分类
        try:
            category = Category.get(Category.id == request.categoryId)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("商品分类不存在")
            return goods_pb2.GoodsInfoResponse()
        try:
            brand = Brands.get(Brands.id == request.brandId)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("品牌不存在")
            return goods_pb2.GoodsInfoResponse()
        
        goods = Goods() #创建商品实例 将grpc消息转换为数据库模型存储到数据库
        goods.brand = brand
        goods.category = category
        goods.name = request.name
        goods.goods_sn = request.goodsSn
        goods.market_price = request.marketPrice
        goods.shop_price = request.shopPrice
        goods.goods_brief = request.goodsBrief
        goods.ship_free = request.shipFree
        goods.images = list(request.images)
        goods.desc_images = list(request.descImages)
        goods.goods_front_image = request.goodsFrontImage
        goods.is_new = request.isNew
        goods.is_hot = request.isHot
        goods.on_sale = request.onSale

        goods.save()

        #TODO 此处完善库存的设置 - 分布式事务
        return self.convert_model_to_message(goods)
    
    @logger.catch # 更新商品
    def UpdateGoods(self,request:goods_pb2.CreateGoodsInfo,context):
        #更新商品
        #先处理外键 商品分类和品牌分类和查询商品是否存在
        try:
            category = Category.get(Category.id == request.categoryId)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("商品分类不存在")
            return goods_pb2.GoodsInfoResponse()
       
        try:
            brand = Brands.get(Brands.id == request.brandId)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("品牌不存在")
            return goods_pb2.GoodsInfoResponse()

        try:
            goods = Goods.get(Goods.id == request.id)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("商品不存在")
            return goods_pb2.GoodsInfoResponse()

        goods.brand = brand
        goods.category = category
        goods.name = request.name
        goods.goods_sn = request.goodsSn
        goods.market_price = request.marketPrice
        goods.shop_price = request.shopPrice
        goods.goods_brief = request.goodsBrief
        goods.ship_free = request.shipFree
        goods.images = list(request.images)
        goods.desc_images = list(request.descImages)
        goods.goods_front_image = request.goodsFrontImage
        goods.is_new = request.isNew
        goods.is_hot = request.isHot
        goods.on_sale = request.onSale

        goods.save()

        #TODO 此处完善库存的设置 - 分布式事务
        return self.convert_model_to_message(goods)

    @logger.catch 
    def UpdateGoodsStatus(self, request: goods_pb2.GoodsStatusRequest, context):
        # 更新商品状态（仅更新 is_new/is_hot/on_sale）
        try:
            goods = Goods.get(Goods.id == request.id)
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("商品不存在")
            return empty_pb2.Empty()

        goods.is_new = request.isNew
        goods.is_hot = request.isHot
        goods.on_sale = request.onSale
        goods.save()

        return empty_pb2.Empty()