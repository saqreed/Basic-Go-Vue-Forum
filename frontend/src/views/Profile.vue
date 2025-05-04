<template>
  <div class="profile">
    <h1>Профиль</h1>
    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="profile-content">
      <div class="profile-info">
        <h2>Информация о пользователе</h2>
        <p><strong>Имя пользователя:</strong> {{ user.username }}</p>
        <p><strong>Email:</strong> {{ user.email }}</p>
      </div>
      <div class="profile-stats">
        <h2>Статистика</h2>
        <p><strong>Количество постов:</strong> {{ stats.posts }}</p>
        <p><strong>Количество комментариев:</strong> {{ stats.comments }}</p>
      </div>
      <div class="change-password">
        <h2>Смена пароля</h2>
        <form @submit.prevent="handleChangePassword">
          <div class="form-group">
            <label for="current-password">Текущий пароль</label>
            <input
              type="password"
              id="current-password"
              v-model="currentPassword"
              required
            />
          </div>
          <div class="form-group">
            <label for="new-password">Новый пароль</label>
            <input
              type="password"
              id="new-password"
              v-model="newPassword"
              required
            />
          </div>
          <div class="form-group">
            <label for="confirm-password">Подтвердите новый пароль</label>
            <input
              type="password"
              id="confirm-password"
              v-model="confirmPassword"
              required
            />
          </div>
          <button type="submit" :disabled="loading">Изменить пароль</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'Profile',
  setup() {
    const authStore = useAuthStore()
    const loading = ref(false)
    const error = ref(null)
    const user = ref({ username: '', email: '' })
    const stats = ref({ posts: 0, comments: 0 })
    const currentPassword = ref('')
    const newPassword = ref('')
    const confirmPassword = ref('')

    const fetchProfile = async () => {
      loading.value = true
      error.value = null
      try {
        console.log('Fetching profile with token:', authStore.token ? 'exists' : 'not found')
        const response = await fetch('http://localhost:8081/api/profile', {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        console.log('Profile response status:', response.status)
        if (!response.ok) {
          const errorText = await response.text()
          console.error('Profile error response:', errorText)
          throw new Error('Ошибка при загрузке профиля')
        }
        user.value = await response.json()
        console.log('Profile data loaded:', user.value)
      } catch (err) {
        console.error('Profile fetch error:', err)
        error.value = err.message
      } finally {
        loading.value = false
      }
    }

    const fetchStats = async () => {
      try {
        console.log('Fetching stats with token:', authStore.token ? 'exists' : 'not found')
        const response = await fetch('http://localhost:8081/api/profile/stats', {
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          }
        })
        console.log('Stats response status:', response.status)
        if (!response.ok) {
          const errorText = await response.text()
          console.error('Stats error response:', errorText)
          throw new Error('Ошибка при загрузке статистики')
        }
        stats.value = await response.json()
        console.log('Stats data loaded:', stats.value)
      } catch (err) {
        console.error('Stats fetch error:', err)
      }
    }

    const handleChangePassword = async () => {
      if (newPassword.value !== confirmPassword.value) {
        error.value = 'Пароли не совпадают'
        return
      }

      loading.value = true
      error.value = null
      try {
        const response = await fetch('http://localhost:8081/api/profile/password', {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${authStore.token}`
          },
          body: JSON.stringify({
            currentPassword: currentPassword.value,
            newPassword: newPassword.value
          })
        })
        if (!response.ok) {
          throw new Error('Ошибка при смене пароля')
        }
        currentPassword.value = ''
        newPassword.value = ''
        confirmPassword.value = ''
        alert('Пароль успешно изменен')
      } catch (err) {
        error.value = err.message
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      fetchProfile()
      fetchStats()
    })

    return {
      loading,
      error,
      user,
      stats,
      currentPassword,
      newPassword,
      confirmPassword,
      handleChangePassword,
      fetchProfile,
      fetchStats
    }
  }
}
</script>

<style scoped>
.profile {
  max-width: 800px;
  width: 100%;
  margin: 0 auto;
  padding: 20px;
}

.profile-content {
  display: grid;
  gap: 20px;
}

.profile-info,
.profile-stats,
.change-password {
  background: var(--card-bg);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--border-color);
}

h1 {
  font-size: 2rem;
  margin-bottom: 1.5rem;
  color: var(--text-color);
}

h2 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: var(--text-color);
  font-size: 1.5rem;
}

p {
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--card-bg);
  color: var(--text-color);
  font-size: 1rem;
}

button {
  background: var(--primary-color);
  color: white;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

button:hover {
  background: #357abd;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: var(--text-color);
}

.error {
  color: #e74c3c;
  padding: 1rem;
  background: #fde8e8;
  border-radius: 4px;
  margin-bottom: 1rem;
}

@media (max-width: 768px) {
  .profile {
    padding: 10px;
  }
  
  h1 {
    font-size: 1.5rem;
  }
  
  h2 {
    font-size: 1.25rem;
  }
}
</style> 