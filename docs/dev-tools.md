# 开发工具清单

## Go 工具

| 工具名称           | 核心用途                 | 安装方式                                                          |
| ------------------ | ------------------------ | ----------------------------------------------------------------- |
| protoc             | Protocol Buffers 编译器  | 官网下载二进制包，配置环境变量                                    |
| protoc-gen-go      | protoc Go 代码生成插件   | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`  |
| protoc-gen-go-grpc | gRPC Go 代码生成插件     | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| air                | 代码热重载，保存自动重启 | `go install github.com/cosmtrek/air@latest`                       |
| mc                 | MinIO 命令行客户端       | 官网下载二进制包或 `docker pull minio/mc`                         |

### Go 依赖包

| 包名                                                                           | 用途                      | 安装方式                                                                                     |
| ------------------------------------------------------------------------------ | ------------------------- | -------------------------------------------------------------------------------------------- |
| `github.com/smartwalle/alipay/v3`                                              | 支付宝 v3 SDK             | `go get github.com/smartwalle/alipay/v3@v3.2.30`                                             |
| `go.opentelemetry.io/otel`                                                     | OpenTelemetry SDK         | `go get go.opentelemetry.io/otel@latest`                                                     |
| `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin` | Gin OpenTelemetry 中间件  | `go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin@latest` |
| `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`  | gRPC OpenTelemetry 拦截器 | `go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc@latest`  |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`              | OTLP gRPC 导出器          | `go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@latest`              |
| `github.com/alibaba/sentinel-golang`                                           | Sentinel 限流熔断         | `go get github.com/alibaba/sentinel-golang@latest`                                           |

## Python 工具

| 工具名称                               | 核心用途                    | 安装方式                                             |
| -------------------------------------- | --------------------------- | ---------------------------------------------------- |
| grpcio-tools                           | Protobuf 代码生成（Python） | `pip install grpcio-tools`                           |
| python-consul                          | Consul 客户端               | `pip install python-consul`                          |
| nacos-sdk-python                       | Nacos 客户端（Python）      | `pip install nacos-sdk-python<3.0.0`                 |
| redis                                  | Redis 客户端                | `pip install redis`                                  |
| python-redis-lock                      | Redis 分布式锁              | `pip install python-redis-lock`                      |
| rocketmq-client-python                 | RocketMQ 消息队列客户端     | `pip install rocketmq-client-python`                 |
| opentelemetry-api                      | OpenTelemetry API（Python） | `pip install opentelemetry-api`                      |
| opentelemetry-sdk                      | OpenTelemetry SDK（Python） | `pip install opentelemetry-sdk`                      |
| opentelemetry-exporter-otlp-proto-grpc | OTLP gRPC 导出器（Python）  | `pip install opentelemetry-exporter-otlp-proto-grpc` |

## Proto 文件生成命令

### Go 端

在 proto 文件所在目录下执行：

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative xxx.proto
```

### Python 端

在 proto 文件所在目录下执行：

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. xxx.proto
```

## 一键启停脚本

项目提供一键启停脚本，使用 tmux 管理所有微服务进程，方便本地开发和调试。

| 脚本            | 平台    | 说明                                       |
| --------------- | ------- | ------------------------------------------ |
| `start-all.sh`  | Linux   | 使用 tmux 会话管理，一键启动所有微服务     |
| `stop-all.sh`   | Linux   | 销毁 tmux 会话 + 兜底清理残留进程          |
| `start-all.ps1` | Windows | 使用 PowerShell 多窗口，一键启动所有微服务 |
| `stop-all.ps1`  | Windows | 关闭所有微服务窗口及其子进程               |

**依赖**：Linux 需安装 `tmux`（终端复用器）。

```bash
# Linux
tmux --version       # 检查 tmux 是否安装
./start-all.sh       # 一键启动所有微服务
./stop-all.sh        # 一键停止所有微服务

# Windows PowerShell
.\start-all.ps1      # 一键启动所有微服务
.\stop-all.ps1       # 一键停止所有微服务
```

## 中间件 / 基础设施

项目运行依赖以下中间件，建议使用 Docker 部署：
Docker 部署的步骤和方法详见我的:【有道云笔记】0.前期准备https://share.note.youdao.com/s/RA68Zm9P

| 中间件   | 用途            | 默认端口        | Docker 启动示例                                                                                        |
| -------- | --------------- | --------------- | ------------------------------------------------------------------------------------------------------ |
| MySQL    | 数据持久化存储  | 3306            | `docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8.0`                        |
| Redis    | 缓存 / 分布式锁 | 6379            | `docker run -d --name redis -p 6379:6379 redis:latest`                                                 |
| Consul   | 服务注册与发现  | 8500            | `docker run -d --name consul -p 8500:8500 consul:latest`                                               |
| Nacos    | 配置中心        | 8848/9848       | `docker run -d --name nacos -p 8848:8848 -p 9848:9848 nacos/nacos-server:latest`                       |
| MinIO    | 对象存储        | 9000/9001       | `docker run -d --name minio -p 9000:9000 -p 9001:9001 minio/minio server /data`                        |
| Jaeger   | 分布式链路追踪  | 16686/4317/4318 | `docker run -d --name jaeger -p 16686:16686 -p 4317:4317 -p 4318:4318 jaegertracing/all-in-one:latest` |
| RocketMQ | 消息队列        | 9876            | 需启动 NameServer + Broker，详见 RocketMQ 官方文档                                                     |
| Kong     | API 网关        | 8000/8001/8443  | 统一入口、路由分发、JWT 认证，Konga 可视化面板管理                                                     |

> RocketMQ 需要同时启动 NameServer（默认 9876 端口）和 Broker，order_srv 和 inventory_srv 依赖 RocketMQ 进行分布式事务和异步消息处理。
> Kong 推荐使用 Docker 部署，通过 Konga（管理面板，默认端口 1337）可视化管理路由和 JWT 插件。Sentinel 限流熔断规则在代码中硬编码，无需外部控制台。