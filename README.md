# AtopMall 电商微服务项目

基于 Go + Python 双语言混合开发的电商微服务项目。微服务层使用 Python + gRPC 实现业务逻辑，Web API 层使用 Go + Gin 对外提供 HTTP 接口，通过 gRPC 进行服务间通信，使用 Consul 实现服务注册与发现，使用 Nacos 作为配置中心统一管理配置。

## 一、项目结构

```
AtopMall/
├── atopmall_srvs/                    # 微服务层（Python + gRPC）
│   ├── common/                       # 公共模块
│   │   └── register/                 # 服务注册（Consul）
│   │       ├── base.py               # 注册接口抽象基类
│   │       └── consul.py             # Consul 注册实现
│   ── user_srv/                     # 用户服务
│       ├── handler/                  # gRPC 服务实现
│       ├── model/                    # Peewee ORM 数据模型
│       ├── proto/                    # Protobuf 定义及生成代码
│       ├── settings/                 # 配置管理
│       │   ├── settings.py           # Nacos 配置加载（同步版）
│       │   └── settings_asnyc.py     # Nacos 配置加载（异步版 + 配置变更监听）
│       ├── tests/                    # gRPC 客户端测试
│       ├── requirements.txt          # Python 依赖
│       └── server.py                 # gRPC 服务入口（含 Consul 注册 + 优雅退出）
│
├── atopmall_web/                     # Web API 层（Go + Gin）
│   └── user_web/                     # 用户 Web 服务
│       ├── api/                      # HTTP 接口实现
│       │   ├── captcha.go            # 图片验证码接口
│       │   ├── email_code.go         # 邮箱验证码接口
│       │   ├── user.go               # 用户登录/注册/列表接口
│       │   └── redis_test/           # Redis 连接测试工具
│       │       └── main.go
│       ├── config/                   # 配置结构体定义
│       ├── forms/                    # 请求表单验证
│       ├── global/                   # 全局变量（配置、Redis、翻译器、gRPC 客户端）
│       │   └── responselist/         # 响应结构体定义
│       ├── initialize/               # 初始化（配置加载、路由、日志、Redis、Consul 服务发现）
│       │   ├── config.go             # Viper + Nacos 配置加载
│       │   ├── logger.go             # Zap 日志初始化
│       │   ├── redis.go              # Redis 客户端初始化
│       │   ├── router.go             # Gin 路由注册
│       │   ├── src_conn.go           # Consul 服务发现 + gRPC 客户端初始化（含负载均衡）
│       │   └── validator_trans.go    # 表单验证器中文翻译
│       ├── middlewares/              # 中间件（JWT、CORS、权限）
│       ├── models/                   # 请求模型定义
│       ├── proto/                    # Protobuf 定义及生成代码
│       ├── router/                   # 路由分组
│       │   ├── base.go               # 基础路由（验证码相关）
│       │   └── user.go               # 用户路由（登录/注册/列表）
│       ├── utils/                    # 工具函数
│       │   └── addr_port.go          # 动态可用端口获取
│       ├── validator/                # 自定义验证器
│       ├── config-debug_templ.yaml   # 调试配置模板（更早的配置模板信息）
│       ├── config-debug.yaml         # nacos调试配置（仅含 Nacos 连接信息）
│       ├── config-pro.yaml           # nacos生产配置
│       └── main.go                   # 服务入口
│
└── README.md
```

## 二、技术栈总览

| 分类         | 技术选型                | 说明                                       |
| ------------ | ----------------------- | ------------------------------------------ |
| 开发语言     | Go 1.22+ / Python 3.13+ | Go 负责 Web API 层，Python 负责微服务层    |
| 微服务通信   | gRPC + Protobuf         | 服务间远程调用                             |
| 服务注册发现 | Consul                  | 微服务注册与健康检查                       |
| 配置中心     | Nacos                   | 统一配置管理，支持配置变更实时推送         |
| Web 框架     | Gin                     | Go HTTP 接口层开发                         |
| Python ORM   | Peewee                  | Python 数据库操作（含连接池 + 断线重连）   |
| Go ORM       | GORM（待集成）          | Go 数据库操作                              |
| Python 日志  | Loguru                  | Python 端日志组件                          |
| Go 日志      | Zap                     | Go 端高性能结构化日志                      |
| 配置管理     | Viper                   | YAML 配置文件加载与管理（本地 Nacos 连接） |
| 数据库       | MySQL                   | 数据存储                                   |
| 缓存         | Redis                   | 验证码存储、会话管理                       |
| JWT 认证     | golang-jwt/v5           | Token 生成与验证                           |
| 图片验证码   | base64Captcha           | 登录防暴力破解                             |
| 邮件服务     | jordan-wright/email     | SMTP 邮箱验证码发送                        |
| 表单验证     | go-playground/validator | 请求参数校验                               |

## 三、已完成功能

### 用户微服务（Python gRPC）

| 功能             | 接口            | 说明                                     |
| ---------------- | --------------- | ---------------------------------------- |
| 获取用户列表     | GetUserList     | 支持分页查询                             |
| 通过 id 查询用户 | GetUserById     | 验证用户是否存在                         |
| 通过 email 查询  | GetUserByEmail  | 验证用户是否存在                         |
| 根据 mobile 查询 | GetUserByMobile | 登录/注册时验证用户是否存在              |
| 创建用户         | CreateUser      | 注册新用户，密码使用 PBKDF2 加密         |
| 更新用户信息     | UpdateUser      | 检查用户是否存在后更新                   |
| 密码校验         | CheckPassWord   | 使用 passlib 验证 PBKDF2 哈希            |
| Consul 服务注册  | -               | 启动时自动注册到 Consul                  |
| gRPC 健康检查    | -               | 注册健康检查服务，Consul 定期探测        |
| 优雅退出         | -               | 监听 SIGINT/SIGTERM 信号，退出时注销服务 |
| Nacos 配置管理   | -               | 从 Nacos 拉取配置，支持配置变更监听      |
| 动态端口分配     | -               | 默认端口为 0 时自动获取可用端口          |

### 用户 Web 服务（Go Gin）

#### API 接口

| 功能       | 接口                 | 方法 | 说明                                           |
| ---------- | -------------------- | ---- | ---------------------------------------------- |
| 图片验证码 | /u/v1/base/captcha   | GET  | 生成 base64 格式验证码图片                     |
| 邮箱验证码 | /u/v1/base/send-code | POST | 发送注册验证码到邮箱，Redis 存储 5 分钟        |
| 密码登录   | /u/v1/user/pwd_login | POST | 手机号+密码+图片验证码登录，返回 JWT Token     |
| 用户注册   | /u/v1/user/register  | POST | 邮箱验证码+手机号+密码注册，注册成功返回 Token |
| 用户列表   | /u/v1/user/list      | GET  | 获取用户列表（需 JWT 认证 + 管理员权限）       |

#### 中间件

| 中间件           | 说明                            |
| ---------------- | ------------------------------- |
| JWT 认证中间件   | 解析 x-token 头部，验证登录状态 |
| 管理员权限中间件 | 验证用户角色是否为管理员        |
| CORS 中间件      | 跨域请求支持                    |

#### 工具模块

| 模块               | 说明                                     |
| ------------------ | ---------------------------------------- |
| 动态端口获取       | `utils/addr_port.go` 获取系统可用端口    |
| Redis 测试工具     | `api/redis_test/` 独立测试 Redis 连接    |
| 响应结构体         | `global/responselist/` 统一 API 响应格式 |
| 负载均衡 gRPC 连接 | `initialize/src_conn.go` 带负载均衡策略  |

## 四、API 路由结构

```
/u/v1/
├── base/                          # 基础服务（无需登录）
│   ├── GET  captcha               # 获取图片验证码
│   └── POST send-code             # 发送邮箱验证码
│
└── user/                          # 用户服务
    ├── POST pwd_login             # 密码登录（无需登录）
    ├── POST register              # 用户注册（无需登录）
    └── GET  list                  # 用户列表（需 JWT + 管理员）
```

## 五、开发工具清单

### Go 工具

| 工具名称           | 核心用途                 | 安装方式                                                          |
| ------------------ | ------------------------ | ----------------------------------------------------------------- |
| protoc             | Protocol Buffers 编译器  | 官网下载二进制包，配置环境变量                                    |
| protoc-gen-go      | protoc Go 代码生成插件   | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`  |
| protoc-gen-go-grpc | gRPC Go 代码生成插件     | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| air                | 代码热重载，保存自动重启 | `go install github.com/cosmtrek/air@latest`                       |

### Python 工具

| 工具名称         | 核心用途                    | 安装方式                             |
| ---------------- | --------------------------- | ------------------------------------ |
| grpcio-tools     | Protobuf 代码生成（Python） | `pip install grpcio-tools`           |
| python-consul    | Consul 客户端               | `pip install python-consul`          |
| nacos-sdk-python | Nacos 客户端（Python）      | `pip install nacos-sdk-python<3.0.0` |

### Proto 文件生成命令

**Go 端**（在 proto 文件所在目录下执行）：

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative xxx.proto
```

**Python 端**（在 proto 文件所在目录下执行）：

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. xxx.proto
```

## 六、快速开始

### 1. 环境准备

1. 安装 Go 1.22+ 并配置 GOPATH 环境变量
2. 安装 Python 3.10+ 并创建虚拟环境
3. 安装上表中所有开发工具
4. 本地启动 MySQL 数据库或者虚拟机安装 Docker 拉取 MySQL 镜像使用
5. 本地启动 Redis 服务
6. 本地启动 Consul 服务（服务注册与发现）
7. Dokcer启动 Nacos 服务（配置中心），并在 Nacos 中创建对应的配置

> 没有开发经验的可以参考我的有道云笔记:
> 【有道云笔记】项目前期准备
> https://share.note.youdao.com/s/QJFUWhau

### 2. Nacos 配置中心准备

在 Nacos 控制台中创建以下配置：

| 服务     | Data ID       | Group | 配置内容                              |
| -------- | ------------- | ----- | ------------------------------------- |
| user_srv | user-srv.json | dev   | MySQL、Consul、服务名称等配置         |
| user_web | user-web.json | dev   | MySQL、Redis、Consul、JWT、邮箱等配置 |

> user_web的nacos配置可参考 `config-debug_templ.yaml`文件

### 3. 启动用户微服务（Python gRPC）

```bash
cd atopmall_srvs/user_srv
pip install -r requirements.txt
python -m server
```

> 默认监听端口：50051，启动后自动注册到 Consul，配置从 Nacos 拉取

### 4. 启动用户 Web 服务（Go Gin）

```bash
cd atopmall_web/user_web
# 复制配置模板并修改（仅需配置 Nacos 连接信息）
cp config-debug_templ.yaml config-debug.yaml
go mod tidy
go run main.go
```

> 默认监听端口：8081，启动后从 Nacos 拉取业务配置，从 Consul 发现用户服务地址

## 七、配置说明

### 本地配置文件

项目使用 Viper 管理本地配置，支持多环境：

| 文件                    | 用途                                | 是否提交 Git |
| ----------------------- | ----------------------------------- | ------------ |
| config-debug_templ.yaml | 调试配置模板用于nacos | ✅ 是        |
| config-debug.yaml       | nacos连接调试配置 | ✅ 是    |
| config-pro.yaml         | nacos连接生产配置                 | ✅ 是        |

### Nacos 配置中心

业务配置（MySQL、Redis、Consul、JWT、邮箱等）统一存放在 Nacos 配置中心：

| 服务     | Data ID       | Group | 说明                  |
| -------- | ------------- | ----- | --------------------- |
| user_srv | user-srv.json | dev   | 用户微服务业务配置    |
| user_web | user-web.json | dev   | 用户 Web 服务业务配置 |

### 配置加载流程

```
启动 → Viper 读取本地 config-debug.yaml（Nacos 连接信息）
     → 连接 Nacos 配置中心
     → 拉取业务配置（MySQL、Redis、Consul 等）
     → 解析配置到全局变量
     → 注册配置变更监听（配置修改后实时生效）
```

## 八、服务注册与发现流程

![alt text](docs/image/consul注册服务简单图示.png)

1. **user_srv** 启动时通过 `python-consul` 注册到 Consul，包含 GRPC 健康检查
2. **user_web** 启动时从 Consul 查询 user_srv 的地址和端口
3. **user_web** 建立 gRPC 长连接（支持负载均衡策略），后续请求复用该连接
4. **user_srv** 异常退出时，Consul 自动注销该服务实例

## 九、用户注册流程

```
前端 → 获取图片验证码 → 填写注册信息（手机号、密码、邮箱）
     → 请求发送邮箱验证码 → 后端生成验证码存入 Redis（5分钟有效期）
     → 用户收到邮件，填写验证码 → 提交注册
     → 后端校验：手机号是否已存在 → 邮箱验证码是否正确
     → 调用 gRPC CreateUser 创建用户 → 生成 JWT Token 返回
```

## 十、配置中心架构图

```
┌──────────────────────────────────────────────────────────────────┐
│                        Nacos 配置中心                              │
│  ┌─────────────────┐  ┌─────────────────┐                        │
│  │  user-srv.json   │  │  user-web.json   │                        │
│  │  (Python 配置)   │  │  (Go 配置)       │                        │
│  │  - MySQL         │  │  - MySQL         │                        │
│  │  - Consul        │  │  - Redis         │                        │
│  │  - 服务名称       │  │  - Consul        │                        │
│  │                  │  │  - JWT           │                        │
│  │                  │  │  - 邮箱 SMTP     │                        │
│  └────────┬────────  └────────┬────────                        │
│           │ 配置推送            │ 配置推送                          │
└───────────┼────────────────────┼──────────────────────────────────┘
            │                    │
            ▼                    ▼
┌───────────────────┐  ┌────────────────────┐
│    user_srv       │  │    user_web         │
│  (Python gRPC)    │  │    (Go Gin)         │
│                   │  │                     │
│  settings.py      │  │  initialize/        │
│  ↓ Nacos 拉取配置  │  │  config.go          │
│  ↓ 配置变更监听    │  │  ↓ Nacos 拉取配置    │
│  ↓ 初始化 DB      │  │  ↓ 配置变更监听      │
│                   │  │  ↓ 初始化各组件      │
│  Consul 注册      │  │  Consul 发现        │
│  ↓ 注册服务       │  │  ↓ 获取服务地址      │
│  ↓ 健康检查       │  │  ↓ gRPC 长连接      │
│  ↓ 优雅退出       │  │  ↓ 负载均衡         │
└───────────────────┘  └────────────────────┘
```

## 十一、各服务 README

> 每个微服务将拥有独立的 README 文档，开发中...

| 服务                     | 语言   | 状态   |
| ------------------------ | ------ | ------ |
| user_srv（用户微服务）   | Python | 开发中 |
| user_web（用户 Web API） | Go     | 开发中 |
