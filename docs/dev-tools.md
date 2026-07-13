# 开发工具清单

## Go 工具

| 工具名称           | 核心用途                 | 安装方式                                                          |
| ------------------ | ------------------------ | ----------------------------------------------------------------- |
| protoc             | Protocol Buffers 编译器  | 官网下载二进制包，配置环境变量                                    |
| protoc-gen-go      | protoc Go 代码生成插件   | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`  |
| protoc-gen-go-grpc | gRPC Go 代码生成插件     | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| air                | 代码热重载，保存自动重启 | `go install github.com/cosmtrek/air@latest`                       |
| mc                 | MinIO 命令行客户端       | 官网下载二进制包或 `docker pull minio/mc`                         |

## Python 工具

| 工具名称              | 核心用途                    | 安装方式                             |
| --------------------- | --------------------------- | ------------------------------------ |
| grpcio-tools          | Protobuf 代码生成（Python） | `pip install grpcio-tools`           |
| python-consul         | Consul 客户端               | `pip install python-consul`          |
| nacos-sdk-python      | Nacos 客户端（Python）      | `pip install nacos-sdk-python<3.0.0` |
| redis                 | Redis 客户端                | `pip install redis`                  |
| python-redis-lock     | Redis 分布式锁              | `pip install python-redis-lock`      |

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
