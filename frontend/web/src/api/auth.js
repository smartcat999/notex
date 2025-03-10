import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data,
  })
}

export function register(data) {
  return request({
    url: '/auth/register',
    method: 'post',
    data,
  })
}

export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post',
  })
}

export function getProfile() {
  return request({
    url: '/auth/profile',
    method: 'get',
  })
}

export function updateProfile(data) {
  return request({
    url: '/auth/profile',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/auth/change-password',
    method: 'post',
    data,
  })
}

export function sendVerificationCode(email) {
  return request({
    url: '/auth/email/send-verification',
    method: 'post',
    data: { email },
  })
}

export function verifyEmail(data) {
  return request({
    url: '/auth/email/verify',
    method: 'post',
    data,
  })
}

export function updateEmail(data) {
  return request({
    url: '/auth/email/update',
    method: 'post',
    data,
  })
}

export function sendPasswordReset(email) {
  return request({
    url: '/auth/password-reset/send',
    method: 'post',
    data: { email },
  })
}

export function resetPassword(data) {
  return request({
    url: '/auth/password-reset/verify',
    method: 'post',
    data,
  })
} 