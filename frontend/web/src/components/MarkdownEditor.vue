<template>
  <div class="markdown-editor">
    <div class="editor-toolbar">
      <el-button-group>
        <el-button @click="insertMarkdown('**', '**')" size="small">
          <el-icon><Edit /></el-icon>
        </el-button>
        <el-button @click="insertMarkdown('*', '*')" size="small">
          <el-icon><EditPen /></el-icon>
        </el-button>
        <el-button @click="insertMarkdown('# ')" size="small">
          <el-icon><Document /></el-icon>
        </el-button>
        <el-button @click="insertMarkdown('> ')" size="small">
          <el-icon><ChatDotRound /></el-icon>
        </el-button>
        <el-button @click="insertMarkdown('```\n', '\n```')" size="small">
          <el-icon><Monitor /></el-icon>
        </el-button>
        <el-button @click="insertMarkdown('- ')" size="small">
          <el-icon><List /></el-icon>
        </el-button>
        <el-button @click="insertMarkdown('1. ')" size="small">
          <el-icon><Sort /></el-icon>
        </el-button>
      </el-button-group>
    </div>
    <div class="editor-container">
      <div class="editor-section">
        <el-input
          v-model="content"
          type="textarea"
          :rows="20"
          placeholder="请输入 Markdown 内容..."
          @input="handleInput"
          resize="none"
        />
      </div>
      <div class="preview-section">
        <div class="markdown-preview" v-html="renderedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import {
  Edit,
  EditPen,
  Document,
  ChatDotRound,
  Monitor,
  List,
  Sort
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const content = ref(props.modelValue)

// 配置 marked
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value
    }
    return hljs.highlightAuto(code).value
  },
  breaks: true
})

// 渲染 Markdown 内容
const renderedContent = computed(() => {
  return marked(content.value)
})

// 处理输入
const handleInput = (value) => {
  emit('update:modelValue', value)
}

// 插入 Markdown 语法
const insertMarkdown = (prefix, suffix = '') => {
  const textarea = document.querySelector('.el-textarea__inner')
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end)
  
  const newContent = content.value.substring(0, start) +
    prefix + selectedText + suffix +
    content.value.substring(end)
  
  content.value = newContent
  emit('update:modelValue', newContent)
  
  // 恢复光标位置
  setTimeout(() => {
    textarea.focus()
    textarea.setSelectionRange(
      start + prefix.length,
      end + prefix.length
    )
  }, 0)
}

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
  content.value = newValue
})
</script>

<style scoped lang="scss">
.markdown-editor {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;

  .editor-toolbar {
    padding: 8px;
    border-bottom: 1px solid #dcdfe6;
    background-color: #f5f7fa;
  }

  .editor-container {
    display: flex;
    height: 500px;

    .editor-section {
      flex: 1;
      border-right: 1px solid #dcdfe6;
      overflow: hidden;

      :deep(.el-textarea__inner) {
        height: 100%;
        border: none;
        border-radius: 0;
        resize: none;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
        font-size: 14px;
        line-height: 1.6;
        padding: 16px;
      }
    }

    .preview-section {
      flex: 1;
      padding: 16px;
      overflow-y: auto;
      background-color: #fff;

      .markdown-preview {
        font-size: 14px;
        line-height: 1.6;

        :deep(h1) {
          font-size: 2em;
          margin-bottom: 0.5em;
          padding-bottom: 0.3em;
          border-bottom: 1px solid #eaecef;
        }

        :deep(h2) {
          font-size: 1.5em;
          margin-bottom: 0.5em;
          padding-bottom: 0.3em;
          border-bottom: 1px solid #eaecef;
        }

        :deep(h3) {
          font-size: 1.25em;
          margin-bottom: 0.5em;
        }

        :deep(p) {
          margin-bottom: 1em;
        }

        :deep(ul), :deep(ol) {
          padding-left: 2em;
          margin-bottom: 1em;
        }

        :deep(li) {
          margin-bottom: 0.5em;
        }

        :deep(blockquote) {
          margin: 0;
          padding: 0 1em;
          color: #6a737d;
          border-left: 0.25em solid #dfe2e5;
        }

        :deep(pre) {
          padding: 16px;
          overflow: auto;
          background-color: #f6f8fa;
          border-radius: 4px;
          margin-bottom: 1em;
        }

        :deep(code) {
          padding: 0.2em 0.4em;
          margin: 0;
          font-size: 85%;
          background-color: rgba(27, 31, 35, 0.05);
          border-radius: 3px;
        }

        :deep(img) {
          max-width: 100%;
          box-sizing: border-box;
        }

        :deep(table) {
          border-spacing: 0;
          border-collapse: collapse;
          margin-bottom: 1em;
          width: 100%;

          th, td {
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
          }

          tr {
            background-color: #fff;
            border-top: 1px solid #c6cbd1;

            &:nth-child(2n) {
              background-color: #f6f8fa;
            }
          }
        }
      }
    }
  }
}
</style> 