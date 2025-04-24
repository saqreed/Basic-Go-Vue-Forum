<template>
  <div class="admin-users">
    <h1>User Management</h1>
    
    <div v-if="loading" class="loading">
      Loading users...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else class="users-list">
      <div v-for="user in users" :key="user.id" class="user-card">
        <div class="user-info">
          <h3>{{ user.username }}</h3>
          <p>Email: {{ user.email }}</p>
          <p>Role: {{ user.role }}</p>
          <p>Created: {{ formatDate(user.created_at) }}</p>
        </div>
        
        <div class="user-actions">
          <button 
            @click="toggleRole(user)" 
            class="role-button"
            :class="{ 'admin': user.role === 'admin' }"
          >
            {{ user.role === 'admin' ? 'Remove Admin' : 'Make Admin' }}
          </button>
          <button 
            @click="deleteUser(user.id)" 
            class="delete-button"
          >
            Delete User
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const users = ref([])
const loading = ref(true)
const error = ref(null)

const fetchUsers = async () => {
  try {
    loading.value = true
    const response = await fetch('http://localhost:8081/api/admin/users', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch users')
    }
    
    users.value = await response.json()
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const toggleRole = async (user) => {
  try {
    const newRole = user.role === 'admin' ? 'user' : 'admin'
    const response = await fetch(`http://localhost:8081/api/admin/users/${user.id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ role: newRole })
    })
    
    if (!response.ok) {
      throw new Error('Failed to update user role')
    }
    
    // Update local state
    user.role = newRole
  } catch (err) {
    error.value = err.message
  }
}

const deleteUser = async (userId) => {
  if (!confirm('Are you sure you want to delete this user?')) {
    return
  }
  
  try {
    const response = await fetch(`http://localhost:8081/api/admin/users/${userId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to delete user')
    }
    
    // Remove user from local state
    users.value = users.value.filter(user => user.id !== userId)
  } catch (err) {
    error.value = err.message
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(fetchUsers)
</script>

<style scoped>
.admin-users {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.users-list {
  display: grid;
  gap: 1.5rem;
  margin-top: 2rem;
}

.user-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info h3 {
  margin: 0 0 0.5rem;
  color: #333;
}

.user-info p {
  margin: 0.25rem 0;
  color: #666;
}

.user-actions {
  display: flex;
  gap: 1rem;
}

.role-button, .delete-button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.role-button {
  background-color: #4CAF50;
  color: white;
}

.role-button.admin {
  background-color: #f44336;
}

.delete-button {
  background-color: #f44336;
  color: white;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.error {
  color: #f44336;
}
</style> 