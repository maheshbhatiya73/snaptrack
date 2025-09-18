<template>
  <div class="min-h-screen bg-slate-50">
    <div class="bg-white border-b border-slate-200">
      <div class="px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-slate-900">Dashboard</h1>
            <p class="mt-1 text-sm text-slate-600">Monitor your backup system status and manage your data</p>
          </div>
          <div class="flex items-center space-x-3">
            <button 
              @click="refreshDashboard"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              {{ loading ? 'Refreshing...' : 'Refresh' }}
            </button>
            <button 
              @click="toggleAutoRefresh"
              :class="[
                'inline-flex items-center px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200',
                autoRefreshEnabled 
                  ? 'bg-green-600 text-white hover:bg-green-700' 
                  : 'bg-slate-200 text-slate-700 hover:bg-slate-300'
              ]"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              {{ autoRefreshEnabled ? 'Auto-refresh ON' : 'Auto-refresh OFF' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <main class="px-6 py-6">
      <div v-if="loading && !stats" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-200 border-t-blue-600"></div>
          <p class="mt-6 text-lg text-slate-600 font-medium">Loading dashboard...</p>
        </div>
      </div>

      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6 mb-8">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-lg font-semibold text-red-800">Error loading dashboard</h3>
            <p class="text-red-700 mt-2">{{ error }}</p>
            <button 
              @click="refreshDashboard"
              class="mt-4 inline-flex items-center px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>

      <div v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 rounded-lg" :class="systemStatus === 'online' ? 'bg-green-100' : 'bg-red-100'">
                <svg class="w-6 h-6" :class="systemStatus === 'online' ? 'text-green-600' : 'text-red-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">System Status</p>
                <p class="text-2xl font-bold" :class="systemStatus === 'online' ? 'text-green-600' : 'text-red-600'">
                  {{ systemStatus === 'online' ? 'Online' : 'Offline' }}
                </p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-blue-100 rounded-lg">
                <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Total Backups</p>
                <p class="text-2xl font-bold text-blue-600">{{ formatNumber(stats?.total_backups || 0) }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 rounded-lg">
                <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Storage Used</p>
                <p class="text-2xl font-bold text-purple-600">{{ formatStorage(stats?.storage_used || 0) }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-amber-100 rounded-lg">
                <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Last Backup</p>
                <p class="text-2xl font-bold text-amber-600">{{ formatLastBackup(stats?.last_backup) }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-slate-200">
          <div class="px-6 py-4 border-b border-slate-200">
            <h2 class="text-lg font-semibold text-slate-900">Recent Activity</h2>
          </div>
          <div class="p-6">
            <div v-if="recentActivity.length === 0" class="text-center py-8 text-slate-500">
              <svg class="w-12 h-12 mx-auto mb-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
              </svg>
              <p>No recent activity</p>
            </div>
            <div v-else class="space-y-4">
              <div 
                v-for="activity in recentActivity" 
                :key="activity.id"
                class="flex items-center space-x-3"
              >
                <div 
                  class="w-2 h-2 rounded-full"
                  :class="getActivityColor(activity.type)"
                ></div>
                <div class="flex-1">
                  <p class="text-sm text-slate-900">{{ activity.message }}</p>
                  <p class="text-xs text-slate-500">{{ formatTimeAgo(activity.timestamp) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  isAuthenticated, 
  fetchDashboardStats, 
  fetchRecentActivity, 
  fetchSystemStatus 
} from '~/lib/api'

const router = useRouter()

const stats = ref(null)
const recentActivity = ref([])
const systemStatus = ref('offline')
const loading = ref(false)
const error = ref(null)
const autoRefreshEnabled = ref(true)
const autoRefreshInterval = ref(null)

const loadDashboardData = async () => {
  try {
    loading.value = true
    error.value = null
    
    const [statsData, activityData, statusData] = await Promise.all([
      fetchDashboardStats().catch(() => ({ total_backups: 0, storage_used: 0, last_backup: null })),
      fetchRecentActivity().catch(() => []),
      fetchSystemStatus().catch(() => ({ status: 'offline' }))
    ])
    
    stats.value = statsData
    recentActivity.value = activityData.activities || activityData || []
    systemStatus.value = statusData.status || 'offline'
  } catch (err) {
    error.value = err.message
    console.error('Failed to load dashboard data:', err)
  } finally {
    loading.value = false
  }
}

const refreshDashboard = () => {
  loadDashboardData()
}

const startAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
  }
  
  if (autoRefreshEnabled.value) {
    autoRefreshInterval.value = setInterval(async () => {
      if (!loading.value && document.visibilityState === 'visible') {
        await loadDashboardData()
      }
    }, 30000)
  }
}

const stopAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
    autoRefreshInterval.value = null
  }
}

const toggleAutoRefresh = () => {
  autoRefreshEnabled.value = !autoRefreshEnabled.value
  if (autoRefreshEnabled.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatStorage = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatLastBackup = (timestamp) => {
  if (!timestamp) return 'Never'
  
  const now = new Date()
  const backupTime = new Date(timestamp)
  const diffMs = now - backupTime
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`
  
  return backupTime.toLocaleDateString()
}

const formatTimeAgo = (timestamp) => {
  if (!timestamp) return 'Unknown'
  
  const now = new Date()
  const activityTime = new Date(timestamp)
  const diffMs = now - activityTime
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins} minutes ago`
  if (diffHours < 24) return `${diffHours} hours ago`
  if (diffDays < 7) return `${diffDays} days ago`
  
  return activityTime.toLocaleDateString()
}

const getActivityColor = (type) => {
  const colors = {
    success: 'bg-green-500',
    info: 'bg-blue-500',
    warning: 'bg-amber-500',
    error: 'bg-red-500',
    backup: 'bg-purple-500',
    system: 'bg-slate-500'
  }
  return colors[type] || 'bg-slate-500'
}

onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  loadDashboardData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>