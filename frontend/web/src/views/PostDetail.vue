<template>
  <div class="post-detail-container">
    <div class="post-header" :class="{ 'no-cover': !post.cover }">
      <div class="post-cover" v-if="post.cover">
        <img :src="post.cover" alt="Post Cover">
      </div>
      <div class="header-content">
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
        </div>
        <!-- 作者信息 -->
        <AuthorCard v-if="post.author" :author="post.author" />
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
    </div>

    <div class="post-content markdown-body" ref="contentRef" v-html="renderedContent"></div>

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
            <router-link 
              :to="`/users/${comment.user?.id}`" 
              class="avatar-link"
            >
              <el-avatar :size="40" :src="comment.user?.avatar">
                {{ comment.user?.username?.charAt(0) }}
              </el-avatar>
            </router-link>
            <div class="comment-info">
              <router-link 
                :to="`/users/${comment.user?.id}`" 
                class="username"
              >
                {{ comment.user?.username }}
              </router-link>
              <span class="author-tag" v-if="comment.user?.id === post.author?.id">作者</span>
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
import AuthorCard from '@/components/AuthorCard.vue'

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
const copiedIndex = ref(-1)

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

// 更新复制代码功能
const handleCopy = (code, index) => {
  navigator.clipboard.writeText(code).then(() => {
    copiedIndex.value = index
    setTimeout(() => {
      copiedIndex.value = -1
    }, 2000)
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
const updateCodeSections = () => {
  if (!contentRef.value) return

  const pres = contentRef.value.querySelectorAll('pre')
  pres.forEach(pre => {
    if (!pre.querySelector('.pre-wrapper')) {
      const code = pre.querySelector('code')
      
      // 创建包装容器
      const wrapper = document.createElement('div')
      wrapper.className = 'pre-wrapper'
      
      // 移动代码到包装容器
      if (code) {
        pre.removeChild(code)
        wrapper.appendChild(code)
      }
      pre.appendChild(wrapper)

      // 添加复制按钮
      if (!pre.querySelector('.copy-button')) {
        const button = document.createElement('button')
        button.className = 'copy-button'
        button.innerHTML = `
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
          复制
        `
        button.addEventListener('click', () => {
          const text = code.textContent
          navigator.clipboard.writeText(text).then(() => {
            button.innerHTML = `
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              已复制
            `
            button.classList.add('copied')
            setTimeout(() => {
              button.innerHTML = `
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                </svg>
                复制
              `
              button.classList.remove('copied')
            }, 2000)
          }).catch(() => {
            ElMessage.error('复制失败')
          })
        })
        pre.appendChild(button)
      }
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
  padding: 40px 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  color: #24292f;
  line-height: 1.6;

  .post-header {
    position: relative;
    margin-bottom: 32px;
    border-radius: 20px;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.98);
    box-shadow: 0 8px 30px rgba(43, 88, 118, 0.1);
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 12px 40px rgba(43, 88, 118, 0.15);
    }

    &.no-cover {
      background: linear-gradient(135deg, rgba(43, 88, 118, 0.02), rgba(78, 67, 118, 0.02));
      padding: 40px;
      
      .header-content {
        position: relative;
        padding: 0;
        background: none;

        h1 {
          font-size: 2em;
          line-height: 1.4;
          margin: 0 0 20px;
          letter-spacing: -0.02em;
        }

        .post-meta {
          margin-bottom: 24px;
        }
      }
    }

    .post-cover {
      width: 100%;
      height: 320px;
      position: relative;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.6s ease;
      }

      &::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 160px;
        background: linear-gradient(
          to bottom,
          rgba(255, 255, 255, 0),
          rgba(255, 255, 255, 0.9) 70%,
          rgba(255, 255, 255, 0.98)
        );
        pointer-events: none;
      }

      &:hover img {
        transform: scale(1.03);
      }
    }

    .header-content {
      padding: 0 40px 32px;
      display: flex;
      flex-direction: column;
      align-items: center;

      h1 {
        font-size: 1.75em;
        margin: 0 0 24px;
        font-weight: 600;
        line-height: 1.4;
        letter-spacing: -0.01em;
        background: linear-gradient(135deg, #2B5876 30%, #4E4376);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        text-align: center;
        width: 100%;
        max-width: 800px;
      }

      .post-meta {
        width: 100%;
        color: #6b7280;
        font-size: 0.9em;
        margin-bottom: 24px;
        display: flex;
        justify-content: center;
        flex-wrap: wrap;
        gap: 16px;

        span {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          padding: 6px 12px;
          background: rgba(43, 88, 118, 0.06);
          border-radius: 20px;
          transition: all 0.3s ease;

          &:hover {
            background: rgba(43, 88, 118, 0.1);
            transform: translateY(-1px);
          }

          .el-icon {
            font-size: 1.1em;
            color: #2B5876;
          }
        }
      }

      :deep(.author-card) {
        width: 100%;
        margin: 0 0 24px;
        display: flex;
        justify-content: center;
        border: none;
        padding: 0;

        .author-content {
          background: rgba(43, 88, 118, 0.03);
          padding: 12px 20px;
          border-radius: 16px;
          transition: all 0.3s ease;

          &:hover {
            background: rgba(43, 88, 118, 0.06);
            transform: translateY(-1px);
          }
        }
      }

      .post-tags {
        width: 100%;
        display: flex;
        justify-content: center;
        gap: 10px;
        flex-wrap: wrap;

        .tag {
          margin: 0;
          padding: 6px 16px;
          font-size: 0.9em;
          font-weight: 500;
          border-radius: 20px;
          transition: all 0.3s ease;
          cursor: pointer;
          border: none;
          background: rgba(43, 88, 118, 0.06);
          color: #2B5876;
          backdrop-filter: blur(4px);

          &:hover {
            transform: translateY(-2px);
            background: rgba(43, 88, 118, 0.1);
            box-shadow: 0 4px 12px rgba(43, 88, 118, 0.1);
          }
        }
      }
    }
  }

  .post-content.markdown-body {
    margin-bottom: 40px;
    background-color: #fff;
    padding: 40px;
    border-radius: 16px;
    box-shadow: 0 4px 30px rgba(0, 0, 0, 0.04);
    position: relative;
    font-size: 16px;
    line-height: 1.8;
    color: #2c3e50;

    h1, h2, h3, h4, h5, h6 {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
      margin: 2em 0 1em;
      font-weight: 600;
      line-height: 1.3;
      letter-spacing: -0.02em;

      &:first-child {
        margin-top: 0.5em;
      }
    }

    h1 {
      font-size: 2.4em;
      color: #1a365d;
      margin: 1em 0 0.8em;
      position: relative;
      font-weight: 700;

      &::after {
        content: '';
        position: absolute;
        bottom: -0.3em;
        left: 0;
        width: 80px;
        height: 4px;
        background: linear-gradient(90deg, #1a365d, rgba(26, 54, 93, 0.1));
        border-radius: 2px;
      }
    }

    h2 {
      font-size: 1.8em;
      color: #2d4a6d;
      margin: 1.8em 0 0.8em;
      position: relative;
      font-weight: 600;

      &::after {
        content: '';
        position: absolute;
        bottom: -0.3em;
        left: 0;
        width: 60px;
        height: 3px;
        background: linear-gradient(90deg, rgba(45, 74, 109, 0.8), rgba(45, 74, 109, 0.1));
        border-radius: 1.5px;
      }
    }

    h3 {
      font-size: 1.4em;
      color: #34557a;
      margin: 1.5em 0 0.8em;
      position: relative;
      padding-left: 1em;
      font-weight: 600;

      &::before {
        content: '';
        position: absolute;
        left: 0;
        top: 0.2em;
        bottom: 0.2em;
        width: 3px;
        background: linear-gradient(180deg, #34557a, rgba(52, 85, 122, 0.2));
        border-radius: 1.5px;
      }
    }

    h4 {
      font-size: 1.2em;
      color: #3d5a7a;
      margin: 1.2em 0 0.6em;
      font-weight: 600;
    }

    h5 {
      font-size: 1.1em;
      color: #456789;
      margin: 1em 0 0.5em;
      font-weight: 500;
    }

    h6 {
      font-size: 1em;
      color: #516b88;
      margin: 1em 0 0.5em;
      font-weight: 500;
    }

    p {
      margin: 1em 0;
      line-height: 1.8;
    }

    img {
      max-width: 100%;
      height: auto;
      border-radius: 12px;
      margin: 1.5em auto;
      display: block;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
      transition: transform 0.3s ease, box-shadow 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 24px rgba(0, 0, 0, 0.15);
      }
    }

    pre {
      background: #f8fafc;
      border-radius: 8px;
      padding: 16px;
      margin: 1.5em 0;
      position: relative;
      border: 1px solid #e2e8f0;

      .pre-wrapper {
        overflow-x: auto;
        position: relative;

        /* 自定义滚动条样式 */
        &::-webkit-scrollbar {
          height: 8px;
          background-color: transparent;
        }

        &::-webkit-scrollbar-thumb {
          background: #cbd5e0;
          border-radius: 4px;
          
          &:hover {
            background: #a0aec0;
          }
        }

        &::-webkit-scrollbar-track {
          background: transparent;
          border-radius: 4px;
        }
      }

      code {
        padding-right: 85px; /* 为复制按钮预留空间 */
        display: inline-block;
        min-width: 100%;
      }

      &::before {
        content: attr(data-language);
        position: absolute;
        top: 8px;
        right: 85px;
        font-size: 0.7em;
        color: #94a3b8;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        font-weight: 500;
        background: #f1f5f9;
        padding: 2px 6px;
        border-radius: 3px;
        border: 1px solid #e2e8f0;
        z-index: 2;
      }

      .copy-button {
        position: absolute;
        top: 6px;
        right: 6px;
        height: 24px;
        padding: 0 10px;
        font-size: 0.75em;
        color: #64748b;
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        display: flex;
        align-items: center;
        gap: 4px;
        opacity: 0;
        transform: translateY(-4px);
        z-index: 3;
        font-family: system-ui, -apple-system, sans-serif;
        font-weight: 500;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
        backdrop-filter: blur(8px);

        svg {
          width: 12px;
          height: 12px;
          stroke: currentColor;
          stroke-width: 2;
          transition: all 0.2s ease;
        }

        &:hover {
          background: #f1f5f9;
          color: #0f172a;
          border-color: #cbd5e0;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05),
                     0 1px 2px rgba(0, 0, 0, 0.1);
          transform: translateY(0);
        }

        &:active {
          transform: translateY(0) scale(0.96);
          background: #e2e8f0;
        }

        &.copied {
          background: #2B5876;
          color: white;
          border-color: #1a365d;
          box-shadow: 0 1px 3px rgba(43, 88, 118, 0.2);

          svg {
            stroke: white;
            transform: scale(1.1);
          }

          &:hover {
            background: #34557a;
            border-color: #2B5876;
            box-shadow: 0 2px 4px rgba(43, 88, 118, 0.2),
                       0 1px 2px rgba(43, 88, 118, 0.1);
          }
        }
      }

      &:hover {
        .copy-button {
          opacity: 1;
          transform: translateY(0);
        }
      }
    }

    code {
      color: #2d4a6d;
      background: rgba(45, 74, 109, 0.06);
      padding: 0.2em 0.4em;
      border-radius: 4px;
      font-size: 0.9em;
    }

    blockquote {
      margin: 1.5em 0;
      padding: 1em 1.5em;
      border-left: 4px solid #2d4a6d;
      background: rgba(45, 74, 109, 0.03);
      border-radius: 0 8px 8px 0;
      color: #4a5568;
      font-style: italic;

      p {
        margin: 0;
      }
    }

    ul, ol {
      margin: 0;
      padding: 0;

      li {
        line-height: 1.5;
        color: inherit;
        margin: 0;
        padding: 0;
        border: none;
        background: none;
        box-shadow: none;
        transition: none;

        p {
          margin: 0;
          padding: 0;
          border: none;
          background: none;
          box-shadow: none;
        }
      }
    }

    ol {
      list-style: decimal;
      padding-left: 1.5em;
    }

    ul {
      list-style: disc;
      padding-left: 1.5em;
    }

    table {
      width: 100%;
      margin: 1.5em 0;
      border-collapse: collapse;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

      th, td {
        padding: 12px 16px;
        border: 1px solid #e2e8f0;
        text-align: left;
      }

      th {
        background: #f8fafc;
        font-weight: 600;
        color: #2c3e50;
      }

      tr:nth-child(even) {
        background: #f8fafc;
      }

      tr:hover {
        background: #f1f5f9;
      }
    }

    hr {
      margin: 2em 0;
      border: none;
      border-top: 2px solid #edf2f7;
    }

    a {
      color: #2d4a6d;
      text-decoration: none;
      border-bottom: 1px dashed #2d4a6d;
      transition: all 0.3s ease;

      &:hover {
        color: #1a365d;
        border-bottom-style: solid;
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
      background: linear-gradient(180deg, #2B5876, #4E4376);
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
        border-color: #2B5876;
        box-shadow: 0 0 0 3px rgba(43, 88, 118, 0.1);
      }
    }

    .el-button {
      margin-top: 15px;
      float: right;
      padding: 8px 16px;
      font-size: 0.95em;
      border-radius: 8px;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      border: none;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
      }
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
        margin-bottom: 12px;

        .avatar-link {
          text-decoration: none;
          margin-right: 12px;
          
          &:hover {
            opacity: 0.9;
          }
        }

        .comment-info {
          display: flex;
          align-items: center;
          gap: 8px;

          .username {
            font-weight: 500;
            color: #2c3e50;
            text-decoration: none;

            &:hover {
              color: #3699FF;
            }
          }

          .author-tag {
            font-size: 0.75em;
            padding: 2px 6px;
            background: rgba(43, 88, 118, 0.08);
            color: #2B5876;
            border-radius: 4px;
            font-weight: 500;
          }

          .time {
            color: #718096;
            font-size: 0.9em;
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