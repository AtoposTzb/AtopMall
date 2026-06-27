import os
import sys
from peewee import *
from datetime import datetime

from passlib.hash import pbkdf2_sha256
#获取当前文件的上一级目录,即user_srv目录，方便vscode|traeIDE终端运行，其实就是找包的路径，不然会报“module not found”错误
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from settings.settings import DB
# print(BASE_DIR)

#删除：物理删除和逻辑删除(软删除)
#通过save方法做了修改如何确保update_time字段更新时间更新而不是add_time字段
class BaseModel(Model):
    add_time = DateTimeField(default=datetime.now,verbose_name="添加时间")
    update_time = DateTimeField(default=datetime.now,verbose_name="更新时间")
    is_delete = BooleanField(default=False,verbose_name="是否删除")

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
            return super().update(is_delete=True) #表示逻辑删除，将is_delete字段设置为True
        
    def delete_instance(self, permanently=False,recursive=False, delete_nullable=False):
        if permanently:
            return self.delete(permanently).where(self._pk_expr()).execute() #表示物理删除，将数据从数据库中删除
        else:
            self.is_delete = True
            return self.save() #表示逻辑删除，将is_delete字段设置为True
    #拦截查询操作，只查询未删除的数据
    @classmethod
    def select(cls, *fields):
        return super().select(*fields).where(cls.is_delete == False)
    #数据库操作类，指定数据库为DB
    class Meta:
        database = DB

class User(BaseModel):
    #用户模型
    GENDER_CHOICES = {
        ("female","女"),
        ("male","男"),
    }  
    ROLE_CHOICES = {
        (1,"普通用户"),
        (2,"管理员"),
    } 
    mobile = CharField(max_length=11,index=True,unique=True,verbose_name="手机号码")
    password = CharField(max_length=100,verbose_name="密码") #1.密文，2.密文不可反解
    nick_name = CharField(max_length=20,null=True,verbose_name="昵称")
    email = CharField(max_length=50,null=True,verbose_name="邮箱")
    gender = CharField(max_length=6,choices=GENDER_CHOICES,null=True,verbose_name="性别")
    role = IntegerField(default=1,choices=ROLE_CHOICES,verbose_name="角色")
    head_url = CharField(max_length=255,null=True,verbose_name="头像")
    birthday = DateField(null=True,verbose_name="生日")
    address = CharField(max_length=200,null=True,verbose_name="地址")
    desc = TextField(null=True,verbose_name="个人简介")

if __name__ == "__main__":
    DB.create_tables([User])
    #密码：1.对称加密，2.非对称加密 用户无法知道原始密码

    for i in range(10):
        user = User()
        user.nick_name = f"yolo{i}"
        user.mobile=f"134000000{i}"
        user.email="3488447218@qq.com"
        user.password = pbkdf2_sha256.hash("123456")
        user.save()
    #验证
    for user in User.select():
        print(pbkdf2_sha256.verify("123456",user.password))
