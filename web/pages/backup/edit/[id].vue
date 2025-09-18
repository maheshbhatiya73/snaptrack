<template>
  <div class="min-h-screen bg-slate-50">
    <!-- Page Header -->
    <div class="bg-white border-b border-slate-200">
      <div class="px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button 
              @click="goBack"
              class="p-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <div>
              <h1 class="text-2xl font-bold text-slate-900">Edit Backup</h1>
              <p class="mt-1 text-sm text-slate-600">Modify backup configuration and settings</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <main class="px-6 py-6">
      <div class="max-w-4xl mx-auto">
        <!-- Loading State -->
        <div v-if="loading && !backup" class="flex items-center justify-center py-12">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p class="mt-4 text-slate-600">Loading backup...</p>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6 mb-6">
          <div class="flex items-center">
            <svg class="w-5 h-5 text-red-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div>
              <h3 class="text-sm font-medium text-red-800">Error loading backup</h3>
              <p class="text-sm text-red-700 mt-1">{{ error }}</p>
            </div>
          </div>
        </div>

        <!-- Success Message -->
        <div v-if="success" class="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
          <div class="flex items-center">
            <svg class="w-5 h-5 text-green-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div>
              <h3 class="text-sm font-medium text-green-800">Success</h3>
              <p class="text-sm text-green-700 mt-1">{{ success }}</p>
            </div>
          </div>
        </div>

        <!-- Backup Form -->
        <BackupForm
          v-if="backup"
          :backup="backup"
          :loading="loading"
          @submit="handleSubmit"
          @cancel="goBack"
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { isAuthenticated, fetchBackup, updateBackup } from '~/lib/api'
import BackupForm from '~/components/BackupForm.vue'

const router = useRouter()
const route = useRoute()

// Reactive data
const backup = ref(null)
const loading = ref(false)
const error = ref(null)
const success = ref(null)

// Methods
const goBack = () => {
  router.push('/backup')
}

const loadBackup = async () => {
  try {
    loading.value = true
    error.value = null
    
    const backupId = route.params.id
    backup.value = await fetchBackup(backupId)
    
  } catch (err) {
    error.value = err.message
    console.error('Failed to load backup:', err)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async (formData) => {
  try {
    loading.value = true
    error.value = null
    success.value = null

    const backupId = route.params.id
    await updateBackup(backupId, formData)
    
    success.value = 'Backup updated successfully!'
    
    // Redirect to backup list after a short delay
    setTimeout(() => {
      router.push('/backup')
    }, 2000)
    
  } catch (err) {
    error.value = err.message
    console.error('Failed to update backup:', err)
  } finally {
    loading.value = false
  }
}

// Check authentication and load backup
onMounted(() => {
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  loadBackup()
})
</script>
