<template>
  <FileUploader
    :modelValue="modelValue"
    @update:modelValue="$emit('update:modelValue', $event)"
    :accept="'image/jpeg,image/jpg,image/png,image/gif'"
    :max-size="5 * 1024 * 1024"
    :upload-text="'点击或拖拽图片上传'"
    :upload-tip="'支持 jpg、jpeg、png、gif 格式'"
    @upload-success="handleUploadSuccess"
    @upload-error="handleUploadError"
  />
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
import FileUploader from '@/components/common/FileUploader.vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const handleUploadSuccess = (url) => {
  emit('update:modelValue', url)
}

const handleUploadError = (error) => {
  console.error('Upload error:', error)
}
</script> 