<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 左侧个人信息 -->
      <el-col :xs="24" :md="8">
        <el-card class="profile-card">
          <div class="profile-header">
            <div class="avatar-upload" @click="handleAvatarClick">
              <FileUploader
                ref="fileUploader"
                v-model="avatarUrl"
                :is-avatar="true"
                :accept="'image/jpeg,image/jpg,image/png,image/gif'"
                :max-size="2 * 1024 * 1024"
                :upload-text="'更换头像'"
                :upload-tip="'支持 jpg、jpeg、png、gif 格式'"
                class="uploader-component"
              >
                <template #upload-area>
                  <div class="avatar-content">
                    <el-avatar 
                      :size="100" 
                      :src="userStore.user?.avatar"
                      class="avatar-image"
                    >
                      {{ userStore.user?.username?.charAt(0) }}
                    </el-avatar>
                    <div class="avatar-hover-mask">
                      <el-icon><Upload /></el-icon>
                      <span>更换头像</span>
                    </div>
                  </div>
                </template>
                <template #preview="{ url }">
                  <div class="avatar-content">
                    <el-avatar 
                      :size="100" 
                      :src="url"
                      class="avatar-image"
                    >
                      {{ userStore.user?.username?.charAt(0) }}
                    </el-avatar>
                    <div class="avatar-hover-mask">
                      <el-icon><Upload /></el-icon>
                      <span>更换头像</span>
                    </div>
                  </div>
                </template>
              </FileUploader>
            </div>
            <div class="user-info">
              <router-link 
                :to="{ name: 'user-profile', params: { id: userStore.user?.id }}" 
                class="username-link"
              >
                <h1 class="username">{{ userStore.user?.username }}</h1>
              </router-link>
              <p class="bio">{{ userStore.user?.bio || '这个用户很懒，还没有写简介' }}</p>
            </div>
          </div>
          <div class="profile-stats">
            <div class="stat-item">
              <span class="number">{{ userStore.user?.post_count || 0 }}</span>
              <span class="label">文章</span>
            </div>
            <div class="stat-item">
              <span class="number">{{ userStore.user?.comment_count || 0 }}</span>
              <span class="label">评论</span>
            </div>
            <div class="stat-item">
              <span class="number">{{ userStore.user?.view_count || 0 }}</span>
              <span class="label">浏览</span>
            </div>
          </div>
        </el-card>

        <el-card class="profile-menu">
          <el-menu
            :default-active="activeMenu"
            class="profile-menu-list"
            @select="handleMenuSelect"
          >
            <el-menu-item index="info">
              <el-icon><User /></el-icon>
              <span>基本信息</span>
            </el-menu-item>
            <el-menu-item index="posts">
              <el-icon><Document /></el-icon>
              <span>我的文章</span>
            </el-menu-item>
            <el-menu-item index="comments">
              <el-icon><ChatDotRound /></el-icon>
              <span>我的评论</span>
            </el-menu-item>
            <el-menu-item index="security">
              <el-icon><Lock /></el-icon>
              <span>安全设置</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>

      <!-- 右侧内容 -->
      <el-col :xs="24" :md="16">
        <!-- 基本信息 -->
        <el-card v-if="activeMenu === 'info'" class="content-card">
          <template #header>
            <h3>基本信息</h3>
          </template>
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
            @submit.prevent="handleSubmit"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" disabled />
            </el-form-item>
            <el-form-item label="个人简介" prop="bio">
              <el-input
                v-model="form.bio"
                type="textarea"
                :rows="4"
                placeholder="介绍一下你自己..."
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="handleSubmit">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 我的文章 -->
        <el-card v-if="activeMenu === 'posts'" class="content-card">
          <template #header>
            <h3>我的文章</h3>
          </template>
          <div class="posts-list">
            <el-card v-for="post in userPosts" :key="post.id" class="post-card" shadow="hover">
              <div class="post-content">
                <div class="post-main">
                  <h4>
                    <router-link :to="`/posts/${post.id}`">{{ post.title }}</router-link>
                  </h4>
                  <p class="post-summary">{{ post.summary }}</p>
                  <div class="post-meta">
                    <span>
                      <el-icon><Calendar /></el-icon>
                      {{ formatDate(post.created_at) }}
                    </span>
                    <span>
                      <el-icon><View /></el-icon>
                      {{ post.views }}
                    </span>
                    <span>
                      <el-icon><ChatDotRound /></el-icon>
                      {{ post.comment_count }}
                    </span>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
          <div class="pagination" v-if="total > 0">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 30, 50]"
              layout="total, sizes, prev, pager, next"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </el-card>

        <!-- 我的评论 -->
        <el-card v-if="activeMenu === 'comments'" class="content-card">
          <template #header>
            <h3>我的评论</h3>
          </template>
          <div class="comments-list">
            <div v-for="comment in userComments" :key="comment.id" class="comment-item">
              <div class="comment-content">
                {{ comment.content }}
              </div>
              <div class="comment-meta">
                <router-link :to="`/posts/${comment.post_id}`" class="post-title">
                  {{ comment.post_title }}
                </router-link>
                <span class="time">{{ formatDate(comment.created_at) }}</span>
              </div>
            </div>
          </div>
          <div class="pagination" v-if="commentTotal > 0">
            <el-pagination
              v-model:current-page="commentPage"
              v-model:page-size="commentPageSize"
              :total="commentTotal"
              :page-sizes="[10, 20, 30, 50]"
              layout="total, sizes, prev, pager, next"
              @size-change="handleCommentSizeChange"
              @current-change="handleCommentCurrentChange"
            />
          </div>
        </el-card>

        <!-- 安全设置 -->
        <el-card v-if="activeMenu === 'security'" class="content-card">
          <template #header>
            <h3>安全设置</h3>
          </template>
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
            @submit.prevent="handlePasswordSubmit"
          >
            <el-form-item label="当前密码" prop="currentPassword">
              <el-input
                v-model="passwordForm.currentPassword"
                type="password"
                show-password
              />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                show-password
              />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="passwordLoading"
                @click="handlePasswordSubmit"
              >
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import {
  User,
  Document,
  ChatDotRound,
  Lock,
  Calendar,
  View,
  Upload,
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword, getProfile } from '@/api/auth'
import { getUserPosts, getUserComments } from '@/api/posts'
import { ElMessage } from 'element-plus'
import FileUploader from '@/components/common/FileUploader.vue'
import axios from 'axios'

const userStore = useUserStore()
const activeMenu = ref('info')
const loading = ref(false)
const passwordLoading = ref(false)

// 基本信息表单
const formRef = ref(null)
const form = ref({
  username: '',
  email: '',
  bio: '',
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3到20个字符之间', trigger: 'blur' },
  ],
  bio: [
    { max: 500, message: '个人简介不能超过500个字符', trigger: 'blur' },
  ],
}

// 密码表单
const passwordFormRef = ref(null)
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.value.newPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validatePass, trigger: 'blur' },
  ],
}

// 文章列表
const userPosts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 评论列表
const userComments = ref([])
const commentTotal = ref(0)
const commentPage = ref(1)
const commentPageSize = ref(10)

// 头像上传相关
const avatarUrl = ref(userStore.user?.avatar || '')

// 添加 ref
const fileUploader = ref(null)

// 处理头像点击
const handleAvatarClick = () => {
  if (fileUploader.value) {
    fileUploader.value.triggerUpload()
  }
}

// 监听头像变化
watch(avatarUrl, async (newUrl) => {
  if (newUrl) {
    try {
      // 构建更新请求数据
      const updateData = {
        username: form.value.username,
        bio: form.value.bio,
        avatar: newUrl
      }
      
      // 发送更新请求
      await updateProfile(updateData)
      // 更新 store 中的用户信息
      userStore.setUser({
        ...userStore.user,
        avatar: newUrl
      })
      ElMessage.success('头像更新成功')
    } catch (error) {
      console.error('Failed to update avatar:', error)
      ElMessage.error('头像更新失败')
    }
  }
})

const handleMenuSelect = (index) => {
  activeMenu.value = index
  if (index === 'posts') {
    fetchUserPosts()
  } else if (index === 'comments') {
    fetchUserComments()
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true
    await updateProfile(form.value)
    await userStore.fetchUserProfile()
    ElMessage.success('个人信息更新成功')
  } catch (error) {
    console.error('Failed to update profile:', error)
    ElMessage.error('个人信息更新失败')
  } finally {
    loading.value = false
  }
}

const handlePasswordSubmit = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    passwordLoading.value = true
    await changePassword({
      old_password: passwordForm.value.currentPassword,
      new_password: passwordForm.value.newPassword,
    })
    // 清空表单
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    }
    // 显示成功提示
    ElMessage.success('密码修改成功')
  } catch (error) {
    // 显示错误提示
    ElMessage.error('密码修改失败')
  } finally {
    passwordLoading.value = false
  }
}

const fetchUserPosts = async () => {
  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize.value
    }
    const response = await getUserPosts(params)
    userPosts.value = response.items || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch user posts:', error)
    ElMessage.error('获取文章列表失败')
  }
}

const fetchUserComments = async () => {
  try {
    const params = {
      page: commentPage.value,
      per_page: commentPageSize.value,
    }
    const response = await getUserComments(params)
    userComments.value = response.items || []
    commentTotal.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch user comments:', error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchUserPosts()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchUserPosts()
}

const handleCommentSizeChange = (val) => {
  commentPageSize.value = val
  fetchUserComments()
}

const handleCommentCurrentChange = (val) => {
  commentPage.value = val
  fetchUserComments()
}

// 获取用户信息
const fetchUserProfile = async () => {
  try {
    const response = await getProfile()
    if (response && response.user) {
      userStore.setUser(response.user)
      // 更新表单数据
      form.value.username = response.user.username
      form.value.email = response.user.email
      form.value.bio = response.user.bio || ''
    }
  } catch (error) {
    console.error('Failed to fetch user profile:', error)
    ElMessage.error('获取用户信息失败')
  }
}

onMounted(() => {
  fetchUserProfile()
})

// 处理头像上传成功
const handleAvatarUpload = async (url) => {
  try {
    const response = await axios.put('/api/user/profile', {
      avatar: url
    }, {
      headers: {
        Authorization: 'Bearer ' + userStore.token
      }
    })
    
    if (response.data) {
      userStore.setAvatar(url)
      ElMessage.success('头像更新成功')
    }
  } catch (error) {
    console.error('Failed to update avatar:', error)
    ElMessage.error('头像更新失败')
  }
}

// 处理上传错误
const handleUploadError = (error) => {
  console.error('Upload error:', error)
  ElMessage.error('上传失败，请重试')
}

// 保存个人信息
const handleSave = async () => {
  try {
    const response = await axios.put('/api/user/profile', {
      nickname: userInfo.value.nickname
    }, {
      headers: {
        Authorization: 'Bearer ' + userStore.token
      }
    })
    
    if (response.data) {
      userStore.setNickname(userInfo.value.nickname)
      ElMessage.success('保存成功')
    }
  } catch (error) {
    console.error('Failed to save profile:', error)
    ElMessage.error('保存失败')
  }
}
</script>

<style lang="scss" scoped>
.profile-container {
  max-width: 1000px;
  margin: 32px auto;
  padding: 0 20px;

  .profile-header {
    background: #ffffff;
    padding: 32px 24px;
    border-radius: 16px;
    color: #2c3e50;
    text-align: center;
    margin-bottom: 32px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(135deg, rgba(54, 153, 255, 0.08), rgba(91, 141, 239, 0.08));
      pointer-events: none;
    }

    .avatar-upload {
      position: relative;
      width: 100px;
      height: 100px;
      margin: 0 auto 24px;
      cursor: pointer;
      border-radius: 50%;
      overflow: hidden;

      :deep(.uploader-component) {
        width: 100%;
        height: 100%;

        .el-upload {
          width: 100%;
          height: 100%;
        }

        .el-upload-dragger {
          width: 100%;
          height: 100%;
          padding: 0;
          border: none;
          background: none;
        }
      }

      .avatar-content {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;

        .avatar-image {
          width: 100%;
          height: 100%;
          border: 2px solid #fff;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
          background: linear-gradient(135deg, #f0f5ff, #e6f0ff);
          transition: all 0.3s ease;

          :deep(.el-avatar__inner) {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(135deg, #f0f5ff, #e6f0ff);
            font-size: 32px;
            font-weight: 500;
            color: #2B5876;
          }

          &:hover {
            border-color: #3699FF;
            box-shadow: 0 6px 16px rgba(54, 153, 255, 0.15);
          }
        }

        .avatar-hover-mask {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          background: rgba(43, 88, 118, 0.75);
          backdrop-filter: blur(2px);
          color: #fff;
          opacity: 0;
          transition: all 0.3s ease;
          gap: 6px;
          z-index: 1;

          .el-icon {
            font-size: 20px;
          }

          span {
            font-size: 12px;
            font-weight: 500;
          }
        }

        &:hover {
          .avatar-hover-mask {
            opacity: 1;
          }
        }
      }
    }

    .user-info {
      .username-link {
        background: #f0f5ff;
        color: #2c3e50;
        border: 1px solid rgba(54, 153, 255, 0.1);
        
        &:hover {
          background: #e6f0ff;
          border-color: rgba(54, 153, 255, 0.2);
          transform: translateY(-1px);
        }

        .username {
          color: #2c3e50;
          text-shadow: none;
          font-weight: 600;
        }
      }

      .bio {
        color: #64748b;
        font-size: 0.95em;
        line-height: 1.6;
        margin: 8px 0;
        padding: 0;
        transition: all 0.3s ease;
        opacity: 0.85;
        
        &:hover {
          color: #2c3e50;
          opacity: 1;
        }
      }
    }

    .profile-stats {
      background: #f8fafc;
      border: 1px solid rgba(0, 0, 0, 0.06);

      &:hover {
        background: #f0f5ff;
        border-color: rgba(54, 153, 255, 0.1);
      }

      .stat-item {
        .number {
          color: #2c3e50;
          text-shadow: none;
          font-weight: 600;
        }

        .label {
          color: #718096;
          text-shadow: none;
        }
      }
    }
  }

  .profile-content {
    display: grid;
    gap: 24px;
    grid-template-columns: 1fr 300px;

    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }

    .main-content {
      .content-card {
        background: rgba(255, 255, 255, 0.98);
        border-radius: 16px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
        overflow: hidden;
        transition: all 0.3s ease;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
          background: rgba(255, 255, 255, 1);
        }

        .card-header {
          padding: 20px 24px;
          border-bottom: 1px solid rgba(0, 0, 0, 0.04);
          background: linear-gradient(to right, rgba(43, 88, 118, 0.05), rgba(43, 88, 118, 0.02));

          h2 {
            margin: 0;
            font-size: 1.4em;
            font-weight: 600;
            background: linear-gradient(135deg, #2B5876, #4E4376);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            letter-spacing: 0.02em;
          }
        }

        .card-content {
          padding: 24px;
        }
      }
    }

    .sidebar {
      .widget {
        background: rgba(255, 255, 255, 0.98);
        border-radius: 16px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
        overflow: hidden;
        margin-bottom: 24px;
        transition: all 0.3s ease;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
          background: rgba(255, 255, 255, 1);
        }

        .widget-header {
          padding: 16px 20px;
          border-bottom: 1px solid rgba(0, 0, 0, 0.04);
          background: linear-gradient(to right, rgba(43, 88, 118, 0.05), rgba(43, 88, 118, 0.02));

          h3 {
            margin: 0;
            font-size: 1.1em;
            font-weight: 600;
            color: #2c3e50;
          }
        }

        .widget-content {
          padding: 20px;
        }
      }
    }
  }

  .action-button {
    background: linear-gradient(135deg, #2B5876, #4E4376);
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 8px;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
    }

    .el-icon {
      margin-right: 6px;
    }
  }

  .el-tabs {
    :deep(.el-tabs__nav-wrap::after) {
      background-color: rgba(43, 88, 118, 0.1);
    }

    :deep(.el-tabs__item) {
      color: #4a5568;
      
      &.is-active {
        color: #2B5876;
      }

      &:hover {
        color: #2B5876;
      }
    }

    :deep(.el-tabs__active-bar) {
      background: linear-gradient(90deg, #2B5876, #4E4376);
    }
  }
}

.profile-card {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(43, 88, 118, 0.08);
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(43, 88, 118, 0.12);
  }

  .profile-header {
    text-align: center;
    margin-bottom: 20px;
    padding: 24px 20px;
    background: linear-gradient(135deg, rgba(43, 88, 118, 0.05), rgba(78, 67, 118, 0.05));

    .avatar-upload {
      position: relative;
      width: 100px;
      height: 100px;
      margin: 0 auto 24px;
      cursor: pointer;
      border-radius: 50%;
      overflow: hidden;

      :deep(.uploader-component) {
        width: 100%;
        height: 100%;

        .el-upload {
          width: 100%;
          height: 100%;
        }

        .el-upload-dragger {
          width: 100%;
          height: 100%;
          padding: 0;
          border: none;
          background: none;
        }
      }

      .avatar-content {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;

        .avatar-image {
          width: 100%;
          height: 100%;
          border: 2px solid #fff;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
          background: linear-gradient(135deg, #f0f5ff, #e6f0ff);
          transition: all 0.3s ease;

          :deep(.el-avatar__inner) {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(135deg, #f0f5ff, #e6f0ff);
            font-size: 32px;
            font-weight: 500;
            color: #2B5876;
          }

          &:hover {
            border-color: #3699FF;
            box-shadow: 0 6px 16px rgba(54, 153, 255, 0.15);
          }
        }

        .avatar-hover-mask {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          background: rgba(43, 88, 118, 0.75);
          backdrop-filter: blur(2px);
          color: #fff;
          opacity: 0;
          transition: all 0.3s ease;
          gap: 6px;
          z-index: 1;

          .el-icon {
            font-size: 20px;
          }

          span {
            font-size: 12px;
            font-weight: 500;
          }
        }

        &:hover {
          .avatar-hover-mask {
            opacity: 1;
          }
        }
      }
    }

    h2 {
      margin: 12px 0 4px;
      font-size: 1.4em;
      font-weight: 600;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .email {
      color: #606266;
      margin: 0;
      font-size: 0.9em;
      opacity: 0.8;
    }
  }

  .profile-stats {
    display: flex;
    justify-content: space-around;
    text-align: center;
    padding: 0 20px 24px;

    .stat-item {
      .number {
        display: block;
        font-size: 1.6em;
        font-weight: 600;
        background: linear-gradient(135deg, #2B5876, #4E4376);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin-bottom: 4px;
      }

      .label {
        color: #606266;
        font-size: 0.9em;
      }
    }
  }
}

.profile-menu {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(43, 88, 118, 0.08);
  overflow: hidden;

  :deep(.el-menu-item) {
    &.is-active {
      background: linear-gradient(135deg, rgba(43, 88, 118, 0.1), rgba(78, 67, 118, 0.1));
      color: #2B5876;
    }

    &:hover {
      background: linear-gradient(135deg, rgba(43, 88, 118, 0.05), rgba(78, 67, 118, 0.05));
    }
  }
}

.content-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(43, 88, 118, 0.08);
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(43, 88, 118, 0.12);
  }

  :deep(.el-card__header) {
    padding: 16px 20px;
    border-bottom: 1px solid rgba(43, 88, 118, 0.08);
    background: linear-gradient(135deg, rgba(43, 88, 118, 0.05), rgba(78, 67, 118, 0.05));

    h3 {
      margin: 0;
      font-size: 1.2em;
      font-weight: 600;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }

  :deep(.el-button--primary) {
    background: linear-gradient(135deg, #2B5876, #4E4376);
    border: none;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
    }
  }

  :deep(.el-input__inner:focus) {
    border-color: #2B5876;
    box-shadow: 0 0 0 2px rgba(43, 88, 118, 0.1);
  }
}

.posts-list {
  .post-card {
    margin-bottom: 20px;

    .post-content {
      .post-main {
        h4 {
          margin: 0 0 8px;
          font-size: 1.2em;

          a {
            color: var(--el-text-color-primary);
            text-decoration: none;

            &:hover {
              color: var(--el-color-primary);
            }
          }
        }

        .post-summary {
          color: var(--el-text-color-regular);
          margin-bottom: 8px;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .post-meta {
          color: var(--el-text-color-secondary);
          font-size: 0.9em;

          span {
            display: inline-flex;
            align-items: center;
            margin-right: 16px;

            .el-icon {
              margin-right: 4px;
            }
          }
        }
      }
    }
  }
}

.comments-list {
  .comment-item {
    padding: 16px;
    border-bottom: 1px solid rgba(43, 88, 118, 0.08);
    transition: all 0.3s ease;

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background: rgba(43, 88, 118, 0.02);
    }

    .comment-content {
      color: #4b5563;
      margin-bottom: 12px;
      line-height: 1.6;
    }

    .comment-meta {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 0.9em;

      .post-title {
        color: #2B5876;
        text-decoration: none;
        font-weight: 500;
        transition: all 0.3s ease;
        padding: 4px 8px;
        margin: -4px -8px;
        border-radius: 4px;

        &:hover {
          color: #4E4376;
          background: rgba(43, 88, 118, 0.08);
          text-decoration: none;
          transform: translateX(4px);
        }
      }

      .time {
        color: #606266;
        font-size: 0.9em;
        opacity: 0.8;
      }
    }
  }
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style> 