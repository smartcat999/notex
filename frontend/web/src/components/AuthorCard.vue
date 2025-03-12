<template>
  <div class="author-card">
    <div class="author-content">
      <router-link :to="`/users/${author.id}`" class="author-avatar">
        <template v-if="author.avatar">
          <img :src="author.avatar" :alt="author.username">
        </template>
        <template v-else>
          <div class="default-avatar">
            {{ author.username.charAt(0).toUpperCase() }}
          </div>
        </template>
      </router-link>
      <div class="author-info">
        <router-link :to="`/users/${author.id}`" class="author-name">
          {{ author.username }}
        </router-link>
        <p class="author-bio" v-if="author.bio">{{ author.bio }}</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AuthorCard',
  props: {
    author: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>
.author-card {
  width: 100%;
  padding: 16px 0;
  margin: 16px 0;
  border-top: 1px solid rgba(43, 88, 118, 0.08);
  border-bottom: 1px solid rgba(43, 88, 118, 0.08);
  display: flex;
  justify-content: center;
}

.author-content {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;
  max-width: 320px;
  padding: 0;
  justify-content: center;
}

.author-avatar {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  transition: transform 0.3s ease;
  display: block;

  &:hover {
    transform: scale(1.05);
  }

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .default-avatar {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #2B5876, #4E4376);
    color: #fff;
    font-size: 1.2em;
    font-weight: 500;
  }
}

.author-info {
  flex: 0 1 auto;
  min-width: 0;
  overflow: hidden;

  .author-name {
    font-size: 0.95em;
    font-weight: 600;
    color: #2c3e50;
    text-decoration: none;
    margin-bottom: 4px;
    display: block;
    transition: color 0.3s ease;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: left;

    &:hover {
      color: #2B5876;
    }
  }

  .author-bio {
    color: #606266;
    font-size: 0.85em;
    line-height: 1.5;
    margin: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    text-align: left;
  }
}
</style> 