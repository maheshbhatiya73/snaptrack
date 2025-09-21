<template>
  <div class="fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg border-r border-gray-200">
    <!-- Logo Section -->
    <div class="flex items-center justify-center h-16 px-6 border-b border-gray-200">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 bg-gray-900 rounded-lg flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        </div>
        <div>
          <span class="text-xl font-bold text-gray-900">SnapTrack</span>
          <p class="text-xs text-gray-500">Backup Management</p>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-4 py-6 space-y-2 sidebar-scroll">
      <div class="space-y-1">
        <NuxtLink
          to="/dashboard"
          class="group flex items-center px-3 py-3 text-sm font-medium rounded-lg transition-all duration-200"
          :class="isActive('/dashboard') ? 'bg-gray-900 text-white' : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'"
        >
          <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v6H8V5z"/>
          </svg>
          Dashboard
        </NuxtLink>

        <NuxtLink
          to="/servers"
          class="group flex items-center px-3 py-3 text-sm font-medium rounded-lg transition-all duration-200"
          :class="isActive('/servers') ? 'bg-gray-900 text-white' : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'"
        >
          <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
          </svg>
          Servers
        </NuxtLink>

        <NuxtLink
          to="/backup"
          class="group flex items-center px-3 py-3 text-sm font-medium rounded-lg transition-all duration-200"
          :class="isActive('/backup') ? 'bg-gray-900 text-white' : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'"
        >
          <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
          </svg>
          Backup
        </NuxtLink>
      </div>
    </nav>

    <!-- User Section -->
    <div class="p-4 border-t border-gray-200">
      <div class="flex items-center space-x-3 mb-4">
        <div class="w-10 h-10 bg-gray-900 rounded-lg flex items-center justify-center">
          <span class="text-sm font-semibold text-white">
            {{ userInitials }}
          </span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate">
            {{ userInfo?.username || userInfo?.email || 'User' }}
          </p>
          <p class="text-xs text-gray-500">Administrator</p>
        </div>
      </div>
      <button
        @click="handleLogout"
        class="w-full flex items-center px-3 py-2.5 text-sm font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all duration-200"
      >
        <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
        </svg>
        Sign out
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getAuthData, logoutUser, isAuthenticated } from '~/lib/api'

const router = useRouter()
const route = useRoute()
const userInfo = ref(null)

const userInitials = computed(() => {
  if (!userInfo.value) return 'U'
  const name = userInfo.value.username || userInfo.value.email || 'User'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

const isActive = (path) => {
  return route.path === path
}

onMounted(() => {
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
