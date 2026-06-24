# AtopMall 电商微服务项目

基于 Go + Python 双语言混合开发的电商微服务项目。微服务层使用 Python + gRPC 实现业务逻辑，Web API 层使用 Go + Gin 对外提供 HTTP 接口，通过 gRPC 进行服务间通信，使用 Consul 实现服务注册与发现。

## 一、项目结构

```
AtopMall/
├── atopmall_srvs/                    # 微服务层（Python + gRPC）
│   ├── common/                       # 公共模块
│   │   └── register/                 # 服务注册（Consul）
│   │       ├── base.py               # 注册接口抽象基类
│   │       └── consul.py             # Consul 注册实现
│   └── user_srv/                     # 用户服务
│       ├── handler/                  # gRPC 服务实现
│       ├── model/                    # Peewee ORM 数据模型
│       ├── proto/                    # Protobuf 定义及生成代码
│       ├── settings/                 # 数据库配置
│       ├── tests/                    # gRPC 客户端测试
│       ├── requirements.txt          # Python 依赖
│       └── server.py                 # gRPC 服务入口（含 Consul 注册）
│
├── atopmall_web/                     # Web API 层（Go + Gin）
│   └── user_web/                     # 用户 Web 服务
│       ├── api/                      # HTTP 接口实现
│       │   ├── captcha.go            # 图片验证码接口
│       │   ├── email_code.go         # 邮箱验证码接口
│       │   └── user.go               # 用户登录/注册/列表接口
│       ├── config/                   # 配置结构体定义
│       ├── forms/                    # 请求表单验证
│       ├── global/                   # 全局变量（配置、Redis、翻译器、gRPC 客户端）
│       ├── initialize/               # 初始化（配置加载、路由、日志、Redis、Consul 服务发现）
│       │   ├── config.go             # Viper 配置加载
│       │   ├── logger.go             # Zap 日志初始化
│       │   ├── redis.go              # Redis 客户端初始化
│       │   ├── router.go             # Gin 路由注册
│       │   ├── src_conn.go           # Consul 服务发现 + gRPC 客户端初始化
│       │   └── validator_trans.go    # 表单验证器中文翻译
│       ├── middlewares/              # 中间件（JWT、CORS、权限）
│       ├── models/                   # 请求模型定义
│       ├── proto/                    # Protobuf 定义及生成代码
│       ├── router/                   # 路由分组
│       ├── validator/                # 自定义验证器
│       ├── config-debug_templ.yaml   # 调试配置模板（可提交）
│       ├── config-pro.yaml           # 生产配置模板
│       └── main.go                   # 服务入口
│
└── README.md
```

## 二、技术栈总览

| 分类         | 技术选型                | 说明                                    |
| ------------ | ----------------------- | --------------------------------------- |
| 开发语言     | Go 1.22+ / Python 3.13+ | Go 负责 Web API 层，Python 负责微服务层 |
| 微服务通信   | gRPC + Protobuf         | 服务间远程调用                          |
| 服务注册发现 | Consul                  | 微服务注册与健康检查                    |
| Web 框架     | Gin                     | Go HTTP 接口层开发                      |
| Python ORM   | Peewee                  | Python 数据库操作                       |
| Go ORM       | GORM（待集成）          | Go 数据库操作                           |
| Python 日志  | Loguru                  | Python 端日志组件                       |
| Go 日志      | Zap                     | Go 端高性能结构化日志                   |
| 配置管理     | Viper                   | YAML 配置文件加载与管理                 |
| 数据库       | MySQL                   | 数据存储                                |
| 缓存         | Redis                   | 验证码存储、会话管理                    |
| JWT 认证     | golang-jwt/v5           | Token 生成与验证                        |
| 图片验证码   | base64Captcha           | 登录防暴力破解                          |
| 邮件服务     | jordan-wright/email     | SMTP 邮箱验证码发送                     |
| 表单验证     | go-playground/validator | 请求参数校验                            |

## 三、已完成功能

### 用户微服务（Python gRPC）

| 功能             | 接口            | 说明                             |
| ---------------- | --------------- | -------------------------------- |
| 获取用户列表     | GetUserList     | 支持分页查询                     |
| 通过 id 查询用户 | GetUserById     | 验证用户是否存在                 |
| 通过 email 查询  | GetUserByEmail  | 验证用户是否存在                 |
| 根据 mobile 查询 | GetUserByMobile | 登录时验证用户是否存在           |
| 创建用户         | CreateUser      | 注册新用户，密码使用 PBKDF2 加密 |
| 更新用户信息     | UpdateUser      | 检查用户是否存在后更新           |
| 密码校验         | CheckPassWord   | 使用 passlib 验证 PBKDF2 哈希    |
| Consul 服务注册  | -               | 启动时自动注册到 Consul          |

### 用户 Web 服务（Go Gin）

| 功能             | 接口                | 说明                                               |
| ---------------- | ------------------- | -------------------------------------------------- |
| 图片验证码       | GET /captcha        | 生成 base64 格式验证码图片                         |
| 邮箱验证码       | POST /email/code    | 发送注册验证码到邮箱，Redis 存储 5 分钟            |
| 密码登录         | POST /user/login    | 手机号+密码+图片验证码登录，返回 JWT Token         |
| 用户注册         | POST /user/register | 邮箱验证码+手机号+密码注册                         |
| 用户列表         | GET /user/list      | 获取用户列表（需 JWT 认证 + 管理员权限）           |
| JWT 认证中间件   | -                   | 解析 x-token 头部，验证登录状态                    |
| 管理员权限中间件 | -                   | 验证用户角色是否为管理员                           |
| CORS 中间件      | -                   | 跨域请求支持                                       |
| Consul 服务发现  | -                   | 启动时从 Consul 获取用户服务地址，建立 gRPC 长连接 |

## 四、开发工具清单

### Go 工具

| 工具名称           | 核心用途                 | 安装方式                                                          |
| ------------------ | ------------------------ | ----------------------------------------------------------------- |
| protoc             | Protocol Buffers 编译器  | 官网下载二进制包，配置环境变量                                    |
| protoc-gen-go      | protoc Go 代码生成插件   | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`  |
| protoc-gen-go-grpc | gRPC Go 代码生成插件     | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| air                | 代码热重载，保存自动重启 | `go install github.com/cosmtrek/air@latest`                       |

### Python 工具

| 工具名称      | 核心用途                    | 安装方式                    |
| ------------- | --------------------------- | --------------------------- |
| grpcio-tools  | Protobuf 代码生成（Python） | `pip install grpcio-tools`  |
| python-consul | Consul 客户端               | `pip install python-consul` |

### Proto 文件生成命令

**Go 端**（在 proto 文件所在目录下执行）：

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative xxx.proto
```

**Python 端**（在 proto 文件所在目录下执行）：

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. xxx.proto
```

## 五、快速开始

### 1. 环境准备

1. 安装 Go 1.22+ 并配置 GOPATH 环境变量
2. 安装 Python 3.10+ 并创建虚拟环境
3. 安装上表中所有开发工具
4. 本地启动 MySQL 数据库或者虚拟机安装 Docker 拉取 MySQL 镜像使用
5. 本地启动 Redis 服务
6. 本地启动 Consul 服务（服务注册与发现）

> 没有开发经验的可以参考我的有道云笔记:
> 【有道云笔记】项目前期准备
> https://share.note.youdao.com/s/QJFUWhau

### 2. 启动用户微服务（Python gRPC）

```bash
cd atopmall_srvs/user_srv
pip install -r requirements.txt
python -m server
```

> 默认监听端口：50051，启动后自动注册到 Consul

### 3. 启动用户 Web 服务（Go Gin）

```bash
cd atopmall_web/user_web
# 复制配置模板并修改
cp config-debug_templ.yaml config-debug.yaml
go mod tidy
go run main.go
```

> 默认监听端口：8081，启动后从 Consul 发现用户服务地址

## 六、配置说明

项目使用 Viper 管理配置，支持多环境：

| 文件                    | 用途                       | 是否提交 Git |
| ----------------------- | -------------------------- | ------------ |
| config-debug_templ.yaml | 调试配置模板               | ✅ 是        |
| config-debug.yaml       | 本地调试配置（含敏感信息） | ❌ 否        |
| config-pro.yaml         | 生产配置模板               | ✅ 是        |

配置项包含：MySQL、Redis、Consul、JWT、邮箱 SMTP 等。

## 七、服务注册与发现流程

```
──────────────┐     注册      ┌──────────┐     发现      ──────────────┐
│  user_srv    │ ──────────▶  │  Consul  │  ◀──────────  │  user_web    │
│ (Python gRPC)│  启动时注册    │          │  启动时查询    │  (Go Gin)    │
└──────────────┘              └──────────┘               └──────────────┘
                                   │
                              健康检查（GRPC）
                              异常自动注销
```

1. **user_srv** 启动时通过 `python-consul` 注册到 Consul，包含 GRPC 健康检查
2. **user_web** 启动时从 Consul 查询 user_srv 的地址和端口
3. **user_web** 建立 gRPC 长连接，后续请求复用该连接
4. **user_srv** 异常退出时，Consul 自动注销该服务实例

## 八、各服务 README

> 每个微服务将拥有独立的 README 文档，开发中...

| 服务                     | 语言   | 状态   |
| ------------------------ | ------ | ------ |
| user_srv（用户微服务）   | Python | 开发中 |
| user_web（用户 Web API） | Go     | 开发中 |
