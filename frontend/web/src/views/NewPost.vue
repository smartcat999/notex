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
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>

        <el-form-item label="分类" prop="category_id">
          <div class="category-select">
            <el-select 
              v-model="form.category_id" 
              placeholder="请选择分类"
              class="category-input"
              value-key="id"
            >
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
            <el-button type="primary" @click="showCreateCategory">
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
            >
              <el-option
                v-for="tag in tags"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
            <el-button type="primary" @click="showCreateTag">
              <el-icon><Plus /></el-icon>
              创建标签
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input
            v-model="form.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要"
          />
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <div class="editor-container">
            <!-- 这里后续可以集成 Markdown 编辑器 -->
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="15"
              placeholder="请输入文章内容"
            />
          </div>
        </el-form-item>

        <el-form-item>
          <div class="form-actions">
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSubmit">发布</el-button>
            <el-button type="info" @click="handleSaveDraft">存为草稿</el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <!-- 创建分类对话框 -->
    <el-dialog
      v-model="categoryDialogVisible"
      title="创建分类"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="categoryFormRef"
        :model="categoryForm"
        :rules="categoryRules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="categoryForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="categoryDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateCategory">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 创建标签对话框 -->
    <el-dialog
      v-model="tagDialogVisible"
      title="创建标签"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="tagFormRef"
        :model="tagForm"
        :rules="tagRules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="tagDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateTag">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { 
  getCategories, 
  getTags, 
  createPost, 
  createCategory,
  createTag 
} from '@/api/posts'

const router = useRouter()
const formRef = ref(null)
const categoryFormRef = ref(null)
const tagFormRef = ref(null)
const categories = ref([])
const tags = ref([])
const categoryDialogVisible = ref(false)
const tagDialogVisible = ref(false)

// 文章表单数据
const form = ref({
  title: '',
  category_id: '',
  tag_ids: [],
  summary: '',
  content: '',
  status: 'published'
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

// 发布文章
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    await createPost(form.value)
    ElMessage.success('文章发布成功')
    router.push('/posts')
  } catch (error) {
    console.error('Failed to create post:', error)
    ElMessage.error('文章发布失败')
  }
}

// 保存草稿
const handleSaveDraft = async () => {
  form.value.status = 'draft'
  await handleSubmit()
}

// 取消编辑
const handleCancel = () => {
  router.back()
}

// 在组件挂载时获取分类和标签数据
onMounted(async () => {
  await Promise.all([
    fetchCategories(),
    fetchTags()
  ])
})
</script>

<style lang="scss" scoped>
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
      }
    }

    .editor-container {
      border: 1px solid var(--el-border-color);
      border-radius: 4px;
      padding: 16px;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 16px;
      margin-top: 24px;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}
</style> 