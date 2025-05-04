import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './auth'

export const useCommentsStore = defineStore('comments', {
  state: () => ({
    comments: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchComments(postId) {
      this.loading = true
      this.error = null
      try {
        const authStore = useAuthStore()
        const response = await axios.get(`http://localhost:8081/api/posts/${postId}/comments`, {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        this.comments = response.data || []
      } catch (error) {
        console.error('Failed to fetch comments:', error)
        this.error = error.response?.data || error.message
        this.comments = []
      } finally {
        this.loading = false
      }
    },

    async addComment(postId, content) {
      this.loading = true
      this.error = null
      try {
        const authStore = useAuthStore()
        if (!authStore.token) {
          throw new Error('Not authenticated')
        }
        
        const response = await axios.post(
          `http://localhost:8081/api/posts/${postId}/comments`,
          { content },
          {
            headers: {
              'Authorization': `Bearer ${authStore.token}`,
              'Content-Type': 'application/json'
            }
          }
        )
        
        if (response.data) {
          this.comments = [response.data, ...(this.comments || [])]
        }
        return response.data
      } catch (error) {
        console.error('Failed to add comment:', error)
        this.error = error.response?.data || error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateComment(commentId, content) {
      this.loading = true
      this.error = null
      try {
        const authStore = useAuthStore()
        const response = await axios.put(
          `http://localhost:8081/api/comments/${commentId}`,
          { content },
          {
            headers: {
              'Authorization': `Bearer ${authStore.token}`,
              'Content-Type': 'application/json'
            }
          }
        )
        const index = (this.comments || []).findIndex(comment => comment.id === commentId)
        if (index !== -1) {
          this.comments[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Failed to update comment:', error)
        this.error = error.response?.data || error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteComment(commentId) {
      this.loading = true
      this.error = null
      try {
        const authStore = useAuthStore()
        await axios.delete(`http://localhost:8081/api/comments/${commentId}`, {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        this.comments = (this.comments || []).filter(comment => comment.id !== commentId)
      } catch (error) {
        console.error('Failed to delete comment:', error)
        this.error = error.response?.data || error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    clearComments() {
      this.comments = []
      this.error = null
      this.loading = false
    }
  }
})