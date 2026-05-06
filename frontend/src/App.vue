<template>
  <div id="app" :class="deviceClass">
    <el-container :style="{ height: appHeight }">
      <el-header class="app-header">
        <div class="header-content">
          <div class="logo" :class="{ 'logo-mobile': isMobileView }">
            <h2>HomeMusic</h2>
          </div>
          <div class="menu-toggle" v-if="isMobileView" @click="toggleSidebar">
            <i class="el-icon-menu"></i>
          </div>
          <div class="user-actions">
            <el-dropdown>
              <i class="el-icon-user"></i>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item>个人中心</el-dropdown-item>
                  <el-dropdown-item>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside 
          :class="['app-sidebar', { 'sidebar-mobile': isMobileView, 'sidebar-hidden': sidebarHidden }]" 
          :style="{ width: sidebarWidth }"
        >
          <el-menu 
            :default-openeds="['1']" 
            :collapse="isMobileView && sidebarHidden"
            :style="{ minHeight: '100%' }"
          >
            <el-sub-menu index="1">
              <template #title>
                <i class="el-icon-menu"></i>
                <span v-show="!isMobileView || !sidebarHidden">音乐库</span>
              </template>
              <el-menu-item index="1-1">全部歌曲</el-menu-item>
              <el-menu-item index="1-2">我的歌单</el-menu-item>
            </el-sub-menu>
            <el-sub-menu index="2">
              <template #title>
                <i class="el-icon-setting"></i>
                <span v-show="!isMobileView || !sidebarHidden">设置</span>
              </template>
              <el-menu-item index="2-1">系统设置</el-menu-item>
            </el-sub-menu>
          </el-menu>
        </el-aside>
        <el-main class="app-main">
          <router-view />
        </el-main>
      </el-container>
      <el-footer class="app-footer" v-if="!isMobileView">
        HomeMusic © 2024 家庭音乐播放器
      </el-footer>
    </el-container>
    <PlayerBar />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import PlayerBar from '@/components/PlayerBar.vue'
import { getDeviceType, isMobile, BREAKPOINTS } from '@/utils/responsive'

const sidebarHidden = ref(true)
const windowWidth = ref(window.innerWidth)

const deviceType = computed(() => getDeviceType())
const deviceClass = computed(() => `device-${deviceType.value}`)
const isMobileView = computed(() => windowWidth.value <= BREAKPOINTS.mobile)
const isTabletView = computed(() => 
  windowWidth.value > BREAKPOINTS.mobile && windowWidth.value <= BREAKPOINTS.tablet
)
const isDesktopView = computed(() => windowWidth.value > BREAKPOINTS.tablet)
const sidebarWidth = computed(() => {
  if (isMobileView.value && sidebarHidden.value) {
    return '0px'
  }
  return isMobileView.value ? '200px' : '240px'
})
const appHeight = computed(() => isMobileView.value ? 'calc(100vh - 60px)' : '100vh')

function toggleSidebar() {
  sidebarHidden.value = !sidebarHidden.value
}

function handleResize() {
  windowWidth.value = window.innerWidth
  // 在移动端，窗口大小改变时自动隐藏侧边栏
  if (window.innerWidth > BREAKPOINTS.mobile) {
    sidebarHidden.value = false
  } else if (window.innerWidth <= BREAKPOINTS.mobile && windowWidth.value > BREAKPOINTS.mobile) {
    sidebarHidden.value = true
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  // 初始判断设备类型
  handleResize()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.app-header {
  padding: 0;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  height: 60px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.logo h2 {
  margin: 0;
  color: #409EFF;
}

.logo-mobile h2 {
  font-size: 1.2rem;
}

.menu-toggle {
  font-size: 1.5rem;
  cursor: pointer;
  padding: 10px;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.app-sidebar {
  transition: all 0.3s ease;
  background-color: #f5f7fa;
  border-right: 1px solid #dfe4ed;
}

.app-sidebar.sidebar-mobile {
  position: fixed;
  z-index: 1000;
  height: calc(100vh - 60px);
  top: 60px;
}

.app-sidebar.sidebar-hidden {
  transform: translateX(-100%);
}

.app-main {
  padding: 20px;
}

.app-footer {
  text-align: center;
  padding: 10px;
  color: #909399;
  font-size: 14px;
  border-top: 1px solid #e6e6e6;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .app-header {
    height: 50px;
  }
  
  .header-content {
    padding: 0 10px;
  }
  
  .logo-mobile h2 {
    font-size: 1.1rem;
  }
  
  .app-main {
    padding: 10px;
  }
  
  .user-actions {
    gap: 8px;
  }
  
  .menu-toggle {
    font-size: 1.2rem;
  }
  
  .device-mobile .el-aside {
    position: fixed;
    z-index: 1000;
  }
  
  .el-footer.app-footer {
    display: none;
  }
}

/* 平板/车机适配 */
@media (min-width: 769px) and (max-width: 1024px) {
  .app-sidebar {
    width: 200px !important;
  }
  
  .app-header {
    height: 55px;
  }
  
  .header-content {
    padding: 0 15px;
  }
  
  .app-main {
    padding: 15px;
  }
}

/* 桌面端优化 */
@media (min-width: 1200px) {
  .app-sidebar {
    width: 250px !important;
  }
}
</style>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

/* 导入响应式样式 */
@import '@/styles/responsive.css';

.app-main {
  padding-bottom: 120px;
}
</style>