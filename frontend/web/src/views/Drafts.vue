<template>
  <div class="drafts-container">
    <div class="header">
      <h1>草稿箱</h1>
      <div class="search-box">
        <el-input
          v-model="searchQuery"
          placeholder="搜索草稿..."
          @input="handleSearch"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <div class="drafts-list" v-loading="loading">
      <template v-if="drafts.length > 0">
        <div v-for="draft in drafts" :key="draft.id" class="draft-card">
          <div class="draft-content">
            <h2 class="draft-title">{{ draft.title || '无标题草稿' }}</h2>
            <div class="draft-meta">
              <span class="draft-date">最后编辑: {{ formatDate(draft.updated_at) }}</span>
              <span class="draft-category" v-if="draft.category">
                {{ draft.category }}
              </span>
            </div>
            <p class="draft-summary">{{ draft.summary || '暂无摘要' }}</p>
          </div>
          <div class="draft-actions">
            <el-button type="primary" @click="handleEdit(draft)">编辑</el-button>
            <el-button type="success" @click="handlePublish(draft)">发布</el-button>
            <el-button type="danger" @click="handleDelete(draft)">删除</el-button>
          </div>
        </div>
      </template>
      <el-empty v-else description="暂无草稿" />
    </div>

    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getDrafts, deleteDraft, publishDraft } from '@/api/drafts'

const router = useRouter()
const loading = ref(false)
const drafts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const searchTimeout = ref(null)

const fetchDrafts = async () => {
  loading.value = true
  try {
    const response = await getDrafts({
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value
    })
    drafts.value = response.items || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to fetch drafts:', error)
    ElMessage.error('获取草稿列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }
  searchTimeout.value = setTimeout(() => {
    currentPage.value = 1
    fetchDrafts()
  }, 300)
}

const handleSizeChange = (size) => {
  pageSize.value = size
  fetchDrafts()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchDrafts()
}

const handleEdit = (draft) => {
  router.push({
    name: 'NewPost',
    query: { draft_id: draft.id }
  })
}

const handlePublish = async (draft) => {
  try {
    await ElMessageBox.confirm('确定要发布这篇草稿吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await publishDraft(draft.id)
    ElMessage.success('发布成功')
    fetchDrafts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to publish draft:', error)
      ElMessage.error('发布失败')
    }
  }
}

const handleDelete = async (draft) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇草稿吗？此操作不可恢复！', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteDraft(draft.id)
    ElMessage.success('删除成功')
    fetchDrafts()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete draft:', error)
      ElMessage.error('删除失败')
    }
  }
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchDrafts()
})
</script>

<style scoped>
.drafts-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
  color: var(--el-text-color-primary);
}

.search-box {
  width: 300px;
}

.drafts-list {
  min-height: 200px;
}

.draft-card {
  background: var(--el-bg-color);
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.3s ease;
}

.draft-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
}

.draft-content {
  flex: 1;
  margin-right: 20px;
}

.draft-title {
  margin: 0 0 8px;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

.draft-meta {
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.draft-date {
  margin-right: 16px;
}

.draft-category {
  background: var(--el-color-primary-light-9);
  padding: 2px 8px;
  border-radius: 4px;
}

.draft-summary {
  margin: 0;
  font-size: 14px;
  color: var(--el-text-color-regular);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.draft-actions {
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .search-box {
    width: 100%;
  }

  .draft-card {
    flex-direction: column;
  }

  .draft-content {
    margin-right: 0;
    margin-bottom: 16px;
  }

  .draft-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style> 