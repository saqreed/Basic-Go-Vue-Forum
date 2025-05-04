import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './auth'

export const usePostsStore = defineStore('posts', {
  state: () => ({
    posts: [],
    currentPost: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchPosts() {
      this.loading = true
      try {
        const authStore = useAuthStore()
        const response = await axios.get('http://localhost:8081/api/posts', {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        this.posts = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.error = error.message
        this.posts = []
      } finally {
        this.loading = false
      }
    },

    async fetchPost(id) {
      this.loading = true
      try {
        const authStore = useAuthStore()
        const response = await axios.get(`http://localhost:8081/api/posts/${id}`, {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        this.currentPost = response.data
      } catch (error) {
        this.error = error.message
        this.currentPost = null
      } finally {
        this.loading = false
      }
    },

    async createPost(postData) {
      this.loading = true
      try {
        const authStore = useAuthStore()
        console.log('Creating post with token:', authStore.token ? 'exists' : 'not found')
        const response = await axios.post('http://localhost:8081/api/posts', postData, {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        
        this.posts = [response.data, ...(this.posts || [])]
        return response.data
      } catch (error) {
        console.error('Create post error:', error)
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updatePost(id, postData) {
      this.loading = true
      try {
        const authStore = useAuthStore()
        const response = await axios.put(`http://localhost:8081/api/posts/${id}`, postData, {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        
        this.posts = this.posts.map(post => 
          post.id === id ? response.data : post
        )
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async deletePost(postId) {
      try {
        await axios.delete(`http://localhost:8081/api/posts/${postId}`)
        // Удаляем пост из локального состояния
        this.posts = this.posts.filter(post => post.id !== postId)
        if (this.currentPost && this.currentPost.id === postId) {
          this.currentPost = null
        }
      } catch (error) {
        if (error.response) {
          // Сервер вернул ошибку
          throw new Error(error.response.data || 'Failed to delete post')
        } else if (error.request) {
          // Запрос был сделан, но ответ не получен
          throw new Error('No response from server')
        } else {
          // Ошибка при настройке запроса
          throw new Error('Failed to send request')
        }
      }
    }
  }
}) 