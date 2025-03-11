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

// 使用 computed 获取用户头像和用户名
const userAvatar = computed(() => userStore.user?.avatar || '')
const userName = computed(() => userStore.user?.username || '用户')

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
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  z-index: 100;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.98);
  }

  .navbar-container {
    max-width: 1200px;
    height: 100%;
    margin: 0 auto;
    padding: 0 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .logo {
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 10px;
    transition: transform 0.3s ease;

    &:hover {
      transform: translateY(-1px);
    }

    .logo-text {
      font-size: 1.6em;
      font-weight: 700;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: -0.5px;
    }
  }

  .nav-links {
    display: flex;
    align-items: center;
    gap: 4px;

    .nav-item {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      color: #4a5568;
      text-decoration: none;
      font-size: 0.95em;
      border-radius: 10px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      font-weight: 500;

      .nav-icon {
        font-size: 1.15em;
        transition: transform 0.3s ease;
      }

      &:hover {
        color: #2B5876;
        background: rgba(43, 88, 118, 0.08);

        .nav-icon {
          transform: translateY(-1px);
        }
      }

      &.active {
        color: #2B5876;
        background: rgba(43, 88, 118, 0.1);
        font-weight: 600;

        &::after {
          content: '';
          position: absolute;
          bottom: 6px;
          left: 50%;
          transform: translateX(-50%);
          width: 16px;
          height: 2px;
          background: #2B5876;
          border-radius: 2px;
        }
      }
    }
  }

  .nav-right {
    display: flex;
    align-items: center;
    gap: 20px;

    .write-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 9px 20px;
      border-radius: 10px;
      font-weight: 500;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      border: none;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      box-shadow: 0 2px 12px rgba(43, 88, 118, 0.2);
      color: #fff;

      .el-icon {
        font-size: 1.15em;
        transition: transform 0.3s ease;
      }

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 16px rgba(43, 88, 118, 0.25);

        .el-icon {
          transform: rotate(-12deg);
        }
      }
    }

    .user-dropdown {
      .user-info {
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 4px;
        border-radius: 50%;
        transition: all 0.3s ease;

        &:hover {
          background: rgba(0, 0, 0, 0.04);
          transform: translateY(-1px);
        }

        .el-avatar {
          border: 2px solid #fff;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
        }
      }
    }
  }
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  font-size: 0.95em;
  
  .el-icon {
    font-size: 1.15em;
    margin-right: 2px;
  }

  &:hover {
    background: rgba(43, 88, 118, 0.08);
    color: #2B5876;
  }
}

@media (max-width: 768px) {
  .navbar {
    .navbar-container {
      padding: 0 16px;
    }

    .nav-links {
      .nav-item {
        padding: 8px 12px;
        font-size: 0.9em;

        .nav-icon {
          font-size: 1.1em;
        }
      }
    }

    .nav-right {
      gap: 12px;

      .write-btn {
        padding: 8px 16px;
        font-size: 0.9em;
      }
    }
  }
}
</style> 