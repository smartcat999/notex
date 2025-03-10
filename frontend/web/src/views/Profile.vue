<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 左侧个人信息 -->
      <el-col :xs="24" :md="8">
        <el-card class="profile-card">
          <div class="profile-header">
            <el-avatar :size="100" :src="userStore.user?.avatar">
              {{ userStore.user?.username?.charAt(0) }}
            </el-avatar>
            <h2>{{ userStore.user?.username }}</h2>
            <p class="email">{{ userStore.user?.email }}</p>
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
import { ref, onMounted } from 'vue'
import {
  User,
  Document,
  ChatDotRound,
  Lock,
  Calendar,
  View,
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword, getProfile } from '@/api/auth'
import { getUserPosts, getUserComments } from '@/api/posts'
import { ElMessage } from 'element-plus'

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
  } catch (error) {
    console.error('Failed to update profile:', error)
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
      current_password: passwordForm.value.currentPassword,
      new_password: passwordForm.value.newPassword,
    })
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    }
  } catch (error) {
    console.error('Failed to change password:', error)
  } finally {
    passwordLoading.value = false
  }
}

const fetchUserPosts = async () => {
  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize.value,
    }
    const response = await getUserPosts(params)
    userPosts.value = response.items || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch user posts:', error)
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
    ElMessage.error('获取用户信息失败')
  }
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<style lang="scss" scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  margin-bottom: 20px;

  .profile-header {
    text-align: center;
    margin-bottom: 20px;

    h2 {
      margin: 12px 0 4px;
    }

    .email {
      color: var(--el-text-color-secondary);
      margin: 0;
    }
  }

  .profile-stats {
    display: flex;
    justify-content: space-around;
    text-align: center;

    .stat-item {
      .number {
        display: block;
        font-size: 1.5em;
        font-weight: bold;
        color: var(--el-color-primary);
      }

      .label {
        color: var(--el-text-color-secondary);
        font-size: 0.9em;
      }
    }
  }
}

.profile-menu {
  .profile-menu-list {
    border-right: none;
  }
}

.content-card {
  h3 {
    margin: 0;
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
    border-bottom: 1px solid var(--el-border-color-light);

    &:last-child {
      border-bottom: none;
    }

    .comment-content {
      color: var(--el-text-color-regular);
      margin-bottom: 8px;
    }

    .comment-meta {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 0.9em;

      .post-title {
        color: var(--el-color-primary);
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }

      .time {
        color: var(--el-text-color-secondary);
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