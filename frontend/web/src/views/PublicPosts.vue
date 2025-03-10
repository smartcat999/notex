<template>
  <div class="public-posts-container">
    <header class="site-header">
      <div class="header-content">
        <div class="logo-section">
          <h1 class="logo">Notex</h1>
          <p class="slogan">记录点滴，分享所得</p>
        </div>
        <div class="auth-actions">
          <router-link to="/login">
            <el-button type="primary" size="large" class="login-button">登录</el-button>
          </router-link>
        </div>
      </div>
    </header>

    <div class="hero-section">
      <div class="hero-content">
        <h2 class="hero-title">欢迎来到 Notex</h2>
        <p class="hero-subtitle">在这里，你可以记录和分享你的学习笔记、技术心得和生活感悟</p>
      </div>
    </div>

    <main class="main-content">
      <div class="content-header">
        <h2>最新文章</h2>
      </div>

      <div class="posts-grid" v-loading="loading">
        <el-card v-for="post in posts" :key="post.id" class="post-card" shadow="hover">
          <div class="post-cover" v-if="post.cover">
            <img :src="post.cover" :alt="post.title">
          </div>
          <div class="post-content">
            <h3 class="post-title">
              <router-link :to="`/posts/${post.id}`">{{ post.title }}</router-link>
            </h3>
            <p class="post-summary">{{ post.summary }}</p>
            <div class="post-meta">
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(post.published_at || post.created_at) }}
              </span>
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ post.views }} 次浏览
              </span>
              <span class="meta-item" v-if="post.category">
                <el-icon><Folder /></el-icon>
                {{ post.category }}
              </span>
            </div>
            <div class="post-footer">
              <div class="post-tags">
                <el-tag
                  v-for="tag in post.tags"
                  :key="tag.id"
                  size="small"
                  effect="plain"
                >
                  {{ tag.name }}
                </el-tag>
              </div>
              <router-link :to="`/posts/${post.id}`" class="read-more">
                阅读更多
                <el-icon><ArrowRight /></el-icon>
              </router-link>
            </div>
          </div>
        </el-card>
      </div>

      <el-empty 
        v-if="!loading && posts.length === 0" 
        description="暂无文章"
        :image-size="200"
      />

      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 30, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Calendar, View, Folder, ArrowRight } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { getPublicPosts } from '@/api/posts'
import { ElMessage } from 'element-plus'

const posts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

const fetchPosts = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    const response = await getPublicPosts(params)
    console.log('API Response:', response)
    posts.value = response.items || []
    total.value = response.total || 0
    console.log('Posts after assignment:', posts.value)
  } catch (error) {
    console.error('Failed to fetch posts:', error)
    ElMessage.error('获取文章列表失败')
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchPosts()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchPosts()
}

onMounted(() => {
  fetchPosts()
})
</script>

<style lang="scss" scoped>
.public-posts-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.site-header {
  background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
  padding: 1.5rem 0;
  color: white;

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .logo-section {
      .logo {
        margin: 0;
        font-size: 2.5rem;
        font-weight: 700;
        letter-spacing: -1px;
      }

      .slogan {
        margin: 0.5rem 0 0;
        font-size: 1.1rem;
        opacity: 0.9;
      }
    }

    .auth-actions {
      .login-button {
        font-size: 1.1rem;
        padding: 12px 30px;
        background: white;
        color: var(--el-color-primary);
        border: none;

        &:hover {
          background: rgba(255, 255, 255, 0.9);
          transform: translateY(-1px);
        }
      }
    }
  }
}

.hero-section {
  background: white;
  padding: 4rem 0;
  text-align: center;
  border-bottom: 1px solid var(--el-border-color-light);

  .hero-content {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 20px;

    .hero-title {
      font-size: 2.5rem;
      font-weight: 700;
      margin: 0 0 1rem;
      background: linear-gradient(to right, var(--el-color-primary), var(--el-color-primary-light-3));
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .hero-subtitle {
      font-size: 1.25rem;
      color: var(--el-text-color-secondary);
      margin: 0;
    }
  }
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 3rem 20px;

  .content-header {
    margin-bottom: 2rem;
    text-align: center;

    h2 {
      font-size: 2rem;
      font-weight: 600;
      margin: 0;
      color: var(--el-text-color-primary);
    }
  }
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;

  .post-card {
    border-radius: 12px;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-5px);
    }

    .post-cover {
      height: 200px;
      overflow: hidden;
      border-radius: 8px;
      margin: -20px -20px 20px;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .post-content {
      .post-title {
        margin: 0 0 1rem;
        font-size: 1.25rem;
        
        a {
          color: var(--el-text-color-primary);
          text-decoration: none;
          
          &:hover {
            color: var(--el-color-primary);
          }
        }
      }

      .post-summary {
        margin: 0 0 1rem;
        color: var(--el-text-color-regular);
        font-size: 0.95rem;
        line-height: 1.6;
      }

      .post-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        margin-bottom: 1rem;
        color: var(--el-text-color-secondary);
        font-size: 0.875rem;

        .meta-item {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }

      .post-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .post-tags {
          display: flex;
          gap: 0.5rem;
          flex-wrap: wrap;
        }

        .read-more {
          display: flex;
          align-items: center;
          gap: 4px;
          color: var(--el-color-primary);
          text-decoration: none;
          font-size: 0.95rem;
          transition: all 0.3s ease;

          &:hover {
            gap: 8px;
          }
        }
      }
    }
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 3rem;
}
</style> 