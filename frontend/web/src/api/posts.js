import request from '@/utils/request'

// 获取文章列表
export function getPosts(params) {
  return request({
    url: '/posts',
    method: 'get',
    params: {
      page: params.page || 1,
      page_size: params.page_size || 10,
      search: params.search,
      category_id: params.category_id,
      tag_id: params.tag_id,
      sort: params.sort || 'newest',
      status: params.status
    }
  })
}

// 获取文章详情
export function getPost(id) {
  return request({
    url: `/public/posts/${id}`,
    method: 'get'
  })
}

// 创建文章
export function createPost(data) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}

// 更新文章
export function updatePost(id, data) {
  return request({
    url: `/posts/${id}`,
    method: 'put',
    data
  })
}

// 删除文章
export function deletePost(id) {
  return request({
    url: `/posts/${id}`,
    method: 'delete'
  })
}

// 获取最新文章
export function getRecentPosts(params) {
  return request({
    url: '/posts/recent',
    method: 'get',
    params
  })
}

// 获取分类列表
export function getCategories(params) {
  return request({
    url: '/public/categories',
    method: 'get',
    params
  })
}

// 获取热门分类
export function getTopCategories(params) {
  return request({
    url: '/public/categories/top',
    method: 'get',
    params
  })
}

// 获取标签列表
export function getTags(params) {
  return request({
    url: '/public/tags',
    method: 'get',
    params
  })
}

// 获取热门标签
export function getTopTags(params) {
  return request({
    url: '/public/tags/top',
    method: 'get',
    params
  })
}

// 获取分类下的文章列表
export function getCategoryPosts(categoryId, params) {
  return request({
    url: `/posts`,
    method: 'get',
    params: {
      ...params,
      category: categoryId
    }
  })
}

// 获取标签下的文章列表
export const getTagPosts = (tagId, params) => {
  return request({
    url: '/posts',
    method: 'get',
    params: {
      ...params,
      tag_id: tagId
    }
  })
}

// 创建分类
export function createCategory(data) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

// 创建标签
export function createTag(data) {
  return request({
    url: '/tags',
    method: 'post',
    data
  })
}

// 获取公开文章列表（匿名访问）
export function getPublicPosts(params) {
  return request({
    url: '/public/posts',
    method: 'get',
    params: {
      page: params.page || 1,
      page_size: params.per_page || 10,
      user_id: params.user_id,
      search: params.search,
      category_id: params.category_id,
      tag_id: params.tag_id,
      sort: params.sort || 'newest',
      status: 'published'
    }
  })
}

// 获取文章评论列表
export function getComments(postId, params) {
  return request({
    url: `/public/posts/${postId}/comments`,
    method: 'get',
    params
  })
}

// 创建评论
export function createComment(postId, data) {
  return request({
    url: `/posts/${postId}/comments`,
    method: 'post',
    data
  })
}

// 删除评论
export function deleteComment(postId, commentId) {
  return request({
    url: `/posts/${postId}/comments/${commentId}`,
    method: 'delete'
  })
}

// 获取文章归档列表
export function getArchives() {
  return request({
    url: '/public/posts/archives',
    method: 'get'
  })
}

// 获取指定归档日期的文章列表
export function getPostsByArchive(yearMonth) {
  return request({
    url: `/public/posts/archives/${yearMonth}`,
    method: 'get'
  })
}

// 获取用户的文章列表
export function getUserPosts(params) {
  return request({
    url: '/posts',
    method: 'get',
    params: {
      page: params.page,
      per_page: params.per_page,
      user: 'current'  // 使用 current 表示当前用户
    }
  })
}

// 获取用户评论列表
export function getUserComments(params) {
  return request({
    url: '/users/comments',
    method: 'get',
    params
  })
}

// 获取用户的公开文章列表
export function getUserPublicPosts(userId, params) {
  return request({
    url: `/users/${userId}/posts`,
    method: 'get',
    params: {
      page: params.page || 1,
      per_page: params.per_page || 10
    }
  })
}

// 获取评论回复列表
export function getReplies(commentId) {
  return request({
    url: `/public/comments/${commentId}/replies`,
    method: 'get'
  })
}