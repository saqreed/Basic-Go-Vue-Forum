import { defineStore } from 'pinia'
import axios from 'axios'

export const useCommentsStore = defineStore('comments', {
  state: () => ({
    comments: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchComments(postId) {
      this.loading = true
      try {
        const response = await axios.get(`http://localhost:8081/api/posts/${postId}/comments`)
        this.comments = response.data
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async addComment(postId, content) {
      this.loading = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.post(
          `http://localhost:8081/api/posts/${postId}/comments`,
          { content },
          {
            headers: {
              Authorization: `Bearer ${token}`
            }
          }
        )
        this.comments.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateComment(commentId, content) {
      this.loading = true
      try {
        const token = localStorage.getItem('token')
        const response = await axios.put(
          `http://localhost:8081/api/comments/${commentId}`,
          { content },
          {
            headers: {
              Authorization: `Bearer ${token}`
            }
          }
        )
        const index = this.comments.findIndex(comment => comment.id === commentId)
        if (index !== -1) {
          this.comments[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteComment(commentId) {
      this.loading = true
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`http://localhost:8081/api/comments/${commentId}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        this.comments = this.comments.filter(comment => comment.id !== commentId)
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    clearComments() {
      this.comments = []
      this.error = null
    }
  }
})