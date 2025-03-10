import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/views/Dashboard.vue'),
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/users/UserList.vue'),
        },
        {
          path: 'posts',
          name: 'posts',
          component: () => import('@/views/posts/PostList.vue'),
        },
        {
          path: 'posts/create',
          name: 'post-create',
          component: () => import('@/views/posts/PostForm.vue'),
        },
        {
          path: 'posts/:id/edit',
          name: 'post-edit',
          component: () => import('@/views/posts/PostForm.vue'),
        },
        {
          path: 'categories',
          name: 'categories',
          component: () => import('@/views/categories/CategoryList.vue'),
        },
        {
          path: 'tags',
          name: 'tags',
          component: () => import('@/views/tags/TagList.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/Settings.vue'),
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFound.vue'),
    },
  ],
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !userStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && userStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router 