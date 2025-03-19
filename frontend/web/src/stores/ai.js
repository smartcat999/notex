import { defineStore } from 'pinia'
import { ref, reactive, computed, watch } from 'vue'
import aiService from '@/services/aiService'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { useUserStore } from './user'

export const useAIStore = defineStore('ai', () => {
  // 状态
  const isLoading = ref(false)
  const isStreaming = ref(false)
  const currentModel = ref('')
  const currentProvider = ref('')
  const defaultModel = ref('')
  const messages = ref([])
  const messageHistory = ref([])
  const currentHistoryIndex = ref(0)
  const editingMessageId = ref(null)
  const availableModels = ref([])
  const userStore = useUserStore()
  const initialized = ref(false)
  const modelsLoaded = ref(false)
  const isInitializing = ref(false)
  const messageSequence = ref(0)
  const conversationThreads = ref({})
  
  // 图像生成相关状态
  const imageModels = ref([])
  const currentImageModel = ref('')
  const isGeneratingImage = ref(false)
  const generatedImages = ref([])

  // 消息类型常量
  const MESSAGE_TYPE = {
    USER: 'user',
    AI: 'assistant'
  }

  // 消息状态常量
  const MESSAGE_STATUS = {
    // 用户消息状态
    USER_SENDING: 'user_sending',
    USER_SENT: 'user_sent',
    USER_CANCELLED: 'user_cancelled',
    USER_EDITING: 'user_editing',

    // AI消息状态
    AI_PREPARING: 'ai_preparing',
    AI_RESPONDING: 'ai_responding',
    AI_RESPONDED: 'ai_responded',
    AI_CANCELLED: 'ai_cancelled'
  }

  // 获取API基础URL
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api'

  // 获取认证头
  const getAuthHeaders = () => {
    const token = userStore.token
    return token ? { Authorization: `Bearer ${token}` } : {}
  }

  // 监听用户登录状态变化
  watch(() => userStore.isAuthenticated, async (isAuthenticated, oldValue) => {
    if (isAuthenticated && !initialized.value && !isInitializing.value) {
      await initialize()
    } else if (!isAuthenticated && initialized.value) {
      // 用户登出，重置状态
      initialized.value = false
      defaultModel.value = ''
      currentModel.value = ''
      currentProvider.value = ''
    }
  })

  // 初始化函数
  async function initialize() {
    if (initialized.value) {
      return
    }
    
    try {
      isLoading.value = true
      
      // 从本地存储加载默认设置
      const localSettings = JSON.parse(localStorage.getItem('aiDefaultSettings') || '{}')
      if (localSettings.defaultImageModel) {
        currentImageModel.value = localSettings.defaultImageModel
      }
      
      // 加载后端默认设置
      const defaultSettings = await aiService.getDefaultSetting()
      if (defaultSettings) {
        currentModel.value = defaultSettings.defaultModel || ''
        // 如果后端有默认图像模型，覆盖本地设置
        if (defaultSettings.defaultImageModel) {
          currentImageModel.value = defaultSettings.defaultImageModel
          // 更新本地存储
          const settings = JSON.parse(localStorage.getItem('aiDefaultSettings') || '{}')
          settings.defaultImageModel = defaultSettings.defaultImageModel
          localStorage.setItem('aiDefaultSettings', JSON.stringify(settings))
        }
      }
      
      // 加载可用模型
      await loadAvailableModels()
      
      // 加载图像模型
      await loadImageModels()
      
      initialized.value = true
    } catch (error) {
      console.error('初始化AI Store失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 加载可用模型
  async function loadAvailableModels() {
    if (modelsLoaded.value && availableModels.value.length > 0) {
      return
    }
    
    try {
      isLoading.value = true
      
      // 获取可用模型列表
      const response = await axios.get(`${apiBaseUrl}/ai/available-models`, {
        headers: getAuthHeaders()
      })
      
      if (response.data && response.data.providers) {
        // 处理提供商和模型数据
        const providers = response.data.providers
        const allModels = []
        
        providers.forEach(provider => {
          if (provider.models && provider.models.length > 0) {
            provider.models.forEach(model => {
              allModels.push({
                id: model.modelId,
                name: model.name,
                provider: provider.providerId,
                providerName: provider.name,
                description: model.description,
                isPaid: model.isPaid,
                type: model.type || 'text' // 默认为文本模型
              })
            })
          }
        })
        
        // 更新状态
        availableModels.value = allModels
        
        // 更新文本模型列表
        updateTextModels()
        
        // 加载图像模型
        try {
          await loadImageModels()
        } catch (err) {
          console.error('加载图像模型失败:', err)
        }
      }
    } catch (error) {
      console.error('加载可用模型失败:', error)
      ElMessage.error('加载可用模型失败')
    } finally {
      isLoading.value = false
    }
  }
  
  // 加载图像生成模型
  async function loadImageModels() {
    try {
      // 获取所有图像类型的模型
      const models = await aiService.getModelsByType('image')
      
      // 更新图像模型列表
      imageModels.value = models.map(model => ({
        id: model.modelId,
        name: model.name,
        provider: model.provider,
        providerName: model.providerName || model.provider,
        description: model.description,
        isPaid: model.isPaid,
        type: 'image'
      }))
      
      // 如果没有当前图像模型但有可用模型，设置第一个为默认
      if (!currentImageModel.value && imageModels.value.length > 0) {
        currentImageModel.value = imageModels.value[0].id
        // 保存到本地存储
        const settings = JSON.parse(localStorage.getItem('aiDefaultSettings') || '{}')
        settings.defaultImageModel = currentImageModel.value
        localStorage.setItem('aiDefaultSettings', JSON.stringify(settings))
      }
      
      return imageModels.value
    } catch (error) {
      console.error('加载图像模型失败:', error)
      throw error
    }
  }
  
  // 生成图像
  async function generateImage(prompt, options = {}) {
    if (!prompt.trim()) {
      ElMessage.warning('请输入图像描述')
      return null
    }
    
    if (!currentImageModel.value) {
      ElMessage.warning('请选择图像生成模型')
      return null
    }
    
    // 找到当前模型信息
    const model = imageModels.value.find(m => m.id === currentImageModel.value)
    if (!model) {
      ElMessage.warning('选择的模型不可用')
      return null
    }

    try {
      // 获取API设置
      const settings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
      const providerSettings = settings[model.provider]
      
      if (!providerSettings?.apiKey) {
        ElMessage.warning(`请先在设置中配置${model.providerName || model.provider}的API密钥`)
        return null
      }

      isGeneratingImage.value = true
      generatedImages.value = []
      
      // 创建请求参数
      const params = {
        prompt,
        model: currentImageModel.value,
        provider: model.provider,
        n: options.n || 1,
        size: options.size || '1024x1024',
        ...options
      }
      
      // 发送请求
      const images = await aiService.generateImage(params)
      generatedImages.value = images
      
      return images
    } catch (error) {
      console.error('生成图像失败:', error)
      ElMessage.error(`生成图像失败: ${error.message || '未知错误'}`)
      return null
    } finally {
      isGeneratingImage.value = false
    }
  }
  
  // 设置当前图像模型
  function setCurrentImageModel(modelId) {
    currentImageModel.value = modelId
    // 也可以保存到本地存储
    localStorage.setItem('currentImageModel', modelId)
  }
  
  // 获取文本模型（过滤出类型为text的模型）
  const textModels = ref([])

  // 更新文本模型列表
  function updateTextModels() {
    textModels.value = availableModels.value.filter(model => model.type === 'text' || !model.type)
  }

  // 加载默认模型
  async function loadDefaultModel() {
    if (initialized.value && defaultModel.value) {
      return
    }
    
    try {
      const response = await axios.get(`${apiBaseUrl}/ai/default-setting`, {
        headers: getAuthHeaders()
      })
      
      if (response.data && response.data.defaultModel) {
        defaultModel.value = response.data.defaultModel
        currentModel.value = response.data.defaultModel
        
        const model = availableModels.value.find(m => m.id === currentModel.value)
        if (model) {
          currentProvider.value = model.provider
        }
      } else {
        if (availableModels.value.length > 0) {
          defaultModel.value = availableModels.value[0].id
          currentModel.value = availableModels.value[0].id
          currentProvider.value = availableModels.value[0].provider
        }
      }
    } catch (error) {
      console.error('加载默认模型失败:', error)
      if (availableModels.value.length > 0) {
        defaultModel.value = availableModels.value[0].id
        currentModel.value = availableModels.value[0].id
        currentProvider.value = availableModels.value[0].provider
      }
    }
  }

  // 保存默认模型
  async function saveDefaultModel(modelId) {
    if (!modelId) return
    
    try {
      await aiService.saveDefaultSetting(modelId)
      defaultModel.value = modelId
      currentModel.value = modelId
      localStorage.setItem('currentModel', modelId)
      return true
    } catch (error) {
      console.error('保存默认模型失败:', error)
      return false
    }
  }

  // 保存默认图像模型
  async function saveDefaultImageModel(modelId) {
    if (!modelId) return
    
    try {
      // 更新本地状态
      currentImageModel.value = modelId
      
      // 保存到后端
      await aiService.saveDefaultSetting({
        defaultImageModel: modelId
      })
      
      // 保存到本地存储
      const settings = JSON.parse(localStorage.getItem('aiDefaultSettings') || '{}')
      settings.defaultImageModel = modelId
      localStorage.setItem('aiDefaultSettings', JSON.stringify(settings))
      
    } catch (error) {
      console.error('保存默认图像模型失败:', error)
      // 即使保存到后端失败，也保持本地状态
      const settings = JSON.parse(localStorage.getItem('aiDefaultSettings') || '{}')
      settings.defaultImageModel = modelId
      localStorage.setItem('aiDefaultSettings', JSON.stringify(settings))
    }
  }

  // 测试提供商连接
  async function testProviderConnection(providerId, apiKey, endpoint, modelId) {
    try {
      isLoading.value = true
      const response = await axios.post(`${apiBaseUrl}/ai/test-connection`, {
        provider: providerId,
        apiKey: apiKey,
        endpoint: endpoint,
        model: modelId || getDefaultModelForProvider(providerId)
      }, {
        headers: getAuthHeaders()
      })
      return response.status === 200
    } catch (error) {
      console.error('测试连接失败:', error)
      return false
    } finally {
      isLoading.value = false
    }
  }

  // 获取提供商的默认模型
  function getDefaultModelForProvider(providerId) {
    const model = availableModels.value.find(m => m.provider === providerId)
    return model ? model.id : null
  }

  // 设置当前模型
  function setCurrentModel(modelId) {
    if (!modelId) return
    
    currentModel.value = modelId
    const model = availableModels.value.find(m => m.id === modelId)
    if (model) {
      currentProvider.value = model.provider
    }
  }

  // 清空消息
  function clearMessages() {
    messages.value = []
    messageHistory.value = []
    currentHistoryIndex.value = 0
    messageSequence.value = 0
    conversationThreads.value = {}
  }

  // 用于取消请求的控制器
  let abortController = null;

  // 发送消息
  async function sendMessage(content, parentId = -1, streamCallback = null) {
    if (!content.trim() || isLoading.value) return

    messageSequence.value++
    const currentSequence = messageSequence.value
    
    // 用于跟踪响应文本长度
    let previousResponseLength = 0
    let previousStreamCallbackResponseLength = 0

    const userMessage = {
      id: Date.now(),
      role: MESSAGE_TYPE.USER,
      content: content.trim(),
      timestamp: new Date(),
      parentId: parentId,
      sequence: currentSequence,
      status: MESSAGE_STATUS.USER_SENDING
    }

    messages.value.push(userMessage)
    
    if (parentId !== -1) {
      if (!conversationThreads.value[parentId]) {
        conversationThreads.value[parentId] = []
      }
      conversationThreads.value[parentId].push(userMessage.id)
    }

    messageSequence.value++
    
    const aiMessage = {
      id: Date.now() + 1,
      role: MESSAGE_TYPE.AI,
      content: '',
      timestamp: new Date(),
      parentId: userMessage.id,
      isStreaming: true,
      sequence: messageSequence.value,
      providerId: currentProvider.value,
      modelId: currentModel.value,
      status: MESSAGE_STATUS.AI_PREPARING
    }

    messages.value.push(aiMessage)
    
    if (!conversationThreads.value[userMessage.id]) {
      conversationThreads.value[userMessage.id] = []
    }
    conversationThreads.value[userMessage.id].push(aiMessage.id)

    saveHistory()

    try {
      userMessage.status = MESSAGE_STATUS.USER_SENT

      isLoading.value = true
      isStreaming.value = true

      const formattedMessages = buildMessageHistory(userMessage.id)
      
      abortController = new AbortController()
      
      if (streamCallback) {
        aiMessage.status = MESSAGE_STATUS.AI_RESPONDING
        
        const response = await axios.post(
          `${apiBaseUrl}/ai/chat`,
          {
            provider: currentProvider.value,
            model: currentModel.value,
            messages: formattedMessages,
            stream: true
          },
          {
            headers: getAuthHeaders(),
            responseType: 'text',
            signal: abortController.signal,
            onDownloadProgress: (progressEvent) => {
              // 获取新增的响应文本，而不是全部文本
              const currentResponseText = progressEvent.event.target.responseText
              const newResponseText = currentResponseText.substring(previousStreamCallbackResponseLength || 0)
              previousStreamCallbackResponseLength = currentResponseText.length
              
              if (newResponseText) {
                try {
                  const lines = newResponseText.split('\n')
                  let newContent = ''
                  
                  for (const line of lines) {
                    if (line.startsWith('data:')) {
                      const data = line.slice(5).trim()
                      if (data === '[DONE]') continue
                      
                      try {
                        const parsed = JSON.parse(data)
                        if (parsed.choices && parsed.choices[0]) {
                          const delta = parsed.choices[0].delta || {}
                          if (delta.content) {
                            newContent += delta.content
                            
                            aiMessage.content += delta.content
                            aiMessage.status = MESSAGE_STATUS.AI_RESPONDING
                            messages.value = [...messages.value]
                            
                            if (streamCallback) {
                              streamCallback({
                                content: delta.content,
                                fullContent: aiMessage.content,
                                messageId: aiMessage.id,
                                isDone: false,
                                status: aiMessage.status
                              })
                            }
                          }
                        }
                      } catch (e) {
                        // 忽略不完整的JSON解析错误
                        if (!line.includes('{') || !line.includes('}')) {
                          continue
                        }
                        console.error('解析SSE数据失败:', e)
                      }
                    }
                  }
                } catch (error) {
                  console.error('处理流式响应失败:', error)
                }
              }
            }
          }
        )
        
        aiMessage.status = MESSAGE_STATUS.AI_RESPONDED
        
        if (streamCallback) {
          streamCallback({
            content: '',
            fullContent: aiMessage.content,
            messageId: aiMessage.id,
            isDone: true,
            status: aiMessage.status
          })
        }
      } else {
        aiMessage.status = MESSAGE_STATUS.AI_RESPONDING
        
        const response = await axios.post(`${apiBaseUrl}/ai/chat`, {
          provider: currentProvider.value,
          model: currentModel.value,
          messages: formattedMessages,
          stream: true
        }, {
          headers: getAuthHeaders(),
          responseType: 'text',
          signal: abortController.signal,
          onDownloadProgress: (progressEvent) => {
            // 获取新增的响应文本，而不是全部文本
            const currentResponseText = progressEvent.event.target.responseText
            const newResponseText = currentResponseText.substring(previousResponseLength || 0)
            previousResponseLength = currentResponseText.length
            
            if (newResponseText) {
              try {
                const lines = newResponseText.split('\n')
                let newContent = ''
                
                for (const line of lines) {
                  if (line.startsWith('data:')) {
                    const data = line.slice(5).trim()
                    if (data === '[DONE]') continue
                    
                    try {
                      const parsed = JSON.parse(data)
                      if (parsed.choices && parsed.choices[0]) {
                        const delta = parsed.choices[0].delta || {}
                        if (delta.content) {
                          newContent += delta.content
                        }
                      }
                    } catch (e) {
                      // 忽略不完整的JSON解析错误
                      if (!line.includes('{') || !line.includes('}')) {
                        continue
                      }
                      console.error('解析SSE数据失败:', e)
                    }
                  }
                }
                
                if (newContent) {
                  aiMessage.content += newContent
                  messages.value = [...messages.value]
                }
              } catch (error) {
                console.error('处理流式响应失败:', error)
              }
            }
          }
        })
        
        aiMessage.status = MESSAGE_STATUS.AI_RESPONDED
      }
    } catch (error) {
      console.error('发送消息失败:', error)
      aiMessage.content = '发送消息失败，请重试。'
      aiMessage.status = MESSAGE_STATUS.AI_CANCELLED
    } finally {
      isLoading.value = false
      isStreaming.value = false
      aiMessage.isStreaming = false
      abortController = null
      saveHistory()
    }

    return aiMessage
  }

  // 构建消息历史，确保保持正确的对话顺序
  function buildMessageHistory(currentMessageId) {
    const currentMessage = messages.value.find(msg => msg.id === currentMessageId)
    if (!currentMessage) {
      return []
    }
    
    let orderedMessages = []
    let processedIds = new Set()
    
    function buildChain(msgId) {
      if (msgId === -1 || processedIds.has(msgId)) return
      
      const msg = messages.value.find(m => m.id === msgId)
      if (!msg) return
      
      if (msg.parentId !== -1) {
        buildChain(msg.parentId)
      }
      
      if (!processedIds.has(msgId)) {
        orderedMessages.push({
          role: msg.role,
          content: msg.content,
          id: msg.id,
          sequence: msg.sequence || 0
        })
        processedIds.add(msgId)
      }
    }
    
    buildChain(currentMessageId)
    
    orderedMessages.sort((a, b) => a.sequence - b.sequence)
    
    const filteredMessages = orderedMessages
      .filter(msg => {
        const originalMsg = messages.value.find(m => m.id === msg.id)
        return !originalMsg || !originalMsg.isStreaming
      })
      .map(msg => ({
        role: msg.role,
        content: msg.content
      }))
    
    return filteredMessages
  }

  // 获取完整对话线程
  function getConversationThread(messageId) {
    const thread = []
    const visited = new Set()
    
    function traverseThread(currentId) {
      if (visited.has(currentId) || !currentId) return
      visited.add(currentId)
      
      const message = messages.value.find(m => m.id === currentId)
      if (message) {
        thread.push(message)
        
        const children = conversationThreads.value[currentId] || []
        children.forEach(childId => traverseThread(childId))
      }
    }
    
    traverseThread(messageId)
    
    return thread.sort((a, b) => a.sequence - b.sequence)
  }

  // 保存历史记录
  function saveHistory() {
    if (currentHistoryIndex.value < messageHistory.value.length - 1) {
      messageHistory.value = messageHistory.value.slice(0, currentHistoryIndex.value + 1)
    }
    
    messageHistory.value.push([...messages.value])
    currentHistoryIndex.value = messageHistory.value.length - 1
  }

  // 撤销操作
  function undo() {
    if (currentHistoryIndex.value > 0) {
      currentHistoryIndex.value--
      messages.value = [...messageHistory.value[currentHistoryIndex.value]]
    }
  }

  // 重做操作
  function redo() {
    if (currentHistoryIndex.value < messageHistory.value.length - 1) {
      currentHistoryIndex.value++
      messages.value = [...messageHistory.value[currentHistoryIndex.value]]
    }
  }

  // 开始编辑消息
  function startEditing(messageId) {
    editingMessageId.value = messageId
    
    const message = messages.value.find(m => m.id === messageId)
    if (message && message.role === MESSAGE_TYPE.USER) {
      message.status = MESSAGE_STATUS.USER_EDITING
      messages.value = [...messages.value]
    }
  }

  // 取消编辑消息
  function cancelEditing() {
    if (editingMessageId.value) {
      const message = messages.value.find(m => m.id === editingMessageId.value)
      if (message && message.role === MESSAGE_TYPE.USER) {
        message.status = MESSAGE_STATUS.USER_SENT
        messages.value = [...messages.value]
      }
    }
    
    editingMessageId.value = null
  }

  // 保存编辑的消息
  async function saveEdit(messageId, content) {
    if (!content.trim()) return

    const messageIndex = messages.value.findIndex(msg => msg.id === messageId)
    if (messageIndex === -1) {
      return
    }

    const message = messages.value[messageIndex]
    
    // 用于跟踪响应文本长度
    let previousResponseLength = 0
    
    if (message.role === MESSAGE_TYPE.USER) {
      message.content = content.trim()
      message.status = MESSAGE_STATUS.USER_SENT
      
      const childResponses = messages.value.filter(msg => msg.parentId === messageId && msg.role === MESSAGE_TYPE.AI)
      
      if (childResponses.length > 0) {
        const streamingMessage = childResponses.find(msg => msg.status === MESSAGE_STATUS.AI_RESPONDING)
        if (streamingMessage) {
          cancelCurrentRequest()
        }
        
        const idsToRemove = new Set()
        
        function collectIdsToRemove(id) {
          if (idsToRemove.has(id)) return
          idsToRemove.add(id)
          
          const children = conversationThreads.value[id] || []
          children.forEach(childId => collectIdsToRemove(childId))
        }
        
        childResponses.forEach(response => collectIdsToRemove(response.id))
        
        messages.value = messages.value.filter(msg => !idsToRemove.has(msg.id))
        
        Object.keys(conversationThreads.value).forEach(key => {
          if (idsToRemove.has(parseInt(key))) {
            delete conversationThreads.value[key]
          } else {
            conversationThreads.value[key] = conversationThreads.value[key].filter(
              id => !idsToRemove.has(id)
            )
          }
        })
      } else {
        const streamingMessageIndex = messages.value.findIndex(msg => msg.status === MESSAGE_STATUS.AI_RESPONDING)
        if (streamingMessageIndex !== -1) {
          cancelCurrentRequest()
        }
      }
      
      saveHistory()
      
      messageSequence.value++
      
      const aiMessage = {
        id: Date.now(),
        role: MESSAGE_TYPE.AI,
        content: '',
        timestamp: new Date(),
        parentId: messageId,
        isStreaming: true,
        sequence: messageSequence.value,
        providerId: currentProvider.value,
        modelId: currentModel.value,
        status: MESSAGE_STATUS.AI_PREPARING
      }
      
      messages.value.push(aiMessage)
      
      if (!conversationThreads.value[messageId]) {
        conversationThreads.value[messageId] = []
      }
      conversationThreads.value[messageId].push(aiMessage.id)
      
      saveHistory()
      
      try {
        isLoading.value = true
        isStreaming.value = true
        
        aiMessage.status = MESSAGE_STATUS.AI_RESPONDING
        
        const formattedMessages = buildMessageHistory(messageId)
        
        abortController = new AbortController()
        
        const headers = getAuthHeaders()
        
        const response = await axios.post(`${apiBaseUrl}/ai/chat`, {
          provider: currentProvider.value,
          model: currentModel.value,
          messages: formattedMessages,
          stream: true
        }, {
          headers,
          responseType: 'text',
          signal: abortController.signal,
          onDownloadProgress: (progressEvent) => {
            // 获取新增的响应文本，而不是全部文本
            const currentResponseText = progressEvent.event.target.responseText
            const newResponseText = currentResponseText.substring(previousResponseLength || 0)
            previousResponseLength = currentResponseText.length
            
            if (newResponseText) {
              try {
                const lines = newResponseText.split('\n')
                let newContent = ''
                
                for (const line of lines) {
                  if (line.startsWith('data:')) {
                    const data = line.slice(5).trim()
                    if (data === '[DONE]') continue
                    
                    try {
                      const parsed = JSON.parse(data)
                      if (parsed.choices && parsed.choices[0]) {
                        const delta = parsed.choices[0].delta || {}
                        if (delta.content) {
                          newContent += delta.content
                        }
                      }
                    } catch (e) {
                      // 忽略不完整的JSON解析错误
                      if (!line.includes('{') || !line.includes('}')) {
                        continue
                      }
                      console.error('解析SSE数据失败:', e)
                    }
                  }
                }
                
                if (newContent) {
                  aiMessage.content += newContent
                  messages.value = [...messages.value]
                }
              } catch (error) {
                console.error('处理流式响应失败:', error)
              }
            }
          }
        })
        
        aiMessage.isStreaming = false
        aiMessage.status = MESSAGE_STATUS.AI_RESPONDED
      } catch (error) {
        console.error('发送消息失败:', error)
        aiMessage.content = '抱歉，我无法生成回复。'
        aiMessage.isStreaming = false
        aiMessage.status = MESSAGE_STATUS.AI_CANCELLED
      } finally {
        isLoading.value = false
        isStreaming.value = false
        abortController = null
        
        saveHistory()
      }
    }
  }

  // 取消当前请求
  function cancelCurrentRequest() {
    if (abortController) {
      abortController.abort()
      abortController = null
    }
    
    if (isStreaming.value) {
      isLoading.value = false
      isStreaming.value = false
      
      const streamingMessageIndex = messages.value.findIndex(msg => msg.status === MESSAGE_STATUS.AI_RESPONDING)
      if (streamingMessageIndex !== -1) {
        const streamingMessage = messages.value[streamingMessageIndex]
        streamingMessage.isStreaming = false
        streamingMessage.status = MESSAGE_STATUS.AI_CANCELLED
        
        if (!editingMessageId.value) {
          streamingMessage.content += '\n\n[已中断生成]'
        }
      }
    }
    
    if (editingMessageId.value) {
      const editingMessage = messages.value.find(msg => msg.id === editingMessageId.value)
      if (editingMessage) {
        const relatedReplies = messages.value.filter(msg => 
          msg.parentId === editingMessageId.value && 
          msg.role === MESSAGE_TYPE.AI
        )
        
        relatedReplies.forEach(reply => {
          reply.isStreaming = false
          reply.status = MESSAGE_STATUS.AI_CANCELLED
          if (reply.content.trim() === '') {
            reply.content = '[已取消]'
          }
        })
      }
    }
    
    saveHistory()
  }

  // 取消用户消息
  function cancelUserMessage(messageId) {
    const message = messages.value.find(msg => msg.id === messageId)
    if (message && message.role === MESSAGE_TYPE.USER) {
      message.status = MESSAGE_STATUS.USER_CANCELLED
      
      const childResponses = messages.value.filter(msg => msg.parentId === messageId && msg.role === MESSAGE_TYPE.AI)
      childResponses.forEach(response => {
        if (response.status === MESSAGE_STATUS.AI_RESPONDING) {
          cancelCurrentRequest()
        }
        response.status = MESSAGE_STATUS.AI_CANCELLED
      })
      
      messages.value = [...messages.value]
      
      saveHistory()
    }
  }

  // 判断消息是否正在编辑
  function isEditing(messageId) {
    return editingMessageId.value === messageId
  }

  // 判断消息是否被取消
  function isCancelled(messageId) {
    const message = messages.value.find(msg => msg.id === messageId)
    if (!message) return false
    
    if (message.role === MESSAGE_TYPE.USER) {
      return message.status === MESSAGE_STATUS.USER_CANCELLED
    } else {
      return message.status === MESSAGE_STATUS.AI_CANCELLED
    }
  }

  // 获取消息状态
  function getMessageStatus(messageId) {
    const message = messages.value.find(msg => msg.id === messageId)
    return message ? message.status : null
  }

  // 初始化
  if (userStore.isAuthenticated && !initialized.value && !isInitializing.value) {
    initialize()
  } else if (!modelsLoaded.value) {
    loadAvailableModels()
    modelsLoaded.value = true
  }

  return {
    // 状态
    isLoading,
    isStreaming,
    currentModel,
    currentProvider,
    defaultModel,
    messages,
    messageHistory,
    currentHistoryIndex,
    editingMessageId,
    availableModels,
    initialized,
    modelsLoaded,
    conversationThreads,
    
    // 图像生成相关
    imageModels,
    currentImageModel,
    isGeneratingImage,
    generatedImages,
    textModels,
    
    // 类型常量
    MESSAGE_TYPE,
    MESSAGE_STATUS,

    // 方法
    initialize,
    loadAvailableModels,
    loadDefaultModel,
    setCurrentModel,
    clearMessages,
    sendMessage,
    undo,
    redo,
    startEditing,
    cancelEditing,
    saveEdit,
    cancelCurrentRequest,
    saveDefaultModel,
    saveDefaultImageModel,
    testProviderConnection,
    getConversationThread,
    cancelUserMessage,
    isEditing,
    isCancelled,
    getMessageStatus,
    
    // 图像生成相关方法
    loadImageModels,
    generateImage,
    setCurrentImageModel,
    updateTextModels
  }
}) 