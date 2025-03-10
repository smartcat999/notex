import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    // console.log('response', res)
    
    // 对于草稿详情接口，直接返回数据
    if (response.config.url.match(/\/api\/drafts\/\d+$/) && response.config.method === 'get') {
      return res
    }
    
    // 如果响应成功，直接返回数据
    return res
  },
  error => {
    console.error('Request error:', error)
    ElMessage({
      message: error.response?.data?.error || '请求失败',
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default request 