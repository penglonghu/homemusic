import axios from 'axios'

export function getHealthStatus() {
  return axios({
    url: '/api/health',
    method: 'get'
  })
}
