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
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  background: rgba(255, 255, 255, 0.98);
  transition: all 0.3s ease;
  border: none;
  overflow: hidden;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 1);
  }

  .card-header {
    padding: 16px 20px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    background: linear-gradient(to right, rgba(43, 88, 118, 0.05), rgba(43, 88, 118, 0.02));

    span {
      font-size: 1.05em;
      font-weight: 600;
      background: linear-gradient(135deg, #2B5876, #4E4376);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: 0.02em;
    }
  }

  .archive-list {
    padding: 12px;

    .archive-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 6px 10px;
      text-decoration: none;
      border-radius: 8px;
      margin-bottom: 4px;
      transition: all 0.3s ease;
      background: rgba(255, 255, 255, 0.5);
      gap: 8px;

      &:hover {
        background: rgba(43, 88, 118, 0.06);
        transform: translateX(4px);

        .archive-date {
          color: #2B5876;
        }

        .archive-count {
          background: rgba(43, 88, 118, 0.12);
        }
      }

      .archive-date {
        font-size: 0.85em;
        color: #2c3e50;
        font-weight: 500;
        transition: color 0.3s ease;
        white-space: nowrap;
      }

      .archive-count {
        font-size: 0.8em;
        color: #2B5876;
        background: rgba(43, 88, 118, 0.08);
        padding: 1px 6px;
        border-radius: 10px;
        transition: all 0.3s ease;
        white-space: nowrap;
        flex-shrink: 0;
      }
    }
  }
}
</style> 