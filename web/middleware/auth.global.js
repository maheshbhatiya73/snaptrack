// auth.global.js
export default defineNuxtRouteMiddleware((to) => {
  // Only run on client
  if (!process.client) return

  const authData = JSON.parse(localStorage.getItem('snapstack_auth') || 'null')

  // Check if auth data exists and is valid
  const isAuthenticated = authData && authData.token && authData.user

  // If not authenticated, redirect to login
  if (!isAuthenticated && !to.path.startsWith('/auth')) {
    return navigateTo('/auth/login')
  }

  // If authenticated and trying to access login, redirect to dashboard
  if (isAuthenticated && to.path.startsWith('/auth')) {
    return navigateTo('/dashboard')
  }
})
