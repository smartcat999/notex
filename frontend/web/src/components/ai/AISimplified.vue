<template>
  <div class="ai-simplified-container">
    <div class="ai-header">
      <div class="model-selector-area">
        <h2>AI 助手</h2>
        <el-select v-model="selectedModel" placeholder="选择AI模型" class="model-selector" :loading="aiStore.isLoading">
          <el-option
            v-for="model in aiStore.availableModels"
            :key="model.id"
            :label="model.name"
            :value="model.id"
          />
        </el-select>
      </div>
      <el-button type="danger" @click="clearChat" class="clear-button" plain>
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
          <div
            v-else 
            class="ai-avatar"
            :style="{
              background: getAvatarBackground(message.providerId || aiStore.currentProvider)
            }"
          >
            AI
          </div>
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
      <div class="input-container">
        <el-input
          v-model="userInput"
          type="textarea"
          :rows="3"
          placeholder="输入您的问题..."
          resize="none"
          @keydown.enter.prevent="sendMessage"
          :disabled="aiStore.isLoading"
          class="message-textarea"
        />
        <div class="input-actions">
          <el-button 
            v-if="aiStore.isLoading"
            @click="handleStop"
            class="stop-button"
            type="warning"
          >
            <el-icon><CircleClose /></el-icon>
            停止
          </el-button>
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
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useAIStore } from '@/stores/ai'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Delete, Edit, Position, ChatDotRound, CircleClose } from '@element-plus/icons-vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const aiStore = useAIStore()
const userStore = useUserStore()
const chatContainer = ref(null)
const selectedModel = ref('')
const userInput = ref('')
const editingContent = ref('')
const userAvatar = computed(() => userStore.user?.avatar || '')
const userInitials = computed(() => {
  const name = userStore.user?.name || ''
  return name.substring(0, 2).toUpperCase()
})

// 将消息按照时间戳排序
const orderedMessages = computed(() => {
  return [...aiStore.messages].sort((a, b) => a.timestamp - b.timestamp)
})

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

// 发送消息
const sendMessage = async () => {
  if (!userInput.value.trim() || aiStore.isLoading) return
  
  const content = userInput.value
  userInput.value = ''
  
  // 保存当前型号到AI Store
  if (selectedModel.value) {
    aiStore.setCurrentModel(selectedModel.value)
  }
  
  try {
    await aiStore.sendMessage(content)
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error('发送消息失败')
  }
}

// 停止生成
const handleStop = () => {
  aiStore.cancelCurrentRequest()
}

// 清空聊天
const clearChat = () => {
  aiStore.clearMessages()
}

// 开始编辑消息
const startEditing = (message) => {
  editingContent.value = message.content
  aiStore.startEditing(message.id)
}

// 取消编辑
const cancelEditing = () => {
  editingContent.value = ''
  aiStore.cancelEditing()
}

// 保存编辑
const saveEdit = async () => {
  if (!editingContent.value.trim()) return
  
  const messageId = aiStore.editingMessageId
  if (!messageId) return
  
  try {
    await aiStore.saveEdit(messageId, editingContent.value)
    editingContent.value = ''
    
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('保存编辑失败:', error)
    ElMessage.error('保存编辑失败')
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

// 检查消息是否正在编辑
const isEditing = (messageId) => {
  return aiStore.isEditing(messageId)
}

// 渲染Markdown
const renderMarkdown = (content) => {
  if (!content) return ''
  
  try {
    // 配置marked
    marked.setOptions({
      highlight: function(code, lang) {
        if (lang && hljs.getLanguage(lang)) {
          return hljs.highlight(code, { language: lang }).value
        }
        return hljs.highlightAuto(code).value
      },
      breaks: true,
      gfm: true
    })
    
    // 渲染markdown
    const renderedContent = marked(content)
    
    // 在下一个tick更新代码块样式
    nextTick(() => {
      updateCodeSections()
    })
    
    return renderedContent
  } catch (error) {
    console.error('渲染Markdown失败:', error)
    return `<p>${content}</p>`
  }
}

// 更新代码块
const updateCodeSections = () => {
  const codeBlocks = document.querySelectorAll('.ai-simplified-container pre code')
  codeBlocks.forEach((block) => {
    // 如果已经高亮过，则跳过
    if (block.classList.contains('hljs')) return
    hljs.highlightElement(block)
    
    // 获取语言
    const className = Array.from(block.classList).find(cls => cls.startsWith('language-'))
    const language = className ? className.replace('language-', '') : ''
    
    // 添加语言标签
    const pre = block.parentElement
    if (language && !pre.querySelector('.language-label')) {
      pre.setAttribute('data-language', language)
      const languageLabel = document.createElement('span')
      languageLabel.className = 'language-label'
      languageLabel.textContent = language
      pre.appendChild(languageLabel)
    }
    
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
        navigator.clipboard.writeText(block.textContent).then(() => {
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
  })
}

// 获取不同提供商的头像背景色
const getAvatarBackground = (providerId) => {
  const colors = {
    'openai': 'linear-gradient(135deg, #10a37f, #0a8a6c)',
    'anthropic': 'linear-gradient(135deg, #b73999, #8f2d79)',
    'google': 'linear-gradient(135deg, #4285f4, #34a853)',
    'deepseek': 'linear-gradient(135deg, #ff5c35, #d64c2d)',
    'custom': 'linear-gradient(135deg, #8b5cf6, #6366f1)',
  }
  return colors[providerId] || colors.custom
}

// 监听消息变化，更新代码块和滚动位置
watch(() => aiStore.messages.length, () => {
  nextTick(() => {
    updateCodeSections()
    scrollToBottom()
  })
})

// 监听消息内容变化，更新代码块
watch(() => 
  aiStore.messages.map(m => m.content).join(''), 
  () => {
    nextTick(() => {
      updateCodeSections()
    })
  }
)

onMounted(async () => {
  // 如果AIStore已初始化，直接使用当前模型
  if (aiStore.initialized) {
    selectedModel.value = aiStore.currentModel
  } 
  // 如果AIStore有可用模型但未初始化
  else if (aiStore.availableModels.length > 0) {
    if (aiStore.defaultModel) {
      selectedModel.value = aiStore.defaultModel
    } else {
      selectedModel.value = aiStore.availableModels[0].id
    }
  }
  // 如果AIStore未初始化，等待初始化
  else {
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
    selectedModel.value = aiStore.currentModel
  }
})

// 监听AI存储默认模型变化
watch(() => aiStore.defaultModel, (newDefaultModel) => {
  if (newDefaultModel && (!selectedModel.value || selectedModel.value !== newDefaultModel)) {
    selectedModel.value = newDefaultModel
  }
})

// 监听AI存储当前模型变化
watch(() => aiStore.currentModel, (newCurrentModel) => {
  if (newCurrentModel && selectedModel.value !== newCurrentModel) {
    selectedModel.value = newCurrentModel
  }
})

// 监听模型选择变化
watch(selectedModel, async (newValue, oldValue) => {
  if (newValue && newValue !== oldValue) {
    try {
      aiStore.setCurrentModel(newValue)
      await aiStore.saveDefaultModel(newValue)
    } catch (error) {
      console.error('保存模型设置失败:', error)
    }
  }
})
</script>

<style lang="scss" scoped>
.ai-simplified-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8fafc;
  color: #1e293b;
  
  .ai-header {
    padding: 16px 28px;
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    position: sticky;
    top: 0;
    z-index: 10;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
    
    .model-selector-area {
      display: flex;
      align-items: center;
      gap: 16px;
      
      h2 {
        margin: 0;
        font-size: 1.5rem;
        font-weight: 600;
        background: linear-gradient(135deg, #0f172a 0%, #334155 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: 0.5px;
      }
    }
    
    .model-selector {
      width: 240px;
      
      :deep(.el-input__wrapper) {
        background-color: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        box-shadow: none;
        
        &:hover {
          border-color: #94a3b8;
        }
        
        &.is-focus {
          border-color: #334155;
          box-shadow: 0 0 0 3px rgba(51, 65, 85, 0.1);
        }
      }
    }
    
    .clear-button {
      height: 40px;
      border-radius: 8px;
      border: none;
      background-color: rgba(220, 38, 38, 0.05);
      color: #dc2626;
      font-weight: 500;
      padding: 0 16px;
      transition: all 0.2s ease;
      
      &:hover {
        background-color: rgba(220, 38, 38, 0.1);
        color: #b91c1c;
      }
      
      &:active {
        background-color: rgba(220, 38, 38, 0.15);
      }
    }
  }

  .chat-messages {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 24px;
    
    .empty-chat {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: #94a3b8;
      
      .empty-icon {
        font-size: 64px;
        margin-bottom: 16px;
        opacity: 0.5;
      }
      
      p {
        font-size: 18px;
        opacity: 0.8;
      }
    }
    
    .message {
      display: flex;
      gap: 16px;
      max-width: 85%;
      
      &.user {
        align-self: flex-end;
        flex-direction: row-reverse;
        
        .message-content {
          align-items: flex-end;
          
          .message-text {
            background-color: #2B5876;
            color: white;
            border-radius: 12px 2px 12px 12px;
            
            .user-message {
              white-space: pre-wrap;
            }
          }
          
          .message-time {
            text-align: right;
          }
        }
      }
      
      &.assistant {
        align-self: flex-start;
        
        .message-content {
          align-items: flex-start;
          
          .message-text {
            background-color: #ffffff;
            color: #1e293b;
            border-radius: 2px 12px 12px 12px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
            
            .ai-message {
              :deep(p) {
                margin: 0.5em 0;
                
                &:first-child {
                  margin-top: 0;
                }
                
                &:last-child {
                  margin-bottom: 0;
                }
              }
              
              :deep(pre) {
                margin: 1em 0;
                padding: 16px;
                background-color: #f8fafc;
                border-radius: 8px;
                border: 1px solid #e2e8f0;
                position: relative;
                
                code {
                  font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace;
                  font-size: 0.9em;
                  line-height: 1.5;
                  tab-size: 2;
                  padding-right: 70px;
                }
                
                .language-label {
                  position: absolute;
                  top: 8px;
                  left: 16px;
                  font-size: 0.7em;
                  color: #64748b;
                  font-weight: 500;
                  text-transform: uppercase;
                  letter-spacing: 0.05em;
                  background: #f1f5f9;
                  padding: 2px 8px;
                  border-radius: 4px;
                  border: 1px solid #e2e8f0;
                  z-index: 2;
                }
                
                .copy-button {
                  position: absolute;
                  top: 8px;
                  right: 8px;
                  height: 24px;
                  min-width: 60px;
                  padding: 0 8px;
                  font-size: 0.75em;
                  color: #64748b;
                  background: #f8fafc;
                  border: 1px solid #e2e8f0;
                  border-radius: 4px;
                  cursor: pointer;
                  transition: all 0.2s ease;
                  display: flex;
                  align-items: center;
                  justify-content: center;
                  gap: 4px;
                  
                  &:hover {
                    background: #f1f5f9;
                    color: #0f172a;
                  }
                  
                  &.copied {
                    background: #2B5876;
                    color: white;
                  }
                }
              }
              
              :deep(ul), :deep(ol) {
                padding-left: 1.5em;
                margin: 1em 0;
              }
              
              :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
                margin: 1.25em 0 0.5em;
                font-weight: 600;
                
                &:first-child {
                  margin-top: 0;
                }
              }
              
              :deep(blockquote) {
                border-left: 4px solid #2B5876;
                padding-left: 16px;
                margin: 1em 0;
                color: #475569;
              }
              
              :deep(a) {
                color: #2B5876;
                text-decoration: none;
                
                &:hover {
                  text-decoration: underline;
                }
              }
              
              :deep(img) {
                max-width: 100%;
                border-radius: 8px;
              }
              
              :deep(table) {
                border-collapse: collapse;
                width: 100%;
                margin: 1em 0;
                
                th, td {
                  border: 1px solid #e2e8f0;
                  padding: 8px 12px;
                  text-align: left;
                }
                
                th {
                  background-color: #f8fafc;
                  font-weight: 500;
                }
                
                tr:nth-child(even) {
                  background-color: #f8fafc;
                }
              }
            }
            
            .streaming-content {
              .typing-indicator {
                display: inline-flex;
                align-items: center;
                gap: 4px;
                height: 16px;
                margin-left: 4px;
                
                .dot {
                  width: 4px;
                  height: 4px;
                  background-color: #64748b;
                  border-radius: 50%;
                  animation: typing 1.5s infinite ease-in-out;
                  
                  &:nth-child(2) {
                    animation-delay: 0.2s;
                  }
                  
                  &:nth-child(3) {
                    animation-delay: 0.4s;
                  }
                }
              }
            }
          }
        }
      }
      
      .message-avatar {
        flex-shrink: 0;
        
        .ai-avatar {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-weight: 600;
          font-size: 14px;
        }
      }
      
      .message-content {
        display: flex;
        flex-direction: column;
        gap: 4px;
        
        .message-container {
          display: flex;
          gap: 8px;
          align-items: flex-start;
          width: 100%;
        }
        
        .message-text {
          padding: 12px 16px;
          font-size: 15px;
          line-height: 1.5;
          max-width: 100%;
          overflow-x: auto;
        }
        
        .message-actions {
          opacity: 0;
          transition: opacity 0.2s ease;
        }
        
        &:hover .message-actions {
          opacity: 1;
        }
        
        .message-edit {
          width: 100%;
          
          .edit-actions {
            display: flex;
            justify-content: flex-end;
            gap: 8px;
            margin-top: 8px;
          }
        }
        
        .message-time {
          font-size: 12px;
          color: #94a3b8;
          margin-top: 4px;
        }
      }
    }
  }
  
  .chat-input {
    padding: 20px 28px;
    background: #ffffff;
    border-top: 1px solid #e2e8f0;
    box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.03);
    
    .input-container {
      display: flex;
      flex-direction: column;
      gap: 12px;
      
      .message-textarea {
        width: 100%;
        
        :deep(.el-textarea__inner) {
          resize: none;
          border: 1px solid #e2e8f0;
          border-radius: 8px;
          padding: 14px 16px;
          min-height: 24px;
          max-height: 200px;
          line-height: 1.5;
          font-size: 15px;
          
          &:focus {
            border-color: #334155;
            box-shadow: 0 0 0 3px rgba(51, 65, 85, 0.1);
          }
        }
      }
      
      .input-actions {
        display: flex;
        justify-content: flex-end;
        gap: 12px;
        
        .stop-button {
          height: 40px;
          padding: 0 16px;
          background-color: rgba(234, 88, 12, 0.1);
          color: #ea580c;
          border: none;
          display: flex;
          align-items: center;
          gap: 4px;
          border-radius: 8px;
          font-size: 14px;
          font-weight: 500;
          
          .el-icon {
            font-size: 16px;
          }
          
          &:hover {
            background-color: rgba(234, 88, 12, 0.15);
            color: #c2410c;
          }
          
          &:active {
            background-color: rgba(234, 88, 12, 0.2);
          }
        }
        
        .send-button {
          height: 40px;
          min-width: 84px;
          padding: 0 16px;
          background: #0f172a;
          border: none;
          border-radius: 8px;
          box-shadow: 0 4px 12px rgba(15, 23, 42, 0.2);
          transition: all 0.2s ease;
          font-weight: 500;
          font-size: 14px;
          letter-spacing: 0.3px;
          display: flex;
          align-items: center;
          gap: 4px;
          
          .el-icon {
            font-size: 16px;
          }
          
          &:hover:not(:disabled) {
            background: #1e293b;
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(15, 23, 42, 0.25);
          }
          
          &:active:not(:disabled) {
            transform: translateY(0);
            box-shadow: 0 4px 12px rgba(15, 23, 42, 0.12);
          }
          
          &:disabled {
            background: #94a3b8;
            opacity: 0.8;
            box-shadow: none;
          }
        }
      }
    }
  }
}

@keyframes typing {
  0%, 100% {
    transform: translateY(0);
    opacity: 0.5;
  }
  50% {
    transform: translateY(-4px);
    opacity: 1;
  }
}

:deep(.el-button--primary) {
  background: #0f172a;
  border: none;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.2);
  font-weight: 500;
  height: 40px;
  
  &:hover:not(:disabled) {
    background: #1e293b;
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(15, 23, 42, 0.25);
  }
  
  &:active:not(:disabled) {
    transform: translateY(0);
    box-shadow: 0 4px 8px rgba(15, 23, 42, 0.12);
  }
  
  &:disabled {
    background: #94a3b8;
    opacity: 0.8;
    box-shadow: none;
  }
}
</style> 