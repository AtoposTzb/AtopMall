import os
import sys
from peewee import *

from passlib.hash import pbkdf2_sha256
#获取当前文件的上一级目录,即user_srv目录，方便vscode|traeIDE终端运行，其实就是找包的路径，不然会报“module not found”错误
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0,BASE_DIR)
from settings.settings import DB
# print(BASE_DIR)

class BaseModel(Model):
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
