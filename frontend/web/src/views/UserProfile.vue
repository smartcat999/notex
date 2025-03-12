<template>
  <div class="user-profile-container">
    <el-card class="profile-card">
      <div class="profile-header">
        <div class="avatar-section">
          <el-avatar 
            :size="120" 
            :src="user?.avatar"
            class="avatar-image"
          >
            {{ user?.username?.charAt(0) }}
          </el-avatar>
        </div>
        <div class="user-info">
          <h1 class="username">{{ user?.username }}</h1>
          <p class="bio">{{ user?.bio || '这个用户很懒，还没有写简介' }}</p>
        </div>
        <div class="profile-stats">
          <div class="stat-item">
            <span class="number">{{ user?.post_count || 0 }}</span>
            <span class="label">文章</span>
          </div>
          <div class="stat-item">
            <span class="number">{{ user?.comment_count || 0 }}</span>
            <span class="label">评论</span>
          </div>
          <div class="stat-item">
            <span class="number">{{ user?.view_count || 0 }}</span>
            <span class="label">浏览</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 用户文章列表 -->
    <el-card class="posts-card">
      <template #header>
        <div class="card-header">
          <h3>{{ user?.username }} 的文章</h3>
        </div>
      </template>
      <div class="posts-list">
        <el-empty v-if="!posts.length" description="暂无文章" />
        <div v-else v-for="post in posts" :key="post.id" class="post-item">
          <router-link :to="`/posts/${post.id}`" class="post-title">
            {{ post.title }}
          </router-link>
          <p class="post-summary">{{ post.summary }}</p>
          <div class="post-meta">
            <span class="time">
              <el-icon><Calendar /></el-icon>
              {{ formatDate(post.created_at) }}
            </span>
            <span class="views">
              <el-icon><View /></el-icon>
              {{ post.views }}
            </span>
            <span class="comments">
              <el-icon><ChatDotRound /></el-icon>
              {{ post.comment_count }}
            </span>
          </div>
        </div>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Calendar, View, ChatDotRound } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { getUserHome } from '@/api/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const user = ref(null)
const posts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 获取用户主页数据
const fetchUserHome = async () => {
  try {
    const response = await getUserHome(route.params.id, {
      page: currentPage.value,
      per_page: pageSize.value
    })
    user.value = response.user
    posts.value = response.posts.items || []
    total.value = response.posts.total || 0
  } catch (error) {
    console.error('Failed to fetch user home:', error)
    ElMessage.error('获取用户主页失败')
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchUserHome()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchUserHome()
}

onMounted(() => {
  fetchUserHome()
})
</script>

<style lang="scss" scoped>
.user-profile-container {
  max-width: 1000px;
  margin: 32px auto;
  padding: 0 20px;

  .profile-card {
    margin-bottom: 24px;
    border-radius: 16px;
    overflow: hidden;
    background: #ffffff;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);

    .profile-header {
      padding: 32px 24px;
      text-align: center;
      background: linear-gradient(135deg, rgba(54, 153, 255, 0.08), rgba(91, 141, 239, 0.08));
      position: relative;

      .avatar-section {
        margin-bottom: 24px;

        .avatar-image {
          border: 4px solid #ffffff;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
          background: linear-gradient(135deg, #f0f5ff, #e6f0ff);
          font-size: 48px;
          color: #3699FF;
        }
      }

      .user-info {
        margin-bottom: 24px;

        .username {
          font-size: 1.8em;
          font-weight: 600;
          color: #2c3e50;
          margin: 0 0 12px;
        }

        .bio {
          color: #4a5568;
          font-size: 1.1em;
          margin: 0;
          padding: 0 20px;
          line-height: 1.6;
        }
      }

      .profile-stats {
        display: flex;
        justify-content: center;
        gap: 48px;
        padding: 16px 24px;
        background: rgba(255, 255, 255, 0.5);
        border-radius: 12px;
        border: 1px solid rgba(0, 0, 0, 0.06);

        .stat-item {
          .number {
            display: block;
            font-size: 1.8em;
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 4px;
          }

          .label {
            color: #718096;
            font-size: 0.95em;
          }
        }
      }
    }
  }

  .posts-card {
    border-radius: 16px;
    overflow: hidden;
    background: #ffffff;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);

    .card-header {
      h3 {
        margin: 0;
        font-size: 1.2em;
        font-weight: 600;
        color: #2c3e50;
      }
    }

    .posts-list {
      .post-item {
        padding: 20px;
        border-bottom: 1px solid rgba(0, 0, 0, 0.06);

        &:last-child {
          border-bottom: none;
        }

        .post-title {
          font-size: 1.2em;
          font-weight: 500;
          color: #2c3e50;
          text-decoration: none;
          display: block;
          margin-bottom: 8px;

          &:hover {
            color: #3699FF;
          }
        }

        .post-summary {
          color: #4a5568;
          margin: 0 0 12px;
          line-height: 1.6;
        }

        .post-meta {
          display: flex;
          gap: 16px;
          color: #718096;
          font-size: 0.9em;

          span {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }
      }
    }

    .pagination {
      padding: 20px;
      display: flex;
      justify-content: center;
    }
  }
}
</style> 