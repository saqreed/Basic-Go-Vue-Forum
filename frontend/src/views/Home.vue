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
  width: 100%;
  margin: 0 auto;
  padding: 1.5rem;
  text-align: center;
}

h1 {
  margin-bottom: 1.5rem;
  color: var(--text-color);
  font-size: 2rem;
}

.content {
  background-color: var(--card-bg);
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--border-color);
}

ul {
  list-style: none;
  padding: 0;
  margin: 1.5rem 0;
  text-align: left;
}

li {
  margin: 1rem 0;
  color: var(--text-color);
  font-size: 1.1rem;
  padding-left: 1.5rem;
  position: relative;
}

li::before {
  content: "•";
  color: var(--primary-color);
  position: absolute;
  left: 0;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 1.5rem;
  flex-wrap: wrap;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
  min-width: 200px;
}

.primary {
  background-color: var(--primary-color);
  color: white;
}

.primary:hover {
  background-color: #357abd;
  transform: translateY(-1px);
}

.secondary {
  background-color: var(--card-bg);
  color: var(--text-color);
  border: 1px solid var(--border-color);
}

.secondary:hover {
  background-color: var(--border-color);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .home {
    padding: 1rem;
  }

  h1 {
    font-size: 1.5rem;
  }

  .content {
    padding: 1rem;
  }

  li {
    font-size: 1rem;
  }

  .btn {
    width: 100%;
    min-width: unset;
  }

  .actions {
    flex-direction: column;
    gap: 0.75rem;
  }
}
</style> 