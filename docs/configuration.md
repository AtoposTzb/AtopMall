# 配置说明

## 本地配置文件

项目使用 Viper 管理本地配置，支持多环境：

| 文件                    | 用途                  | 是否提交 Git |
| ----------------------- | --------------------- | ------------ |
| config-debug_templ.yaml | 调试配置模板用于nacos | ✅ 是         |
| config-debug.yaml       | nacos连接调试配置     | ✅ 是         |
| config-pro.yaml         | nacos连接生产配置     | ✅ 是         |

## Nacos 配置中心

业务配置（MySQL、Redis、Consul、JWT、邮箱等）统一存放在 Nacos 配置中心：

| 服务          | Data ID            | Group | 说明                                        |
| ------------- | ------------------ | ----- | ------------------------------------------- |
| user_srv      | user-srv.json      | dev   | 用户微服务业务配置                          |
| user_web      | user-web.json      | dev   | 用户 Web 服务业务配置                       |
| goods_srv     | goods-srv.json     | dev   | 商品微服务业务配置                          |
| goods_web     | goods-web.json     | dev   | 商品 Web 服务业务配置                       |
| order_srv     | order-srv.json     | dev   | MySQL、Consul、RocketMQ、商品/库存服务名    |
| order_web     | order-web.json     | dev   | Consul、JWT、支付宝、订单/商品/库存服务地址 |
| inventory_srv | inventory-srv.json | dev   | MySQL、Redis、RocketMQ、Consul              |
| userop_srv    | userop-srv.json    | dev   | 用户操作微服务业务配置                      |
| userop_web    | userop-web.json    | dev   | 用户操作 Web 服务业务配置                   |
| oss_web       | oss-web.json       | dev   | 文件存储服务业务配置（MinIO、Consul）       |

## 配置加载流程

```
启动 → Viper 读取本地 config-debug.yaml（Nacos 连接信息）
     → 连接 Nacos 配置中心
     → 拉取业务配置（MySQL、Redis、Consul、RocketMQ、支付宝 等）
     → 解析配置到全局变量
     → 注册配置变更监听（配置修改后实时生效）
     → （Python 微服务）初始化 RocketMQ Producer/Consumer 客户端
```

## RocketMQ 配置

order_srv 和 inventory_srv 依赖 RocketMQ 进行分布式事务和异步解耦，配置从 Nacos 的 `rocketmq` 字段拉取。

**配置字段**（在 Nacos 的 order-srv.json / inventory-srv.json 中）：

```json
{
  "rocketmq": {
    "host": "192.168.1.106",
    "port": 9876
  }
}
```

| 字段            | 说明                     | 默认值     |
| --------------- | ------------------------ | ---------- |
| `rocketmq.host` | RocketMQ NameServer 地址 | 无（必填） |
| `rocketmq.port` | NameServer 端口          | 9876       |

**涉及 Topic 说明**：

| Topic           | 生产者    | 消费者        | 消息类型 | 用途                      |
| --------------- | --------- | ------------- | -------- | ------------------------- |
| `order_srv`     | order_srv | order_srv     | 事务消息 | 订单创建+库存扣减一致性   |
| `order_timeout` | order_srv | order_srv     | 延时消息 | 订单超时（1分钟）自动取消 |
| `order_reback`  | order_srv | inventory_srv | 普通消息 | 取消订单时通知库存归还    |

**启动依赖**：order_srv 和 inventory_srv 启动前需确保 RocketMQ NameServer 和 Broker 已启动。

## 支付宝配置

order_web 集成支付宝网页支付，配置从 Nacos 的 `alipay` 字段拉取。

**配置字段**（在 Nacos 的 order-web.json 中）：

```json
{
  "alipay": {
    "app_id": "your_app_id",
    "private_key": "your_private_key",
    "ali_public_key": "your_ali_public_key",
    "notify_url": "http://your-domain.com/o/v1/pay/alipay/notify",
    "return_url": "http://your-domain.com/o/v1/pay/alipay/return",
    "is_production": false,
    "product_code": "FAST_INSTANT_TRADE_PAY"
  }
}
```

| 字段                    | 说明                                                        |
| ----------------------- | ----------------------------------------------------------- |
| `alipay.app_id`         | 支付宝应用 ID（沙箱环境或正式环境）                         |
| `alipay.private_key`    | 应用私钥                                                    |
| `alipay.ali_public_key` | 支付宝公钥                                                  |
| `alipay.notify_url`     | 异步回调通知地址（支付宝服务端调用）                        |
| `alipay.return_url`     | 同步回调跳转地址（用户支付完成后跳转）                      |
| `alipay.is_production`  | 是否生产环境（`false` 为沙箱环境）                          |
| `alipay.product_code`   | 销售产品码，固定为 `FAST_INSTANT_TRADE_PAY`（电脑网站支付） |
