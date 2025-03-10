import request from '@/utils/request'

// 获取系统设置
export function getSettings() {
  return request({
    url: '/api/settings',
    method: 'get'
  })
}

// 更新系统设置
export function updateSettings(type, data) {
  return request({
    url: `/api/settings/${type}`,
    method: 'put',
    data
  })
}

// 测试邮件发送
export function testEmail() {
  return request({
    url: '/api/settings/test-email',
    method: 'post'
  })
}

// 上传文件
export function uploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/api/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
} 