<template>
    <div class="post-form-container">
      <div class="post-form">
        <h1>{{ isEditing ? 'Редактировать пост' : 'Создать новый пост' }}</h1>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="title">Заголовок</label>
            <input
              type="text"
              id="title"
              v-model="title"
              required
              placeholder="Введите заголовок поста"
            />
          </div>
          <div class="form-group">
            <label for="content">Содержание</label>
            <textarea
              id="content"
              v-model="content"
              required
              placeholder="Введите содержание поста"
              rows="10"
            ></textarea>
          </div>
          <button type="submit" class="submit-btn">
            {{ isEditing ? 'Сохранить изменения' : 'Создать пост' }}
          </button>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  import axios from 'axios'
  import { useRouter, useRoute } from 'vue-router'
  import { ref, onMounted } from 'vue'
  
  export default {
    name: 'CreatePost',
    setup() {
      const router = useRouter()
      const route = useRoute()
      const title = ref('')
      const content = ref('')
      const isEditing = ref(false)
      const postId = ref(null)
  
      onMounted(async () => {
        if (route.params.id) {
          isEditing.value = true
          postId.value = route.params.id
          try {
            const response = await axios.get(`http://localhost:8081/api/posts/${postId.value}`)
            title.value = response.data.title
            content.value = response.data.content
          } catch (error) {
            console.error('Error fetching post:', error)
            router.push('/posts')
          }
        }
      })
  
      const handleSubmit = async () => {
        try {
          const postData = {
            title: title.value,
            content: content.value
          }
  
          if (isEditing.value) {
            await axios.put(`http://localhost:8081/api/posts/${postId.value}`, postData)
          } else {
            await axios.post('http://localhost:8081/api/posts', postData)
          }
  
          router.push('/posts')
        } catch (error) {
          console.error('Error saving post:', error)
        }
      }
  
      return {
        title,
        content,
        handleSubmit,
        isEditing
      }
    }
  }
  </script>
  
  <style scoped>
  .post-form-container {
    display: flex;
    justify-content: center;
    min-height: 80vh;
    padding: 1.5rem;
  }
  
  .post-form {
    width: 100%;
    max-width: 800px;
    padding: 1.5rem;
    background: var(--card-bg);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
  }
  
  h1 {
    margin: 0 0 1.5rem;
    color: var(--text-color);
    font-size: 1.5rem;
  }
  
  .form-group {
    margin-bottom: 1.25rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    color: var(--text-color);
    font-weight: 500;
  }
  
  input, textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    font-size: 1rem;
    background: var(--card-bg);
    color: var(--text-color);
    transition: border-color 0.2s;
  }
  
  input:focus, textarea:focus {
    outline: none;
    border-color: var(--primary-color);
  }
  
  textarea {
    resize: vertical;
    min-height: 200px;
  }
  
  .submit-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    background-color: var(--primary-color);
    color: white;
  }
  
  .submit-btn:hover {
    background-color: #357abd;
    transform: translateY(-1px);
  }
  
  @media (max-width: 768px) {
    .post-form-container {
      padding: 1rem;
    }
  
    .post-form {
      padding: 1rem;
    }
  
    h1 {
      font-size: 1.25rem;
    }
  }
  </style>