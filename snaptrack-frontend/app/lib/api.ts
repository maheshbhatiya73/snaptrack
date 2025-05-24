interface LoginResponse {
  success: boolean;
  message?: string;
  token?: string;
}

interface VerifyResponse {
  valid: boolean;
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  try {
    const response = await fetch("http://localhost:8000/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error(`HTTP error: ${response.status}`);
    }

    const data: LoginResponse = await response.json();
    return data;
  } catch (err) {
    throw new Error("Failed to connect to server");
  }
}

export async function verifyToken(token: string): Promise<boolean> {
  try {
    const response = await fetch("http://localhost:8000/api/verify", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
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