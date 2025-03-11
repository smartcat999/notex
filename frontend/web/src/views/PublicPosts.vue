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
        <h2 class="hero-title">开始你的创作之旅</h2>
      </div>
    </div>

    <main class="main-content">
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
    // console.log('API Response:', response)
    posts.value = response.items || []
    total.value = response.total || 0
    // console.log('Posts after assignment:', posts.value)
  } catch (error) {
    // console.error('Failed to fetch posts:', error)
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
  background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
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
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
        background: rgba(255, 255, 255, 0.15);
        color: white;
        border: 1px solid rgba(255, 255, 255, 0.3);
        backdrop-filter: blur(10px);
        transition: all 0.3s ease;

        &:hover {
          background: rgba(255, 255, 255, 0.25);
          transform: translateY(-1px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
      }
    }
  }
}

.hero-section {
  background: linear-gradient(to bottom, rgba(43, 88, 118, 0.03), rgba(78, 67, 118, 0.02));
  padding: 1.5rem 0;
  text-align: center;
  border-bottom: 1px solid rgba(43, 88, 118, 0.1);

  .hero-content {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 20px;

    .hero-title {
      font-size: 1.6rem;
      font-weight: 600;
      margin: 0;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 20px;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;

  .post-card {
    border-radius: 12px;
    transition: all 0.3s ease;
    border: none;
    background: rgba(255, 255, 255, 0.98);
    box-shadow: 0 2px 12px rgba(43, 88, 118, 0.08);
    height: 100%;
    display: flex;
    flex-direction: column;

    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 6px 20px rgba(43, 88, 118, 0.12);
      background: rgba(255, 255, 255, 1);
    }

    .post-cover {
      height: 200px;
      overflow: hidden;
      border-radius: 8px 8px 0 0;
      margin: -20px -20px 20px;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .post-content {
      flex: 1;
      display: flex;
      flex-direction: column;

      .post-title {
        margin: 0 0 1rem;
        font-size: 1.25rem;
        
        a {
          color: #2c3e50;
          text-decoration: none;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
          line-height: 1.5;
          
          &:hover {
            color: #2B5876;
          }
        }
      }

      .post-summary {
        margin: 0 0 1rem;
        color: #4a5568;
        font-size: 0.95rem;
        line-height: 1.6;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
        flex: 1;
      }

      .post-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        margin-bottom: 1rem;
        color: #718096;
        font-size: 0.875rem;

        .meta-item {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }

      .post-footer {
        margin-top: auto;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .post-tags {
          display: flex;
          gap: 0.5rem;
          flex-wrap: wrap;

          .el-tag {
            background: rgba(43, 88, 118, 0.06);
            border: none;
            color: #2B5876;
            transition: all 0.3s ease;

            &:hover {
              background: rgba(43, 88, 118, 0.1);
              transform: translateY(-1px);
            }
          }
        }

        .read-more {
          display: flex;
          align-items: center;
          gap: 4px;
          color: #2B5876;
          text-decoration: none;
          font-size: 0.9rem;
          font-weight: 500;
          transition: all 0.3s ease;
          padding: 4px 8px;
          border-radius: 4px;

          &:hover {
            background: rgba(43, 88, 118, 0.08);
            transform: translateX(4px);
          }

          .el-icon {
            transition: transform 0.3s ease;
          }

          &:hover .el-icon {
            transform: translateX(2px);
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .posts-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
    
    .post-card {
      .post-cover {
        height: 180px;
      }
      
      .post-content {
        .post-title {
          font-size: 1.2rem;
        }
        
        .post-summary {
          font-size: 0.9rem;
          -webkit-line-clamp: 2;
        }
        
        .post-meta {
          font-size: 0.8rem;
        }
      }
    }
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 3rem;

  :deep(.el-pagination) {
    --el-pagination-hover-color: #2B5876;
    
    .el-pager li {
      &.is-active {
        background: linear-gradient(135deg, #2B5876, #4E4376);
        color: white;
      }
      
      &:hover:not(.is-active) {
        color: #2B5876;
      }
    }
  }
}
</style> 