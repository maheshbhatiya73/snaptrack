const API_BASE = `${process.env.NUXT_PUBLIC_BACKEND_URL || 'http://localhost:8080'}/api`

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

export async function fetchBackups() {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch backups')
  }

  return res.json()
}

export async function fetchBackup(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/${id}`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch backup')
  }

  return res.json()
}

export async function createBackup(backupData) {
  const authData = getAuthData()
  
  const payload = {
    ...backupData,
    server_ids: Array.isArray(backupData.server_ids) ? backupData.server_ids : []
  }
  
  const res = await fetch(`${API_BASE}/backups`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to create backup')
  }

  return res.json()
}

export async function updateBackup(id, backupData) {
  const authData = getAuthData()
  
  const payload = {
    ...backupData,
    server_ids: Array.isArray(backupData.server_ids) ? backupData.server_ids : []
  }
  
  const res = await fetch(`${API_BASE}/backups/${id}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to update backup')
  }

  return res.json()
}

export async function deleteBackup(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to delete backup')
  }

  // For DELETE requests, server may return 204 No Content with no body
  if (res.status === 204) {
    return { success: true }
  }

  return res.json().catch(() => ({ success: true }))
}

export async function executeBackup(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/${id}/execute`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to execute backup')
  }

  return res.json()
}

export async function fetchServers() {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/servers`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch servers')
  }

  return res.json()
}

export async function fetchServer(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/servers/${id}`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch server')
  }

  return res.json()
}
export async function createServer(serverData) {
  const authData = getAuthData();

  try {
    const res = await fetch(`${API_BASE}/servers`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authData.token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(serverData),
    });

    // Always parse JSON safely
    const data = await res.json().catch(() => ({}));

    // Instead of throwing -> return { success: false }
    if (!res.ok) {
      return {
        success: false,
        message: data.message || data.error || "Failed to create server",
      };
    }

    return {
      success: true,
      data,
    };
  } catch (err) {
    return {
      success: false,
      message: err.message || "Network error",
    };
  }
}


export async function updateServer(id, serverData) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/servers/${id}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(serverData)
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to update server')
  }

  return res.json()
}

export async function deleteServer(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/servers/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to delete server')
  }

  // For DELETE requests, server may return 204 No Content with no body
  if (res.status === 204) {
    return { success: true }
  }

  return res.json().catch(() => ({ success: true }))
}

export async function testServerConnection(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/servers/${id}/test`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to test server connection')
  }

  return res.json()
}

export async function fetchDashboardStats() {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/dashboard/stats`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch dashboard stats')
  }

  return res.json()
}

export async function fetchRecentActivity(limit = 10, offset = 0) {
  try {
    const authData = getAuthData()

    const res = await fetch(`${API_BASE}/dashboard/recent-activity?limit=${limit}&offset=${offset}`, {
      headers: {
        'Authorization': `Bearer ${authData?.token ?? ''}`,
        'Content-Type': 'application/json'
      }
    })

    if (!res.ok) return []

    return res.json()
  } catch (err) {
    console.error('Failed to fetch recent activity:', err)
    return []
  }
}


export async function fetchSystemStatus() {
  try {
    const authData = getAuthData()

    const res = await fetch(`${API_BASE}/`, {
      headers: {
        'Authorization': `Bearer ${authData?.token ?? ''}`,
        'Content-Type': 'application/json'
      }
    })

    if (!res.ok) {
      return { status: 'offline' }
    }

    const data = await res.json()

    // if API responds with expected message → online
    if (data?.message) {
      return { status: 'online' }
    }

    return { status: 'offline' }
  } catch (err) {
    // Network failure or server down
    return { status: 'offline' }
  }
}


export async function fetchServerStatuses() {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/servers/status`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch server statuses')
  }

  return res.json()
}

export async function checkServerNameExists(name, excludeId = null) {
  const authData = getAuthData()
  const servers = await fetchServers()

  // Check if any server has the same name (excluding the current server if editing)
  return servers.some(server => server.name === name && server.id !== excludeId)
}

export async function validateServerPath(serverId, path) {
  const authData = getAuthData()
  // If serverId is null, validate locally on the backend host
  const url = serverId == null ? `${API_BASE}/local/validate-path` : `${API_BASE}/servers/${serverId}/validate-path`
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ path })
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to validate path')
  }

  return res.json()
}

export async function fetchRunningBackups() {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/processes/running`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch running backups')
  }

  return res.json()
}

export async function fetchBackupProgress(backupId) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/${backupId}/progress`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    throw new Error('Failed to fetch backup progress')
  }

  return res.json()
}

export async function deleteAllProcesses() {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/processes`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to delete all processes')
  }

  if (res.status === 204) {
    return { success: true }
  }

  return res.json().catch(() => ({ success: true }))
}

export async function deleteProcess(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/backups/processes/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error.message || 'Failed to delete process')
  }

  if (res.status === 204) {
    return { success: true }
  }

  return res.json().catch(() => ({ success: true }))
}

// Monitor helpers
export async function fetchMonitorSnapshot(id) {
  const authData = getAuthData()
  const res = await fetch(`${API_BASE}/monitor/${id}/snapshot`, {
    headers: {
      'Authorization': `Bearer ${authData.token}`,
      'Content-Type': 'application/json'
    }
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.message || 'Failed to fetch snapshot')
  }
  return res.json()
}

export function monitorWsUrl(id) {
  const httpBase = process.env.NUXT_PUBLIC_BACKEND_URL || 'http://localhost:8080'
  const wsBase = httpBase.replace(/^http/,'ws')
  return `${wsBase}/api/monitor/${id}/ws`
}

export function monitorBatchWsUrl() {
  const httpBase = process.env.NUXT_PUBLIC_BACKEND_URL || 'http://localhost:8080'
  const wsBase = httpBase.replace(/^http/,'ws')
  return `${wsBase}/api/monitor/ws`
}
