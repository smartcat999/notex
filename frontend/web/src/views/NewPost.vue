<template>
  <div class="new-post-container">
    <div class="page-header">
      <div class="title-section">
        <h1>写文章</h1>
      </div>
    </div>

    <div class="post-form">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" @input="handleFormChange" />
        </el-form-item>

        <el-form-item label="分类" prop="category_id">
          <div class="category-select">
            <el-select 
              v-model="form.category_id" 
              placeholder="请选择分类"
              class="category-input"
              value-key="id"
              @change="handleFormChange"
            >
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
            <el-button type="primary" plain size="default" @click="showCreateCategory">
              <el-icon><Plus /></el-icon>
              创建分类
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="标签" prop="tag_ids">
          <div class="tag-select">
            <el-select
              v-model="form.tag_ids"
              multiple
              filterable
              class="tag-input"
              placeholder="请选择标签"
              value-key="id"
              @change="handleFormChange"
            >
              <el-option
                v-for="tag in tags"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
            <el-button type="primary" plain size="default" @click="showCreateTag">
              <el-icon><Plus /></el-icon>
              创建标签
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="封面图" prop="cover">
          <file-upload v-model="form.cover" @update:modelValue="handleFormChange" />
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input
            v-model="form.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要"
            @input="handleFormChange"
          />
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <markdown-editor v-model="form.content" @update:modelValue="handleFormChange" />
        </el-form-item>

        <el-form-item>
          <div class="form-actions">
            <el-button plain @click="handleCancel">取消</el-button>
            <el-button type="info" plain @click="handlePreview">
              <el-icon><View /></el-icon>
              预览
            </el-button>
            <el-button type="primary" @click="handleSubmit">发布</el-button>
            <el-button type="info" plain @click="handleSaveDraft">存为草稿</el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      title="文章预览"
      width="800px"
      :close-on-click-modal="false"
      class="preview-dialog"
      destroy-on-close
      fullscreen
    >
      <div class="preview-content">
        <div class="post-header" :class="{ 'no-cover': !form.cover }">
          <div class="post-cover" v-if="form.cover">
            <img :src="form.cover" :alt="form.title">
          </div>
          <div class="header-content">
            <h1>{{ form.title }}</h1>
            <div class="post-meta">
              <span>
                <el-icon><Calendar /></el-icon>
                {{ formatDate(new Date()) }}
              </span>
              <span>
                <el-icon><View /></el-icon>
                0
              </span>
              <span>
                <el-icon><ChatDotRound /></el-icon>
                0
              </span>
            </div>
            <!-- 作者信息 -->
            <AuthorCard :author="userStore.user" />
            <div v-if="form.tag_ids?.length" class="post-tags">
              <el-tag
                v-for="tagId in form.tag_ids"
                :key="tagId"
                size="small"
                class="tag"
              >
                {{ getTagName(tagId) }}
              </el-tag>
            </div>
          </div>
        </div>

        <div class="post-content markdown-body">
          <MarkdownPreview :content="form.content" />
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="previewDialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 创建分类对话框 -->
    <el-dialog
      v-model="categoryDialogVisible"
      title="创建分类"
      width="460px"
      :close-on-click-modal="false"
      class="custom-dialog"
      destroy-on-close
    >
      <el-form
        ref="categoryFormRef"
        :model="categoryForm"
        :rules="categoryRules"
        label-position="top"
      >
        <el-form-item label="名称" prop="name">
          <el-input 
            v-model="categoryForm.name" 
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="categoryForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
            maxlength="200"
            show-word-limit
            resize="none"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="categoryDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateCategory">
            创建分类
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 创建标签对话框 -->
    <el-dialog
      v-model="tagDialogVisible"
      title="创建标签"
      width="460px"
      :close-on-click-modal="false"
      class="custom-dialog"
      destroy-on-close
    >
      <el-form
        ref="tagFormRef"
        :model="tagForm"
        :rules="tagRules"
        label-position="top"
      >
        <el-form-item label="名称" prop="name">
          <el-input 
            v-model="tagForm.name" 
            placeholder="请输入标签名称"
            maxlength="30"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="tagDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateTag">
            创建标签
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, View, Calendar, User, Folder, ChatDotRound } from '@element-plus/icons-vue'
import { 
  getCategories, 
  getTags, 
  createPost, 
  createCategory,
  createTag 
} from '@/api/posts'
import { createDraft, getDraft, updateDraft, publishDraft } from '@/api/drafts'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import FileUpload from '@/components/FileUpload.vue'
import { formatDate } from '@/utils/date'
import { useUserStore } from '@/stores/user'
import AuthorCard from '@/components/AuthorCard.vue'

const router = useRouter()
const route = useRoute()
const formRef = ref(null)
const categoryFormRef = ref(null)
const tagFormRef = ref(null)
const categories = ref([])
const tags = ref([])
const categoryDialogVisible = ref(false)
const tagDialogVisible = ref(false)
const isEditMode = ref(false)
const hasChanges = ref(false)
const autoSaveTimer = ref(null)
const previewDialogVisible = ref(false)

const userStore = useUserStore()

// 文章表单数据
const form = ref({
  title: '',
  category_id: '',
  tag_ids: [],
  summary: '',
  content: '',
  cover: '',
})

// 分类表单数据
const categoryForm = ref({
  name: '',
  description: ''
})

// 标签表单数据
const tagForm = ref({
  name: ''
})

// 文章表单验证规则
const rules = {
  title: [{ required: true, message: '请输入文章标题', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入文章内容', trigger: 'blur' }]
}

// 分类表单验证规则
const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

// 标签表单验证规则
const tagRules = {
  name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }]
}

// 监听表单变化
const handleFormChange = () => {
  hasChanges.value = true
  startAutoSave()
}

// 自动保存
const startAutoSave = () => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
  autoSaveTimer.value = setTimeout(async () => {
    if (hasChanges.value) {
      await handleSaveDraft()
      hasChanges.value = false
    }
  }, 60000) // 每60秒自动保存一次
}

// 路由离开守卫
const handleBeforeUnload = (e) => {
  if (hasChanges.value) {
    e.preventDefault()
    e.returnValue = ''
  }
}

// 注册路由守卫
onBeforeUnmount(() => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload)
  Promise.all([
    fetchCategories(),
    fetchTags()
  ])
  
  const draftId = route.query.draft_id
  if (draftId) {
    loadDraft(draftId)
  }
})

// 取消编辑
const handleCancel = async () => {
  if (hasChanges.value) {
    try {
      await ElMessageBox.confirm(
        '是否将当前内容保存为草稿？',
        '提示',
        {
          confirmButtonText: '保存草稿',
          cancelButtonText: '放弃更改',
          type: 'warning',
        }
      )
      await handleSaveDraft()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('操作失败')
        return
      }
    }
  }
  router.back()
}

// 保存草稿
const handleSaveDraft = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    if (isEditMode.value) {
      await updateDraft(route.query.draft_id, form.value)
      ElMessage.success('草稿更新成功')
    } else {
      await createDraft(form.value)
      ElMessage.success('草稿保存成功')
    }
    hasChanges.value = false
    router.push('/drafts')
    return true
  } catch (error) {
    console.error('Failed to save draft:', error)
    ElMessage.error(isEditMode.value ? '草稿更新失败' : '草稿保存失败')
    return false
  }
}

// 发布文章
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    if (isEditMode.value) {
      await publishDraft(route.query.draft_id)
      ElMessage.success('草稿发布成功')
    } else {
      const timestamp = new Date().getTime()
      const slug = `${generateSlug(form.value.title)}-${timestamp}`
      await createPost({
        ...form.value,
        status: 'published',
        slug: slug
      })
      ElMessage.success('文章发布成功')
    }
    hasChanges.value = false
    router.push('/posts')
  } catch (error) {
    console.error('Failed to publish:', error)
    ElMessage.error(isEditMode.value ? '草稿发布失败' : '文章发布失败')
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await getCategories()
    categories.value = response.items || []
  } catch (error) {
    console.error('Failed to fetch categories:', error)
    ElMessage.error('获取分类列表失败')
    categories.value = []
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await getTags()
    tags.value = response.items || []
  } catch (error) {
    console.error('Failed to fetch tags:', error)
    ElMessage.error('获取标签列表失败')
    tags.value = []
  }
}

// 监听分类选择变化
const handleCategoryChange = (value) => {
  console.log('Selected category:', value)
  const selectedCategory = categories.value.find(c => c.id === value)
  console.log('Selected category object:', selectedCategory)
}

// 监听标签选择变化
const handleTagChange = (value) => {
  console.log('Selected tags:', value)
  const selectedTags = tags.value.filter(t => value.includes(t.id))
  console.log('Selected tag objects:', selectedTags)
}

// 显示创建分类对话框
const showCreateCategory = () => {
  categoryForm.value = {
    name: '',
    description: ''
  }
  categoryDialogVisible.value = true
}

// 显示创建标签对话框
const showCreateTag = () => {
  tagForm.value = {
    name: ''
  }
  tagDialogVisible.value = true
}

// 创建分类
const handleCreateCategory = async () => {
  if (!categoryFormRef.value) return

  try {
    await categoryFormRef.value.validate()
    await createCategory(categoryForm.value)
    ElMessage.success('分类创建成功')
    categoryDialogVisible.value = false
    await fetchCategories() // 重新获取分类列表
  } catch (error) {
    console.error('Failed to create category:', error)
    ElMessage.error(error.response?.data?.message || '创建分类失败')
  }
}

// 创建标签
const handleCreateTag = async () => {
  if (!tagFormRef.value) return

  try {
    await tagFormRef.value.validate()
    await createTag(tagForm.value)
    ElMessage.success('标签创建成功')
    tagDialogVisible.value = false
    await fetchTags() // 重新获取标签列表
  } catch (error) {
    console.error('Failed to create tag:', error)
    ElMessage.error(error.response?.data?.message || '创建标签失败')
  }
}

// 生成 slug 的辅助函数
const generateSlug = (title) => {
  return title
    .toLowerCase()
    .trim()
    .replace(/[^\w\u4e00-\u9fa5]+/g, '-') // 将非字母数字和中文字符替换为连字符
    .replace(/^-+|-+$/g, '') // 移除首尾的连字符
    .substring(0, 100) // 限制长度
}

// 加载草稿数据
const loadDraft = async (draftId) => {
  try {
    const draft = await getDraft(draftId)
    form.value = {
      title: draft.title,
      category_id: draft.category_id,
      tag_ids: draft.tags.map(tag => tag.id),
      summary: draft.summary,
      content: draft.content,
      cover: draft.cover,
    }
    isEditMode.value = true
  } catch (error) {
    console.error('Failed to load draft:', error)
    ElMessage.error('加载草稿失败')
    router.push('/drafts')
  }
}

// 获取分类名称
const getCategoryName = (categoryId) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : ''
}

// 获取标签名称
const getTagName = (tagId) => {
  const tag = tags.value.find(t => t.id === tagId)
  return tag ? tag.name : ''
}

// 处理预览
const handlePreview = () => {
  previewDialogVisible.value = true
}
</script>

<style lang="scss" scoped>
@import '@/styles/post-content.scss';

.new-post-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .page-header {
    margin-bottom: 24px;

    .title-section {
      h1 {
        font-size: 2rem;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin: 0;
        line-height: 1.2;
      }
    }
  }

  .post-form {
    background: var(--el-bg-color);
    border-radius: 8px;
    padding: 2rem;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);

    :deep(.el-form-item__label) {
      font-weight: 500;
    }

    .category-select,
    .tag-select {
      display: flex;
      gap: 16px;
      align-items: flex-start;

      .category-input,
      .tag-input {
        flex: 0 1 300px;
        min-width: 200px;
        max-width: 300px;
      }

      :deep(.el-select) {
        width: 100%;
      }

      :deep(.el-select__tags) {
        flex-wrap: wrap;
        max-width: 100%;
        padding: 0 4px;
      }

      .el-button {
        flex: 0 0 auto;
        white-space: nowrap;
        
        &.el-button--primary {
          --el-button-bg-color: transparent;
          --el-button-border-color: var(--slate-6, #e2e8f0);
          --el-button-text-color: var(--slate-11, #64748b);
          --el-button-hover-text-color: var(--slate-12, #1e293b);
          --el-button-hover-bg-color: var(--slate-3, #f1f5f9);
          --el-button-hover-border-color: var(--slate-7, #cbd5e1);
          --el-button-active-bg-color: var(--slate-4, #e2e8f0);
          
          .el-icon {
            margin-right: 4px;
            font-size: 16px;
          }
        }
      }
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 16px;
      margin-top: 24px;

      :deep(.el-button) {
        padding: 12px 24px;
        font-weight: 500;
        transition: all 0.3s ease;
        border-radius: 8px;
        
        &.el-button--primary {
          --el-button-bg-color: transparent;
          background: linear-gradient(to right, #3B4B66, #2B3A53);
          border: 1px solid rgba(255, 255, 255, 0.1);
          color: rgba(255, 255, 255, 0.9);
          box-shadow: 
            0 1px 2px rgba(0, 0, 0, 0.1),
            0 1px 1px rgba(0, 0, 0, 0.06);
          backdrop-filter: blur(8px);
          
          &:hover {
            transform: translateY(-1px);
            background: linear-gradient(to right, #445373, #324161);
            border-color: rgba(255, 255, 255, 0.15);
            color: rgba(255, 255, 255, 1);
            box-shadow: 
              0 4px 12px rgba(0, 0, 0, 0.1),
              0 2px 4px rgba(0, 0, 0, 0.08);
          }

          &:active {
            transform: translateY(0);
            background: linear-gradient(to right, #324161, #263450);
            box-shadow: 
              0 1px 2px rgba(0, 0, 0, 0.1),
              0 1px 1px rgba(0, 0, 0, 0.06);
          }
        }
        
        &.el-button--info {
          &.is-plain {
            --el-button-bg-color: transparent;
            --el-button-border-color: var(--slate-6, #e2e8f0);
            --el-button-text-color: var(--slate-11, #64748b);
            --el-button-hover-text-color: var(--slate-12, #1e293b);
            --el-button-hover-bg-color: var(--slate-3, #f1f5f9);
            --el-button-hover-border-color: var(--slate-7, #cbd5e1);
          }
        }

        &.is-plain {
          --el-button-hover-text-color: var(--slate-12, #1e293b);
          --el-button-hover-bg-color: var(--slate-3, #f1f5f9);
          --el-button-hover-border-color: var(--slate-7, #cbd5e1);
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

:deep(.custom-dialog) {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 
    0px 10px 38px -10px rgba(22, 23, 24, 0.35),
    0px 10px 20px -15px rgba(22, 23, 24, 0.2);

  .el-dialog__header {
    margin: 0;
    padding: 20px 24px;
    border-bottom: 1px solid rgba(22, 23, 24, 0.06);

    .el-dialog__title {
      font-size: 16px;
      font-weight: 600;
      color: var(--slate-12, #1f2937);
      line-height: 1.4;
    }

    .el-dialog__close {
      font-size: 16px;
      color: var(--slate-11, #6b7280);
      transition: all 0.2s ease;

      &:hover {
        color: var(--slate-12, #1f2937);
        transform: scale(1.1);
      }
    }
  }

  .el-dialog__body {
    padding: 24px;

    .el-form-item {
      margin-bottom: 20px;

      &:last-child {
        margin-bottom: 0;
      }

      .el-form-item__label {
        padding: 0 0 8px;
        font-size: 14px;
        font-weight: 500;
        color: var(--slate-12, #1f2937);
        line-height: 1.4;
      }

      .el-input__wrapper,
      .el-textarea__wrapper {
        box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.08);
        border-radius: 6px;
        padding: 8px 12px;
        transition: all 0.2s ease;

        &:hover {
          box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.12);
        }

        &.is-focus {
          box-shadow: 0 0 0 2px var(--el-color-primary);
        }
      }

      .el-input__inner,
      .el-textarea__inner {
        font-size: 14px;
        color: var(--slate-12, #1f2937);
        
        &::placeholder {
          color: var(--slate-9, #9ca3af);
        }
      }

      .el-input__count {
        font-size: 12px;
        color: var(--slate-10, #6b7280);
        background: transparent;
      }
    }
  }

  .el-dialog__footer {
    padding: 16px 24px;
    border-top: 1px solid rgba(22, 23, 24, 0.06);
    margin-top: 0;

    .dialog-footer {
      display: flex;
      justify-content: flex-end;
      gap: 12px;

      .el-button {
        padding: 8px 16px;
        font-size: 14px;
        font-weight: 500;
        border-radius: 6px;
        transition: all 0.2s ease;

        &--default {
          color: var(--slate-11, #6b7280);
          border-color: var(--slate-7, #e5e7eb);
          background: transparent;

          &:hover {
            color: var(--slate-12, #1f2937);
            border-color: var(--slate-8, #d1d5db);
            background: var(--slate-3, #f9fafb);
          }

          &:active {
            background: var(--slate-4, #f3f4f6);
          }
        }

        &--primary {
          background: var(--slate-3, #f9fafb);
          border-color: var(--slate-7, #e5e7eb);
          color: var(--slate-11, #6b7280);
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

          &:hover {
            background: var(--slate-4, #f3f4f6);
            border-color: var(--slate-8, #d1d5db);
            color: var(--slate-12, #1f2937);
            transform: translateY(-1px);
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
          }

          &:active {
            transform: translateY(0);
            background: var(--slate-5, #e5e7eb);
          }
        }
      }
    }
  }
}

.preview-dialog {
  :deep(.el-dialog__body) {
    padding: 0;
  }

  .preview-content {
    max-width: 900px;
    margin: 0 auto;
    padding: 40px 20px;

    .post-meta {
      width: 100%;
      color: #6b7280;
      font-size: 0.9em;
      margin-bottom: 24px;
      display: flex;
      justify-content: center;
      flex-wrap: wrap;
      gap: 16px;

      span {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        padding: 6px 12px;
        background: rgba(43, 88, 118, 0.06);
        border-radius: 20px;
        transition: all 0.3s ease;

        &:hover {
          background: rgba(43, 88, 118, 0.1);
          transform: translateY(-1px);
        }

        .el-icon {
          font-size: 1.1em;
          color: #2B5876;
        }
      }
    }

    :deep(.author-card) {
      width: 100%;
      margin: 0 0 24px;
      display: flex;
      justify-content: center;
      border: none;
      padding: 0;

      .author-content {
        background: rgba(43, 88, 118, 0.03);
        padding: 12px 20px;
        border-radius: 16px;
        transition: all 0.3s ease;

        &:hover {
          background: rgba(43, 88, 118, 0.06);
          transform: translateY(-1px);
        }
      }
    }
  }
}
</style> 