<template>
  <div class="word-polish-container">
    <div class="polish-header">
      <h2>Word文档润色</h2>
      <div class="header-actions">
        <el-select v-model="selectedModel" placeholder="选择AI模型" class="model-selector">
          <el-option
            v-for="model in aiStore.availableModels"
            :key="model.id"
            :label="model.name"
            :value="model.id"
          />
        </el-select>
      </div>
      <el-upload
        class="upload-area"
        drag
        action="#"
        :auto-upload="false"
        :on-change="handleFileChange"
        accept=".doc,.docx"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽文件到此处或 <em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 .doc, .docx 格式的Word文档
          </div>
        </template>
      </el-upload>
    </div>

    <div class="polish-content" v-if="fileContent">
      <div class="content-preview">
        <h3>原文预览</h3>
        <div class="preview-box" v-html="fileContent"></div>
      </div>

      <div class="polish-options">
        <el-form :model="polishOptions" label-width="120px">
          <el-form-item label="润色风格">
            <el-select v-model="polishOptions.style" placeholder="选择润色风格">
              <el-option label="学术论文" value="academic" />
              <el-option label="商务文档" value="business" />
              <el-option label="创意写作" value="creative" />
              <el-option label="技术文档" value="technical" />
              <el-option label="自定义" value="custom" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="polishOptions.style === 'custom'" label="自定义风格">
            <el-input
              v-model="polishOptions.customStyle"
              type="textarea"
              :rows="3"
              placeholder="请输入自定义润色风格要求，例如：'简洁明了，重点突出，适合技术文档阅读'"
            />
          </el-form-item>

          <el-form-item label="润色重点">
            <el-checkbox-group v-model="polishOptions.focus">
              <el-checkbox value="grammar">语法优化</el-checkbox>
              <el-checkbox value="vocabulary">词汇提升</el-checkbox>
              <el-checkbox value="structure">结构优化</el-checkbox>
              <el-checkbox value="tone">语气调整</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="语言风格">
            <el-radio-group v-model="polishOptions.language">
              <el-radio value="formal">正式</el-radio>
              <el-radio value="casual">随意</el-radio>
              <el-radio value="neutral">中性</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>

      <div class="action-buttons">
        <el-button type="primary" @click="startPolish" v-if="!isPolishing">
          开始润色
        </el-button>
        <el-button type="danger" @click="cancelPolish" v-if="isPolishing">
          取消润色
        </el-button>
        <el-button @click="resetForm">重置</el-button>
      </div>

      <div class="polish-status" v-if="isPolishing">
        <el-icon class="is-loading"><loading /></el-icon>
        <span>正在润色中...</span>
      </div>

      <div class="polish-result" v-if="polishedContent">
        <h3>润色结果</h3>
        <div class="result-box">
          <div class="result-content" v-html="polishedContent"></div>
          <div class="result-actions">
            <el-button type="primary" @click="downloadResult">
              <el-icon><download /></el-icon>
              下载润色后的文档
            </el-button>
            <el-button @click="copyToClipboard">
              <el-icon><document-copy /></el-icon>
              复制内容
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { UploadFilled, Download, DocumentCopy, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useAIStore } from '@/stores/ai'
import { marked } from 'marked'
import mammoth from 'mammoth'
import { Document, Packer, Paragraph, TextRun, HeadingLevel } from 'docx'
import { saveAs } from 'file-saver'

const aiStore = useAIStore()
const fileContent = ref('')
const polishedContent = ref('')
const isPolishing = ref(false)
const selectedModel = ref('')

const polishOptions = reactive({
  style: 'academic',
  focus: ['grammar', 'vocabulary'],
  language: 'formal',
  customStyle: ''
})

// 监听模型选择变化
watch(selectedModel, (newModel) => {
  aiStore.setCurrentModel(newModel)
})

// 初始化时设置默认模型
onMounted(() => {
  if (aiStore.defaultModel) {
    selectedModel.value = aiStore.defaultModel
  } else if (aiStore.availableModels.length > 0) {
    selectedModel.value = aiStore.availableModels[0].id
  }
  aiStore.setCurrentModel(selectedModel.value)
})

const handleFileChange = async (file) => {
  if (!file) return
  
  try {
    const arrayBuffer = await file.raw.arrayBuffer()
    const result = await mammoth.convertToHtml({ arrayBuffer })
    fileContent.value = result.value
    polishedContent.value = ''
  } catch (error) {
    console.error('Error reading file:', error)
    ElMessage.error('文件读取失败，请确保上传的是有效的Word文档')
    fileContent.value = ''
    polishedContent.value = ''
  }
}

const cancelPolish = async () => {
  try {
    await aiStore.cancelCurrentRequest()
    isPolishing.value = false
    ElMessage.success('已取消润色')
  } catch (error) {
    console.error('Cancel error:', error)
    ElMessage.error('取消失败，请重试')
  }
}

const startPolish = async () => {
  if (!fileContent.value) {
    ElMessage.warning('请先上传文档')
    return
  }

  const content = fileContent.value.trim()
  if (!content) {
    ElMessage.warning('文档内容为空，请重新上传')
    return
  }

  // 检查是否选择了模型
  if (!aiStore.currentModel) {
    ElMessage.warning('请先选择AI模型')
    return
  }

  isPolishing.value = true
  polishedContent.value = '' // 清空之前的结果
  
  try {
    const prompt = generatePolishPrompt()
    if (!prompt || !prompt.trim()) {
      ElMessage.error('生成提示词失败，请重试')
      return
    }

    // 使用流式响应
    await aiStore.sendMessage(prompt.trim(), -1, (content) => {
      // 实时更新润色结果
      polishedContent.value = marked(content)
    })
  } catch (error) {
    if (error.name === 'AbortError') {
      return
    }
    ElMessage.error('润色失败，请重试')
  } finally {
    isPolishing.value = false
  }
}

const generatePolishPrompt = () => {
  const styleMap = {
    academic: '学术论文',
    business: '商务文档',
    creative: '创意写作',
    technical: '技术文档',
    custom: polishOptions.customStyle || '自定义风格'
  }

  const focusMap = {
    grammar: '语法优化',
    vocabulary: '词汇提升',
    structure: '结构优化',
    tone: '语气调整'
  }

  const languageMap = {
    formal: '正式',
    casual: '随意',
    neutral: '中性'
  }

  const content = fileContent.value.trim() // 确保内容不为空
  if (!content) {
    ElMessage.error('文档内容为空，请重新上传')
    return null
  }

  // 验证自定义风格
  if (polishOptions.style === 'custom' && !polishOptions.customStyle.trim()) {
    ElMessage.error('请输入自定义润色风格要求')
    return null
  }

  return `请对以下文本进行润色，要求：
1. 润色风格：${styleMap[polishOptions.style]}
2. 润色重点：${polishOptions.focus.map(f => focusMap[f]).join('、')}
3. 语言风格：${languageMap[polishOptions.language]}
4. 保持原文的核心意思不变
5. 输出润色后的完整文本，保持原有的格式（包括标题、段落、列表等）
6. 使用Markdown格式输出，包括：
   - 使用 # 表示标题
   - 使用 * 或 _ 表示强调
   - 使用 > 表示引用
   - 使用 - 或 1. 表示列表
   - 使用 | 表示表格
7. 请逐段输出润色结果，每段之间保持适当的空行

原文：
${content}`
}

const downloadResult = async () => {
  try {
    // 将 HTML 转换为纯文本
    const tempDiv = document.createElement('div')
    tempDiv.innerHTML = polishedContent.value
    const textContent = tempDiv.textContent || tempDiv.innerText

    // 创建 Word 文档
    const doc = new Document({
      sections: [{
        properties: {},
        children: textContent.split('\n').map(line => {
          // 处理标题
          if (line.startsWith('# ')) {
            return new Paragraph({
              text: line.replace('# ', ''),
              heading: HeadingLevel.HEADING_1,
              spacing: { after: 200 }
            })
          } else if (line.startsWith('## ')) {
            return new Paragraph({
              text: line.replace('## ', ''),
              heading: HeadingLevel.HEADING_2,
              spacing: { after: 200 }
            })
          } else if (line.startsWith('### ')) {
            return new Paragraph({
              text: line.replace('### ', ''),
              heading: HeadingLevel.HEADING_3,
              spacing: { after: 200 }
            })
          }
          // 处理普通段落
          else if (line.trim()) {
            return new Paragraph({
              children: [
                new TextRun({
                  text: line,
                  size: 24 // 12pt
                })
              ],
              spacing: { after: 200 }
            })
          }
          // 处理空行
          else {
            return new Paragraph({
              spacing: { after: 200 }
            })
          }
        })
      }]
    })

    // 生成文档并下载
    const blob = await Packer.toBlob(doc)
    saveAs(blob, '润色后的文档.docx')
    ElMessage.success('文档下载成功')
  } catch (error) {
    console.error('下载文档失败:', error)
    ElMessage.error('下载文档失败，请重试')
  }
}

const copyToClipboard = async () => {
  try {
    const markdownText = polishedContent.value.replace(/<[^>]+>/g, '')
    await navigator.clipboard.writeText(markdownText)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败，请手动复制')
  }
}

const resetForm = () => {
  fileContent.value = ''
  polishedContent.value = ''
  polishOptions.style = 'academic'
  polishOptions.focus = ['grammar', 'vocabulary']
  polishOptions.language = 'formal'
  polishOptions.customStyle = ''
}
</script>

<style scoped lang="scss">
.word-polish-container {
  padding: 24px;
  background-color: #171717;
  border-radius: 12px;
  min-height: 100vh;

  .polish-header {
    margin-bottom: 24px;
    
    h2 {
      color: #e0e0e0;
      margin-bottom: 16px;
    }

    .header-actions {
      margin-bottom: 16px;
    }

    .upload-area {
      width: 100%;
      
      :deep(.el-upload-dragger) {
        background-color: #1f1f1f;
        border: 2px dashed #2a2a2a;
        border-radius: 8px;
        transition: all 0.3s ease;
        
        &:hover {
          border-color: #409EFF;
        }
        
        .el-icon--upload {
          font-size: 48px;
          color: #409EFF;
          margin-bottom: 16px;
        }
        
        .el-upload__text {
          color: #e0e0e0;
          
          em {
            color: #409EFF;
          }
        }
        
        .el-upload__tip {
          color: #a0a0a0;
        }
      }
    }
  }

  .polish-content {
    .content-preview, .polish-result {
      margin-bottom: 24px;
      
      h3 {
        color: #e0e0e0;
        margin-bottom: 16px;
      }
      
      .preview-box, .result-box {
        background-color: #1f1f1f;
        border: 1px solid #2a2a2a;
        border-radius: 8px;
        padding: 16px;
        color: #e0e0e0;
        min-height: 200px;
        
        .result-content {
          margin-bottom: 16px;
        }
        
        .result-actions {
          display: flex;
          gap: 12px;
          justify-content: flex-end;
        }
      }
    }

    .polish-options {
      background-color: #1f1f1f;
      border: 1px solid #2a2a2a;
      border-radius: 8px;
      padding: 24px;
      margin-bottom: 24px;
      
      :deep(.el-form-item__label) {
        color: #e0e0e0;
      }
      
      :deep(.el-input__wrapper),
      :deep(.el-select__wrapper),
      :deep(.el-checkbox__label),
      :deep(.el-radio__label) {
        background-color: #252525;
        border-color: #2a2a2a;
        color: #e0e0e0;
        
        &:hover {
          border-color: #409EFF;
        }
        
        &.is-focus {
          border-color: #409EFF;
          box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
        }
      }
    }

    .action-buttons {
      display: flex;
      gap: 12px;
      justify-content: center;
      margin-bottom: 24px;
    }

    .polish-status {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      margin-bottom: 24px;
      color: #409EFF;
      font-size: 14px;

      .el-icon {
        font-size: 16px;
      }
    }
  }
}

:deep(.markdown-body) {
  color: #e0e0e0;
  font-size: 14px;
  line-height: 1.6;
  
  p {
    margin: 8px 0;
  }
  
  pre {
    background-color: #252525;
    border-radius: 4px;
    padding: 12px;
    margin: 8px 0;
    overflow-x: auto;
    
    code {
      color: #e0e0e0;
      font-family: 'Fira Code', monospace;
    }
  }
  
  code {
    background-color: #252525;
    padding: 2px 4px;
    border-radius: 4px;
    font-family: 'Fira Code', monospace;
  }
}

.preview-box {
  background-color: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 16px;
  color: #e0e0e0;
  min-height: 200px;
  
  :deep(h1) {
    font-size: 24px;
    font-weight: bold;
    margin: 16px 0;
    color: #ffffff;
  }
  
  :deep(h2) {
    font-size: 20px;
    font-weight: bold;
    margin: 14px 0;
    color: #ffffff;
  }
  
  :deep(h3) {
    font-size: 18px;
    font-weight: bold;
    margin: 12px 0;
    color: #ffffff;
  }
  
  :deep(p) {
    margin: 8px 0;
    line-height: 1.6;
  }
  
  :deep(ul), :deep(ol) {
    margin: 8px 0;
    padding-left: 24px;
    
    li {
      margin: 4px 0;
    }
  }
  
  :deep(strong) {
    color: #ffffff;
    font-weight: bold;
  }
  
  :deep(em) {
    font-style: italic;
  }
  
  :deep(blockquote) {
    border-left: 4px solid #409EFF;
    margin: 8px 0;
    padding-left: 16px;
    color: #a0a0a0;
  }
  
  :deep(table) {
    border-collapse: collapse;
    width: 100%;
    margin: 8px 0;
    
    th, td {
      border: 1px solid #2a2a2a;
      padding: 8px;
      text-align: left;
    }
    
    th {
      background-color: #252525;
      color: #ffffff;
    }
  }
}

.result-box {
  background-color: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 16px;
  color: #e0e0e0;
  min-height: 200px;
  
  .result-content {
    margin-bottom: 16px;
    
    :deep(h1) {
      font-size: 24px;
      font-weight: bold;
      margin: 16px 0;
      color: #ffffff;
    }
    
    :deep(h2) {
      font-size: 20px;
      font-weight: bold;
      margin: 14px 0;
      color: #ffffff;
    }
    
    :deep(h3) {
      font-size: 18px;
      font-weight: bold;
      margin: 12px 0;
      color: #ffffff;
    }
    
    :deep(p) {
      margin: 8px 0;
      line-height: 1.6;
    }
    
    :deep(ul), :deep(ol) {
      margin: 8px 0;
      padding-left: 24px;
      
      li {
        margin: 4px 0;
      }
    }
    
    :deep(strong) {
      color: #ffffff;
      font-weight: bold;
    }
    
    :deep(em) {
      font-style: italic;
    }
    
    :deep(blockquote) {
      border-left: 4px solid #409EFF;
      margin: 8px 0;
      padding-left: 16px;
      color: #a0a0a0;
    }
    
    :deep(table) {
      border-collapse: collapse;
      width: 100%;
      margin: 8px 0;
      
      th, td {
        border: 1px solid #2a2a2a;
        padding: 8px;
        text-align: left;
      }
      
      th {
        background-color: #252525;
        color: #ffffff;
      }
    }
  }
  
  .result-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #2a2a2a;
  }
}
</style> 