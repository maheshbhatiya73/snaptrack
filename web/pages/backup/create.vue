<template>
  <div class="min-h-screen bg-slate-50 flex flex-col">
    <header class="bg-white border-b border-slate-200">
      <div class="px-6 py-4 flex items-center justify-between">
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
            <h1 class="text-2xl font-bold text-slate-900">Create New Backup</h1>
            <p class="mt-1 text-sm text-slate-600">Configure and schedule a new backup job</p>
          </div>
        </div>
      </div>
    </header>

    <main class="flex-1 px-6 py-6">
      <div class="max-w-6xl mx-auto">
        <transition name="fade">
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-red-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <div>
                <h3 class="text-sm font-medium text-red-800">Error</h3>
                <p class="text-sm text-red-700 mt-1">{{ error }}</p>
              </div>
            </div>
          </div>
        </transition>

        <transition name="fade">
          <div v-if="success" class="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-green-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <div>
                <h3 class="text-sm font-medium text-green-800">Success</h3>
                <p class="text-sm text-green-700 mt-1">{{ success }}</p>
              </div>
            </div>
          </div>
        </transition>

        <BackupForm
          :loading="loading"
          @submit="handleSubmit"
          @cancel="goBack"
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { isAuthenticated, createBackup } from '~/lib/api'
import BackupForm from '~/components/BackupForm.vue'

const router = useRouter()

const loading = ref(false)
const error = ref(null)
const success = ref(null)

const goBack = () => router.push('/backup')

const handleSubmit = async (formData) => {
  try {
    loading.value = true
    error.value = null
    success.value = null

    await createBackup(formData)

    success.value = 'Backup created successfully!'
    setTimeout(() => {
      router.push('/backup')
    }, 2000)

  } catch (err) {
    error.value = err?.message || 'Failed to create backup.'
    console.error('Backup error:', err)
  } finally {
    loading.value = false
  }
}

if (!isAuthenticated()) {
  router.push('/auth/login')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
