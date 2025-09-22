<template>
  <div class="bg-white border border-gray-200 rounded-lg shadow-sm overflow-hidden">
    <div class="p-6">
      <div class="flex items-start justify-between mb-4">
        <div class="flex-1">
          <div class="flex items-center space-x-3 mb-2">
            <div class="w-10 h-10 bg-gray-900 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">{{ process.backup?.name || `Backup #${process.backup_id}` }}</h3>
              <p class="text-sm text-gray-600">{{ process.backup?.source || 'Unknown source' }} → {{ process.backup?.destination || 'Unknown destination' }}</p>
            </div>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <button
            @click="$emit('remove', process.id)"
            class="inline-flex items-center px-2 py-1 border border-red-300 rounded-md text-xs font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          >
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            Remove
          </button>
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="getStatusClass(process.status)">
            {{ getStatusText(process.status) }}
          </span>
        </div>
      </div>

      <!-- Progress Bar -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-700">Progress</span>
          <span class="text-sm text-gray-500">{{ process.progress }}%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="h-2 rounded-full transition-all duration-300 ease-out"
            :class="getProgressBarClass(process.status)"
            :style="{ width: process.progress + '%' }"
          ></div>
        </div>
      </div>

      <!-- Process Details -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Type</p>
          <p class="text-sm font-semibold text-gray-900 capitalize">{{ process.backup?.type || 'Unknown' }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">File Type</p>
          <p class="text-sm font-semibold text-gray-900 uppercase">{{ process.backup?.file_type || 'Unknown' }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Size</p>
          <p class="text-sm font-semibold text-gray-900">{{ formatBytes(process.bytes_processed || 0) }}</p>
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Speed</p>
          <p class="text-sm font-semibold text-gray-900">{{ formatSpeed(process.speed_bps) }}</p>
        </div>
      </div>

      <!-- Current File -->
      <div v-if="process.current_file" class="mb-4">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">Current File</p>
        <p class="text-sm text-gray-900 truncate">{{ process.current_file }}</p>
      </div>

      <!-- Message -->
      <div class="mb-4">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">Status</p>
        <p class="text-sm text-gray-900">{{ process.message }}</p>
      </div>

      <!-- ETA -->
      <div v-if="process.eta_seconds && process.status === 'running'" class="mb-4">
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">Estimated Time Remaining</p>
        <p class="text-sm text-gray-900">{{ formatDuration(process.eta_seconds) }}</p>
      </div>

      <!-- Servers -->
      <div v-if="process.backup?.servers && process.backup.servers.length > 0" class="border-t border-gray-200 pt-4">
        <div class="flex items-center space-x-2 mb-3">
          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
          </svg>
          <p class="text-xs font-medium text-gray-500 uppercase tracking-wide">Target Servers</p>
        </div>
        <div class="flex flex-wrap gap-2">
          <span v-for="server in process.backup.servers" :key="server.id"
                class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-medium bg-gray-100 text-gray-800 border border-gray-200">
            <span v-if="server.type === 'local'" class="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
            <span v-else class="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
            {{ server.name }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  process: {
    type: Object,
    required: true
  }
})

const getStatusClass = (status) => {
  const statusClasses = {
    'pending': 'bg-gray-100 text-gray-800',
    'running': 'bg-blue-100 text-blue-800',
    'completed': 'bg-green-100 text-green-800',
    'failed': 'bg-red-100 text-red-800',
    'error': 'bg-red-100 text-red-800'
  }
  return statusClasses[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const statusTexts = {
    'pending': 'Pending',
    'running': 'Running',
    'completed': 'Completed',
    'failed': 'Failed',
    'error': 'Error'
  }
  return statusTexts[status] || status
}

const getProgressBarClass = (status) => {
  const classes = {
    'pending': 'bg-gray-400',
    'running': 'bg-blue-500',
    'completed': 'bg-green-500',
    'failed': 'bg-red-500',
    'error': 'bg-red-500'
  }
  return classes[status] || 'bg-gray-400'
}

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatSpeed = (bps) => {
  if (!bps || bps === 0) return 'N/A'
  return formatBytes(bps) + '/s'
}

const formatDuration = (seconds) => {
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ${seconds % 60}s`
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}h ${minutes}m`
}
</script>