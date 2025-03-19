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

// 提供商API端点导出
export const getProviderEndpoint = (providerId) => {
  const endpoints = {
    'openai': 'https://api.openai.com/v1/chat/completions',
    'anthropic': 'https://api.anthropic.com/v1/messages',
    'google': 'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent',
    'deepseek': 'https://api.deepseek.com/v1/chat/completions',
    'custom': '' // 自定义端点需要从设置中获取
  }
  return endpoints[providerId] || ''
}

// 获取图像生成API端点
export const getImageGenerationEndpoint = (providerId) => {
  const endpoints = {
    'openai': 'https://api.openai.com/v1/images/generations',
    'stabilityai': 'https://api.stability.ai/v1/generation/stable-diffusion-xl-1024-v1-0/text-to-image',
    'custom': '' // 自定义端点需要从设置中获取
  }
  return endpoints[providerId] || ''
}

// 获取提供商API密钥
const getProviderApiKey = (providerId) => {
  const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
  const providerSettings = allSettings[providerId] || {}
  return providerSettings.apiKey || ''
}

/**
 * 获取指定类型的模型列表
 * @param {string} type - 模型类型，如 'text' 或 'image'
 * @returns {Promise} - 返回模型列表
 */
export const getModelsByType = async (type) => {
  try {
    const response = await axios.get(`/api/ai/models/type/${type}`, {
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders()
      }
    })
    return response.data.models || []
  } catch (error) {
    console.error(`获取${type}类型模型失败:`, error)
    throw error
  }
}

/**
 * 生成图像
 * @param {Object} params - 图像生成参数
 * @param {string} params.prompt - 图像描述提示
 * @param {string} params.model - 模型ID
 * @param {string} params.provider - 提供商ID
 * @param {number} params.n - 生成图像数量
 * @param {string} params.size - 图像尺寸
 * @returns {Promise} - 返回生成的图像URLs
 */
export const generateImage = async (params) => {
  try {
    // 验证提示文本
    if (!params.prompt || !params.prompt.trim()) {
      throw new Error('图像描述不能为空')
    }

    // 获取API密钥
    const allSettings = JSON.parse(localStorage.getItem('aiProviderSettings') || '{}')
    const providerSettings = allSettings[params.provider] || {}
    const apiKey = providerSettings.apiKey || params.apiKey

    if (!apiKey) {
      throw new Error('未配置API密钥')
    }

    // 准备请求数据
    const requestData = {
      provider: params.provider,
      model: params.model,
      prompt: params.prompt,
      n: params.n || 1,
      size: params.size || '1024x1024',
      apiKey: apiKey
    }

    // 如果有自定义端点，添加到请求数据
    if (providerSettings.endpoint && params.provider === 'custom') {
      requestData.endpoint = providerSettings.endpoint
    }

    // 发送请求
    const response = await axios.post('/api/ai/generate-image', requestData, {
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders()
      }
    })

    // 处理不同提供商的响应格式
    if (params.provider === 'openai') {
      return response.data.data.map(item => item.url)
    } else if (params.provider === 'stabilityai') {
      return response.data.artifacts.map(item => `data:image/png;base64,${item.base64}`)
    } else {
      // 通用处理
      return response.data.images || response.data.urls || []
    }
  } catch (error) {
    console.error('生成图像失败:', error)
    throw error
  }
}

// 获取认证头
const getAuthHeaders = () => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
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

// 保存默认模型设置
async function saveDefaultSetting(defaultModel, defaultImageModel) {
  // 构建请求数据
  const data = {};
  
  // 根据传入的参数设置请求数据
  if (defaultModel) {
    data.defaultModel = defaultModel;
  }
  
  if (defaultImageModel) {
    data.defaultImageModel = defaultImageModel;
  }
  
  // 发送请求
  const response = await makeAuthenticatedRequest('post', `/ai/default-setting`, data);
  return response;
}

/**
 * 获取所有可用的AI模型
 * @returns {Promise<Array>} - 返回可用的AI模型列表
 */
async function getAvailableModels() {
  try {
    const response = await axios.get('/api/ai/available-models', {
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders()
      }
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
      
      return allModels
    }
    
    return []
  } catch (error) {
    console.error('获取可用模型失败:', error)
    throw error
  }
}

/**
 * 获取默认模型设置
 * @returns {Promise<Object>} - 返回默认模型设置
 */
async function getDefaultSetting() {
  try {
    const response = await axios.get('/api/ai/default-setting', {
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders()
      }
    })
    
    return response.data;
  } catch (error) {
    console.error('获取默认模型设置失败:', error)
    throw error
  }
}

/**
 * 发送认证请求的辅助函数
 * @param {string} method - 请求方法
 * @param {string} url - 请求URL
 * @param {Object} data - 请求数据
 * @returns {Promise<Object>} - 返回响应数据
 */
async function makeAuthenticatedRequest(method, url, data = null) {
  const apiBaseUrl = '/api'; // 可以从环境变量获取
  const fullUrl = `${apiBaseUrl}${url}`;
  
  try {
    const config = {
      method,
      url: fullUrl,
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders()
      }
    };
    
    if (data && (method.toLowerCase() === 'post' || method.toLowerCase() === 'put')) {
      config.data = data;
    }
    
    const response = await axios(config);
    return response.data;
  } catch (error) {
    console.error(`${method.toUpperCase()} ${url} 请求失败:`, error);
    throw error;
  }
}

export default {
  sendChat,
  testProviderConnection,
  getModelsByType,
  generateImage,
  getProviderEndpoint,
  getImageGenerationEndpoint,
  saveDefaultSetting,
  getAvailableModels,
  getDefaultSetting
} 