<template>
  <div class="tags-container">
    <h1>文章标签</h1>

    <!-- 标签云 -->
    <div class="tags-cloud">
      <el-tag
        v-for="tag in tags"
        :key="tag.id"
        :type="getRandomTagType()"
        class="tag"
        :class="{ active: selectedTag?.id === tag.id }"
        effect="plain"
        @click="selectTag(tag)"
      >
        {{ tag.name }}
        <span class="tag-count">({{ tag.post_count }})</span>
      </el-tag>
    </div>

    <!-- 标签下的文章列表 -->
    <div class="tag-posts" v-if="selectedTag">
      <h2>标签：{{ selectedTag.name }}</h2>
      <div class="posts-list">
        <el-card v-for="post in posts" :key="post.id" class="post-card" shadow="hover">
          <div class="post-content">
            <div class="post-main">
              <h3>
                <router-link :to="`/posts/${post.id}`">{{ post.title }}</router-link>
              </h3>
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
            <div class="post-cover" v-if="post.cover">
              <el-image
                :src="post.cover"
                :alt="post.title"
                fit="cover"
                class="cover-image"
              />
            </div>
          </div>
          <div class="post-footer">
            <div class="post-tags">
              <el-tag
                v-for="tag in post.tags"
                :key="tag.id"
                size="small"
                class="tag"
              >
                {{ tag.name }}
              </el-tag>
            </div>
            <el-button
              type="primary"
              link
              @click="$router.push(`/posts/${post.id}`)"
            >
              阅读更多
            </el-button>
          </div>
        </el-card>
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
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Calendar, View, ChatDotRound } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { getTags, getPosts } from '@/api/posts'

const route = useRoute()
const router = useRouter()
const tags = ref([])
const selectedTag = ref(null)
const posts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const getRandomTagType = () => {
  const types = ['primary', 'success', 'warning', 'danger', 'info']
  return types[Math.floor(Math.random() * types.length)]
}

const fetchTags = async () => {
  try {
    const response = await getTags()
    tags.value = response.items
    
    // 如果 URL 中有 tag_id，选中对应的标签
    const tagId = Number(route.query.tag_id)
    if (tagId) {
      const tag = response.items.find(t => t.id === tagId)
      if (tag) {
        selectedTag.value = tag
        await fetchPosts()
      }
    }
  } catch (error) {
    console.error('Failed to fetch tags:', error)
  }
}

const fetchPosts = async () => {
  if (!selectedTag.value) return

  try {
    const response = await getPosts({
      page: currentPage.value,
      page_size: pageSize.value,
      tag_id: selectedTag.value.id
    })
    
    posts.value = response.items || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch posts:', error)
  }
}

const selectTag = (tag) => {
  selectedTag.value = tag
  currentPage.value = 1
  router.push({ query: { tag_id: tag.id.toString() } })
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

onMounted(fetchTags)

// 监听路由变化，处理标签切换
watch(
  () => route.query.tag_id,
  (newTagId) => {
    if (!newTagId) {
      selectedTag.value = null
      posts.value = []
      total.value = 0
      return
    }

    // 如果有标签 ID，找到对应的标签并加载文章
    const tagId = Number(newTagId)
    const tag = tags.value.find(t => t.id === tagId)
    if (tag) {
      selectedTag.value = tag
      fetchPosts()
    }
  }
)
</script>

<style lang="scss" scoped>
.tags-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  h1 {
    margin-bottom: 30px;
  }
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 40px;

  .tag {
    cursor: pointer;
    transition: all 0.3s;
    padding: 8px 16px;
    font-size: 1.1em;

    &:hover {
      transform: translateY(-2px);
    }

    &.active {
      transform: scale(1.1);
    }

    .tag-count {
      margin-left: 4px;
      font-size: 0.9em;
    }
  }
}

.tag-posts {
  h2 {
    margin-bottom: 20px;
  }

  .posts-list {
    .post-card {
      margin-bottom: 20px;

      .post-content {
        display: flex;
        gap: 20px;

        .post-main {
          flex: 1;

          h3 {
            margin: 0 0 12px;
            font-size: 1.3em;

            a {
              color: inherit;
              text-decoration: none;
              transition: color 0.3s;

              &:hover {
                color: var(--el-color-primary);
              }
            }
          }

          .post-summary {
            color: #666;
            margin-bottom: 12px;
          }

          .post-meta {
            display: flex;
            gap: 16px;
            color: #999;
            font-size: 0.9em;

            span {
              display: flex;
              align-items: center;
              gap: 4px;
            }
          }
        }

        .post-cover {
          width: 200px;
          height: 120px;
          overflow: hidden;
          border-radius: 4px;

          .cover-image {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }
        }
      }

      .post-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 16px;
        padding-top: 16px;
        border-top: 1px solid #eee;

        .post-tags {
          display: flex;
          gap: 8px;
        }
      }
    }
  }
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}
</style> 