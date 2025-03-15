<template>
  <div class="ai-toolbox-container">
    <div class="ai-toolbox-header">
      <h1>AI 工具箱</h1>
      <p>使用多种AI模型，提升您的工作效率</p>
    </div>

    <el-tabs v-model="activeTab" class="ai-toolbox-tabs">
      <el-tab-pane label="AI 聊天" name="chat">
        <AIChat />
      </el-tab-pane>
      <el-tab-pane label="AI 设置" name="settings">
        <AISettings />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import AIChat from '@/components/ai/AIChat.vue'
import AISettings from '@/components/ai/AISettings.vue'

const activeTab = ref('chat')
</script>

<style scoped lang="scss">
.ai-toolbox-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  padding-top: 88px;
}

.ai-toolbox-header {
  margin-bottom: 24px;
  
  h1 {
    font-size: 28px;
    font-weight: 700;
    margin-bottom: 8px;
    background: linear-gradient(135deg, #2B5876, #4E4376);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  
  p {
    color: #6b7280;
    font-size: 16px;
  }
}

.ai-toolbox-tabs {
  background: white;
  border-radius: 12px;
  box-shadow: 
    0 1px 3px rgba(0, 0, 0, 0.02),
    0 1px 2px rgba(0, 0, 0, 0.04);
  padding: 20px;
  
  :deep(.el-tabs__header) {
    margin-bottom: 20px;
  }
  
  :deep(.el-tabs__nav-wrap::after) {
    height: 1px;
    background-color: rgba(0, 0, 0, 0.05);
  }
  
  :deep(.el-tabs__item) {
    font-size: 16px;
    font-weight: 500;
    padding: 0 20px;
    height: 40px;
    line-height: 40px;
    
    &.is-active {
      color: #2B5876;
    }
  }
  
  :deep(.el-tabs__active-bar) {
    background-color: #2B5876;
    height: 3px;
    border-radius: 3px;
  }
}
</style> 