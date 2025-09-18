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
          <div>
            <label for="source" class="block text-sm font-medium text-slate-700 mb-2">
              Source Path *
            </label>
            <input
              id="source"
              v-model="formData.source"
              type="text"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="/var/www"
            />
          </div>

          <div>
            <label for="destination" class="block text-sm font-medium text-slate-700 mb-2">
              Destination Path *
            </label>
            <input
              id="destination"
              v-model="formData.destination"
              type="text"
              required
              class="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="/mnt/backups/backup.tar.gz"
            />
          </div>

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
            <MultiSelectDropdown
              v-else
              v-model="formData.server_ids"
              :options="serverOptions"
              placeholder="Select target servers"
              searchable
            />
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
import { fetchServers } from '~/lib/api'
import MultiSelectDropdown from '~/components/MultiSelectDropdown.vue'

const props = defineProps({
  backup: {
    type: Object,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])

const isEdit = !!props.backup
const loadingServers = ref(false)
const servers = ref([])

const formData = reactive({
  name: '',
  type: 'full',
  source: '',
  destination: '',
  file_type: 'tar',
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

const isFormValid = computed(() => {
  return formData.name.trim() !== '' &&
         formData.source.trim() !== '' &&
         formData.destination.trim() !== '' &&
         formData.server_ids.length > 0
})

watch(() => props.backup, (newBackup) => {
  if (newBackup) {
    Object.assign(formData, {
      name: newBackup.name || '',
      type: newBackup.type || 'full',
      source: newBackup.source || '',
      destination: newBackup.destination || '',
      file_type: newBackup.file_type || 'tar',
      server_ids: newBackup.server_ids || [],
      schedule_type: newBackup.schedule_type || 'one_time'
    })
  }
}, { immediate: true })

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

const handleSubmit = () => {
  if (!isFormValid.value) return
  emit('submit', formData)
}

onMounted(() => {
  loadServers()
})
</script>
