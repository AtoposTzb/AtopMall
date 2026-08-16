# AtopMall 前端项目 AI 代码生成 Prompt 模板

## 一、项目架构说明

### 1.1 项目结构
```
atopmall_frontend/
├── mall/                    # 商城前台（用户端）✅ 已完成
│   ├── src/
│   │   ├── api/            # API 接口封装
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 公共组件
│   │   ├── composables/    # 组合式函数（跨组件状态共享）
│   │   ├── router/         # 路由配置
│   │   ├── store/          # 状态管理（Pinia）
│   │   ├── views/          # 页面视图
│   │   └── utils/          # 工具函数
│   └── public/             # 公共资源
│
└── admin/                   # 管理后台 🚧 待开发
    ├── src/
    │   ├── api/            # API 接口封装
    │   ├── assets/         # 静态资源
    │   ├── components/     # 公共组件
    │   ├── router/         # 路由配置
    │   ├── store/          # 状态管理（Pinia）
    │   ├── views/          # 页面视图
    │   └── utils/          # 工具函数
    └── public/             # 公共资源
```

### 1.2 技术栈
- **框架**: Vue 3.4+ (Composition API + `<script setup>`)
- **构建工具**: Vite 5.x
- **语言**: TypeScript 5.x
- **UI 组件库**: Element Plus 2.x
- **HTTP 客户端**: Axios 1.x
- **路由**: Vue Router 4.x
- **状态管理**: Pinia 2.x
- **CSS 预处理器**: Sass

### 1.3 API 文档位置
```
\AtopMall\atopmall接口文档\atopmall_all_api.openapi.json
```

---

## 二、Prompt 模板使用规范

### 2.1 输入格式
每个 Prompt 必须包含以下部分：
```markdown
## 任务目标
[明确说明要生成的模块/功能]

## 技术约束
- 框架版本：Vue 3.4+ / Vite 5.x / TypeScript 5.x
- UI 组件库：Element Plus 2.x
- 状态管理：Pinia 2.x
- HTTP 客户端：Axios 1.x
- 代码风格：Composition API + <script setup>

## 接口文档
[粘贴相关的 OpenAPI JSON 片段]

## 输出要求
1. 文件路径：[指定生成的文件位置]
2. 代码结构：[组件/模块的组织方式]
3. 功能清单：[必须实现的功能点]
4. 错误处理：[异常情况的处理方式]

## 上下文信息
- 项目类型：[商城前台 mall / 管理后台 admin]
- 依赖模块：[已生成的相关文件]
- 特殊要求：[业务逻辑、UI 规范等]
```

### 2.2 输出格式
AI 生成的代码必须遵循以下结构：
```typescript
// 1. 导入依赖
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getGoodsList } from '@/api/goods'
import { useUserStore } from '@/store/user'

// 2. 类型定义
interface DataType {
  // 数据结构
}

// 3. 响应式数据
const state = reactive<DataType>({
  // 初始状态
})

// 4. 方法定义
const handleAction = async () => {
  try {
    // 业务逻辑
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 5. 生命周期
onMounted(() => {
  // 初始化逻辑
})
```

### 2.3 参数说明

| 参数名     | 类型   | 必填 | 说明                         |
| ---------- | ------ | ---- | ---------------------------- |
| 任务目标   | string | ✅    | 明确要生成的模块或功能       |
| 技术约束   | object | ✅    | 框架版本、技术栈要求         |
| 接口文档   | object | ✅    | OpenAPI JSON 片段            |
| 输出要求   | object | ✅    | 文件路径、代码结构、功能清单 |
| 上下文信息 | object | ⚠️    | 项目类型、依赖模块、特殊要求 |

### 2.4 上下文管理机制

#### 2.4.1 上下文传递规则
1. **首轮对话**：提供完整的项目架构说明和技术栈约束
2. **后续对话**：仅需提供变更部分和依赖关系
3. **跨模块引用**：明确说明依赖的接口、组件、工具函数

#### 2.4.2 上下文示例
```markdown
## 上下文信息
- 项目类型：商城前台 mall
- 依赖模块：
  - src/utils/request.ts（已生成）
  - src/api/goods.ts（已生成）
  - src/store/cart.ts（待生成）
- 特殊要求：
  - 购物车数量实时更新
  - 支持批量删除
  - 空状态展示
```

### 2.5 错误处理规则

#### 2.5.1 代码生成错误处理
```typescript
// 1. API 请求错误
try {
  const response = await api.getData()
  return response.data
} catch (error) {
  ElMessage.error('网络错误')
  throw error
}

// 2. 表单验证错误
const valid = await formRef.value?.validate().catch(() => false)
if (!valid) return

// 3. 未登录处理（弹窗模式）
if (!userStore.isAuthenticated) {
  const authModal = useAuthModal()
  authModal.open('login')
  return
}
```

#### 2.5.2 生成失败处理流程
1. **检查接口文档**：确认 OpenAPI JSON 格式正确
2. **验证依赖关系**：确保依赖的模块已生成
3. **调整 Prompt**：简化任务目标，分步生成
4. **手动补充**：对于复杂逻辑，先生成框架再补充细节

---

## 三、分步生成 Prompt 模板

### Prompt 1：项目初始化配置

```markdown
## 任务目标
生成 Vue 3 + Vite + TypeScript 项目的基础配置文件

## 技术约束
- 框架版本：Vue 3.4+ / Vite 5.x / TypeScript 5.x
- UI 组件库：Element Plus 2.x
- 状态管理：Pinia 2.x
- HTTP 客户端：Axios 1.x
- CSS 预处理器：Sass

## 输出要求
生成以下文件：
1. package.json（包含所有依赖）
2. vite.config.ts（Vite 配置）
3. tsconfig.json（TypeScript 配置，使用 paths 代替 baseUrl）
4. src/main.ts（入口文件，全局注册 Element Plus Icons）
5. src/App.vue（根组件，包含 AuthModal）
6. src/env.d.ts（类型声明）
7. index.html（HTML 模板）

## 上下文信息
- 项目类型：[mall / admin]
- 特殊要求：
  - 配置路径别名 @ -> ./src/*（不使用已弃用的 baseUrl）
  - 配置代理解决跨域（单一网关代理所有微服务）
  - 启用 Sass
  - 全局注册 Element Plus 图标
```

### Prompt 2：HTTP 请求封装

```markdown
## 任务目标
封装 Axios HTTP 请求工具，支持请求拦截、响应拦截、错误处理

## 技术约束
- 框架版本：Vue 3.4+ / TypeScript 5.x
- HTTP 客户端：Axios 1.x
- 认证方式：JWT Token（Header: x-token）

## 输出要求
生成文件：src/utils/request.ts

功能要求：
1. 创建 Axios 实例，配置 timeout
2. 请求拦截器：自动添加 x-token 到 Header
3. 响应拦截器：统一处理错误码，提取 response.data
4. 错误处理：网络错误、超时、业务错误
5. Token 过期处理：清除 Token 并提示
6. 支持 GET/POST/PUT/DELETE/PATCH 方法封装

## 上下文信息
- 项目类型：[mall / admin]
- 依赖模块：无
- 特殊要求：
  - Token 存储在 localStorage
  - 错误提示使用 ElMessage
  - 支持文件上传（multipart/form-data）
  - 不使用 baseURL，由 Vite proxy 处理路径映射
```

### Prompt 3：API 接口封装（商品模块）

```markdown
## 任务目标
封装商品相关 API 接口

## 技术约束
- 框架版本：Vue 3.4+ / TypeScript 5.x
- HTTP 客户端：Axios（使用 src/utils/request.ts）

## 接口文档
[粘贴 OpenAPI JSON 中的商品接口片段]

## 输出要求
生成文件：src/api/goods.ts

功能要求：
1. 定义 TypeScript 类型接口（GoodsItem, GoodsListParams, GoodsListResponse）
2. 封装商品列表接口（支持多条件过滤）
3. 封装商品详情接口
4. 导出所有接口方法

## 上下文信息
- 项目类型：[mall / admin]
- 依赖模块：src/utils/request.ts
- 特殊要求：
  - 类型定义完整，字段名与后端 Go 服务 form tag 一致
  - 商品列表参数：p（页码）、pnum（每页条数）、q（搜索）、c（分类）、ishot、isnew、istab
  - 图片字段使用 front_image、images、desc_images
```

### Prompt 4：状态管理（购物车）

```markdown
## 任务目标
实现购物车状态管理模块

## 技术约束
- 框架版本：Vue 3.4+ / TypeScript 5.x
- 状态管理：Pinia 2.x

## 接口文档
[粘贴购物车接口 JSON]

## 输出要求
生成文件：src/store/cart.ts

功能要求：
1. 定义购物车状态类型（CartItem, CartState）
2. 实现购物车列表获取
3. 实现添加商品到购物车
4. 实现修改购物车商品数量/选中状态
5. 实现删除购物车商品
6. 计算属性：totalCount（总商品数）、checkedCount（选中数）、checkedTotalPrice（选中总价）
7. 支持批量操作（全选、批量删除）

## 上下文信息
- 项目类型：商城前台 mall
- 依赖模块：
  - src/utils/request.ts
  - src/api/cart.ts（需先生成）
- 特殊要求：
  - 使用 Composition API 风格
  - 实时计算总价
  - 加入购物车后 ElMessage 提示成功
```

### Prompt 5：页面组件（商品列表）

```markdown
## 任务目标
实现商品列表页面组件

## 技术约束
- 框架版本：Vue 3.4+ / TypeScript 5.x
- UI 组件库：Element Plus 2.x
- 代码风格：Composition API + <script setup>

## 接口文档
[粘贴商品列表接口 JSON]

## 输出要求
生成文件：src/views/goods/GoodsList.vue

功能要求：
1. 商品列表展示（卡片布局，4 列网格）
2. 搜索功能（关键词、分类、价格区间）
3. 筛选功能（新品、热销、推荐标签）
4. 分页功能（Element Plus Pagination）
5. 商品卡片：图片、名称、价格、标签（新品/热销/推荐）
6. 点击跳转商品详情
7. 空状态展示
8. 加载状态

## 上下文信息
- 项目类型：商城前台 mall
- 依赖模块：
  - src/api/goods.ts
  - src/api/category.ts
- 特殊要求：
  - 响应式布局
  - 搜索参数：p（分页）、q（关键词）、c（分类）、pmin/pmax（价格区间）、ishot/isnew/istab（标签）
  - 分类侧边栏支持展开/收起子分类
```

### Prompt 6：路由配置

```markdown
## 任务目标
配置前端路由

## 技术约束
- 框架版本：Vue 3.4+ / TypeScript 5.x
- 路由：Vue Router 4.x

## 输出要求
生成文件：src/router/index.ts

功能要求：
1. 路由模式：history 模式
2. 路由守卫：登录验证（弹窗模式，非页面跳转）
3. 路由懒加载
4. 路由元信息（title、requiresAuth）

路由清单：
- /：首页
- /login：登录页（全屏，保留直接 URL 访问）
- /register：注册页（全屏，保留直接 URL 访问）
- /goods：商品列表
- /goods/:id：商品详情
- /cart：购物车（requiresAuth）
- /checkout：结算页（requiresAuth）
- /order：订单列表（requiresAuth）
- /order/:id：订单详情（requiresAuth）
- /user：用户中心（requiresAuth）
- /user/address：地址管理（requiresAuth）
- /user/favorite：收藏列表（requiresAuth）

## 上下文信息
- 项目类型：[mall / admin]
- 依赖模块：所有 views 组件、src/composables/useAuthModal.ts
- 特殊要求：
  - 需要登录的路由添加 meta.requiresAuth
  - 路由守卫未登录时弹出登录弹窗，登录成功后自动跳转目标页
  - 路由切换时设置页面标题
```

---

## 四、完整生成流程

### 4.1 生成顺序
按照以下顺序执行 Prompt，确保依赖关系正确：

```
1. 项目初始化配置（Prompt 1）
   ↓
2. HTTP 请求封装（Prompt 2）
   ↓
3. API 接口封装（按模块）
   - 商品模块（Prompt 3）
   - 用户模块
   - 购物车模块
   - 订单模块
   - 地址模块
   - 收藏模块
   ↓
4. 状态管理（按模块）
   - 用户状态
   - 购物车状态（Prompt 4）
   - 订单状态
   ↓
5. 组合式函数
   - useAuthModal（登录弹窗状态管理）
   ↓
6. 公共组件
   - AppHeader（导航栏，含搜索、登录弹窗触发）
   - AuthModal（登录/注册弹窗）
   - CategoryNav（分类导航）
   ↓
7. 页面组件（按功能）
   - 首页
   - 商品列表（Prompt 5）
   - 商品详情
   - 购物车
   - 结算页
   - 订单页
   - 用户中心
   ↓
8. 路由配置（Prompt 6）
   ↓
9. 样式优化
```

### 4.2 上下文传递示例

```markdown
## 第 1 轮：项目初始化
[使用 Prompt 1]

## 第 2 轮：HTTP 封装
[使用 Prompt 2]
上下文：项目已初始化，package.json 已生成

## 第 3 轮：商品 API
[使用 Prompt 3]
上下文：request.ts 已生成，可直接使用

## 第 4 轮：购物车状态
[使用 Prompt 4]
上下文：
- request.ts 已生成
- src/api/cart.ts 已生成（需先生成）
- 依赖的接口方法：getCartList, addToCart, updateCartItem, deleteCartItem

## 第 5 轮：商品列表页
[使用 Prompt 5]
上下文：
- src/api/goods.ts 已生成
- src/api/category.ts 已生成
```

### 4.3 错误处理策略

#### 4.3.1 生成失败
```markdown
问题：AI 生成的代码有语法错误
解决：
1. 检查 TypeScript 类型定义
2. 检查导入路径是否正确
3. 检查 Element Plus 组件用法
4. 手动修复或重新生成
```

#### 4.3.2 依赖缺失
```markdown
问题：提示找不到模块
解决：
1. 确认依赖文件已生成
2. 检查文件路径是否正确
3. 按顺序重新生成依赖模块
```

#### 4.3.3 接口不匹配
```markdown
问题：API 调用失败
解决：
1. 检查 OpenAPI JSON 是否最新
2. 检查前端参数名是否与后端 Go 服务的 form tag 一致
3. 检查后端服务是否启动
4. 检查 Vite proxy 跨域配置
```

---

## 五、最佳实践

### 5.1 Prompt 编写技巧
1. **明确具体**：不要说"生成一个页面"，而要说"生成商品列表页，包含搜索、筛选、分页"
2. **提供上下文**：说明依赖的模块、已生成的文件
3. **分步生成**：复杂功能拆分成多个小任务
4. **验证结果**：生成后检查代码结构和类型定义

### 5.2 代码组织规范
```typescript
// 1. 导入顺序
// 1.1 Vue 核心
import { ref, reactive, onMounted } from 'vue'
// 1.2 第三方库
import { ElMessage } from 'element-plus'
// 1.3 项目模块
import { getGoodsList } from '@/api/goods'
import { useUserStore } from '@/store/user'
import { useAuthModal } from '@/composables/useAuthModal'

// 2. 类型定义
interface GoodsItem {
  id: number
  name: string
  shop_price: number
  front_image: string
  goods_brief: string
  is_hot: boolean
  is_new: boolean
}

// 3. 响应式数据
const goodsList = ref<GoodsItem[]>([])
const loading = ref(false)

// 4. 方法
const loadGoods = async () => {
  loading.value = true
  try {
    const res = await getGoodsList({ p: 1, pnum: 12 })
    goodsList.value = (res as any).data || []
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

// 5. 生命周期
onMounted(() => {
  loadGoods()
})
```

### 5.3 性能优化建议
1. **路由懒加载**：`() => import('@/views/GoodsList.vue')`
2. **图片懒加载**：使用 `v-lazy` 指令
3. **列表虚拟滚动**：大数据量使用 `vue-virtual-scroller`
4. **防抖节流**：搜索输入使用防抖
5. **缓存策略**：商品分类数据缓存到 localStorage

---

## 六、常见问题 FAQ

### Q1: 如何处理跨域问题？
**A**: 在 vite.config.ts 中配置代理，所有微服务路径通过单一网关代理：
```typescript
server: {
  proxy: {
    '/u/v1':  { target: 'http://192.168.1.106:8000', changeOrigin: true },
    '/g/v1':  { target: 'http://192.168.1.106:8000', changeOrigin: true },
    '/o/v1':  { target: 'http://192.168.1.106:8000', changeOrigin: true },
    '/op/v1': { target: 'http://192.168.1.106:8000', changeOrigin: true },
    '/oss/v1':{ target: 'http://192.168.1.106:8000', changeOrigin: true }
  }
}
```

### Q2: 如何处理 Token 过期？
**A**: 在响应拦截器中统一处理：
```typescript
service.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      ElMessage.error('登录已过期，请重新登录')
    }
    return Promise.reject(error)
  }
)
```

### Q3: 如何实现文件上传？
**A**: 使用 FormData：
```typescript
const uploadFile = async (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  const response = await request.post('/oss/v1/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return response.url
}
```

### Q4: 如何实现权限控制？
**A**: 使用路由守卫 + 弹窗登录：
```typescript
router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth) {
    const userStore = useUserStore()
    if (!userStore.isAuthenticated) {
      const authModal = useAuthModal()
      authModal.open('login')
      // 登录成功后自动跳转
      const unwatch = setInterval(() => {
        if (userStore.isAuthenticated) {
          clearInterval(unwatch)
          next({ path: to.fullPath, replace: true })
        }
      }, 200)
      return
    }
  }
  next()
})
```

### Q5: 前端请求参数与后端不匹配？
**A**: 检查以下映射关系，确保前端参数名与后端 Go 服务 `form` tag 一致：
| 功能           | 前端参数                | 后端参数                |
| -------------- | ----------------------- | ----------------------- |
| 分页页码       | `p`                     | `pn`                    |
| 搜索关键词     | `q`                     | `skey`                  |
| 分类筛选       | `c`                     | `ctg`                   |
| 品牌筛选       | `pn`                    | `b`                     |
| 价格区间       | `pmin`/`pmax`           | `pmin`/`pmax`           |
| 热门/新品/推荐 | `ishot`/`isnew`/`istab` | `ishot`/`isnew`/`istab` |

---

## 七、附录

### 7.1 API 接口分类

| 模块     | 路径前缀 | 说明                     |
| -------- | -------- | ------------------------ |
| 用户     | /u/v1    | 用户登录、注册、信息管理 |
| 商品     | /g/v1    | 商品、分类、品牌、轮播图 |
| 订单     | /o/v1    | 订单、购物车             |
| 用户操作 | /op/v1   | 地址、留言、收藏         |
| 文件存储 | /oss/v1  | MinIO 文件上传           |

### 7.2 认证方式
- **Header**: `x-token: <JWT Token>`
- **Token 获取**: 登录成功后返回
- **Token 存储**: localStorage
- **Token 过期**: 响应拦截器统一处理

### 7.3 关键设计模式

#### 登录弹窗模式
使用模块级 `ref` 实现跨组件状态共享，登录/注册以模态框形式弹出，不离开当前页面：

**composables/useAuthModal.ts**：
```typescript
import { ref } from 'vue'
export type AuthMode = 'login' | 'register'
const visible = ref(false)
const mode = ref<AuthMode>('login')
export function useAuthModal() {
  const open = (m?: AuthMode) => { mode.value = m || 'login'; visible.value = true }
  const close = () => { visible.value = false }
  const switchMode = (m: AuthMode) => { mode.value = m }
  return { visible, mode, open, close, switchMode }
}
```

#### 购物车即时反馈
- 商品详情页底部显示购物车商品数量（`cartStore.totalCount`），点击可直达购物车
- 立即购买前检查购物车中是否已有该商品，避免重复叠加
- 加购成功显示 ElMessage 提示

### 7.4 错误码说明
| 错误码 | 说明         | 处理方式                 |
| ------ | ------------ | ------------------------ |
| 200    | 成功         | 正常处理                 |
| 400    | 请求参数错误 | 提示用户检查输入         |
| 401    | 未授权       | 提示登录过期，清除 Token |
| 403    | 禁止访问     | 提示权限不足             |
| 404    | 资源不存在   | 提示资源未找到           |
| 500    | 服务器错误   | 提示系统繁忙             |

---

**文档版本**: v2.0  
**最后更新**: 2026-08-16  
**维护者**: yolo和AI