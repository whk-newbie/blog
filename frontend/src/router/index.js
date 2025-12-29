import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import i18n from '@/locales'

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
        meta: { titleKey: 'nav.home' }
      },
      {
        path: 'articles',
        name: 'Articles',
        component: () => import('@/views/Articles.vue'),
        meta: { titleKey: 'nav.articles' }
      },
      {
        path: 'article/:slug',
        name: 'ArticleDetail',
        component: () => import('@/views/ArticleDetail.vue'),
        meta: { titleKey: 'article.title' }
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
        meta: { titleKey: 'nav.dashboard', requiresAuth: true }
      },
      {
        path: 'password',
        name: 'ChangePassword',
        component: () => import('@/views/admin/ChangePassword.vue'),
        meta: { titleKey: 'user.changePassword', requiresAuth: true }
      },
      {
        path: 'articles',
        name: 'AdminArticles',
        component: () => import('@/views/admin/Articles.vue'),
        meta: { titleKey: 'nav.articleManage', requiresAuth: true }
      },
      {
        path: 'articles/new',
        name: 'CreateArticle',
        component: () => import('@/views/admin/ArticleEditor.vue'),
        meta: { titleKey: 'nav.newArticle', requiresAuth: true }
      },
      {
        path: 'articles/edit/:id',
        name: 'EditArticle',
        component: () => import('@/views/admin/ArticleEditor.vue'),
        meta: { titleKey: 'article.editArticle', requiresAuth: true }
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: () => import('@/views/admin/Categories.vue'),
        meta: { titleKey: 'nav.categories', requiresAuth: true }
      },
      {
        path: 'tags',
        name: 'AdminTags',
        component: () => import('@/views/admin/Tags.vue'),
        meta: { titleKey: 'nav.tags', requiresAuth: true }
      },
      {
        path: 'stats',
        name: 'VisitStats',
        component: () => import('@/views/admin/VisitStats.vue'),
        meta: { titleKey: 'nav.visits', requiresAuth: true }
      },
      {
        path: 'fingerprints',
        name: 'Fingerprints',
        component: () => import('@/views/admin/Fingerprints.vue'),
        meta: { titleKey: 'nav.fingerprints', requiresAuth: true }
      },
      {
        path: 'crawler',
        name: 'CrawlerMonitor',
        component: () => import('@/views/admin/CrawlerMonitor.vue'),
        meta: { titleKey: 'nav.crawler', requiresAuth: true }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { titleKey: 'app.notFound' }
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
  // 设置页面标题 - 使用i18n
  const titleKey = to.meta.titleKey || (to.meta.title ? `nav.${to.meta.title}` : null)
  if (titleKey) {
    const title = i18n.global.t(titleKey)
    document.title = `${title} - ${i18n.global.t('app.blogSystem')}`
  } else {
    document.title = i18n.global.t('app.blogSystem')
  }
  
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

