<template>
  <div class="archives-container">
    <el-card class="archives-card">
      <template #header>
        <div class="card-header">
          <h2>文章归档</h2>
          <span class="total-count">共 {{ totalPosts }} 篇文章</span>
        </div>
      </template>

      <div class="archives-list">
        <div v-for="archive in archives" :key="archive.date" class="archive-item">
          <div class="archive-header" @click="toggleArchive(archive.date)">
            <el-icon>
              <component :is="expandedArchives[archive.date] ? 'ArrowDown' : 'ArrowRight'" />
            </el-icon>
            <span class="archive-date">{{ formatDate(archive.date) }}</span>
            <span class="archive-count">{{ archive.count }} 篇文章</span>
          </div>
          
          <el-collapse-transition>
            <div v-show="expandedArchives[archive.date]" class="archive-posts">
              <div v-if="archivePosts[archive.date]" class="posts-list">
                <div v-for="post in archivePosts[archive.date]" :key="post.id" class="post-item">
                  <router-link :to="'/posts/' + post.id" class="post-link">
                    <span class="post-date">{{ formatPostDate(post.published_at) }}</span>
                    <span class="post-title">{{ post.title }}</span>
                  </router-link>
                </div>
              </div>
              <div v-else class="loading-posts">
                <el-icon class="is-loading"><Loading /></el-icon>
                加载中...
              </div>
            </div>
          </el-collapse-transition>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getArchives, getPostsByArchive } from '@/api/posts'
import { ArrowRight, ArrowDown, Loading } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const archives = ref([])
const archivePosts = ref({})
const expandedArchives = ref({})

const totalPosts = computed(() => {
  return archives.value.reduce((total, archive) => total + archive.count, 0)
})

const formatDate = (date) => {
  const [year, month] = date.split('-')
  return `${year}年${month}月`
}

const formatPostDate = (date) => {
  return dayjs(date).format('MM-DD')
}

const toggleArchive = async (date) => {
  expandedArchives.value[date] = !expandedArchives.value[date]
  
  // 如果展开且还没有加载文章，则加载该月份的文章
  if (expandedArchives.value[date] && !archivePosts.value[date]) {
    try {
      const response = await getPostsByArchive(date)
      archivePosts.value[date] = response.items || []
    } catch (error) {
      console.error('Failed to load archive posts:', error)
      ElMessage.error('加载文章失败')
    }
  }
}

onMounted(async () => {
  try {
    const response = await getArchives()
    archives.value = response.items || []
  } catch (error) {
    console.error('Failed to load archives:', error)
    ElMessage.error('加载归档失败')
  }
})
</script>

<style lang="scss" scoped>
.archives-container {
  max-width: 800px;
  margin: 32px auto;
  padding: 0 20px;
}

.archives-card {
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  background: rgba(255, 255, 255, 0.98);
  transition: all 0.3s ease;
  border: none;
  overflow: hidden;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
    background: rgba(255, 255, 255, 1);
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    background: linear-gradient(to right, rgba(43, 88, 118, 0.05), rgba(43, 88, 118, 0.02));
    border-bottom: 1px solid rgba(0, 0, 0, 0.04);

    h2 {
      margin: 0;
      font-size: 1.4em;
      font-weight: 600;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: 0.02em;
    }

    .total-count {
      color: #4a5568;
      font-size: 0.85em;
      background: rgba(43, 88, 118, 0.08);
      padding: 4px 12px;
      border-radius: 6px;
      transition: all 0.3s ease;

      &:hover {
        background: rgba(43, 88, 118, 0.12);
      }
    }
  }
}

.archive-item {
  margin: 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);

  &:last-child {
    border-bottom: none;
  }

  .archive-header {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 16px 24px;
    transition: all 0.3s ease;

    &:hover {
      background: rgba(43, 88, 118, 0.04);

      .archive-date {
        color: #2B5876;
      }
    }

    .el-icon {
      margin-right: 12px;
      font-size: 1.1em;
      color: var(--el-text-color-secondary);
      transition: transform 0.3s ease;
    }

    .archive-date {
      font-size: 1.05em;
      font-weight: 500;
      color: #2c3e50;
      margin-right: 12px;
      transition: color 0.3s ease;
    }

    .archive-count {
      font-size: 0.82em;
      color: var(--el-text-color-secondary);
      background: rgba(43, 88, 118, 0.08);
      padding: 2px 8px;
      border-radius: 10px;
    }
  }

  .archive-posts {
    padding: 0 24px 16px 48px;

    .post-item {
      margin: 12px 0;

      .post-link {
        display: flex;
        align-items: center;
        text-decoration: none;
        padding: 8px 12px;
        border-radius: 8px;
        transition: all 0.3s ease;
        background: rgba(255, 255, 255, 0.5);

        &:hover {
          background: rgba(43, 88, 118, 0.06);
          transform: translateX(4px);
          
          .post-title {
            color: #2B5876;
          }
        }

        .post-date {
          font-size: 0.85em;
          color: #4a5568;
          margin-right: 16px;
          min-width: 60px;
          font-family: Monaco, monospace;
        }

        .post-title {
          color: #2c3e50;
          font-size: 0.95em;
          transition: color 0.3s ease;
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }

    .loading-posts {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 20px;
      color: var(--el-text-color-secondary);
      font-size: 0.9em;
      
      .el-icon {
        margin-right: 8px;
        animation: rotating 2s linear infinite;
      }
    }
  }
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style> 
