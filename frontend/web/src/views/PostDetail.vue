<template>
  <div class="post-detail-container">
    <div class="post-header">
      <h1>{{ post.title }}</h1>
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
        <span>
          <el-icon><User /></el-icon>
          {{ post.author?.username }}
        </span>
      </div>
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
    </div>

    <div class="post-content">
      <div class="markdown-body" v-html="renderedContent"></div>
    </div>

    <!-- 评论区 -->
    <div class="comments-section">
      <h2>评论 ({{ post.comment_count }})</h2>
      
      <!-- 评论表单 -->
      <div class="comment-form" v-if="userStore.isAuthenticated">
        <el-input
          v-model="commentContent"
          type="textarea"
          :rows="3"
          placeholder="写下你的评论..."
          :maxlength="500"
          show-word-limit
        />
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmitComment"
        >
          发表评论
        </el-button>
      </div>
      <div v-else class="login-prompt">
        <p>
          请
          <router-link to="/login">登录</router-link>
          后发表评论
        </p>
      </div>

      <!-- 评论列表 -->
      <div class="comments-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <div class="comment-header">
            <el-avatar :size="40" :src="comment.user?.avatar">
              {{ comment.user?.username?.charAt(0) }}
            </el-avatar>
            <div class="comment-info">
              <span class="username">{{ comment.user?.username }}</span>
              <span class="time">{{ formatDate(comment.created_at) }}</span>
            </div>
          </div>
          <div class="comment-content">
            {{ comment.content }}
          </div>
        </div>
      </div>

      <!-- 评论分页 -->
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
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Calendar, View, ChatDotRound, User } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { useUserStore } from '@/stores/user'
import { getPost, getComments, createComment } from '@/api/posts'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import { ElMessage } from 'element-plus'

const route = useRoute()
const userStore = useUserStore()
const post = ref({})
const comments = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const commentContent = ref('')
const submitting = ref(false)

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch (__) {}
    }
    return ''
  },
})

const renderedContent = computed(() => {
  return md.render(post.value.content || '')
})

const fetchPost = async () => {
  try {
    const response = await getPost(route.params.id)
    post.value = response || {}
  } catch (error) {
    console.error('Failed to fetch post:', error)
    ElMessage.error('获取文章详情失败')
  }
}

const fetchComments = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    const response = await getComments(route.params.id, params)
    comments.value = response.items || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch comments:', error)
    ElMessage.error('获取评论失败')
    comments.value = []
    total.value = 0
  }
}

const handleSubmitComment = async () => {
  if (!commentContent.value.trim()) {
    ElMessage.warning('评论内容不能为空')
    return
  }

  try {
    submitting.value = true
    const response = await createComment(route.params.id, {
      content: commentContent.value.trim()
    })
    ElMessage.success('评论发表成功')
    commentContent.value = ''
    // 重新获取评论列表
    currentPage.value = 1
    await fetchComments()
    // 更新文章的评论数
    post.value.comment_count = (post.value.comment_count || 0) + 1
  } catch (error) {
    console.error('Failed to create comment:', error)
    ElMessage.error('评论发表失败')
  } finally {
    submitting.value = false
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchComments()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchComments()
}

onMounted(() => {
  fetchPost()
  fetchComments()
})
</script>

<style lang="scss" scoped>
.post-detail-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.post-header {
  margin-bottom: 40px;
  text-align: center;

  h1 {
    font-size: 2em;
    margin-bottom: 16px;
  }

  .post-meta {
    color: var(--el-text-color-secondary);
    font-size: 0.9em;
    margin-bottom: 16px;

    span {
      display: inline-flex;
      align-items: center;
      margin: 0 12px;

      .el-icon {
        margin-right: 4px;
      }
    }
  }

  .post-tags {
    .tag {
      margin: 0 4px;
    }
  }
}

.post-content {
  margin-bottom: 40px;

  :deep(.markdown-body) {
    font-size: 1.1em;
    line-height: 1.8;
    color: var(--el-text-color-primary);

    h1, h2, h3, h4, h5, h6 {
      margin: 1.5em 0 1em;
      font-weight: 600;
    }

    p {
      margin-bottom: 1em;
    }

    code {
      background-color: var(--el-fill-color-light);
      padding: 0.2em 0.4em;
      border-radius: 3px;
      font-family: monospace;
    }

    pre {
      background-color: var(--el-fill-color-light);
      padding: 1em;
      border-radius: 4px;
      overflow-x: auto;
      margin: 1em 0;

      code {
        background-color: transparent;
        padding: 0;
      }
    }

    blockquote {
      margin: 1em 0;
      padding: 0.5em 1em;
      border-left: 4px solid var(--el-border-color);
      color: var(--el-text-color-secondary);
    }

    img {
      max-width: 100%;
      height: auto;
      margin: 1em 0;
      border-radius: 4px;
    }
  }
}

.comments-section {
  h2 {
    margin-bottom: 20px;
  }

  .comment-form {
    margin-bottom: 30px;

    .el-button {
      margin-top: 12px;
      float: right;
    }
  }

  .login-prompt {
    text-align: center;
    padding: 20px;
    background-color: var(--el-fill-color-light);
    border-radius: 4px;
    margin-bottom: 30px;

    a {
      color: var(--el-color-primary);
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .comments-list {
    .comment-item {
      margin-bottom: 20px;
      padding-bottom: 20px;
      border-bottom: 1px solid var(--el-border-color-light);

      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
        padding-bottom: 0;
      }

      .comment-header {
        display: flex;
        align-items: center;
        margin-bottom: 12px;

        .comment-info {
          margin-left: 12px;

          .username {
            font-weight: 500;
            color: var(--el-text-color-primary);
            margin-right: 8px;
          }

          .time {
            color: var(--el-text-color-secondary);
            font-size: 0.9em;
          }
        }
      }

      .comment-content {
        color: var(--el-text-color-regular);
        line-height: 1.6;
      }
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 30px;
  }
}
</style> 