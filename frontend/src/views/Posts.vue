<template>
    <div class="posts-container">
      <div class="posts-header">
        <h1>Посты</h1>
        <router-link to="/posts/create" class="create-post-btn">
          Создать пост
        </router-link>
      </div>
  
      <div v-if="loading" class="loading">Загрузка постов...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="!posts || posts.length === 0" class="no-posts">
        Пока нет постов. Создайте первый!
      </div>
      <div v-else class="posts-list">
        <div v-for="post in posts" :key="post.id" class="post-card">
          <div class="post-header">
            <h2 class="post-title">
              <router-link :to="'/post/' + post.id">{{ post.title }}</router-link>
            </h2>
            <span class="post-author">Автор: {{ post.author?.username || 'Неизвестный автор' }}</span>
          </div>
          <div class="post-content">{{ post.content }}</div>
          <div class="post-footer">
            <span class="post-date">{{ formatDate(post.created_at) }}</span>
            <div class="post-stats">
              <span class="comments-count">
                {{ post.comments_count || 0 }} комментариев
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { onMounted } from 'vue'
  import { usePostsStore } from '../stores/posts'
  import { computed } from 'vue'
  
  const postsStore = usePostsStore()
  const loading = computed(() => postsStore.loading)
  const error = computed(() => postsStore.error)
  const posts = computed(() => postsStore.posts || [])
  
  onMounted(async () => {
    await loadPosts()
  })
  
  async function loadPosts() {
    try {
      await postsStore.fetchPosts()
    } catch (err) {
      console.error('Failed to load posts:', err)
    }
  }
  
  function formatDate(dateString) {
    return new Date(dateString).toLocaleString()
  }
  </script>
  
  <style scoped>
  .posts-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .posts-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  .create-post-btn {
    padding: 0.5rem 1rem;
    background-color: #4CAF50;
    color: white;
    text-decoration: none;
    border-radius: 4px;
    transition: background-color 0.3s;
  }
  
  .create-post-btn:hover {
    background-color: #45a049;
  }
  
  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .post-card {
    padding: 1.5rem;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .post-header {
    margin-bottom: 1rem;
  }
  
  .post-title {
    margin: 0;
    font-size: 1.5rem;
  }
  
  .post-title a {
    color: #333;
    text-decoration: none;
  }
  
  .post-title a:hover {
    color: #4CAF50;
  }
  
  .post-author {
    color: #666;
    font-size: 0.9rem;
  }
  
  .post-content {
    margin-bottom: 1rem;
    line-height: 1.6;
  }
  
  .post-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #666;
    font-size: 0.9rem;
  }
  
  .post-stats {
    display: flex;
    gap: 1rem;
  }
  
  .loading, .error, .no-posts {
    text-align: center;
    padding: 2rem;
    color: #666;
  }
  
  .error {
    color: #f44336;
  }
  </style>