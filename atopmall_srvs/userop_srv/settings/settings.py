import json
from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin
from nacos import NacosClient
#不配置的话，会出现连接超时的问题
#使用peewee的连接池，使用ReconnectMixin防止连接超时断开导致查询失败

class ReconnectMysqlDatebase(ReconnectMixin,PooledMySQLDatabase):
    pass
# python的mro

# Nacos 配置
NACOS={
    "host":"192.168.1.106",
    "port":8848,
    "namespace":"3f453a3a-56b0-4896-b2c3-94617d1721ab",
    "data_id":"userop-srv.json",
    "group":"dev",
    "user":"nacos",
    "password":"nacos"
}
#根据自己的Nacos配置,添加鉴权信息
client = NacosClient(
    server_addresses=f"{NACOS['host']}:{NACOS['port']}",
    namespace=NACOS["namespace"]
)
data = client.get_config(
    NACOS["data_id"],
    NACOS["group"]
)
data = json.loads(data)
#监听配置修改
def update_config(args):
    print("配置文件发生变化")
    print(args)

#consul配置
CONSUL_HOST = data["consul"]["host"]
CONSUL_PORT = data["consul"]["port"]

#服务配置
SERVICE_NAME = data["name"]
SERVICE_ID = data["name"]
SERVICE_TAGS = data["tags"]

#mysql配置
MYSQL_DB = data["mysql"]["db"]
MYSQL_HOST = data["mysql"]["host"]
MYSQL_PORT = data["mysql"]["port"]
MYSQL_USER = data["mysql"]["user"]
MYSQL_PASSWORD = data["mysql"]["password"]

#redis配置
# REDIS_HOST = data["redis"]["host"]
# REDIS_PORT = data["redis"]["port"]
# # REDIS_PASSWORD = data["redis"]["password"]
# REDIS_DB = data["redis"]["db"]

#数据库全局连接池(mysql和redis连接池)
DB = ReconnectMysqlDatebase(database=MYSQL_DB,host=MYSQL_HOST,port=MYSQL_PORT,user=MYSQL_USER,password=MYSQL_PASSWORD)
# pool = redis.ConnectionPool(host=REDIS_HOST,port=REDIS_PORT,db=REDIS_DB)
# REDIS_CLIENT = redis.StrictRedis(connection_pool=pool)