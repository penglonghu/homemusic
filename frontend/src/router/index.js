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

  if (initStatus && initStatus.is_initialized === false && to.path !== '/onboarding') {
    next('/onboarding')
    return
  }

  if (initStatus && initStatus.is_initialized === true && to.path === '/onboarding') {
    next('/login')
    return
  }

  next()
})

export default router
