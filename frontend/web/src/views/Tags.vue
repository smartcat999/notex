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
  margin: 32px auto;
  padding: 0 20px;

  h1 {
    margin-bottom: 30px;
    font-size: 1.8em;
    font-weight: 600;
    text-align: center;
    background: linear-gradient(120deg, #2c3e50, #3498db);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: 0.02em;
    position: relative;
    display: inline-block;
    left: 50%;
    transform: translateX(-50%);

    &::after {
      content: '';
      position: absolute;
      bottom: 4px;
      left: -8px;
      right: -8px;
      height: 12px;
      background-color: rgba(64, 158, 255, 0.15);
      border-radius: 6px;
      z-index: -1;
      transform: skew(-12deg);
    }
  }
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 40px;
  justify-content: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
    background: rgba(255, 255, 255, 1);
  }

  .tag {
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    padding: 8px 16px;
    font-size: 1em;
    border: none;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
    border-radius: 8px;

    &:hover {
      transform: translateY(-2px) scale(1.05);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
    }

    &.active {
      transform: scale(1.1);
      box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
    }

    .tag-count {
      margin-left: 4px;
      font-size: 0.85em;
      opacity: 0.8;
      background: rgba(0, 0, 0, 0.1);
      padding: 2px 8px;
      border-radius: 10px;
    }
  }
}

.tag-posts {
  h2 {
    margin-bottom: 20px;
    font-size: 1.4em;
    font-weight: 600;
    background: linear-gradient(120deg, #2c3e50, #3498db);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    padding: 16px 20px;
    border-radius: 10px;
    background-color: rgba(255, 255, 255, 0.98);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
      background-color: rgba(255, 255, 255, 1);
    }
  }

  .posts-list {
    .post-card {
      margin-bottom: 20px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      border: none;
      border-radius: 10px;
      background: rgba(255, 255, 255, 0.98);
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
        background: rgba(255, 255, 255, 1);

        .post-cover {
          transform: perspective(1000px) rotateY(-2deg) translateZ(10px);
          box-shadow: 6px 6px 16px rgba(0, 0, 0, 0.15);

          .cover-image {
            transform: scale(1.05);
          }
        }
      }

      .post-content {
        display: flex;
        gap: 20px;

        .post-main {
          flex: 1;

          h3 {
            margin: 0 0 12px;
            font-size: 1.15em;

            a {
              color: inherit;
              text-decoration: none;
              transition: color 0.3s ease;

              &:hover {
                color: var(--el-color-primary);
              }
            }
          }

          .post-summary {
            color: #666;
            margin-bottom: 12px;
            font-size: 0.9em;
            line-height: 1.6;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
          }

          .post-meta {
            display: flex;
            gap: 16px;
            color: #999;
            font-size: 0.82em;

            span {
              display: flex;
              align-items: center;
              gap: 4px;
              padding: 2px 8px;
              border-radius: 6px;
              background: rgba(0, 0, 0, 0.04);
              transition: all 0.3s ease;

              &:hover {
                background: rgba(0, 0, 0, 0.06);
              }

              .el-icon {
                font-size: 1.1em;
                opacity: 0.85;
              }
            }
          }
        }

        .post-cover {
          width: 200px;
          height: 120px;
          overflow: hidden;
          border-radius: 8px;
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

      .post-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 16px;
        padding-top: 16px;
        border-top: 1px solid rgba(0, 0, 0, 0.04);

        .post-tags {
          display: flex;
          gap: 8px;

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

        .el-button {
          font-size: 0.85em;
          background: rgba(64, 158, 255, 0.08);
          padding: 4px 12px;
          border-radius: 6px;
          transition: all 0.3s ease;

          &:hover {
            background: rgba(64, 158, 255, 0.12);
            transform: translateX(2px);
          }
        }
      }
    }
  }
}

.pagination {
  margin-top: 30px;
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
  .tags-container {
    padding: 0 16px;

    h1 {
      font-size: 1.6em;
    }
  }

  .tag-posts {
    .posts-list {
      .post-card {
        .post-content {
          flex-direction: column;

          .post-cover {
            order: -1;
            width: 100%;
            height: 160px;
            margin: -16px -16px 16px;
            border-radius: 10px 10px 0 0;
          }
        }
      }
    }
  }
}
</style> 