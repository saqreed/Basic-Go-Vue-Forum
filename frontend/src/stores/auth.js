import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    isAuthenticated: false
  }),

  actions: {
    async login(email, password) {
      try {
        const response = await axios.post('http://localhost:8081/api/login', {
          email,
          password
        }, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
        
        if (response.data.token) {
          this.token = response.data.token
          this.isAuthenticated = true
          localStorage.setItem('token', this.token)
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          return true
        }
        return false
      } catch (error) {
        throw error
      }
    },

    async register(username, email, password) {
      try {
        console.log('Sending registration request with data:', {
          username,
          email,
          password,
          passwordLength: password.length
        })

        const response = await axios.post('http://localhost:8081/api/register', {
          username,
          email,
          password
        }, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
        
        console.log('Registration response:', response.data)
        
        if (response.data.token) {
          this.token = response.data.token
          this.isAuthenticated = true
          localStorage.setItem('token', this.token)
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          return true
        }
        return false
      } catch (error) {
        console.error('Registration error details:', {
          message: error.message,
          response: error.response?.data,
          requestData: {
            username,
            email,
            passwordLength: password.length
          }
        })
        throw error
      }
    },

    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
    },

    // Инициализация при загрузке приложения
    initialize() {
      const token = localStorage.getItem('token')
      if (token) {
        this.token = token
        this.isAuthenticated = true
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      }
    }
  }
}) 