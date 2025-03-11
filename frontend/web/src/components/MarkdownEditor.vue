<template>
  <div class="markdown-editor">
    <div class="editor-toolbar">
      <div class="toolbar-left">
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

      <div class="toolbar-right">
        <el-radio-group v-model="viewMode" size="small">
          <el-radio-button label="edit">编辑</el-radio-button>
          <el-radio-button label="split">分屏</el-radio-button>
          <el-radio-button label="preview">预览</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <div class="editor-container" :class="viewMode">
      <div class="editor-section" v-show="viewMode !== 'preview'">
        <el-input
          v-model="content"
          type="textarea"
          :rows="20"
          placeholder="请输入 Markdown 内容..."
          @input="handleInput"
          resize="none"
        />
      </div>
      <div class="preview-section" v-show="viewMode !== 'edit'">
        <div class="markdown-preview" v-html="renderedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
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
const viewMode = ref('split') // 默认分屏模式

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
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;

  .editor-toolbar {
    padding: 12px;
    border-bottom: 1px solid #dcdfe6;
    background: linear-gradient(to bottom, #f8fafc, #f1f5f9);
    display: flex;
    justify-content: space-between;
    align-items: center;

    .toolbar-left {
      .el-button-group {
        margin-right: 8px;
        border-radius: 6px;
        overflow: hidden;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

        .el-button {
          border: 1px solid #e2e8f0;
          background: #fff;
          color: #475569;
          font-size: 14px;
          height: 32px;
          padding: 0 12px;
          transition: all 0.2s ease;

          &:hover {
            background: #f8fafc;
            color: var(--el-color-primary);
            transform: translateY(-1px);
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
          }

          &:active {
            transform: translateY(0);
          }

          .el-icon {
            font-size: 16px;
          }

          & + .el-button {
            border-left: none;
          }
        }
      }
    }

    .toolbar-right {
      .el-radio-group {
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
      }
    }
  }

  .editor-container {
    display: flex;
    height: calc(100vh - 300px);
    min-height: 500px;
    background: #fff;

    &.edit {
      .editor-section {
        flex: 1;
      }
    }

    &.preview {
      .preview-section {
        flex: 1;
      }
    }

    &.split {
      .editor-section,
      .preview-section {
        flex: 0 0 50%;
        border: none;
      }

      .editor-section {
        border-right: 1px solid #e2e8f0;
      }
    }

    .editor-section {
      background: #fafafa;
      position: relative;

      &::before {
        content: 'Markdown';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        padding: 4px 16px;
        background: #f1f5f9;
        color: #64748b;
        font-size: 12px;
        font-weight: 500;
        border-bottom: 1px solid #e2e8f0;
        z-index: 1;
      }

      :deep(.el-textarea__inner) {
        height: 100%;
        border: none;
        border-radius: 0;
        resize: none;
        font-family: 'JetBrains Mono', 'Fira Code', 'Monaco', monospace;
        font-size: 14px;
        line-height: 1.6;
        padding: 48px 16px 16px;
        color: #334155;
        background: #fafafa;
        box-shadow: none;
        
        &:focus {
          box-shadow: none;
        }

        &::placeholder {
          color: #94a3b8;
        }
      }
    }

    .preview-section {
      padding: 48px 16px 16px;
      overflow-y: auto;
      background: #fff;
      position: relative;

      &::before {
        content: '预览';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        padding: 4px 16px;
        background: #f1f5f9;
        color: #64748b;
        font-size: 12px;
        font-weight: 500;
        border-bottom: 1px solid #e2e8f0;
        z-index: 1;
      }

      .markdown-preview {
        font-size: 14px;
        line-height: 1.6;
        color: #334155;
        max-width: 100%;

        :deep(h1) {
          font-size: 2em;
          margin: 0.5em 0 0.5em;
          padding-bottom: 0.3em;
          border-bottom: 1px solid #e2e8f0;
        }

        :deep(h2) {
          font-size: 1.5em;
          margin: 0.5em 0 0.5em;
          padding-bottom: 0.3em;
          border-bottom: 1px solid #e2e8f0;
        }

        :deep(h3) {
          font-size: 1.25em;
          margin: 0.5em 0 0.5em;
        }

        :deep(p) {
          margin: 0.5em 0 1em;
          line-height: 1.7;
        }

        :deep(ul), :deep(ol) {
          padding-left: 2em;
          margin: 0.5em 0 1em;
        }

        :deep(li) {
          margin: 0.5em 0;
        }

        :deep(blockquote) {
          margin: 0 0 1em;
          padding: 0.5em 1em;
          color: #475569;
          border-left: 4px solid #e2e8f0;
          background: #f8fafc;
          border-radius: 4px;
        }

        :deep(pre) {
          margin: 0 0 1em;
          padding: 1em;
          background: #1e293b;
          border-radius: 6px;
          overflow-x: auto;
          
          code {
            color: #e2e8f0;
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
            line-height: 1.5;
          }
        }

        :deep(code):not(pre code) {
          padding: 0.2em 0.4em;
          margin: 0 0.2em;
          font-size: 0.9em;
          background: #f1f5f9;
          border-radius: 4px;
          color: #ef4444;
        }

        :deep(img) {
          max-width: 100%;
          border-radius: 8px;
          box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }

        :deep(table) {
          width: 100%;
          margin: 1em 0;
          border-collapse: collapse;

          th, td {
            padding: 0.75em 1em;
            border: 1px solid #e2e8f0;
          }

          th {
            background: #f8fafc;
            font-weight: 600;
          }

          tr:nth-child(even) {
            background: #f8fafc;
          }
        }
      }
    }
  }
}
</style> 