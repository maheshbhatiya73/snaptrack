const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

type LoginResponse = {
  message: any;
  success: any;
  token: string;
  user: {
    username: string;
    role: string;
  };
};

type Backup = {
  id: string;
  app: string;
  type: 'manual' | 'automatic';
  size: string;
  status: 'completed' | 'pending' | 'scheduled' | 'started';
  sourcePath: string;
  destinationPath: string;
  fileType: 'zip' | 'tar.gz';
  schedule: {
    kind: 'one-time' | 'hourly';
    date: string;
  };
  nextRun: string;
  logs: Array<{
    id: string;
    backupId: string;
    status: string;
    message: string;
    createdAt: string;
  }>;
  createdAt: string;
};

type BackupResponse = {
  schedule: any;
  fileType: string;
  destinationPath: string;
  sourcePath: string;
  type: string;
  app: string;
  message: string;
  success: boolean;
  data: Backup;
};

type BackupListResponse = {
  message: string;
  success: boolean;
  data: Backup[];
};

export async function login(username: string, password: string): Promise<LoginResponse> {
  const res = await fetch(`${BASE_URL}/api/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Login failed');
  }

  return res.json();
}

export async function createBackup(backup: any): Promise<BackupResponse> {
  const res = await fetch(`${BASE_URL}/api/backups`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(backup),
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Failed to create backup');
  }

  return res.json();
}

export async function updateBackup(id: any, backup: any): Promise<BackupResponse> {
  const res = await fetch(`${BASE_URL}/api/backups/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(backup),
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Failed to update backup');
  }

  return res.json();
}

export async function deleteBackup(id: string): Promise<{ message: string; success: boolean }> {
  const res = await fetch(`${BASE_URL}/api/backups/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Failed to delete backup');
  }

  return res.json();
}

export async function getAllBackups(): Promise<BackupListResponse> {
  const res = await fetch(`${BASE_URL}/api/backups`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Failed to fetch backups');
  }

  return res.json();
}

export async function getBackupById(id: any): Promise<BackupResponse> {
  const res = await fetch(`${BASE_URL}/api/backups/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Failed to fetch backup');
  }

  return res.json();
}
