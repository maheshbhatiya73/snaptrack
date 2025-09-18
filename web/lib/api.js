const API_BASE = 'http://localhost:8080/api'

export async function loginUser({ username, password }) {
  const res = await fetch(`${API_BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  })

  const data = await res.json().catch(() => ({}))

  if (!res.ok) {
    throw new Error(data.message || 'Login failed')
  }

  return data 
}

export function logoutUser() {
  localStorage.removeItem('snapstack_auth')
}

export function getAuthData() {
  if (typeof window === 'undefined') return null
  return JSON.parse(localStorage.getItem('snapstack_auth') || 'null')
}

export function isAuthenticated() {
  const authData = getAuthData()
  return authData && authData.token && authData.user
}
