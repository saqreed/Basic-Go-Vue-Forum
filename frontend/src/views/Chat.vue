<template>
  <div class="chat-container">
    <div class="chat-header">
      <h1>Общий чат</h1>
      <div v-if="!authStore.isAuthenticated" class="login-prompt">
        <p>Пожалуйста, <router-link to="/login">войдите</router-link> чтобы участвовать в чате.</p>
      </div>
      <div v-else-if="!isConnected" class="connection-status">
        <p>Соединение потеряно. Попытка переподключения...</p>
      </div>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div v-for="message in messages" :key="message.id" class="message" :class="{ 'own-message': isOwnMessage(message) }">
        <div class="message-header">
          <span class="message-author">{{ message.user.username }}</span>
          <span class="message-time">{{ formatTime(message.created_at) }}</span>
        </div>
        <div class="message-content">{{ message.content }}</div>
      </div>
    </div>

    <div v-if="authStore.isAuthenticated" class="chat-input">
      <textarea
        v-model="newMessage"
        placeholder="Введите сообщение..."
        @keyup.enter="sendMessage"
        :disabled="!isConnected"
      ></textarea>
      <button @click="sendMessage" :disabled="!newMessage.trim() || !isConnected">
        {{ isConnected ? 'Отправить' : 'Отключено' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const messages = ref([])
const newMessage = ref('')
const socket = ref(null)
const isConnected = ref(false)
const messagesContainer = ref(null)

onMounted(() => {
  connectWebSocket()
})

onUnmounted(() => {
  if (socket.value) {
    socket.value.close()
  }
})

function connectWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const token = localStorage.getItem('token')
  if (!token) {
    console.error('No authentication token found')
    return
  }
  
  // Используем правильный URL для WebSocket соединения
  const wsUrl = `${protocol}//${window.location.hostname}:8081/ws/chat?token=${token}`
  console.log('Connecting to WebSocket:', wsUrl)
  
  try {
    socket.value = new WebSocket(wsUrl)
    
    socket.value.onopen = () => {
      console.log('WebSocket connected')
      isConnected.value = true
    }
    
    socket.value.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        messages.value.unshift(message)
        scrollToBottom()
      } catch (error) {
        console.error('Error parsing message:', error)
      }
    }
    
    socket.value.onclose = (event) => {
      console.log('WebSocket disconnected:', event.code, event.reason)
      isConnected.value = false
      // Попытка переподключения через 5 секунд
      setTimeout(connectWebSocket, 5000)
    }
    
    socket.value.onerror = (error) => {
      console.error('WebSocket error:', error)
      isConnected.value = false
    }

    // Устанавливаем обработчик пинга
    socket.value.onping = () => {
      socket.value.pong()
    }
  } catch (error) {
    console.error('Error creating WebSocket:', error)
    isConnected.value = false
  }
}

function sendMessage() {
  if (!newMessage.value.trim() || !isConnected.value) return
  
  try {
    const message = {
      content: newMessage.value.trim()
    }
    
    if (socket.value.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify(message))
      newMessage.value = ''
    } else {
      console.error('WebSocket is not open')
      isConnected.value = false
    }
  } catch (error) {
    console.error('Error sending message:', error)
    isConnected.value = false
  }
}

function isOwnMessage(message) {
  return message.user.id === authStore.user?.id
}

function formatTime(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function scrollToBottom() {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.chat-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 1rem;
  height: calc(100vh - 2rem);
  display: flex;
  flex-direction: column;
}

.chat-header {
  margin-bottom: 1rem;
}

.connection-status {
  text-align: center;
  padding: 0.5rem;
  background-color: #fff3cd;
  color: #856404;
  border-radius: 4px;
  margin-top: 0.5rem;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column-reverse;
}

.message {
  background-color: white;
  padding: 0.75rem;
  border-radius: 8px;
  margin-bottom: 0.75rem;
  max-width: 80%;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.own-message {
  align-self: flex-end;
  background-color: #e3f2fd;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.message-author {
  font-weight: 500;
  color: #2c3e50;
}

.message-time {
  color: #666;
}

.message-content {
  color: #333;
  line-height: 1.4;
}

.chat-input {
  display: flex;
  gap: 0.5rem;
}

.chat-input textarea {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: none;
  height: 3rem;
}

.chat-input textarea:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.chat-input button {
  padding: 0.75rem 1.5rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.chat-input button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.login-prompt {
  text-align: center;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 4px;
  margin-top: 1rem;
}

.login-prompt a {
  color: #4CAF50;
  text-decoration: none;
}

.login-prompt a:hover {
  text-decoration: underline;
}
</style> 