<template>
  <div class="register-container">
    <div class="register-form">
      <h2>Register</h2>
      <form @submit.prevent="handleRegister">
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <div class="form-group">
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            required
            placeholder="Enter your username"
          />
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input
            type="email"
            id="email"
            v-model="email"
            required
            placeholder="Enter your email"
          />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            required
            placeholder="Enter your password (min 6 characters)"
          />
        </div>
        <div class="form-group">
          <label for="confirmPassword">Confirm Password</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            required
            placeholder="Confirm your password"
          />
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? 'Registering...' : 'Register' }}
        </button>
        <p class="login-link">
          Already have an account? <router-link to="/login">Login</router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const router = useRouter()
const authStore = useAuthStore()

const validatePassword = (password) => {
  if (password.length < 6) {
    return 'Password must be at least 6 characters long'
  }
  return ''
}

const handleRegister = async () => {
  error.value = ''
  
  // Валидация пароля
  const passwordError = validatePassword(password.value)
  if (passwordError) {
    error.value = passwordError
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  // Логируем данные перед отправкой
  console.log('Registration data:', {
    username: username.value,
    email: email.value,
    password: password.value,
    passwordLength: password.value.length
  })

  loading.value = true
  try {
    const success = await authStore.register(username.value, email.value, password.value)
    if (success) {
      router.push('/login')
    }
  } catch (err) {
    error.value = err.response?.data || 'An error occurred during registration'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.register-form {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.error-message {
  color: #dc3545;
  margin-bottom: 1rem;
  padding: 0.5rem;
  background-color: #f8d7da;
  border-radius: 4px;
  text-align: center;
}

h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: #666;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

button {
  width: 100%;
  padding: 0.75rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

button:hover {
  background-color: #45a049;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.login-link {
  text-align: center;
  margin-top: 1rem;
  color: #666;
}

.login-link a {
  color: #4CAF50;
  text-decoration: none;
}

.login-link a:hover {
  text-decoration: underline;
}
</style> 