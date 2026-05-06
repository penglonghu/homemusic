<!-- 
  响应式容器组件
  根据设备类型自动调整布局
-->
<template>
  <div 
    :class="[
      'responsive-container',
      `device-${deviceType}`,
      {
        'is-mobile': isMobileView,
        'is-tablet': isTabletView,
        'is-desktop': isDesktopView
      },
      customClass
    ]"
    :style="containerStyles"
  >
    <slot :device-type="deviceType" :is-mobile="isMobileView" :is-tablet="isTabletView" :is-desktop="isDesktopView" />
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { getDeviceType, BREAKPOINTS } from '@/utils/responsive'

const props = defineProps({
  customClass: {
    type: String,
    default: ''
  },
  adaptiveLayout: {
    type: Boolean,
    default: true
  }
})

const windowWidth = ref(window.innerWidth)

const deviceType = computed(() => getDeviceType())
const isMobileView = computed(() => windowWidth.value <= BREAKPOINTS.mobile)
const isTabletView = computed(() => 
  windowWidth.value > BREAKPOINTS.mobile && windowWidth.value <= BREAKPOINTS.tablet
)
const isDesktopView = computed(() => windowWidth.value > BREAKPOINTS.tablet)

const containerStyles = computed(() => {
  if (!props.adaptiveLayout) return {}
  
  return isMobileView.value 
    ? { flexDirection: 'column', alignItems: 'stretch' }
    : {}
})

function handleResize() {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  handleResize() // 初始检测
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.responsive-container {
  display: flex;
  transition: all 0.3s ease;
}

.device-mobile {
  flex-direction: column;
}

.device-tablet {
  flex-direction: row;
}

.device-desktop {
  flex-direction: row;
}
</style>