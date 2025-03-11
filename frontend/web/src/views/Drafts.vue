<template>
  <div class="drafts-container">
    <div class="page-header">
      <h1>草稿箱</h1>
      <div class="header-actions">
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
          <div class="draft-header">
            <h2 class="draft-title">
              <a :href="draft.url" target="_blank">{{ draft.title || '无标题草稿' }}</a>
            </h2>
            <div class="draft-actions">
              <el-button type="primary" @click="handleEdit(draft)">编辑</el-button>
              <el-button type="success" @click="handlePublish(draft)">发布</el-button>
              <el-button type="danger" @click="handleDelete(draft)">删除</el-button>
            </div>
          </div>
          <div class="draft-meta">
            <span class="draft-date">最后编辑: {{ formatDate(draft.updated_at) }}</span>
            <span class="draft-category" v-if="draft.category">
              {{ draft.category }}
            </span>
          </div>
          <p class="draft-summary">{{ draft.summary || '暂无摘要' }}</p>
          <div class="draft-footer">
            <div class="draft-tags">
              <el-tag v-for="tag in draft.tags" :key="tag" :type="tag.type" :effect="tag.effect">
                {{ tag.name }}
              </el-tag>
            </div>
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
  max-width: 900px;
  margin: 32px auto;
  padding: 0 20px;

  .page-header {
    margin-bottom: 32px;
    text-align: center;

    h1 {
      margin: 0 0 16px;
      font-size: 2.2em;
      font-weight: 700;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: 0.02em;
      position: relative;
      display: inline-block;

      &::after {
        content: '';
        position: absolute;
        bottom: 4px;
        left: -8px;
        right: -8px;
        height: 12px;
        background-color: rgba(43, 88, 118, 0.15);
        border-radius: 6px;
        z-index: -1;
        transform: skew(-12deg);
      }
    }

    .header-actions {
      display: flex;
      justify-content: center;
      gap: 16px;
      margin-top: 24px;

      .action-button {
        padding: 12px 24px;
        border-radius: 8px;
        transition: all 0.3s ease;
        background: linear-gradient(135deg, #2B5876, #4E4376);
        color: white;
        border: none;
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 1em;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
        }

        .el-icon {
          font-size: 1.2em;
        }
      }
    }
  }

  .drafts-list {
    .draft-card {
      background: rgba(255, 255, 255, 0.98);
      border-radius: 12px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
      padding: 24px;
      margin-bottom: 20px;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
        background: rgba(255, 255, 255, 1);
      }

      .draft-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 16px;

        .draft-title {
          margin: 0;
          font-size: 1.3em;
          font-weight: 600;
          color: #2c3e50;

          a {
            color: inherit;
            text-decoration: none;
            transition: all 0.3s ease;

            &:hover {
              color: #2B5876;
              text-shadow: 0 0 1px rgba(43, 88, 118, 0.15);
            }
          }
        }

        .draft-actions {
          display: flex;
          gap: 12px;

          .el-button {
            padding: 8px 16px;
            border: none;
            transition: all 0.3s ease;
            font-weight: 500;
            font-size: 0.95em;

            &.el-button--primary {
              background: linear-gradient(135deg, #2B5876, #4E4376);
              color: white;

              &:hover {
                transform: translateY(-1px);
                box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
              }
            }

            &.el-button--success {
              background: linear-gradient(135deg, #2ecc71, #27ae60);
              color: white;

              &:hover {
                transform: translateY(-1px);
                box-shadow: 0 4px 12px rgba(46, 204, 113, 0.2);
              }
            }

            &.el-button--danger {
              background: linear-gradient(135deg, #e74c3c, #c0392b);
              color: white;

              &:hover {
                transform: translateY(-1px);
                box-shadow: 0 4px 12px rgba(231, 76, 60, 0.2);
              }
            }
          }
        }
      }

      .draft-meta {
        display: flex;
        align-items: center;
        gap: 16px;
        color: #718096;
        font-size: 0.9em;
        margin-bottom: 12px;

        .meta-item {
          display: flex;
          align-items: center;
          gap: 6px;

          .el-icon {
            font-size: 1.1em;
            opacity: 0.85;
          }
        }
      }

      .draft-summary {
        color: #4a5568;
        font-size: 0.95em;
        line-height: 1.6;
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .draft-footer {
        margin-top: 16px;
        padding-top: 16px;
        border-top: 1px solid rgba(0, 0, 0, 0.04);

        .draft-tags {
          display: flex;
          gap: 8px;
          flex-wrap: wrap;

          .el-tag {
            background: rgba(43, 88, 118, 0.06);
            border: none;
            color: #2B5876;
            transition: all 0.3s ease;

            &:hover {
              background: rgba(43, 88, 118, 0.1);
              transform: translateY(-1px);
            }
          }
        }
      }
    }
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    background: rgba(43, 88, 118, 0.03);
    border-radius: 16px;
    margin-top: 32px;

    .empty-icon {
      font-size: 3em;
      color: #2B5876;
      margin-bottom: 16px;
      opacity: 0.5;
    }

    .empty-text {
      color: #4a5568;
      font-size: 1.1em;
      margin-bottom: 24px;
    }

    .create-button {
      background: linear-gradient(135deg, #2B5876, #4E4376);
      color: white;
      border: none;
      padding: 12px 28px;
      border-radius: 8px;
      font-size: 1em;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
      }

      .el-icon {
        margin-right: 8px;
      }
    }
  }

  .pagination {
    margin-top: 24px;
    display: flex;
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .draft-card {
    flex-direction: column;
  }

  .draft-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .draft-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style> 