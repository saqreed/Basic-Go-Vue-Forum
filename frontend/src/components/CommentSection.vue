<template>
  <div class="comment-section">
    <h3>Комментарии</h3>
    
    <!-- Форма добавления комментария -->
    <div class="add-comment">
      <textarea
        v-model="newComment"
        placeholder="Напишите комментарий..."
        rows="3"
      ></textarea>
      <button @click="submitComment" :disabled="!newComment.trim()">
        Отправить
      </button>
    </div>

    <!-- Список комментариев -->
    <div v-if="loading" class="loading">Загрузка комментариев...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="comments.length === 0" class="no-comments">
      Пока нет комментариев. Будьте первым!
    </div>
    <div v-else class="comments-list">
      <div v-for="comment in comments" :key="comment.id" class="comment">
        <div class="comment-header">
          <span class="author">{{ comment.author.username }}</span>
          <span class="date">{{ formatDate(comment.created_at) }}</span>
        </div>
        <div class="comment-content">
          <div v-if="!comment.isEditing">{{ comment.content }}</div>
          <textarea
            v-else
            v-model="comment.editContent"
            rows="3"
          ></textarea>
        </div>
        <div class="comment-actions" v-if="isAuthor(comment)">
          <button
            v-if="!comment.isEditing"
            @click="startEdit(comment)"
            class="edit-btn"
          >
            Редактировать
          </button>
          <button
            v-else
            @click="saveEdit(comment)"
            :disabled="!comment.editContent.trim()"
            class="save-btn"
          >
            Сохранить
          </button>
          <button
            v-if="!comment.isEditing"
            @click="deleteComment(comment.id)"
            class="delete-btn"
          >
            Удалить
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useCommentsStore } from '../stores/comments'
import { useAuthStore } from '../stores/auth'

const props = defineProps({
  postId: {
    type: Number,
    required: true
  }
})

const commentsStore = useCommentsStore()
const authStore = useAuthStore()
const newComment = ref('')
const loading = ref(false)
const error = ref(null)

const comments = computed(() => commentsStore.comments)

onMounted(async () => {
  await loadComments()
})

async function loadComments() {
  loading.value = true
  try {
    await commentsStore.fetchComments(props.postId)
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

async function submitComment() {
  if (!newComment.value.trim()) return
  
  try {
    await commentsStore.addComment(props.postId, newComment.value)
    newComment.value = ''
  } catch (err) {
    error.value = err.message
  }
}

function startEdit(comment) {
  comment.isEditing = true
  comment.editContent = comment.content
}

async function saveEdit(comment) {
  if (!comment.editContent.trim()) return
  
  try {
    await commentsStore.updateComment(comment.id, comment.editContent)
    comment.isEditing = false
  } catch (err) {
    error.value = err.message
  }
}

async function deleteComment(commentId) {
  if (!confirm('Вы уверены, что хотите удалить этот комментарий?')) return
  
  try {
    await commentsStore.deleteComment(commentId)
  } catch (err) {
    error.value = err.message
  }
}

function isAuthor(comment) {
  return authStore.user && authStore.user.id === comment.author.id
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleString()
}
</script>

<style scoped>
.comment-section {
  margin-top: 2rem;
}

.add-comment {
  margin-bottom: 2rem;
}

.add-comment textarea {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.comment {
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #f9f9f9;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  color: #666;
}

.comment-content {
  margin-bottom: 0.5rem;
}

.comment-actions {
  display: flex;
  gap: 0.5rem;
}

button {
  padding: 0.25rem 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.edit-btn {
  background-color: #4CAF50;
  color: white;
}

.save-btn {
  background-color: #2196F3;
  color: white;
}

.delete-btn {
  background-color: #f44336;
  color: white;
}

.loading, .error, .no-comments {
  text-align: center;
  padding: 1rem;
  color: #666;
}

.error {
  color: #f44336;
}
</style> 