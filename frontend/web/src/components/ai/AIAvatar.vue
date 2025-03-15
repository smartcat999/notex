<template>
  <el-avatar :size="size" :class="['ai-avatar', providerId]">
    <template v-if="customIcon">
      {{ customIcon }}
    </template>
    <template v-else>
      <img v-if="avatarSrc" :src="avatarSrc" alt="AI Avatar" />
      <span v-else>AI</span>
    </template>
  </el-avatar>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelId: {
    type: String,
    default: ''
  },
  providerId: {
    type: String,
    default: ''
  },
  size: {
    type: Number,
    default: 36
  },
  customIcon: {
    type: String,
    default: ''
  }
})

// 根据提供商和模型获取头像
const avatarSrc = computed(() => {
  if (!props.providerId) return null
  
  // 这里可以根据不同的提供商和模型返回不同的头像
  switch (props.providerId) {
    case 'openai':
      return '/images/ai/openai-avatar.png'
    case 'anthropic':
      return '/images/ai/anthropic-avatar.png'
    case 'google':
      return '/images/ai/google-avatar.png'
    case 'deepseek':
      return '/images/ai/deepseek-avatar.svg'
    case 'custom':
      return '/images/ai/custom-avatar.png'
    default:
      return null
  }
})
</script>

<style scoped lang="scss">
.ai-avatar {
  background: linear-gradient(135deg, #2B5876, #4E4376);
  color: white;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &.openai {
    background: linear-gradient(135deg, #10a37f, #0e8a6d);
  }
  
  &.anthropic {
    background: linear-gradient(135deg, #b83280, #9b2c6f);
  }
  
  &.google {
    background: linear-gradient(135deg, #4285f4, #34a853, #fbbc05, #ea4335);
  }
  
  &.deepseek {
    background: linear-gradient(135deg, #1a73e8, #0d47a1);
  }
  
  &.custom {
    background: linear-gradient(135deg, #6b7280, #4b5563);
  }
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style> 