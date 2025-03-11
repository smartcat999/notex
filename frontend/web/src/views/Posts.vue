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
              <div class="post-cover" v-if="post.cover">
                <img :src="post.cover" :alt="post.title">
              </div>
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
  margin: 24px auto;
  padding: 0 20px;

  .page-header {
    margin-bottom: 36px;
    padding: 20px 0;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(to right, rgba(64, 158, 255, 0), rgba(64, 158, 255, 0.5), rgba(64, 158, 255, 0));
    }

    .title-section {
      margin-bottom: 28px;
      text-align: center;
      
      .title-wrapper {
        display: inline-block;
        position: relative;
        margin-bottom: 12px;

        h1 {
          margin: 0;
          font-size: 2.4em;
          font-weight: 700;
          color: #2c3e50;
          letter-spacing: 0.05em;
          position: relative;
          z-index: 1;
        }

        .title-decoration {
          position: absolute;
          bottom: 8px;
          left: -8px;
          right: -8px;
          height: 12px;
          background-color: rgba(64, 158, 255, 0.1);
          z-index: 0;
        }
      }

      .subtitle {
        margin: 0;
        font-size: 1.15em;
        color: #5c6b77;
        font-weight: 300;
        letter-spacing: 0.03em;
      }
    }

    .header-actions {
      display: flex;
      justify-content: center;
      gap: 16px;
      margin-top: 24px;

      .search-input {
        width: 360px;
      }

      .el-select {
        width: 140px;
      }
    }
  }

  .posts-container {
    .posts-list {
      display: flex;
      flex-direction: column;
      gap: 20px;

      .post-card {
        transition: transform 0.3s ease, box-shadow 0.3s ease;
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }

        .post-cover {
          height: 200px;
          overflow: hidden;
          border-radius: 8px;
          margin: -20px -20px 16px;

          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: transform 0.3s ease;

            &:hover {
              transform: scale(1.05);
            }
          }
        }

        .post-content {
          .post-main {
            .post-title {
              margin: 0 0 12px;
              font-size: 1.4em;
              font-weight: 600;
              line-height: 1.4;

              a {
                color: #2c3e50;
                text-decoration: none;
                transition: color 0.3s ease;

                &:hover {
                  color: #409eff;
                }
              }
            }

            .post-summary {
              margin: 0 0 16px;
              color: #5c6b77;
              font-size: 0.95em;
              line-height: 1.6;
              display: -webkit-box;
              -webkit-line-clamp: 3;
              -webkit-box-orient: vertical;
              overflow: hidden;
            }

            .post-meta {
              display: flex;
              gap: 16px;
              align-items: center;
              margin-bottom: 16px;
              font-size: 0.9em;
              color: #8c9ba5;

              span {
                display: flex;
                align-items: center;
                gap: 4px;

                .el-icon {
                  font-size: 1.1em;
                }

                &.post-category {
                  color: #409eff;
                  font-weight: 500;
                }
              }
            }
          }

          .post-footer {
            margin-top: 16px;
            padding-top: 16px;
            border-top: 1px solid #ebeef5;

            .post-tags {
              display: flex;
              flex-wrap: wrap;
              gap: 8px;

              .el-tag {
                border-radius: 4px;
                cursor: pointer;
                transition: all 0.3s ease;

                &:hover {
                  transform: translateY(-1px);
                }
              }
            }
          }
        }
      }

      .pagination {
        margin-top: 24px;
        display: flex;
        justify-content: center;
      }
    }
  }

  .sidebar {
    display: flex;
    flex-direction: column;
    gap: 20px;
    position: sticky;
    top: 20px;
    height: fit-content;
    transition: transform 0.2s ease;
    will-change: transform;
    z-index: 1;

    &:hover {
      transform: translateY(-2px);
    }

    .widget {
      .widget-header {
        h3 {
          margin: 0;
          font-size: 1.1em;
          font-weight: 600;
        }
      }
    }

    .tags-cloud {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .tag-item {
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          transform: translateY(-1px);
        }
      }
    }

    .archive-list {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        a {
          display: flex;
          justify-content: space-between;
          padding: 8px 0;
          color: #606266;
          text-decoration: none;
          transition: all 0.3s ease;

          &:hover {
            color: #409eff;
            padding-left: 8px;
          }

          .count {
            color: #909399;
          }
        }

        &:not(:last-child) {
          border-bottom: 1px solid #ebeef5;
        }
      }
    }
  }
}
</style> 