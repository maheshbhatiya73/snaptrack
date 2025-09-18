<template>
  <div class="min-h-screen bg-slate-50">
    <!-- Header -->
    <header class="bg-white border-b border-slate-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center space-x-4">
            <span class="text-sm text-slate-600">Welcome, {{ userInfo?.username || userInfo?.email || 'User' }}</span>
            <button
              @click="handleLogout"
              class="text-slate-600 hover:text-slate-900 text-sm font-medium"
            >
              Sign out
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <!-- Page Title -->
        <div class="mb-8">
          <h1 class="text-3xl font-bold text-slate-900">Dashboard</h1>
          <p class="mt-2 text-slate-600">Monitor your backup system status and manage your data</p>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-green-100 rounded-lg">
                <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">System Status</p>
                <p class="text-2xl font-bold text-green-600">Online</p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-blue-100 rounded-lg">
                <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Total Backups</p>
                <p class="text-2xl font-bold text-blue-600">1,247</p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 rounded-lg">
                <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Storage Used</p>
                <p class="text-2xl font-bold text-purple-600">2.4 TB</p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
            <div class="flex items-center">
              <div class="p-2 bg-amber-100 rounded-lg">
                <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-slate-600">Last Backup</p>
                <p class="text-2xl font-bold text-amber-600">2m ago</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white rounded-xl shadow-sm border border-slate-200">
          <div class="px-6 py-4 border-b border-slate-200">
            <h2 class="text-lg font-semibold text-slate-900">Recent Activity</h2>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="flex items-center space-x-3">
                <div class="w-2 h-2 bg-green-500 rounded-full"></div>
                <div class="flex-1">
                  <p class="text-sm text-slate-900">Backup completed successfully</p>
                  <p class="text-xs text-slate-500">2 minutes ago</p>
                </div>
              </div>
              <div class="flex items-center space-x-3">
                <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
                <div class="flex-1">
                  <p class="text-sm text-slate-900">New backup scheduled</p>
                  <p class="text-xs text-slate-500">15 minutes ago</p>
                </div>
              </div>
              <div class="flex items-center space-x-3">
                <div class="w-2 h-2 bg-amber-500 rounded-full"></div>
                <div class="flex-1">
                  <p class="text-sm text-slate-900">Storage optimization completed</p>
                  <p class="text-xs text-slate-500">1 hour ago</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAuthData, logoutUser, isAuthenticated } from '~/lib/api'

const router = useRouter()
const userInfo = ref(null)

onMounted(() => {
  // Get user info from localStorage
  if (!isAuthenticated()) {
    router.push('/auth/login')
    return
  }
  
  const authData = getAuthData()
  if (authData) {
    userInfo.value = authData.user
  }
})

const handleLogout = () => {
  logoutUser()
  router.push('/auth/login')
}
</script>
