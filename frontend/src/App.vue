<script setup>
import Navbar from './components/Navbar.vue'
</script>

<template>
  <div class="app">
    <Navbar />
    <main class="container">
      <router-view />
    </main>
  </div>
</template>

<style>
:root {
  --primary-color: #4a90e2;
  --text-color: #333;
  --bg-color: #f5f5f5;
  --card-bg: #fff;
  --border-color: #ddd;
}

@media (prefers-color-scheme: dark) {
  :root {
    --text-color: #f5f5f5;
    --bg-color: #1a1a1a;
    --card-bg: #2d2d2d;
    --border-color: #404040;
  }
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  background-color: var(--bg-color);
  color: var(--text-color);
  line-height: 1.6;
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.container {
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem;
  flex: 1;
}

@media (max-width: 768px) {
  .container {
    padding: 1rem;
  }
}
</style>
