import request from '@/utils/request'

// 获取草稿列表
export function getDrafts(params) {
  return request({
    url: '/drafts',
    method: 'get',
    params
  })
}

// 获取单个草稿详情
export function getDraft(id) {
  return request({
    url: `/drafts/${id}`,
    method: 'get'
  })
}

// 创建草稿
export function createDraft(data) {
  return request({
    url: '/drafts',
    method: 'post',
    data
  })
}

// 更新草稿
export function updateDraft(id, data) {
  return request({
    url: `/drafts/${id}`,
    method: 'put',
    data
  })
}

// 删除草稿
export function deleteDraft(id) {
  return request({
    url: `/drafts/${id}`,
    method: 'delete'
  })
}

// 发布草稿
export function publishDraft(id) {
  return request({
    url: `/drafts/${id}/publish`,
    method: 'post'
  })
} 