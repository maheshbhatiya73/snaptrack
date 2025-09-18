<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
      <div class="px-6 py-4 border-b border-slate-200">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold text-slate-900">
            {{ isEdit ? 'Edit Server' : 'Add New Server' }}
          </h2>
          <div class="flex gap-4">
          <div v-if="formData.type === 'remote' && isEdit" class="">
            <button
              type="button"
              @click="testConnection"
              :disabled="testingConnection"
              class="inline-flex items-center px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 disabled:opacity-50"
            >
              <svg v-if="testingConnection" class="w-4 h-4 mr-2 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              {{ testingConnection ? 'Testing...' : 'Test Connection' }}
            </button>
          </div>
          <button 
            @click="$emit('close')"
            class="p-2 text-slate-400 hover:text-slate-600 transition-colors duration-200"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
          </div>
        </div>
      </div>

      <form @submit.prevent="handleSubmit" class="px-6 py-6">
        <div class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              Server Name *
            </label>
            
            <input
              v-model="formData.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter server name"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              Server Type *
            </label>
            <select
              v-model="formData.type"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="remote">Remote Server</option>
              <option value="local">Local Server</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              Host/IP Address *
            </label>
            <input
              v-model="formData.host"
              type="text"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="192.168.1.100 or hostname.com"
            />
          </div>

          <div v-if="formData.type === 'remote'" class="space-y-4">
            <h3 class="text-lg font-medium text-slate-900">SSH Configuration</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-2">
                  SSH Username *
                </label>
                <input
                  v-model="formData.ssh_user"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="ubuntu, root, etc."
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-slate-700 mb-2">
                  SSH Port *
                </label>
                <input
                  v-model="formData.ssh_port"
                  type="number"
                  required
                  min="1"
                  max="65535"
                  class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="22"
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-700 mb-2">
                SSH Key Path *
              </label>
              <input
                v-model="formData.ssh_key_path"
                type="text"
                required
                class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="/home/user/.ssh/id_rsa"
              />
              <p class="mt-1 text-sm text-slate-500">
                Path to your private SSH key file
              </p>
            </div>
          </div>

          <div>
            <label class="flex items-center">
              <input
                v-model="formData.enabled"
                type="checkbox"
                class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-slate-300 rounded"
              />
              <span class="ml-2 text-sm text-slate-700">Server is enabled</span>
            </label>
          </div>

          
        </div>

        <div class="flex items-center justify-end space-x-3 pt-6 border-t border-slate-200 mt-6">
          <button
            type="button"
            @click="$emit('close')"
            class="px-4 py-2 text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50 transition-colors duration-200"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors duration-200"
          >
            {{ loading ? 'Saving...' : (isEdit ? 'Update Server' : 'Create Server') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { createServer, updateServer, testServerConnection } from '~/lib/api'

const props = defineProps({
  server: {
    type: Object,
    default: null
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'success'])

const loading = ref(false)
const testingConnection = ref(false)

const formData = reactive({
  name: '',
  host: '',
  type: 'remote',
  ssh_user: '',
  ssh_port: 22,
  ssh_key_path: '',
  enabled: true
})

watch(() => props.server, (newServer) => {
  if (newServer) {
    Object.assign(formData, {
      name: newServer.name || '',
      host: newServer.host || '',
      type: newServer.type || 'remote',
      ssh_user: newServer.ssh_user || '',
      ssh_port: newServer.ssh_port || 22,
      ssh_key_path: newServer.ssh_key_path || '',
      enabled: newServer.enabled !== false
    })
  }
}, { immediate: true })

const handleSubmit = async () => {
  try {
    loading.value = true
    
    const serverData = {
      name: formData.name,
      host: formData.host,
      type: formData.type,
      enabled: formData.enabled
    }

    if (formData.type === 'remote') {
      serverData.ssh_user = formData.ssh_user
      serverData.ssh_port = parseInt(formData.ssh_port)
      serverData.ssh_key_path = formData.ssh_key_path
    }

    if (props.isEdit) {
      await updateServer(props.server.id, serverData)
    } else {
      await createServer(serverData)
    }

    emit('success')
    emit('close')
  } catch (error) {
    console.error('Failed to save server:', error)
    alert(error.message || 'Failed to save server')
  } finally {
    loading.value = false
  }
}

const testConnection = async () => {
  if (!props.isEdit || !props.server?.id) return

  try {
    testingConnection.value = true
    const result = await testServerConnection(props.server.id)
    
    if (result.success) {
      alert('Connection test successful!')
    } else {
      alert(`Connection test failed: ${result.message || 'Unknown error'}`)
    }
  } catch (error) {
    console.error('Connection test failed:', error)
    alert(`Connection test failed: ${error.message}`)
  } finally {
    testingConnection.value = false
  }
}
</script>
