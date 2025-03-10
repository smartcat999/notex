<template>
  <el-container class="layout-container">
    <el-header>
      <nav class="nav-container">
        <div class="nav-left">
          <router-link to="/" class="logo">
            Notex
          </router-link>
        </div>
        <div class="nav-center">
          <el-menu
            mode="horizontal"
            :router="true"
            :default-active="activeMenu"
          >
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/posts">文章</el-menu-item>
            <el-menu-item index="/categories">分类</el-menu-item>
            <el-menu-item index="/tags">标签</el-menu-item>
          </el-menu>
        </div>
        <div class="nav-right">
          <template v-if="userStore.isAuthenticated">
            <router-link to="/posts/new">
              <el-button
                type="primary"
                class="write-btn"
              >
                <el-icon><Plus /></el-icon>
                写文章
              </el-button>
            </router-link>
            <el-dropdown @command="handleCommand">
              <span class="user-dropdown">
                {{ userStore.user?.username }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="settings">设置</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <router-link to="/login" class="login-btn">
              登录
            </router-link>
            <router-link to="/register" class="register-btn">
              注册
            </router-link>
          </template>
        </div>
      </nav>
    </el-header>

    <el-main>
      <router-view />
    </el-main>

    <el-footer>
      <div class="footer-content">
        <p>&copy; {{ new Date().getFullYear() }} Notex. All rights reserved.</p>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ArrowDown, Plus } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      userStore.logoutUser()
      router.push('/login')
      break
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  min-height: 100vh;
}

.nav-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 20px;
}

.nav-left {
  .logo {
    font-size: 24px;
    font-weight: bold;
    color: var(--el-color-primary);
    text-decoration: none;
  }
}

.nav-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 16px;

  a {
    text-decoration: none;
  }

  .write-btn {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .user-dropdown {
    display: flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    color: var(--el-text-color-primary);
  }

  .login-btn,
  .register-btn {
    padding: 8px 16px;
    border-radius: 4px;
    text-decoration: none;
    transition: all 0.3s;
  }

  .login-btn {
    color: var(--el-color-primary);
    border: 1px solid var(--el-color-primary);

    &:hover {
      background-color: var(--el-color-primary-light-9);
    }
  }

  .register-btn {
    background-color: var(--el-color-primary);
    color: white;

    &:hover {
      background-color: var(--el-color-primary-dark-2);
    }
  }
}

.footer-content {
  text-align: center;
  color: var(--el-text-color-secondary);
  padding: 20px 0;
}
</style> 