from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin
#不配置的话，会出现连接超时的问题
#使用peewee的连接池，使用ReconnectMixin防止连接超时断开导致查询失败

class ReconnectMysqlDatebase(ReconnectMixin,PooledMySQLDatabase):
    pass
# python的mro

MYSQL_DB = "atopmall_user_srv"
MYSQL_HOST = "192.168.1.106"
MYSQL_PORT = 3306
MYSQL_USER = "root"
MYSQL_PASSWORD = "123456"

DB = ReconnectMysqlDatebase(MYSQL_DB,host=MYSQL_HOST,port=MYSQL_PORT,user=MYSQL_USER,password=MYSQL_PASSWORD)

#consul配置
CONSUL_HOST = "192.168.1.106"
CONSUL_PORT = 8500

#服务配置
SERVICE_NAME = "user_srv"
SERVICE_ID = "user_srv"
SERVICE_TAGS = ["atopmall","python","srv","yolo_t"]
