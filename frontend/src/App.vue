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
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  flex: 1;
}
</style>
