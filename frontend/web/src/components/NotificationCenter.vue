<template>
  <el-popover
    placement="bottom-end"
    :width="400"
    trigger="click"
    popper-class="notification-popover"
  >
    <template #reference>
      <el-badge :value="unreadCount || ''" :hidden="!unreadCount" class="notification-badge">
        <el-button class="notification-btn" :class="{ 'has-unread': unreadCount > 0 }">
          <el-icon><Bell /></el-icon>
        </el-button>
      </el-badge>
    </template>

    <div class="notification-header">
      <h3>通知中心</h3>
      <div class="notification-actions">
        <el-radio-group v-model="currentType" size="small" @change="handleTypeChange">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="post_comment">评论</el-radio-button>
          <el-radio-button label="comment_reply">回复</el-radio-button>
        </el-radio-group>
        <el-button 
          v-if="notifications.length > 0"
          link 
          type="primary" 
          @click="handleMarkAllAsRead"
        >
          全部已读
        </el-button>
      </div>
    </div>

    <div class="notification-list">
      <template v-if="notifications.length > 0">
        <div
          v-for="item in notifications"
          :key="item.id"
          class="notification-item"
          :class="{ 'unread': !item.is_read }"
          @click="handleNotificationClick(item)"
        >
          <el-avatar :size="40" :src="item.sender_avatar">
            {{ item.sender_username?.charAt(0) }}
          </el-avatar>
          <div class="notification-content">
            <div class="notification-message">
              <span class="sender">{{ item.sender_username }}</span>
              {{ item.content }}
            </div>
            <div class="notification-time">
              {{ formatTime(item.created_at) }}
            </div>
          </div>
        </div>
      </template>
      <div v-else class="empty-state">
        暂无通知
      </div>
    </div>
  </el-popover>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getNotifications, markAsRead, markAllAsRead, getUnreadCount } from '@/api/notification'
import { formatTime } from '@/utils/time'

const router = useRouter()
const notifications = ref([])
const unreadCount = ref(0)
const currentType = ref('')

const fetchNotifications = async () => {
  try {
    const params = { 
      page: 1, 
      page_size: 10,
      type: currentType.value || undefined
    }
    const { items } = await getNotifications(params)
    notifications.value = items
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  }
}

const fetchUnreadCount = async () => {
  try {
    const { count } = await getUnreadCount()
    unreadCount.value = count
  } catch (error) {
    console.error('Failed to fetch unread count:', error)
  }
}

const handleTypeChange = () => {
  fetchNotifications()
}

const handleNotificationClick = async (notification) => {
  if (!notification.is_read) {
    try {
      await markAsRead(notification.id)
      notification.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (error) {
      console.error('Failed to mark notification as read:', error)
    }
  }

  // 根据通知类型跳转到相应页面
  if (notification.type === 'post_comment' || notification.type === 'comment_reply') {
    router.push(`/posts/${notification.post_id}#comment-${notification.comment_id}`)
  }
}

const handleMarkAllAsRead = async () => {
  try {
    await markAllAsRead()
    notifications.value.forEach(notification => {
      notification.is_read = true
    })
    unreadCount.value = 0
    ElMessage.success('已将全部通知标记为已读')
  } catch (error) {
    console.error('Failed to mark all as read:', error)
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  fetchNotifications()
  fetchUnreadCount()
})
</script>

<style lang="scss" scoped>
.notification-badge {
  :deep(.el-badge__content) {
    z-index: 1;
  }
}

.notification-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  padding: 0;
  border: none;
  background: transparent;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(0, 0, 0, 0.04);
  }

  &.has-unread {
    color: var(--el-color-primary);
  }

  .el-icon {
    font-size: 20px;
  }
}

.notification-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);

  h3 {
    margin: 0 0 12px 0;
    font-size: 16px;
    font-weight: 600;
  }

  .notification-actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
  }
}

.notification-list {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background: var(--el-fill-color-light);
  }

  &.unread {
    background: var(--el-color-primary-light-9);

    &:hover {
      background: var(--el-color-primary-light-8);
    }

    .notification-message {
      font-weight: 500;
    }
  }
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-message {
  font-size: 14px;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
  line-height: 1.4;

  .sender {
    font-weight: 600;
    margin-right: 4px;
  }
}

.notification-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.empty-state {
  padding: 32px;
  text-align: center;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

:deep(.notification-popover) {
  padding: 0;
  border-radius: 12px;
}
</style> 