# 配置说明

## 本地配置文件

项目使用 Viper 管理本地配置，支持多环境：

| 文件                    | 用途                  | 是否提交 Git |
| ----------------------- | --------------------- | ------------ |
| config-debug_templ.yaml | 调试配置模板用于nacos | ✅ 是        |
| config-debug.yaml       | nacos连接调试配置     | ✅ 是        |
| config-pro.yaml         | nacos连接生产配置     | ✅ 是        |

## Nacos 配置中心

业务配置（MySQL、Redis、Consul、JWT、邮箱等）统一存放在 Nacos 配置中心：

| 服务          | Data ID            | Group | 说明                                    |
| ------------- | ------------------ | ----- | --------------------------------------- |
| user_srv      | user-srv.json      | dev   | 用户微服务业务配置                      |
| user_web      | user-web.json      | dev   | 用户 Web 服务业务配置                   |
| goods_srv     | goods-srv.json     | dev   | 商品微服务业务配置                      |
| goods_web     | goods-web.json     | dev   | 商品 Web 服务业务配置                   |
| oss-web       | oss-web.json       | dev   | 文件存储服务业务配置（MinIO、Consul）   |
| order_srv     | order-srv.json     | dev   | 订单微服务业务配置（MySQL、Consul、商品/库存服务名） |
| inventory_srv | inventory-srv.json | dev   | 库存微服务业务配置（MySQL、Redis、Consul） |

## 配置加载流程

```
启动 → Viper 读取本地 config-debug.yaml（Nacos 连接信息）
     → 连接 Nacos 配置中心
     → 拉取业务配置（MySQL、Redis、Consul 等）
     → 解析配置到全局变量
     → 注册配置变更监听（配置修改后实时生效）
```
