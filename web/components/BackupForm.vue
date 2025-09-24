<template>
  <div class="max-w-4xl mx-auto">
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
        <h3 class="text-lg font-semibold text-slate-900 mb-4">Backup Information</h3>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label for="name" class="block text-sm font-medium text-slate-700 mb-2">
              Name *
            </label>
            <input
              id="name"
              v-model="formData.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter backup name"
            />
          </div>

          <div>
            <label for="type" class="block text-sm font-medium text-slate-700 mb-2">
              Type *
            </label>
            <select
              id="type"
              v-model="formData.type"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="full">Full</option>
              <option value="incremental">Incremental</option>
            </select>
          </div>

          <div>
            <label for="file_type" class="block text-sm font-medium text-slate-700 mb-2">
              File Type *
            </label>
            <select
              id="file_type"
              v-model="formData.file_type"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="tar">TAR</option>
              <option value="zip">ZIP</option>
              <option value="raw">RAW</option>
            </select>
          </div>

          <div>
            <label for="schedule_type" class="block text-sm font-medium text-slate-700 mb-2">
              Schedule Type *
            </label>
            <select
              id="schedule_type"
              v-model="formData.schedule_type"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="one_time">One Time</option>
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
            </select>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
        <h3 class="text-lg font-semibold text-slate-900 mb-4">Paths & Servers</h3>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-slate-700 mb-2">
              Target Servers *
            </label>
            <div v-if="loadingServers" class="text-center py-4">
              <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
              <p class="mt-2 text-sm text-slate-600">Loading servers...</p>
            </div>
            <div v-else-if="servers.length === 0" class="text-center py-4 text-slate-500">
              <p class="text-sm">No servers available</p>
            </div>
            <CustomDropdown
              v-if="singleServer"
              v-model="formData.server_id"
              :options="serverOptions"
              placeholder="Select target server"
              searchable
            />
            <MultiSelectDropdown
              v-else
              v-model="formData.server_ids"
              :options="serverOptions"
              placeholder="Select target servers"
              searchable
            />
            <div v-if="serverConnectionErrors.length > 0" class="mt-2">
              <div class="text-sm text-red-600">
                <p class="font-medium">Server connection issues:</p>
                <ul class="list-disc list-inside mt-1">
                  <li v-for="error in serverConnectionErrors" :key="error.serverId">
                    {{ getServerName(error.serverId) }}: {{ error.error }}
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <div v-if="hasSelectedServer">
            <label for="source" class="block text-sm font-medium text-slate-700 mb-2">
              Source Path * <span class="text-xs text-slate-500" v-if="selectedServerType === 'remote'">(validated on local)</span>
            </label>
            <div class="relative">
              <input
                id="source"
                v-model="formData.source"
                type="text"
                required
                :class="[
                  'w-full px-3 py-2 pr-10 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
                  sourceValid === null ? 'border-slate-300' :
                  sourceValid ? 'border-green-500' : 'border-red-500'
                ]"
                :placeholder="selectedServerType === 'remote' ? '/var/www (local source)' : '/var/www'"
              />
              <div class="absolute inset-y-0 right-0 flex items-center pr-3">
                <div v-if="validatingSource" class="w-4 h-4">
                  <svg class="animate-spin text-blue-500" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </div>
                <div v-else-if="sourceValid === true" class="w-4 h-4 text-green-500">
                  <svg fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                  </svg>
                </div>
                <div v-else-if="sourceValid === false" class="w-4 h-4 text-red-500">
                  <svg fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                  </svg>
                </div>
              </div>
            </div>
            <p v-if="sourceError" class="mt-1 text-sm text-red-600">{{ sourceError }}</p>
          </div>

          <div v-if="hasSelectedServer">
            <label for="destination" class="block text-sm font-medium text-slate-700 mb-2">
              Destination Path * <span class="text-xs text-slate-500" v-if="selectedServerType === 'remote'">(validated on remote)</span>
            </label>
            <div class="relative">
              <input
                id="destination"
                v-model="formData.destination"
                type="text"
                required
                :class="[
                  'w-full px-3 py-2 pr-10 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
                  destinationValid === null ? 'border-slate-300' :
                  destinationValid ? 'border-green-500' : 'border-red-500'
                ]"
                :placeholder="selectedServerType === 'remote' ? '/path/on/remote/server' : '/mnt/backups/backup.tar.gz'"
              />
              <div class="absolute inset-y-0 right-0 flex items-center pr-3">
                <div v-if="validatingDestination" class="w-4 h-4">
                  <svg class="animate-spin text-blue-500" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </div>
                <div v-else-if="destinationValid === true" class="w-4 h-4 text-green-500">
                  <svg fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                  </svg>
                </div>
                <div v-else-if="destinationValid === false" class="w-4 h-4 text-red-500">
                  <svg fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                  </svg>
                </div>
              </div>
            </div>
            <p v-if="destinationError" class="mt-1 text-sm text-red-600">{{ destinationError }}</p>
          </div>
          
        </div>
      </div>

      <div class="flex items-center justify-end space-x-3 pt-4">
        <button
          type="button"
          @click="$emit('cancel')"
          class="px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="loading || !isFormValid"
          class="px-6 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
        >
          <span v-if="loading" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
          {{ isEdit ? 'Update Backup' : 'Create Backup' }}
        </button>
      </div>
    </form>
  </div>
</template>


<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { fetchServers, validateServerPath, testServerConnection } from '~/lib/api'
import MultiSelectDropdown from '~/components/MultiSelectDropdown.vue'
import CustomDropdown from '~/components/CustomDropdown.vue'

const props = defineProps({
  backup: {
    type: Object,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  },
  singleServer: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])

const isEdit = !!props.backup
const loadingServers = ref(false)
const servers = ref([])

// Validation states
const validatingSource = ref(false)
const validatingDestination = ref(false)
const sourceValid = ref(null) // null = not validated, true = valid, false = invalid
const destinationValid = ref(null)
const sourceError = ref('')
const destinationError = ref('')
const serverConnectionErrors = ref([]) // Array of {serverId, error}

const formData = reactive({
  name: '',
  type: 'full',
  source: '',
  destination: '',
  file_type: 'tar',
  server_id: null,
  server_ids: [],
  schedule_type: 'one_time'
})

const serverOptions = computed(() => {
  return servers.value.map(server => ({
    value: server.id,
    label: server.name,
    description: `${server.host} - ${server.type || 'Unknown'}`
  }))
})

const selectedServers = computed(() => {
  return props.singleServer
    ? (formData.server_id !== null ? [formData.server_id] : [])
    : formData.server_ids
})

const selectedServerType = computed(() => {
  if (selectedServers.value.length === 0) return null
  const first = servers.value.find(s => s.id === selectedServers.value[0])
  return first ? (first.type || null) : null
})

const hasSelectedServer = computed(() => {
  return props.singleServer
    ? formData.server_id !== null
    : formData.server_ids.length > 0
})

const isFormValid = computed(() => {
  const hasServer = props.singleServer ? formData.server_id !== null : formData.server_ids.length > 0
  return formData.name.trim() !== '' &&
          formData.source.trim() !== '' &&
          formData.destination.trim() !== '' &&
          hasServer &&
          sourceValid.value !== false &&
          destinationValid.value !== false &&
          serverConnectionErrors.value.length === 0
})

const getServerName = (serverId) => {
  const server = servers.value.find(s => s.id === serverId)
  return server ? server.name : `Server ${serverId}`
}

watch(() => props.backup, (newBackup) => {
  if (newBackup) {
    Object.assign(formData, {
      name: newBackup.name || '',
      type: newBackup.type || 'full',
      source: newBackup.source || '',
      destination: newBackup.destination || '',
      file_type: newBackup.file_type || 'tar',
      server_ids: newBackup.server_ids || [],
      server_id: Array.isArray(newBackup.server_ids) && newBackup.server_ids.length > 0 ? newBackup.server_ids[0] : (newBackup.server_id || null),
      schedule_type: newBackup.schedule_type || 'one_time'
    })
  }
}, { immediate: true })

// Watch for path changes and validate
watch(() => formData.source, (newSource) => {
  if (newSource && newSource.trim() !== '') {
    setTimeout(() => validatePath(newSource, true), 500)
  } else {
    sourceValid.value = null
    sourceError.value = ''
  }
})

watch(() => formData.destination, (newDestination) => {
  if (newDestination && newDestination.trim() !== '') {
    setTimeout(() => validatePath(newDestination, false), 500)
  } else {
    destinationValid.value = null
    destinationError.value = ''
  }
})

// Watch for server selection changes and validate connections
watch(() => formData.server_ids, (newServerIds) => {
  if (!props.singleServer) {
    if (newServerIds.length > 0) {
      setTimeout(() => validateServerConnections(), 300)
    } else {
      serverConnectionErrors.value = []
    }
  }
}, { deep: true })

watch(() => formData.server_id, (newServerId) => {
  if (props.singleServer) {
    if (newServerId !== null) {
      setTimeout(() => validateServerConnections(), 300)
    } else {
      serverConnectionErrors.value = []
    }
  }
})

const loadServers = async () => {
  try {
    loadingServers.value = true
    servers.value = await fetchServers()
  } catch (error) {
    console.error('Failed to load servers:', error)
  } finally {
    loadingServers.value = false
  }
}

const validatePath = async (path, isSource = true) => {
  if (!path || path.trim() === '') return;

  const validating = isSource ? validatingSource : validatingDestination;
  const valid = isSource ? sourceValid : destinationValid;
  const error = isSource ? sourceError : destinationError;

  validating.value = true;
  valid.value = null;
  error.value = '';

  try {
    let allValid = true;
    const errors = [];

    // Determine validation targets based on selected server type
    const serverIds = selectedServers.value;
    const serverType = selectedServerType.value;

    let targetServers = [];
    if (serverType === 'remote') {
      if (isSource) {
        // Source path validated locally when backing up from local -> remote
        targetServers = [null];
      } else {
        // Destination path validated on selected remote server(s)
        targetServers = serverIds;
      }
    } else {
      // Local server or not selected yet -> validate locally
      targetServers = [null];
    }

    // If destination for remote but no server selected yet, postpone validation gracefully
    if (!isSource && serverType === 'remote' && serverIds.length === 0) {
      validating.value = false;
      valid.value = null;
      error.value = 'Select a remote server to validate destination path';
      return;
    }

    for (const serverId of targetServers) {
      try {
        // serverId = null means local path
        const res = await validateServerPath(serverId, path);
        if (!res.valid) {
          allValid = false;
          errors.push(serverId ? `${getServerName(serverId)}: ${res.message}` : `Local: ${res.message}`);
        }
      } catch (e) {
        allValid = false;
        errors.push(serverId ? `${getServerName(serverId)}: ${e.message || 'Validation failed'}` : `Local: ${e.message || 'Validation failed'}`);
      }
    }

    valid.value = allValid;
    if (!allValid) error.value = errors.join('; ');

  } catch (err) {
    valid.value = false;
    error.value = err.message || 'Failed to validate path';
  } finally {
    validating.value = false;
  }
};



const validateServerConnections = async () => {
  serverConnectionErrors.value = []

  const selectedServers = props.singleServer
    ? (formData.server_id !== null ? [formData.server_id] : [])
    : formData.server_ids

  for (const serverId of selectedServers) {
    try {
      const result = await testServerConnection(serverId)
      if (!result.success) {
        serverConnectionErrors.value.push({
          serverId,
          error: result.message
        })
      }
    } catch (err) {
      serverConnectionErrors.value.push({
        serverId,
        error: err.message || 'Connection test failed'
      })
    }
  }
}

const handleSubmit = async () => {
  if (!isFormValid.value) return

  // Validate paths before submission
  await validatePath(formData.source, true)
  await validatePath(formData.destination, false)

  // Check server connections for remote servers
  await validateServerConnections()

  // If there are validation errors, don't submit
  if (sourceValid.value === false || destinationValid.value === false || serverConnectionErrors.value.length > 0) {
    return
  }

  const payload = {
    ...formData,
    server_ids: props.singleServer
      ? (formData.server_id !== null ? [formData.server_id] : [])
      : formData.server_ids
  }
  delete payload.server_id

  emit('submit', payload)
}

onMounted(() => {
  loadServers()
})
</script>
