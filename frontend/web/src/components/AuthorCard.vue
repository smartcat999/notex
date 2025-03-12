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
  padding: 0;
  margin: 16px 0;
  background: transparent;
  border: none;
  box-shadow: none;
  transition: all 0.3s ease;
  display: flex;
  justify-content: center;
}

.author-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background: rgba(43, 88, 118, 0.03);
  border-radius: 24px;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(43, 88, 118, 0.06);
    transform: translateY(-1px);
  }
}

.author-avatar {
  flex: 0 0 40px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  transition: all 0.3s ease;
  display: block;
  box-shadow: 0 2px 8px rgba(43, 88, 118, 0.12);
  border: 2px solid #fff;

  &:hover {
    transform: scale(1.05) rotate(5deg);
    box-shadow: 0 4px 12px rgba(43, 88, 118, 0.2);
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
    font-size: 1.4em;
    font-weight: 600;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.1);
  }
}

.author-info {
  flex: 0 1 auto;
  min-width: 0;
  max-width: 150px;
  display: flex;
  flex-direction: column;
  justify-content: center;

  .author-name {
    font-size: 0.95em;
    font-weight: 600;
    color: #2c3e50;
    text-decoration: none;
    margin: 0;
    display: block;
    transition: all 0.3s ease;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    &:hover {
      color: #2B5876;
    }
  }

  .author-bio {
    display: none;
  }
}
</style> 