<template>
  <div class="chat-container">
    <div class="messages" ref="messagesContainer">
      <div v-for="message in messages" :key="message.id" class="message" :class="{ 'my-message': message.user_id === currentUser.id }">
        <div v-if="message.reply_to" class="reply-to">
          <span class="reply-username">{{ message.reply_to.username }}</span>
          <span class="reply-content">{{ message.reply_to.content }}</span>
        </div>
        <div class="message-content">
          <div class="message-header">
            <span class="username">{{ message.username }}</span>
            <span class="time">{{ formatTime(message.created_at) }}</span>
          </div>
          <p class="text">{{ message.content }}</p>
        </div>
        <button class="reply-button" @click="startReply(message)">Reply</button>
      </div>
    </div>
    <div class="input-container">
      <div v-if="replyingTo" class="reply-preview">
        <span>Replying to {{ replyingTo.username }}:</span>
        <span class="reply-content">{{ replyingTo.content }}</span>
        <button class="cancel-reply" @click="cancelReply">×</button>
      </div>
      <input
        v-model="newMessage"
        @keyup.enter="sendMessage"
        placeholder="Type a message..."
        :disabled="!isConnected"
      />
      <button @click="sendMessage" :disabled="!isConnected || !newMessage.trim()">
        Send
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useUserStore } from '../stores/user'
import { formatTime } from '../utils/time'

const userStore = useUserStore()
const currentUser = userStore.user

const messages = ref([])
const newMessage = ref('')
const socket = ref(null)
const isConnected = ref(false)
const messagesContainer = ref(null)
const replyingTo = ref(null)

const connectWebSocket = () => {
  const token = localStorage.getItem('token')
  if (!token) return

  socket.value = new WebSocket(`ws://localhost:8080/ws/chat?token=${token}`)

  socket.value.onopen = () => {
    console.log('WebSocket connected')
    isConnected.value = true
  }

  socket.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    if (Array.isArray(data)) {
      messages.value = data
    } else {
      messages.value.push(data)
    }
    scrollToBottom()
  }

  socket.value.onclose = () => {
    console.log('WebSocket disconnected')
    isConnected.value = false
    setTimeout(connectWebSocket, 1000)
  }

  socket.value.onerror = (error) => {
    console.error('WebSocket error:', error)
    isConnected.value = false
  }
}

const sendMessage = () => {
  if (!newMessage.value.trim() || !socket.value) return

  const message = {
    content: newMessage.value.trim(),
    reply_to_id: replyingTo.value?.id
  }

  socket.value.send(JSON.stringify(message))
  newMessage.value = ''
  replyingTo.value = null
}

const startReply = (message) => {
  replyingTo.value = {
    id: message.id,
    username: message.username,
    content: message.content
  }
}

const cancelReply = () => {
  replyingTo.value = null
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

onMounted(() => {
  connectWebSocket()
})

onUnmounted(() => {
  if (socket.value) {
    socket.value.close()
  }
})
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.message {
  display: flex;
  flex-direction: column;
  max-width: 70%;
  padding: 10px;
  border-radius: 8px;
  background-color: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.my-message {
  align-self: flex-end;
  background-color: #e3f2fd;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.username {
  font-weight: bold;
  color: #1976d2;
}

.time {
  font-size: 0.8em;
  color: #666;
}

.text {
  margin: 0;
  word-break: break-word;
}

.input-container {
  padding: 20px;
  background-color: white;
  border-top: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.reply-preview {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-size: 0.9em;
}

.reply-content {
  color: #666;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cancel-reply {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 1.2em;
  padding: 0 5px;
}

.cancel-reply:hover {
  color: #f44336;
}

.reply-to {
  font-size: 0.9em;
  color: #666;
  margin-bottom: 5px;
  padding: 5px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.reply-username {
  color: #1976d2;
  font-weight: bold;
  margin-right: 5px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  font-size: 1em;
}

button {
  padding: 10px 20px;
  background-color: #1976d2;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1em;
}

button:disabled {
  background-color: #bdbdbd;
  cursor: not-allowed;
}

.reply-button {
  align-self: flex-end;
  background: none;
  color: #1976d2;
  padding: 5px 10px;
  font-size: 0.9em;
}

.reply-button:hover {
  background-color: #e3f2fd;
}
</style> 