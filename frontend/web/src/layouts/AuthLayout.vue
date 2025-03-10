<template>
  <div class="auth-layout">
    <header class="auth-header">
      <div class="header-content">
        <router-link to="/public" class="logo">
          <h1>Notex</h1>
        </router-link>
      </div>
    </header>

    <main class="auth-main">
      <div class="auth-card">
        <h2 class="auth-title">登录</h2>
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const title = computed(() => {
  switch (route.name) {
    case 'login':
      return '登录'
    case 'register':
      return '注册'
    case 'forgot-password':
      return '忘记密码'
    case 'reset-password':
      return '重置密码'
    default:
      return ''
  }
})

const showLoginLink = computed(() => {
  return ['register', 'forgot-password', 'reset-password'].includes(route.name)
})

const showRegisterLink = computed(() => {
  return ['login', 'forgot-password', 'reset-password'].includes(route.name)
})
</script>

<style lang="scss" scoped>
.auth-layout {
  min-height: 100vh;
  background: linear-gradient(
    135deg,
    rgba(var(--el-color-primary-rgb), 0.05) 0%,
    rgba(var(--el-color-primary-rgb), 0.02) 100%
  );
}

.auth-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 50;
  backdrop-filter: blur(12px);
  background: rgba(255, 255, 255, 0.8);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem 2rem;

    .logo {
      text-decoration: none;

      h1 {
        font-size: 1.75rem;
        font-weight: 700;
        margin: 0;
        background: linear-gradient(to right, var(--el-color-primary), var(--el-color-primary-light-3));
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: -0.5px;
      }
    }
  }
}

.auth-main {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;
  background: var(--el-bg-color);
  border-radius: 12px;
  box-shadow: 0 4px 24px -4px rgba(0, 0, 0, 0.08);

  .auth-title {
    font-size: 1.75rem;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0 0 2rem;
    text-align: center;
  }
}

@media (max-width: 768px) {
  .auth-header .header-content {
    padding: 1rem;
  }

  .auth-main {
    padding: 1rem;
  }

  .auth-card {
    padding: 1.5rem;
  }
}
</style> 