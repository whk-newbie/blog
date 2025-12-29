import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  // 公开页面 - 文章
  {
    path: '/articles',
    name: 'Articles',
    component: () => import('@/views/Articles.vue'),
    meta: { title: '文章列表' }
  },
  {
    path: '/article/:slug',
    name: 'ArticleDetail',
    component: () => import('@/views/ArticleDetail.vue'),
    meta: { title: '文章详情' }
  },
  // 管理后台 - 登录
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/Login.vue'),
    meta: { title: '管理员登录' }
  },
  // 管理后台 - 仪表盘
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/admin/Dashboard.vue'),
    meta: { title: '管理后台', requiresAuth: true }
  },
  // 管理后台 - 修改密码
  {
    path: '/admin/password',
    name: 'ChangePassword',
    component: () => import('@/views/admin/ChangePassword.vue'),
    meta: { title: '修改密码', requiresAuth: true }
  },
  // 管理后台 - 文章管理
  {
    path: '/admin/articles',
    name: 'AdminArticles',
    component: () => import('@/views/admin/Articles.vue'),
    meta: { title: '文章管理', requiresAuth: true }
  },
  {
    path: '/admin/articles/create',
    name: 'CreateArticle',
    component: () => import('@/views/admin/ArticleEditor.vue'),
    meta: { title: '新建文章', requiresAuth: true }
  },
  {
    path: '/admin/articles/edit/:id',
    name: 'EditArticle',
    component: () => import('@/views/admin/ArticleEditor.vue'),
    meta: { title: '编辑文章', requiresAuth: true }
  },
  // 管理后台 - 分类管理
  {
    path: '/admin/categories',
    name: 'AdminCategories',
    component: () => import('@/views/admin/Categories.vue'),
    meta: { title: '分类管理', requiresAuth: true }
  },
  // 管理后台 - 标签管理
  {
    path: '/admin/tags',
    name: 'AdminTags',
    component: () => import('@/views/admin/Tags.vue'),
    meta: { title: '标签管理', requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面未找到' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 博客系统` : '博客系统'
  
  const userStore = useUserStore()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    if (!userStore.isLoggedIn()) {
      ElMessage.warning('请先登录')
      next({
        path: '/admin/login',
        query: { redirect: to.fullPath }
      })
      return
    }
  }
  
  // 如果已登录，访问登录页则跳转到管理后台
  if (to.path === '/admin/login' && userStore.isLoggedIn()) {
    next('/admin')
    return
  }
  
  next()
})

export default router

