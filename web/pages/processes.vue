<template>
  <div class="min-h-screen bg-white">
    <!-- Header Section -->
    <div class="bg-white border-b border-gray-200">
      <div class="px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Backup Processes</h1>
            <p class="mt-2 text-gray-600">Monitor running backup operations in real-time</p>
          </div>
          <div class="flex space-x-3">
            <button
              @click="removeAllProcesses"
              :disabled="loading || processes.length === 0"
              class="inline-flex items-center px-4 py-2 border border-red-300 rounded-md shadow-sm text-sm font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
              Remove All
            </button>
            <button
              @click="refreshProcesses"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              {{ loading ? 'Refreshing...' : 'Refresh' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <main class="px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading && processes.length === 0" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
          <p class="mt-6 text-lg text-gray-600 font-medium">Loading backup processes...</p>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-white border border-red-200 rounded-lg p-6 mb-8">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-lg font-semibold text-red-800">Error loading processes</h3>
            <p class="text-gray-600 mt-2">{{ error }}</p>
            <button
              @click="refreshProcesses"
              class="mt-4 inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div v-else>
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div class="bg-white border border-gray-200 rounded-lg p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                  <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Running</p>
                <p class="text-2xl font-bold text-gray-900">{{ stats.running }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white border border-gray-200 rounded-lg p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-green-100 rounded-lg flex items-center justify-center">
                  <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Completed</p>
                <p class="text-2xl font-bold text-gray-900">{{ stats.completed }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white border border-gray-200 rounded-lg p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-red-100 rounded-lg flex items-center justify-center">
                  <svg class="w-4 h-4 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Failed</p>
                <p class="text-2xl font-bold text-gray-900">{{ stats.failed }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white border border-gray-200 rounded-lg p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-yellow-100 rounded-lg flex items-center justify-center">
                  <svg class="w-4 h-4 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Pending</p>
                <p class="text-2xl font-bold text-gray-900">{{ stats.pending }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="processes.length === 0" class="text-center py-20">
          <div class="w-24 h-24 bg-gray-100 rounded-lg flex items-center justify-center mx-auto mb-6">
            <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-3">No backup processes running</h3>
          <p class="text-gray-500 mb-8 max-w-md mx-auto">Start a backup to see real-time progress monitoring here</p>
          <NuxtLink
            to="/backup"
            class="inline-flex items-center px-6 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Start Backup
          </NuxtLink>
        </div>

        <!-- Process List -->
        <div v-else class="space-y-6">
          <BackupProcessCard
            v-for="process in processes"
            :key="process.id"
            :process="process"
            @remove="removeProcess"
          />
        </div>
      </div>
    </main>

    <!-- Confirmation Modal -->
    <ConfirmationModal
      v-if="showConfirmModal"
      :message="confirmMessage"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useRuntimeConfig } from '#app'
import { isAuthenticated, fetchRunningBackups, deleteAllProcesses, deleteProcess } from '~/lib/api'
import BackupProcessCard from '~/components/BackupProcessCard.vue'
import ConfirmationModal from '~/components/ConfirmationModal.vue'

const router = useRouter()
const config = useRuntimeConfig()
const processes = ref([])
const loading = ref(false)
const error = ref(null)
let wsConnection = null
const showConfirmModal = ref(false)
const confirmMessage = ref('')
const confirmAction = ref(null)

const stats = computed(() => {
  const stats = {
    running: 0,
    completed: 0,
    failed: 0,
    pending: 0
  }

  processes.value.forEach(process => {
    if (process.status === 'running') stats.running++
    else if (process.status === 'completed') stats.completed++
    else if (process.status === 'failed') stats.failed++
    else if (process.status === 'pending') stats.pending++
  })

  return stats
})

const loadProcesses = async () => {
  try {
    loading.value = true
    error.value = null
    processes.value = await fetchRunningBackups() || []
  } catch (err) {
    error.value = err.message
    console.error('Failed to load processes:', err)
  } finally {
    loading.value = false
  }
}

const refreshProcesses = () => {
  loadProcesses()
}

const removeAllProcesses = () => {
  confirmMessage.value = 'Are you sure you want to remove all backup processes? This action cannot be undone.'
  confirmAction.value = async () => {
    try {
      await deleteAllProcesses()
      processes.value = []
      showConfirmModal.value = false
    } catch (err) {
      console.error('Failed to remove all processes:', err)
      alert('Failed to remove all processes: ' + err.message)
      showConfirmModal.value = false
    }
  }
  showConfirmModal.value = true
}

const removeProcess = (processId) => {
  confirmMessage.value = 'Are you sure you want to remove this backup process?'
  confirmAction.value = async () => {
    try {
      await deleteProcess(processId)
      const index = processes.value.findIndex(p => p.id === processId)
      if (index >= 0) {
        processes.value.splice(index, 1)
      }
      showConfirmModal.value = false
    } catch (err) {
      console.error('Failed to remove process:', err)
      alert('Failed to remove process: ' + err.message)
      showConfirmModal.value = false
    }
  }
  showConfirmModal.value = true
}

const confirmDelete = () => {
  if (confirmAction.value) {
    confirmAction.value()
  }
}

const cancelDelete = () => {
  showConfirmModal.value = false
  confirmAction.value = null
}

const connectWebSocket = () => {
  try {
    const backendUrl = config.public.backendUrl
    const wsUrl = backendUrl.replace(/^http/, 'ws') + '/ws/backups'

    wsConnection = new WebSocket(wsUrl)

    wsConnection.onopen = () => {
      console.log('WebSocket connected for backup processes')
    }

    wsConnection.onmessage = (event) => {
      try {
        const progressData = JSON.parse(event.data)
        updateProcessProgress(progressData)
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err)
      }
    }

    wsConnection.onclose = () => {
      console.log('WebSocket disconnected')
      // Attempt to reconnect after a delay
      setTimeout(connectWebSocket, 5000)
    }

    wsConnection.onerror = (err) => {
      console.error('WebSocket error:', err)
    }
  } catch (err) {
    console.error('Failed to connect WebSocket:', err)
  }
}

const updateProcessProgress = (progressData) => {
  const existingIndex = processes.value.findIndex(p => p.backup_id === progressData.backup_id)

  if (existingIndex >= 0) {
    // Update existing process
    processes.value[existingIndex] = {
      ...processes.value[existingIndex],
      ...progressData,
      backup: {
        ...progressData.backup,
        servers: progressData.servers
      }
    }
  } else {
    // Add new process
    processes.value.unshift({
      ...progressData,
      backup: {
        ...progressData.backup,
        servers: progressData.servers
      }
    })
  }

  // Remove completed processes after a delay, keep failed for user to see
  if (progressData.status === 'completed') {
    setTimeout(() => {
      const index = processes.value.findIndex(p => p.id === progressData.id)
      if (index >= 0) {
        processes.value.splice(index, 1)
      }
    }, 10000) // Remove after 10 seconds
  }
}

onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }

  loadProcesses()
  connectWebSocket()
})

onUnmounted(() => {
  if (wsConnection) {
    wsConnection.close()
  }
})
</script>