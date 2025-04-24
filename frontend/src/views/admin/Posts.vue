<template>
  <div class="admin-posts">
    <h1>Post Management</h1>
    
    <div v-if="loading" class="loading">
      Loading posts...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else class="posts-list">
      <div v-for="post in posts" :key="post.id" class="post-card">
        <div class="post-info">
          <h3>{{ post.title }}</h3>
          <p class="post-content">{{ post.content }}</p>
          <div class="post-meta">
            <p>Author: {{ post.author.username }}</p>
            <p>Created: {{ formatDate(post.created_at) }}</p>
          </div>
        </div>
        
        <div class="post-actions">
          <button 
            @click="deletePost(post.id)" 
            class="delete-button"
          >
            Delete Post
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const posts = ref([])
const loading = ref(true)
const error = ref(null)

const fetchPosts = async () => {
  try {
    loading.value = true
    const response = await fetch('http://localhost:8081/api/admin/posts', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch posts')
    }
    
    posts.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const deletePost = async (postId) => {
  if (!confirm('Are you sure you want to delete this post?')) {
    return
  }
  
  try {
    const response = await fetch(`http://localhost:8081/api/admin/posts/${postId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to delete post')
    }
    
    // Remove post from local state
    posts.value = posts.value.filter(post => post.id !== postId)
  } catch (err) {
    error.value = err.message
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(fetchPosts)
</script>

<style scoped>
.admin-posts {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.posts-list {
  display: grid;
  gap: 1.5rem;
  margin-top: 2rem;
}

.post-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.post-info {
  flex: 1;
}

.post-info h3 {
  margin: 0 0 0.5rem;
  color: #333;
}

.post-content {
  margin: 0.5rem 0;
  color: #666;
  max-width: 600px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.post-meta {
  margin-top: 1rem;
  color: #666;
  font-size: 0.9rem;
}

.post-meta p {
  margin: 0.25rem 0;
}

.post-actions {
  margin-left: 1rem;
}

.delete-button {
  padding: 0.5rem 1rem;
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.error {
  color: #f44336;
}
</style> 