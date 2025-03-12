<template>
  <div class="file-uploader" :class="{ 'is-avatar': isAvatar }">
    <div class="upload-container">
      <slot name="upload-area" v-if="!modelValue">
        <div :class="['upload-area', { 'is-avatar': isAvatar }]" @click="triggerUpload">
          <el-icon class="upload-icon"><Upload /></el-icon>
          <div class="upload-text">{{ uploadText }}</div>
          <div class="upload-tip">{{ uploadTip }}</div>
        </div>
      </slot>
      <slot name="preview" v-if="modelValue" :url="modelValue">
        <div :class="['preview-area', { 'is-avatar': isAvatar }]">
          <img :src="modelValue" class="preview-image" />
          <div class="preview-actions">
            <el-button type="primary" link @click.stop="triggerUpload">
              <el-icon><Upload /></el-icon>
              重新上传
            </el-button>
            <el-button type="danger" link @click.stop="handleRemove">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </div>
      </slot>
    </div>
    <input
      ref="fileInput"
      type="file"
      :accept="accept"
      style="display: none"
      @change="handleFileChange"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Delete } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { uploadFile } from '@/utils/upload'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  isAvatar: {
    type: Boolean,
    default: false
  },
  accept: {
    type: String,
    default: 'image/*'
  },
  maxSize: {
    type: Number,
    default: 5 * 1024 * 1024 // 默认5MB
  },
  uploadText: {
    type: String,
    default: '点击或拖拽文件上传'
  },
  uploadTip: {
    type: String,
    default: '支持 jpg、png、gif 格式'
  },
  compressOptions: {
    type: Object,
    default: () => ({
      maxWidth: 1920,
      maxHeight: 1080,
      quality: 0.8,
      maxSizeMB: 1
    })
  }
})

const emit = defineEmits(['update:modelValue'])
const userStore = useUserStore()
const fileInput = ref(null)

// 压缩图片
const compressImage = async (file) => {
  // 如果不是图片文件，直接返回原文件
  if (!file.type.startsWith('image/')) {
    return file
  }

  // 如果文件小于1MB且不是头像，直接返回
  if (file.size <= props.compressOptions.maxSizeMB * 1024 * 1024 && !props.isAvatar) {
    return file
  }

  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = (e) => {
      const img = new Image()
      img.src = e.target.result
      img.onload = () => {
        const canvas = document.createElement('canvas')
        let width = img.width
        let height = img.height

        // 计算缩放比例
        if (props.isAvatar) {
          // 头像固定尺寸为 200x200
          width = height = 200
        } else {
          // 普通图片按比例缩放
          const maxWidth = props.compressOptions.maxWidth
          const maxHeight = props.compressOptions.maxHeight
          if (width > maxWidth || height > maxHeight) {
            const ratio = Math.min(maxWidth / width, maxHeight / height)
            width *= ratio
            height *= ratio
          }
        }

        canvas.width = width
        canvas.height = height
        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, width, height)

        // 转换为Blob
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('Canvas to Blob failed'))
              return
            }
            // 创建新的文件对象
            const compressedFile = new File([blob], file.name, {
              type: file.type,
              lastModified: Date.now()
            })
            resolve(compressedFile)
          },
          file.type,
          props.compressOptions.quality
        )
      }
      img.onerror = (error) => reject(error)
    }
    reader.onerror = (error) => reject(error)
  })
}

// 触发文件选择
const triggerUpload = () => {
  fileInput.value.click()
}

// 处理文件选择
const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  try {
    // 压缩图片
    const processedFile = await compressImage(file)
    
    const fileUrl = await uploadFile(processedFile, userStore.token, {
      maxSize: props.maxSize,
      allowedTypes: props.accept.split(',')
    })
    
    emit('update:modelValue', fileUrl)
    ElMessage.success('上传成功')
  } catch (error) {
    ElMessage.error(error.message || '上传失败，请重试')
  } finally {
    // 清空文件输入框，以便能够重新选择同一个文件
    event.target.value = ''
  }
}

// 删除文件
const handleRemove = () => {
  emit('update:modelValue', '')
}

// 暴露方法给父组件
defineExpose({
  triggerUpload
})
</script>

<style lang="scss" scoped>
.file-uploader {
  width: 100%;
  
  &.is-avatar {
    width: 120px;
    height: 120px;
  }
  
  :deep(.el-upload) {
    width: 100%;
    
    .el-upload-dragger {
      width: 100%;
      height: 100%;
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
    
    &.is-avatar {
      height: 120px;
      border-radius: 50%;
      
      .upload-icon {
        font-size: 24px;
      }
      
      .upload-text {
        font-size: 12px;
      }
      
      .upload-tip {
        display: none;
      }
    }
    
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
    height: 200px;
    position: relative;
    border-radius: 8px;
    overflow: hidden;
    
    &.is-avatar {
      height: 120px;
      border-radius: 50%;
      
      .preview-image {
        height: 120px;
      }
    }
    
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