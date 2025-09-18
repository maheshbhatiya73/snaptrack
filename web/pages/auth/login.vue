<template>
  <div class="h-screen w-screen flex items-center justify-center bg-gradient-to-br from-slate-50 via-white to-slate-100 p-4">
    <div class="w-full max-w-md">
      <div class="bg-white/80 backdrop-blur-sm border border-white/20 rounded-2xl shadow-2xl p-8">
        <div class="text-center mb-8">
          <img src="../../assets/images/logo.png" class="mx-auto w-24 h-24" />
          <p class="text-slate-600 mt-4 text-sm">Secure access to your backup system</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">
          <div class="space-y-2">
            <label class="text-sm font-medium text-slate-700">Username</label>
            <input
              v-model="username"
              type="text"
              placeholder="Enter username"
              class="w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium text-slate-700">Password</label>
            <div class="relative">
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Password"
                class="w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
              <button type="button" @click="showPassword = !showPassword" class="absolute inset-y-0 right-0 px-3 text-slate-400">
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3 bg-blue-600 cursor-pointer text-white font-semibold rounded-xl shadow-lg hover:bg-blue-700 disabled:opacity-50 transition-all"
          >
            {{ isLoading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>
      </div>
    </div>

    <Toast v-if="toastMessage" :message="toastMessage" :type="toastType" />
  </div>
</template>

<script setup>
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
