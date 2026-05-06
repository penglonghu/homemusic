<template>
  <div class="onboarding-page">
    <el-card class="onboarding-card">
      <div class="title">
        <h2>欢迎使用 HomeMusic</h2>
        <p>首次访问，请完成初始化创建管理员账号和配置音乐目录。</p>
      </div>

      <el-steps :active="activeStep" align-center>
        <el-step title="管理员" />
        <el-step title="音乐目录" />
        <el-step title="完成" />
      </el-steps>

      <div class="content">
        <div v-if="activeStep === 0">
          <el-form ref="adminFormRef" :model="adminForm" :rules="adminRules" label-width="120px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="adminForm.username" placeholder="请输入管理员用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="adminForm.password" type="password" placeholder="请输入管理员密码" show-password />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="adminForm.email" placeholder="请输入管理员邮箱（可选）" />
            </el-form-item>
          </el-form>
        </div>

        <div v-else-if="activeStep === 1">
          <el-form ref="musicFormRef" :model="musicForm" :rules="musicRules" label-width="120px">
            <el-form-item label="音乐目录" prop="music_path">
              <el-input v-model="musicForm.music_path" placeholder="请输入音乐目录绝对路径" />
            </el-form-item>
            <div class="hint">请填入本地音乐目录，HomeMusic 将在此目录下扫描音乐文件。</div>
          </el-form>
        </div>

        <div v-else>
          <div class="summary">
            <h3>确认信息</h3>
            <p>管理员：<strong>{{ adminForm.username }}</strong></p>
            <p>音乐目录：<strong>{{ musicForm.music_path }}</strong></p>
            <el-alert title="确认无误后完成初始化" type="success" show-icon />
          </div>
        </div>
      </div>

      <div class="actions">
        <el-button @click="prevStep" :disabled="activeStep === 0">上一步</el-button>
        <el-button v-if="activeStep < 2" type="primary" @click="nextStep" :loading="loading">下一步</el-button>
        <el-button v-else type="primary" @click="submitInit" :loading="loading">完成初始化</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { execInit } from '@/api/init'

const router = useRouter()
const activeStep = ref(0)
const loading = ref(false)
const adminFormRef = ref(null)
const musicFormRef = ref(null)

const adminForm = reactive({
  username: '',
  password: '',
  email: '',
})

const musicForm = reactive({
  music_path: '',
})

const adminRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名长度在3-32个字符之间', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 64, message: '密码长度在6-64个字符之间', trigger: 'blur' },
  ],
  email: [
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' },
  ],
}

const musicRules = {
  music_path: [
    { required: true, message: '请输入音乐目录', trigger: 'blur' },
    { min: 1, max: 255, message: '路径长度不能超过255个字符', trigger: 'blur' },
  ],
}

const nextStep = async () => {
  if (activeStep.value === 0) {
    await adminFormRef.value.validate(async (valid) => {
      if (!valid) return
      activeStep.value++
    })
    return
  }

  if (activeStep.value === 1) {
    await musicFormRef.value.validate(async (valid) => {
      if (!valid) return
      activeStep.value++
    })
  }
}

const prevStep = () => {
  if (activeStep.value > 0) {
    activeStep.value--
  }
}

const submitInit = async () => {
  loading.value = true
  try {
    await execInit({
      username: adminForm.username,
      password: adminForm.password,
      email: adminForm.email,
      music_path: musicForm.music_path,
    })
    ElMessage.success('初始化完成，正在进入系统')
    router.push('/login')
  } catch (error) {
    ElMessage.error(error.message || '初始化失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.onboarding-page {
  min-height: calc(100vh - 40px);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.onboarding-card {
  width: 100%;
  max-width: 760px;
  padding: 30px;
}

.title {
  text-align: center;
  margin-bottom: 24px;
}

.title h2 {
  margin-bottom: 10px;
}

.content {
  margin-top: 24px;
}

.summary {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 6px;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>
