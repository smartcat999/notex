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
        <el-menu-item index="/ai/ai-writing">
          <el-icon><Edit /></el-icon>
          <span>AI写作</span>
        </el-menu-item>
        <el-menu-item index="/ai/ai-document">
          <el-icon><Document /></el-icon>
          <span>AI文档</span>
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
import { ChatDotRound, Setting, Document, Edit } from '@element-plus/icons-vue'

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
  background-color: #f8fafc;
  color: #1e293b;
  
  .ai-sidebar {
    width: 260px;
    background-color: #ffffff;
    border-right: 1px solid #e2e8f0;
    display: flex;
    flex-direction: column;
    
    .sidebar-header {
      padding: 20px;
      border-bottom: 1px solid #e2e8f0;
      background: #ffffff;
      
      .sidebar-title {
        font-size: 20px;
        font-weight: 600;
        background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin: 0;
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
        color: #64748b;
        margin: 4px 0;
        border-radius: 8px;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        
        .el-icon {
          color: #64748b;
          margin-right: 12px;
          font-size: 18px;
          transition: all 0.3s ease;
        }
        
        span {
          font-size: 15px;
        }
        
        &:hover {
          background-color: rgba(139, 92, 246, 0.1);
          color: #8b5cf6;
          transform: translateX(4px);
          
          .el-icon {
            color: #8b5cf6;
          }
        }
        
        &.is-active {
          background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
          color: #ffffff;
          box-shadow: 0 2px 8px rgba(139, 92, 246, 0.25);
          
          .el-icon {
            color: #ffffff;
          }
        }
      }
    }
  }
  
  .ai-content {
    flex: 1;
    background-color: #f8fafc;
    overflow: hidden;
    position: relative;
  }
}
</style> 