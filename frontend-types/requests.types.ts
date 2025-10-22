/**
 * API request types
 * Generated from Go backend request models
 */

import { Base } from './base.types';

/**
 * LoginRequest represents admin login credentials
 */
export interface LoginRequest extends Base {
  username: string;  // Min 3, max 50 characters
  password: string;  // Min 8 characters
}

/**
 * RefreshTokenRequest represents token refresh request
 */
export interface RefreshTokenRequest extends Base {
  refreshToken: string;  // Valid refresh token
}

/**
 * LogoutRequest represents logout request
 */
export interface LogoutRequest extends Base {
  refreshToken: string;  // Refresh token to revoke
}

/**
 * CreateAdminRequest represents admin user creation request
 */
export interface CreateAdminRequest extends Base {
  username: string;                    // Min 3, max 50 characters
  email: string;                       // Valid email, max 100 characters
  password: string;                    // Min 8 characters
  name: string;                        // Min 3, max 100 characters
  role: string;                        // Min 2, max 50 characters (e.g., super_admin, admin, moderator)
  permissions?: string[];              // Array of permission strings
  claims?: Record<string, any>;       // Custom claims as key-value pairs
}
