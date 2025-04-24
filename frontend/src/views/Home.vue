<template>
  <div class="home">
    <h1>Добро пожаловать на наш форум!</h1>
    <div class="content">
      <p>Здесь вы можете:</p>
      <ul>
        <li>Читать интересные посты</li>
        <li>Делиться своими мыслями</li>
        <li>Общаться с другими участниками</li>
        <li>Создавать собственные темы для обсуждения</li>
      </ul>
      <div class="actions">
        <router-link to="/posts" class="btn primary">Просмотреть посты</router-link>
        <router-link to="/posts/create" class="btn secondary">Создать пост</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
// Компонент домашней страницы не требует дополнительной логики
</script>

<style scoped>
.home {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}

h1 {
  margin-bottom: 2rem;
  color: #2c3e50;
}

.content {
  background-color: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

ul {
  list-style: none;
  padding: 0;
  margin: 2rem 0;
}

li {
  margin: 1rem 0;
  color: #2c3e50;
  font-size: 1.1rem;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 2rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: bold;
  transition: background-color 0.3s;
}

.primary {
  background-color: #4CAF50;
  color: white;
}

.primary:hover {
  background-color: #45a049;
}

.secondary {
  background-color: #2196F3;
  color: white;
}

.secondary:hover {
  background-color: #1e88e5;
}
</style> 