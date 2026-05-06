import { reactive, readonly } from 'vue'

const STORAGE_KEY = 'homemusic-player-state'

const state = reactive({
  playlist: [],
  currentIndex: 0,
  currentTrack: null,
  playing: false,
  currentTime: 0,
  duration: 0,
  mode: 'order',
  volume: 0.8
})

function saveState() {
  const payload = {
    currentIndex: state.currentIndex,
    currentTrack: state.currentTrack,
    currentTime: state.currentTime,
    mode: state.mode,
    volume: state.volume
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(payload))
}

function loadState() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return
    }
    const payload = JSON.parse(raw)
    if (payload && payload.currentTrack) {
      state.currentIndex = payload.currentIndex || 0
      state.currentTrack = payload.currentTrack
      state.currentTime = payload.currentTime || 0
      state.mode = payload.mode || 'order'
      state.volume = typeof payload.volume === 'number' ? payload.volume : 0.8
    }
  } catch (error) {
    console.warn('加载播放状态失败', error)
  }
}

function setPlaylist(playlist = [], index = 0) {
  state.playlist = playlist
  state.currentIndex = Math.min(Math.max(index, 0), playlist.length-1)
  state.currentTrack = playlist[state.currentIndex] || null
  state.currentTime = 0
  state.duration = 0
  state.playing = !!state.currentTrack
  saveState()
}

function playTrack(track, index) {
  if (!track) {
    return
  }
  state.playlist = state.playlist.length ? state.playlist : [track]
  state.currentIndex = typeof index === 'number' ? index : state.playlist.findIndex(item => item.id === track.id)
  if (state.currentIndex < 0) {
    state.currentIndex = 0
  }
  state.currentTrack = state.playlist[state.currentIndex]
  state.playing = true
  state.currentTime = 0
  state.duration = 0
  saveState()
}

function togglePlay() {
  state.playing = !state.playing
  saveState()
}

function setMode(mode) {
  if (['order', 'repeat', 'shuffle'].includes(mode)) {
    state.mode = mode
    saveState()
  }
}

function updateTime(time) {
  state.currentTime = time
  saveState()
}

function updateDuration(duration) {
  state.duration = duration
}

function setVolume(volume) {
  state.volume = Math.min(Math.max(volume, 0), 1)
  saveState()
}

function next() {
  if (!state.playlist.length) {
    return null
  }
  if (state.mode === 'shuffle') {
    const nextIndex = Math.floor(Math.random() * state.playlist.length)
    state.currentIndex = nextIndex
  } else {
    state.currentIndex += 1
    if (state.currentIndex >= state.playlist.length) {
      if (state.mode === 'repeat') {
        state.currentIndex = 0
      } else {
        state.currentIndex = state.playlist.length - 1
        state.playing = false
      }
    }
  }
  state.currentTrack = state.playlist[state.currentIndex]
  state.currentTime = 0
  saveState()
  return state.currentTrack
}

function prev() {
  if (!state.playlist.length) {
    return null
  }
  if (state.mode === 'shuffle') {
    const prevIndex = Math.floor(Math.random() * state.playlist.length)
    state.currentIndex = prevIndex
  } else {
    state.currentIndex -= 1
    if (state.currentIndex < 0) {
      state.currentIndex = state.mode === 'repeat' ? state.playlist.length - 1 : 0
    }
  }
  state.currentTrack = state.playlist[state.currentIndex]
  state.currentTime = 0
  saveState()
  return state.currentTrack
}

loadState()

export default function usePlayerStore() {
  return {
    state: readonly(state),
    setPlaylist,
    playTrack,
    togglePlay,
    setMode,
    updateTime,
    updateDuration,
    setVolume,
    next,
    prev,
    saveState
  }
}
