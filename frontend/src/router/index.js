import { createRouter, createWebHistory } from 'vue-router'
import { checkInit } from '@/api/init'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/onboarding',
    name: 'Onboarding',
    component: () => import('@/views/Onboarding.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

let initChecked = false
let initStatus = null

function isAuthenticated() {
  return !!localStorage.getItem('homemusic-token')
}

router.beforeEach(async (to, from, next) => {
  if (!initChecked) {
    try {
      const res = await checkInit()
      initStatus = res.data
    } catch (error) {
      console.error('初始化状态检查失败', error)
    }
    initChecked = true
  }

  if (initStatus && initStatus.is_initialized === false) {
    if (to.path !== '/onboarding') {
      next('/onboarding')
      return
    }
    next()
    return
  }

  if (initStatus && initStatus.is_initialized === true) {
    if (to.path === '/onboarding') {
      next('/login')
      return
    }
    if (to.path !== '/login' && !isAuthenticated()) {
      next('/login')
      return
    }
    if (to.path === '/login' && isAuthenticated()) {
      next('/')
      return
    }
  }

  next()
})

export default router
