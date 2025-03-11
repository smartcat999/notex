<template>
  <div class="categories-container">
    <div class="page-header">
      <div class="title-section">
        <div class="title-wrapper">
          <h1>文章分类</h1>
        </div>
        <p class="subtitle">探索不同领域的知识</p>
      </div>
    </div>

    <div class="categories-content" :class="{ 'empty': categories.length === 0 }">
      <!-- 分类列表 -->
      <template v-if="categories.length > 0">
        <div class="categories-list">
          <div
            v-for="category in categories"
            :key="category.id"
            class="category-card"
            :class="{ active: selectedCategory?.id === category.id }"
            @click="selectCategory(category)"
          >
            <div class="category-info">
              <div class="category-header">
                <h3>{{ category.name }}</h3>
                <span class="post-count">{{ category.post_count }} 篇</span>
              </div>
              <p class="description">{{ category.description || '暂无描述' }}</p>
              <div class="category-meta">
                <span class="update-time">
                  <el-icon><Clock /></el-icon>
                  {{ formatDate(category.updated_at) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 分类下的文章列表 -->
        <div class="category-posts" v-if="selectedCategory">
          <div class="section-header">
            <div class="header-content">
              <h2>{{ selectedCategory.name }}</h2>
              <p class="section-description">{{ selectedCategory.description || '暂无描述' }}</p>
            </div>
          </div>

          <div class="posts-list" v-loading="loading">
            <div v-for="post in posts" :key="post.id" class="post-card">
              <div class="post-content">
                <div class="post-main">
                  <div class="post-header">
                    <router-link :to="`/posts/${post.id}`" class="post-title">
                      {{ post.title }}
                    </router-link>
                    <div class="post-meta">
                      <span>
                        <el-icon><Calendar /></el-icon>
                        {{ formatDate(post.published_at || post.created_at) }}
                      </span>
                      <span v-if="post.views !== undefined">
                        <el-icon><View /></el-icon>
                        {{ post.views || 0 }} 次浏览
                      </span>
                    </div>
                  </div>
                  <p class="post-summary">{{ post.summary }}</p>
                  <div class="post-footer">
                    <div class="post-tags">
                      <span
                        v-for="tag in post.tags"
                        :key="tag.id"
                        class="tag"
                      >
                        {{ tag.name }}
                      </span>
                    </div>
                    <router-link :to="`/posts/${post.id}`" class="read-more">
                      阅读更多
                      <el-icon><ArrowRight /></el-icon>
                    </router-link>
                  </div>
                </div>
                <div class="post-cover" v-if="post.cover">
                  <el-image
                    :src="post.cover"
                    :alt="post.title"
                    fit="cover"
                    class="cover-image"
                  />
                </div>
              </div>
            </div>

            <div v-if="!loading && posts.length === 0" class="empty-posts">
              <el-empty description="该分类下暂无文章" />
            </div>
          </div>

          <!-- 分页 -->
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
        </div>
      </template>

      <div v-else class="empty-state">
        <el-empty>
          <template #image>
            <div class="empty-icon">
              <el-icon><Folder /></el-icon>
            </div>
          </template>
          <template #description>
            <span>暂无分类</span>
          </template>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Calendar, View, Clock, Folder, ArrowRight } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { getCategories, getPosts } from '@/api/posts'
import { ElMessage } from 'element-plus'

const categories = ref([])
const selectedCategory = ref(null)
const posts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

const fetchCategories = async () => {
  try {
    const response = await getCategories()
    categories.value = response.items || []
  } catch (error) {
    console.error('Failed to fetch categories:', error)
    ElMessage.error('获取分类列表失败')
    categories.value = []
  }
}

const fetchPosts = async () => {
  if (!selectedCategory.value) return
  
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      category_id: selectedCategory.value.id
    }
    const response = await getPosts(params)
    posts.value = response.items || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch posts:', error)
    ElMessage.error('获取文章列表失败')
    posts.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const selectCategory = (category) => {
  selectedCategory.value = category
  currentPage.value = 1
  fetchPosts()
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
  fetchCategories()
})
</script>

<style lang="scss" scoped>
.categories-container {
  max-width: 1200px;
  margin: 32px auto;
  padding: 0 20px;

  .page-header {
    margin-bottom: 40px;
    text-align: center;

    .title-section {
      .title-wrapper {
        display: inline-block;
        position: relative;
        margin-bottom: 16px;

        h1 {
          font-size: 1.8em;
          font-weight: 600;
          background: linear-gradient(120deg, #2c3e50, #3498db);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          margin: 0;
          line-height: 1.2;
          letter-spacing: 0.02em;
          position: relative;
          z-index: 1;
        }

        &::after {
          content: '';
          position: absolute;
          bottom: 4px;
          left: -8px;
          right: -8px;
          height: 12px;
          background-color: rgba(64, 158, 255, 0.15);
          border-radius: 6px;
          z-index: 0;
          transform: skew(-12deg);
        }
      }

      .subtitle {
        margin: 0.5rem 0 0;
        color: var(--el-text-color-secondary);
        font-size: 0.95em;
      }
    }
  }
}

.categories-content {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 2rem;

  &.empty {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 60vh;
  }

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.categories-list {
  .category-card {
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-light);
    border-radius: 10px;
    padding: 1rem;
    margin-bottom: 1rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

    &:hover {
      border-color: var(--el-color-primary-light-5);
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
    }

    &.active {
      background: var(--el-color-primary-light-9);
      border-color: var(--el-color-primary);
    }

    .category-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;

      h3 {
        margin: 0;
        font-size: 0.95em;
        font-weight: 500;
        color: var(--el-text-color-primary);
      }

      .post-count {
        font-size: 0.82em;
        color: var(--el-color-primary);
        background: rgba(64, 158, 255, 0.08);
        padding: 4px 10px;
        border-radius: 6px;
        transition: all 0.3s ease;

        &:hover {
          background: rgba(64, 158, 255, 0.12);
        }
      }
    }

    .description {
      margin: 0.5rem 0;
      font-size: 0.85em;
      color: var(--el-text-color-secondary);
      line-height: 1.5;
    }

    .category-meta {
      margin-top: 0.5rem;
      font-size: 0.75em;
      color: var(--el-text-color-secondary);

      .update-time {
        display: flex;
        align-items: center;
        gap: 0.25rem;
      }
    }
  }
}

.category-posts {
  .section-header {
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    background: linear-gradient(to right, rgba(64, 158, 255, 0.05), rgba(64, 158, 255, 0.02));
    padding: 20px;
    border-radius: 10px;

    .header-content {
      h2 {
        margin: 0;
        font-size: 1.4em;
        font-weight: 600;
        background: linear-gradient(120deg, #2c3e50, #3498db);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
      }

      .section-description {
        margin: 0.5rem 0 0;
        color: var(--el-text-color-secondary);
        font-size: 0.85em;
      }
    }
  }

  .post-card {
    background: var(--el-bg-color);
    border: none;
    border-radius: 10px;
    padding: 1.5rem;
    margin-bottom: 1rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);

      .post-cover {
        transform: perspective(1000px) rotateY(-2deg) translateZ(10px);
        box-shadow: 6px 6px 16px rgba(0, 0, 0, 0.15);

        img {
          transform: scale(1.05);
        }
      }
    }

    .post-content {
      display: flex;
      gap: 1.5rem;

      .post-main {
        flex: 1;

        .post-header {
          margin-bottom: 1rem;

          .post-title {
            display: block;
            font-size: 1.15em;
            font-weight: 500;
            color: var(--el-text-color-primary);
            text-decoration: none;
            margin-bottom: 0.5rem;
            transition: color 0.3s ease;

            &:hover {
              color: var(--el-color-primary);
            }
          }

          .post-meta {
            display: flex;
            gap: 1rem;
            font-size: 0.82em;
            color: var(--el-text-color-secondary);

            span {
              display: flex;
              align-items: center;
              gap: 0.25rem;
            }
          }
        }

        .post-summary {
          margin: 0;
          font-size: 0.9em;
          color: var(--el-text-color-regular);
          line-height: 1.6;
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .post-footer {
          margin-top: 1rem;
          display: flex;
          justify-content: space-between;
          align-items: center;

          .post-tags {
            display: flex;
            gap: 0.5rem;
            flex-wrap: wrap;

            .tag {
              font-size: 0.75em;
              color: var(--el-color-primary);
              background: rgba(64, 158, 255, 0.08);
              padding: 4px 10px;
              border-radius: 6px;
              transition: all 0.3s ease;

              &:hover {
                background: rgba(64, 158, 255, 0.12);
                transform: translateY(-1px);
              }
            }
          }

          .read-more {
            display: flex;
            align-items: center;
            gap: 0.25rem;
            font-size: 0.85em;
            color: var(--el-color-primary);
            text-decoration: none;
            transition: all 0.3s ease;
            background: rgba(64, 158, 255, 0.08);
            padding: 4px 12px;
            border-radius: 6px;

            &:hover {
              background: rgba(64, 158, 255, 0.12);
              gap: 0.5rem;
            }
          }
        }
      }

      .post-cover {
        width: 200px;
        height: 140px;
        border-radius: 8px;
        overflow: hidden;
        position: relative;
        transform: perspective(1000px) rotateY(0deg) translateZ(0);
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: linear-gradient(
            to bottom,
            rgba(0, 0, 0, 0),
            rgba(0, 0, 0, 0.02) 40%,
            rgba(0, 0, 0, 0.08)
          );
          z-index: 1;
          pointer-events: none;
          transition: opacity 0.3s ease;
        }

        .cover-image {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
        }
      }
    }
  }
}

.empty-state {
  text-align: center;
  padding: 40px;
  
  .empty-icon {
    display: inline-flex;
    padding: 2rem;
    border-radius: 50%;
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
    font-size: 2rem;
    margin-bottom: 1rem;
    transition: all 0.3s ease;

    &:hover {
      transform: scale(1.05);
      background: var(--el-color-primary-light-8);
    }
  }
}

.pagination {
  margin-top: 2rem;
  display: flex;
  justify-content: center;

  :deep(.el-pagination) {
    --el-pagination-button-bg-color: rgba(255, 255, 255, 0.9);
    --el-pagination-hover-color: var(--el-color-primary);
    
    .el-pager li {
      border-radius: 6px;
      transition: all 0.3s ease;
      
      &:hover, &.is-active {
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
      }
    }
  }
}

@media (max-width: 768px) {
  .categories-content {
    grid-template-columns: 1fr;
  }
}
</style> 