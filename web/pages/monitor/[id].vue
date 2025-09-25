<!-- /mahesh/maheshstore/cloud/snaptrack/web/pages/monitor/[id].vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
    <!-- Loading State -->
    <div v-if="loading" class="min-h-screen flex items-center justify-center">
      <div class="text-center">
        <div class="relative">
          <div class="w-20 h-20 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mx-auto"></div>
          <i
            class="i-lucide-server w-8 h-8 text-blue-600 absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2"></i>
        </div>
        <p class="mt-6 text-lg text-gray-600 font-medium">Connecting to server...</p>
        <div class="mt-2 flex justify-center space-x-1">
          <div v-for="i in [0, 1, 2]" :key="i" class="w-2 h-2 bg-blue-600 rounded-full animate-bounce"
            :style="{ animationDelay: `${i * 0.1}s` }"></div>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error"
      class="min-h-screen bg-gradient-to-br from-red-50 to-white flex items-center justify-center p-4">
      <div class="bg-white border border-red-200 rounded-2xl p-8 max-w-md w-full shadow-lg">
        <div class="text-center">
          <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <i class="i-lucide-activity w-8 h-8 text-red-600"></i>
          </div>
          <h3 class="text-xl font-semibold text-red-800 mb-2">Connection Error</h3>
          <p class="text-gray-600 mb-6">{{ error }}</p>
          <button @click="reconnect"
            class="px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors">
            Try Again
          </button>
        </div>
      </div>
    </div>

    <!-- Main Dashboard -->
    <div v-else>
      <!-- Header -->
      <div class="bg-white/80 backdrop-blur-sm border-b border-gray-200/50 sticky top-0 z-10">
        <div class="px-8 py-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">

              <div>
                <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
                  Server Monitor
                </h1>
              </div>
            </div>
            <div class="flex items-center space-x-4">
              <!-- Connection Status -->
              <div class="flex items-center space-x-3">
                <div class="flex items-center space-x-2">
                  <div :class="`w-3 h-3 rounded-full ${connectionStatusConfig.color} ${connectionStatusConfig.pulse}`">
                  </div>
                  <span class="text-sm font-medium text-gray-700">{{ connectionStatusConfig.text }}</span>
                </div>
                <button v-if="connectionStatus !== 'connected'" @click="reconnect"
                  class="flex items-center space-x-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium rounded-lg transition-colors">
                  <i class="i-lucide-refresh-cw w-3 h-3"></i>
                  <span>Reconnect</span>
                </button>
              </div>
              <NuxtLink to="/dashboard"
                class="flex items-center space-x-2 px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white font-medium rounded-xl transition-colors">
                <span>Back to Dashboard</span>
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <main class="px-8 py-8 space-y-8">
        <!-- Server Info & CPU Circle -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Server Info -->
          <div class="lg:col-span-2 bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-2xl p-6 shadow-lg">
            <!-- Server Header -->
            <div class="flex items-center justify-between mb-6">
              <div>
                <h2 class="text-xl font-bold text-gray-800">{{ server?.name || 'Unnamed Server' }}</h2>
                <p class="text-gray-500">{{ server?.host }}</p>
              </div>
              <span
                :class="`px-3 py-1 rounded-full text-sm font-medium ${server?.type === 'remote' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'}`">
                {{ server?.type === 'remote' ? 'Remote Server' : 'Local Server' }}
              </span>
            </div>

            <!-- Server Metrics -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div class="text-center p-4 bg-gradient-to-r from-green-50 to-green-100 rounded-xl">
                <i class="i-lucide-clock w-6 h-6 text-green-600 mx-auto mb-2"></i>
                <div class="text-sm text-gray-600">Uptime</div>
                <div class="text-lg font-bold text-gray-800">{{ formatUptime(metrics?.uptime_seconds) }}</div>
              </div>
              <div class="text-center p-4 bg-gradient-to-r from-blue-50 to-blue-100 rounded-xl">
                <i class="i-lucide-zap w-6 h-6 text-blue-600 mx-auto mb-2"></i>
                <div class="text-sm text-gray-600">Load Average</div>
                <div class="text-lg font-bold text-gray-800">{{ formatLoad(metrics) }}</div>
              </div>
              <div class="text-center p-4 bg-gradient-to-r from-purple-50 to-purple-100 rounded-xl">
                <i class="i-lucide-activity w-6 h-6 text-purple-600 mx-auto mb-2"></i>
                <div class="text-sm text-gray-600">Last Update</div>
                <div class="text-lg font-bold text-gray-800">
                  {{ lastUpdate ? new Date(lastUpdate).toLocaleTimeString() : '—' }}
                </div>
              </div>
            </div>

            <!-- Additional Server Details -->
            <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600">
              <div>
                <strong>SSH User:</strong> {{ server?.ssh_user || '—' }}
              </div>
              <div>
                <strong>SSH Port:</strong> {{ server?.ssh_port || '—' }}
              </div>
              <div>
                <strong>Transfer Type:</strong> {{ server?.transferType || '—' }}
              </div>
              <div>
                <strong>Status:</strong> {{ server?.enabled ? 'Enabled' : 'Disabled' }}
              </div>
            </div>
          </div>


          <!-- CPU Circle -->
          <div class="bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-2xl p-6 shadow-lg">
            <div class="text-center">
              <div class="flex items-center justify-center mb-4">
                <i class="i-lucide-cpu w-6 h-6 text-blue-600 mr-2"></i>
                <h3 class="text-lg font-semibold text-gray-800">CPU Usage</h3>
              </div>
              <div class="flex justify-center mb-4">
                <!-- Circular Progress for CPU -->
                <div class="relative inline-flex items-center justify-center">
                  <svg :width="cpuCircle.size" :height="cpuCircle.size" class="transform -rotate-90">
                    <circle :cx="cpuCircle.size / 2" :cy="cpuCircle.size / 2" :r="cpuCircle.radius"
                      stroke="rgba(156, 163, 175, 0.2)" :stroke-width="cpuCircle.strokeWidth" fill="none" />
                    <circle :cx="cpuCircle.size / 2" :cy="cpuCircle.size / 2" :r="cpuCircle.radius"
                      :stroke="cpuCircle.color" :stroke-width="cpuCircle.strokeWidth" fill="none"
                      :stroke-dasharray="cpuCircle.circumference" :stroke-dashoffset="cpuCircle.offset"
                      stroke-linecap="round" class="transition-all duration-300 ease-in-out"
                      :style="{ filter: 'drop-shadow(0 0 6px ' + cpuCircle.shadowColor + ')' }" />
                  </svg>
                  <div class="absolute inset-0 flex flex-col items-center justify-center">
                    <span class="text-2xl font-bold text-gray-800">{{ (metrics?.cpu_percent || 0).toFixed(1) }}%</span>
                  </div>
                </div>
              </div>
              <div class="text-sm text-gray-600">
                Real-time CPU utilization
              </div>
            </div>
          </div>
        </div>

        <!-- Metric Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Memory Usage Card -->
          <div
            :class="`bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all duration-300 transform ${memoryCard.isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center space-x-3">
                <h3 class="text-sm font-semibold text-gray-700">Memory Usage</h3>
              </div>
            </div>
            <div class="flex items-center justify-between">
              <div>
                <div class="text-2xl font-bold text-gray-800 mb-1">{{ human(metrics?.mem_used_bytes) }}</div>
                <div class="text-sm text-gray-500">{{ `${human(metrics?.mem_total_bytes)} total` }}</div>
              </div>
              <!-- Circular Progress for Memory -->
              <div class="relative inline-flex items-center justify-center">
                <svg :width="metricCircle.size" :height="metricCircle.size" class="transform -rotate-90">
                  <circle :cx="metricCircle.size / 2" :cy="metricCircle.size / 2" :r="metricCircle.radius"
                    stroke="rgba(156, 163, 175, 0.2)" :stroke-width="metricCircle.strokeWidth" fill="none" />
                  <circle :cx="metricCircle.size / 2" :cy="metricCircle.size / 2" :r="metricCircle.radius"
                    :stroke="memoryCircle.color" :stroke-width="metricCircle.strokeWidth" fill="none"
                    :stroke-dasharray="metricCircle.circumference" :stroke-dashoffset="memoryCircle.offset"
                    stroke-linecap="round" class="transition-all duration-300 ease-in-out"
                    :style="{ filter: 'drop-shadow(0 0 6px ' + memoryCircle.shadowColor + ')' }" />
                </svg>
                <div class="absolute inset-0 flex flex-col items-center justify-center">
                  <span class="text-2xl font-bold text-gray-800">{{ memoryCard.percentage.toFixed(1) }}%</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Disk Usage Card -->
          <div
            :class="`bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all duration-300 transform ${diskCard.isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center space-x-3">
                <h3 class="text-sm font-semibold text-gray-700">Disk Usage</h3>
              </div>
            </div>
            <div class="flex items-center justify-between">
              <div>
                <div class="text-2xl font-bold text-gray-800 mb-1">{{ human(metrics?.disk_used_bytes) }}</div>
                <div class="text-sm text-gray-500">{{ `${human(metrics?.disk_total_bytes)} total` }}</div>
              </div>
              <!-- Circular Progress for Disk -->
              <div class="relative inline-flex items-center justify-center">
                <svg :width="metricCircle.size" :height="metricCircle.size" class="transform -rotate-90">
                  <circle :cx="metricCircle.size / 2" :cy="metricCircle.size / 2" :r="metricCircle.radius"
                    stroke="rgba(156, 163, 175, 0.2)" :stroke-width="metricCircle.strokeWidth" fill="none" />
                  <circle :cx="metricCircle.size / 2" :cy="metricCircle.size / 2" :r="metricCircle.radius"
                    :stroke="diskCircle.color" :stroke-width="metricCircle.strokeWidth" fill="none"
                    :stroke-dasharray="metricCircle.circumference" :stroke-dashoffset="diskCircle.offset"
                    stroke-linecap="round" class="transition-all duration-300 ease-in-out"
                    :style="{ filter: 'drop-shadow(0 0 6px ' + diskCircle.shadowColor + ')' }" />
                </svg>
                <div class="absolute inset-0 flex flex-col items-center justify-center">
                  <span class="text-2xl font-bold text-gray-800">{{ diskCard.percentage.toFixed(1) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Process Table -->
        <div class="bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-2xl shadow-lg overflow-hidden">
          <div class="bg-gradient-to-r from-gray-50 to-gray-100 px-6 py-4 border-b border-gray-200">
            <h3 class="text-lg font-semibold text-gray-800 flex items-center">
              <i class="i-lucide-activity w-5 h-5 mr-2 text-blue-600"></i>
              Top Processes
            </h3>
          </div>

          <div v-if="!sortedProcesses.length" class="p-6 text-center text-gray-500">
            No processes running
          </div>

          <div v-else class="overflow-x-auto max-h-96">
            <table class="min-w-full">
              <thead class="bg-gray-50 sticky top-0">
                <tr>
                  <th
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                    @click="handleSort('pid')">
                    PID
                    <i :class="sortBy === 'pid' ? (sortDesc ? 'i-lucide-arrow-down' : 'i-lucide-arrow-up') : ''"
                      class="w-4 h-4 inline-block ml-1"></i>
                  </th>
                  <th
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                    @click="handleSort('name')">
                    Name
                    <i :class="sortBy === 'name' ? (sortDesc ? 'i-lucide-arrow-down' : 'i-lucide-arrow-up') : ''"
                      class="w-4 h-4 inline-block ml-1"></i>
                  </th>
                  <th
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                    @click="handleSort('cpu_percent')">
                    CPU %
                    <i :class="sortBy === 'cpu_percent' ? (sortDesc ? 'i-lucide-arrow-down' : 'i-lucide-arrow-up') : ''"
                      class="w-4 h-4 inline-block ml-1"></i>
                  </th>
                  <th
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
                    @click="handleSort('memory_percent')">
                    Memory %
                    <i :class="sortBy === 'memory_percent' ? (sortDesc ? 'i-lucide-arrow-down' : 'i-lucide-arrow-up') : ''"
                      class="w-4 h-4 inline-block ml-1"></i>
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="process in sortedProcesses" :key="process.pid">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ process.pid }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ process.name }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ process.cpu_percent?.toFixed(1) }}%
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ process.memory_percent?.toFixed(1)
                    }}%</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from '#app';
import { fetchServers, monitorWsUrl } from '~/lib/api';

// Utility function to convert hex to RGBA
const hexToRgba = (hex, opacity) => {
  hex = hex.replace('#', '');
  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);
  return `rgba(${r}, ${g}, ${b}, ${opacity})`;
};

// Main Component Logic
const loading = ref(true);
const route = useRoute();
const error = ref(null);
const server = ref({});
const metrics = ref(null);
const lastUpdate = ref(null);
const connectionStatus = ref('disconnected');
const socketRef = ref(null);
const serverId = computed(() => route.params.id);

// Utility Functions
const human = (bytes) => {
  if (bytes == null) return '—';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let i = 0, v = Number(bytes);
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++; }
  return `${v.toFixed(1)} ${units[i]}`;
};

const formatUptime = (seconds) => {
  if (seconds == null) return '—';
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  return `${d}d ${h}h ${m}m`;
};

const formatLoad = (m) => {
  if (!m) return '—';
  const { load1, load5, load15 } = m;
  if (load1 == null) return '—';
  return `${load1.toFixed(2)} / ${load5.toFixed(2)} / ${load15.toFixed(2)}`;
};

// Connection Status Config
const statusConfigs = {
  connected: { color: 'bg-green-500', text: 'Connected', pulse: 'animate-pulse' },
  connecting: { color: 'bg-yellow-500', text: 'Connecting...', pulse: 'animate-pulse' },
  disconnected: { color: 'bg-gray-500', text: 'Disconnected', pulse: '' },
  error: { color: 'bg-red-500', text: 'Connection Error', pulse: 'animate-pulse' }
};
const connectionStatusConfig = computed(() => statusConfigs[connectionStatus.value] || statusConfigs.disconnected);

// CPU Circle Logic
const cpuPercentage = computed(() => metrics.value?.cpu_percent || 0);
const cpuCircle = computed(() => {
  const size = 140;
  const strokeWidth = 10;
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (cpuPercentage.value / 100) * circumference;
  const color = cpuPercentage.value > 85 ? '#ef4444' : cpuPercentage.value > 65 ? '#f59e0b' : '#10b981';
  const shadowColor = hexToRgba(color, 0.25);
  return { size, strokeWidth, radius, circumference, offset, color, shadowColor };
});

async function load() {
  loading.value = true;
  try {
    const list = await fetchServers();
    const arr = Array.isArray(list) ? list : list?.servers || [];

    // Pick the server by route param first
    server.value = arr.find(s => s.id == serverId.value) || arr.find(s => s.enabled) || arr[0] || null;
  } catch (e) {
    error.value = e?.message || 'Failed to load server';
  } finally {
    loading.value = false;
  }
}


// Metric Circle Defaults
const metricCircle = {
  size: 80,
  strokeWidth: 6,
  radius: (80 - 6) / 2,
  circumference: ((80 - 6) / 2) * 2 * Math.PI
};

// Memory Card Logic
const memoryCard = reactive({
  percentage: computed(() => metrics.value?.mem_used_bytes && metrics.value?.mem_total_bytes ? (metrics.value.mem_used_bytes / metrics.value.mem_total_bytes) * 100 : 0),
  isVisible: false
});
const memoryCircle = computed(() => {
  const offset = metricCircle.circumference - (memoryCard.percentage / 100) * metricCircle.circumference;
  const color = memoryCard.percentage > 85 ? '#ef4444' : memoryCard.percentage > 65 ? '#f59e0b' : '#10b981';
  const shadowColor = hexToRgba(color, 0.25);
  return { offset, color, shadowColor };
});

// Disk Card Logic
const diskCard = reactive({
  percentage: computed(() => metrics.value?.disk_used_bytes && metrics.value?.disk_total_bytes ? (metrics.value.disk_used_bytes / metrics.value.disk_total_bytes) * 100 : 0),
  isVisible: false
});
const diskCircle = computed(() => {
  const offset = metricCircle.circumference - (diskCard.percentage / 100) * metricCircle.circumference;
  const color = diskCard.percentage > 85 ? '#ef4444' : diskCard.percentage > 65 ? '#f59e0b' : '#10b981';
  const shadowColor = hexToRgba(color, 0.25);
  return { offset, color, shadowColor };
});

// Process Table Logic
const sortBy = ref('cpu_percent');
const sortDesc = ref(true);
const sortedProcesses = computed(() => {
  return [...(metrics.value?.processes || [])].sort((a, b) => {
    const aVal = Number(a[sortBy.value] ?? 0);
    const bVal = Number(b[sortBy.value] ?? 0);
    return sortDesc.value ? bVal - aVal : aVal - bVal;
  });
});
const handleSort = (field) => {
  if (sortBy.value === field) {
    sortDesc.value = !sortDesc.value;
  } else {
    sortBy.value = field;
    sortDesc.value = true;
  }
};

// WebSocket Connection
function connect() {
  if (!serverId.value) return;

  if (socketRef.value) {
    socketRef.value.close();
    socketRef.value = null;
  }

  const url = monitorWsUrl(serverId.value);
  connectionStatus.value = 'connecting';

  const ws = new WebSocket(url);
  socketRef.value = ws;

  ws.onopen = () => {
    connectionStatus.value = 'connected';
  };

  ws.onmessage = (ev) => {
    try {
      const data = JSON.parse(ev.data);
      metrics.value = { ...metrics.value, ...data };
      lastUpdate.value = Date.now();
    } catch (e) {
      console.error("Failed to parse WebSocket message", e);
    }
  };

  ws.onclose = () => {
    connectionStatus.value = 'disconnected';
    socketRef.value = null;
  };

  ws.onerror = () => {
    connectionStatus.value = 'error';
    error.value = 'Failed to connect to the server';
  };
}

const reconnect = () => {
  if (socketRef.value) {
    socketRef.value.close();
  }
  connect();
};

// Lifecycle Hooks
onMounted(async () => {
  await load()

  loading.value = true;
  setTimeout(() => {
    loading.value = false;
    connect();
    memoryCard.isVisible = true;
    diskCard.isVisible = true;
  }, 1500);
});

onUnmounted(() => {
  if (socketRef.value) {
    socketRef.value.close();
  }
});
</script>

<style scoped>
.animate-pulse {
  animation: pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.5;
  }
}

.animate-bounce {
  animation: bounce 1s infinite;
}

@keyframes bounce {

  0%,
  100% {
    transform: translateY(-25%);
  }

  50% {
    transform: translateY(0);
  }
}
</style>