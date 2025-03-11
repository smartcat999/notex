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
  background: linear-gradient(135deg, rgba(43, 88, 118, 0.02), rgba(78, 67, 118, 0.02));
  padding: 20px;

  .not-found-content {
    text-align: center;
    background: rgba(255, 255, 255, 0.98);
    padding: 48px;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(43, 88, 118, 0.08);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 32px rgba(43, 88, 118, 0.12);
    }
    
    h1 {
      font-size: 8rem;
      font-weight: 700;
      margin: 0;
      line-height: 1;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      animation: pulse 2s infinite;
    }

    h2 {
      font-size: 2rem;
      color: #2c3e50;
      margin: 1rem 0;
      font-weight: 600;
    }

    p {
      color: #606266;
      margin-bottom: 2rem;
      font-size: 1.1rem;
    }

    .el-button {
      padding: 12px 32px;
      font-size: 1.1rem;
      font-weight: 500;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      border: none;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
      }

      &:active {
        transform: translateY(0);
      }
    }
  }
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
  100% {
    opacity: 1;
  }
}
</style> 