<template>
  <div class="home-container">
    <el-card shadow="hover" class="music-card">
      <div class="tools-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索歌曲、歌手、专辑"
          clearable
          @clear="resetSearch"
          @keyup.enter="onSearch"
          :style="inputStyle"
        />
        <el-button type="primary" @click="onSearch">搜索</el-button>
        <el-button type="success" @click="scanMusicLibrary" :loading="scanning">扫描音乐</el-button>
      </div>
      <el-table
        :data="musicList"
        v-loading="loading"
        stripe
        :style="{ width: '100%', maxHeight: tableMaxHeight }"
        :max-height="tableMaxHeightNum"
      >
        <el-table-column prop="title" label="歌曲" :min-width="columnWidth.title" />
        <el-table-column prop="artist" label="歌手" :width="columnWidth.artist" v-if="!isMobileView" />
        <el-table-column prop="album" label="专辑" :width="columnWidth.album" v-if="!isMobileView" />
        <el-table-column prop="duration" label="时长" :width="columnWidth.duration">
          <template #default="{ row }">
            {{ formatDuration(row.duration) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="columnWidth.action">
          <template #default="{ row, $index }">
            <el-button type="primary" size="small" @click="play(row, $index)">播放</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-wrapper" :class="{ 'mobile-pagination': isMobileView }">
        <el-pagination
          background
          :layout="paginationLayout"
          :page-size="pageSize"
          :current-page="page"
          :total="total"
          @current-change="onPageChange"
        />
      </div>
    </el-card>
    <el-card shadow="hover" class="recent-card" v-loading="recentLoading">
      <div class="section-title">最近播放</div>
      <el-table :data="recentPlays" stripe :style="{ width: '100%', maxHeight: recentTableMaxHeight }" :max-height="recentTableMaxHeightNum">
        <el-table-column label="歌曲" :min-width="recentColumnWidth.title">
          <template #default="{ row }">
            {{ row.music && row.music.title ? row.music.title : '未知歌曲' }}
          </template>
        </el-table-column>
        <el-table-column label="歌手" :width="recentColumnWidth.artist" v-if="!isMobileView">
          <template #default="{ row }">
            {{ row.music && row.music.artist ? row.music.artist : '未知歌手' }}
          </template>
        </el-table-column>
        <el-table-column label="专辑" :width="recentColumnWidth.album" v-if="!isMobileView">
          <template #default="{ row }">
            {{ row.music && row.music.album ? row.music.album : '未知专辑' }}
          </template>
        </el-table-column>
        <el-table-column label="播放时间" :width="recentColumnWidth.duration" v-if="!isMobileView">
          <template #default="{ row }">
            {{ formatPlayedAt(row.played_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="recentColumnWidth.action">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="playRecent(row)">播放</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { fetchMusicList, fetchRecentPlays, recordPlayHistory, searchMusic, scanMusic } from '@/api/music'
import usePlayerStore from '@/store/player'
import { getDeviceType, BREAKPOINTS } from '@/utils/responsive'

const musicList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const loading = ref(false)
const keyword = ref('')
const recentPlays = ref([])
const recentLoading = ref(false)

const player = usePlayerStore()

const scanning = ref(false)
const windowWidth = ref(window.innerWidth)

// 计算属性
const deviceType = computed(() => getDeviceType())
const isMobileView = computed(() => windowWidth.value <= BREAKPOINTS.mobile)

const inputStyle = computed(() => {
  return isMobileView.value 
    ? { width: '100%', marginBottom: '10px' } 
    : { width: '320px' }
})

const columnWidth = computed(() => {
  if (isMobileView.value) {
    return {
      title: 120,
      artist: 0,
      album: 0,
      duration: 80,
      action: 80
    }
  } else {
    return {
      title: 220,
      artist: 140,
      album: 140,
      duration: 100,
      action: 140
    }
  }
})

const recentColumnWidth = computed(() => {
  if (isMobileView.value) {
    return {
      title: 120,
      artist: 0,
      album: 0,
      duration: 0,
      action: 80
    }
  } else {
    return {
      title: 220,
      artist: 140,
      album: 140,
      duration: 180,
      action: 120
    }
  }
})

const paginationLayout = computed(() => {
  return isMobileView.value ? 'prev, next' : 'prev, pager, next'
})

const tableMaxHeight = computed(() => {
  return isMobileView.value ? '40vh' : '60vh'
})

const tableMaxHeightNum = computed(() => {
  return isMobileView.value ? 40 * window.innerHeight / 100 : 60 * window.innerHeight / 100
})

const recentTableMaxHeight = computed(() => {
  return isMobileView.value ? '30vh' : '40vh'
})

const recentTableMaxHeightNum = computed(() => {
  return isMobileView.value ? 30 * window.innerHeight / 100 : 40 * window.innerHeight / 100
})

function handleResize() {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadMusicList()
  loadRecentPlays()
})

function formatDuration(seconds) {
  if (!seconds || seconds <= 0) {
    return '00:00'
  }
  const min = Math.floor(seconds / 60)
  const sec = Math.floor(seconds % 60)
  return `${String(min).padStart(2, '0')}:${String(sec).padStart(2, '0')}`
}

function formatPlayedAt(timestamp) {
  if (!timestamp) {
    return '-'
  }
  const date = new Date(timestamp)
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `${month}-${day} ${hour}:${minute}`
}

async function loadMusicList() {
  loading.value = true
  try {
    const res = await fetchMusicList(page.value, pageSize.value)
    musicList.value = res.data.musics || []
    total.value = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载歌曲列表失败')
  } finally {
    loading.value = false
  }
}

async function loadRecentPlays() {
  recentLoading.value = true
  try {
    const res = await fetchRecentPlays(1, 10)
    recentPlays.value = res.data.records || []
  } catch (error) {
    ElMessage.error('加载最近播放失败')
  } finally {
    recentLoading.value = false
  }
}

async function onSearch() {
  if (!keyword.value.trim()) {
    page.value = 1
    await loadMusicList()
    return
  }

  loading.value = true
  try {
    const res = await searchMusic(keyword.value.trim(), page.value, pageSize.value)
    musicList.value = res.data.musics || []
    total.value = res.data.total || 0
  } catch (error) {
    ElMessage.error('搜索歌曲失败')
  } finally {
    loading.value = false
  }
}

function resetSearch() {
  keyword.value = ''
  page.value = 1
  loadMusicList()
}

function onPageChange(newPage) {
  page.value = newPage
  if (keyword.value.trim()) {
    onSearch()
  } else {
    loadMusicList()
  }
}

async function play(track, index) {
  player.setPlaylist(musicList.value, index)
  try {
    await recordPlayHistory(track.id)
  } catch (error) {
    console.warn('记录播放历史失败', error)
  }
  ElMessage.success(`正在播放：${track.title || '未知歌曲'}`)
  loadRecentPlays()
}

async function playRecent(record) {
  const track = record.music
  if (!track) {
    ElMessage.warning('最近播放数据不可用')
    return
  }
  player.setPlaylist([track], 0)
  try {
    await recordPlayHistory(track.id)
  } catch (error) {
    console.warn('记录播放历史失败', error)
  }
  ElMessage.success(`播放最近歌曲：${track.title || '未知歌曲'}`)
}

async function scanMusicLibrary() {
  try {
    await ElMessageBox.confirm(
      '确定要扫描音乐目录吗？这将搜索您的音乐目录并导入所有支持的音频文件。',
      '扫描音乐',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    scanning.value = true
    const res = await scanMusic()
    
    if (res.code === 200) {
      ElMessage.success(res.data.message || '音乐扫描已启动')
      // 扫描是异步的，稍等片刻后刷新列表
      setTimeout(() => {
        loadMusicList()
        loadRecentPlays()
      }, 3000) // 3秒后刷新列表
    } else {
      ElMessage.error(res.message || '扫描请求失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('扫描音乐失败:', error)
      ElMessage.error('扫描音乐失败')
    }
  } finally {
    scanning.value = false
  }
}

onMounted(() => {
  loadMusicList()
  loadRecentPlays()
})
</script>

<style scoped>
.home-container {
  padding: 20px;
}
.music-card {
  min-height: calc(100vh - 180px);
}
.tools-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.pagination-wrapper {
  margin-top: 16px;
  text-align: right;
}
.mobile-pagination {
  text-align: center !important;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .home-container {
    padding: 10px;
  }
  
  .music-card {
    min-height: calc(100vh - 160px);
  }
  
  .tools-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .tools-bar > * {
    margin-bottom: 10px;
  }
  
  .pagination-wrapper.mobile-pagination {
    text-align: center;
  }
  
  .el-table .cell {
    padding: 0 5px !important;
  }
}

/* 平板/车机适配 */
@media (min-width: 769px) and (max-width: 1024px) {
  .home-container {
    padding: 15px;
  }
  
  .tools-bar {
    gap: 8px;
  }
  
  .el-table .cell {
    padding: 0 6px !important;
  }
}
</style>