<template>
  <div class="admin-comments">
    <h1>Comment Management</h1>
    
    <div v-if="loading" class="loading">
      Loading comments...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else class="comments-list">
      <div v-for="comment in comments" :key="comment.id" class="comment-card">
        <div class="comment-info">
          <p class="comment-content">{{ comment.content }}</p>
          <div class="comment-meta">
            <p>Author: {{ comment.author.username }}</p>
            <p>Post: {{ comment.post.title }}</p>
            <p>Created: {{ formatDate(comment.created_at) }}</p>
          </div>
        </div>
        
        <div class="comment-actions">
          <button 
            @click="deleteComment(comment.id)" 
            class="delete-button"
          >
            Delete Comment
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
const comments = ref([])
const loading = ref(true)
const error = ref(null)

const fetchComments = async () => {
  try {
    loading.value = true
    const response = await fetch('http://localhost:8081/api/admin/comments', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch comments')
    }
    
    comments.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const deleteComment = async (commentId) => {
  if (!confirm('Are you sure you want to delete this comment?')) {
    return
  }
  
  try {
    const response = await fetch(`http://localhost:8081/api/admin/comments/${commentId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to delete comment')
    }
    
    // Remove comment from local state
    comments.value = comments.value.filter(comment => comment.id !== commentId)
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

onMounted(fetchComments)
</script>

<style scoped>
.admin-comments {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.comments-list {
  display: grid;
  gap: 1.5rem;
  margin-top: 2rem;
}

.comment-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.comment-info {
  flex: 1;
}

.comment-content {
  margin: 0 0 0.5rem;
  color: #333;
  font-size: 1.1rem;
}

.comment-meta {
  margin-top: 1rem;
  color: #666;
  font-size: 0.9rem;
}

.comment-meta p {
  margin: 0.25rem 0;
}

.comment-actions {
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