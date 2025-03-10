<template>
  <div class="app">
    <Navbar v-if="shouldShowNavbar" />
    <main :class="{ 'main-content': shouldShowNavbar }">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const userStore = useUserStore()

const shouldShowNavbar = computed(() => {
  return userStore.isAuthenticated && !isPublicRoute.value
})

const isPublicRoute = computed(() => {
  return route.path === '/public' || route.path === '/login'
})
</script>

<style lang="scss">
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