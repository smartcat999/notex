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

    <div class="post-content markdown-body" ref="contentRef">
      <div v-html="renderedContent"></div>
      <!-- 复制按钮组件 -->
      <button
        v-for="(code, index) in codeSections"
        :key="index"
        class="copy-button"
        :style="{
          top: `${code.top}px`,
          left: `${code.left}px`
        }"
        @click="handleCopy(code.text)"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
        </svg>
        复制
      </button>
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
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { Calendar, View, ChatDotRound, User } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'
import { useUserStore } from '@/stores/user'
import { getPost, getComments, createComment } from '@/api/posts'
import { marked } from 'marked'
import hljs from 'highlight.js'
import { ElMessage } from 'element-plus'
import '@/assets/styles/markdown.scss'

const route = useRoute()
const userStore = useUserStore()
const post = ref({})
const comments = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const commentContent = ref('')
const submitting = ref(false)
const contentRef = ref(null)
const codeSections = ref([])

// 配置 marked
marked.setOptions({
  highlight: function(code, lang) {
    const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext';
    const html = hljs.highlight(code, { language }).value;
    return `<pre data-language="${language}"><code>${html}</code></pre>`;
  },
  breaks: true,
  gfm: true,
  headerIds: true,
  langPrefix: 'hljs language-'
})

// 复制代码功能
const handleCopy = (code) => {
  navigator.clipboard.writeText(code).then(() => {
    ElMessage({
      message: '复制成功',
      type: 'success',
      duration: 2000,
      offset: 80
    })
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 更新代码块位置
const updateCodeSections = async () => {
  await nextTick()
  if (!contentRef.value) return

  const pres = contentRef.value.querySelectorAll('pre')
  codeSections.value = Array.from(pres).map(pre => {
    const code = pre.querySelector('code')
    const rect = pre.getBoundingClientRect()
    const contentRect = contentRef.value.getBoundingClientRect()
    return {
      text: code.textContent,
      top: rect.top - contentRect.top + 8,
      left: rect.right - contentRect.left - 80 // 调整按钮的水平位置
    }
  })
}

// 渲染 Markdown 内容
const renderedContent = computed(() => {
  const content = marked(post.value.content || '')
  nextTick(() => {
    updateCodeSections()
  })
  return content
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
  window.addEventListener('resize', updateCodeSections)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateCodeSections)
})
</script>

<style lang="scss">
.post-detail-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 60px 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  color: #24292f;
  line-height: 1.6;

  .post-content {
    margin-bottom: 40px;
    background-color: #fff;
    padding: 40px;
    border-radius: 16px;
    box-shadow: 0 4px 30px rgba(0, 0, 0, 0.04);
    position: relative;

    pre {
      position: relative;
      margin: 1em 0;
    }

    .copy-button {
      position: absolute;
      padding: 0.4em 0.8em;
      font-size: 0.8em;
      color: #1a1a1a;
      background: #fff;
      border: 1px solid rgba(0, 0, 0, 0.1);
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s ease;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      display: flex;
      align-items: center;
      gap: 4px;
      z-index: 10;
      font-weight: 500;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      
      svg {
        width: 12px;
        height: 12px;
        stroke: currentColor;
      }

      &:hover {
        background: #f0f0f0;
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }

      &:active {
        transform: translateY(0);
        background: #e0e0e0;
      }
    }
  }
}

.post-header {
  margin-bottom: 40px;
  text-align: center;
  position: relative;
  padding: 30px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.04);

  &::after {
    content: '';
    position: absolute;
    bottom: -2px;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, 
      var(--el-color-primary) 0%,
      var(--el-color-primary-light-3) 50%,
      var(--el-color-primary-light-5) 100%
    );
    border-radius: 0 0 16px 16px;
    opacity: 0.6;
  }

  h1 {
    font-size: 2em;
    margin-bottom: 20px;
    color: #1f2937;
    font-weight: 700;
    line-height: 1.3;
    letter-spacing: -0.02em;
    background: linear-gradient(120deg, #1f2937 0%, #4b5563 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .post-meta {
    color: #6b7280;
    font-size: 0.9em;
    margin-bottom: 20px;
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 16px;

    span {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 4px 10px;
      background: #f9fafb;
      border-radius: 20px;
      transition: all 0.3s ease;

      &:hover {
        background: #f3f4f6;
        transform: translateY(-1px);
      }

      .el-icon {
        font-size: 1em;
        color: var(--el-color-primary);
      }
    }
  }

  .post-tags {
    display: flex;
    justify-content: center;
    gap: 8px;
    flex-wrap: wrap;

    .tag {
      margin: 0;
      padding: 4px 12px;
      font-size: 0.85em;
      font-weight: 500;
      border-radius: 20px;
      transition: all 0.3s ease;
      cursor: pointer;
      border: none;
      background: var(--el-color-primary-light-9);
      color: var(--el-color-primary-dark-2);

      &:hover {
        transform: translateY(-2px);
        background: var(--el-color-primary-light-8);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
      }
    }
  }
}

.comments-section {
  background: #fff;
  padding: 25px;
  border-radius: 16px;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.04);

  h2 {
    margin-bottom: 20px;
    font-size: 1.3em;
    font-weight: 700;
    color: #1f2937;
    position: relative;
    padding-left: 15px;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 4px;
      height: 20px;
      background: linear-gradient(180deg, var(--el-color-primary), var(--el-color-primary-light-5));
      border-radius: 3px;
    }
  }

  .comment-form {
    margin-bottom: 25px;
    background: #f9fafb;
    padding: 15px;
    border-radius: 12px;
    border: 1px solid #e5e7eb;

    .el-input__inner {
      border: 2px solid #e5e7eb;
      transition: all 0.3s ease;
      border-radius: 8px;
      padding: 8px;
      font-size: 0.95em;

      &:focus {
        border-color: var(--el-color-primary);
        box-shadow: 0 0 0 3px var(--el-color-primary-light-8);
      }
    }

    .el-button {
      margin-top: 15px;
      float: right;
      padding: 8px 16px;
      font-size: 0.95em;
      border-radius: 8px;
    }
  }

  .login-prompt {
    text-align: center;
    padding: 30px;
    background: linear-gradient(120deg, #f9fafb, #f3f4f6);
    border-radius: 12px;
    margin-bottom: 30px;
    border: 2px dashed #e5e7eb;

    p {
      font-size: 1em;
      color: #4b5563;
    }

    a {
      color: var(--el-color-primary);
      text-decoration: none;
      font-weight: 600;
      border-bottom: 2px solid transparent;
      transition: all 0.3s ease;
      padding: 2px 4px;
      margin: 0 2px;
      border-radius: 4px;

      &:hover {
        background: var(--el-color-primary-light-9);
        border-bottom-color: var(--el-color-primary);
      }
    }
  }

  .comments-list {
    .comment-item {
      margin-bottom: 20px;
      padding-bottom: 20px;
      border-bottom: 1px solid #e5e7eb;
      transition: all 0.3s ease;

      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
        padding-bottom: 0;
      }

      &:hover {
        transform: translateX(8px);
      }

      .comment-header {
        display: flex;
        align-items: center;
        margin-bottom: 10px;

        .el-avatar {
          border: 3px solid #fff;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
          transition: all 0.3s ease;
          width: 36px;
          height: 36px;

          &:hover {
            transform: scale(1.1);
          }
        }

        .comment-info {
          margin-left: 15px;

          .username {
            font-weight: 600;
            color: #1f2937;
            margin-right: 16px;
            font-size: 0.9em;
          }

          .time {
            color: #6b7280;
            font-size: 0.8em;
          }
        }
      }

      .comment-content {
        color: #4b5563;
        line-height: 1.6;
        padding-left: 51px;
        font-size: 0.9em;
      }
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 30px;

    :deep(.el-pagination) {
      --el-pagination-font-size: 1.05em;
      --el-pagination-button-height: 36px;
      --el-pagination-button-width: 36px;

      .el-pagination__total,
      .el-pagination__jump {
        font-size: 1em;
      }

      .btn-prev,
      .btn-next {
        background: #f9fafb;
        border-radius: 8px;
        margin: 0 4px;
      }

      .el-pager li {
        background: #f9fafb;
        border-radius: 8px;
        margin: 0 4px;
        min-width: 36px;

        &.active {
          background: var(--el-color-primary);
        }

        &:hover {
          color: var(--el-color-primary);
        }
      }
    }
  }
}
</style> 