<template>
  <div class="ai-layout">
    <div class="ai-sidebar">
      <div class="sidebar-header">
        <h2 class="sidebar-title">AI 工具箱</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="ai-menu"
        @select="handleSelect"
        :router="true"
      >
        <el-menu-item index="/ai/chat">
          <el-icon><ChatDotRound /></el-icon>
          <span>聊天</span>
        </el-menu-item>
        <el-menu-item index="/ai/word-polish">
          <el-icon><Document /></el-icon>
          <span>Word优化</span>
        </el-menu-item>
        <el-menu-item index="/ai/settings">
          <el-icon><Setting /></el-icon>
          <span>设置</span>
        </el-menu-item>
      </el-menu>
    </div>
    <div class="ai-content">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChatDotRound, Setting, Document } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const activeMenu = computed(() => route.path)

const handleSelect = (key) => {
  router.push(key)
}
</script>

<style scoped lang="scss">
.ai-layout {
  display: flex;
  min-height: calc(100vh - 64px);
  background-color: #171717;
  color: #e0e0e0;
  
  .ai-sidebar {
    width: 260px;
    background-color: #1f1f1f;
    border-right: 1px solid #2a2a2a;
    display: flex;
    flex-direction: column;
    
    .sidebar-header {
      padding: 20px;
      border-bottom: 1px solid #2a2a2a;
      background: linear-gradient(180deg, #1f1f1f 0%, #252525 100%);
      
      .sidebar-title {
        font-size: 1.2em;
        font-weight: 600;
        color: #ffffff;
        margin: 0;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
      }
    }
    
    .ai-menu {
      flex: 1;
      border-right: none;
      background-color: transparent;
      padding: 8px;
      
      :deep(.el-menu-item) {
        height: 44px;
        line-height: 44px;
        color: #a0a0a0;
        margin: 4px 0;
        border-radius: 8px;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        
        .el-icon {
          color: #a0a0a0;
          margin-right: 12px;
          font-size: 18px;
          transition: all 0.3s ease;
        }
        
        span {
          font-size: 14px;
        }
        
        &:hover {
          background-color: #2a2a2a;
          color: #ffffff;
          transform: translateX(4px);
          
          .el-icon {
            color: #409EFF;
          }
        }
        
        &.is-active {
          background: linear-gradient(90deg, #2B5876 0%, #4E4376 100%);
          color: #ffffff;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
          
          .el-icon {
            color: #ffffff;
          }
        }
      }
    }
  }
  
  .ai-content {
    flex: 1;
    background-color: #171717;
    overflow: hidden;
    position: relative;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(90deg, 
        rgba(43, 88, 118, 0.3) 0%, 
        rgba(78, 67, 118, 0.3) 100%);
    }
  }
}
</style> 