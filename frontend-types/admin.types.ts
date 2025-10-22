/**
 * Admin authentication types
 * Generated from Go backend domain models
 */

import { Data } from './base.types';

/**
 * AdminUser represents an administrative user
 */
export interface AdminUser {
  id: number;
  username: string;
  email: string;
  name: string;
  role: string;
  permissions: string[];
  claims: Data;
  isActive: boolean;
  isLocked: boolean;
  failedAttempts: number;
  lastLogin?: string;       // ISO 8601 timestamp
  lastLoginIp?: string;
  createdAt: string;        // ISO 8601 timestamp
  updatedAt: string;        // ISO 8601 timestamp
}

/**
 * AdminUserInfo represents safe admin user info for responses
 */
export interface AdminUserInfo {
  id: number;
  username: string;
  email: string;
  name: string;
  role: string;
  permissions: string[];
  claims?: Data;
}

/**
 * TokenPairResponse represents the response with token pair
 */
export interface TokenPairResponse {
  accessToken: string;
  refreshToken: string;
  tokenType: string;        // Usually "Bearer"
  expiresIn: number;        // Seconds until expiration
  expiresAt: string;        // ISO 8601 timestamp
  user: AdminUserInfo | null;
}

/**
 * TokenClaims represents validated token claims
 */
export interface TokenClaims {
  userId: number;
  username: string;
  email: string;
  name: string;
  role: string;
  permissions: string[];
  claims: Data;
}

/**
 * RefreshToken represents a JWT refresh token
 */
export interface RefreshToken {
  id: number;
  adminId: number;
  token: string;
  tokenFamily: string;
  userAgent?: string;
  ipAddress?: string;
  expiresAt: string;        // ISO 8601 timestamp
  isRevoked: boolean;
  revokedAt?: string;       // ISO 8601 timestamp
  revokedReason?: string;
  createdAt: string;        // ISO 8601 timestamp
}

/**
 * APIKey represents an API key for external integrations
 */
export interface APIKey {
  id: number;
  name: string;
  value: string;
  type: string;
  claims: Data;
  rateLimit: number;
  allowedIps: string[];
  permissions: string[];
  isActive: boolean;
  expiresAt?: string;       // ISO 8601 timestamp
  lastUsedAt?: string;      // ISO 8601 timestamp
  createdBy?: number;
  createdAt: string;        // ISO 8601 timestamp
  updatedAt: string;        // ISO 8601 timestamp
}

/**
 * AuthLog represents an authentication event log
 */
export interface AuthLog {
  id: number;
  userId?: number;
  username: string;
  action: string;
  status: string;
  ipAddress?: string;
  userAgent?: string;
  details: Data;
  createdAt: string;        // ISO 8601 timestamp
}
