<template>
  <nav class="navbar">
    <div class="navbar-container">
      <router-link to="/" class="logo">
        <span class="logo-text">Notex</span>
      </router-link>

      <div class="nav-links">
        <router-link 
          v-for="item in navItems" 
          :key="item.path" 
          :to="item.path"
          class="nav-item"
          active-class="active"
        >
          <el-icon class="nav-icon">
            <component :is="item.icon" />
          </el-icon>
          {{ item.name }}
        </router-link>
      </div>

      <div class="nav-right">
        <el-button 
          type="primary" 
          class="write-btn"
          @click="$router.push('/posts/new')"
        >
          <el-icon><EditPen /></el-icon>
          写文章
        </el-button>

        <el-dropdown trigger="click" class="user-dropdown">
          <div class="user-info">
            <el-avatar :size="32" :src="userAvatar">
              {{ userInitials }}
            </el-avatar>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="$router.push('/profile')">
                <el-icon><User /></el-icon>
                个人中心
              </el-dropdown-item>
              <el-dropdown-item @click="$router.push('/drafts')">
                <el-icon><Document /></el-icon>
                草稿箱
              </el-dropdown-item>
              <el-dropdown-item @click="$router.push('/settings')">
                <el-icon><Setting /></el-icon>
                设置
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  House,
  Document,
  Folder,
  Collection,
  EditPen,
  User,
  Setting,
  SwitchButton
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const navItems = [
  { name: '首页', path: '/', icon: House },
  { name: '文章', path: '/posts', icon: Document },
  { name: '归档', path: '/archives', icon: Collection },
  { name: '分类', path: '/categories', icon: Folder },
  { name: '标签', path: '/tags', icon: Collection }
]

// 模拟用户数据，实际应该从状态管理或API获取
const userAvatar = ref('')
const userName = ref('用户')

const userInitials = computed(() => {
  return userName.value.charAt(0).toUpperCase()
})

const handleLogout = async () => {
  try {
    await userStore.logoutUser()
    ElMessage.success('退出登录成功')
    router.push('/public')
  } catch (error) {
    ElMessage.error('退出登录失败')
  }
}
</script>

<style scoped lang="scss">
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  z-index: 100;

  .navbar-container {
    max-width: 1200px;
    height: 100%;
    margin: 0 auto;
    padding: 0 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .logo {
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 8px;

    .logo-text {
      font-size: 1.5em;
      font-weight: 600;
      color: #409eff;
      letter-spacing: -0.5px;
    }
  }

  .nav-links {
    display: flex;
    align-items: center;
    gap: 8px;

    .nav-item {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 8px 16px;
      color: #606266;
      text-decoration: none;
      font-size: 0.95em;
      border-radius: 8px;
      transition: all 0.3s ease;

      .nav-icon {
        font-size: 1.1em;
      }

      &:hover {
        color: #409eff;
        background: rgba(64, 158, 255, 0.1);
      }

      &.active {
        color: #409eff;
        background: rgba(64, 158, 255, 0.1);
        font-weight: 500;
      }
    }
  }

  .nav-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .write-btn {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 8px 16px;
      border-radius: 8px;
      transition: all 0.3s ease;

      .el-icon {
        font-size: 1.1em;
      }

      &:hover {
        transform: translateY(-1px);
      }
    }

    .user-dropdown {
      .user-info {
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 2px;
        border-radius: 50%;
        transition: all 0.3s ease;

        &:hover {
          background: rgba(0, 0, 0, 0.05);
        }
      }
    }
  }
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  
  .el-icon {
    font-size: 1.1em;
  }
}
</style> 