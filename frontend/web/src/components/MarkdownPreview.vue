<template>
  <div class="markdown-preview" v-html="renderedContent"></div>
</template>

<script setup>
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

// 配置 markdown-it
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
    return '' // 使用默认的转义
  }
})

// 渲染 markdown 内容
const renderedContent = computed(() => {
  return md.render(props.content || '')
})
</script>

<style lang="scss" scoped>
.markdown-preview {
  font-size: 16px;
  line-height: 1.8;
  color: var(--el-text-color-primary);

  :deep(h1),
  :deep(h2),
  :deep(h3),
  :deep(h4),
  :deep(h5),
  :deep(h6) {
    margin: 1.6em 0 0.6em;
    font-weight: 600;
    line-height: 1.25;
  }

  :deep(h1) {
    font-size: 2em;
    border-bottom: 1px solid var(--el-border-color-light);
    padding-bottom: 0.3em;
  }

  :deep(h2) {
    font-size: 1.5em;
    border-bottom: 1px solid var(--el-border-color-light);
    padding-bottom: 0.3em;
  }

  :deep(h3) { font-size: 1.25em; }
  :deep(h4) { font-size: 1em; }
  :deep(h5) { font-size: 0.875em; }
  :deep(h6) { font-size: 0.85em; }

  :deep(p) {
    margin: 1em 0;
    line-height: 1.8;
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

  :deep(ul),
  :deep(ol) {
    padding-left: 2em;
    margin: 1em 0;
  }

  :deep(li) {
    margin: 0.5em 0;
  }

  :deep(code) {
    padding: 0.2em 0.4em;
    margin: 0;
    font-size: 0.85em;
    background: var(--el-fill-color-light);
    border-radius: 3px;
    font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace;
  }

  :deep(pre) {
    margin: 1em 0;
    padding: 1em;
    overflow: auto;
    background: var(--el-fill-color-light);
    border-radius: 6px;
    
    code {
      padding: 0;
      background: transparent;
      font-size: 0.9em;
      line-height: 1.5;
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