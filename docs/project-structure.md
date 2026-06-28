# AtopMall 项目结构

## 一、整体目录结构

```
AtopMall/
├── atopmall_srvs/                    # 微服务层（Python + gRPC）
│   ├── common/                       # 公共模块
│   │   └── register/                 # 服务注册（Consul）
│   │       ├── base.py               # 注册接口抽象基类
│   │       └── consul.py             # Consul 注册实现
│   │
│   ├── user_srv/                     # 用户服务
│   │   ├── handler/                  # gRPC 服务实现
│   │   │   └── user.py               # 用户服务（列表/按ID/按Email/按Mobile/创建/更新/密码校验）
│   │   ├── model/                    # Peewee ORM 数据模型
│   │   │   ── models.py             # User 模型 + BaseModel（逻辑删除/连接池/断线重连）
│   │   ├── proto/                    # Protobuf 定义及生成代码
│   │   │   ├── user.proto            # 用户服务 Protobuf 定义
│   │   │   ├── user_pb2.py           # 生成的消息类
│   │   │   └── user_pb2_grpc.py      # 生成的 gRPC 服务类
│   │   ├── settings/                 # 配置管理
│   │   │   ├── settings.py           # Nacos 配置加载（同步版）
│   │   │   ── settings_asnyc.py     # Nacos 配置加载（异步版 + 配置变更监听）
│   │   ├── tests/                    # gRPC 客户端测试
│   │   │   └── user.py               # 用户服务 gRPC 测试
│   │   ├── requirements.txt          # Python 依赖
│   │   └── server.py                 # gRPC 服务入口（含 Consul 注册 + 优雅退出）
│   │
│   └── goods_srv/                    # 商品服务
│       ├── handler/                  # gRPC 服务实现
│       │   ├── goods.py              # 商品服务（列表/批量/创建/删除/更新/详情）
│       │   ├── category.py           # 分类服务（全部分类/子分类/创建/删除/更新）
│       │   ├── brands.py             # 品牌服务（列表/创建/删除/更新）
│       │   ├── banners.py            # 轮播图服务（列表/创建/删除/更新）
│       │   └── category_brand.py     # 品牌分类服务（列表/按分类查/创建/删除/更新）
│       ├── model/                    # Peewee ORM 数据模型
│       │   └── models.py             # Goods/Category/Brands/Banner/GoodsCategoryBrand 模型
│       ├── proto/                    # Protobuf 定义及生成代码
│       │   ├── goods.proto           # 商品服务 Protobuf 定义（5 个 service）
│       │   ├── goods_pb2.py          # 生成的消息类
│       │   ── goods_pb2_grpc.py     # 生成的 gRPC 服务类
│       ├── settings/                 # 配置管理（Nacos 配置加载）
│       │   └── settings.py           # Nacos 配置加载 + 配置变更监听
│       ├── tests/                    # gRPC 客户端测试
│       │   └── goods.py              # 商品服务 gRPC 测试
│       ├── requirements.txt          # Python 依赖
│       └── server.py                 # gRPC 服务入口（含 Consul 注册 + 优雅退出）
│
├── atopmall_web/                     # Web API 层（Go + Gin）
│   ├── user_web/                     # 用户 Web 服务
│   │   ├── api/                      # HTTP 接口实现
│   │   │   ├── captcha.go            # 图片验证码接口
│   │   │   ├── email_code.go         # 邮箱验证码接口
│   │   │   ├── user.go               # 用户登录/注册/列表接口
│   │   │   └── redis_test/           # Redis 连接测试工具
│   │   │       └── main.go
│   │   ├── config/                   # 配置结构体定义
│   │   │   └── config.go             # ServerConfig / ConsulInfo / NacosInfo 等结构体
│   │   ├── forms/                    # 请求表单验证
│   │   │   ├── user.go               # 登录/注册表单验证
│   │   │   └── email_code.go         # 邮箱验证码表单验证
│   │   ├── global/                   # 全局变量（配置、Redis、翻译器、gRPC 客户端）
│   │   │   ├── global.go             # ServerConfig / Trans / RDB / UserSrvClient
│   │   │   └── responselist/         # 响应结构体定义
│   │   │       └── user.go           # UserResponse 等响应结构体
│   │   ├── initialize/               # 初始化（配置加载、路由、日志、Redis、Consul 服务发现）
│   │   │   ├── config.go             # Viper + Nacos 配置加载
│   │   │   ├── logger.go             # Zap 日志初始化
│   │   │   ├── redis.go              # Redis 客户端初始化
│   │   │   ├── router.go             # Gin 路由注册
│   │   │   ├── src_conn.go           # Consul 服务发现 + gRPC 客户端初始化（含负载均衡）
│   │   │   └── validator_trans.go    # 表单验证器中文翻译
│   │   ├── middlewares/              # 中间件（JWT、CORS、权限）
│   │   │   ├── admin.go              # 管理员权限中间件
│   │   │   ├── cors.go               # CORS 跨域中间件
│   │   │   └── jwt.go                # JWT 认证中间件
│   │   ├── models/                   # 请求模型定义
│   │   │   └── user.go               # CustomClaims 等模型
│   │   ├── proto/                    # Protobuf 定义及生成代码
│   │   │   ├── user.proto            # 用户服务 Protobuf 定义（与 user_srv 共享）
│   │   │   ├── user.pb.go            # 生成的 Go 消息类
│   │   │   └── user_grpc.pb.go       # 生成的 Go gRPC 服务类
│   │   ├── router/                   # 路由分组
│   │   │   ├── base.go               # 基础路由（验证码相关）
│   │   │   └── user.go               # 用户路由（登录/注册/列表）
│   │   ├── utils/                    # 工具函数
│   │   │   ── addr_port.go          # 动态可用端口获取
│   │   ├── validator/                # 自定义验证器
│   │   ├── config-debug_templ.yaml   # 调试配置模板（可提交）
│   │   ├── config-debug.yaml         # Nacos 连接调试配置（含敏感信息，不提交）
│   │   ├── config-pro.yaml           # Nacos 连接生产配置
│   │   └── main.go                   # 服务入口
│   │
│   └── goods_web/                    # 商品 Web 服务
│       ├── api/                      # HTTP 接口实现
│       │   └── goods/
│       │       └── goods.go            # 商品列表接口（多条件过滤）
│       ├── config/                   # 配置结构体定义
│       │   └── config.go             # GoodsSrvConfig / JWTConfig / ConsulConfig / ServerConfig / NacosConfig
│       ├── forms/                    # 请求表单验证
│       │   └── goods.go              # 商品表单验证（待扩展）
│       ├── global/                   # 全局变量（配置、翻译器、gRPC 客户端）
│       │   └── global.go             # ServerConfig / Trans / GoodsSrvClient
│       ├── initialize/               # 初始化（配置加载、路由、日志、Consul 服务发现）
│       │   ├── config.go             # Viper + Nacos 配置加载
│       │   ├── logger.go             # Zap 日志初始化
│       │   ├── router.go             # Gin 路由注册（含 /health 健康检查）
│       │   ├── src_conn.go           # Consul 服务发现 + gRPC 客户端初始化（含负载均衡）
│       │   └── validator_trans.go    # 表单验证器中文翻译
│       ├── middlewares/              # 中间件（JWT、CORS、权限）
│       │   ├── admin.go              # 管理员权限中间件
│       │   ├── cors.go               # CORS 跨域中间件
│       │   └── jwt.go                # JWT 认证中间件
│       ├── models/                   # 请求模型定义
│       │   └── request.go            # 请求模型
│       ├── proto/                    # Protobuf 定义及生成代码
│       │   ├── goods.proto           # 商品服务 Protobuf 定义（与 goods_srv 共享）
│       │   ├── goods.pb.go           # 生成的 Go 消息类
│       │   └── goods_grpc.pb.go      # 生成的 Go gRPC 服务类
│       ├── router/                   # 路由分组
│       │   └── goods.go              # 商品路由（商品列表）
│       ├── utils/                    # 工具函数
│       │   ├── addr_port.go          # 动态可用端口获取
│       │   └── register/
│       │       └── consul/
│       │           └── register.go   # Consul 服务注册（接口 + 实现）
│       ├── config-debug_templ.yaml   # 调试配置模板（可提交）
│       ├── config-debug.yaml         # Nacos 连接调试配置（含敏感信息，不提交）
│       ├── config-pro.yaml           # Nacos 连接生产配置
│       ── main.go                   # 服务入口（初始化 + Consul 注册 + 启动）
│
├── docs/                             # 项目文档
│   ├── project-structure.md          # 项目结构详细说明（本文档）
│   ├── features.md                   # 已完成功能清单
│   ├── dev-tools.md                  # 开发工具清单
│   ├── configuration.md              # 配置说明
│   ├── design/                       # 设计文档
│   │   └── diagrams/                 # 设计图（drawio）
│   ── image/                        # 文档配图
│
├── .gitignore                        # Git 忽略配置
├── README.md                         # 项目主文档
└── go.work                           # Go 工作区配置
```

## 二、目录说明

### atopmall_srvs/ — 微服务层

Python + gRPC 实现的微服务层，每个服务独立目录，共享 `common/` 公共模块。

| 目录               | 说明                                                       |
| ------------------ | ---------------------------------------------------------- |
| `common/register/` | Consul 服务注册公共模块，提供抽象基类和 Consul 实现        |
| `user_srv/`        | 用户微服务，提供用户 CRUD、密码校验等 gRPC 接口            |
| `goods_srv/`       | 商品微服务，提供商品/分类/品牌/轮播图/品牌分类等 gRPC 接口 |

每个微服务目录结构统一：

```
xxx_srv/
├── handler/          # gRPC 服务实现（业务逻辑）
├── model/            # Peewee ORM 数据模型
├── proto/            # Protobuf 定义及生成代码
├── settings/         # 配置管理（Nacos）
├── tests/            # gRPC 客户端测试
├── requirements.txt  # Python 依赖
└── server.py         # 服务入口
```

### atopmall_web/ — Web API 层

Go + Gin 实现的 HTTP 接口层，通过 gRPC 调用微服务层。

| 目录         | 说明                                  |
| ------------ | ------------------------------------- |
| `user_web/`  | 用户 Web 服务（登录/注册/验证码）     |
| `goods_web/` | 商品 Web 服务（商品列表/分类/品牌等） |

每个 Web 服务目录结构统一：

```
xxx_web/
├── api/              # HTTP 接口实现（Handler 层）
├── config/           # 配置结构体定义
├── forms/            # 请求表单验证（Validator）
├── global/           # 全局变量（配置、Redis、翻译器、gRPC 客户端）
├── initialize/       # 服务启动初始化流程
├── middlewares/      # 中间件（JWT、CORS、权限）
├── models/           # 请求模型定义
├── proto/            # Protobuf 定义及生成代码（与微服务层共享 .proto 文件）
├── router/           # 路由分组
── utils/            # 工具函数
├── validator/        # 自定义验证器
├── config-*.yaml     # 配置文件
└── main.go           # 服务入口
```

### docs/ — 项目文档

| 文件/目录              | 说明                          |
| ---------------------- | ----------------------------- |
| `project-structure.md` | 项目结构详细说明（本文档）    |
| `features.md`          | 已完成功能清单                |
| `dev-tools.md`         | 开发工具清单和 Proto 生成命令 |
| `configuration.md`     | 配置说明和加载流程            |
| `design/diagrams/`     | 设计图（drawio 格式）         |
| `image/`               | 文档配图                      |

## 三、关键文件说明

### 微服务层关键文件

| 文件                   | 说明                                                                  |
| ---------------------- | --------------------------------------------------------------------- |
| `server.py`            | 服务入口，包含 gRPC 服务器启动、Consul 注册、优雅退出、Nacos 配置监听 |
| `settings/settings.py` | Nacos 配置加载，解析 MySQL、Consul 等服务配置                         |
| `model/models.py`      | Peewee ORM 模型定义，BaseModel 提供逻辑删除、连接池、断线重连         |
| `handler/*.py`         | gRPC 服务实现，包含业务逻辑和错误处理                                 |

### Web API 层关键文件

| 文件                                | 说明                                                        |
| ----------------------------------- | ----------------------------------------------------------- |
| `main.go`                           | 服务入口，包含初始化流程、路由注册、Consul 注册、服务器启动 |
| `initialize/config.go`              | Viper + Nacos 配置加载，支持配置变更监听                    |
| `initialize/src_conn.go`            | Consul 服务发现 + gRPC 客户端初始化（含负载均衡）           |
| `initialize/router.go`              | Gin 路由注册，含 `/health` 健康检查端点                     |
| `middlewares/jwt.go`                | JWT 认证中间件，解析 x-token 头部                           |
| `router/*.go`                       | 路由分组，定义 API 路径和中间件挂载                         |
| `utils/register/consul/register.go` | Consul 服务注册（接口定义 + 实现），支持 HTTP 健康检查      |
