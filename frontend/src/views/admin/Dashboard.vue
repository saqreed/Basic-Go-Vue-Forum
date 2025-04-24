<template>
  <div class="admin-dashboard">
    <h1>Admin Dashboard</h1>
    
    <div class="admin-nav">
      <router-link to="/admin/users" class="admin-nav-item">
        <h2>Users</h2>
        <p>Manage forum users</p>
      </router-link>
      
      <router-link to="/admin/posts" class="admin-nav-item">
        <h2>Posts</h2>
        <p>Manage forum posts</p>
      </router-link>
      
      <router-link to="/admin/comments" class="admin-nav-item">
        <h2>Comments</h2>
        <p>Manage forum comments</p>
      </router-link>
    </div>
  </div>
</template>

<script setup>
// Компонент не требует дополнительной логики
</script>

<style scoped>
.admin-dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.admin-nav {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.admin-nav-item {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s, box-shadow 0.2s;
}

.admin-nav-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.admin-nav-item h2 {
  margin: 0 0 1rem;
  color: #333;
}

.admin-nav-item p {
  margin: 0;
  color: #666;
}
</style> 