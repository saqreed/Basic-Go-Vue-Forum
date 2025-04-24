import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/posts',
    name: 'Posts',
    component: () => import('../views/Posts.vue')
  },
  {
    path: '/posts/:id',
    name: 'PostDetail',
    component: () => import('../views/PostDetail.vue')
  },
  {
    path: '/posts/create',
    name: 'CreatePost',
    component: () => import('../views/CreatePost.vue')
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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Добавляем проверку на админские права и аутентификацию
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresAdmin && (!authStore.isAuthenticated || authStore.user?.role !== 'admin')) {
    next('/login')
  } else {
    next()
  }
})

export default router 