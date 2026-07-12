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
    "namespace":"3af4ccb8-3f1a-4532-8b52-8b0d33addc9c",
    "data_id":"inventory-srv.json",
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

DB = ReconnectMysqlDatebase(database=data["mysql"]["db"],host=data["mysql"]["host"],port=data["mysql"]["port"],user=data["mysql"]["user"],password=data["mysql"]["password"])
#consul配置
CONSUL_HOST = data["consul"]["host"]
CONSUL_PORT = data["consul"]["port"]

#服务配置
SERVICE_NAME = data["name"]
SERVICE_ID = data["name"]
SERVICE_TAGS = data["tags"]
