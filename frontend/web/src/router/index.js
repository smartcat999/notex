import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import AILayout from '@/components/ai/AILayout.vue'
import AIChat from '@/components/ai/AIChat.vue'
import AISettings from '@/components/ai/AISettings.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/posts',
      name: 'Posts',
      component: () => import('@/views/Posts.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/archives',
      name: 'Archives',
      component: () => import('@/views/Archives.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/posts/new',
      name: 'NewPost',
      component: () => import('@/views/NewPost.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/posts/:id',
      name: 'PostDetail',
      component: () => import('@/views/PostDetail.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/public',
      name: 'PublicPosts',
      component: () => import('@/views/PublicPosts.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/login',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'Login',
          component: () => import('@/views/auth/Login.vue')
        }
      ],
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'Register',
          component: () => import('@/views/auth/Register.vue')
        }
      ],
      meta: { requiresAuth: false }
    },
    {
      path: '/categories',
      name: 'Categories',
      component: () => import('@/views/Categories.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/tags',
      name: 'Tags',
      component: () => import('@/views/Tags.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/Profile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/users/:id',
      name: 'user-profile',
      component: () => import('@/views/UserProfile.vue')
    },
    {
      path: '/drafts',
      name: 'Drafts',
      component: () => import('@/views/Drafts.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/ai',
      component: AILayout,
      redirect: '/ai/chat',
      children: [
        {
          path: 'chat',
          name: 'AIChat',
          component: () => import('@/components/ai/AIChat.vue'),
          meta: {
            requiresAuth: true,
            title: 'AI 助手'
          }
        },
        {
          path: 'settings',
          name: 'AISettings',
          component: () => import('@/components/ai/AISettings.vue'),
          meta: {
            requiresAuth: true,
            title: 'AI 设置'
          }
        },
        {
          path: 'ai-document',
          name: 'WordPolish',
          component: () => import('@/components/ai/WordPolish.vue'),
          meta: {
            requiresAuth: true,
            title: 'AI 文档'
          }
        },
        {
          path: 'ai-writing',
          name: 'AIWriting',
          component: () => import('@/components/ai/AIWriting.vue'),
          meta: {
            requiresAuth: true,
            title: 'AI 写作'
          }
        },
        {
          path: 'image-generator',
          name: 'ImageGenerator',
          component: () => import('@/components/ai/ImageGeneratorPage.vue'),
          meta: {
            requiresAuth: true,
            title: 'AI 图像生成'
          }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    // 如果需要认证但未登录，重定向到公共页面
    next('/public')
  } else if (to.path === '/login' && userStore.isAuthenticated) {
    // 如果已登录还访问登录页，重定向到首页
    next('/')
  } else {
    next()
  }
})

export default router 