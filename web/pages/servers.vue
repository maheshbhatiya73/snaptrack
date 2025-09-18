<template>
  <div class="min-h-screen bg-slate-50">
    <div class="bg-white border-b border-slate-200">
      <div class="px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-slate-900">Servers</h1>
            <p class="mt-1 text-sm text-slate-600">Manage and monitor your backup servers</p>
          </div>
          <div class="flex space-x-3">
            <button 
              @click="refreshServers"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              {{ loading ? 'Refreshing...' : 'Refresh' }}
            </button>
            <button 
              @click="showAddServerForm"
              class="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors duration-200"
            >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Add Server
          </button>
          </div>
        </div>
      </div>
    </div>

    <main class="px-6 py-6">
      <div v-if="loading && servers.length === 0" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-200 border-t-blue-600"></div>
          <p class="mt-6 text-lg text-slate-600 font-medium">Loading servers...</p>
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
            <h3 class="text-lg font-semibold text-red-800">Error loading servers</h3>
            <p class="text-red-700 mt-2">{{ error }}</p>
            <button 
              @click="refreshServers"
              class="mt-4 inline-flex items-center px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>

      <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
          <div class="flex items-center">
            <div class="p-2 bg-green-100 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-slate-600">Active Servers</p>
                <p class="text-2xl font-bold text-green-600">{{ activeServersCount }}</p>
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
                <p class="text-sm font-medium text-slate-600">Inactive Servers</p>
                <p class="text-2xl font-bold text-red-600">{{ inactiveServersCount }}</p>
              </div>
          </div>
        </div>

        <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
          <div class="flex items-center">
            <div class="p-2 bg-blue-100 rounded-lg">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
              </svg>
            </div>
            <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Remote Servers</p>
                <p class="text-2xl font-bold text-blue-600">{{ remoteServersCount }}</p>
              </div>
          </div>
        </div>

        <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
          <div class="flex items-center">
            <div class="p-2 bg-purple-100 rounded-lg">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
              </svg>
            </div>
            <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Total Servers</p>
                <p class="text-2xl font-bold text-purple-600">{{ servers.length }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="mb-6 flex items-center space-x-4">
          <div class="relative flex-1 max-w-md">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search servers..."
              class="w-full pl-10 pr-4 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <svg class="w-4 h-4 absolute left-3 top-3 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
      </div>
          <select v-model="typeFilter" class="px-4 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="">All Types</option>
            <option value="remote">Remote</option>
            <option value="local">Local</option>
          </select>
          <select v-model="statusFilter" class="px-4 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="">All Status</option>
            <option value="enabled">Enabled</option>
            <option value="disabled">Disabled</option>
          </select>
        </div>

        <div v-if="filteredServers.length === 0" class="text-center py-16 text-slate-500">
          <div class="w-20 h-20 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg class="w-10 h-10 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
                    </svg>
                  </div>
          <h3 class="text-xl font-semibold text-slate-900 mb-2">No servers found</h3>
          <p class="text-slate-600 mb-6">Get started by adding your first server</p>
          <button 
            @click="showAddServerForm"
            class="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                    </svg>
            Add Your First Server
                  </button>
                </div>
        <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
          <ServerCard
            v-for="server in filteredServers"
            :key="server.id"
            :server="server"
            @edit="editServer"
            @delete="deleteServerHandler"
          />
        </div>
      </div>
    </main>

    <ServerForm
      v-if="showForm"
      :server="editingServer"
      :is-edit="!!editingServer"
      @close="closeForm"
      @success="handleFormSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { isAuthenticated, fetchServers, deleteServer, fetchServerStatuses } from '~/lib/api'
import ServerCard from '~/components/ServerCard.vue'
import ServerForm from '~/components/ServerForm.vue'

const router = useRouter()

const servers = ref([])
const loading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const showForm = ref(false)
const editingServer = ref(null)
const autoRefreshInterval = ref(null)
const autoRefreshEnabled = ref(true)

const filteredServers = computed(() => {
  let filtered = servers.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(server => 
      server.name.toLowerCase().includes(query) ||
      server.host.toLowerCase().includes(query) ||
      (server.ssh_user && server.ssh_user.toLowerCase().includes(query))
    )
  }

  if (typeFilter.value) {
    filtered = filtered.filter(server => server.type === typeFilter.value)
  }

  if (statusFilter.value) {
    const enabled = statusFilter.value === 'enabled'
    filtered = filtered.filter(server => server.enabled === enabled)
  }

  return filtered
})

const activeServersCount = computed(() => 
  servers.value.filter(server => server.enabled).length
)

const inactiveServersCount = computed(() => 
  servers.value.filter(server => !server.enabled).length
)

const remoteServersCount = computed(() => 
  servers.value.filter(server => server.type === 'remote').length
)

const loadServers = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await fetchServers()
    servers.value = response.servers || response || []
    
    if (response.server_statuses) {
      servers.value = servers.value.map(server => ({
        ...server,
        connection_status: response.server_statuses[server.id] || 'unknown',
        last_checked: response.server_statuses[server.id] ? new Date().toISOString() : server.last_checked
      }))
    }
  } catch (err) {
    error.value = err.message
    console.error('Failed to load servers:', err)
  } finally {
    loading.value = false
  }
}

const refreshServers = () => {
  loadServers()
}

const startAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
  }
  
  if (autoRefreshEnabled.value) {
    autoRefreshInterval.value = setInterval(async () => {
      if (!loading.value && document.visibilityState === 'visible') {
        await loadServerStatuses()
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

const loadServerStatuses = async () => {
  try {
    const response = await fetchServerStatuses()
    
    if (response.server_statuses) {
      servers.value = servers.value.map(server => ({
        ...server,
        connection_status: response.server_statuses[server.id] || 'unknown',
        last_checked: response.server_statuses[server.id] ? new Date().toISOString() : server.last_checked
      }))
    }
  } catch (err) {
    console.error('Failed to refresh server statuses:', err)
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

const showAddServerForm = () => {
  editingServer.value = null
  showForm.value = true
}

const editServer = (server) => {
  editingServer.value = server
  showForm.value = true
}

const closeForm = () => {
  showForm.value = false
  editingServer.value = null
}

const handleFormSuccess = () => {
  loadServers()
}

const deleteServerHandler = async (id) => {
  if (!confirm('Are you sure you want to delete this server? This action cannot be undone.')) {
    return
  }

  try {
    await deleteServer(id)
    await loadServers()
  } catch (err) {
    error.value = err.message
    console.error('Failed to delete server:', err)
  }
}

onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  loadServers()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>
