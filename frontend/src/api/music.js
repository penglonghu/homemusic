import request from '@/utils/request'

export function fetchMusicList(page = 1, pageSize = 20) {
  return request({
    url: '/music/list',
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  })
}

export function searchMusic(keyword, page = 1, pageSize = 20) {
  return request({
    url: '/music/search',
    method: 'get',
    params: {
      keyword,
      page,
      page_size: pageSize
    }
  })
}

export function streamMusicUrl(id) {
  return `/api/music/stream/${id}`
}

export function recordPlayHistory(musicId) {
  return request({
    url: '/music/recent',
    method: 'post',
    data: {
      music_id: musicId
    }
  })
}

export function fetchRecentPlays(page = 1, pageSize = 20) {
  return request({
    url: '/music/recent',
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  })
}

export function scanMusic() {
  return request({
    url: '/music/scan',
    method: 'post'
  })
}