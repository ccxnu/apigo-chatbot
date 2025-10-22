# API Chatbot Frontend Types & Client

TypeScript types and API client generated from the Go backend for seamless frontend integration.

## Installation

### Option 1: Copy to React Project

Copy the `frontend-types` directory to your React project:

```bash
# From your React project root
cp -r /path/to/apigo-chatbot/frontend-types ./src/api
```

### Option 2: Use as NPM Package (Recommended)

Link the package locally:

```bash
# In the frontend-types directory
npm link

# In your React project
npm link @api-chatbot/types
```

### Option 3: Install from Git

Add to your `package.json`:

```json
{
  "dependencies": {
    "@api-chatbot/types": "file:../apigo-chatbot/frontend-types"
  }
}
```

## Usage

### 1. Basic API Client Setup

```typescript
// src/config/api.ts
import { ApiClient, ApiError } from '@api-chatbot/types';

export const api = new ApiClient({
  baseUrl: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  timeout: 30000,

  // Get access token from your auth store
  getAccessToken: () => {
    return localStorage.getItem('accessToken');
  },

  // Handle unauthorized (401) responses
  onUnauthorized: () => {
    // Clear tokens and redirect to login
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    window.location.href = '/login';
  },

  // Global error handler
  onError: (error: ApiError) => {
    console.error('API Error:', error.code, error.info);
    // You can show toast notifications here
  },
});
```

### 2. React Hook for Authentication

```typescript
// src/hooks/useAuth.ts
import { useState } from 'react';
import { api } from '../config/api';
import { TokenPairResponse, ApiError } from '@api-chatbot/types';

export function useAuth() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (username: string, password: string) => {
    setLoading(true);
    setError(null);

    try {
      const result = await api.login(username, password);

      if (result.success && result.data) {
        // Store tokens
        localStorage.setItem('accessToken', result.data.accessToken);
        localStorage.setItem('refreshToken', result.data.refreshToken);

        // Store user info
        localStorage.setItem('user', JSON.stringify(result.data.user));

        return result.data;
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.info);
      } else {
        setError('An unexpected error occurred');
      }
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    const refreshToken = localStorage.getItem('refreshToken');
    if (!refreshToken) return;

    try {
      await api.logout(refreshToken);
    } catch (err) {
      console.error('Logout error:', err);
    } finally {
      // Clear local storage regardless of API call success
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('user');
    }
  };

  const refreshToken = async () => {
    const refreshToken = localStorage.getItem('refreshToken');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const result = await api.refreshToken(refreshToken);

      if (result.success && result.data) {
        localStorage.setItem('accessToken', result.data.accessToken);
        localStorage.setItem('refreshToken', result.data.refreshToken);
        return result.data;
      }
    } catch (err) {
      // If refresh fails, clear everything and redirect to login
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('user');
      throw err;
    }
  };

  return {
    login,
    logout,
    refreshToken,
    loading,
    error,
  };
}
```

### 3. Login Component Example

```typescript
// src/components/LoginForm.tsx
import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { useNavigate } from 'react-router-dom';

export function LoginForm() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { login, loading, error } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const result = await login(username, password);
      console.log('Logged in:', result.user);
      navigate('/dashboard');
    } catch (err) {
      // Error is already handled by useAuth hook
      console.error('Login failed');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Admin Login</h2>

      {error && <div className="error">{error}</div>}

      <div>
        <label>Username</label>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          minLength={3}
          maxLength={50}
        />
      </div>

      <div>
        <label>Password</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={8}
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? 'Logging in...' : 'Login'}
      </button>
    </form>
  );
}
```

### 4. Automatic Token Refresh with Axios Interceptor

```typescript
// src/config/axios-setup.ts
import axios from 'axios';
import { api } from './api';

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

// Setup axios interceptor for automatic token refresh
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers['Authorization'] = 'Bearer ' + token;
          return axios(originalRequest);
        }).catch(err => {
          return Promise.reject(err);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = localStorage.getItem('refreshToken');

      if (!refreshToken) {
        processQueue(error, null);
        isRefreshing = false;
        return Promise.reject(error);
      }

      try {
        const result = await api.refreshToken(refreshToken);

        if (result.success && result.data) {
          const newAccessToken = result.data.accessToken;
          localStorage.setItem('accessToken', newAccessToken);
          localStorage.setItem('refreshToken', result.data.refreshToken);

          axios.defaults.headers.common['Authorization'] = 'Bearer ' + newAccessToken;
          originalRequest.headers['Authorization'] = 'Bearer ' + newAccessToken;

          processQueue(null, newAccessToken);
          isRefreshing = false;

          return axios(originalRequest);
        }
      } catch (err) {
        processQueue(err, null);
        isRefreshing = false;

        // Clear tokens and redirect to login
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        window.location.href = '/login';

        return Promise.reject(err);
      }
    }

    return Promise.reject(error);
  }
);
```

### 5. Auth Context Provider

```typescript
// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, useEffect } from 'react';
import { AdminUserInfo } from '@api-chatbot/types';

interface AuthContextType {
  user: AdminUserInfo | null;
  isAuthenticated: boolean;
  setUser: (user: AdminUserInfo | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AdminUserInfo | null>(null);

  useEffect(() => {
    // Load user from localStorage on mount
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);

  const isAuthenticated = user !== null && !!localStorage.getItem('accessToken');

  return (
    <AuthContext.Provider value={{ user, isAuthenticated, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuthContext must be used within AuthProvider');
  }
  return context;
}
```

### 6. Protected Route Component

```typescript
// src/components/ProtectedRoute.tsx
import { Navigate } from 'react-router-dom';
import { useAuthContext } from '../contexts/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredPermissions?: string[];
}

export function ProtectedRoute({ children, requiredPermissions }: ProtectedRouteProps) {
  const { isAuthenticated, user } = useAuthContext();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Check permissions if specified
  if (requiredPermissions && requiredPermissions.length > 0) {
    const hasPermission = requiredPermissions.every(permission =>
      user?.permissions.includes(permission)
    );

    if (!hasPermission) {
      return <Navigate to="/unauthorized" replace />;
    }
  }

  return <>{children}</>;
}
```

### 7. Complete App Setup

```typescript
// src/App.tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { LoginForm } from './components/LoginForm';
import { Dashboard } from './pages/Dashboard';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginForm />} />

          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin"
            element={
              <ProtectedRoute requiredPermissions={['admin.manage']}>
                <AdminPanel />
              </ProtectedRoute>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}
```

## Type Definitions

### Base Request Fields

All API requests automatically include these fields via `generateBaseFields()`:

```typescript
{
  idSession: string;      // Auto-generated session ID
  idRequest: string;      // UUID v4 for each request
  process: string;        // Always "web-frontend"
  idDevice: string;       // Persistent device ID
  deviceAddress: string;  // Set by backend
  dateProcess: string;    // ISO 8601 timestamp
}
```

### API Response Format

All responses follow this structure:

```typescript
{
  success: boolean;
  code: string;
  info: string;
  data: T;
}
```

## Error Handling

The API client throws `ApiError` instances:

```typescript
try {
  await api.login(username, password);
} catch (err) {
  if (err instanceof ApiError) {
    console.error('Error code:', err.code);
    console.error('Error message:', err.info);
    console.error('HTTP status:', err.statusCode);
  }
}
```

## Environment Variables

Create a `.env` file in your React project:

```env
REACT_APP_API_URL=http://localhost:8080
```

## Adding More Endpoints

To add more endpoints to the API client, edit `api-client.ts`:

```typescript
// In the ApiClient class
async getDocuments(): Promise<Result<Document[]>> {
  const request = {
    ...generateBaseFields(),
    // Add any additional parameters
  };

  return this.request<Document[]>(
    '/api/v1/documents/get-all',
    'POST',
    request,
    true // requiresAuth
  );
}
```

## Type Safety

All types are strongly typed and match your Go backend exactly:

- Request validation mirrors Go struct tags
- Response types match Go domain models
- Enum values are type-safe
- Null safety with TypeScript strict mode

## Best Practices

1. Always use the API client instead of direct fetch calls
2. Store tokens securely (consider using httpOnly cookies for production)
3. Implement token refresh before access token expires
4. Handle errors gracefully with user-friendly messages
5. Use TypeScript strict mode for maximum type safety
6. Keep types in sync when backend changes

## License

MIT
