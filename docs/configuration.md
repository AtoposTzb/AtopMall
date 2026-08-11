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
     → 拉取业务配置（MySQL、Redis、Consul、RocketMQ、Jaeger、支付宝 等）
     → 解析配置到全局变量
     → 注册配置变更监听（配置修改后实时生效）
     → （Go Web 服务）初始化 Jaeger TracerProvider
     → （Python 微服务）初始化 Jaeger Tracer（init_tracer）
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

## Jaeger 配置

所有服务（Go Web 和 Python 微服务）均支持 OpenTelemetry + Jaeger 全链路追踪，配置从 Nacos 的 `jaeger` 字段拉取。

**配置字段**（在 Nacos 的各服务配置文件中）：

```json
{
  "jaeger": {
    "host": "192.168.1.106",
    "port": 4317,
    "name": "user_web"
  }
}
```

| 字段          | 说明                               | 默认值     |
| ------------- | ---------------------------------- | ---------- |
| `jaeger.host` | Jaeger OTLP gRPC 地址              | 无（必填） |
| `jaeger.port` | OTLP gRPC 端口（Jaeger 默认 4317） | 4317       |
| `jaeger.name` | 服务名称，在 Jaeger UI 中显示      | 无（必填） |

**实现机制**：

| 层级               | 实现方式                                                                                             |
| ------------------ | ---------------------------------------------------------------------------------------------------- |
| Go Web 服务（Gin） | `initialize/jaeger_trace.go` 初始化 TracerProvider + Propagator，`router.go` 挂载 `otelgin` 中间件   |
| Go gRPC 客户端     | `src_conn.go` 中所有 gRPC 连接添加 `otelgrpc.NewClientHandler()`，API 层透传 `ctx.Request.Context()` |
| Python 微服务      | `common/jaeger_trace/trace.py` 统一 `init_tracer()` 函数，`server.py` 启动时调用                     |
| Python gRPC 服务端 | `grpc.aio` 框架自动集成 OpenTelemetry 拦截器，过滤健康检查 Span                                      |

**健康检查 Span 过滤**：Consul 每 5 秒触发一次 gRPC 健康检查，会产生大量无意义的 trace 数据。Python 端通过 `_FilterHealthExporter` 过滤 `/grpc.health.v1.Health/Check` 的 Span；Go 端将健康检查路由放在 `otelgin` 中间件之前，避免产生 Span。

**启动依赖**：Jaeger 为可选组件，不启动不影响业务功能。启动后所有服务的请求链路自动上报。

## Kong API 网关配置

Kong 作为统一入口网关，负责路由分发和 JWT 认证，前端请求通过 Kong 转发至各 Web 服务。

### 路由映射

| Kong Path | 目标服务   | 端口  | strip_path | 说明                              |
| --------- | ---------- | ----- | ---------- | --------------------------------- |
| `/g`      | goods_web  | 18082 | true       | 商品服务（商品/分类/品牌/轮播图） |
| `/u`      | user_web   | 18081 | true       | 用户服务（登录/注册/验证码）      |
| `/o`      | order_web  | 18084 | true       | 订单服务（订单/购物车/支付）      |
| `/op`     | userop_web | 18085 | true       | 用户操作服务（留言/收藏/地址）    |
| `/oss`    | oss_web    | 18083 | true       | 文件存储服务（MinIO 预签名直传）  |

> 路由前缀统一为 `/v1`，Kong 通过 `strip_path=true` 剥离 Path 前缀后转发。例如：`/g/v1/goods` → Kong 剥离 `/g` → 转发到 goods_web 的 `/v1/goods`

### JWT 认证

Kong 的 JWT 插件会自动在 Token 前添加 `Bearer ` 前缀，但直连微服务时前端可能不添加前缀。代码中已通过 `strings.HasPrefix` 兼容两种场景：

```go
// middlewares/jwt.go
if strings.HasPrefix(token, "Bearer ") {
    token = strings.Split(token, " ")[1]
}
```

**配置要点**：
- 在 Konga 中为各 Service 添加 JWT 插件
- 为 Consumer 创建 JWT Credential（key + secret）
- 前端请求统一在 Header 中携带 `Authorization: Bearer <token>`

## Sentinel 限流熔断配置

所有 Web 服务在启动时通过 `initialize/sentinel.go` 加载限流和熔断规则，使用哨兵 Sentinel 保护服务稳定性。

### 配置入口

每个 Web 服务在 `main.go` 中初始化 Sentinel：

```go
// goods_web/main.go
initialize.InitSentinel()
```

### 限流规则（Flow Rules）

所有规则定义在 `initialize/sentinel.go` 中，策略为 QPS 限流（基于滑动窗口）。

| 服务       | 规则数 | 阈值策略                                    |
| ---------- | ------ | ------------------------------------------- |
| goods_web  | 25 条  | 查询类 10次/6s，操作类 3次/6s               |
| order_web  | 8 条   | 查询类 10次/6s，操作类 3-5次/6s             |
| user_web   | 7 条   | 验证码 1次/60s，登录 3次/6s，其余 5-10次/6s |
| userop_web | 10 条  | 查询类 10次/6s，操作类 3-5次/6s             |
| oss_web    | 3 条   | 上传 Token 5次/6s，清理 1次/60s             |

**规则示例**：

```go
{
    Resource:          "goods-list",
    TokenCalculateStrategy: flow.Direct,
    ControlBehavior:        flow.Reject,
    Threshold:         10,
    StatIntervalInMs:   6000,
}
```

### 熔断规则（Circuit Breaker Rules）

与限流规则共用同一资源名，使用错误比例策略。

| 参数             | 值           | 说明                            |
| ---------------- | ------------ | ------------------------------- |
| Strategy         | `ErrorRatio` | 错误比例触发                    |
| Threshold        | 0.5          | 错误率超过 50% 打开熔断         |
| MinRequestAmount | 5            | 统计窗口内至少 5 个请求才触发   |
| StatIntervalMs   | 10000        | 10 秒滑动窗口统计               |
| RetryTimeoutMs   | 5000         | 熔断后 5 秒进入半开状态尝试恢复 |

**规则示例**：

```go
{
    Resource:         "goods-list",
    Strategy:         circuitbreaker.ErrorRatio,
    Threshold:        0.5,
    MinRequestAmount: 5,
    StatIntervalMs:   10000,
    RetryTimeoutMs:   5000,
}
```

### 代码实现

**限流入口**（在 handler 函数中）：

```go
e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
if b != nil {
    ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "请求频率过快,请稍后重试"})
    return
}
defer e.Exit()
```

**熔断错误追踪**（在所有 gRPC 调用错误处）：

```go
rsp, err := global.GoodsSrvCli.Goods.GoodsList(ctx.Request.Context(), req)
if err != nil {
    sentinel.TraceError(e, err) // 追踪错误，累计错误率
    api.HandleGrpcErrorToHttpError(err, ctx)
    return
}
```

### 熔断效果

当下游 gRPC 服务（Srv 层）故障时：

```
正常情况：请求 → 限流检查 → gRPC 调用 → 返回结果
故障情况：请求 → 限流检查 → 熔断已打开 → 立即返回 429（不等待 gRPC 超时）
恢复探测：5s 后 → 半开状态 → 放行少量请求 → 成功则恢复，失败则继续熔断
```

**保护效果**：防止下游服务故障导致 Web 层连接池耗尽、线程堆积和雪崩效应。熔断打开后，请求在 ms 级别快速失败，不会因为 gRPC 超时等待数秒。