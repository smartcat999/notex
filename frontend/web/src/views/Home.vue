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
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 1.5em;
      font-weight: 600;
    }

    .view-all {
      display: flex;
      align-items: center;
      gap: 4px;
      color: #409eff;
      text-decoration: none;
      font-size: 0.9em;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .post-card {
    transition: transform 0.3s ease;

    &:hover {
      transform: translateY(-2px);
    }

    .post-cover {
      height: 200px;
      overflow: hidden;
      border-radius: 4px;
      margin: -20px -20px 16px;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .post-content {
      .post-title {
        margin: 0 0 12px;
        font-size: 1.25em;
        font-weight: 600;
        color: #303133;
      }

      .post-excerpt {
        margin: 0 0 16px;
        color: #606266;
        font-size: 0.95em;
        line-height: 1.6;
      }

      .post-meta {
        display: flex;
        gap: 16px;
        color: #909399;
        font-size: 0.85em;

        span {
          display: flex;
          align-items: center;
          gap: 4px;

          .el-icon {
            font-size: 1.1em;
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
    .widget-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      h3 {
        margin: 0;
        font-size: 1.1em;
        font-weight: 600;
      }

      .view-all {
        display: flex;
        align-items: center;
        gap: 4px;
        color: #409eff;
        text-decoration: none;
        font-size: 0.85em;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }

  .categories-list {
    list-style: none;
    padding: 0;
    margin: 0;

    .category-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;
      color: #606266;
      text-decoration: none;
      transition: all 0.3s ease;

      &:not(:last-child) {
        border-bottom: 1px solid #ebeef5;
      }

      &:hover {
        color: #409eff;
        padding-left: 8px;
      }

      .category-count {
        color: #909399;
        font-size: 0.9em;
      }
    }
  }

  .tags-cloud {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    .tag-item {
      color: #606266;
      text-decoration: none;
      padding: 4px 8px;
      border-radius: 4px;
      background: #f4f4f5;
      transition: all 0.3s ease;
      display: flex;
      align-items: center;
      gap: 4px;

      &:hover {
        color: #fff;
        background: #409eff;
      }

      .tag-count {
        font-size: 0.9em;
        opacity: 0.8;
      }
    }
  }
}
</style> 