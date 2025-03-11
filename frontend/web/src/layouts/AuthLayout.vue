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
        <h2 class="auth-title">{{ title }}</h2>
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
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f5f7fa 0%, #f8f9fb 100%);
}

.auth-header {
  padding: 1.5rem 0;
  background: transparent;

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;

    .logo {
      text-decoration: none;
      display: inline-block;

      h1 {
        font-size: 1.6rem;
        font-weight: 700;
        margin: 0;
        background: linear-gradient(135deg, #2B5876, #4E4376);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: -0.5px;
      }
    }
  }
}

.auth-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  position: relative;
  margin-top: -2rem;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(43, 88, 118, 0.03) 0%, rgba(78, 67, 118, 0.02) 50%, transparent 70%);
    border-radius: 50%;
  }
}

.auth-card {
  width: 100%;
  max-width: 360px;
  background: white;
  border-radius: 16px;
  position: relative;
  z-index: 1;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.04);
  overflow: hidden;

  .auth-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #2c3e50;
    margin: 0;
    padding: 1rem 1.5rem;
    text-align: center;
    position: relative;
  }
}

@media (max-width: 768px) {
  .auth-header .header-content {
    padding: 0 1rem;
  }

  .auth-main {
    padding: 1.5rem 1rem;
  }
}
</style> 