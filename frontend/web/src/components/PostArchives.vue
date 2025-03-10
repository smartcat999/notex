<template>
  <el-card class="archive-card">
    <template #header>
      <div class="card-header">
        <span>文章归档</span>
      </div>
    </template>
    <div class="archive-list">
      <el-skeleton v-if="loading" :rows="5" animated />
      <template v-else>
        <router-link
          v-for="archive in archives"
          :key="archive.date"
          :to="'/archives'"
          class="archive-item"
        >
          <span class="archive-date">{{ formatDate(archive.date) }}</span>
          <span class="archive-count">({{ archive.count }})</span>
        </router-link>
      </template>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getArchives } from '@/api/posts'
import { ElMessage } from 'element-plus'

const archives = ref([])
const loading = ref(true)

const formatDate = (date) => {
  const [year, month] = date.split('-')
  return `${year}年${month}月`
}

onMounted(async () => {
  try {
    const response = await getArchives()
    archives.value = response.items || []
  } catch (error) {
    console.error('Failed to load archives:', error)
    ElMessage.error('加载归档失败')
  } finally {
    loading.value = false
  }
})
</script>

<style lang="scss" scoped>
.archive-card {
  margin-bottom: 20px;

  .card-header {
    font-weight: 500;
  }

  .archive-list {
    .archive-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 0;
      text-decoration: none;
      color: var(--el-text-color-primary);
      transition: all 0.3s;

      &:not(:last-child) {
        border-bottom: 1px solid var(--el-border-color-lighter);
      }

      &:hover {
        color: var(--el-color-primary);
        padding-left: 8px;
      }

      .archive-date {
        font-size: 0.9rem;
      }

      .archive-count {
        font-size: 0.85rem;
        color: var(--el-text-color-secondary);
      }
    }
  }
}
</style> 