/**
 * HomeMusic 多端适配模块
 * 实现响应式布局，自动适配PC/手机/车机浏览器
 */

// 响应式工具函数
export { 
  getDeviceType, 
  isMobile, 
  isTablet, 
  isDesktop, 
  BREAKPOINTS 
} from './utils/responsive'

// 响应式组合式API
export { 
  responsiveEvents,
  responsiveConfig,
  responsiveClasses
} from './composables/useResponsive'

// 响应式组件
export { default as ResponsiveContainer } from './components/ResponsiveContainer.vue'

console.log('HomeMusic 响应式模块已加载')
console.log('当前设备类型:', getDeviceType())