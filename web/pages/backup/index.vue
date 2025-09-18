<template>
  <div class="min-h-screen bg-slate-50">
    <div class="bg-white border-b border-slate-200">
      <div class="px-6 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-slate-900">Backups</h1>
            <p class="mt-1 text-slate-600">Manage your backup operations</p>
          </div>
          <div class="flex space-x-3">
            <button 
              @click="refreshBackups"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              {{ loading ? 'Refreshing...' : 'Refresh' }}
            </button>
            <button 
              @click="createNewBackup"
              class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
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

    <main class="px-6 py-8">
      <div v-if="loading && backups.length === 0" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-200 border-t-blue-600"></div>
          <p class="mt-6 text-lg text-slate-600 font-medium">Loading backups...</p>
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
            <h3 class="text-lg font-semibold text-red-800">Error loading backups</h3>
            <p class="text-red-700 mt-2">{{ error }}</p>
            <button 
              @click="refreshBackups"
              class="mt-4 inline-flex items-center px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>

      <div v-else>
        <div class="mb-6 flex items-center space-x-4">
          <div class="relative flex-1 max-w-md">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search backups..."
              class="w-full pl-10 pr-4 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <svg class="w-4 h-4 absolute left-3 top-3 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
          </div>
          <select v-model="statusFilter" class="px-4 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="running">Running</option>
            <option value="completed">Completed</option>
            <option value="failed">Failed</option>
          </select>
        </div>

        <div v-if="filteredBackups.length === 0" class="text-center py-16 text-slate-500">
          <div class="w-20 h-20 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg class="w-10 h-10 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-slate-900 mb-2">No backups found</h3>
          <p class="text-slate-600 mb-6">Get started by creating your first backup</p>
          <button 
            @click="createNewBackup"
            class="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Create Your First Backup
          </button>
        </div>
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { isAuthenticated, fetchBackups, deleteBackup, executeBackup } from '~/lib/api'
import BackupCard from '~/components/BackupCard.vue'

const router = useRouter()

const backups = ref([])
const loading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const statusFilter = ref('')

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

const deleteBackupHandler = async (id) => {
  if (!confirm('Are you sure you want to delete this backup?')) {
    return
  }

  try {
    await deleteBackup(id)
    await loadBackups()
  } catch (err) {
    error.value = err.message
    console.error('Failed to delete backup:', err)
  }
}

const executeBackupHandler = async (id) => {
  try {
    await executeBackup(id)
    await loadBackups()
  } catch (err) {
    error.value = err.message
    console.error('Failed to execute backup:', err)
  }
}

onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  loadBackups()
})
</script>
