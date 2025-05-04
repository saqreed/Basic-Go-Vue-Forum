<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <router-link to="/" class="navbar-item">Forum</router-link>
    </div>
    <div class="navbar-menu">
      <div class="navbar-end">
        <template v-if="!authStore.isAuthenticated">
          <router-link to="/login" class="navbar-item">Login</router-link>
          <router-link to="/register" class="navbar-item">Register</router-link>
        </template>
        <template v-else>
          <router-link to="/chat" class="navbar-item">Chat</router-link>
          <router-link to="/posts/create" class="navbar-item">Create Post</router-link>
          <router-link to="/profile" class="navbar-item">Profile</router-link>
          <router-link v-if="authStore.user?.role === 'admin'" to="/admin" class="navbar-item">Admin Panel</router-link>
          <a @click="logout" class="navbar-item">Logout</a>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: #f8f9fa;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.navbar-brand {
  font-size: 1.5rem;
  font-weight: bold;
}

.navbar-menu {
  display: flex;
}

.navbar-end {
  display: flex;
  gap: 1rem;
}

.navbar-item {
  text-decoration: none;
  color: #333;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.navbar-item:hover {
  background-color: #e9ecef;
}
</style> 