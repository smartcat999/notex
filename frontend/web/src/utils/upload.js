import axios from 'axios'
import CryptoJS from 'crypto-js'
import { ElMessage } from 'element-plus'

// 上传配置缓存
let uploadConfig = null

// 获取上传配置
const getUploadConfig = async (token) => {
  if (uploadConfig) return uploadConfig

  try {
    const response = await axios.get('/api/upload/config', {
      headers: {
        Authorization: 'Bearer ' + token
      }
    })
    if (response.data) {
      uploadConfig = response.data
      return uploadConfig
    }
  } catch (error) {
    console.error('Failed to get upload config:', error)
    throw new Error('获取上传配置失败')
  }
}

// 生成OSS Policy
const generateOSSPolicy = (config) => {
  const date = new Date()
  date.setHours(date.getHours() + 1)
  
  const policyText = {
    expiration: date.toISOString(),
    conditions: [
      ['content-length-range', 0, config.maxSize],
      ['starts-with', '$key', ''],
      { bucket: config.bucket }
    ]
  }
  
  return btoa(JSON.stringify(policyText))
}

// 计算OSS签名
const calculateOSSSignature = (policy, credentials) => {
  return CryptoJS.enc.Base64.stringify(
    CryptoJS.HmacSHA1(policy, credentials.accessKeySecret)
  )
}

// 构造OSS上传数据
const getOSSUploadData = (file, credentials, config) => {
  const policy = generateOSSPolicy(config)
  const signature = calculateOSSSignature(policy, credentials)
  
  return {
    'key': credentials.fileKey,
    'policy': policy,
    'OSSAccessKeyId': credentials.accessKeyId,
    'success_action_status': '200',
    'signature': signature,
    'x-oss-security-token': credentials.securityToken,
    'Content-Type': file.type,
  }
}

// 获取上传凭证
const getUploadCredentials = async (file, token) => {
  try {
    const response = await axios.get('/api/upload/credentials', {
      headers: {
        'Authorization': 'Bearer ' + token
      },
      params: {
        filename: file.name,
        contentType: file.type
      }
    })
    
    if (!response.data) {
      throw new Error('获取上传凭证失败')
    }
    
    return response.data
  } catch (error) {
    console.error('Failed to get upload credentials:', error)
    throw new Error('获取上传凭证失败')
  }
}

// 验证文件
const validateFile = (file, config) => {
  // 验证文件类型
  if (!config.allowedTypes.includes(file.type)) {
    throw new Error('文件类型不支持')
  }
  
  // 验证文件大小
  if (file.size > config.maxSize) {
    throw new Error(`文件大小不能超过 ${config.maxSize / 1024 / 1024}MB`)
  }

  return true
}

// 处理上传结果
const processUploadResult = (response, credentials, config) => {
  let fileUrl = ''
  
  if (config.directUpload) {
    switch (config.type) {
      case 'oss':
        if (credentials && credentials.fileKey) {
          const prefix = config.urlPrefix.endsWith('/') 
            ? config.urlPrefix 
            : config.urlPrefix + '/'
          fileUrl = prefix + credentials.fileKey
        }
        break
      case 'cos':
        fileUrl = response.url || response.Location
        break
      case 'minio':
        fileUrl = response.url || (config.urlPrefix + response.key)
        break
      default:
        fileUrl = response.url
    }
  } else {
    fileUrl = response.url
  }
  
  if (!fileUrl) {
    throw new Error('上传失败：无效的响应数据')
  }
  
  return fileUrl
}

// 上传文件
export const uploadFile = async (file, token, options = {}) => {
  try {
    // 获取上传配置
    const config = await getUploadConfig(token)
    
    // 合并配置
    const finalConfig = {
      ...config,
      maxSize: options.maxSize || config.maxSize,
      allowedTypes: options.allowedTypes || config.allowedTypes
    }
    
    // 验证文件
    validateFile(file, finalConfig)
    
    let credentials = null
    let formData = new FormData()
    
    // 如果是直传，获取凭证
    if (finalConfig.directUpload) {
      credentials = await getUploadCredentials(file, token)
      
      if (finalConfig.type === 'oss') {
        const ossData = getOSSUploadData(file, credentials, finalConfig)
        Object.entries(ossData).forEach(([key, value]) => {
          formData.append(key, value)
        })
      }
    }
    
    // 添加文件到 FormData
    formData.append('file', file)
    
    // 发送上传请求
    const response = await axios.post(
      finalConfig.directUpload ? finalConfig.uploadUrl : '/api/upload/file',
      formData,
      {
        headers: {
          'Authorization': 'Bearer ' + token,
          ...finalConfig.headers
        }
      }
    )
    
    // 处理上传结果
    const fileUrl = processUploadResult(response.data, credentials, finalConfig)
    return fileUrl
    
  } catch (error) {
    console.error('Upload error:', error)
    throw error
  }
} 