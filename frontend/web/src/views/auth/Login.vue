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
  padding: 1.5rem;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #2c3e50;
  font-size: 0.9rem;
  padding-bottom: 8px;
}

:deep(.el-input__wrapper) {
  box-shadow: none !important;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #f8fafc;
  transition: all 0.2s ease;

  &:hover {
    border-color: #2B5876;
    background: white;
  }

  &.is-focus {
    border-color: #2B5876;
    background: white;
    box-shadow: 0 2px 8px rgba(43, 88, 118, 0.08) !important;
  }
}

:deep(.el-input__inner) {
  height: 40px;
  color: #2c3e50;

  &::placeholder {
    color: #94a3b8;
  }
}

:deep(.el-checkbox__input) {
  .el-checkbox__inner {
    border-color: #e2e8f0;
    transition: all 0.2s ease;

    &:hover {
      border-color: #2B5876;
    }
  }

  &.is-checked .el-checkbox__inner {
    background: #2B5876;
    border-color: #2B5876;
  }
}

:deep(.el-checkbox__label) {
  color: #64748b;
  font-size: 0.9rem;
}

.form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 1.25rem 0;
}

.forgot-password {
  color: #2B5876;
  font-size: 0.9rem;
  text-decoration: none;

  &:hover {
    color: #4E4376;
    text-decoration: underline;
  }
}

.submit-btn {
  width: 100%;
  height: 40px;
  color: white;
  font-size: 0.95rem;
  font-weight: 500;
  background: #2B5876;
  border: none;
  border-radius: 8px;
  margin: 1.25rem 0;
  transition: all 0.2s ease;

  &:hover {
    background: #4E4376;
    transform: translateY(-1px);
  }

  &:active {
    transform: translateY(0);
  }
}

.register-link {
  text-align: center;
  font-size: 0.9rem;
  color: #64748b;

  a {
    color: #2B5876;
    text-decoration: none;
    font-weight: 500;
    margin-left: 4px;

    &:hover {
      color: #4E4376;
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