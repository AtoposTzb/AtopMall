import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useAuthModal } from '@/composables/useAuthModal'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/goods',
    name: 'GoodsList',
    component: () => import('@/views/goods/GoodsList.vue'),
    meta: { title: '商品列表' }
  },
  {
    path: '/goods/:id',
    name: 'GoodsDetail',
    component: () => import('@/views/goods/GoodsDetail.vue'),
    meta: { title: '商品详情' }
  },
  {
    path: '/cart',
    name: 'Cart',
    component: () => import('@/views/Cart.vue'),
    meta: { title: '购物车', requiresAuth: true }
  },
  {
    path: '/checkout',
    name: 'Checkout',
    component: () => import('@/views/Checkout.vue'),
    meta: { title: '结算', requiresAuth: true }
  },
  {
    path: '/order',
    name: 'OrderList',
    component: () => import('@/views/order/OrderList.vue'),
    meta: { title: '订单列表', requiresAuth: true }
  },
  {
    path: '/order/:id',
    name: 'OrderDetail',
    component: () => import('@/views/order/OrderDetail.vue'),
    meta: { title: '订单详情', requiresAuth: true }
  },
  {
    path: '/user',
    name: 'UserCenter',
    component: () => import('@/views/user/UserCenter.vue'),
    meta: { title: '用户中心', requiresAuth: true },
    children: [
      {
        path: '',
        name: 'UserProfile',
        component: () => import('@/views/user/UserProfile.vue'),
        meta: { title: '个人信息' }
      },
      {
        path: 'address',
        name: 'UserAddress',
        component: () => import('@/views/user/UserAddress.vue'),
        meta: { title: '地址管理' }
      },
      {
        path: 'favorite',
        name: 'UserFavorite',
        component: () => import('@/views/user/UserFavorite.vue'),
        meta: { title: '我的收藏' }
      },
      {
        path: 'cart',
        name: 'UserCart',
        component: () => import('@/views/Cart.vue'),
        meta: { title: '购物车' }
      },
      {
        path: 'orders',
        name: 'UserOrders',
        component: () => import('@/views/order/OrderList.vue'),
        meta: { title: '我的订单' }
      },
      {
        path: 'message',
        name: 'UserMessage',
        component: () => import('@/views/user/UserMessage.vue'),
        meta: { title: '我的留言' }
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
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || 'AtopMall'} - AtopMall 商城`

  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    const userStore = useUserStore()
    if (!userStore.isAuthenticated) {
      const authModal = useAuthModal()
      authModal.open('login')
      // 监听登录成功，自动跳转目标页
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

export default router