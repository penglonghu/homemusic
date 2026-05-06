/**
 * HomeMusic 响应式布局模块
 * 实现多端适配，自动适配PC/手机/车机浏览器
 */

// 导出响应式工具函数
export { 
  getDeviceType, 
  isMobile, 
  isTablet, 
  isDesktop, 
  BREAKPOINTS 
} from './utils/responsive'

// 导出响应式组件
export { default as ResponsiveContainer } from './components/ResponsiveContainer.vue'

// 导出响应式样式类
export const responsiveClasses = {
  // 设备检测类
  deviceClass: (deviceType) => `device-${deviceType}`,
  
  // 响应式辅助类
  getResponsiveClass: (baseClass, modifiers = {}) => {
    const classes = [baseClass]
    
    Object.entries(modifiers).forEach(([modifier, condition]) => {
      if (condition) {
        classes.push(`${baseClass}--${modifier}`)
      }
    })
    
    return classes.join(' ')
  }
}

// 响应式事件处理器
export const responsiveEvents = {
  // 添加响应式事件监听
  addResponsiveListener: (callback) => {
    const handler = () => callback(getDeviceType())
    window.addEventListener('resize', handler)
    // 立即执行一次
    handler()
    return () => window.removeEventListener('resize', handler)
  },
  
  // 获取当前视口信息
  getViewportInfo: () => ({
    width: window.innerWidth,
    height: window.innerHeight,
    deviceType: getDeviceType(),
    isPortrait: window.innerHeight > window.innerWidth,
    isLandscape: window.innerWidth > window.innerHeight
  })
}

// 导出响应式配置
export const responsiveConfig = {
  breakpoints: BREAKPOINTS,
  defaultTransition: 'all 0.3s ease',
  mobileOptimized: true
}