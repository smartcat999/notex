<template>
  <div class="home">
    <el-row :gutter="24">
      <el-col :span="16">
        <div class="content-section">
          <div class="section-header">
            <h2>最新文章</h2>
            <router-link :to="{ name: 'Posts' }" class="view-all">
              查看全部
              <el-icon><ArrowRight /></el-icon>
            </router-link>
          </div>

          <div class="posts-list" v-loading="loading.posts">
            <el-card v-for="post in recentPosts" :key="post.id" class="post-card">
              <div class="post-cover" v-if="post.cover">
                <img :src="post.cover" :alt="post.title">
              </div>
              <div class="post-content">
                <h3 class="post-title">
                  <router-link :to="`/posts/${post.id}`">{{ post.title }}</router-link>
                </h3>
                <p class="post-excerpt">{{ post.excerpt }}</p>
                <div class="post-meta">
                  <span class="post-date">
                    <el-icon><Calendar /></el-icon>
                    {{ formatDate(post.created_at) }}
                  </span>
                  <span class="post-category" v-if="post.category">
                    <el-icon><Folder /></el-icon>
                    {{ post.category }}
                  </span>
                  <span class="post-views">
                    <el-icon><View /></el-icon>
                    {{ post.views }}
                  </span>
                </div>
              </div>
            </el-card>

            <el-empty
              v-if="!loading.posts && recentPosts.length === 0"
              description="暂无文章"
            />
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="sidebar">
          <el-card class="widget categories-widget" v-loading="loading.categories">
            <template #header>
              <div class="widget-header">
                <h3>分类</h3>
                <router-link :to="{ name: 'Categories' }" class="view-all">
                  更多
                  <el-icon><ArrowRight /></el-icon>
                </router-link>
              </div>
            </template>
            <ul class="categories-list">
              <li v-for="category in topCategories" :key="category.id">
                <router-link :to="{ name: 'Posts', query: { category_id: category.id } }" class="category-item">
                  <span class="category-name">{{ category.name }}</span>
                  <span class="category-count">{{ category.post_count }}</span>
                </router-link>
              </li>
            </ul>
            <el-empty
              v-if="!loading.categories && topCategories.length === 0"
              description="暂无分类"
            />
          </el-card>

          <el-card class="widget tags-widget" v-loading="loading.tags">
            <template #header>
              <div class="widget-header">
                <h3>标签</h3>
                <router-link :to="{ name: 'Tags' }" class="view-all">
                  更多
                  <el-icon><ArrowRight /></el-icon>
                </router-link>
              </div>
            </template>
            <div class="tags-cloud">
              <router-link
                v-for="tag in topTags"
                :key="tag.id"
                :to="{ name: 'Tags', query: { tag_id: tag.id } }"
                class="tag-item"
                :style="{ fontSize: tagSize(tag.post_count) + 'px' }"
              >
                {{ tag.name }}
                <span class="tag-count">({{ tag.post_count }})</span>
              </router-link>
            </div>
            <el-empty
              v-if="!loading.tags && topTags.length === 0"
              description="暂无标签"
            />
          </el-card>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { 
  Calendar,
  Folder,
  View,
  ArrowRight
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getRecentPosts, getTopCategories, getTopTags } from '@/api/posts'

const recentPosts = ref([])
const topCategories = ref([])
const topTags = ref([])

const loading = ref({
  posts: false,
  categories: false,
  tags: false
})

const fetchRecentPosts = async () => {
  loading.value.posts = true
  try {
    const response = await getRecentPosts()
    recentPosts.value = response.items || []
  } catch (error) {
    console.error('Failed to fetch recent posts:', error)
    ElMessage.error('获取最新文章失败')
    recentPosts.value = []
  } finally {
    loading.value.posts = false
  }
}

const fetchTopCategories = async () => {
  loading.value.categories = true
  try {
    const response = await getTopCategories({ limit: 5 })
    topCategories.value = response.items || []
  } catch (error) {
    console.error('Failed to fetch top categories:', error)
    ElMessage.error('获取分类列表失败')
    topCategories.value = []
  } finally {
    loading.value.categories = false
  }
}

const fetchTopTags = async () => {
  loading.value.tags = true
  try {
    const response = await getTopTags({ limit: 10 })
    topTags.value = response.items || []
  } catch (error) {
    console.error('Failed to fetch top tags:', error)
    ElMessage.error('获取标签列表失败')
    topTags.value = []
  } finally {
    loading.value.tags = false
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

const tagSize = (count) => {
  if (topTags.value.length === 0) return 14
  const minSize = 12
  const maxSize = 20
  const minCount = Math.min(...topTags.value.map(t => t.post_count))
  const maxCount = Math.max(...topTags.value.map(t => t.post_count))
  if (minCount === maxCount) return 14
  return minSize + (count - minCount) * (maxSize - minSize) / (maxCount - minCount)
}

onMounted(async () => {
  await Promise.all([
    fetchRecentPosts(),
    fetchTopCategories(),
    fetchTopTags()
  ])
})
</script>

<style scoped lang="scss">
.home {
  max-width: 1200px;
  margin: 24px auto;
  padding: 0 20px;
}

.content-section {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h2 {
      margin: 0;
      font-size: 1.4em;
      font-weight: 600;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: 0.02em;
      position: relative;
      display: inline-block;

      &::after {
        content: '';
        position: absolute;
        bottom: 2px;
        left: -8px;
        right: -8px;
        height: 10px;
        background-color: rgba(43, 88, 118, 0.12);
        border-radius: 6px;
        z-index: -1;
        transform: skew(-12deg);
      }
    }

    .view-all {
      display: flex;
      align-items: center;
      gap: 4px;
      color: #2B5876;
      text-decoration: none;
      font-size: 0.9em;
      padding: 6px 12px;
      border-radius: 6px;
      transition: all 0.3s ease;
      background: rgba(43, 88, 118, 0.08);

      &:hover {
        background: rgba(43, 88, 118, 0.12);
        transform: translateX(2px);
      }
    }
  }
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .post-card {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border: none;
    border-radius: 10px;
    overflow: hidden;
    padding: 16px;
    background: rgba(255, 255, 255, 0.98);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
      background: rgba(255, 255, 255, 1);

      .post-cover {
        transform: perspective(1000px) rotateY(-2deg) translateZ(10px);
        box-shadow: 6px 6px 16px rgba(0, 0, 0, 0.15);

        img {
          transform: scale(1.05);
        }
      }
    }

    .post-cover {
      height: 220px;
      overflow: hidden;
      border-radius: 6px;
      margin: -16px -16px 16px;
      position: relative;
      transform: perspective(1000px) rotateY(0deg) translateZ(0);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);

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

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
      }
    }

    .post-content {
      .post-title {
        margin: 0 0 12px;
        font-size: 1.15em;
        font-weight: 600;
        line-height: 1.4;

        a {
          color: #2c3e50;
          text-decoration: none;
          transition: all 0.3s ease;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;

          &:hover {
            color: #2B5876;
            text-shadow: 0 0 1px rgba(43, 88, 118, 0.15);
          }
        }
      }

      .post-excerpt {
        margin: 0 0 14px;
        color: #606266;
        font-size: 0.9em;
        line-height: 1.6;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        opacity: 0.9;
      }

      .post-meta {
        display: flex;
        gap: 12px;
        color: #909399;
        font-size: 0.82em;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
          padding: 2px 0;

          .el-icon {
            font-size: 1.1em;
            opacity: 0.85;
          }

          &.post-category {
            color: #2B5876;
            font-weight: 500;
            background: rgba(43, 88, 118, 0.08);
            padding: 2px 8px;
            border-radius: 4px;
            transition: all 0.3s ease;

            &:hover {
              background: rgba(43, 88, 118, 0.12);
            }
          }
        }
      }
    }
  }
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .widget {
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
    background: rgba(255, 255, 255, 0.98);
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
      background: rgba(255, 255, 255, 1);
    }

    .widget-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 14px 16px;
      border-bottom: 1px solid rgba(0, 0, 0, 0.04);
      background: linear-gradient(to right, rgba(43, 88, 118, 0.05), rgba(43, 88, 118, 0.02));

      h3 {
        margin: 0;
        font-size: 1em;
        font-weight: 600;
        color: #2c3e50;
      }

      .view-all {
        display: flex;
        align-items: center;
        gap: 4px;
        color: #2B5876;
        text-decoration: none;
        font-size: 0.85em;
        padding: 4px 8px;
        border-radius: 4px;
        transition: all 0.3s ease;
        background: rgba(43, 88, 118, 0.06);

        &:hover {
          background: rgba(43, 88, 118, 0.1);
          transform: translateX(2px);
        }
      }
    }
  }

  .categories-list {
    list-style: none;
    padding: 12px 16px;
    margin: 0;

    .category-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 10px;
      color: #606266;
      text-decoration: none;
      transition: all 0.3s ease;
      border-radius: 6px;

      &:not(:last-child) {
        border-bottom: 1px solid rgba(0, 0, 0, 0.04);
      }

      &:hover {
        color: var(--el-color-primary);
        background: rgba(64, 158, 255, 0.06);
        padding-left: 14px;
      }

      .category-count {
        color: #909399;
        font-size: 0.85em;
        background: rgba(0, 0, 0, 0.04);
        padding: 2px 8px;
        border-radius: 10px;
      }
    }
  }

  .tags-cloud {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding: 12px 16px;

    .tag-item {
      color: #4a5568;
      text-decoration: none;
      padding: 4px 10px;
      border-radius: 4px;
      background: rgba(43, 88, 118, 0.06);
      transition: all 0.3s ease;
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 0.9em;

      &:hover {
        color: #fff;
        background: linear-gradient(135deg, #2B5876, #4E4376);
        transform: translateY(-1px);
      }

      .tag-count {
        font-size: 0.85em;
        opacity: 0.8;
        background: rgba(255, 255, 255, 0.2);
        padding: 1px 6px;
        border-radius: 8px;
      }
    }
  }
}
</style> 