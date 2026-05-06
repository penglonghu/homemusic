import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('homemusic-token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

request.interceptors.response.use(
  (response) => {
    const { code, msg, data } = response.data
    if (code === 0) {
      return { data, msg }
    }
    if (code === 401) {
      localStorage.removeItem('homemusic-token')
    }
    ElMessage.error(msg || '操作失败')
    return Promise.reject(new Error(msg || '操作失败'))
  },
  (error) => {
    const message = error.response?.data?.msg || error.message || '网络错误'
    if (error.response?.status === 401) {
      localStorage.removeItem('homemusic-token')
    }
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default request
