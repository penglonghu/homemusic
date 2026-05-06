<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="logo">HomeMusic</div>
      <el-form ref="formRef" :model="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="login-btn" @click="login">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login as loginApi } from '@/api/auth'

const router = useRouter()
const formRef = ref(null)
const form = ref({
  username: '',
  password: ''
})

const login = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  try {
    const res = await loginApi({
      username: form.value.username,
      password: form.value.password,
      remember: true
    })
    if (res.data && res.data.token) {
      localStorage.setItem('homemusic-token', res.data.token)
      ElMessage.success('登录成功')
      router.push('/')
      return
    }
    ElMessage.error('登录失败，请重试')
  } catch (error) {
    ElMessage.error(error.message || '登录失败')
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
.login-card {
  width: 400px;
}
.logo {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
  color: #409eff;
}
.login-btn {
  width: 100%;
}
</style>
