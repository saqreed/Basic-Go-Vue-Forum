import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    isAuthenticated: false
  }),

  actions: {
    async fetchUserData() {
      try {
        const response = await axios.get('http://localhost:8081/api/profile', {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        })
        this.user = response.data
        this.isAuthenticated = true
      } catch (error) {
        console.error('Failed to fetch user data:', error)
        this.logout()
      }
    },

    async login(email, password) {
      try {
        console.log('Attempting login with:', { email })
        const response = await axios.post('http://localhost:8081/api/login', {
          email,
          password
        }, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
        
        if (response.data.token) {
          console.log('Login successful, received token:', response.data.token)
          this.token = response.data.token
          this.isAuthenticated = true
          localStorage.setItem('token', this.token)
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          console.log('Token stored and axios headers updated')
          
          // Fetch user data after successful login
          await this.fetchUserData()
          return true
        }
        return false
      } catch (error) {
        console.error('Login error:', error)
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
          console.log('Registration successful, received token:', response.data.token)
          this.token = response.data.token
          this.isAuthenticated = true
          localStorage.setItem('token', this.token)
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          console.log('Token stored and axios headers updated')
          
          // Fetch user data after successful registration
          await this.fetchUserData()
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
      console.log('Logging out, clearing token and auth state')
      this.user = null
      this.token = null
      this.isAuthenticated = false
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
      console.log('Token cleared and axios headers removed')
    },

    // Инициализация при загрузке приложения
    async initialize() {
      console.log('Initializing auth store')
      const token = localStorage.getItem('token')
      console.log('Token from localStorage:', token ? 'exists' : 'not found')
      if (token) {
        console.log('Token found, setting up auth state')
        this.token = token
        this.isAuthenticated = true
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        console.log('Auth store initialized with token:', token)
        
        // Fetch user data on initialization
        await this.fetchUserData()
      } else {
        console.log('No token found, auth store not initialized')
      }
    }
  }
}) 