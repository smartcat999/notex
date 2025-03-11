<template>
  <div class="app">
    <Navbar v-if="shouldShowNavbar" />
    <main :class="{ 'main-content': shouldShowNavbar }">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const userStore = useUserStore()

// 在组件挂载时检查并加载用户信息
onMounted(async () => {
  const token = localStorage.getItem('token')
  if (token) {
    await userStore.fetchUserProfile()
  }
})

const shouldShowNavbar = computed(() => {
  return userStore.isAuthenticated && !isPublicRoute.value
})

const isPublicRoute = computed(() => {
  return route.path === '/public' || route.path === '/login'
})
</script>

<style lang="scss">
html {
  overflow-y: scroll;
}

::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
  &:hover {
    background: #a8a8a8;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.app {
  min-height: 100vh;
}

.main-content {
  padding-top: 64px;
  min-height: calc(100vh - 64px);
}
</style> 