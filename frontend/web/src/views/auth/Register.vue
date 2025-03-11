<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-position="top"
    @submit.prevent="handleSubmit"
  >
    <el-form-item label="用户名" prop="username">
      <el-input
        v-model="form.username"
        placeholder="请输入用户名"
        :prefix-icon="User"
      />
    </el-form-item>

    <el-form-item label="邮箱" prop="email">
      <el-input
        v-model="form.email"
        type="email"
        placeholder="请输入邮箱"
        :prefix-icon="Message"
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

    <el-form-item label="确认密码" prop="confirmPassword">
      <el-input
        v-model="form.confirmPassword"
        type="password"
        placeholder="请再次输入密码"
        :prefix-icon="Lock"
        show-password
      />
    </el-form-item>

    <el-form-item>
      <el-checkbox v-model="form.agreement">
        我已阅读并同意
        <el-link type="primary" :underline="false">服务条款</el-link>
        和
        <el-link type="primary" :underline="false">隐私政策</el-link>
      </el-checkbox>
    </el-form-item>

    <el-button
      type="primary"
      native-type="submit"
      :loading="userStore.loading"
      class="submit-btn"
    >
      注册
    </el-button>

    <div class="login-link">
      已有账号？
      <router-link to="/login">立即登录</router-link>
    </div>
  </el-form>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { User, Lock, Message } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreement: false,
})

const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const validateAgreement = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请阅读并同意服务条款和隐私政策'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3到20个字符之间', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validatePass, trigger: 'blur' },
  ],
  agreement: [
    { validator: validateAgreement, trigger: 'change' },
  ],
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    await userStore.registerUser({
      username: form.username,
      email: form.email,
      password: form.password,
    })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data.message || '注册失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.el-form {
  padding: 1.5rem;
  
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

  :deep(.el-link) {
    color: #2B5876;
    font-weight: 500;
    padding: 2px 4px;
    border-radius: 4px;
    transition: all 0.3s ease;
    
    &:hover {
      color: #4E4376;
      background: rgba(43, 88, 118, 0.08);
    }
  }
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 12px;
  margin: 1.5rem 0;
  background: linear-gradient(135deg, #2B5876, #4E4376);
  border: none;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(43, 88, 118, 0.15);
  }

  &:active {
    transform: translateY(0);
  }
}

.login-link {
  text-align: center;
  font-size: 0.9rem;
  color: #64748b;

  a {
    color: #2B5876;
    text-decoration: none;
    margin-left: 0.25rem;
    font-weight: 600;
    transition: all 0.3s ease;
    position: relative;
    padding: 4px 8px;
    border-radius: 6px;

    &:hover {
      color: #4E4376;
      background: rgba(43, 88, 118, 0.08);
    }
  }
}
</style>