<template>
  <div class="bg-white rounded-xl shadow-sm border border-slate-200 hover:shadow-md transition-shadow duration-200">
    <div class="p-6">
      <div class="flex items-start justify-between mb-4">
        <div class="flex-1">
          <h3 class="text-lg font-semibold text-slate-900 mb-1">{{ backup.name }}</h3>
          <p class="text-sm text-slate-600">{{ backup.source }} → {{ backup.destination }}</p>
        </div>
        <div class="flex items-center space-x-2">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="getStatusClass(backup.status)">
            {{ getStatusText(backup.status) }}
          </span>
          <div class="relative">
            <button @click="toggleMenu" class="p-1 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"/>
              </svg>
            </button>
            <div v-if="showMenu" class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border border-slate-200">
              <div class="py-1">
                <button @click="viewDetails" class="block w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-100">
                  View Details
                </button>
                <button @click="editBackup" class="block w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-100">
                  Edit
                </button>
                <button @click="executeBackup" class="block w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-100">
                  Execute Now
                </button>
                <button @click="deleteBackup" class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50">
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">Schedule</p>
          <p class="text-sm text-slate-900">{{ getScheduleText(backup.schedule_type) }}</p>
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">Type</p>
          <p class="text-sm text-slate-900 capitalize">{{ backup.type }}</p>
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">File Type</p>
          <p class="text-sm text-slate-900 uppercase">{{ backup.file_type }}</p>
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">Size</p>
          <p class="text-sm text-slate-900">{{ formatBytes(backup.size_bytes) }}</p>
        </div>
      </div>

      <div class="flex items-center justify-between text-sm text-slate-500">
        <div class="flex items-center space-x-4">
          <span>Created: {{ formatDate(backup.created_at) }}</span>
          <span v-if="backup.completed_at && backup.completed_at !== '0001-01-01T05:53:28+05:53'">
            Last run: {{ formatDate(backup.completed_at) }}
          </span>
        </div>
        <div v-if="backup.duration_sec > 0" class="text-slate-400">
          {{ formatDuration(backup.duration_sec) }}
        </div>
      </div>

      <!-- Servers -->
      <div v-if="backup.servers && backup.servers.length > 0" class="mt-4 pt-4 border-t border-slate-200">
        <p class="text-xs font-medium text-slate-500 uppercase tracking-wide mb-2">Servers</p>
        <div class="flex flex-wrap gap-2">
          <span v-for="server in backup.servers" :key="server.id" 
                class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-blue-100 text-blue-800">
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
