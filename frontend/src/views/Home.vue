<template>
  <div class="home-container">
    <el-card shadow="hover">
      <div class="title">HomeMusic 服务状态</div>
      <div class="status" :class="{'success': status === 'ok', 'error': status !== 'ok'}">
        服务状态: {{ status }}
      </div>
      <div class="service">服务名称: {{ service }}</div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getHealthStatus } from '@/api/health'
import { ElMessage } from 'element-plus'

const status = ref('')
const service = ref('')

onMounted(() => {
  getHealthStatus().then(res => {
    if (res.data.code === 0) {
      status.value = res.data.data.status
      service.value = res.data.data.service
      ElMessage.success('服务连接成功')
    }
  }).catch(() => {
    status.value = 'error'
    ElMessage.error('服务连接失败')
  })
})
</script>

<style scoped>
.home-container {
  padding: 20px;
}
.title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 20px;
}
.status {
  margin-bottom: 10px;
  font-size: 16px;
}
.success {
  color: #67c23a;
}
.error {
  color: #f56c6c;
}
.service {
  font-size: 14px;
  color: #666;
}
</style>
