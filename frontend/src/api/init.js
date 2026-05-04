import request from '@/utils/request'

// 检查初始化状态
export function checkInit() {
  return request({
    url: '/api/v1/init/check',
    method: 'get'
  })
}

// 执行初始化
export function execInit(data) {
  return request({
    url: '/api/v1/init/exec',
    method: 'post',
    data
  })
}
