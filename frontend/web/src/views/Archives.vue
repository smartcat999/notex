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
  margin: 20px auto;
  padding: 0 20px;
}

.archives-card {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    h2 {
      margin: 0;
      font-size: 1.5rem;
      color: var(--el-text-color-primary);
    }

    .total-count {
      color: var(--el-text-color-secondary);
      font-size: 0.9rem;
    }
  }
}

.archive-item {
  margin-bottom: 16px;

  .archive-header {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 8px;
    border-radius: 4px;
    transition: background-color 0.3s;

    &:hover {
      background-color: var(--el-fill-color-light);
    }

    .el-icon {
      margin-right: 8px;
      font-size: 1.2em;
      color: var(--el-text-color-secondary);
    }

    .archive-date {
      font-size: 1.1rem;
      font-weight: 500;
      color: var(--el-text-color-primary);
      margin-right: 12px;
    }

    .archive-count {
      font-size: 0.9rem;
      color: var(--el-text-color-secondary);
    }
  }

  .archive-posts {
    padding-left: 32px;
    margin-top: 8px;

    .post-item {
      margin: 8px 0;

      .post-link {
        display: flex;
        align-items: center;
        text-decoration: none;
        padding: 6px;
        border-radius: 4px;
        transition: all 0.3s;

        &:hover {
          background-color: var(--el-fill-color-lighter);
          
          .post-title {
            color: var(--el-color-primary);
          }
        }

        .post-date {
          font-size: 0.9rem;
          color: var(--el-text-color-secondary);
          margin-right: 16px;
          min-width: 60px;
        }

        .post-title {
          color: var(--el-text-color-primary);
          transition: color 0.3s;
        }
      }
    }

    .loading-posts {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 20px;
      color: var(--el-text-color-secondary);
      
      .el-icon {
        margin-right: 8px;
      }
    }
  }
}
</style> 