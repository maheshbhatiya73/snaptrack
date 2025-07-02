// libs/api.tsx
const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

type LoginResponse = {
  message: any
  success: any;
  token: string;
  user: {
    username: string;
    role: string;
  };
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
