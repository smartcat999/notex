<template>
  <div class="markdown-preview" v-html="renderedContent" ref="contentRef"></div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted, nextTick } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { ElMessage } from 'element-plus'

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

const contentRef = ref(null)

// 配置 markdown-it
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        const html = hljs.highlight(str, { language: lang }).value
        return `<pre data-language="${lang}"><code>${html}</code></pre>`
      } catch (__) {}
    }
    return `<pre><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// 渲染 markdown 内容
const renderedContent = computed(() => {
  const content = md.render(props.content || '')
  nextTick(() => {
    updateCodeSections()
  })
  return content
})

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
            ElMessage({
              message: '复制成功',
              type: 'success',
              duration: 2000,
              offset: 80
            })
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

onMounted(() => {
  updateCodeSections()
  window.addEventListener('resize', updateCodeSections)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateCodeSections)
})
</script>

<style lang="scss" scoped>
.markdown-preview {
  font-size: 16px;
  line-height: 1.8;
  color: #24292f;

  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    position: relative;
    margin: 1.8em 0 0.9em;
    font-weight: 600;
    line-height: 1.3;
    color: #1f2937;
    padding: 0.8em 1em;
    border-radius: 12px;
    background: #f0f7ff;
    transition: all 0.3s ease;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
    border-left: 4px solid #8b5cf6;

    &::before {
      position: absolute;
      left: -30px;
      top: 50%;
      transform: translateY(-50%);
      background: #8b5cf6;
      color: white;
      padding: 3px 6px;
      border-radius: 4px;
      font-size: 11px;
      font-weight: 600;
    }

    &:first-child {
      margin-top: 0.5em;
    }
  }

  :deep(h1) {
    font-size: 1.8em;
    background: linear-gradient(135deg, #f0f7ff 0%, #ffffff 100%);
    border-bottom: 3px solid #8b5cf6;
    margin-top: 1.2em;
    &::before { content: 'H1'; }
  }

  :deep(h2) {
    font-size: 1.5em;
    background: linear-gradient(135deg, #f5f9ff 0%, #ffffff 100%);
    border-bottom: 2px solid #9f75ff;
    &::before { content: 'H2'; }
  }

  :deep(h3) {
    font-size: 1.3em;
    background: linear-gradient(135deg, #f8faff 0%, #ffffff 100%);
    &::before { content: 'H3'; }
  }

  :deep(h4) {
    font-size: 1.2em;
    color: #4a5568;
    margin: 1.2em 0 0.6em;
    font-weight: 600;
  }

  :deep(h5) {
    font-size: 1.1em;
    color: #4a5568;
    margin: 1em 0 0.5em;
    font-weight: 500;
  }

  :deep(h6) {
    font-size: 1em;
    color: #4a5568;
    margin: 1em 0 0.5em;
    font-weight: 500;
  }

  :deep(p) {
    margin: 1em 0;
    line-height: 1.7;
    color: #374151;
    font-size: 0.95em;
    padding: 0.6em 1em;
    background: #fafafa;
    border-radius: 8px;
    border-left: 4px solid #e5e7eb;
    transition: all 0.3s ease;

    &:hover {
      border-left-color: #8b5cf6;
      background: #f9fafb;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
    }
  }

  :deep(ul), :deep(ol) {
    padding: 0 0 0 1.2em;
    margin: 0.5em 0;
    background: none;
    border: none;
    border-radius: 0;

    li {
      margin: 0.2em 0;
      padding: 0;
      line-height: 1.6;
      color: #374151;
      background: none;
      border: none;
      box-shadow: none;

      &::marker {
        color: #2B5876;
      }

      p {
        margin: 0;
        padding: 0;
        background: none;
        border: none;
      }
    }
  }

  :deep(ol) {
    list-style: none;
    counter-reset: markdown-counter;
    padding-left: 2em;
    border: none;

    li {
      position: relative;
      counter-increment: markdown-counter;
      margin: 0.5em 0;
      padding: 0.5em 0.8em;
      border: none !important;
      background: none;
      box-shadow: none;
      transition: none;
      transform: none;
      outline: none;
      text-decoration: none;

      &:hover, &:focus, &:active {
        transform: none;
        border: none !important;
        outline: none;
        text-decoration: none;
      }

      &::before {
        content: counter(markdown-counter) ".";
        position: absolute;
        left: -2em;
        width: 1.5em;
        text-align: right;
        color: #4a5568;
        font-weight: 500;
        border: none;
        background: none;
        box-shadow: none;
      }

      > * {
        border: none !important;
        margin: 0;
        padding: 0;
        background: none;
        box-shadow: none;
        transition: none;
        transform: none;
        outline: none;
        text-decoration: none;
      }
    }
  }

  :deep(ul) {
    list-style: disc;
    padding-left: 1.5em;

    li {
      border: none;
      outline: none;
      text-decoration: none;
      margin: 0.5em 0;
      padding: 0;
      background: none;
      box-shadow: none;
      transition: none;
      transform: none;

      &:hover {
        transform: none;
        border: none;
        outline: none;
        text-decoration: none;
      }

      p {
        margin: 0;
        padding: 0;
        border: none;
        background: none;
        box-shadow: none;
        transition: none;
        transform: none;
        outline: none;
        text-decoration: none;

        &:hover {
          transform: none;
          border: none;
          outline: none;
          text-decoration: none;
        }
      }
    }
  }

  :deep(blockquote) {
    margin: 1em 0;
    padding: 0.5em 1em;
    color: var(--el-text-color-regular);
    background: var(--el-fill-color-light);
    border-left: 4px solid var(--el-border-color);
    border-radius: 4px;

    p {
      margin: 0.5em 0;
    }
  }

  :deep(code) {
    padding: 36px 85px 8px 8px;
    display: inline-block;
    min-width: 100%;
    font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace;
    font-size: 0.9em;
    line-height: 1.5;
    background: transparent;
    border-radius: 0;
  }

  :deep(pre) {
    background: #f8fafc;
    border-radius: 8px;
    padding: 12px;
    margin: 1.5em 0;
    position: relative;
    border: 1px solid #e2e8f0;

    .pre-wrapper {
      overflow-x: auto;
      position: relative;
      margin-top: 0;

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
      padding: 36px 85px 8px 8px;
      display: inline-block;
      min-width: 100%;
      font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace;
      font-size: 0.9em;
      line-height: 1.5;
      background: transparent;
      border-radius: 0;
    }

    &[data-language]::before {
      content: attr(data-language);
      position: absolute;
      top: 8px;
      left: 12px;
      height: 20px;
      font-size: 12px;
      color: #64748b;
      font-weight: 500;
      text-transform: lowercase;
      letter-spacing: 0.02em;
      background: #f1f5f9;
      padding: 2px 8px;
      border-radius: 4px;
      border: 1px solid #e2e8f0;
      z-index: 2;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
      transition: all 0.2s ease;
      line-height: 1;

      &:hover {
        color: #475569;
        background: #e2e8f0;
        border-color: #cbd5e0;
      }
    }

    .copy-button {
      position: absolute;
      top: 6px;
      right: 8px;
      height: 24px;
      min-width: 65px;
      padding: 0 10px;
      font-size: 12px;
      color: #64748b;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 4px;
      opacity: 0;
      transform: translateY(-4px);
      z-index: 3;
      font-family: system-ui, -apple-system, sans-serif;
      font-weight: 500;
      box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
      backdrop-filter: blur(8px);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;

      svg {
        width: 12px;
        height: 12px;
        stroke: currentColor;
        stroke-width: 2;
        transition: all 0.2s ease;
        flex-shrink: 0;
      }

      span {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
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

  :deep(table) {
    margin: 1em 0;
    border-collapse: collapse;
    width: 100%;
    
    th, td {
      padding: 0.5em 1em;
      border: 1px solid var(--el-border-color);
    }
    
    th {
      background: var(--el-fill-color-light);
      font-weight: 600;
    }
    
    tr:nth-child(2n) {
      background: var(--el-fill-color-lighter);
    }
  }

  :deep(hr) {
    height: 1px;
    margin: 2em 0;
    border: none;
    background: var(--el-border-color-light);
  }

  :deep(img) {
    max-width: 100%;
    height: auto;
    margin: 1em 0;
    border-radius: 8px;
  }

  :deep(a) {
    color: var(--el-color-primary);
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
}
</style> 