# AtopMall 电商微服务项目

基于 Go + Python 双语言混合开发的电商微服务项目。微服务层使用 Python + gRPC 实现业务逻辑，Web API 层使用 Go + Gin 对外提供 HTTP 接口，通过 gRPC 进行服务间通信。



## 一、项目结构

```
AtopMall/
├── atopmall_srvs/              # 微服务层（Python + gRPC）
│   └── user_srv/               # 用户服务
|   └── 开发中...
│
├── atopmall_web/               # Web API 层（Go + Gin）
│   └── user_web/               # 用户 Web 服务
|   └── 开发中...
│
└── README.md
```

## 二、技术栈总览

| 分类 | 技术选型 | 说明 |
|------|----------|------|
| 开发语言 | Go + Python | Go 负责 Web API 层，Python 负责微服务层 |
| 微服务通信 | gRPC + Protobuf | 服务间远程调用 |
| Web 框架 | Gin | Go HTTP 接口层开发 |
| Python ORM | Peewee | Python 数据库操作 |
| Python 日志 | Loguru | Python 端日志组件 |
| Go 日志 | Zap | Go 端高性能结构化日志 |
| 配置管理 | Viper | YAML 配置文件加载与管理 |
| 数据库 | MySQL | 数据存储 |



## 三、开发工具清单

### Go 工具

| 工具名称 | 核心用途 | 安装方式 |
|----------|----------|----------|
| protoc | Protocol Buffers 编译器 | 官网下载二进制包，配置环境变量 |
| protoc-gen-go | protoc Go 代码生成插件 | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` |
| protoc-gen-go-grpc | gRPC Go 代码生成插件 | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| air | 代码热重载，保存自动重启 | `go install github.com/cosmtrek/air@latest` |

### Python 工具

| 工具名称 | 核心用途 | 安装方式 |
|----------|----------|----------|
| grpcio-tools | Protobuf 代码生成（Python） | `pip install grpcio-tools` |

### Proto 文件生成命令

**Go 端**（在 proto 文件所在目录下执行）：
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative xxx.proto
```

**Python 端**（在 proto 文件所在目录下执行）：

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. xxx.proto
```

## 四、快速开始

### 1. 环境准备
1. 安装 Go 1.22+ 并配置 GOPATH 环境变量
2. 安装 Python 3.10+ 并创建虚拟环境
3. 安装上表中所有开发工具
4. 本地启动 MySQL 数据库或者虚拟机安装Docekr拉取MySQL镜像使用
> 没有开发经验的可以参考我的有道云笔记:
【有道云笔记】项目前期准备
https://share.note.youdao.com/s/QJFUWhau

### 2. 启动用户微服务（Python gRPC）
```bash
cd atopmall_srvs/user_srv
pip install -r requirements.txt
python -m server
```
> 默认监听端口：50051

### 3. 启动用户 Web 服务（Go Gin）
```bash
cd atopmall_web
go mod tidy
cd user_web
go run main.go
```
> 默认监听端口：8081

## 五、各服务 README

> 每个微服务将拥有独立的 README 文档，开发中...

| 服务 | 语言 | 状态 |
|------|------|------|
| user_srv（用户微服务） | Python | 开发中 |
| user_web（用户 Web API） | Go | 开发中 |
