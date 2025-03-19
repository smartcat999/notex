<template>
  <div class="ai-chat-container">
    <div class="ai-chat-header">
      <el-select v-model="selectedModel" placeholder="选择AI模型" class="model-selector" :loading="aiStore.isLoading">
        <el-option
          v-for="model in aiStore.availableModels"
          :key="model.id"
          :label="model.name"
          :value="model.id"
        />
      </el-select>
      <el-button type="primary" @click="clearChat" plain>
        <el-icon><Delete /></el-icon>
        清空对话
      </el-button>
    </div>

    <!-- 快捷工具栏 -->
    <div class="quick-tools">
      <div class="tools-container">
        <el-button class="tool-button" @click="openSearchDialog">
          <el-icon><Search /></el-icon>
          <span>搜全网</span>
        </el-button>
        <el-button class="tool-button" @click="openReadDocDialog">
          <el-icon><Document /></el-icon>
          <span>读文档</span>
        </el-button>
        <el-button class="tool-button" @click="openImageDialog">
          <el-icon><Picture /></el-icon>
          <span>生成图像</span>
        </el-button>
        <el-button class="tool-button" @click="openQuickCreateDialog">
          <el-icon><EditPen /></el-icon>
          <span>快速创作</span>
        </el-button>
        <el-button class="tool-button" @click="openCreatePPTDialog">
          <el-icon><Promotion /></el-icon>
          <span>生成PPT</span>
        </el-button>
        <el-button class="tool-button" @click="openLongTextWritingDialog">
          <el-icon><Document /></el-icon>
          <span>长文写作</span>
        </el-button>
        <el-button class="tool-button more-button">
          <el-icon><More /></el-icon>
          <span>更多</span>
        </el-button>
      </div>
    </div>

    <div class="chat-messages" ref="chatContainer">
      <!-- 在聊天消息区域顶部显示文档上传区域 -->
      <div v-if="showDocUploader" class="inline-doc-uploader">
        <div class="uploader-header">
          <h3>读文档</h3>
          <el-button type="text" @click="showDocUploader = false">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
        
        <!-- 文件上传区域 -->
        <div class="doc-upload-container">
          <div 
            class="doc-upload-area" 
            @click="triggerFileUpload"
            @dragover.prevent 
            @drop.prevent="handleFileDrop"
          >
            <input 
              type="file" 
              ref="fileInput" 
              class="file-input" 
              multiple 
              @change="handleFileSelect" 
              accept=".pdf,.docx,.xlsx,.pptx,.txt,.png,.jpg,.jpeg" 
            />
            <div class="upload-inner">
              <div class="upload-text">
                <p class="upload-title">在此处拖放文件</p>
                <p class="upload-desc">支持文件类型: pdf、docx、xlsx、pptx、txt、png、jpg、jpeg 等</p>
              </div>
            </div>
          </div>
          
          <!-- 单独的按钮区域 -->
          <div class="upload-button-area">
            <el-button class="local-file-button" @click.stop="triggerFileInputClick">
              <el-icon><Upload /></el-icon>
              本地文件
            </el-button>
          </div>
        </div>
        
        <!-- 上传的文件列表 -->
        <div v-if="uploadedFiles.length > 0" class="uploaded-files">
          <div v-for="(file, index) in uploadedFiles" :key="index" class="file-item">
            <div class="file-icon" :class="getFileTypeClass(file.name)">
              <el-icon><component :is="getFileTypeIcon(file.name)" /></el-icon>
            </div>
            <div class="file-info">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-size">{{ formatFileSize(file.size) }}</div>
            </div>
            <el-button type="text" @click="removeFile(index)" class="remove-file">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        
        <!-- 文件操作按钮 -->
        <div class="doc-actions">
          <el-button type="primary" @click="confirmDocUpload" :disabled="!uploadedFiles.length || isProcessingFile" :loading="isProcessingFile">
            {{ isProcessingFile ? '处理中...' : '确认使用文件' }}
          </el-button>
          <el-button @click="showDocUploader = false" :disabled="isProcessingFile">
            取消
          </el-button>
        </div>
      </div>
      
      <!-- 图像生成区域 -->
      <div v-if="showImageGenerator" class="inline-image-generator">
        <div class="generator-header">
          <h3>生成图像</h3>
          <el-button type="text" @click="showImageGenerator = false">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
        
        <!-- 使用新的图像生成组件 -->
        <ImageGenerator @send="handleImageSend" @close="showImageGenerator = false" />
      </div>

      <div v-if="orderedMessages.length === 0 && !showDocUploader && !showImageGenerator" class="empty-chat">
        <el-icon class="empty-icon"><ChatDotRound /></el-icon>
        <p>开始与AI助手对话</p>
      </div>
      
      <!-- 上传文件成功提示 -->
      <div v-if="activeDocFiles.length > 0" class="active-files-info">
        <div class="active-files-header">
          <div class="active-files-title">
            <el-icon><Document /></el-icon>
            <span>当前使用的文件</span>
          </div>
          <el-button type="text" @click="clearActiveFiles">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
        <div class="active-files-list">
          <div v-for="(fileName, index) in activeDocFiles" :key="index" class="active-file-item">
            {{ fileName }}
          </div>
        </div>
      </div>
      
      <div 
        v-for="message in orderedMessages" 
        :key="message.id"
        class="message"
        :class="[message.role]"
      >
        <div class="message-avatar">
          <el-avatar v-if="message.role === 'user'" :size="36" :src="userAvatar">
            {{ userInitials }}
          </el-avatar>
          <AIAvatar 
            v-else 
            :size="36" 
            :provider-id="message.providerId || currentProvider"
            :model-id="message.modelId || selectedModel"
          />
        </div>
        <div class="message-content">
          <template v-if="!isEditing(message.id)">
            <div class="message-container">
              <div class="message-text" :class="{ 'streaming': message.isStreaming }">
                <template v-if="message.role === 'user'">
                  <div class="user-message">{{ message.content }}</div>
                </template>
                <template v-else>
                  <div class="ai-message">
                    <div v-if="message.isStreaming" class="streaming-content">
                      <span v-html="renderMarkdown(message.content)"></span>
                      <span class="typing-indicator">
                        <span class="dot"></span>
                        <span class="dot"></span>
                        <span class="dot"></span>
                      </span>
                    </div>
                    <div v-else v-html="renderMarkdown(message.content)"></div>
                  </div>
                </template>
              </div>
              <div class="message-actions" v-if="message.role === 'user'">
                <el-button-group>
                  <el-button size="small" @click="startEditing(message)">
                    <el-icon><Edit /></el-icon>
                  </el-button>
                </el-button-group>
              </div>
            </div>
          </template>
          <template v-else>
            <div class="message-edit">
              <el-input
                v-model="editingContent"
                type="textarea"
                :rows="3"
                @keyup.enter.ctrl="saveEdit"
              />
              <div class="edit-actions">
                <el-button size="small" @click="cancelEditing">取消</el-button>
                <el-button size="small" type="primary" @click="saveEdit">保存</el-button>
              </div>
            </div>
          </template>
          <div class="message-time">{{ formatTime(message.timestamp) }}</div>
        </div>
      </div>
    </div>

    <div class="chat-input">
      <div class="input-wrapper">
        <el-input
          v-model="userInput"
          type="textarea"
          :rows="3"
          placeholder="输入您的问题..."
          resize="none"
          @keydown.enter.prevent="sendMessage"
          :disabled="aiStore.isLoading"
        />
        <el-button 
          type="primary" 
          :disabled="!userInput.trim() || aiStore.isLoading" 
          @click="sendMessage"
          class="send-button"
          :loading="aiStore.isLoading"
        >
          <el-icon><Position /></el-icon>
          发送
        </el-button>
      </div>
      <div class="input-actions">
        <div></div>
        <el-button 
          v-if="aiStore.isLoading"
          @click="handleStop"
          class="stop-btn"
        >
          停止生成
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { useAIStore } from '@/stores/ai'
import { ElMessage } from 'element-plus'
import { ChatDotRound, Position, Delete, Edit, Back, Right, Picture, Download, 
  Search, Document, EditPen, Promotion, More, Upload, Close, Files, VideoPlay,
  PictureFilled, Reading, Warning } from '@element-plus/icons-vue'
import AIAvatar from '@/components/ai/AIAvatar.vue'
import ImageGenerator from '@/components/ai/ImageGenerator.vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import axios from 'axios'
import * as XLSX from 'xlsx'
import * as pdfjs from 'pdfjs-dist'
import { getDocument } from 'pdfjs-dist'
import mammoth from 'mammoth'
import JSZip from 'jszip'
import { useRouter } from 'vue-router'

// 设置PDF.js worker路径
// 注意：需要在项目中安装pdfjs-dist
pdfjs.GlobalWorkerOptions.workerSrc = `//cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjs.version}/pdf.worker.min.js`

const userStore = useUserStore()
const aiStore = useAIStore()
const router = useRouter()
const userInput = ref('')
const chatContainer = ref(null)
const selectedModel = ref('')
const editingMessageId = ref(null)
const editingContent = ref('')

// 用户头像
const userAvatar = computed(() => userStore.user?.avatar || '')

// 用户名首字母
const userInitials = computed(() => {
  if (!userStore.user?.name) return ''
  return userStore.user.name.substring(0, 1).toUpperCase()
})

// 当前提供商
const currentProvider = computed(() => aiStore.currentProvider || '')

// 消息列表，按序列号排序
const orderedMessages = computed(() => {
  return [...aiStore.messages].sort((a, b) => (a.sequence || 0) - (b.sequence || 0))
})

// 历史记录撤销和重做状态
const canUndo = computed(() => aiStore.currentHistoryIndex > 0)
const canRedo = computed(() => aiStore.currentHistoryIndex < aiStore.messageHistory.length - 1)

// 图像生成相关状态
const imageDialogVisible = ref(false)
const imagePrompt = ref('')
const imageSize = ref('512x512')
const imageCount = ref(1)
const generatedImages = ref([])
const isGeneratingImage = ref(false)
const imageError = ref('')
const showImageGenerator = ref(false)
const showApiKeyPrompt = ref(false) // 控制API密钥配置引导UI的显示

// 新增对话框控制变量
const searchDialogVisible = ref(false)
const showDocUploader = ref(false)
const docPrompt = ref('')
const uploadedFiles = ref([])
const quickCreateDialogVisible = ref(false)
const createPPTDialogVisible = ref(false)
const longTextWritingDialogVisible = ref(false)
const isDeepSeekEnabled = ref(false)
const activeDocFiles = ref([]) // 当前活跃的文件列表
const fileContents = ref({}) // 存储文件名到文件内容的映射
const isProcessingFile = ref(false) // 文件处理状态

// 获取API基础URL
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api'

// 文件上传相关
const fileInput = ref(null)

// 初始化
onMounted(async () => {
  // 如果AIStore已初始化，直接使用当前模型
  if (aiStore.initialized) {
    selectedModel.value = aiStore.currentModel
    // console.log('AIChat: AIStore已初始化，使用当前模型:', selectedModel.value)
  } 
  // 如果AIStore有可用模型但未初始化
  else if (aiStore.availableModels.length > 0) {
  if (aiStore.defaultModel) {
    selectedModel.value = aiStore.defaultModel
      // console.log('AIChat: 使用默认模型:', selectedModel.value)
    } else {
    selectedModel.value = aiStore.availableModels[0].id
      // console.log('AIChat: 使用第一个可用模型:', selectedModel.value)
    }
  }
  // 如果AIStore未初始化，等待初始化
  else {
    // console.log('AIChat: AIStore未初始化，等待初始化完成')
    // 用户已登录但AIStore未初始化，手动初始化
    if (userStore.isAuthenticated && !aiStore.initialized && !aiStore.isInitializing) {
      await aiStore.initialize()
    }
  }
  
  // 异步初始化
  nextTick(() => {
    if (chatContainer.value && aiStore.messages.length > 0) {
      scrollToBottom()
    }
  })
})

// 监听AI存储初始化状态变化
watch(() => aiStore.initialized, (isInitialized) => {
  if (isInitialized) {
    // console.log('AIChat: AIStore初始化状态变化，重新设置模型')
    selectedModel.value = aiStore.currentModel
  }
})

// 监听AI存储默认模型变化
watch(() => aiStore.defaultModel, (newDefaultModel) => {
  if (newDefaultModel && (!selectedModel.value || selectedModel.value !== newDefaultModel)) {
    // console.log('AIChat: 默认模型变化:', newDefaultModel)
    selectedModel.value = newDefaultModel
  }
})

// 监听AI存储当前模型变化
watch(() => aiStore.currentModel, (newCurrentModel) => {
  if (newCurrentModel && selectedModel.value !== newCurrentModel) {
    // console.log('AIChat: 当前模型变化:', newCurrentModel)
    selectedModel.value = newCurrentModel
  }
})

// 监听消息变化
watch(() => aiStore.messages, () => {
  // 在消息变化时检查是否需要滚动
  nextTick(() => {
    if (chatContainer.value) {
      const isAtBottom = isUserAtBottom()
      if (isAtBottom) {
        scrollToBottom()
      }
    }
  })
}, { deep: true })

// 添加特定监听流式消息内容变化的监听器
watch(() => {
  // 查找流式消息
  const streamingMsg = aiStore.messages.find(msg => msg.isStreaming)
  return streamingMsg ? streamingMsg.content : null
}, (newContent) => {
  if (newContent && chatContainer.value) {
    // 仅当用户在底部时自动滚动
    const isAtBottom = isUserAtBottom()
    if (isAtBottom) {
      scrollToBottom()
    }
  }
})

// 监听选择的模型变化
watch(selectedModel, (newModel) => {
  if (newModel) {
    aiStore.setCurrentModel(newModel)
  }
})

// 检查用户是否在滚动容器底部
function isUserAtBottom() {
  if (!chatContainer.value) return true
  
  const scrollHeight = chatContainer.value.scrollHeight
  const clientHeight = chatContainer.value.clientHeight
  const scrollTop = chatContainer.value.scrollTop
  
  // 如果用户离底部在50px以内，视为在底部
  return scrollTop + clientHeight >= scrollHeight - 50
}

// 创建 markdown-it 实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch (__) {}
    }
    return '' // 使用默认的转义
  }
})

// 渲染Markdown内容
const renderMarkdown = (content) => {
  if (!content) return ''
  return md.render(content)
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

// 发送消息
const sendMessage = async () => {
  if (!userInput.value.trim() || aiStore.isLoading) return
  
  const message = userInput.value
  userInput.value = '' // 清空输入框
  
  try {
    let enhancedMessage = message
    // 如果有活跃文件，将文件内容作为上下文加入到消息中
    if (activeDocFiles.value.length > 0) {
      let fileContexts = ''
      
      for (const fileName of activeDocFiles.value) {
        const content = fileContents.value[fileName]
        if (content) {
          // 限制每个文件的内容长度，避免超出token限制
          const truncatedContent = content.length > 8000 
            ? content.substring(0, 8000) + '... [文件内容过长，已截断]' 
            : content
          
          fileContexts += `------- 文件 ${fileName} 的内容 -------\n${truncatedContent}\n\n`
        }
      }
      
      // 构建增强的提示，包含文件内容作为上下文
      enhancedMessage = `我需要您帮我分析以下文件，并回答我的问题:\n\n${fileContexts}\n我的问题是: ${message}`
    }
    
    // 发送增强的消息
    await aiStore.sendMessage(enhancedMessage)
    scrollToBottom()
  } catch (error) {
    ElMessage.error('获取AI回复失败')
    userInput.value = message // 如果发送失败，恢复输入内容
  }
}

// 清空对话
const clearChat = () => {
  aiStore.clearMessages()
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      const scrollHeight = chatContainer.value.scrollHeight
      chatContainer.value.scrollTo({
        top: scrollHeight,
        behavior: 'smooth'
      })
    }
  })
}

// 开始编辑消息
const startEditing = (message) => {
  if (aiStore.isLoading) return
  
  editingMessageId.value = message.id
  editingContent.value = message.content
  
  // 聚焦到编辑框
  nextTick(() => {
    const textarea = document.querySelector('.message-edit .el-textarea__inner')
    if (textarea) {
      textarea.focus()
    }
  })
}

// 检查是否正在编辑
const isEditing = (messageId) => {
  return editingMessageId.value === messageId
}

// 取消编辑
const cancelEditing = () => {
  editingMessageId.value = null
  editingContent.value = ''
}

// 保存编辑
const saveEdit = async () => {
  if (!editingContent.value.trim()) {
    ElMessage.warning('消息内容不能为空')
    return
  }
  
  const messageId = editingMessageId.value
  const content = editingContent.value
  
  // 清除编辑状态
  editingMessageId.value = null
  editingContent.value = ''
  
  // 保存编辑
  await aiStore.saveEdit(messageId, content)
  
  // 滚动到底部
  scrollToBottom()
}

// 撤销操作
const undo = () => {
  if (canUndo.value && !aiStore.isLoading) {
    aiStore.undo()
  }
}

// 重做操作
const redo = () => {
  if (canRedo.value && !aiStore.isLoading) {
    aiStore.redo()
  }
}

// 停止生成
const handleStop = () => {
  aiStore.cancelCurrentRequest()
}

// 打开图像生成对话框
const openImageDialog = () => {
  // 显示图像生成区域，保留聊天消息
  showImageGenerator.value = true
  
  // 确保在显示图像生成区域后滚动到顶部
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  })
}

// 处理图像发送到聊天
const handleImageSend = (markdownContent) => {
  // 创建包含图像的用户消息
  aiStore.sendMessage(markdownContent)
  
  // 关闭图像生成区域
  showImageGenerator.value = false
  
  // 滚动到底部
  nextTick(() => {
    scrollToBottom()
  })
}

// 打开搜索网对话框
const openSearchDialog = () => {
  searchDialogVisible.value = true
}

// 打开读文档功能
const openReadDocDialog = () => {
  // 显示文档上传区域，保留聊天消息
  showDocUploader.value = true
  // 重置文档相关状态
  docPrompt.value = ''
  
  // 确保在显示文档上传区域后滚动到顶部
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  })
}

// 打开快速创作对话框
const openQuickCreateDialog = () => {
  quickCreateDialogVisible.value = true
}

// 打开生成PPT对话框
const openCreatePPTDialog = () => {
  createPPTDialogVisible.value = true
}

// 打开长文写作对话框
const openLongTextWritingDialog = () => {
  longTextWritingDialogVisible.value = true
}

// 触发文件上传
const triggerFileUpload = (event) => {
  // 只有当点击的是上传区域自身时才触发文件上传
  if (event.target.classList.contains('doc-upload-area') || event.target.closest('.doc-upload-area') && !event.target.closest('.uploaded-file-item')) {
    triggerFileInputClick()
  }
}

// 直接点击文件输入框
const triggerFileInputClick = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

// 处理文件拖放
const handleFileDrop = (event) => {
  const files = event.dataTransfer.files
  if (files.length > 0) {
    handleFiles(files)
  }
}

// 处理文件选择
const handleFileSelect = (event) => {
  const files = event.target.files
  if (files.length > 0) {
    handleFiles(files)
  }
}

// 处理文件
const handleFiles = (files) => {
  for (let i = 0; i < files.length; i++) {
    // 检查文件类型
    const file = files[i]
    const extension = file.name.split('.').pop().toLowerCase()
    const validExtensions = ['pdf', 'docx', 'xlsx', 'pptx', 'txt', 'png', 'jpg', 'jpeg']
    
    if (validExtensions.includes(extension)) {
      // 添加到上传文件列表
      uploadedFiles.value.push(file)
    } else {
      ElMessage.warning(`不支持的文件类型: ${file.name}`)
    }
  }
  
  // 清空input，允许再次选择相同的文件
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 移除文件
const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1)
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 处理文档提示提交
const handleDocPromptSubmit = async () => {
  if (!docPrompt.value.trim() || uploadedFiles.value.length === 0) {
    return
  }
  
  try {
    // 构建发送文件的消息内容
    const fileNames = uploadedFiles.value.map(file => file.name).join(', ')
    const userMessage = `关于文件 ${fileNames} 的问题: ${docPrompt.value}`
    
    // 创建文件上传消息
    await aiStore.sendMessage(userMessage)
    
    // 关闭文档上传区域
    showDocUploader.value = false
    
    // 清空提示和文件列表
    docPrompt.value = ''
    uploadedFiles.value = []
    
    // 滚动到底部查看回复
    scrollToBottom()
  } catch (error) {
    console.error('处理文档提示失败:', error)
    ElMessage.error('处理文档提示失败，请重试')
  }
}

// 确认使用文件
const confirmDocUpload = async () => {
  if (uploadedFiles.value.length === 0) return
  
  try {
    isProcessingFile.value = true
    ElMessage.info('正在处理文件，请稍候...')
    
    // 解析所有上传的文件
    for (const file of uploadedFiles.value) {
      const content = await extractFileContent(file)
      if (content) {
        fileContents.value[file.name] = content
      }
    }
    
    // 保存文件名到活跃文件列表
    activeDocFiles.value = uploadedFiles.value.map(file => file.name)
    
    // 关闭文档上传区域
    showDocUploader.value = false
    
    // 提示用户现在可以在聊天框中提问关于文件的问题
    ElMessage.success(`已处理${uploadedFiles.value.length}个文件，您可以在聊天框中提问关于这些文件的问题`)
    
    // 在聊天框中增加一个提示消息
    userInput.value = `请问关于${activeDocFiles.value.join('、')}这些文件，`
    
    // 聚焦输入框
    nextTick(() => {
      const textarea = document.querySelector('.chat-input .el-textarea__inner')
      if (textarea) {
        textarea.focus()
      }
    })
    
    // 清空上传区域的文件列表，以便下次上传
    uploadedFiles.value = []
  } catch (error) {
    console.error('处理文件失败:', error)
    ElMessage.error('处理文件失败，请重试')
  } finally {
    isProcessingFile.value = false
  }
}

// 清除活跃文件
const clearActiveFiles = () => {
  activeDocFiles.value = []
  fileContents.value = {}
  ElMessage.info('已清除当前使用的文件')
}

// 提取文件内容
const extractFileContent = async (file) => {
  const extension = file.name.split('.').pop().toLowerCase()
  
  try {
    switch (extension) {
      case 'txt':
        return await readTextFile(file)
      
      case 'pdf':
        return await readPdfFile(file)
      
      case 'docx':
        return await readDocxFile(file)
      
      case 'xlsx':
      case 'xls':
        return await readExcelFile(file)
      
      case 'pptx':
        return await readPptxFile(file)
      
      case 'jpg':
      case 'jpeg':
      case 'png':
        return `[这是一张图片: ${file.name}，可以询问关于图片的问题]`
      
      default:
        return `[无法解析 ${file.name} 的内容，文件格式不支持]`
    }
  } catch (error) {
    console.error(`解析文件 ${file.name} 失败:`, error)
    return `[解析 ${file.name} 时发生错误: ${error.message}]`
  }
}

// 读取文本文件
const readTextFile = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.onerror = (e) => reject(new Error('读取文本文件失败'))
    reader.readAsText(file)
  })
}

// 读取PDF文件
const readPdfFile = async (file) => {
  try {
    const arrayBuffer = await file.arrayBuffer()
    const pdf = await getDocument({ data: arrayBuffer }).promise
    let content = ''
    
    for (let i = 1; i <= pdf.numPages; i++) {
      const page = await pdf.getPage(i)
      const textContent = await page.getTextContent()
      const pageText = textContent.items.map(item => item.str).join(' ')
      content += `--- 第${i}页 ---\n${pageText}\n\n`
      
      // 如果内容过多，只提取前几页
      if (i >= 10) {
        content += `[PDF文件过长，仅显示前${i}页内容]`
        break
      }
    }
    
    return content
  } catch (error) {
    console.error('解析PDF失败:', error)
    throw new Error('PDF解析失败')
  }
}

// 读取Docx文件
const readDocxFile = async (file) => {
  try {
    const arrayBuffer = await file.arrayBuffer();
    
    // 使用mammoth提取文本内容
    const result = await mammoth.extractRawText({ 
      arrayBuffer: arrayBuffer,
      includeDefaultStyleMap: true
    });
    
    let content = result.value || '';
    
    // 如果内容为空，提供友好提示
    if (!content.trim()) {
      return `[${file.name} 文件似乎没有可提取的文本内容]`;
    }
    
    // 尝试提取文档的元数据
    try {
      const zip = new JSZip();
      const docxData = await zip.loadAsync(arrayBuffer);
      
      // 尝试读取文档属性
      if (docxData.files['docProps/core.xml']) {
        const coreXml = await docxData.file('docProps/core.xml').async('string');
        const titleMatch = coreXml.match(/<dc:title>(.*?)<\/dc:title>/);
        const authorMatch = coreXml.match(/<dc:creator>(.*?)<\/dc:creator>/);
        
        let metaData = '--- 文档信息 ---\n';
        let hasMetaData = false;
        
        if (titleMatch && titleMatch[1]) {
          metaData += `标题: ${titleMatch[1]}\n`;
          hasMetaData = true;
        }
        
        if (authorMatch && authorMatch[1]) {
          metaData += `作者: ${authorMatch[1]}\n`;
          hasMetaData = true;
        }
        
        if (hasMetaData) {
          content = metaData + '\n' + content;
        }
      }
    } catch (metaError) {
      console.error('提取文档元数据失败:', metaError);
      // 忽略元数据错误，继续处理文本内容
    }
    
    return content;
  } catch (error) {
    console.error('解析Docx失败:', error);
    return `[解析 ${file.name} 时发生错误: ${error.message}]`;
  }
}

// 读取Excel文件
const readExcelFile = async (file) => {
  try {
    const arrayBuffer = await file.arrayBuffer()
    const workbook = XLSX.read(arrayBuffer, { type: 'array' })
    let content = ''
    
    // 遍历所有工作表
    for (const sheetName of workbook.SheetNames) {
      const worksheet = workbook.Sheets[sheetName]
      const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 })
      
      content += `--- 工作表: ${sheetName} ---\n`
      
      // 将数据格式化为表格样式的文本
      for (const row of jsonData) {
        content += row.join('\t') + '\n'
      }
      
      content += '\n'
    }
    
    return content
  } catch (error) {
    console.error('解析Excel失败:', error)
    throw new Error('Excel解析失败')
  }
}

// 读取Pptx文件
const readPptxFile = async (file) => {
  try {
    const arrayBuffer = await file.arrayBuffer();
    const zip = new JSZip();
    
    // 加载 pptx 文件 (实际上是 zip 文件)
    const pptxData = await zip.loadAsync(arrayBuffer);
    
    // 提取幻灯片内容
    let slideContent = '';
    let slideCount = 0;
    
    // 遍历 ppt/slides 目录下的所有文件
    const slidePromises = [];
    
    for (const filename in pptxData.files) {
      // 只处理幻灯片 XML 文件
      if (filename.startsWith('ppt/slides/slide') && filename.endsWith('.xml')) {
        const slideNumber = parseInt(filename.replace(/[^0-9]/g, ''));
        slidePromises.push(
          pptxData.file(filename).async('string').then(content => {
            return { slideNumber, content };
          })
        );
        slideCount++;
      }
    }
    
    // 等待所有幻灯片内容加载完成
    const slides = await Promise.all(slidePromises);
    
    // 按幻灯片顺序排序
    slides.sort((a, b) => a.slideNumber - b.slideNumber);
    
    // 从每张幻灯片提取文本
    for (const slide of slides) {
      slideContent += `--- 幻灯片 ${slide.slideNumber} ---\n`;
      
      // 从 XML 中提取文本内容 (简单提取 <a:t> 标签中的内容)
      const textMatches = slide.content.match(/<a:t>([^<]*)<\/a:t>/g);
      if (textMatches) {
        const extractedTexts = textMatches.map(match => {
          return match.replace(/<a:t>/, '').replace(/<\/a:t>/, '');
        });
        
        slideContent += extractedTexts.join('\n') + '\n\n';
      } else {
        slideContent += '[未找到文本内容]\n\n';
      }
    }
    
    if (slideCount === 0) {
      return `[未能从PPT文件中提取幻灯片内容]`;
    }
    
    return slideContent;
  } catch (error) {
    console.error('解析Pptx失败:', error);
    return `[解析 ${file.name} 时发生错误: ${error.message}]`;
  }
}

// 获取文件类型图标
const getFileTypeIcon = (fileName) => {
  const extension = fileName.split('.').pop().toLowerCase()
  switch (extension) {
    case 'txt':
      return 'Files'
    case 'pdf':
      return 'Reading'
    case 'docx':
      return 'Reading'
    case 'xlsx':
      return 'Files'
    case 'pptx':
      return 'Reading'
    case 'jpg':
    case 'jpeg':
    case 'png':
      return 'PictureFilled'
    default:
      return 'Files'
  }
}

// 获取文件类型类名
const getFileTypeClass = (fileName) => {
  const extension = fileName.split('.').pop().toLowerCase()
  switch (extension) {
    case 'txt':
      return 'text-file'
    case 'pdf':
      return 'pdf-file'
    case 'docx':
      return 'docx-file'
    case 'xlsx':
      return 'xlsx-file'
    case 'pptx':
      return 'pptx-file'
    case 'jpg':
    case 'jpeg':
    case 'png':
      return 'image-file'
    default:
      return 'other-file'
  }
}

// 跳转到AI设置页面
const goToAISettings = () => {
  // 关闭当前的图像生成区域
  showImageGenerator.value = false
  // 重置状态
  showApiKeyPrompt.value = false
  // 跳转到AI设置页面
  router.push('/ai/settings')
}
</script>

<style lang="scss" scoped>
.ai-chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8fafc;
  
  .ai-chat-header {
    padding: 16px 24px;
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    position: sticky;
    top: 0;
    z-index: 10;
    
    .model-selector {
      width: 240px;
      
      :deep(.el-input__wrapper) {
        background-color: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        box-shadow: none;
        
        &:hover {
          border-color: #8b5cf6;
        }
        
        &.is-focus {
          border-color: #8b5cf6;
          box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
        }
      }
    }
    
    .el-button {
      border-color: #e2e8f0;
      color: #475569;
        
        &:hover {
        border-color: #8b5cf6;
        color: #8b5cf6;
      }
    }
  }
  
  .quick-tools {
    padding: 16px 24px;
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    
    .tools-container {
      display: flex;
      gap: 8px;
      overflow-x: auto;
      padding-bottom: 6px;
    
    &::-webkit-scrollbar {
        height: 4px;
    }
    
    &::-webkit-scrollbar-track {
        background: #f1f5f9;
        border-radius: 4px;
    }
    
    &::-webkit-scrollbar-thumb {
        background: #cbd5e1;
        border-radius: 4px;
      
      &:hover {
          background: #94a3b8;
        }
      }
      
      .tool-button {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        border-radius: 20px;
        border: 1px solid #e2e8f0;
        background-color: #ffffff;
        color: #475569;
        font-size: 14px;
        white-space: nowrap;
        transition: all 0.2s ease;
        flex-shrink: 0;
        
        .el-icon {
          font-size: 16px;
        }
        
        &:hover {
          background-color: #f1f5f9;
          color: #8b5cf6;
          border-color: #cbd5e1;
        }
        
        &.more-button {
          background-color: #f8fafc;
        }
      }
    }
  }
  
  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    
    .empty-chat {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #94a3b8;
      
      .empty-icon {
        font-size: 48px;
        margin-bottom: 16px;
        color: #8b5cf6;
      }
      
      p {
        font-size: 16px;
      }
    }
    
    .message {
      display: flex;
      gap: 16px;
      margin-bottom: 24px;
      
      &.user {
        flex-direction: row-reverse;
      
      .message-content {
          align-items: flex-end;
          
          .message-container {
            flex-direction: row-reverse;
          }
          
          .message-text {
            background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
            color: #ffffff;
            border-radius: 16px 16px 0 16px;
            
            .user-message {
              white-space: pre-wrap;
            }
          }
          
          .message-actions {
            margin-right: 8px;
          }
          
          .message-time {
            text-align: right;
          }
        }
      }
      
      &.assistant {
        .message-text {
          background-color: #ffffff;
          border: 1px solid #e2e8f0;
          border-radius: 16px 16px 16px 0;
          color: #1e293b;
          
          .ai-message {
              :deep(pre) {
              background: #f1f5f9;
                border-radius: 8px;
                padding: 16px;
              margin: 8px 0;
                
                code {
                  font-family: 'Fira Code', monospace;
                font-size: 14px;
                line-height: 1.5;
                }
              }
              
              :deep(p) {
                margin: 8px 0;
              line-height: 1.6;
              }
              
              :deep(ul), :deep(ol) {
                margin: 8px 0;
              padding-left: 24px;
              }
              
              :deep(a) {
              color: #8b5cf6;
                text-decoration: none;
                
                &:hover {
                  text-decoration: underline;
                }
              }
          }
          
          .streaming-content {
            .typing-indicator {
              display: inline-flex;
              align-items: center;
              gap: 4px;
              margin-left: 8px;
              
              .dot {
                width: 4px;
                height: 4px;
                background-color: #8b5cf6;
                border-radius: 50%;
                animation: typing 1s infinite;
                
                &:nth-child(2) { animation-delay: 0.2s; }
                &:nth-child(3) { animation-delay: 0.4s; }
              }
            }
          }
        }
      }
      
      .message-avatar {
        flex-shrink: 0;
        
        .el-avatar {
          background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
          color: #ffffff;
          font-weight: 600;
        }
      }
      
      .message-content {
        display: flex;
        flex-direction: column;
        max-width: 80%;
        
        .message-container {
          display: flex;
          align-items: flex-start;
          gap: 8px;
        }
        
        .message-text {
          padding: 12px 16px;
          font-size: 15px;
          line-height: 1.6;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
          }
          
          .message-actions {
            opacity: 0;
            transition: opacity 0.2s;
            
            .el-button {
            padding: 6px;
            border-color: #e2e8f0;
              
              &:hover {
              border-color: #8b5cf6;
              color: #8b5cf6;
              }
            }
          }
          
          &:hover .message-actions {
            opacity: 1;
          }
        
        .message-time {
          font-size: 12px;
          color: #94a3b8;
          margin-top: 4px;
        }
        
        .message-edit {
          width: 100%;
          
          .el-input {
            margin-bottom: 8px;
            
            :deep(.el-textarea__inner) {
              border-color: #e2e8f0;
              border-radius: 8px;
              
              &:hover, &:focus {
                border-color: #8b5cf6;
              }
            }
          }
          
          .edit-actions {
            display: flex;
            justify-content: flex-end;
            gap: 8px;
            
            .el-button {
              &--primary {
                background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
                border: none;
          
          &:hover {
                  transform: translateY(-1px);
                  box-shadow: 0 4px 8px rgba(139, 92, 246, 0.2);
                }
              }
            }
          }
        }
      }
    }
  }
  
  .chat-input {
    padding: 24px;
    background: #ffffff;
    border-top: 1px solid #e2e8f0;
    
    .input-wrapper {
      display: flex;
      gap: 16px;
      margin-bottom: 12px;
      
      .el-input {
        :deep(.el-textarea__inner) {
          border-color: #e2e8f0;
          border-radius: 12px;
          padding: 12px 16px;
          min-height: 80px;
          resize: none;
          font-size: 15px;
          
          &:hover, &:focus {
            border-color: #8b5cf6;
          }
        }
      }
      
      .send-button {
        align-self: flex-end;
        height: 40px;
        padding: 0 24px;
        background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
        border: none;
        border-radius: 8px;
        font-size: 15px;
        font-weight: 500;
        
        &:hover:not(:disabled) {
          transform: translateY(-1px);
          box-shadow: 0 4px 8px rgba(139, 92, 246, 0.2);
        }
        
        &:disabled {
          background: #e2e8f0;
          opacity: 0.7;
        }
        
        .el-icon {
          margin-right: 4px;
        }
      }
    }
    
    .input-actions {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .el-button-group {
        .el-button {
          border-color: #e2e8f0;
          color: #475569;
          
          &:hover:not(:disabled) {
            border-color: #8b5cf6;
            color: #8b5cf6;
          }
          
          &:disabled {
            color: #cbd5e1;
            border-color: #e2e8f0;
            background-color: #f8fafc;
          }
        }
      }
    }
  }
}

@keyframes typing {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

.image-generation-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  
  .image-input-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .image-options {
    display: flex;
    gap: 12px;
    
    .el-select {
      width: 50%;
    }
  }
  
  .image-generate-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin: 8px 0;
    
    .image-generate-area {
      border: 1px dashed #e2e8f0;
      border-radius: 8px;
      padding: 24px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      position: relative;
      transition: all 0.3s ease;
      
      &:hover {
        background-color: #f8fafc;
        border-color: #8b5cf6;
      }
      
      .generate-inner {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;
        
        .generate-icon {
          font-size: 36px;
          color: #8b5cf6;
        }
        
        .generate-text {
          text-align: center;
          
          .generate-title {
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 8px;
            color: #1e293b;
          }
          
          .generate-desc {
            font-size: 14px;
            color: #94a3b8;
          }
        }
      }
    }
  }
  
  .generate-button-area {
    display: flex;
    justify-content: center;
    margin-top: 8px;
  }
  
  .generate-image-button {
    padding: 8px 16px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    background-color: #ffffff;
    color: #475569;
    font-size: 14px;
    cursor: pointer;
    
    &:hover {
      background-color: #f1f5f9;
      color: #8b5cf6;
      border-color: #8b5cf6;
    }
    
    .el-icon {
      margin-right: 4px;
    }
  }
  
  .generated-images {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
    margin-top: 16px;
    
    .image-item {
      border-radius: 8px;
      overflow: hidden;
      border: 1px solid #e2e8f0;
      
      img {
        width: 100%;
        height: auto;
        display: block;
      }
      
      .image-actions {
        display: flex;
        justify-content: space-between;
        padding: 8px;
        background: #f8fafc;
        
        .el-button {
          padding: 4px 8px;
          font-size: 12px;
        }
      }
    }
  }
  
  .generation-error {
    color: #dc2626;
    padding: 8px 12px;
    background-color: #fee2e2;
    border-radius: 6px;
    font-size: 14px;
    margin-top: 8px;
  }
}

.inline-doc-uploader {
  background-color: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 16px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  .uploader-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      color: #1e293b;
      margin: 0;
    }
  }
  
  .doc-upload-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 16px;
  }
  
  .doc-upload-area {
    border: 1px dashed #e2e8f0;
    border-radius: 8px;
    padding: 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    position: relative;
    transition: all 0.3s ease;
    flex: 1;
    
    &:hover {
      background-color: #f8fafc;
      border-color: #8b5cf6;
    }
    
    .file-input {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      opacity: 0;
      cursor: pointer;
      z-index: -1; /* 隐藏但保持功能 */
    }
  }
  
  .upload-button-area {
    display: flex;
    justify-content: center;
    margin-top: 8px;
  }
  
  .local-file-button {
    padding: 8px 16px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    background-color: #ffffff;
    color: #475569;
    font-size: 14px;
    cursor: pointer;
    
    &:hover {
      background-color: #f1f5f9;
      color: #8b5cf6;
      border-color: #8b5cf6;
    }
    
    .el-icon {
      margin-right: 4px;
    }
  }
  
  .uploaded-files {
    margin-bottom: 16px;
    max-height: 200px;
    overflow-y: auto;
    
    &::-webkit-scrollbar {
      width: 4px;
    }
    
    &::-webkit-scrollbar-track {
      background: #f1f5f9;
      border-radius: 4px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #cbd5e1;
      border-radius: 4px;
      
      &:hover {
        background: #94a3b8;
      }
    }
    
    .file-item {
      display: flex;
      align-items: center;
      padding: 10px;
      border-radius: 8px;
      background-color: #f8fafc;
      margin-bottom: 8px;
      
      .file-icon {
        width: 40px;
        height: 40px;
        background-color: #e2e8f0;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 12px;
        
        .el-icon {
          font-size: 20px;
          color: #64748b;
        }
        
        &.text-file {
          background-color: #e2e8f0;
          .el-icon {
            color: #64748b;
          }
        }
        
        &.pdf-file {
          background-color: #fee2e2;
          .el-icon {
            color: #b91c1c;
          }
        }
        
        &.docx-file {
          background-color: #dbeafe;
          .el-icon {
            color: #1d4ed8;
          }
        }
        
        &.xlsx-file {
          background-color: #dcfce7;
          .el-icon {
            color: #16a34a;
          }
        }
        
        &.pptx-file {
          background-color: #ffedd5;
          .el-icon {
            color: #c2410c;
          }
        }
        
        &.image-file {
          background-color: #f3e8ff;
          .el-icon {
            color: #9333ea;
          }
        }
        
        &.other-file {
          background-color: #e2e8f0;
          .el-icon {
            color: #64748b;
          }
        }
      }
      
      .file-info {
        flex: 1;
        
        .file-name {
          font-size: 14px;
          font-weight: 500;
          color: #1e293b;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          max-width: 300px;
        }
        
        .file-size {
          font-size: 12px;
          color: #94a3b8;
        }
      }
      
      .remove-file {
        color: #94a3b8;
        
        &:hover {
          color: #ef4444;
        }
      }
    }
  }
  
  .doc-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    
    .el-button {
      padding: 6px 12px;
      border-color: #e2e8f0;
      color: #475569;
      
      &:hover:not(:disabled) {
        border-color: #8b5cf6;
        color: #8b5cf6;
      }
      
      &:disabled {
        color: #cbd5e1;
        border-color: #e2e8f0;
        background-color: #f8fafc;
      }
    }
  }
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .stop-btn {
    font-size: 14px;
    padding: 6px 12px;
    height: 32px;
    border-color: #e2e8f0;
    color: #475569;
    
    &:hover {
      border-color: #ef4444;
      color: #ef4444;
    }
  }
}

.active-files-info {
  background-color: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
  
  .active-files-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .active-files-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 500;
      color: #0284c7;
      
      .el-icon {
        font-size: 16px;
      }
    }
    
    .el-button {
      color: #64748b;
      
      &:hover {
        color: #ef4444;
      }
    }
  }
  
  .active-files-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    
    .active-file-item {
      background-color: #e0f2fe;
      color: #0369a1;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 12px;
    }
  }
}

.inline-image-generator {
  background-color: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 16px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  .generator-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      color: #1e293b;
      margin: 0;
    }
  }
}

.api-key-prompt {
  background-color: #ffffff;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 16px;
  margin-bottom: 16px;
  
  .prompt-content {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .warning-icon {
      font-size: 20px;
      color: #ef4444;
    }
    
    .prompt-text {
      h4 {
        font-size: 16px;
        font-weight: 600;
        color: #1e293b;
        margin: 0;
      }
      
      p {
        font-size: 14px;
        color: #94a3b8;
      }
    }
  }
  
  .prompt-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    
    .el-button {
      &--primary {
        background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
        border: none;
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 4px 8px rgba(139, 92, 246, 0.2);
        }
      }
    }
  }
}
</style> 