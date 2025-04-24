<template>
    <div class="post-detail">
      <div v-if="postsStore.loading" class="loading">
        Loading post...
      </div>
  
      <div v-else-if="postsStore.error" class="error">
        {{ postsStore.error }}
      </div>
  
      <template v-else>
        <div class="post">
          <div class="post-header">
            <h1 class="post-title">{{ postsStore.currentPost.title }}</h1>
            <div class="post-meta">
              <span class="post-author">By {{ postsStore.currentPost.author.username }}</span>
              <span class="post-date">{{ formatDate(postsStore.currentPost.created_at) }}</span>
            </div>
          </div>
          
          <div class="post-content">
            {{ postsStore.currentPost.content }}
          </div>
  
          <div class="post-actions" v-if="isAuthor">
            <button @click="editPost" class="edit-button">Edit</button>
            <button @click="deletePost" class="delete-button">Delete</button>
          </div>
        </div>
  
        <div class="comments-section">
          <h2>Comments</h2>
          
          <div v-if="authStore.isAuthenticated" class="comment-form">
            <textarea
              v-model="newComment"
              placeholder="Write a comment..."
              rows="3"
            ></textarea>
            <button @click="addComment" :disabled="!newComment.trim()">Add Comment</button>
          </div>
  
          <div v-if="comments.length === 0" class="no-comments">
            No comments yet. Be the first to comment!
          </div>
  
          <div v-else class="comments-list">
            <div v-for="comment in comments" :key="comment.id" class="comment">
              <div class="comment-content">
                {{ comment.content }}
              </div>
              <div class="comment-meta">
                <span class="comment-author">{{ comment.author.username }}</span>
                <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </template>
  
  <script setup>
  import { ref, computed, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { usePostsStore } from '../stores/posts'
  import { useAuthStore } from '../stores/auth'
  import { useCommentsStore } from '../stores/comments'
  
  const route = useRoute()
  const router = useRouter()
  const postsStore = usePostsStore()
  const authStore = useAuthStore()
  const commentsStore = useCommentsStore()
  const newComment = ref('')
  
  const isAuthor = computed(() => {
    return authStore.user && postsStore.currentPost && 
           authStore.user.id === postsStore.currentPost.author_id
  })
  
  const comments = computed(() => commentsStore.comments)
  
  onMounted(async () => {
    await postsStore.fetchPost(route.params.id)
    await commentsStore.fetchComments(route.params.id)
  })
  
  const formatDate = (dateString) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
  
  const editPost = () => {
    router.push(`/posts/${route.params.id}/edit`)
  }
  
  const deletePost = async () => {
    if (confirm('Are you sure you want to delete this post?')) {
      try {
        await postsStore.deletePost(route.params.id)
        router.push('/posts')
      } catch (error) {
        console.error('Failed to delete post:', error)
      }
    }
  }
  
  const addComment = async () => {
    if (!newComment.value.trim()) return
  
    try {
      await commentsStore.addComment(route.params.id, newComment.value)
      newComment.value = ''
    } catch (error) {
      console.error('Failed to add comment:', error)
    }
  }
  </script>
  
  <style scoped>
  .post-detail {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .post {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 2rem;
    margin-bottom: 2rem;
  }
  
  .post-header {
    margin-bottom: 1.5rem;
  }
  
  .post-title {
    margin: 0 0 1rem;
    color: #333;
    font-size: 2rem;
  }
  
  .post-meta {
    color: #666;
    font-size: 0.9rem;
  }
  
  .post-content {
    line-height: 1.6;
    color: #333;
    margin-bottom: 1.5rem;
  }
  
  .post-actions {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
  }
  
  .edit-button, .delete-button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
  }
  
  .edit-button {
    background-color: #4CAF50;
    color: white;
  }
  
  .delete-button {
    background-color: #f44336;
    color: white;
  }
  
  .comments-section {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 2rem;
  }
  
  .comments-section h2 {
    margin: 0 0 1.5rem;
    color: #333;
  }
  
  .comment-form {
    margin-bottom: 2rem;
  }
  
  .comment-form textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    margin-bottom: 1rem;
    resize: vertical;
  }
  
  .comment-form button {
    padding: 0.5rem 1rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .comment-form button:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }
  
  .no-comments {
    text-align: center;
    color: #666;
    padding: 2rem;
  }
  
  .comments-list {
    display: grid;
    gap: 1.5rem;
  }
  
  .comment {
    padding: 1rem;
    border: 1px solid #eee;
    border-radius: 4px;
  }
  
  .comment-content {
    margin-bottom: 0.5rem;
    line-height: 1.6;
  }
  
  .comment-meta {
    color: #666;
    font-size: 0.9rem;
  }
  
  .loading, .error {
    text-align: center;
    padding: 2rem;
    color: #666;
  }
  </style>