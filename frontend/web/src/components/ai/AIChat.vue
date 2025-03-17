<template>
  <div class="ai-chat-container">
    <div class="ai-chat-header">
      <el-select v-model="selectedModel" placeholder="选择AI模型" class="model-selector">
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
      <div v-if="messages.length === 0" class="empty-chat">
        <el-icon class="empty-icon"><ChatDotRound /></el-icon>
        <p>开始与AI助手对话</p>
      </div>
      
      <div 
        v-for="message in messages" 
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
            :provider-id="getMessageProviderID(message)"
            :model-id="selectedModel"
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
                      {{ message.content }}
                      <span class="typing-indicator">
                        <span class="dot"></span>
                        <span class="dot"></span>
                        <span class="dot"></span>
                      </span>
                    </div>
                    <div v-else>
                      {{ message.content }}
                    </div>
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
import { marked } from 'marked'
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
  const username = userStore.user?.username || '用户'
  return username.charAt(0).toUpperCase()
})

// 当前模型提供商
const currentModelProvider = computed(() => {
  const model = aiStore.availableModels.find(m => m.id === aiStore.currentModel)
  return model ? model.provider : ''
})

// 计算属性
const messages = computed(() => aiStore.messages)
const isLoading = computed(() => aiStore.isLoading)
const currentModel = computed(() => aiStore.currentModel)
const currentProvider = computed(() => aiStore.currentProvider)
const canUndo = computed(() => aiStore.currentHistoryIndex > 0)
const canRedo = computed(() => aiStore.currentHistoryIndex < aiStore.messageHistory.length - 1)

// 修改 isEditing 方法
const isEditing = (messageId) => {
  return editingMessageId.value === messageId
}

// 初始化时设置默认模型
onMounted(() => {
  if (aiStore.defaultModel) {
    selectedModel.value = aiStore.defaultModel
  } else if (aiStore.availableModels.length > 0) {
    selectedModel.value = aiStore.availableModels[0].id
  }
  aiStore.setCurrentModel(selectedModel.value)
  
  // 添加滚动事件监听
  if (chatContainer.value) {
    chatContainer.value.addEventListener('scroll', () => {
      // 可以在这里添加滚动相关的逻辑
    })
  }
  
  scrollToBottom()
})

// 监听模型变化
watch(selectedModel, (newModel) => {
  aiStore.setCurrentModel(newModel)
})

// 监听消息变化，自动滚动到底部
watch(messages, () => {
  // 只在消息变化时检查是否需要滚动
  nextTick(() => {
    if (chatContainer.value) {
      const scrollHeight = chatContainer.value.scrollHeight
      const clientHeight = chatContainer.value.clientHeight
      const maxScrollTop = scrollHeight - clientHeight
      const isAtBottom = chatContainer.value.scrollTop + clientHeight >= scrollHeight - 50
      
      if (isAtBottom) {
        chatContainer.value.scrollTo({
          top: maxScrollTop,
          behavior: 'smooth'
        })
      }
    }
  })
}, { deep: true })

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

// 格式化消息内容
const formatMessage = (content) => {
  if (!content) return ''
  return md.render(content)
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

// 获取消息的提供商ID
const getMessageProviderID = (message) => {
  if (message.providerId) return message.providerId
  return currentModelProvider.value
}

// 发送消息
const sendMessage = async () => {
  if (!userInput.value.trim() || aiStore.isLoading) return
  
  const message = userInput.value
  userInput.value = '' // 清空输入框
  
  try {
    await aiStore.sendMessage(message)
  } catch (error) {
    ElMessage.error('获取AI回复失败')
    userInput.value = message // 如果发送失败，恢复输入内容
  }
}

// 清空对话
const clearChat = () => {
  aiStore.clearMessages()
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      const scrollHeight = chatContainer.value.scrollHeight
      const clientHeight = chatContainer.value.clientHeight
      const maxScrollTop = scrollHeight - clientHeight
      
      // 只有当用户不在查看历史消息时才自动滚动到底部
      const isAtBottom = chatContainer.value.scrollTop + clientHeight >= scrollHeight - 50
      if (isAtBottom) {
        chatContainer.value.scrollTo({
          top: maxScrollTop,
          behavior: 'smooth'
        })
      }
    }
  })
}

// 修改 startEditing 方法
const startEditing = (message) => {
  if (message.role === 'assistant' && message.isStreaming) {
    ElMessage.warning('AI正在回复中，请等待回复完成')
    return
  }
  
  editingMessageId.value = message.id
  editingContent.value = message.content
  aiStore.startEditing(message.id)
}

// 修改 cancelEditing 方法
const cancelEditing = () => {
  editingMessageId.value = null
  editingContent.value = ''
  aiStore.cancelEditing()
}

// 修改 saveEdit 方法
const saveEdit = async () => {
  if (!editingContent.value.trim()) {
    ElMessage.warning('消息内容不能为空')
    return
  }
  
  try {
    // 保存当前编辑的消息ID和内容
    const messageId = editingMessageId.value
    const content = editingContent.value
    
    // 重置编辑状态，让编辑框消失
    editingMessageId.value = null
    editingContent.value = ''
    
    // 然后保存编辑并重新发送消息
    await aiStore.saveEdit(messageId, content)
  } catch (error) {
    console.error('保存编辑失败:', error)
    ElMessage.error('保存编辑失败')
  }
}

const undo = () => {
  aiStore.undo()
}

const redo = () => {
  aiStore.redo()
}

// 处理停止生成
const handleStop = () => {
  aiStore.cancelCurrentRequest()
}
</script>

<style scoped lang="scss">
.ai-chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #171717;
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  position: relative;
  
  .ai-chat-header {
    flex: none;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    background: linear-gradient(180deg, #1f1f1f 0%, #1a1a1a 100%);
    border-bottom: 1px solid #2a2a2a;
    border-radius: 16px 16px 0 0;
    height: 64px;
    
    .model-selector {
      width: 240px;
      
      :deep(.el-input__wrapper) {
        background-color: #252525;
        border: 1px solid #2a2a2a;
        box-shadow: none;
        transition: all 0.3s ease;
        border-radius: 8px;
        
        &:hover {
          border-color: #2B5876;
        }
        
        &.is-focus {
          border-color: #2B5876;
          box-shadow: 0 0 0 3px rgba(43, 88, 118, 0.1);
        }
        
        .el-input__inner {
          color: #e0e0e0;
        }
      }
    }
    
    .el-button {
      background-color: #252525;
      border: 1px solid #2a2a2a;
      color: #a0a0a0;
      transition: all 0.3s ease;
      border-radius: 8px;
      padding: 8px 16px;
      
      &:hover {
        background-color: #2a2a2a;
        border-color: #2B5876;
        color: #2B5876;
      }
    }
  }
  
  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    background-color: #171717;
    border-radius: 16px;
    margin: 16px;
    margin-bottom: 0;
    height: calc(100vh - 64px - 160px - 32px);
    min-height: 0;
    position: relative;
    
    &::-webkit-scrollbar {
      width: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background: #1f1f1f;
      border-radius: 3px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #2a2a2a;
      border-radius: 3px;
      
      &:hover {
        background: #383838;
      }
    }
    
    .empty-chat {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #9ca3af;
      
      .empty-icon {
        font-size: 48px;
        margin-bottom: 16px;
        color: #2B5876;
      }
      
      p {
        font-size: 16px;
      }
    }
    
    .message {
      display: flex;
      gap: 16px;
      margin-bottom: 24px;
      
      .message-avatar {
        margin: 0 8px;
      }
      
      .message-content {
        flex: 1;
        background-color: #1f1f1f;
        padding: 20px;
        border-radius: 16px;
        border: 1px solid #2a2a2a;
        color: #e0e0e0;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        transition: all 0.3s ease;
        
        &:hover {
          border-color: #2B5876;
          box-shadow: 0 4px 12px rgba(43, 88, 118, 0.1);
        }
        
        .message-container {
          position: relative;
          padding-right: 40px;
          
          .message-text {
            line-height: 1.6;
            font-size: 14px;
            
            &.streaming {
              .typing-indicator {
                display: inline-flex;
                align-items: center;
                margin-left: 4px;
                padding: 0;
                
                .dot {
                  width: 4px;
                  height: 4px;
                  margin: 0 1px;
                  background-color: #2B5876;
                  border-radius: 50%;
                  animation: typing 1s infinite;
                  
                  &:nth-child(2) { animation-delay: 0.2s; }
                  &:nth-child(3) { animation-delay: 0.4s; }
                }
              }
            }
          }
          
          .message-actions {
            position: absolute;
            top: 0;
            right: 0;
            opacity: 0;
            transition: opacity 0.2s;
            
            .el-button {
              padding: 4px 8px;
              color: #a0a0a0;
              
              &:hover {
                color: #2B5876;
              }
            }
          }
          
          &:hover .message-actions {
            opacity: 1;
          }
        }
        
        .message-edit {
          width: 100%;
          
          .el-textarea {
            :deep(.el-textarea__inner) {
              background-color: #252525;
              border-color: #2a2a2a;
              color: #e0e0e0;
              border-radius: 8px;
              
              &:focus {
                border-color: #2B5876;
              }
            }
          }
          
          .edit-actions {
            display: flex;
            justify-content: flex-end;
            gap: 8px;
            margin-top: 8px;
          }
        }
        
        .message-time {
          margin-top: 8px;
          font-size: 12px;
          color: #666;
        }
      }
      
      &.user {
        .message-content {
          background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
          
          &:hover {
            border-color: #2B5876;
          }
          
          .user-message {
            white-space: pre-wrap;
            word-break: break-word;
            font-size: 14px;
            line-height: 1.6;
            color: #e0e0e0;
          }
        }
      }
    }
  }

  .chat-input {
    flex: none;
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 20px 24px;
    background: linear-gradient(0deg, #1f1f1f 0%, #1a1a1a 100%);
    border-top: 1px solid #2a2a2a;
    border-radius: 0 0 16px 16px;
    height: 160px;
    position: sticky;
    bottom: 0;
    z-index: 10;
    
    .input-wrapper {
      display: flex;
      gap: 12px;
      margin-bottom: 8px;
      
      .el-textarea {
        flex: 1;
        
        :deep(.el-textarea__inner) {
          background-color: #252525;
          border: 1px solid #2a2a2a;
          color: #e0e0e0;
          box-shadow: none;
          transition: all 0.3s ease;
          height: 80px;
          border-radius: 8px;
          
          &:hover {
            border-color: #2B5876;
          }
          
          &:focus {
            border-color: #2B5876;
            box-shadow: 0 0 0 3px rgba(43, 88, 118, 0.1);
          }
          
          &::placeholder {
            color: #666;
          }
        }
      }
      
      .send-button {
        align-self: flex-end;
        height: 40px;
        background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
        border: none;
        padding: 0 24px;
        transition: all 0.3s ease;
        border-radius: 8px;
        
        &:hover:not(:disabled) {
          transform: translateY(-1px);
          box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
        }
        
        &:disabled {
          opacity: 0.5;
          background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
        }
        
        .el-icon {
          margin-right: 8px;
        }
      }
    }
    
    .input-actions {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;
    }
  }
}

.typing-indicator {
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
  padding: 0;
  
  .dot {
    width: 4px;
    height: 4px;
    margin: 0 1px;
    background-color: #e0e0e0;
    border-radius: 50%;
    animation: typing 1s infinite;
    
    &:nth-child(2) { animation-delay: 0.2s; }
    &:nth-child(3) { animation-delay: 0.4s; }
  }
}

@keyframes typing {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

.message-text {
  &.streaming {
    .typing-indicator {
      display: inline-flex;
      align-items: center;
      margin-left: 4px;
      padding: 0;
      
      .dot {
        width: 4px;
        height: 4px;
        margin: 0 1px;
        background-color: #e0e0e0;
        border-radius: 50%;
        animation: typing 1s infinite;
        
        &:nth-child(2) { animation-delay: 0.2s; }
        &:nth-child(3) { animation-delay: 0.4s; }
      }
    }
  }
}

@keyframes typing {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

.message {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  
  .message-avatar {
    margin: 0 8px;
  }
  
  .message-content {
    flex: 1;
    background-color: #1f1f1f;
    padding: 16px;
    border-radius: 12px;
    border: 1px solid #2a2a2a;
    color: #e0e0e0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
    
    &:hover {
      border-color: #2B5876;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    
    .message-container {
      position: relative;
      padding-right: 40px;
      
      .message-text {
        line-height: 1.6;
        font-size: 14px;
        
        &.streaming {
          .typing-indicator {
            display: inline-flex;
            align-items: center;
            margin-left: 4px;
            padding: 0;
            
            .dot {
              width: 4px;
              height: 4px;
              margin: 0 1px;
              background-color: #e0e0e0;
              border-radius: 50%;
              animation: typing 1s infinite;
              
              &:nth-child(2) { animation-delay: 0.2s; }
              &:nth-child(3) { animation-delay: 0.4s; }
            }
          }
        }
      }
      
      .message-actions {
        position: absolute;
        top: 0;
        right: 0;
        opacity: 0;
        transition: opacity 0.2s;
        
        .el-button {
          padding: 4px 8px;
          color: #a0a0a0;
          
          &:hover {
            color: #409EFF;
          }
        }
      }
      
      &:hover .message-actions {
        opacity: 1;
      }
    }
    
    .message-edit {
      width: 100%;
      
      .el-textarea {
        :deep(.el-textarea__inner) {
          background-color: #252525;
          border-color: #2a2a2a;
          color: #e0e0e0;
          
          &:focus {
            border-color: #409EFF;
          }
        }
      }
      
      .edit-actions {
        display: flex;
        justify-content: flex-end;
        gap: 8px;
        margin-top: 8px;
      }
    }
    
    .message-time {
      margin-top: 8px;
      font-size: 12px;
      color: #666;
    }
  }
  
  &.user {
    .message-content {
      background: linear-gradient(135deg, #2B5876 0%, #4E4376 100%);
      
      &:hover {
        border-color: #409EFF;
      }
      
      .user-message {
        white-space: pre-wrap;
        word-break: break-word;
        font-size: 14px;
        line-height: 1.6;
        color: #e0e0e0;
      }
    }
  }
}

.markdown-body {
  color: #e0e0e0;
  font-size: 14px;
  line-height: 1.6;
  
  :deep(p) {
    margin: 8px 0;
  }
  
  :deep(pre) {
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
  
  :deep(code) {
    background-color: #252525;
    padding: 2px 4px;
    border-radius: 4px;
    font-family: 'Fira Code', monospace;
  }
  
  :deep(a) {
    color: #409EFF;
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
  
  :deep(ul), :deep(ol) {
    padding-left: 20px;
    margin: 8px 0;
  }
  
  :deep(blockquote) {
    margin: 8px 0;
    padding: 8px 16px;
    border-left: 4px solid #409EFF;
    background-color: #252525;
    
    p {
      margin: 0;
    }
  }
  
  :deep(table) {
    border-collapse: collapse;
    width: 100%;
    margin: 8px 0;
    
    th, td {
      border: 1px solid #383838;
      padding: 8px;
      text-align: left;
    }
    
    th {
      background-color: #252525;
    }
    
    tr:nth-child(even) {
      background-color: #1f1f1f;
    }
  }
}

:deep(.hljs) {
  background: #1a1a1a !important;
  border-radius: 4px;
  padding: 16px !important;
  
  .hljs-keyword,
  .hljs-selector-tag,
  .hljs-title,
  .hljs-section,
  .hljs-doctag,
  .hljs-name,
  .hljs-strong {
    color: #c678dd;
  }

  .hljs-string {
    color: #98c379;
  }

  .hljs-comment {
    color: #5c6370;
  }

  .hljs-number,
  .hljs-literal {
    color: #d19a66;
  }

  .hljs-attribute,
  .hljs-attr {
    color: #e06c75;
  }

  .hljs-variable,
  .hljs-template-variable,
  .hljs-tag,
  .hljs-name,
  .hljs-selector-id,
  .hljs-selector-class {
    color: #61afef;
  }
}

.streaming-content {
  display: inline-block;
  white-space: pre-wrap;
  
  .typing-indicator {
    display: inline-flex;
    align-items: center;
    margin-left: 4px;
    padding: 0;
    
    .dot {
      width: 4px;
      height: 4px;
      margin: 0 1px;
      background-color: #e0e0e0;
      border-radius: 50%;
      animation: typing 1s infinite;
      
      &:nth-child(2) { animation-delay: 0.2s; }
      &:nth-child(3) { animation-delay: 0.4s; }
    }
  }
}

@keyframes typing {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
</style> 