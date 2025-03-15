import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import aiService from '@/services/aiService'
import { ElMessage } from 'element-plus'

export const useAIStore = defineStore('ai', () => {
  // 当前会话消息
  const messages = ref([])
  
  // 当前选择的模型
  const currentModel = ref('')
  
  // 加载状态
  const isLoading = ref(false)
  
  // 可用的AI模型
  const availableModels = reactive([
    { id: 'gpt-3.5-turbo', name: 'GPT-3.5 Turbo', provider: 'openai' },
    { id: 'gpt-4', name: 'GPT-4', provider: 'openai' },
    { id: 'claude-3-opus', name: 'Claude 3 Opus', provider: 'anthropic' },
    { id: 'claude-3-sonnet', name: 'Claude 3 Sonnet', provider: 'anthropic' },
    { id: 'gemini-pro', name: 'Gemini Pro', provider: 'google' },
    { id: 'deepseek-chat', name: 'DeepSeek Chat', provider: 'deepseek' },
    { id: 'deepseek-coder', name: 'DeepSeek Coder', provider: 'deepseek' },
    { id: 'custom-model', name: '自定义模型', provider: 'custom' }
  ])
  
  // 添加历史记录相关状态
  const messageHistory = ref([])
  const currentHistoryIndex = ref(-1)
  const isEditing = ref(false)
  const editingMessageId = ref(null)
  
  // 添加消息到历史记录
  const addToHistory = (messageState) => {
    // 如果当前不在最新位置，删除当前位置之后的所有历史记录
    if (currentHistoryIndex.value < messageHistory.value.length - 1) {
      messageHistory.value = messageHistory.value.slice(0, currentHistoryIndex.value + 1)
    }
    messageHistory.value.push([...messageState])
    currentHistoryIndex.value = messageHistory.value.length - 1
  }
  
  // 撤销上一步
  const undo = () => {
    if (currentHistoryIndex.value > 0) {
      currentHistoryIndex.value--
      // 恢复到历史记录中的状态
      messages.value = [...messageHistory.value[currentHistoryIndex.value]]
    }
  }
  
  // 重做
  const redo = () => {
    if (currentHistoryIndex.value < messageHistory.value.length - 1) {
      currentHistoryIndex.value++
      // 恢复到历史记录中的状态
      messages.value = [...messageHistory.value[currentHistoryIndex.value]]
    }
  }
  
  // 当前请求的控制器
  const currentController = ref(null)

  // 取消当前请求
  const cancelCurrentRequest = () => {
    if (currentController.value) {
      currentController.value.abort()
      currentController.value = null
      isLoading.value = false
    }
  }
  
  // 开始编辑消息
  const startEditing = (messageId) => {
    // 如果正在请求中，取消请求
    cancelCurrentRequest()
    
    isEditing.value = true
    editingMessageId.value = messageId
  }
  
  // 取消编辑
  const cancelEditing = () => {
    isEditing.value = false
    editingMessageId.value = null
  }
  
  // 修改 saveEdit 方法
  const saveEdit = async (messageId, content) => {
    const messageIndex = messages.value.findIndex(m => m.id === messageId)
    if (messageIndex === -1) return
    
    // 更新消息内容
    messages.value[messageIndex].content = content
    
    // 检查是否是最后一条用户消息
    const isLastUserMessage = messages.value
      .filter(m => m.role === 'user')
      .slice(-1)[0]?.id === messageId
    
    if (isLastUserMessage) {
      // 如果存在对应的 AI 回复，则需要删除
      if (messageIndex < messages.value.length - 1) {
        // 删除 AI 回复
        messages.value = messages.value.slice(0, messageIndex + 1)
      }
      
      // 更新历史记录
      addToHistory([...messages.value])
      
      // 重新发送请求
      try {
        await sendMessage(content, messageIndex)
      } catch (error) {
        console.error('重新发送消息失败:', error)
        throw error
      }
    } else {
      // 如果不是最后一条消息，直接更新历史记录
      addToHistory([...messages.value])
    }
    
    // 重置编辑状态
    isEditing.value = false
    editingMessageId.value = null
  }
  
  // 添加 currentProvider 计算属性
  const currentProvider = computed(() => {
    const model = availableModels.find(m => m.id === currentModel.value)
    return model ? model.provider : ''
  })
  
  // 修改 sendMessage 方法，添加可选的 messageIndex 参数
  const sendMessage = async (content, existingMessageIndex = -1, onProgress = null) => {
    try {
      // 验证消息内容
      if (!content || !content.trim()) {
        throw new Error('消息内容不能为空')
      }

      // 验证是否选择了模型
      if (!currentModel.value) {
        throw new Error('请先选择AI模型')
      }

      // 创建或更新消息
      let message
      if (existingMessageIndex === -1) {
        message = {
          id: Date.now(),
          role: 'user',
          content: content.trim(),
          timestamp: new Date().toISOString()
        }
        messages.value = [message] // 只保留当前消息
      } else {
        message = messages.value[existingMessageIndex]
        message.content = content.trim()
        message.timestamp = new Date().toISOString()
        messages.value = messages.value.slice(0, existingMessageIndex + 1) // 只保留到当前消息
      }

      isLoading.value = true

      // 创建 AbortController
      currentController.value = new AbortController()

      // 准备消息数组
      const messageArray = [{
        role: 'user',
        content: content.trim()
      }]

      const response = await aiService.sendChat(
        messageArray,
        currentModel.value,
        currentProvider.value,
        currentController.value.signal,
        (chunk) => {
          // 创建或更新 AI 响应消息
          let aiMessage = messages.value.find(msg => msg.role === 'assistant')
          if (!aiMessage) {
            aiMessage = {
              id: Date.now() + 1,
              role: 'assistant',
              content: '',
              timestamp: new Date().toISOString()
            }
            messages.value.push(aiMessage)
          }
          
          if (chunk) {
            aiMessage.content += chunk
            if (onProgress) {
              onProgress(aiMessage.content)
            }
          }
        }
      )

      if (response) {
        // 确保 AI 响应消息存在
        let aiMessage = messages.value.find(msg => msg.role === 'assistant')
        if (!aiMessage) {
          aiMessage = {
            id: Date.now() + 1,
            role: 'assistant',
            content: '',
            timestamp: new Date().toISOString()
          }
          messages.value.push(aiMessage)
        }
        aiMessage.content = response
      }
    } catch (error) {
      if (error.name === 'AbortError') {
        // 请求被取消，移除 AI 响应消息
        messages.value = messages.value.filter(msg => msg.role !== 'assistant')
        throw error
      }
      throw error
    } finally {
      isLoading.value = false
      currentController.value = null
    }
  }
  
  // 清空对话
  const clearMessages = () => {
    messages.value = []
    messageHistory.value = []
    currentHistoryIndex.value = -1
  }
  
  // 设置当前模型
  const setCurrentModel = (modelId) => {
    currentModel.value = modelId
  }
  
  // 获取模型名称
  const getModelName = (modelId) => {
    const model = availableModels.find(m => m.id === modelId)
    return model ? model.name : modelId
  }
  
  // 获取模型提供商
  const getModelProvider = (modelId) => {
    const model = availableModels.find(m => m.id === modelId)
    return model ? model.provider : 'unknown'
  }
  
  // 测试AI提供商连接
  const testProviderConnection = async (providerId) => {
    return await aiService.testProviderConnection(providerId)
  }
  
  // 默认模型
  const defaultModel = ref('')

  // 从 localStorage 读取默认模型
  const loadDefaultModel = () => {
    const saved = localStorage.getItem('defaultAIModel')
    if (saved) {
      defaultModel.value = saved
      currentModel.value = saved
    }
  }

  // 保存默认模型
  const saveDefaultModel = (modelId) => {
    defaultModel.value = modelId
    localStorage.setItem('defaultAIModel', modelId)
  }

  // 初始化时加载默认模型
  loadDefaultModel()
  
  return {
    messages,
    currentModel,
    isLoading,
    availableModels,
    messageHistory,
    currentHistoryIndex,
    isEditing,
    editingMessageId,
    undo,
    redo,
    startEditing,
    cancelEditing,
    saveEdit,
    sendMessage,
    clearMessages,
    setCurrentModel,
    getModelName,
    getModelProvider,
    testProviderConnection,
    currentProvider,
    defaultModel,
    saveDefaultModel,
    addToHistory,
    cancelCurrentRequest,
  }
}) 