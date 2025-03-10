<template>
  <div class="settings-container">
    <h1>系统设置</h1>

    <el-tabs v-model="activeTab">
      <!-- 基本设置 -->
      <el-tab-pane label="基本设置" name="basic">
        <el-card class="settings-card">
          <template #header>
            <h3>基本设置</h3>
          </template>
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            :rules="basicRules"
            label-width="120px"
            @submit.prevent="handleBasicSubmit"
          >
            <el-form-item label="网站名称" prop="siteName">
              <el-input v-model="basicForm.siteName" />
            </el-form-item>
            <el-form-item label="网站描述" prop="siteDescription">
              <el-input
                v-model="basicForm.siteDescription"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
            <el-form-item label="网站关键词" prop="siteKeywords">
              <el-select
                v-model="basicForm.siteKeywords"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请输入关键词"
              >
                <el-option
                  v-for="keyword in keywords"
                  :key="keyword"
                  :label="keyword"
                  :value="keyword"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="网站 Logo" prop="siteLogo">
              <el-upload
                class="logo-uploader"
                action="/api/upload"
                :show-file-list="false"
                :on-success="handleLogoSuccess"
                :before-upload="beforeLogoUpload"
              >
                <img v-if="basicForm.siteLogo" :src="basicForm.siteLogo" class="logo" />
                <el-icon v-else class="logo-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="basicLoading" @click="handleBasicSubmit">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 邮件设置 -->
      <el-tab-pane label="邮件设置" name="email">
        <el-card class="settings-card">
          <template #header>
            <h3>邮件设置</h3>
          </template>
          <el-form
            ref="emailFormRef"
            :model="emailForm"
            :rules="emailRules"
            label-width="120px"
            @submit.prevent="handleEmailSubmit"
          >
            <el-form-item label="SMTP 服务器" prop="smtpHost">
              <el-input v-model="emailForm.smtpHost" />
            </el-form-item>
            <el-form-item label="SMTP 端口" prop="smtpPort">
              <el-input-number v-model="emailForm.smtpPort" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="SMTP 用户名" prop="smtpUsername">
              <el-input v-model="emailForm.smtpUsername" />
            </el-form-item>
            <el-form-item label="SMTP 密码" prop="smtpPassword">
              <el-input
                v-model="emailForm.smtpPassword"
                type="password"
                show-password
              />
            </el-form-item>
            <el-form-item label="发件人邮箱" prop="fromEmail">
              <el-input v-model="emailForm.fromEmail" />
            </el-form-item>
            <el-form-item label="发件人名称" prop="fromName">
              <el-input v-model="emailForm.fromName" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="emailLoading" @click="handleEmailSubmit">
                保存设置
              </el-button>
              <el-button @click="handleTestEmail">测试邮件</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 安全设置 -->
      <el-tab-pane label="安全设置" name="security">
        <el-card class="settings-card">
          <template #header>
            <h3>安全设置</h3>
          </template>
          <el-form
            ref="securityFormRef"
            :model="securityForm"
            :rules="securityRules"
            label-width="120px"
            @submit.prevent="handleSecuritySubmit"
          >
            <el-form-item label="登录尝试次数" prop="maxLoginAttempts">
              <el-input-number
                v-model="securityForm.maxLoginAttempts"
                :min="1"
                :max="10"
              />
            </el-form-item>
            <el-form-item label="登录锁定时间(分钟)" prop="loginLockDuration">
              <el-input-number
                v-model="securityForm.loginLockDuration"
                :min="1"
                :max="1440"
              />
            </el-form-item>
            <el-form-item label="密码最小长度" prop="minPasswordLength">
              <el-input-number
                v-model="securityForm.minPasswordLength"
                :min="6"
                :max="20"
              />
            </el-form-item>
            <el-form-item label="密码复杂度要求">
              <el-checkbox-group v-model="securityForm.passwordRules">
                <el-checkbox label="uppercase">必须包含大写字母</el-checkbox>
                <el-checkbox label="lowercase">必须包含小写字母</el-checkbox>
                <el-checkbox label="numbers">必须包含数字</el-checkbox>
                <el-checkbox label="special">必须包含特殊字符</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="启用双因素认证">
              <el-switch v-model="securityForm.enable2FA" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="securityLoading" @click="handleSecuritySubmit">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getSettings, updateSettings, testEmail } from '@/api/settings'

const activeTab = ref('basic')
const basicLoading = ref(false)
const emailLoading = ref(false)
const securityLoading = ref(false)

// 基本设置表单
const basicFormRef = ref(null)
const basicForm = ref({
  siteName: '',
  siteDescription: '',
  siteKeywords: [],
  siteLogo: '',
})

const keywords = ref([
  '技术',
  '编程',
  '开发',
  '设计',
  '生活',
  '随笔',
])

const basicRules = {
  siteName: [
    { required: true, message: '请输入网站名称', trigger: 'blur' },
    { max: 50, message: '网站名称不能超过50个字符', trigger: 'blur' },
  ],
  siteDescription: [
    { max: 200, message: '网站描述不能超过200个字符', trigger: 'blur' },
  ],
}

// 邮件设置表单
const emailFormRef = ref(null)
const emailForm = ref({
  smtpHost: '',
  smtpPort: 587,
  smtpUsername: '',
  smtpPassword: '',
  fromEmail: '',
  fromName: '',
})

const emailRules = {
  smtpHost: [
    { required: true, message: '请输入SMTP服务器地址', trigger: 'blur' },
  ],
  smtpPort: [
    { required: true, message: '请输入SMTP端口', trigger: 'blur' },
  ],
  smtpUsername: [
    { required: true, message: '请输入SMTP用户名', trigger: 'blur' },
  ],
  smtpPassword: [
    { required: true, message: '请输入SMTP密码', trigger: 'blur' },
  ],
  fromEmail: [
    { required: true, message: '请输入发件人邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  fromName: [
    { required: true, message: '请输入发件人名称', trigger: 'blur' },
  ],
}

// 安全设置表单
const securityFormRef = ref(null)
const securityForm = ref({
  maxLoginAttempts: 5,
  loginLockDuration: 30,
  minPasswordLength: 8,
  passwordRules: ['lowercase', 'numbers'],
  enable2FA: false,
})

const securityRules = {
  maxLoginAttempts: [
    { required: true, message: '请输入最大登录尝试次数', trigger: 'blur' },
  ],
  loginLockDuration: [
    { required: true, message: '请输入登录锁定时间', trigger: 'blur' },
  ],
  minPasswordLength: [
    { required: true, message: '请输入密码最小长度', trigger: 'blur' },
  ],
}

// Logo 上传
const handleLogoSuccess = (response) => {
  basicForm.value.siteLogo = response.data.url
}

const beforeLogoUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB！')
    return false
  }
  return true
}

// 提交处理
const handleBasicSubmit = async () => {
  if (!basicFormRef.value) return

  try {
    await basicFormRef.value.validate()
    basicLoading.value = true
    await updateSettings('basic', basicForm.value)
    ElMessage.success('基本设置已更新')
  } catch (error) {
    console.error('Failed to update basic settings:', error)
  } finally {
    basicLoading.value = false
  }
}

const handleEmailSubmit = async () => {
  if (!emailFormRef.value) return

  try {
    await emailFormRef.value.validate()
    emailLoading.value = true
    await updateSettings('email', emailForm.value)
    ElMessage.success('邮件设置已更新')
  } catch (error) {
    console.error('Failed to update email settings:', error)
  } finally {
    emailLoading.value = false
  }
}

const handleSecuritySubmit = async () => {
  if (!securityFormRef.value) return

  try {
    await securityFormRef.value.validate()
    securityLoading.value = true
    await updateSettings('security', securityForm.value)
    ElMessage.success('安全设置已更新')
  } catch (error) {
    console.error('Failed to update security settings:', error)
  } finally {
    securityLoading.value = false
  }
}

const handleTestEmail = async () => {
  try {
    await testEmail()
    ElMessage.success('测试邮件发送成功')
  } catch (error) {
    console.error('Failed to send test email:', error)
  }
}

// 获取设置
const fetchSettings = async () => {
  try {
    const response = await getSettings()
    const settings = response.data

    basicForm.value = {
      siteName: settings.basic?.siteName || '',
      siteDescription: settings.basic?.siteDescription || '',
      siteKeywords: settings.basic?.siteKeywords || [],
      siteLogo: settings.basic?.siteLogo || '',
    }

    emailForm.value = {
      smtpHost: settings.email?.smtpHost || '',
      smtpPort: settings.email?.smtpPort || 587,
      smtpUsername: settings.email?.smtpUsername || '',
      smtpPassword: settings.email?.smtpPassword || '',
      fromEmail: settings.email?.fromEmail || '',
      fromName: settings.email?.fromName || '',
    }

    securityForm.value = {
      maxLoginAttempts: settings.security?.maxLoginAttempts || 5,
      loginLockDuration: settings.security?.loginLockDuration || 30,
      minPasswordLength: settings.security?.minPasswordLength || 8,
      passwordRules: settings.security?.passwordRules || ['lowercase', 'numbers'],
      enable2FA: settings.security?.enable2FA || false,
    }
  } catch (error) {
    console.error('Failed to fetch settings:', error)
  }
}

onMounted(() => {
  fetchSettings()
})
</script>

<style lang="scss" scoped>
.settings-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;

  h1 {
    margin-bottom: 30px;
  }
}

.settings-card {
  h3 {
    margin: 0;
  }
}

.logo-uploader {
  :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);

    &:hover {
      border-color: var(--el-color-primary);
    }
  }
}

.logo-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.logo {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: cover;
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style> 