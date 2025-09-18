<template>
  <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
    <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
      <div class="flex items-center">
        <div class="p-2 bg-green-100 rounded-lg">
          <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-slate-600">Successful Backups</p>
          <p class="text-2xl font-bold text-green-600">{{ stats.successful }}</p>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
      <div class="flex items-center">
        <div class="p-2 bg-red-100 rounded-lg">
          <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"/>
          </svg>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-slate-600">Failed Backups</p>
          <p class="text-2xl font-bold text-red-600">{{ stats.failed }}</p>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
      <div class="flex items-center">
        <div class="p-2 bg-blue-100 rounded-lg">
          <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
          </svg>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-slate-600">Total Storage</p>
          <p class="text-2xl font-bold text-blue-600">{{ formatBytes(stats.totalSize) }}</p>
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
          <p class="text-2xl font-bold text-amber-600">{{ formatLastBackup(stats.lastBackup) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  backups: {
    type: Array,
    default: () => []
  }
})

const stats = computed(() => {
  const successful = props.backups.filter(backup => backup.status === 'completed').length
  const failed = props.backups.filter(backup => backup.status === 'failed').length
  const totalSize = props.backups.reduce((sum, backup) => sum + (backup.size_bytes || 0), 0)
  
  const completedBackups = props.backups.filter(backup => 
    backup.status === 'completed' && 
    backup.completed_at && 
    backup.completed_at !== '0001-01-01T05:53:28+05:53'
  )
  
  const lastBackup = completedBackups.length > 0 
    ? completedBackups.sort((a, b) => new Date(b.completed_at) - new Date(a.completed_at))[0]
    : null

  return {
    successful,
    failed,
    totalSize,
    lastBackup
  }
})

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatLastBackup = (lastBackup) => {
  if (!lastBackup) return 'Never'
  
  const now = new Date()
  const backupTime = new Date(lastBackup.completed_at)
  const diffMs = now - backupTime
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)
  
  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`
  
  return backupTime.toLocaleDateString()
}
</script>
