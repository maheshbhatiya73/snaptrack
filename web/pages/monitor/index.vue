<template>
  <div class="min-h-screen bg-white">
    <!-- Header Section -->
    <div class="bg-white border-b border-gray-200">
      <div class="px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Monitor</h1>
            <p class="mt-2 text-gray-600">Live metrics for your backup servers</p>
          </div>
          <div class="flex space-x-3">
            <button
              @click="refreshServers"
              :disabled="loading"
              class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              {{ loading ? 'Refreshing...' : 'Refresh' }}
            </button>
            <button
              @click="toggleAutoRefresh"
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white"
              :class="autoRefreshEnabled ? 'bg-gray-900 hover:bg-gray-800 focus:ring-gray-500' : 'bg-gray-500 hover:bg-gray-600 focus:ring-gray-400'"
            >
              <svg v-if="autoRefreshEnabled" class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
              </svg>
              {{ autoRefreshEnabled ? 'Auto-refresh On' : 'Auto-refresh Off' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <main class="px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading && servers.length === 0" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
          <p class="mt-6 text-lg text-gray-600 font-medium">Loading servers...</p>
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
            <h3 class="text-lg font-semibold text-red-800">Error loading servers</h3>
            <p class="text-gray-600 mt-2">{{ error }}</p>
            <button
              @click="refreshServers"
              class="mt-4 inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div v-else>
        <!-- Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="p-3 bg-green-100 rounded-lg">
                <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Active Servers</p>
                <p class="text-2xl font-bold text-green-600">{{ activeServersCount }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="p-3 bg-red-100 rounded-lg">
                <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Inactive Servers</p>
                <p class="text-2xl font-bold text-red-600">{{ inactiveServersCount }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="p-3 bg-blue-100 rounded-lg">
                <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Remote Servers</p>
                <p class="text-2xl font-bold text-blue-600">{{ remoteServersCount }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="p-3 bg-gray-100 rounded-lg">
                <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-600">Total Servers</p>
                <p class="text-2xl font-bold text-gray-900">{{ servers.length }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Search and Filter -->
        <div class="mb-8 flex items-center space-x-4">
          <div class="relative flex-1 max-w-md">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search servers..."
              class="w-full px-3 py-2 pl-10 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <svg class="w-4 h-4 absolute left-3 top-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
          </div>
          <select v-model="typeFilter" class="w-40 px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
            <option value="">All Types</option>
            <option value="remote">Remote</option>
            <option value="local">Local</option>
          </select>
          <select v-model="statusFilter" class="w-40 px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
            <option value="">All Status</option>
            <option value="enabled">Enabled</option>
            <option value="disabled">Disabled</option>
          </select>
        </div>

        <!-- Empty State -->
        <div v-if="filteredServers.length === 0" class="text-center py-20">
          <div class="w-24 h-24 bg-gray-100 rounded-lg flex items-center justify-center mx-auto mb-6">
            <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-3">No servers found</h3>
          <p class="text-gray-500 mb-8 max-w-md mx-auto">Add servers on the Servers page to monitor metrics</p>
        </div>

        <!-- Server Grid -->
        <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
          <div
            v-for="server in filteredServers"
            :key="server.id"
            class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
            @click="goToMonitor(server)"
          >
            <div class="flex items-start justify-between">
              <div>
                <div class="text-lg font-semibold text-gray-900">{{ server.name }}</div>
                <div class="text-sm text-gray-500">{{ server.host }}</div>
              </div>
              <div class="flex items-center space-x-3">
                <span class="px-2 py-1 rounded text-xs" :class="server.type==='remote' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'">{{ server.type }}</span>
                <!-- CPU Circular Meter -->
                <div class="relative w-12 h-12" @click.stop="goToMonitor(server)">
                  <svg viewBox="0 0 48 48" class="w-12 h-12">
                    <circle cx="24" cy="24" r="18" class="text-gray-200" stroke="currentColor" stroke-width="6" fill="transparent" />
                    <circle
                      cx="24"
                      cy="24"
                      r="18"
                      :stroke="cpuColor(cpuPercentFor(server.id))"
                      stroke-width="6"
                      fill="transparent"
                      stroke-linecap="round"
                      :style="{ transform: 'rotate(-90deg)', transformOrigin: '50% 50%' }"
                      :stroke-dasharray="circumference"
                      :stroke-dashoffset="cpuDashOffset(cpuPercentFor(server.id))"
                    />
                  </svg>
                  <div class="absolute inset-0 flex items-center justify-center">
                    <span class="text-xs font-semibold text-gray-900">{{ cpuPercentFor(server.id) != null ? cpuPercentFor(server.id).toFixed(0)+'%' : '—' }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="mt-4 grid grid-cols-3 gap-4 text-sm">
              <div>
                <div class="text-xs text-gray-500">Status</div>
                <div :class="['font-medium', (server.connection_status||'unknown')==='connected' ? 'text-green-600' : 'text-gray-700']">{{ server.connection_status || 'unknown' }}</div>
              </div>
              <div>
                <div class="text-xs text-gray-500">Enabled</div>
                <div class="font-medium">{{ server.enabled ? 'Yes' : 'No' }}</div>
              </div>
              <div>
                <div class="text-xs text-gray-500">Last Checked</div>
                <div class="font-medium">{{ formatDate(server.last_checked) }}</div>
              </div>
            </div>
            <div class="mt-6">
              <button class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500">
                View Monitor
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { isAuthenticated, fetchServers, monitorBatchWsUrl } from '~/lib/api'

const router = useRouter()

const servers = ref([])
const loading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const autoRefreshInterval = ref(null)
const autoRefreshEnabled = ref(false)
const snapshotById = reactive({})
let socket = null

const circumference = 2 * Math.PI * 18

const filteredServers = computed(() => {
  let filtered = servers.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(server => 
      (server.name && server.name.toLowerCase().includes(query)) ||
      (server.host && server.host.toLowerCase().includes(query)) ||
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

const formatDate = (iso) => {
  if (!iso) return '—'
  try { return new Date(iso).toLocaleString() } catch { return '—' }
}

const cpuPercentFor = (id) => {
  const v = snapshotById[id]?.cpu_percent
  if (v == null) return null
  return Math.max(0, Math.min(100, Number(v)))
}

const cpuDashOffset = (p) => {
  if (p == null) return circumference
  return circumference - (circumference * p) / 100
}

const cpuColor = (p) => {
  if (p == null) return '#9CA3AF'
  if (p > 85) return '#EF4444'
  if (p > 65) return '#F59E0B'
  return '#10B981'
}

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

    // websocket will fill snapshots
  } catch (err) {
    error.value = err.message
    console.error('Failed to load servers:', err)
  } finally {
    loading.value = false
  }
}

const connectWs = () => {
  if (socket) { try { socket.close() } catch {} socket = null }
  const url = monitorBatchWsUrl()
  console.debug('[monitor] connecting ws', url)
  socket = new WebSocket(url)
  socket.onmessage = ev => {
    try {
      const arr = JSON.parse(ev.data)
      if (Array.isArray(arr)) {
        for (const m of arr) {
          if (m && m.server_id != null) snapshotById[m.server_id] = m
        }
      }
    } catch {}
  }
  socket.onopen = () => { console.debug('[monitor] ws connected') }
  socket.onerror = (e) => { console.warn('[monitor] ws error', e) }
  socket.onclose = () => {
    console.debug('[monitor] ws closed, retrying in 3s')
    socket = null
    setTimeout(() => { if (!socket) connectWs() }, 3000)
  }
}

const refreshServers = () => {
  loadServers()
}

const startAutoRefresh = () => {
  // no-op; realtime via websocket
}

const stopAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
    autoRefreshInterval.value = null
  }
}

const loadServerStatuses = async () => {}

const toggleAutoRefresh = () => {
  autoRefreshEnabled.value = !autoRefreshEnabled.value
  if (autoRefreshEnabled.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

const goToMonitor = (server) => {
  router.push(`/monitor/${server.id}`)
}

onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }

  loadServers()
  startAutoRefresh()
  connectWs()
})

onUnmounted(() => {
  stopAutoRefresh()
  if (socket) try { socket.close() } catch {}
})
</script>

<style scoped>
</style>
