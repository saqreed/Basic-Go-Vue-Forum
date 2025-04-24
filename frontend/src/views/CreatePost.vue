<template>
    <div class="post-form-container">
      <div class="post-form">
        <h2>{{ isEdit ? 'Edit Post' : 'Create Post' }}</h2>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="title">Title</label>
            <input
              type="text"
              id="title"
              v-model="title"
              required
              placeholder="Enter post title"
            />
          </div>
          <div class="form-group">
            <label for="content">Content</label>
            <textarea
              id="content"
              v-model="content"
              required
              placeholder="Write your post content..."
              rows="10"
            ></textarea>
          </div>
          <div class="form-actions">
            <button type="submit" :disabled="loading">
              {{ loading ? 'Saving...' : (isEdit ? 'Update' : 'Create') }}
            </button>
            <button type="button" @click="cancel" class="cancel-button">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { usePostsStore } from '../stores/posts'
  import { useAuthStore } from '../stores/auth'
  
  const route = useRoute()
  const router = useRouter()
  const postsStore = usePostsStore()
  const authStore = useAuthStore()
  
  const title = ref('')
  const content = ref('')
  const loading = ref(false)
  const isEdit = computed(() => route.params.id !== undefined)
  
  onMounted(async () => {
    if (isEdit.value) {
      await postsStore.fetchPost(route.params.id)
      if (postsStore.currentPost) {
        title.value = postsStore.currentPost.title
        content.value = postsStore.currentPost.content
      }
    }
  })
  
  const handleSubmit = async () => {
    loading.value = true
    try {
      const postData = {
        title: title.value,
        content: content.value
      }
  
      if (isEdit.value) {
        await postsStore.updatePost(route.params.id, postData)
      } else {
        await postsStore.createPost(postData)
      }
      router.push('/posts')
    } catch (error) {
      console.error('Failed to save post:', error)
    } finally {
      loading.value = false
    }
  }
  
  const cancel = () => {
    router.push('/posts')
  }
  </script>
  
  <style scoped>
  .post-form-container {
    display: flex;
    justify-content: center;
    min-height: 80vh;
  }
  
  .post-form {
    width: 100%;
    max-width: 800px;
    padding: 2rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  h2 {
    margin: 0 0 1.5rem;
    color: #333;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    color: #666;
  }
  
  input, textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }
  
  textarea {
    resize: vertical;
    min-height: 200px;
  }
  
  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }
  
  button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  button[type="submit"] {
    background-color: #4CAF50;
    color: white;
  }
  
  button[type="submit"]:hover {
    background-color: #45a049;
  }
  
  button[type="submit"]:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }
  
  .cancel-button {
    background-color: #f44336;
    color: white;
  }
  
  .cancel-button:hover {
    background-color: #d32f2f;
  }
  </style>