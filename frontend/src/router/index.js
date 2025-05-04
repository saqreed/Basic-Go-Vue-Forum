import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import PostDetail from '../views/PostDetail.vue'
import Chat from '../views/Chat.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/register',
    name: 'register',
    component: Register
  },
  {
    path: '/posts',
    name: 'Posts',
    component: () => import('../views/Posts.vue')
  },
  {
    path: '/posts/:id',
    name: 'post-detail',
    component: PostDetail
  },
  {
    path: '/posts/create',
    name: 'CreatePost',
    component: () => import('../views/CreatePost.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/posts/:id/edit',
    name: 'EditPost',
    component: () => import('../views/CreatePost.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: () => import('../views/admin/Dashboard.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: () => import('../views/admin/Users.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/posts',
    name: 'AdminPosts',
    component: () => import('../views/admin/Posts.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/comments',
    name: 'AdminComments',
    component: () => import('../views/admin/Comments.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/admin/Dashboard.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/chat',
    name: 'chat',
    component: Chat
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Навигационный гвард
router.beforeEach((to, from, next) => {
  const publicPages = ['/login', '/register']
  const authRequired = !publicPages.includes(to.path)
  const authStore = useAuthStore()

  if (authRequired && !authStore.isAuthenticated) {
    return next('/login')
  }

  next()
})

export default router 