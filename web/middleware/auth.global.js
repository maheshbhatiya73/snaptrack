export default defineNuxtRouteMiddleware((to) => {
  if (!process.client) return

  const authData = JSON.parse(localStorage.getItem('snapstack_auth') || 'null')

  const isAuthenticated = authData && authData.token && authData.user

  if (!isAuthenticated && !to.path.startsWith('/auth')) {
    return navigateTo('/auth/login')
  }

  if (isAuthenticated && to.path.startsWith('/auth')) {
    return navigateTo('/dashboard')
  }
})
