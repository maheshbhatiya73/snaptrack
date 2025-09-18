<template>
  <div class="min-h-screen bg-slate-50">
    <!-- Page Header -->
    <div class="bg-white border-b border-slate-200">
      <div class="px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button 
              @click="goBack"
              class="p-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <div>
              <h1 class="text-2xl font-bold text-slate-900">{{ backup?.name || 'Backup Details' }}</h1>
              <p class="mt-1 text-sm text-slate-600">View backup information and execution history</p>
            </div>
          </div>
          <div class="flex items-center space-x-3">
            <button 
              @click="editBackup"
              class="inline-flex items-center px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 transition-colors duration-200"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
              Edit
            </button>
            <button 
              @click="executeBackup"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors duration-200 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h1m4 0h1m-6-8h8a2 2 0 012 2v8a2 2 0 01-2 2H8a2 2 0 01-2-2V6a2 2 0 012-2z"/>
              </svg>
              {{ loading ? 'Executing...' : 'Execute Now' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <main class="px-6 py-6">
      <div class="max-w-6xl mx-auto">
        <!-- Loading State -->
        <div v-if="loading && !backup" class="flex items-center justify-center py-12">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p class="mt-4 text-slate-600">Loading backup details...</p>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6 mb-6">
          <div class="flex items-center">
            <svg class="w-5 h-5 text-red-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div>
              <h3 class="text-sm font-medium text-red-800">Error loading backup</h3>
              <p class="text-sm text-red-700 mt-1">{{ error }}</p>
            </div>
          </div>
        </div>

        <!-- Success Message -->
        <div v-if="success" class="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
          <div class="flex items-center">
            <svg class="w-5 h-5 text-green-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div>
              <h3 class="text-sm font-medium text-green-800">Success</h3>
              <p class="text-sm text-green-700 mt-1">{{ success }}</p>
            </div>
          </div>
        </div>

        <!-- Backup Details -->
        <div v-if="backup" class="space-y-6">
          <!-- Status and Basic Info -->
          <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            <div class="flex items-start justify-between mb-6">
              <div>
                <h2 class="text-xl font-semibold text-slate-900 mb-2">{{ backup.name }}</h2>
                <p class="text-slate-600">{{ backup.source }} → {{ backup.destination }}</p>
              </div>
              <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium"
                    :class="getStatusClass(backup.status)">
                {{ getStatusText(backup.status) }}
              </span>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Backup Type</p>
                <p class="text-lg font-semibold text-slate-900 capitalize">{{ backup.type }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">File Type</p>
                <p class="text-lg font-semibold text-slate-900 uppercase">{{ backup.file_type }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Schedule</p>
                <p class="text-lg font-semibold text-slate-900 capitalize">{{ getScheduleText(backup.schedule_type) }}</p>
              </div>
            </div>
          </div>

          <!-- Servers -->
          <div v-if="backup.servers && backup.servers.length > 0" class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            <h3 class="text-lg font-semibold text-slate-900 mb-4">Target Servers</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div v-for="server in backup.servers" :key="server.id" class="border border-slate-200 rounded-lg p-4">
                <div class="flex items-center justify-between mb-2">
                  <h4 class="font-medium text-slate-900">{{ server.name }}</h4>
                  <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                        :class="server.enabled ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                    {{ server.enabled ? 'Enabled' : 'Disabled' }}
                  </span>
                </div>
                <p class="text-sm text-slate-600">{{ server.host }}:{{ server.ssh_port }}</p>
                <p class="text-xs text-slate-500 mt-1">User: {{ server.ssh_user }}</p>
              </div>
            </div>
          </div>

          <!-- Execution Details -->
          <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            <h3 class="text-lg font-semibold text-slate-900 mb-4">Execution Details</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Size</p>
                <p class="text-lg font-semibold text-slate-900">{{ formatBytes(backup.size_bytes) }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Duration</p>
                <p class="text-lg font-semibold text-slate-900">{{ formatDuration(backup.duration_sec) }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Started At</p>
                <p class="text-lg font-semibold text-slate-900">{{ formatDate(backup.started_at) }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Completed At</p>
                <p class="text-lg font-semibold text-slate-900">{{ formatDate(backup.completed_at) }}</p>
              </div>
            </div>
            
            <div v-if="backup.checksum" class="mt-6 pt-6 border-t border-slate-200">
              <p class="text-sm font-medium text-slate-500 uppercase tracking-wide mb-2">Checksum</p>
              <p class="text-sm font-mono text-slate-900 break-all">{{ backup.checksum }}</p>
            </div>
          </div>

          <!-- Metadata -->
          <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            <h3 class="text-lg font-semibold text-slate-900 mb-4">Metadata</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Created At</p>
                <p class="text-sm text-slate-900">{{ formatDate(backup.created_at) }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Updated At</p>
                <p class="text-sm text-slate-900">{{ formatDate(backup.updated_at) }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Executed By</p>
                <p class="text-sm text-slate-900">{{ backup.executed_by || 'N/A' }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-500 uppercase tracking-wide">Backup ID</p>
                <p class="text-sm font-mono text-slate-900">{{ backup.id }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { isAuthenticated, fetchBackup, executeBackup } from '~/lib/api'

const router = useRouter()
const route = useRoute()

// Reactive data
const backup = ref(null)
const loading = ref(false)
const error = ref(null)
const success = ref(null)

// Methods
const goBack = () => {
  router.push('/backup')
}

const editBackup = () => {
  router.push(`/backup/edit/${route.params.id}`)
}

const loadBackup = async () => {
  try {
    loading.value = true
    error.value = null
    
    const backupId = route.params.id
    backup.value = await fetchBackup(backupId)
    
  } catch (err) {
    error.value = err.message
    console.error('Failed to load backup:', err)
  } finally {
    loading.value = false
  }
}

const executeBackupHandler = async () => {
  try {
    loading.value = true
    error.value = null
    success.value = null

    const backupId = route.params.id
    await executeBackup(backupId)
    
    success.value = 'Backup execution started!'
    
    // Refresh backup details
    await loadBackup()
    
  } catch (err) {
    error.value = err.message
    console.error('Failed to execute backup:', err)
  } finally {
    loading.value = false
  }
}

const getStatusClass = (status) => {
  const statusClasses = {
    'pending': 'bg-yellow-100 text-yellow-800',
    'running': 'bg-blue-100 text-blue-800',
    'completed': 'bg-green-100 text-green-800',
    'failed': 'bg-red-100 text-red-800',
    'cancelled': 'bg-gray-100 text-gray-800'
  }
  return statusClasses[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const statusTexts = {
    'pending': 'Pending',
    'running': 'Running',
    'completed': 'Completed',
    'failed': 'Failed',
    'cancelled': 'Cancelled'
  }
  return statusTexts[status] || status
}

const getScheduleText = (scheduleType) => {
  const scheduleTexts = {
    'one_time': 'One Time',
    'daily': 'Daily',
    'weekly': 'Weekly',
    'monthly': 'Monthly'
  }
  return scheduleTexts[scheduleType] || scheduleType
}

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDuration = (seconds) => {
  if (!seconds || seconds === 0) return 'N/A'
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ${seconds % 60}s`
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}h ${minutes}m`
}

const formatDate = (dateString) => {
  if (!dateString || dateString === '0001-01-01T05:53:28+05:53') return 'Never'
  const date = new Date(dateString)
  return date.toLocaleString()
}

// Check authentication and load backup
onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  loadBackup()
})
</script>
