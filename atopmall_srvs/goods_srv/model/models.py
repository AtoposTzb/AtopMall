from datetime import datetime
from peewee import *
from playhouse.mysql_ext import JSONField
from settings.settings import DB

#删除：物理删除和逻辑删除(软删除)
#通过save方法做了修改如何确保update_time字段更新时间更新而不是add_time字段
class BaseModel(Model):
    add_time = DateTimeField(default=datetime.now,verbose_name="添加时间")
    update_time = DateTimeField(default=datetime.now,verbose_name="更新时间")
    is_deleted = BooleanField(default=False,verbose_name="是否删除")

    #重写save方法，确保update_time字段更新时间更新而不是add_time字段 拦截添加操作
    def save(self, *args, **kwargs): #*args, **kwargs是可变参数，用于接收任意数量的参数和关键字参数
        if self._pk is not None: #_pk是主键，如果主键存在，说明是更新操作，需要更新update_time字段 主键不一定是id
            self.update_time = datetime.now()
        else:
            self.add_time = datetime.now()
        return super().save(*args, **kwargs)

    #以下两个方法是删除方法，一个类方法，一个实例方法都能拦截删除操作
    @classmethod
    def delete(cls,permanently=False): #permanently表示是否永久删除
        if permanently:
            return super().delete() #表示物理删除，将数据从数据库中删除
        else:
            return super().update(is_deleted=True) #表示逻辑删除，将is_deleted字段设置为True
        
    def delete_instance(self, permanently=False,recursive=False, delete_nullable=False):
        if permanently:
            return self.delete(permanently).where(self._pk_expr()).execute() #表示物理删除，将数据从数据库中删除
        else:
            self.is_deleted = True
            return self.save() #表示逻辑删除，将is_deleted字段设置为True
    #拦截查询操作，只查询未删除的数据
    @classmethod
    def select(cls, *fields):
        return super().select(*fields).where(cls.is_deleted == False)

    class Meta:
        database = DB

# 商品分类表
class Category(BaseModel):
    name = CharField(max_length=20,verbose_name="名称")
    parent_category = ForeignKeyField("self",verbose_name="父分类",null=True) #一级类别可以没有父类别
    level = IntegerField(default=1,verbose_name="级别")
    is_tab = BooleanField(default=False,verbose_name="是否显示在首页tab")

class Brands(BaseModel):
    #品牌
    name = CharField(max_length=50, verbose_name="名称", index=True, unique=True)
    logo = CharField(max_length=200, null=True, verbose_name="图标", default="")


class Goods(BaseModel):
    """
    商品， 分布式的事务最好的解决方案 就是不要让分布式事务出现
    """
    category = ForeignKeyField(Category, verbose_name="商品类目", on_delete='CASCADE')
    brand = ForeignKeyField(Brands, verbose_name="品牌", on_delete='CASCADE')
    on_sale = BooleanField(default=True, verbose_name="是否上架")
    goods_sn = CharField(max_length=50, default="", verbose_name="商品唯一货号")
    name = CharField(max_length=100, verbose_name="商品名")
    click_num = IntegerField(default=0, verbose_name="点击数")
    sold_num = IntegerField(default=0, verbose_name="商品销售量")
    fav_num = IntegerField(default=0, verbose_name="收藏数") #库存是电商中一个重要的环节
    market_price = FloatField(default=0, verbose_name="市场价格")
    shop_price = FloatField(default=0, verbose_name="本店价格")
    goods_brief = CharField(max_length=200, verbose_name="商品简短描述")
    ship_free = BooleanField(default=True, verbose_name="是否承担运费")
    images = JSONField(verbose_name="商品轮播图")
    desc_images = JSONField(verbose_name="详情页图片")
    goods_front_image = CharField(max_length=200, verbose_name="封面图")
    is_new = BooleanField(default=False, verbose_name="是否新品")
    is_hot = BooleanField(default=False, verbose_name="是否热销")


class GoodsCategoryBrand(BaseModel):
    #品牌分类
    id = AutoField(primary_key=True, verbose_name="id")
    category = ForeignKeyField(Category, verbose_name="类别")
    brand = ForeignKeyField(Brands, verbose_name="品牌")

    class Meta:
        indexes = (
            #联合主键 约束category和brand的组合唯一 取消默认的联合主键索引
            (("category", "brand"), True),
        )


class Banner(BaseModel):
    """
    轮播的商品
    """
    image = CharField(max_length=200, default="", verbose_name="图片url")
    url = CharField(max_length=200, default="", verbose_name="访问url")
    index = IntegerField(default=0, verbose_name="轮播顺序")

if __name__ == "__main__":
    DB.create_tables([Category,Goods,Brands,GoodsCategoryBrand,Banner]) #创建商品分类表
    # c1 = Category(name="yolo1",level=1)
    # c1.save()
    # c2 = Category(name="yolo2",level=2)
    # c2.save()
    # c1 = Category.get(name="yolo1")
    # print(c1)
    # c1.delete_instance()
    # c1 = Category.get(name="yolo1")
    # print(c1)
    # for c in Category.select():
    #     print(c.id,c.name)
    # Category.delete().where(Category.id ==2).execute() #此时的执行会被拦截而进行逻辑删除