import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const routes = [
  // 公开页面 - 使用主布局
  {
    path: '/',
    component: () => import('@/components/layout/MainLayout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'articles',
        name: 'Articles',
        component: () => import('@/views/Articles.vue'),
        meta: { title: '文章列表' }
      },
      {
        path: 'article/:slug',
        name: 'ArticleDetail',
        component: () => import('@/views/ArticleDetail.vue'),
        meta: { title: '文章详情' }
      }
    ]
  },
  // 管理后台 - 使用布局
  {
    path: '/admin',
    component: () => import('@/components/layout/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Admin',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '仪表盘', requiresAuth: true }
      },
      {
        path: 'password',
        name: 'ChangePassword',
        component: () => import('@/views/admin/ChangePassword.vue'),
        meta: { title: '修改密码', requiresAuth: true }
      },
      {
        path: 'articles',
        name: 'AdminArticles',
        component: () => import('@/views/admin/Articles.vue'),
        meta: { title: '文章管理', requiresAuth: true }
      },
      {
        path: 'articles/new',
        name: 'CreateArticle',
        component: () => import('@/views/admin/ArticleEditor.vue'),
        meta: { title: '新建文章', requiresAuth: true }
      },
      {
        path: 'articles/edit/:id',
        name: 'EditArticle',
        component: () => import('@/views/admin/ArticleEditor.vue'),
        meta: { title: '编辑文章', requiresAuth: true }
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: () => import('@/views/admin/Categories.vue'),
        meta: { title: '分类管理', requiresAuth: true }
      },
      {
        path: 'tags',
        name: 'AdminTags',
        component: () => import('@/views/admin/Tags.vue'),
        meta: { title: '标签管理', requiresAuth: true }
      }
    ]
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
  // 如果未登录但需要认证，允许访问但会在AdminLayout中显示登录弹窗
  if (to.meta.requiresAuth && !userStore.isLoggedIn()) {
    // 允许访问，但会在AdminLayout中触发登录弹窗
    next()
    return
  }
  
  next()
})

export default router

