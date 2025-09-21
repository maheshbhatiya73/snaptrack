<template>
  <div class="min-h-screen bg-white flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="bg-white border border-gray-200 rounded-lg shadow-lg p-8 animate-fade-in">
        <div class="text-center mb-8">
          <div class="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-gray-900 mb-2">Welcome to SnapTrack</h1>
          <p class="text-gray-600 text-sm">Secure access to your backup management system</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">
          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Username</label>
            <input
              v-model="username"
              type="text"
              placeholder="Enter your username"
              class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
              required
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium text-gray-700">Password</label>
            <div class="relative">
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Enter your password"
                class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
                required
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 transition-colors"
              >
                <svg v-if="showPassword" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                </svg>
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
              </button>
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors mt-6"
          >
            <svg v-if="isLoading" class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="-ml-1 mr-3 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"/>
            </svg>
            {{ isLoading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>
      </div>
    </div>

    <Toast v-if="toastMessage" :message="toastMessage" :type="toastType" />
  </div>
</template>

<script setup>
// @ts-ignore
definePageMeta({
  layout: 'auth'
})

import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { loginUser, isAuthenticated } from '~/lib/api'
import Toast from '~/components/Toast.vue'

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const toastMessage = ref('')
const toastType = ref('success')
const router = useRouter()

onMounted(() => {
  if (isAuthenticated()) router.replace('/dashboard')
})

function showToast(message, type = 'success') {
  toastMessage.value = message
  toastType.value = type
  setTimeout(() => {
    toastMessage.value = ''
  }, 3000)
}

const handleLogin = async () => {
  if (isLoading.value) return
  if (!username.value.trim() || !password.value.trim()) {
    showToast(!username.value.trim() ? 'Username is required' : 'Password is required', 'error')
    return
  }

  isLoading.value = true

  try {
    const data = await loginUser({ username: username.value, password: password.value })

    const authData = {
      token: data.token,
      user: data.user || { username: username.value },
      timestamp: Date.now()
    }
    localStorage.setItem('snapstack_auth', JSON.stringify(authData))

    showToast(data.message || 'Login successful', 'success')

    setTimeout(() => {
      router.push('/dashboard')
    }, 1000)

  } catch (err) {
    showToast(err.message || 'Login failed', 'error')
  } finally {
    isLoading.value = false
  }
}
</script>
