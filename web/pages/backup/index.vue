<template>
  <div class="min-h-screen bg-white">
    <!-- Header Section -->
    <div class="bg-white border-b border-gray-200">
      <div class="px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Backups</h1>
            <p class="mt-2 text-gray-600">Manage your backup operations and schedules</p>
          </div>
          <div class="flex space-x-3">
            <button
              @click="refreshBackups"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              {{ loading ? 'Refreshing...' : 'Refresh' }}
            </button>
            <button
              @click="createNewBackup"
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
              </svg>
              New Backup
            </button>
          </div>
        </div>
      </div>
    </div>

    <main class="px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading && backups.length === 0" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
          <p class="mt-6 text-lg text-gray-600 font-medium">Loading backups...</p>
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
            <h3 class="text-lg font-semibold text-red-800">Error loading backups</h3>
            <p class="text-gray-600 mt-2">{{ error }}</p>
            <button
              @click="refreshBackups"
              class="mt-4 inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div v-else>
        <!-- Search and Filter -->
        <div class="mb-8 flex items-center space-x-4">
          <div class="relative flex-1 max-w-md">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search backups..."
              class="w-full px-3 py-2 pl-10 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <svg class="w-4 h-4 absolute left-3 top-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
          </div>
          <select v-model="statusFilter" class="w-48 px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="running">Running</option>
            <option value="completed">Completed</option>
            <option value="failed">Failed</option>
          </select>
        </div>

        <!-- Empty State -->
        <div v-if="filteredBackups.length === 0" class="text-center py-20">
          <div class="w-24 h-24 bg-gray-100 rounded-lg flex items-center justify-center mx-auto mb-6">
            <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-3">No backups found</h3>
          <p class="text-gray-500 mb-8 max-w-md mx-auto">Get started by creating your first backup to protect your data</p>
          <button
            @click="createNewBackup"
            class="inline-flex items-center px-6 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Create Your First Backup
          </button>
        </div>

        <!-- Backup Grid -->
        <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
          <BackupCard
            v-for="backup in filteredBackups"
            :key="backup.id"
            :backup="backup"
            @edit="editBackup"
            @delete="deleteBackupHandler"
            @execute="executeBackupHandler"
          />
        </div>
      </div>
    </main>

    <ConfirmationModal
      v-if="showDeleteConfirm"
      :message="deleteConfirmMessage"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />

    <Toast
      v-if="toastMessage"
      :message="toastMessage"
      :type="toastType"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { isAuthenticated, fetchBackups, deleteBackup, executeBackup } from '~/lib/api'
import BackupCard from '~/components/BackupCard.vue'
import ConfirmationModal from '~/components/ConfirmationModal.vue'
import Toast from '~/components/Toast.vue'

const router = useRouter()

const backups = ref([])
const loading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const statusFilter = ref('')
const showDeleteConfirm = ref(false)
const deleteConfirmMessage = ref('')
const pendingDeleteId = ref(null)
const toastMessage = ref('')
const toastType = ref('success')

const filteredBackups = computed(() => {
  let filtered = backups.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(backup => 
      backup.name.toLowerCase().includes(query) ||
      backup.source.toLowerCase().includes(query) ||
      backup.destination.toLowerCase().includes(query)
    )
  }

  if (statusFilter.value) {
    filtered = filtered.filter(backup => backup.status === statusFilter.value)
  }

  return filtered
})

const loadBackups = async () => {
  try {
    loading.value = true
    error.value = null
    backups.value = await fetchBackups()
  } catch (err) {
    error.value = err.message
    console.error('Failed to load backups:', err)
  } finally {
    loading.value = false
  }
}

const refreshBackups = () => {
  loadBackups()
}

const createNewBackup = () => {
  router.push('/backup/create')
}

const editBackup = (id) => {
  router.push(`/backup/edit/${id}`)
}

const deleteBackupHandler = (id) => {
  const backup = backups.value.find(b => b.id === id)
  deleteConfirmMessage.value = `Are you sure you want to delete "${backup?.name || 'this backup'}"? This action cannot be undone.`
  pendingDeleteId.value = id
  showDeleteConfirm.value = true
}

const confirmDelete = async () => {
  try {
    await deleteBackup(pendingDeleteId.value)
    await loadBackups()
    showToast('Backup deleted successfully', 'success')
  } catch (err) {
    showToast(err.message || 'Failed to delete backup', 'error')
    console.error('Failed to delete backup:', err)
  } finally {
    showDeleteConfirm.value = false
    pendingDeleteId.value = null
  }
}

const cancelDelete = () => {
  showDeleteConfirm.value = false
  pendingDeleteId.value = null
}

const executeBackupHandler = async (id) => {
  try {
    await executeBackup(id)
    await loadBackups()
    showToast('Backup executed successfully', 'success')
  } catch (err) {
    showToast(err.message || 'Failed to execute backup', 'error')
    console.error('Failed to execute backup:', err)
  }
}

const showToast = (message, type = 'success') => {
  toastMessage.value = message
  toastType.value = type
  setTimeout(() => {
    toastMessage.value = ''
  }, 3000)
}

onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  loadBackups()
})
</script>
