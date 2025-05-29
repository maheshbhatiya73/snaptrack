// Types
type FileType = 'zip' | 'tar' | 'tar.gz';

export type Schedule = {
  kind: 'one-time' | 'hourly' | 'daily' | 'weekly' | 'monthly';
  date?: string;
  time?: string;
  dayOfWeek?: string;
  dayOfMonth?: number;
}

export type Backup = {
  id: string;
  app: string;
  type: 'full' | 'incremental' | 'manual'; // Added 'manual' to support provided data
  size: string;
  date: string;
  status: 'success' | 'failed' | 'pending';
  sourcePath: string;
  destinationPath: string;
  fileType: FileType;
  schedule: Schedule;
}

interface LoginResponse {
  success: boolean;
  message?: string;
  token?: string;
}

interface VerifyResponse {
  valid: boolean;
}

interface BackupResponse {
  success: boolean;
  message?: string;
  data?: Backup | Backup[];
  total?: number;
}

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export async function login(username: string, password: string): Promise<LoginResponse> {
  try {
    const response = await fetch(`${BASE_URL}/api/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
      credentials: 'include',
    });

    if (!response.ok) {
      const errorData = await response.json();
      return { success: false, message: errorData.message || `HTTP error: ${response.status}` };
    }

    const data: LoginResponse = await response.json();
    return data;
  } catch (err) {
    return { success: false, message: 'Failed to connect to server' };
  }
}

export async function verifyToken(token: string): Promise<boolean> {
  try {
    const response = await fetch(`${BASE_URL}/api/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      credentials: 'include',
    });

    if (!response.ok) {
      return false;
    }

    const data: VerifyResponse = await response.json();
    return data.valid;
  } catch (err) {
    return false;
  }
}

export async function createBackup(backup: Partial<Backup>, token: string): Promise<BackupResponse> {
  try {
    const response = await fetch(`${BASE_URL}/api/backups`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(backup),
      credentials: 'include',
    });

    if (!response.ok) {
      const errorData = await response.json();
      return { success: false, message: errorData.message || `HTTP error: ${response.status}` };
    }

    const data: Backup = await response.json();
    return { success: true, data };
  } catch (err) {
    return { success: false, message: 'Failed to create backup' };
  }
}

export async function getBackups(token: string, page: number, limit: number): Promise<BackupResponse> {
  try {
    const response = await fetch(`${BASE_URL}/api/backups`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      credentials: 'include',
    });

    if (!response.ok) {
      const errorData = await response.json();
      return { success: false, message: errorData.message || `HTTP error: ${response.status}` };
    }

    const data: Backup[] = await response.json();
    const start = (page - 1) * limit;
    const paginatedData = data.slice(start, start + limit);
    return { success: true, data: paginatedData, total: data.length };
  } catch (err) {
    return { success: false, message: 'Failed to fetch backups' };
  }
}

export async function updateBackup(id: string, backup: Partial<Backup>, token: string): Promise<BackupResponse> {
  try {
    console.log(JSON.stringify(backup))
    const response = await fetch(`${BASE_URL}/api/backups/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      
      body: JSON.stringify(backup),
    });

    if (!response.ok) {
      const errorData = await response.json();
      return { success: false, message: errorData.message || `HTTP error: ${response.status}` };
    }

    const data: Backup = await response.json();
    return { success: true, data };
  } catch (err) {
    return { success: false, message: 'Failed to update backup' };
  }
}

export async function deleteBackup(id: string, token: string): Promise<BackupResponse> {
  try {
    const response = await fetch(`${BASE_URL}/api/backups/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      credentials: 'include',
    });

    if (!response.ok) {
      const errorData = await response.json();
      return { success: false, message: errorData.message || `HTTP error: ${response.status}` };
    }

    return { success: true };
  } catch (err) {
    return { success: false, message: 'Failed to delete backup' };
  }
}