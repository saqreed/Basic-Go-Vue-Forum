import { defineStore } from 'pinia'
import axios from 'axios'

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
        const response = await axios.get('http://localhost:8081/api/posts')
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
        const response = await axios.get(`http://localhost:8081/api/posts/${id}`)
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
        const token = localStorage.getItem('token')
        const response = await axios.post('http://localhost:8081/api/posts', postData, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        
        // Create a new array with the new post at the beginning
        this.posts = [response.data, ...(this.posts || [])]
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updatePost(id, postData) {
      this.loading = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.put(`http://localhost:8081/api/posts/${id}`, postData, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        
        // Create a new array with the updated post
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

    async deletePost(id) {
      this.loading = true
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`http://localhost:8081/api/posts/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        
        // Create a new array without the deleted post
        this.posts = this.posts.filter(post => post.id !== id)
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    }
  }
}) 