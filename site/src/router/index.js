import { createRouter, createWebHashHistory } from 'vue-router'
import { getToken } from '../store/auth'

const routes = [
  { path: '/', name: 'home', component: () => import('../views/Home.vue') },
  { path: '/register', name: 'register', component: () => import('../views/Register.vue') },
  { path: '/login', name: 'login', component: () => import('../views/Login.vue') },
  { path: '/dashboard', name: 'dashboard', component: () => import('../views/Dashboard.vue'), meta: { requiresAuth: true } }
]

const router = createRouter({ history: createWebHashHistory(), routes })

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getToken()) {
    return { name: 'login' }
  }
})

export default router
