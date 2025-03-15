<template>
  <div class="ai-sidebar">
    <el-menu
      :default-active="activeMenu"
      class="ai-menu"
      @select="handleSelect"
    >
      <el-menu-item index="chat">
        <el-icon><ChatDotRound /></el-icon>
        <span>AI 聊天</span>
      </el-menu-item>
      <el-menu-item index="settings">
        <el-icon><Setting /></el-icon>
        <span>AI 设置</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ChatDotRound, Setting } from '@element-plus/icons-vue'

const router = useRouter()
const activeMenu = ref('chat')

const handleSelect = (index) => {
  router.push(`/ai/${index}`)
}
</script>

<style scoped lang="scss">
.ai-sidebar {
  width: 200px;
  height: 100%;
  border-right: 1px solid #e5e7eb;
  
  .ai-menu {
    border-right: none;
    
    .el-menu-item {
      display: flex;
      align-items: center;
      padding: 0 20px;
      height: 50px;
      
      .el-icon {
        margin-right: 12px;
        font-size: 18px;
      }
      
      span {
        font-size: 14px;
      }
      
      &.is-active {
        background-color: #e9f5fe;
        color: #1d4ed8;
        
        &::before {
          content: '';
          position: absolute;
          left: 0;
          top: 0;
          bottom: 0;
          width: 3px;
          background-color: #1d4ed8;
        }
      }
    }
  }
}
</style> 