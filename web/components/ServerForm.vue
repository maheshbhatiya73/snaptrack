<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold text-gray-900">
            {{ isEdit ? 'Edit Server' : 'Add New Server' }}
          </h2>
          <div class="flex gap-4">
            <div v-if="formData.type === 'remote' && isEdit">
              <button type="button" @click="testConnection" :disabled="testingConnection"
                class="inline-flex items-center px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 transition-colors">
                <svg v-if="testingConnection" class="w-4 h-4 mr-2 animate-spin" fill="none" stroke="currentColor"
                  viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ testingConnection ? 'Testing...' : 'Test Connection' }}
              </button>
            </div>
            <button @click="$emit('close')"
              class="p-2 text-gray-400 hover:text-gray-600 transition-colors duration-200">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <form @submit.prevent="handleSubmit" class="px-6 py-6">
        <div class="space-y-6">
          <!-- Server Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Server Name *
            </label>
            <div class="relative">
              <input v-model="formData.name" type="text" required :class="[
                'w-full px-3 py-2 pr-10 border rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors',
                nameValid === null ? 'border-gray-300' :
                  nameValid ? 'border-green-500' : 'border-red-500'
              ]" placeholder="Enter server name" />
              <div class="absolute inset-y-0 right-0 flex items-center pr-3">
                <div v-if="validatingName" class="w-4 h-4">
                  <svg class="animate-spin text-blue-500" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                    </path>
                  </svg>
                </div>
                <div v-else-if="nameValid === true" class="w-4 h-4 text-green-500">
                  <svg fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"></path>
                  </svg>
                </div>
                <div v-else-if="nameValid === false" class="w-4 h-4 text-red-500">
                  <svg fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                      clip-rule="evenodd"></path>
                  </svg>
                </div>
              </div>
            </div>
            <p v-if="nameError" class="mt-1 text-sm text-red-600">{{ nameError }}</p>
          </div>

          <!-- Server Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Server Type *
            </label>
            <select v-model="formData.type" required
              class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors">
              <option value="remote">Remote Server</option>
              <option value="local">Local Server</option>
            </select>
          </div>

          <!-- Transfer Type (Only for Remote) -->
          <div v-if="formData.type === 'remote'">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Transfer Type *
            </label>
            <select v-model="formData.TransferType" required
              class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors">
              <option value="rsync">Rsync</option>
              <option value="scp">SCP</option>
            </select>
          </div>
          <div v-else class="text-sm text-gray-600">
            <p><strong>Transfer Type:</strong> Local (automatic)</p>
          </div>

          <!-- Host (Only for Remote) -->
          <div v-if="formData.type === 'remote'">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Host/IP Address *
            </label>
            <input v-model="formData.host" type="text" required
              class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
              placeholder="192.168.1.100 or hostname.com" />
          </div>
          <div v-else class="text-sm text-gray-600">
            <p><strong>Local Server:</strong> This will use the local machine (localhost) for operations.</p>
          </div>

          <!-- SSH Configuration (Only for Remote) -->
          <div v-if="formData.type === 'remote'" class="space-y-4">
            <h3 class="text-lg font-medium text-gray-900">SSH Configuration</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">SSH Username *</label>
                <input v-model="formData.ssh_user" type="text" required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
                  placeholder="ubuntu, root, etc." />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">SSH Port *</label>
                <input v-model="formData.ssh_port" type="number" required min="1" max="65535"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
                  placeholder="22" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">SSH Key Path *</label>
              <input v-model="formData.ssh_key_path" type="text" required
                class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
                placeholder="/home/user/.ssh/id_rsa" />
              <p class="mt-1 text-sm text-gray-500">Path to your private SSH key file</p>
            </div>
          </div>

          <!-- Enabled Checkbox -->
          <div>
            <label class="flex items-center">
              <input v-model="formData.enabled" type="checkbox"
                class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded" />
              <span class="ml-2 text-sm text-gray-700">Server is enabled</span>
            </label>
          </div>
        </div>

        <!-- Buttons -->
        <div class="flex items-center justify-end space-x-3 pt-6 border-t border-gray-200 mt-6">
          <button type="button" @click="$emit('close')"
            class="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-200">
            Cancel
          </button>
          <button type="submit" :disabled="loading"
            class="px-4 py-2 bg-gray-900 text-white rounded-md hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 disabled:opacity-50 transition-colors duration-200">
            {{ loading ? 'Saving...' : (isEdit ? 'Update Server' : 'Create Server') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { createServer, updateServer, testServerConnection, checkServerNameExists } from '~/lib/api'

const props = defineProps({
  server: { type: Object, default: null },
  isEdit: { type: Boolean, default: false }
})
const emit = defineEmits(['close', 'success', 'error'])

const loading = ref(false)
const testingConnection = ref(false)

// Validation states
const nameValid = ref(null)
const nameError = ref('')
const validatingName = ref(false)

const formData = reactive({
  name: '',
  host: '',
  type: 'remote',
  ssh_user: '',
  ssh_port: 22,
  ssh_key_path: '',
  TransferType: 'rsync',
  enabled: true
})

// Prefill when editing
watch(() => props.server, (newServer) => {
  if (newServer) {
    Object.assign(formData, {
      name: newServer.name || '',
      host: newServer.host || '',
      type: newServer.type || 'remote',
      ssh_user: newServer.ssh_user || '',
      ssh_port: newServer.ssh_port || 22,
      TransferType: newServer.TransferType || (newServer.type === 'local' ? 'local' : 'rsync'),
      ssh_key_path: newServer.ssh_key_path || '',
      enabled: newServer.enabled !== false
    })
  }
}, { immediate: true })

// Validate server name
watch(() => formData.name, (newName) => {
  if (newName && newName.trim() !== '') {
    setTimeout(() => validateServerName(newName), 500)
  } else {
    nameValid.value = null
    nameError.value = ''
  }
})

const validateServerName = async (name) => {
  if (!name || name.trim() === '') return
  validatingName.value = true
  nameValid.value = null
  nameError.value = ''
  try {
    const exists = await checkServerNameExists(name, props.isEdit ? props.server?.id : null)
    nameValid.value = !exists
    if (exists) nameError.value = 'Server name already exists'
  } catch (err) {
    nameValid.value = false
    nameError.value = 'Failed to validate server name'
  } finally {
    validatingName.value = false
  }
}

const handleSubmit = async () => {
  if (!formData.name.trim()) return emit('error', 'Server name is required')
  if (formData.type === 'remote' && (!formData.host || !formData.ssh_user || !formData.ssh_key_path)) {
    return emit('error', 'Please fill in all required fields for remote server')
  }

  if (nameValid.value === null) await validateServerName(formData.name)
  if (validatingName.value) {
    await new Promise(resolve => {
      const unwatch = watch(validatingName, (isValidating) => {
        if (!isValidating) { unwatch(); resolve() }
      })
    })
  }
  if (nameValid.value === false) return emit('error', nameError.value || 'Server name validation failed')

  try {
    loading.value = true
    const serverData = {
      name: formData.name,
      type: formData.type,
      enabled: formData.enabled,
      host: formData.type === 'remote' ? formData.host : 'localhost',
      ssh_user: formData.type === 'remote' ? formData.ssh_user : undefined,
      ssh_port: formData.type === 'remote' ? parseInt(formData.ssh_port) : undefined,
      ssh_key_path: formData.type === 'remote' ? formData.ssh_key_path : undefined,
      TransferType: formData.type === 'remote' ? formData.TransferType : 'local'
    }
    if (props.isEdit) {
      await updateServer(props.server.id, serverData)
      emit('success')
      emit('close')
    } else {
      const result = await createServer(serverData)
      result.success ? (emit('success'), emit('close')) : emit('error', result.message)
    }
  } catch (error) {
    console.error('Failed to save server:', error)
    emit('error', error.message || 'Failed to save server')
  } finally {
    loading.value = false
  }
}

const testConnection = async () => {
  if (!props.isEdit || !props.server?.id) return
  try {
    testingConnection.value = true
    const result = await testServerConnection(props.server.id)
    result.success ? emit('success', 'Connection test successful!')
                   : emit('error', `Connection test failed: ${result.message || 'Unknown error'}`)
  } catch (error) {
    console.error('Connection test failed:', error)
    emit('error', `Connection test failed: ${error.message}`)
  } finally {
    testingConnection.value = false
  }
}

onMounted(() => {
  if (props.isEdit && formData.name) validateServerName(formData.name)
})
</script>
