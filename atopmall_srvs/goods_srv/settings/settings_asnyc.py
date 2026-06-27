import json
import asyncio
import threading
from v2.nacos import NacosConfigService, ClientConfigBuilder, ConfigParam
from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin

# nacos-sdk-python 提供的V2版本，支持异步操作

# ============== Nacos 配置 ==============
NACOS = {
    "host": "192.168.1.106",
    "port": 8848,
    "namespace": "62270fc3-fd01-46fb-b27a-df43f2d99318",
    "data_id": "user-srv.json",
    "group": "dev",
    "user": "nacos",
    "password": "nacos"
}

# 客户端基础配置
client_config = (ClientConfigBuilder()
                 .server_address(f"{NACOS['host']}:{NACOS['port']}")
                 .namespace_id(NACOS["namespace"])
                 .build())


# ============== 全局共享：初始配置 + 就绪信号 ==============
_init_config_content = ""
_init_ready = threading.Event()


# ============== 配置变更回调 ==============
async def config_listener(config_info):
    """Nacos配置变更时自动触发"""
    print("\n========== 配置文件发生变化 ==========")
    print(f"最新配置内容:\n{config_info.content}")


# ============== 后台线程 = 拉取初始配置 + 注册监听 + 常驻 ==============
def _watch_daemon():
    """后台守护线程：一次性完成初始化+监听，全程复用一个客户端"""
    async def _watch_loop():
        global _init_config_content

        # 1. 创建配置客户端（只创建一次）
        config_client = await NacosConfigService.create_config_service(client_config)

        # 2. 拉取初始配置
        _init_config_content = await config_client.get_config(
            ConfigParam(data_id=NACOS["data_id"], group=NACOS["group"])
        )

        # 3. 通知主线程：初始配置已加载完成
        _init_ready.set()

        # 4. 注册配置变更监听器
        param = ConfigParam(data_id=NACOS["data_id"], group=NACOS["group"])
        await config_client.add_listener(param, config_listener)
        print("[Nacos] 配置监听已注册成功，等待变更...")

        # 5. 挂起协程，保持事件循环常驻，持续监听
        stop_event = asyncio.Event()
        await stop_event.wait()

    # 线程内创建并运行事件循环
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    loop.run_until_complete(_watch_loop())


# ============== 数据库连接池（保留原逻辑） ==============
class ReconnectMysqlDatebase(ReconnectMixin, PooledMySQLDatabase):
    pass


# ============== 模块初始化执行 ==============
# 启动后台监听线程
watch_thread = threading.Thread(target=_watch_daemon, daemon=True)
watch_thread.start()

# 阻塞等待初始配置加载完成（最多等10秒，超时抛错）
init_success = _init_ready.wait(timeout=10)
if not init_success:
    raise RuntimeError("Nacos 初始配置加载超时，请检查服务是否可用")

# 解析初始配置，初始化数据库
data = json.loads(_init_config_content)
DB = ReconnectMysqlDatebase(
    database=data["mysql"]["db"],
    host=data["mysql"]["host"],
    port=data["mysql"]["port"],
    user=data["mysql"]["user"],
    password=data["mysql"]["password"]
)

#consul配置
CONSUL_HOST = data["consul"]["host"]
CONSUL_PORT = data["consul"]["port"]

#服务配置
SERVICE_NAME = data["name"]
SERVICE_ID = data["name"]
SERVICE_TAGS = data["tags"]