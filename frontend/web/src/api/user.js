import request from '@/utils/request'

export function getUserPosts(params) {
  return request({
    url: `/users/${params.userId}/posts`,
    method: 'get',
    params: {
      page: params.page,
      per_page: params.pageSize
    }
  })
}

export function getUserComments(params) {
  return request({
    url: `/users/${params.userId}/comments`,
    method: 'get',
    params: {
      page: params.page,
      per_page: params.pageSize
    }
  })
}

// 获取公开用户信息
export function getUserProfile(userId) {
  return request({
    url: `/public/users/${userId}/profile`,
    method: 'get'
  })
}

// 获取用户公开信息和文章列表
export function getUserProfileWithPosts(userId, params) {
  return request({
    url: `/public/users/${userId}/profile`,
    method: 'get',
    params: {
      page: params?.page || 1,
      per_page: params?.per_page || 10
    }
  })
}

// 获取用户主页信息（包含用户资料和文章列表）
export function getUserHome(userId, params) {
  return request({
    url: `/public/users/${userId}/home`,
    method: 'get',
    params: {
      page: params?.page || 1,
      per_page: params?.per_page || 10
    }
  })
} 