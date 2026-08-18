import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/AdminLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '控制台' }
      },
      {
        path: 'goods',
        name: 'GoodsList',
        component: () => import('@/views/goods/GoodsList.vue'),
        meta: { title: '商品管理' }
      },
      {
        path: 'goods/create',
        name: 'GoodsCreate',
        component: () => import('@/views/goods/GoodsEdit.vue'),
        meta: { title: '新增商品' }
      },
      {
        path: 'goods/edit/:id',
        name: 'GoodsEdit',
        component: () => import('@/views/goods/GoodsEdit.vue'),
        meta: { title: '编辑商品' }
      },
      {
        path: 'category',
        name: 'CategoryList',
        component: () => import('@/views/category/CategoryList.vue'),
        meta: { title: '分类管理' }
      },
      {
        path: 'brand',
        name: 'BrandList',
        component: () => import('@/views/brand/BrandList.vue'),
        meta: { title: '品牌管理' }
      },
      {
        path: 'banner',
        name: 'BannerList',
        component: () => import('@/views/banner/BannerList.vue'),
        meta: { title: '轮播图管理' }
      },
      {
        path: 'order',
        name: 'OrderList',
        component: () => import('@/views/order/OrderList.vue'),
        meta: { title: '订单管理' }
      },
      {
        path: 'order/:id',
        name: 'OrderDetail',
        component: () => import('@/views/order/OrderDetail.vue'),
        meta: { title: '订单详情' }
      },
      {
        path: 'user',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'message',
        name: 'MessageList',
        component: () => import('@/views/message/MessageList.vue'),
        meta: { title: '留言管理' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.PROD ? '/admin' : '/'),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title || 'AtopMall'} - 管理后台`

  if (to.meta.requiresAuth) {
    const userStore = useUserStore()
    if (!userStore.isAuthenticated) {
      next({ path: '/login', query: { redirect: to.fullPath } })
      return
    }
    if (!userStore.isAdmin) {
      userStore.logout()
      next({ path: '/login', query: { redirect: to.fullPath } })
      return
    }
  }

  next()
})

export default router