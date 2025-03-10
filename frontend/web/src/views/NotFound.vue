<template>
  <div class="not-found">
    <div class="not-found-content">
      <h1>404</h1>
      <h2>页面不存在</h2>
      <p>{{ countdown }}秒后自动返回首页...</p>
      <el-button type="primary" @click="goHome">立即返回首页</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const countdown = ref(5) // 5秒倒计时
let timer = null

const goHome = () => {
  if (userStore.isAuthenticated) {
    router.push('/')
  } else {
    router.push('/public')
  }
}

onMounted(() => {
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
      goHome()
    }
  }, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style lang="scss" scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  padding: 20px;

  .not-found-content {
    text-align: center;
    
    h1 {
      font-size: 6rem;
      font-weight: 700;
      color: var(--el-color-primary);
      margin: 0;
      line-height: 1;
      background: linear-gradient(to right, var(--el-color-primary), var(--el-color-primary-light-3));
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    h2 {
      font-size: 2rem;
      color: var(--el-text-color-primary);
      margin: 1rem 0;
    }

    p {
      color: var(--el-text-color-secondary);
      margin-bottom: 2rem;
      font-size: 1.1rem;
    }

    .el-button {
      padding: 12px 24px;
      font-size: 1.1rem;
    }
  }
}
</style> 