# AtopMall 前端项目

基于 Vue 3 + Vite + TypeScript + Element Plus 构建的商城前台和管理后台。

前端基于AI开发部署，具体的Prompt 模板参考PROMPT_TEMPLATES.md文档

## 项目结构

```
atopmall_frontend/
├── mall/                    # 商城前台（用户端）✅ 已完成
│   ├── src/
│   │   ├── api/            # API 接口封装
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 公共组件
│   │   ├── composables/    # 组合式函数（全局状态共享）
│   │   ├── router/         # 路由配置
│   │   ├── store/          # 状态管理（Pinia）
│   │   ├── views/          # 页面视图
│   │   └── utils/          # 工具函数
│   └── package.json
│
├── admin/                   # 管理后台 ✅ 已完成
│   ├── src/
│   │   ├── api/            # API 接口封装
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 公共组件
│   │   ├── layout/         # 布局组件
│   │   ├── router/         # 路由配置
│   │   ├── store/          # 状态管理（Pinia）
│   │   ├── views/          # 页面视图
│   │   └── utils/          # 工具函数
│   └── package.json
│
├── deploy/                  # 部署配置
│   └── nginx_atopmall.conf # Nginx 部署配置
│
├── README.md               # 本文件
└── PROMPT_TEMPLATES.md     # AI 代码生成 Prompt 模板
```

## 技术栈

- **框架**: Vue 3.4+ (Composition API + `<script setup>`)
- **构建工具**: Vite 5.x
- **语言**: TypeScript 5.x
- **UI 组件库**: Element Plus 2.x
- **HTTP 客户端**: Axios 1.x
- **路由**: Vue Router 4.x
- **状态管理**: Pinia 2.x
- **CSS 预处理器**: Sass

## 快速开始

### 1. 安装依赖

```bash
# 安装商城前台依赖
cd mall
npm install

# 安装管理后台依赖（待开发）
cd ../admin
npm install
```

### 2. 启动开发服务器

```bash
# 启动商城前台（端口 3000）
cd mall
npm run dev

# 启动管理后台（端口 3001）
cd admin
npm run dev
```

### 3. 访问应用

- 商城前台：http://localhost:3000
- 管理后台：http://localhost:3001

## 功能模块

### 商城前台（mall）✅ 已完成

| 模块      | 功能                                                                       | 状态  |
| --------- | -------------------------------------------------------------------------- | :---: |
| 首页      | 轮播图、商品分类侧边栏（悬浮展开子分类）、新品推荐、热销商品               |   ✅   |
| 商品列表  | 关键词搜索、分类筛选、品牌筛选、价格区间、新品/热销/推荐标签、分页         |   ✅   |
| 商品详情  | 图片展示、商品信息、购物车实时入口、数量选择、加购/立即购买/收藏、回到顶部 |   ✅   |
| 购物车    | 商品管理、数量修改、全选/批量删除、购物车合计、继续购物入口                |   ✅   |
| 结算页    | 地址选择、订单确认                                                         |   ✅   |
| 订单管理  | 订单列表、订单详情                                                         |   ✅   |
| 用户中心  | 地址管理、收藏列表、留言管理（含附件上传/详情弹窗）                        |   ✅   |
| 登录/注册 | 弹窗模态框（不离开当前页面）、邮箱验证码注册、图形验证码登录               |   ✅   |

### 管理后台（admin）✅ 已完成

| 模块       | 功能                                                    | 状态  |
| ---------- | ------------------------------------------------------- | :---: |
| 控制台     | 数据统计、最新订单、热销商品                            |   ✅   |
| 商品管理   | 商品列表、新增/编辑商品、上下架筛选、分类搜索、价格筛选 |   ✅   |
| 分类管理   | 树形结构、增删改查                                      |   ✅   |
| 品牌管理   | 品牌列表、增删改查                                      |   ✅   |
| 轮播图管理 | 轮播图列表、增删改查                                    |   ✅   |
| 订单管理   | 订单列表、订单详情                                      |   ✅   |
| 用户管理   | 用户列表                                                |   ✅   |
| 留言管理   | 留言列表、详情弹窗、附件预览/下载                       |   ✅   |

## API 接口

所有 API 接口已封装在 `src/api/` 目录下，与后端 Go 服务完全对齐。

### 接口分类

| 模块     | 路径前缀 | 说明                     |
| -------- | -------- | ------------------------ |
| 用户     | /u/v1    | 用户登录、注册、信息管理 |
| 商品     | /g/v1    | 商品、分类、品牌、轮播图 |
| 订单     | /o/v1    | 订单、购物车             |
| 用户操作 | /op/v1   | 地址、留言、收藏         |
| 文件存储 | /oss/v1  | MinIO 文件上传           |

### 商品列表 API 参数映射（前端 → 后端）

| 功能       | 前端参数 | 后端参数 |
| ---------- | -------- | -------- |
| 分页页码   | `p`      | `pn`     |
| 每页条数   | `pnum`   | `pnum`   |
| 搜索关键词 | `q`      | `skey`   |
| 分类筛选   | `c`      | `ctg`    |
| 品牌筛选   | `pn`     | `b`      |
| 最低价格   | `pmin`   | `pmin`   |
| 最高价格   | `pmax`   | `pmax`   |
| 是否热门   | `ishot`  | `ishot`  |
| 是否新品   | `isnew`  | `isnew`  |
| 是否推荐   | `istab`  | `istab`  |

> **注意**：前端请求参数名与后端 Go 服务的 `form` tag 参数名一一对应，不要混淆。

### 认证方式

- **Header**: `x-token: <JWT Token>`
- **Token 获取**: 登录成功后返回
- **Token 存储**: localStorage
- **Token 校验**: Axios 请求拦截器自动注入

## 项目配置

### Vite 代理配置

通过单一网关代理所有后端微服务：

```typescript
proxy: {
  '/u/v1':  { target: 'http://192.168.1.106:8000', changeOrigin: true },
  '/g/v1':  { target: 'http://192.168.1.106:8000', changeOrigin: true },
  '/o/v1':  { target: 'http://192.168.1.106:8000', changeOrigin: true },
  '/op/v1': { target: 'http://192.168.1.106:8000', changeOrigin: true },
  '/oss/v1':{ target: 'http://192.168.1.106:8000', changeOrigin: true }
}
```

### 环境变量

创建 `.env` 文件配置环境变量：

```env
VITE_APP_TITLE=AtopMall
VITE_API_BASE_URL=http://localhost:8000
```

## 开发规范

### 代码风格

- 使用 Composition API + `<script setup>` 语法
- 统一使用 TypeScript 类型定义
- 遵循 Vue 3 官方风格指南

### 目录规范

| 目录           | 用途                         |
| -------------- | ---------------------------- |
| `api/`         | 按模块拆分 API 接口          |
| `views/`       | 按功能模块拆分页面           |
| `components/`  | 公共可复用组件               |
| `composables/` | 组合式函数（跨组件状态共享） |
| `store/`       | 按模块拆分 Pinia 状态管理    |
| `utils/`       | 工具函数                     |

### 命名规范

- 文件名：PascalCase（组件）或 camelCase（工具函数）
- 变量名：camelCase
- 常量名：UPPER_SNAKE_CASE
- 类型名：PascalCase

## 关键设计决策

### 登录/注册弹窗化

登录和注册以模态框形式弹出，不跳转页面，用户可随时关闭弹窗继续浏览。路由守卫触发登录弹窗后，登录成功自动跳转目标页。

实现方式：
- `src/composables/useAuthModal.ts` — 模块级 `ref` 实现跨组件状态共享
- `src/components/AuthModal.vue` — 登录/注册弹窗组件（Tab 切换）
- 路由守卫改用弹窗而非 `/login` 重定向

### 购物车实时反馈

商品详情页底部实时显示购物车商品数量，点击可直达购物车。立即购买前检查购物车中是否已有该商品，避免重复叠加。

### 商品分类侧边栏

首页左侧悬浮下拉菜单，CSS `overflow: visible` 确保下拉面板不被裁剪。

## AI 代码生成

项目提供了完整的 Prompt 模板，用于 AI 辅助代码生成。~~详见 [PROMPT_TEMPLATES.md](./PROMPT_TEMPLATES.md)~~

## 常见问题

### Q: 如何修改 API 地址？

A: 修改 `vite.config.ts` 中的 proxy target 地址。

### Q: 如何处理 Token 过期？

A: 在 `src/utils/request.ts` 的响应拦截器中已统一处理，会自动跳转登录页。

### Q: 如何添加新页面？

A: 
1. 在 `src/views/` 创建页面组件
2. 在 `src/router/index.ts` 添加路由配置
3. 如需权限控制，添加 `meta: { requiresAuth: true }`

### Q: 如何添加新的 API 接口？

A:
1. 在 `src/api/` 对应模块文件中添加接口方法
2. 定义 TypeScript 类型
3. 在页面中调用接口

## 部署

### 构建生产版本

```bash
# 构建商城前台
cd mall
npm run build

# 构建管理后台
cd admin
npm run build
```

### 部署到服务器

将 `dist` 目录下的文件部署到 Nginx 服务器。

**部署路径**：
- Mall 用户前台 → `/usr/share/nginx/html/user/`
- Admin 管理后台 → `/usr/share/nginx/html/admin/`

**访问地址**：
- 用户前台：`http://your-server-ip/`
- 管理后台：`http://your-server-ip/admin`

Nginx 配置参考 `deploy/nginx_atopmall.conf`：

```nginx
server {
    listen 80;
    server_name _;

    set $backend_host 192.168.1.106:8000;

    # API 网关代理
    location /u/v1  { proxy_pass http://$backend_host; ... }
    location /g/v1  { proxy_pass http://$backend_host; ... }
    location /o/v1  { proxy_pass http://$backend_host; ... }
    location /op/v1 { proxy_pass http://$backend_host; ... }
    location /oss/v1{ proxy_pass http://$backend_host; ... }

    # 管理后台（子路径 /admin）
    location ^~ /admin {
        alias /usr/share/nginx/html/admin/;
        try_files $uri $uri/ /admin/index.html;
    }

    # 用户前台（根路径 /）
    root /usr/share/nginx/html/user/;
    index index.html index.htm;
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## 许可证

MIT

## 联系方式

如有问题，请联系yolo_t@outlook.com。