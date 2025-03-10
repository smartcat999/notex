<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-position="top"
    @submit.prevent="handleSubmit"
    class="login-form"
  >
    <el-form-item label="用户名" prop="username">
      <el-input
        v-model="form.username"
        placeholder="请输入用户名"
        :prefix-icon="User"
      />
    </el-form-item>

    <el-form-item label="密码" prop="password">
      <el-input
        v-model="form.password"
        type="password"
        placeholder="请输入密码"
        :prefix-icon="Lock"
        show-password
      />
    </el-form-item>

    <div class="form-footer">
      <el-checkbox v-model="form.remember">记住我</el-checkbox>
      <router-link to="/forgot-password" class="forgot-password">
        忘记密码？
      </router-link>
    </div>

    <el-button
      type="primary"
      native-type="submit"
      :loading="userStore.loading"
      class="submit-btn"
    >
      登录
    </el-button>

    <div class="register-link">
      还没有账号？
      <router-link to="/register">立即注册</router-link>
    </div>
  </el-form>
</template>

<style lang="scss" scoped>
.login-form {
  :deep(.el-form-item__label) {
    font-weight: 500;
    color: var(--el-text-color-primary);
  }

  :deep(.el-input__wrapper) {
    box-shadow: 0 0 0 1px var(--el-border-color) !important;
    border-radius: 8px;
    transition: all 0.2s ease;

    &:hover {
      box-shadow: 0 0 0 1px var(--el-color-primary) !important;
    }

    &.is-focus {
      box-shadow: 0 0 0 1px var(--el-color-primary) !important;
    }
  }

  :deep(.el-input__inner) {
    height: 42px;
  }

  :deep(.el-checkbox__label) {
    color: var(--el-text-color-regular);
  }
}

.form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;

  .forgot-password {
    font-size: 0.875rem;
    color: var(--el-color-primary);
    text-decoration: none;
    transition: color 0.2s ease;

    &:hover {
      color: var(--el-color-primary-dark-2);
      text-decoration: underline;
    }
  }
}

.submit-btn {
  width: 100%;
  height: 42px;
  font-size: 1rem;
  font-weight: 500;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(var(--el-color-primary-rgb), 0.15);
  }
}

.register-link {
  text-align: center;
  font-size: 0.875rem;
  color: var(--el-text-color-regular);

  a {
    color: var(--el-color-primary);
    text-decoration: none;
    margin-left: 0.25rem;
    transition: color 0.2s ease;

    &:hover {
      color: var(--el-color-primary-dark-2);
      text-decoration: underline;
    }
  }
}
</style>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)

const form = reactive({
  username: '',
  password: '',
  remember: false,
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3到20个字符之间', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    const result = await userStore.loginUser({
      username: form.username,
      password: form.password,
    })
    
    if (result) {
      ElMessage.success('登录成功')
      await router.push('/posts')
    }
  } catch (error) {
    console.error('Login error:', error)
    if (error.response) {
      ElMessage.error(error.response.data.message || '登录失败')
    } else {
      ElMessage.error('登录过程中发生错误')
    }
  }
}
</script> 