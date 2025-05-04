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
              <router-link :to="'/posts/' + post.id">{{ post.title }}</router-link>
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
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
    padding: 1.5rem;
  }
  
  .posts-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }
  
  h1 {
    margin: 0;
    color: var(--text-color);
    font-size: 1.75rem;
  }
  
  .create-post-btn {
    padding: 0.75rem 1.5rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    transition: all 0.2s;
  }
  
  .create-post-btn:hover {
    background-color: #357abd;
    transform: translateY(-1px);
  }
  
  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .post-card {
    background: var(--card-bg);
    border-radius: 8px;
    padding: 1.25rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    transition: all 0.2s;
  }
  
  .post-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
  
  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }
  
  .post-title {
    margin: 0;
    color: var(--text-color);
    font-size: 1.25rem;
    font-weight: 600;
  }
  
  .post-meta {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }
  
  .post-content {
    color: var(--text-color);
    margin-bottom: 1rem;
    line-height: 1.5;
    max-height: 4.5em;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
  }
  
  .post-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border-color);
  }
  
  .post-stats {
    display: flex;
    gap: 1rem;
    color: var(--text-secondary);
    font-size: 0.875rem;
  }
  
  .stat {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .read-more {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s;
  }
  
  .read-more:hover {
    color: #357abd;
  }
  
  .loading, .error {
    text-align: center;
    padding: 2rem;
    color: var(--text-color);
  }
  
  .error {
    color: #dc3545;
  }
  
  .no-posts {
    text-align: center;
    padding: 1.5rem;
    color: var(--text-color);
    font-size: 1rem;
  }
  
  @media (max-width: 768px) {
    .posts-container {
      padding: 1rem;
    }
  
    .posts-header {
      flex-direction: column;
      gap: 1rem;
      align-items: flex-start;
    }
  
    h1 {
      font-size: 1.5rem;
    }
  
    .create-post-btn {
      width: 100%;
      text-align: center;
    }
  
    .post-card {
      padding: 1rem;
    }
  
    .post-title {
      font-size: 1.125rem;
    }
  
    .post-meta {
      font-size: 0.75rem;
    }
  
    .post-content {
      font-size: 0.9375rem;
    }
  
    .post-footer {
      flex-direction: column;
      gap: 0.75rem;
      align-items: flex-start;
    }
  
    .post-stats {
      width: 100%;
      justify-content: space-between;
    }
  }
  </style>