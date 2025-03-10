import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/admin/auth/login',
    method: 'post',
    data,
  })
}

export function logout() {
  return request({
    url: '/admin/auth/logout',
    method: 'post',
  })
}

export function getProfile() {
  return request({
    url: '/admin/auth/profile',
    method: 'get',
  })
}

export function changePassword(data) {
  return request({
    url: '/admin/auth/change-password',
    method: 'post',
    data,
  })
}

export function updateProfile(data) {
  return request({
    url: '/admin/auth/profile',
    method: 'put',
    data,
  })
} 