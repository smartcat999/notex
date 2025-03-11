<template>
  <div class="posts-page">
    <div class="page-header">
      <div class="title-section">
        <div class="title-wrapper">
          <h1>学习笔记</h1>
          <div class="title-decoration"></div>
        </div>
        <p class="subtitle">记录点滴，分享所得</p>
      </div>
      <div class="header-actions">
        <el-input
          v-model="searchQuery"
          placeholder="搜索文章..."
          class="search-input"
          clearable
          @clear="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
          <template #append>
            <el-button @click="handleSearch">搜索</el-button>
          </template>
        </el-input>

        <el-select
          v-model="selectedCategory"
          placeholder="选择分类"
          clearable
          @change="handleSearch"
        >
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>

        <el-select
          v-model="sortBy"
          placeholder="排序方式"
          @change="handleSearch"
        >
          <el-option label="最新发布" value="newest" />
          <el-option label="最多浏览" value="most_viewed" />
          <el-option label="最多评论" value="most_commented" />
        </el-select>
      </div>
    </div>

    <div class="posts-container">
      <el-row :gutter="24">
        <el-col :span="18">
          <div class="posts-list" v-loading="loading">
            <el-card v-for="post in posts" :key="post.id" class="post-card">
              <div class="post-content">
                <div class="post-main">
                  <h2 class="post-title">
                    <router-link :to="`/posts/${post.id}`">
                      {{ post.title }}
                    </router-link>
                  </h2>
                  <p class="post-summary">{{ post.summary }}</p>
                  <div class="post-meta">
                    <span class="post-date">
                      <el-icon><Calendar /></el-icon>
                      {{ formatDate(post.published_at || post.created_at) }}
                    </span>
                    <span class="post-category" v-if="post.category">
                      <el-icon><Folder /></el-icon>
                      {{ post.category }}
                    </span>
                    <span class="post-views" v-if="post.views !== undefined">
                      <el-icon><View /></el-icon>
                      {{ post.views || 0 }} 次浏览
                    </span>
                  </div>
                </div>
                <div class="post-cover" v-if="post.cover">
                  <img :src="post.cover" :alt="post.title">
                </div>
              </div>
              <div class="post-footer">
                <div class="post-tags" v-if="post.tags && post.tags.length">
                  <el-tag
                    v-for="tag in post.tags"
                    :key="tag.id"
                    size="small"
                    :type="getRandomTagType()"
                    effect="plain"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
              </div>
            </el-card>

            <div class="pagination" v-if="total > 0">
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :total="total"
                :page-sizes="[10, 20, 30, 50]"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handlePageChange"
              />
            </div>

            <el-empty
              v-if="!loading && posts.length === 0"
              description="暂无文章"
            />
          </div>
        </el-col>

        <el-col :span="6">
          <div class="sidebar">
            <el-card class="widget tags-widget">
              <template #header>
                <div class="widget-header">
                  <h3>热门标签</h3>
                </div>
              </template>
              <div class="tags-cloud">
                <el-tag
                  v-for="tag in popularTags"
                  :key="tag.id"
                  :type="tag.type"
                  :effect="selectedTag === tag.id ? 'dark' : 'plain'"
                  class="tag-item"
                  @click="handleTagClick(tag)"
                >
                  {{ tag.name }}
                </el-tag>
              </div>
            </el-card>

            <PostArchives />
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { 
  Search,
  Calendar,
  Folder,
  View,
  ChatDotRound,
  Plus
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { getPosts, getCategories, getTopTags } from '@/api/posts'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import PostArchives from '@/components/PostArchives.vue'

const posts = ref([])
const categories = ref([])
const popularTags = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const selectedCategory = ref('')
const sortBy = ref('newest')
const loading = ref(false)
const selectedTag = ref('')

const route = useRoute()
const router = useRouter()

// 获取文章列表
const fetchPosts = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value,
      category_id: selectedCategory.value,
      tag_id: selectedTag.value,
      sort: sortBy.value,
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

// 获取分类列表
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

// 获取热门标签
const fetchPopularTags = async () => {
  try {
    const response = await getTopTags({ limit: 20 })
    popularTags.value = response.items.map(tag => ({
      ...tag,
      type: getRandomTagType()
    })) || []
  } catch (error) {
    console.error('Failed to fetch popular tags:', error)
    popularTags.value = []
  }
}

// 获取随机标签类型
const getRandomTagType = () => {
  const types = ['primary', 'success', 'warning', 'danger', 'info']
  return types[Math.floor(Math.random() * types.length)]
}

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchPosts()
}

// 处理标签点击
const handleTagClick = (tag) => {
  // 如果当前标签已经选中，则取消选中
  if (selectedTag.value === tag.id) {
    selectedTag.value = ''
  } else {
    selectedTag.value = tag.id
  }
  currentPage.value = 1
  // 更新 URL 参数
  const query = { ...route.query }
  if (selectedTag.value) {
    query.tag_id = selectedTag.value.toString()
  } else {
    delete query.tag_id
  }
  router.push({ query })
}

// 处理分页大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchPosts()
}

// 处理页码变化
const handlePageChange = (val) => {
  currentPage.value = val
  fetchPosts()
}

const init = async () => {
  loading.value = true
  try {
    // 从 URL 参数中获取分类和标签 ID
    const categoryId = route.query.category_id
    const tagId = route.query.tag_id
    
    if (categoryId) {
      selectedCategory.value = parseInt(categoryId)
    }
    if (tagId) {
      selectedTag.value = parseInt(tagId)
    }

    // 并行加载所有数据
    const [categoriesRes, tagsRes, postsRes] = await Promise.all([
      getCategories(),
      getTopTags({ limit: 20 }),
      getPosts({
        page: currentPage.value,
        page_size: pageSize.value,
        search: searchQuery.value,
        category_id: selectedCategory.value,
        tag_id: selectedTag.value,
        sort: sortBy.value,
      })
    ])

    // 处理返回的数据
    if (categoriesRes) {
      categories.value = categoriesRes.items || []
    }
    if (tagsRes) {
      popularTags.value = tagsRes.items.map(tag => ({
        ...tag,
        type: getRandomTagType()
      })) || []
    }
    if (postsRes) {
      posts.value = postsRes.items || []
      total.value = postsRes.total || 0
    }
  } catch (error) {
    console.error('Failed to initialize page:', error)
    ElMessage.error('页面初始化失败')
    categories.value = []
    popularTags.value = []
    posts.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 监听路由参数变化
watch(
  () => route.query,
  (newQuery) => {
    const categoryId = newQuery.category_id
    const tagId = newQuery.tag_id
    
    if (categoryId) {
      selectedCategory.value = parseInt(categoryId)
    } else {
      selectedCategory.value = ''
    }
    
    if (tagId) {
      selectedTag.value = parseInt(tagId)
    } else {
      selectedTag.value = ''
    }
    
    currentPage.value = 1
    fetchPosts()
  }
)

onMounted(() => {
  init()
})
</script>

<style scoped lang="scss">
.posts-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .page-header {
    margin-bottom: 32px;
    padding: 0;
    max-width: 900px;
    margin-left: auto;
    margin-right: auto;

    .title-section {
      text-align: center;
      margin-bottom: 24px;

      .title-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 8px;

        h1 {
          font-size: 1.5em;
          font-weight: 600;
          color: #2c3e50;
          margin: 0;
          position: relative;
          display: inline-block;
          background: linear-gradient(135deg, #2B5876, #4E4376);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
        }
      }

      .subtitle {
        color: #606266;
        font-size: 0.95em;
        margin: 0;
        opacity: 0.8;
      }
    }

    .header-actions {
      display: flex;
      gap: 16px;
      justify-content: center;
      align-items: center;

      .search-input {
        width: 300px;
      }

      .el-select {
        width: 140px;
      }
    }
  }

  .posts-container {
    max-width: 900px;
    margin: 0 auto;

    .posts-list {
      .post-card {
        margin-bottom: 32px;
        padding: 32px;
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
          }
        }

        .post-content {
          display: flex;
          align-items: flex-start;
          gap: 32px;

          .post-main {
            flex: 1;
            min-width: 0;

            .post-title {
              margin: 0 0 20px;
              font-size: 1.25em;
              font-weight: 600;
              line-height: 1.5;

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

            .post-summary {
              margin: 0 0 24px;
              color: #606266;
              font-size: 1em;
              line-height: 1.8;
              display: -webkit-box;
              -webkit-line-clamp: 4;
              -webkit-box-orient: vertical;
              overflow: hidden;
              opacity: 0.9;
            }

            .post-meta {
              margin-top: 16px;
              display: flex;
              gap: 12px;
              align-items: center;
              flex-wrap: wrap;
              font-size: 0.82em;
              color: #909399;

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

          .post-cover {
            flex: 0 0 240px;
            height: 160px;
            border-radius: 6px;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
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

            img {
              width: 100%;
              height: 100%;
              object-fit: cover;
              transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
            }

            &:hover {
              &::before {
                opacity: 0.7;
              }
              
              img {
                transform: scale(1.05);
              }
            }
          }
        }

        .post-footer {
          margin-top: 24px;
          padding-top: 24px;
          border-top: 1px solid rgba(0, 0, 0, 0.04);

          .post-tags {
            display: flex;
            gap: 6px;
            flex-wrap: wrap;

            .el-tag {
              border-radius: 4px;
              font-size: 0.8em;
              padding: 0 8px;
              height: 22px;
              line-height: 20px;
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
        }
      }

      .pagination {
        margin-top: 32px;
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
    }
  }

  .sidebar {
    display: flex;
    flex-direction: column;
    gap: 24px;
    position: sticky;
    top: 24px;
    height: fit-content;
    transition: all 0.3s ease;
    will-change: transform;
    z-index: 1;

    &:hover {
      transform: translateY(-2px);
    }

    .widget {
      border-radius: 10px;
      border-radius: 16px;
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
      background: rgba(255, 255, 255, 0.98);
      overflow: hidden;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
        background: rgba(255, 255, 255, 1);
      }

      .widget-header {
        padding: 16px 20px;
        border-bottom: 1px solid rgba(0, 0, 0, 0.06);
        background: linear-gradient(to right, rgba(43, 88, 118, 0.05), rgba(43, 88, 118, 0.02));

        h3 {
          margin: 0;
          font-size: 1.05em;
          font-weight: 600;
          color: #2c3e50;
        }
      }
    }

    .tags-cloud {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      padding: 16px;

      .tag-item {
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        border: none;
        box-shadow: 0 2px 6px rgba(43, 88, 118, 0.1);
        border-radius: 6px;

        &:hover {
          transform: translateY(-2px) scale(1.05);
          box-shadow: 0 4px 8px rgba(43, 88, 118, 0.15);
          background: linear-gradient(135deg, #2B5876, #4E4376);
          color: white;
        }
      }
    }
  }
}

// 响应式布局优化
@media (max-width: 768px) {
  .posts-page {
    margin: 20px auto;
    
    .page-header {
      margin-bottom: 24px;
      padding: 16px 0;
      
      .title-section {
        .title-wrapper h1 {
          font-size: 2em;
        }
        
        .subtitle {
          font-size: 1em;
        }
      }
      
      .header-actions {
        flex-direction: column;
        gap: 12px;
        
        .search-input {
          width: 100%;
        }
        
        .el-select {
          width: 100%;
        }
      }
    }
    
    .posts-list {
      .post-card {
        padding: 20px;
        
        .post-content {
          flex-direction: column;
          gap: 20px;

          .post-cover {
            order: -1;
            flex: none;
            width: calc(100% + 40px);
            height: 200px;
            margin: -20px -20px 0;
            border-radius: 16px 16px 0 0;
            box-shadow: none;
            
            &:hover {
              transform: none;
              box-shadow: none;
            }
          }
          
          .post-main {
            .post-title {
              font-size: 1.2em;
              margin-bottom: 12px;
            }
            
            .post-summary {
              font-size: 0.95em;
              margin-bottom: 12px;
            }
            
            .post-meta {
              gap: 16px;
              font-size: 0.9em;
            }
          }
        }
      }
    }
  }
}
</style> 