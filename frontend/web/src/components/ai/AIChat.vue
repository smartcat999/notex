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

    <div class="chat-messages" ref="chatContainer">
      <div v-if="orderedMessages.length === 0" class="empty-chat">
        <el-icon class="empty-icon"><ChatDotRound /></el-icon>
        <p>开始与AI助手对话</p>
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
        <el-button-group>
          <el-button :disabled="!canUndo || aiStore.isLoading" @click="undo">
            <el-icon><Back /></el-icon>
          </el-button>
          <el-button :disabled="!canRedo || aiStore.isLoading" @click="redo">
            <el-icon><Right /></el-icon>
          </el-button>
        </el-button-group>
        <el-button 
          v-if="aiStore.isLoading"
          @click="handleStop"
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
import { ChatDotRound, Position, Delete, Edit, Back, Right } from '@element-plus/icons-vue'
import AIAvatar from '@/components/ai/AIAvatar.vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

const userStore = useUserStore()
const aiStore = useAIStore()
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
    await aiStore.sendMessage(message)
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
</style> 