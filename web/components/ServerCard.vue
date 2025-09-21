<template>
  <div class="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow group">
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center space-x-4">
          <div class="flex-shrink-0">
            <div class="w-12 h-12 bg-gray-900 rounded-lg flex items-center justify-center group-hover:scale-105 transition-transform duration-200">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
              </svg>
            </div>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">{{ server.name }}</h3>
            <div class="flex items-center space-x-3 mt-2">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="server.enabled ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                {{ server.enabled ? 'Enabled' : 'Disabled' }}
              </span>
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="server.type === 'remote' ? 'bg-blue-100 text-blue-800' : 'bg-yellow-100 text-yellow-800'">
                {{ server.type === 'remote' ? 'Remote' : 'Local' }}
              </span>
            </div>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click="testConnection"
            :disabled="testingConnection || server.type !== 'remote'"
            class="p-2 text-gray-400 hover:text-green-600 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg hover:bg-gray-50"
            :title="server.type === 'remote' ? 'Test Connection' : 'Test not available for local servers'"
          >
            <svg v-if="testingConnection" class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </button>
          <button
            @click="$emit('edit', server)"
            class="p-2 text-gray-400 hover:text-gray-600 transition-colors duration-200 rounded-lg hover:bg-gray-50"
            title="Edit Server"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
          </button>
          <button
            @click="confirmDelete"
            class="p-2 text-gray-400 hover:text-red-600 transition-colors duration-200 rounded-lg hover:bg-red-50"
            title="Delete Server"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </button>
        </div>
      </div>

      <div class="border-t border-gray-200 mb-6"></div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div>
          <h4 class="text-sm font-semibold text-gray-900 mb-4 flex items-center">
            <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9"/>
            </svg>
            Host Information
          </h4>
          <div class="space-y-3">
            <div class="flex items-center text-sm">
              <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9"/>
              </svg>
              <span class="font-mono text-gray-900">{{ server.host }}</span>
            </div>
            <div v-if="server.type === 'remote'" class="flex items-center text-sm">
              <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
              </svg>
              <span class="text-gray-900">{{ server.ssh_user }}@{{ server.host }}:{{ server.ssh_port }}</span>
            </div>
          </div>
        </div>

        <div v-if="server.type === 'remote'">
          <h4 class="text-sm font-semibold text-gray-900 mb-4 flex items-center">
            <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
            </svg>
            SSH Configuration
          </h4>
          <div class="space-y-3">
            <div class="flex items-center text-sm">
              <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
              </svg>
              <span class="font-mono text-xs text-gray-900 truncate" :title="server.ssh_key_path">
                {{ server.ssh_key_path }}
              </span>
            </div>
          </div>
        </div>

        <div>
          <h4 class="text-sm font-semibold text-gray-900 mb-4 flex items-center">
            <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            Status
          </h4>
          <div class="space-y-3">
            <div class="flex items-center text-sm">
              <div class="w-2 h-2 rounded-full mr-3" :class="server.enabled ? 'bg-green-500' : 'bg-red-500'"></div>
              <span class="text-gray-900">{{ server.enabled ? 'Active' : 'Inactive' }}</span>
            </div>
            <div class="text-xs text-gray-500">
              Created: {{ formatDate(server.created_at) }}
            </div>
          </div>
        </div>

        <div>
          <h4 class="text-sm font-semibold text-gray-900 mb-4 flex items-center">
            <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            Last Updated
          </h4>
          <div class="text-sm text-gray-900">
            {{ formatDate(server.updated_at) }}
          </div>
        </div>
      </div>

      <div v-if="connectionTestResult" class="mt-6 pt-4">
        <div class="border-t border-gray-200 mb-4"></div>
        <div class="flex items-center space-x-3 p-4 rounded-lg" :class="connectionTestResult.success ? 'bg-green-50' : 'bg-red-50'">
          <svg v-if="connectionTestResult.success" class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <svg v-else class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <div>
            <span class="text-sm font-medium" :class="connectionTestResult.success ? 'text-green-800' : 'text-red-800'">
              {{ connectionTestResult.success ? 'Connection successful' : 'Connection failed' }}
            </span>
            <span v-if="!connectionTestResult.success" class="text-sm text-red-700 block mt-1">
              {{ connectionTestResult.message }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { testServerConnection } from '~/lib/api'

const props = defineProps({
  server: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['edit', 'delete'])

const testingConnection = ref(false)
const connectionTestResult = ref(null)

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const testConnection = async () => {
  if (props.server.type !== 'remote') return

  try {
    testingConnection.value = true
    connectionTestResult.value = null
    
    const result = await testServerConnection(props.server.id)
    connectionTestResult.value = result
    
    setTimeout(() => {
      connectionTestResult.value = null
    }, 5000)
  } catch (error) {
    connectionTestResult.value = {
      success: false,
      message: error.message || 'Connection test failed'
    }
    
    setTimeout(() => {
      connectionTestResult.value = null
    }, 5000)
  } finally {
    testingConnection.value = false
  }
}

const confirmDelete = () => {
  emit('delete', props.server.id)
}
</script>
