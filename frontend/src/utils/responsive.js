/**
 * 响应式布局工具函数
 */

// 设备检测
export function getDeviceType() {
  const width = window.innerWidth
  if (width <= 768) {
    return 'mobile' // 手机
  } else if (width <= 1024) {
    return 'tablet' // 平板/车机
  } else {
    return 'desktop' // PC
  }
}

// 断点常量
export const BREAKPOINTS = {
  mobile: 768,
  tablet: 1024,
  desktop: 1200
}

// 检查是否为移动端
export function isMobile() {
  return window.innerWidth <= BREAKPOINTS.mobile
}

// 检查是否为平板或车机
export function isTablet() {
  const width = window.innerWidth
  return width > BREAKPOINTS.mobile && width <= BREAKPOINTS.tablet
}

// 检查是否为桌面端
export function isDesktop() {
  return window.innerWidth > BREAKPOINTS.tablet
}