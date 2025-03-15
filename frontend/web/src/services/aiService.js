import axios from 'axios'

// 创建一个专用的axios实例
const aiApi = axios.create({
  timeout: 60000, // 60秒超时，因为AI请求可能需要更长时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 添加请求拦截器，用于添加认证信息
aiApi.interceptors.request.use(
  config => {
    // 从localStorage获取AI提供商设置
    const providerId = config.providerId
    if (providerId) {
      const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
      const providerSettings = allSettings[providerId]
      
      if (providerSettings?.apiKey) {
        // 根据不同的提供商设置不同的认证头
        switch (providerId) {
          case 'openai':
            config.headers['Authorization'] = `Bearer ${providerSettings.apiKey}`
            break
          case 'anthropic':
            config.headers['x-api-key'] = providerSettings.apiKey
            break
          case 'google':
            config.headers['Authorization'] = `Bearer ${providerSettings.apiKey}`
            break
          case 'deepseek':
            config.headers['Authorization'] = `Bearer ${providerSettings.apiKey}`
            break
          case 'custom':
            config.headers['Authorization'] = `Bearer ${providerSettings.apiKey}`
            break
          default:
            break
        }
      }
      
      // 如果有自定义端点，使用自定义端点
      if (providerSettings?.endpoint && providerId === 'custom') {
        config.baseURL = providerSettings.endpoint
      }
      
      // 删除自定义属性，避免发送到服务器
      delete config.providerId
    }
    
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

/**
 * 发送聊天消息
 * @param {Array} messages - 消息历史
 * @param {string} modelId - 模型ID
 * @param {string} providerId - 提供商ID
 * @param {AbortSignal} signal - 取消信号
 * @param {Function} onProgress - 流式响应的进度回调
 * @returns {Promise} - 返回AI响应
 */
export const sendChat = async (messages, modelId, providerId, signal, onProgress) => {
  try {
    // 验证消息数组
    if (!Array.isArray(messages) || messages.length === 0) {
      throw new Error('消息数组不能为空')
    }

    // 验证每条消息
    for (const msg of messages) {
      if (!msg.content || !msg.content.trim()) {
        throw new Error('消息内容不能为空')
      }
    }

    // 获取API密钥
    const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
    const providerSettings = allSettings[providerId] || {}
    const apiKey = providerSettings.apiKey

    if (!apiKey) {
      throw new Error('未配置API密钥')
    }

    // 获取模型参数
    const modelParams = providerSettings.modelParams?.[modelId] || {}

    // 获取API端点
    const endpoint = getProviderEndpoint(providerId)
    if (!endpoint) {
      throw new Error('未配置API端点')
    }

    // 根据不同的提供商设置不同的请求头
    const headers = {
      'Content-Type': 'application/json'
    }

    switch (providerId) {
      case 'openai':
        headers['Authorization'] = `Bearer ${apiKey}`
        break
      case 'anthropic':
        headers['x-api-key'] = apiKey
        break
      case 'google':
        headers['Authorization'] = `Bearer ${apiKey}`
        break
      case 'deepseek':
        headers['Authorization'] = `Bearer ${apiKey}`
        break
      case 'custom':
        headers['Authorization'] = `Bearer ${apiKey}`
        break
      default:
        throw new Error(`不支持的AI提供商: ${providerId}`)
    }

    // 准备请求体
    const requestBody = {
      messages: formatMessages(messages, providerId),
      model: modelId,
      stream: true,
      ...modelParams
    }

    const response = await fetch(endpoint, {
      method: 'POST',
      headers,
      body: JSON.stringify(requestBody),
      signal
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => null)
      throw new Error(errorData?.error?.message || `HTTP error! status: ${response.status}`)
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''
    let fullContent = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      
      // 保留最后一个不完整的行
      buffer = lines.pop() || ''

      for (const line of lines) {
        const trimmedLine = line.trim()
        if (!trimmedLine || trimmedLine === 'data: [DONE]') continue

        if (trimmedLine.startsWith('data: ')) {
          try {
            const data = JSON.parse(trimmedLine.slice(6))
            if (data.choices?.[0]?.delta?.content) {
              const content = data.choices[0].delta.content
              fullContent += content
              if (onProgress) {
                onProgress(content)
              }
            }
          } catch (e) {
            console.error('解析流式响应数据出错:', e)
          }
        }
      }
    }

    return fullContent
  } catch (error) {
    if (error.name === 'AbortError') {
      throw error
    }
    throw new Error(`发送消息失败: ${error.message}`)
  }
}

// 格式化消息以适应不同提供商的格式
const formatMessages = (messages, providerId) => {
  switch (providerId) {
    case 'anthropic':
      return messages.map(msg => ({
        role: msg.role === 'assistant' ? 'assistant' : 'user',
        content: msg.content
      }))
    case 'google':
      return messages.map(msg => ({
        role: msg.role,
        parts: [{ text: msg.content }]
      }))
    default:
      return messages.map(msg => ({
        role: msg.role,
        content: msg.content
      }))
  }
}

// 获取提供商的API端点
const getProviderEndpoint = (providerId) => {
  switch (providerId) {
    case 'openai':
      return 'https://api.openai.com/v1/chat/completions'
    case 'anthropic':
      return 'https://api.anthropic.com/v1/messages'
    case 'google':
      return 'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent'
    case 'deepseek':
      return 'https://api.deepseek.com/v1/chat/completions'
    case 'custom':
      const settings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
      return settings.custom?.endpoint || ''
    default:
      throw new Error(`不支持的AI提供商: ${providerId}`)
  }
}

// 获取提供商的API密钥
const getProviderApiKey = (providerId) => {
  const settings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
  return settings[providerId]?.apiKey
}

/**
 * 测试AI提供商连接
 * @param {string} providerId - 提供商ID
 * @returns {Promise<boolean>} - 连接是否成功
 */
export const testProviderConnection = async (providerId) => {
  try {
    const testMessage = [{ role: 'user', content: 'Hello, this is a test message.' }]
    
    // 获取提供商的第一个模型
    const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
    const providerSettings = allSettings[providerId] || {}
    
    // 找到第一个启用的模型
    let modelId = null
    if (providerSettings.enabledModels) {
      for (const [id, enabled] of Object.entries(providerSettings.enabledModels)) {
        if (enabled) {
          modelId = id
          break
        }
      }
    }
    
    if (!modelId) {
      throw new Error('未找到启用的模型')
    }
    
    // 发送测试消息
    await sendChat(testMessage, modelId, providerId)
    return true
  } catch (error) {
    console.error(`测试连接失败: ${error.message}`)
    return false
  }
}

export default {
  sendChat,
  testProviderConnection
} 