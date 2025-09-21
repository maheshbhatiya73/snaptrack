<template>
  <div class="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow group">
    <div class="p-6">
      <div class="flex items-start justify-between mb-6">
        <div class="flex-1">
          <div class="flex items-center space-x-3 mb-2">
            <div class="w-10 h-10 bg-gray-900 rounded-lg flex items-center justify-center group-hover:scale-105 transition-transform duration-200">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">{{ backup.name }}</h3>
              <p class="text-sm text-gray-600">{{ backup.source }} → {{ backup.destination }}</p>
            </div>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="getStatusClass(backup.status)">
            {{ getStatusText(backup.status) }}
          </span>
          <div class="relative">
            <button @click="toggleMenu" class="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-50 transition-all duration-200">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/>
              </svg>
            </button>
            <div v-if="showMenu" class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-10">
              <div class="py-2">
                <button @click="viewDetails" class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                  View Details
                </button>
                <button @click="editBackup" class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                  </svg>
                  Edit
                </button>
                <button @click="executeBackup" class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  Execute Now
                </button>
                <div class="border-t border-gray-200 my-1"></div>
                <button @click="deleteBackup" class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Schedule</p>
          <p class="text-sm font-semibold text-gray-900">{{ getScheduleText(backup.schedule_type) }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Type</p>
          <p class="text-sm font-semibold text-gray-900 capitalize">{{ backup.type }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">File Type</p>
          <p class="text-sm font-semibold text-gray-900 uppercase">{{ backup.file_type }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Size</p>
          <p class="text-sm font-semibold text-gray-900">{{ formatBytes(backup.size_bytes) }}</p>
        </div>
      </div>

      <div class="border-t border-gray-200 mb-4"></div>

      <div class="flex items-center justify-between text-sm">
        <div class="flex items-center space-x-6">
          <div class="flex items-center space-x-2">
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <span class="text-gray-500">Created: {{ formatDate(backup.created_at) }}</span>
          </div>
          <div v-if="backup.completed_at && backup.completed_at !== '0001-01-01T05:53:28+05:53'" class="flex items-center space-x-2">
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <span class="text-gray-500">Last run: {{ formatDate(backup.completed_at) }}</span>
          </div>
        </div>
        <div v-if="backup.duration_sec > 0" class="flex items-center space-x-2 text-gray-500">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
          <span>{{ formatDuration(backup.duration_sec) }}</span>
        </div>
      </div>

      <!-- Servers -->
      <div v-if="backup.servers && backup.servers.length > 0" class="mt-6 pt-4">
        <div class="flex items-center space-x-2 mb-3">
          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
          </svg>
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Servers</p>
        </div>
        <div class="flex flex-wrap gap-2">
          <span v-for="server in backup.servers" :key="server.id"
                class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-medium bg-gray-100 text-gray-800 border border-gray-200 hover:bg-gray-200 transition-colors duration-200">
            {{ server.name }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  backup: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['edit', 'delete', 'execute'])

const router = useRouter()
const showMenu = ref(false)

const toggleMenu = () => {
  showMenu.value = !showMenu.value
}

const viewDetails = () => {
  showMenu.value = false
  router.push(`/backup/${props.backup.id}`)
}

const editBackup = () => {
  showMenu.value = false
  router.push(`/backup/edit/${props.backup.id}`)
}

const executeBackup = () => {
  showMenu.value = false
  emit('execute', props.backup.id)
}

const deleteBackup = () => {
  showMenu.value = false
  emit('delete', props.backup.id)
}

const getStatusClass = (status) => {
  const statusClasses = {
    'pending': 'bg-gray-100 text-gray-800',
    'running': 'bg-blue-100 text-blue-800',
    'completed': 'bg-green-100 text-green-800',
    'failed': 'bg-red-100 text-red-800',
    'cancelled': 'bg-yellow-100 text-yellow-800'
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

const formatDate = (dateString) => {
  if (!dateString || dateString === '0001-01-01T05:53:28+05:53') return 'Never'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
}

const formatDuration = (seconds) => {
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ${seconds % 60}s`
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}h ${minutes}m`
}

const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
