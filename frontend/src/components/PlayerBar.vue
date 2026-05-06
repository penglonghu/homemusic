<template>
  <div class="player-bar" :class="{ 'player-mobile': isMobileView }" v-if="player.state.currentTrack">
    <audio ref="audioRef" :src="audioSrc" preload="auto"></audio>
    <div class="player-left">
      <div class="track-info">
        <div class="track-title">{{ player.state.currentTrack.title || '未知歌曲' }}</div>
        <div class="track-meta">{{ player.state.currentTrack.artist || '未知歌手' }} • {{ player.state.currentTrack.album || '未知专辑' }}</div>
      </div>
    </div>
    <div class="player-center">
      <el-button type="text" icon="el-icon-s-backward" @click="onPrev" :size="buttonSize"></el-button>
      <el-button type="text" :icon="player.state.playing ? 'el-icon-s-pause' : 'el-icon-s-play'" @click="onPlayToggle" :size="buttonSize"></el-button>
      <el-button type="text" icon="el-icon-s-forward" @click="onNext" :size="buttonSize"></el-button>
      <el-slider
        class="player-progress"
        v-model="progress"
        :min="0"
        :max="duration"
        @change="onSeek"
        :format-tooltip="formatTime"
        :height="sliderHeight"
      />
      <div class="player-time">{{ formatTime(currentTime) }} / {{ formatTime(duration) }}</div>
    </div>
    <div class="player-right" v-if="!isMobileView">
      <div class="mode-label">模式：{{ modeText }}</div>
      <el-select v-model="mode" placeholder="播放模式" size="small" @change="onModeChange">
        <el-option label="顺序" value="order" />
        <el-option label="循环" value="repeat" />
        <el-option label="随机" value="shuffle" />
      </el-select>
      <el-slider class="volume-slider" v-model="volume" :min="0" :max="1" :step="0.01" @change="onVolumeChange" />
    </div>
    <div class="player-right-mobile" v-else>
      <el-button type="text" icon="el-icon-setting" @click="toggleMoreOptions" :size="buttonSize"></el-button>
      <div v-if="showMoreOptions" class="mobile-options">
        <div class="mode-selector">
          <span>模式：</span>
          <el-select v-model="mode" placeholder="播放模式" size="mini" @change="onModeChange">
            <el-option label="顺序" value="order" />
            <el-option label="循环" value="repeat" />
            <el-option label="随机" value="shuffle" />
          </el-select>
        </div>
        <div class="volume-control">
          <span>音量：</span>
          <el-slider class="volume-slider" v-model="volume" :min="0" :max="1" :step="0.01" @change="onVolumeChange" :vertical="false" height="4px"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import usePlayerStore from '@/store/player'
import { recordPlayHistory, streamMusicUrl } from '@/api/music'
import { getDeviceType, BREAKPOINTS } from '@/utils/responsive'

// 检测是否为移动设备
const isMobileDevice = /Android|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)

// 检查音频上下文是否已被暂停（移动浏览器常见问题）
function unlockAudioContext() {
  if (isMobileDevice && typeof AudioContext !== 'undefined') {
    // 尝试解锁音频上下文
    document.addEventListener('touchstart', function() {
      if (audioRef.value && audioRef.value.paused) {
        // 预加载音频但不播放
        audioRef.value.load()
      }
    }, { once: true }) // 只监听一次
  }
}

// 初始化音频解锁
unlockAudioContext()

const player = usePlayerStore()
const audioRef = ref(null)
const progress = ref(0)
const currentTime = ref(0)
const duration = ref(0)
const mode = ref(player.state.mode)
const volume = ref(player.state.volume)
const showMoreOptions = ref(false)
const windowWidth = ref(window.innerWidth)

const modeText = computed(() => {
  switch (mode.value) {
    case 'repeat':
      return '循环'
    case 'shuffle':
      return '随机'
    default:
      return '顺序'
  }
})

const isMobileView = computed(() => windowWidth.value <= BREAKPOINTS.mobile)
const buttonSize = computed(() => isMobileView.value ? 'small' : 'default')
const sliderHeight = computed(() => isMobileView.value ? '3px' : '6px')

function toggleMoreOptions() {
  showMoreOptions.value = !showMoreOptions.value
}

function handleResize() {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

const audioSrc = computed(() => {
  return player.state.currentTrack ? streamMusicUrl(player.state.currentTrack.id) : ''
})

function formatTime(value) {
  const minutes = Math.floor(value / 60)
  const seconds = Math.floor(value % 60)
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}

function syncPlayerState() {
  if (audioRef.value) {
    // 确保音频源和音量是最新的
    audioRef.value.src = audioSrc.value
    audioRef.value.volume = volume.value
    
    const isMobile = isMobileDevice
    
    if (player.state.playing) {
      // 确保音频已加载
      audioRef.value.load()
      
      if (isMobile) {
        // 在移动端，直接尝试播放，不使用异步处理
        audioRef.value.play().catch(error => {
          console.error('移动设备播放失败:', error)
          player.togglePlay()
        })
      } else {
        audioRef.value.play().catch(error => {
          console.error('播放失败:', error)
          player.togglePlay()
        })
      }
    } else {
      audioRef.value.pause()
    }
  }
}

async function onPlayToggle() {
  const isMobile = isMobileDevice
  
  // 如果当前没有音频源，但需要播放，则先确保音频已准备
  if (player.state.playing && (!audioRef.value || !audioRef.value.src)) {
    player.togglePlay() // 先恢复状态
    try {
      await readyPlay()
    } catch (error) {
      console.error('播放准备失败:', error)
    }
    return
  }
  
  // 对于移动设备，使用同步方式立即尝试播放
  if (isMobile && audioRef.value && audioRef.value.src) {
    player.togglePlay()
    
    // 立即尝试播放，这是移动端的关键 - 必须在用户手势回调中立即执行
    if (player.state.playing) {
      try {
        await audioRef.value.play()
        await recordHistory()
      } catch (error) {
        console.error('移动设备播放失败:', error)
        player.togglePlay() // 恢复状态
        
        // 有时候需要先加载音频再播放
        audioRef.value.load()
        // 尝试在微任务队列中播放（有时这有助于绕过限制）
        queueMicrotask(async () => {
          try {
            await audioRef.value.play()
            // 更新播放状态
            player.state.playing = true
            await recordHistory()
          } catch (err) {
            console.error('即使在微任务中播放也失败:', err)
          }
        })
      }
    } else {
      audioRef.value.pause()
    }
  } else {
    player.togglePlay()
    
    // 确保音频源是最新的
    if (audioRef.value) {
      audioRef.value.src = audioSrc.value
      audioRef.value.volume = volume.value
    }
    
    // 添加防抖，避免在短时间内多次调用
    if (player.state.playing) {
      // 短暂延迟以确保状态已更新
      await new Promise(resolve => setTimeout(resolve, 50));
      syncPlayerState()
      await recordHistory()
    } else {
      // 如果暂停，直接同步状态
      syncPlayerState()
    }
  }
}

// 用于切换音乐的防抖变量
let switchTrackTimeout = null;

async function onPrev() {
  // 清除之前的切换音乐定时器
  if (switchTrackTimeout) {
    clearTimeout(switchTrackTimeout);
  }
  
  const nextTrack = player.prev()
  if (nextTrack) {
    player.updateTime(0)
    
    // 使用防抖，防止快速点击
    switchTrackTimeout = setTimeout(async () => {
      try {
        // 在切换前暂停当前音频
        if (audioRef.value) {
          audioRef.value.pause();
        }
        
        // 在移动端，减少等待时间，尽快准备播放
        if (isMobileDevice) {
          await new Promise(resolve => setTimeout(resolve, 50))
        } else {
          await new Promise(resolve => setTimeout(resolve, 100))
        }
        await readyPlay()
      } catch (error) {
        console.error('切换上一首歌曲失败:', error)
      } finally {
        switchTrackTimeout = null;
      }
    }, 200); // 200ms 防抖
  }
}

async function onNext() {
  // 清除之前的切换音乐定时器
  if (switchTrackTimeout) {
    clearTimeout(switchTrackTimeout);
  }
  
  const nextTrack = player.next()
  if (nextTrack) {
    player.updateTime(0)
    
    // 使用防抖，防止快速点击
    switchTrackTimeout = setTimeout(async () => {
      try {
        // 在切换前暂停当前音频
        if (audioRef.value) {
          audioRef.value.pause();
        }
        
        // 在移动端，减少等待时间，尽快准备播放
        if (isMobileDevice) {
          await new Promise(resolve => setTimeout(resolve, 50))
        } else {
          await new Promise(resolve => setTimeout(resolve, 100))
        }
        await readyPlay()
      } catch (error) {
        console.error('切换下一首歌曲失败:', error)
      } finally {
        switchTrackTimeout = null;
      }
    }, 200); // 200ms 防抖
  }
}

function onModeChange(value) {
  player.setMode(value)
}

function onVolumeChange(value) {
  player.setVolume(value)
  if (audioRef.value) {
    audioRef.value.volume = value
  }
}

function onSeek(value) {
  if (audioRef.value) {
    audioRef.value.currentTime = value
    player.updateTime(value)
  }
}

async function recordHistory() {
  if (!player.state.currentTrack || !player.state.currentTrack.id) {
    return
  }
  try {
    await recordPlayHistory(player.state.currentTrack.id)
  } catch (error) {
    console.warn('记录播放历史失败', error)
  }
}

// 用于防抖的变量
let pendingReadyPlayPromise = null;

async function readyPlay() {
  if (!audioRef.value) {
    return
  }
  
  // 取消之前的待处理请求
  if (pendingReadyPlayPromise) {
    // 如果有正在进行的 readyPlay，取消它
    console.log('取消之前的播放请求');
    return;
  }
  
  // 创建一个新的 promise 来跟踪这个 readyPlay 调用
  pendingReadyPlayPromise = new Promise((resolve, reject) => {
    setTimeout(() => reject(new Error('readyPlay 超时')), 15000); // 15秒超时
  });
  
  try {
    // 暂停当前播放
    audioRef.value.pause()
    
    // 重置时间和进度
    player.updateTime(0)
    currentTime.value = 0
    progress.value = 0
    
    // 设置音量
    audioRef.value.volume = volume.value
    
    // 设置新的音频源
    audioRef.value.src = audioSrc.value
    
    // 等待一段时间确保新音频源被设置
    await new Promise(resolve => setTimeout(resolve, 100));
    
    // 加载新音频
    audioRef.value.load();
    
    // 等待音频元数据加载完成
    await new Promise((resolve, reject) => {
      const handleLoadedMetadata = () => {
        audioRef.value.removeEventListener('loadedmetadata', handleLoadedMetadata);
        audioRef.value.removeEventListener('error', handleError);
        resolve();
      };
      
      const handleError = (e) => {
        audioRef.value.removeEventListener('loadedmetadata', handleLoadedMetadata);
        audioRef.value.removeEventListener('error', handleError);
        reject(new Error(`音频元数据加载失败: ${e.target.error?.message || '未知错误'}`));
      };
      
      audioRef.value.addEventListener('loadedmetadata', handleLoadedMetadata);
      audioRef.value.addEventListener('error', handleError);
      
      // 设置超时
      setTimeout(() => {
        audioRef.value.removeEventListener('loadedmetadata', handleLoadedMetadata);
        audioRef.value.removeEventListener('error', handleError);
        reject(new Error('音频元数据加载超时'));
      }, 5000);
    });
    
    // 设置时间位置
    if (player.state.currentTime && player.state.currentTime > 0) {
      audioRef.value.currentTime = player.state.currentTime
    }
    
    // 如果应该播放，则开始播放
    if (player.state.playing) {
      await audioRef.value.play()
    }
    
    // 清除待处理的 promise
    pendingReadyPlayPromise = null;
  } catch (error) {
    console.error('准备播放失败:', error)
    pendingReadyPlayPromise = null;
    throw error
  }
}

watch(audioSrc, async (newSrc, oldSrc) => {
  if (audioRef.value && newSrc) {
    // 在移动端，减少等待时间
    if (isMobileDevice) {
      await new Promise(resolve => setTimeout(resolve, 50))
    } else {
      await new Promise(resolve => setTimeout(resolve, 100))
    }
    try {
      await readyPlay()
    } catch (error) {
      console.error('音频源变更后准备播放失败:', error)
    }
  }
})

watch(
  () => player.state.currentTrack,
  () => {
    mode.value = player.state.mode
    volume.value = player.state.volume
    progress.value = 0
    currentTime.value = player.state.currentTime
    duration.value = player.state.duration
  }
)

watch(
  () => player.state.playing,
  () => {
    syncPlayerState()
  }
)

// 存储事件处理函数，以便后续移除
let timeupdateHandler = null
let loadedmetadataHandler = null
let endedHandler = null

onMounted(() => {
  if (audioRef.value) {
    audioRef.value.volume = volume.value

    timeupdateHandler = () => {
      currentTime.value = audioRef.value.currentTime
      progress.value = currentTime.value
      player.updateTime(currentTime.value)
    }
    audioRef.value.addEventListener('timeupdate', timeupdateHandler)

    loadedmetadataHandler = () => {
      duration.value = audioRef.value.duration || 0
      player.updateDuration(duration.value)
    }
    audioRef.value.addEventListener('loadedmetadata', loadedmetadataHandler)

    endedHandler = async () => {
      if (player.state.mode === 'repeat') {
        audioRef.value.currentTime = 0
        audioRef.value.play()
      } else {
        const nextTrack = player.next()
        if (nextTrack) {
          player.updateTime(0)
          try {
            await readyPlay()
          } catch (error) {
            console.error('播放结束后准备下一首失败:', error)
          }
        }
      }
    }
    audioRef.value.addEventListener('ended', endedHandler)

    syncPlayerState()
  }
})

onUnmounted(() => {
  if (audioRef.value) {
    if (timeupdateHandler) {
      audioRef.value.removeEventListener('timeupdate', timeupdateHandler)
    }
    if (loadedmetadataHandler) {
      audioRef.value.removeEventListener('loadedmetadata', loadedmetadataHandler)
    }
    if (endedHandler) {
      audioRef.value.removeEventListener('ended', endedHandler)
    }
  }
})
</script>

<style scoped>
.player-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.96);
  border-top: 1px solid #ebeef5;
  z-index: 1000;
  flex-wrap: wrap;
}
.player-left,
.player-center,
.player-right {
  display: flex;
  align-items: center;
}
.track-info {
  min-width: 200px;
}
.track-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
}
.track-meta {
  font-size: 12px;
  color: #909399;
}
.player-progress {
  width: 420px;
  margin: 0 16px;
}
.player-time {
  font-size: 12px;
  color: #606266;
}
.volume-slider {
  width: 140px;
  margin-left: 16px;
}

/* 移动端适配 */
.player-mobile {
  padding: 8px 10px;
  flex-direction: column;
  align-items: stretch;
  height: auto;
}

.player-mobile .player-left {
  order: 1;
  justify-content: center;
  margin-bottom: 8px;
}

.player-mobile .player-center {
  order: 2;
  flex-direction: column;
  align-items: center;
  width: 100%;
  margin: 8px 0;
}

.player-mobile .player-progress {
  width: 90%;
  margin: 8px 0;
}

.player-mobile .player-time {
  margin-top: 4px;
}

.player-mobile .player-right-mobile {
  order: 3;
  display: flex;
  justify-content: center;
  align-items: center;
}

.player-mobile .track-info {
  min-width: auto;
  text-align: center;
}

.player-mobile .track-title {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.player-mobile .track-meta {
  font-size: 11px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.mobile-options {
  position: absolute;
  bottom: 100%;
  right: 0;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 10px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 1001;
  width: 250px;
}

.mobile-options .mode-selector,
.mobile-options .volume-control {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.mobile-options .mode-selector span,
.mobile-options .volume-control span {
  margin-right: 10px;
  font-size: 12px;
  white-space: nowrap;
}

.mobile-options .volume-slider {
  width: 120px !important;
  margin: 0 !important;
}

/* 平板/车机适配 */
@media (min-width: 769px) and (max-width: 1024px) {
  .player-bar {
    padding: 8px 15px;
  }
  
  .player-progress {
    width: 300px;
  }
  
  .volume-slider {
    width: 100px;
  }
}
</style>