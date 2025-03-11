<template>
  <div class="file-upload">
    <el-upload
      :action="uploadUrl"
      :headers="headers"
      :show-file-list="false"
      :before-upload="handleBeforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      accept="image/*"
      :multiple="false"
    >
      <div class="upload-area" v-if="!modelValue">
        <el-icon class="upload-icon"><Upload /></el-icon>
        <div class="upload-text">点击或拖拽图片上传</div>
        <div class="upload-tip">支持 jpg、png、gif 格式</div>
      </div>
      <div v-else class="preview-area">
        <img :src="modelValue" class="preview-image" />
        <div class="preview-actions">
          <el-button type="primary" link @click.stop="handleReupload">
            <el-icon><Upload /></el-icon>
            重新上传
          </el-button>
          <el-button type="danger" link @click.stop="handleRemove">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>
    </el-upload>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Delete } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const userStore = useUserStore()

// 上传地址
const uploadUrl = '/api/upload/file'

// 请求头（包含认证信息）
const headers = computed(() => ({
  'Authorization': `Bearer ${userStore.token}`
}))

// 上传前验证
const handleBeforeUpload = (file) => {
  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('只能上传 jpg、png、gif 格式的图片！')
    return false
  }
  
  // 验证文件大小（限制为 5MB）
  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error('图片大小不能超过 5MB！')
    return false
  }
  
  return true
}

// 上传成功
const handleSuccess = (response) => {
  // 直接使用返回的数据，不需要检查 code
  if (response && response.url) {
    emit('update:modelValue', response.url)
    ElMessage.success('上传成功')
  } else {
    ElMessage.error('上传失败：无效的响应数据')
  }
}

// 上传失败
const handleError = (error) => {
  console.error('Upload error:', error)
  ElMessage.error('上传失败，请重试')
}

// 重新上传
const handleReupload = () => {
  // 保持当前预览，等待新文件上传
}

// 删除图片
const handleRemove = () => {
  emit('update:modelValue', '')
}
</script>

<style lang="scss" scoped>
.file-upload {
  width: 100%;
  
  :deep(.el-upload) {
    width: 100%;
    
    .el-upload-dragger {
      width: 100%;
      height: 200px;
    }
  }
  
  .upload-area {
    width: 100%;
    height: 200px;
    border: 2px dashed var(--el-border-color);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    transition: border-color 0.3s;
    
    &:hover {
      border-color: var(--el-color-primary);
    }
    
    .upload-icon {
      font-size: 32px;
      color: var(--el-text-color-secondary);
      margin-bottom: 8px;
    }
    
    .upload-text {
      font-size: 16px;
      color: var(--el-text-color-regular);
      margin-bottom: 4px;
    }
    
    .upload-tip {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
  
  .preview-area {
    width: 100%;
    position: relative;
    border-radius: 8px;
    overflow: hidden;
    
    &:hover .preview-actions {
      opacity: 1;
    }
    
    .preview-image {
      width: 100%;
      height: 200px;
      object-fit: cover;
    }
    
    .preview-actions {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.5);
      display: flex;
      justify-content: center;
      align-items: center;
      gap: 16px;
      opacity: 0;
      transition: opacity 0.3s;
      
      .el-button {
        color: #fff;
        
        &:hover {
          color: var(--el-color-primary);
        }
      }
    }
  }
}
</style> 