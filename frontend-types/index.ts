/**
 * Main export file for frontend types and API client
 * Import from this file in your React components
 */

// Base types
export type { Base, Result, Data } from './base.types';

// Admin types
export type {
  AdminUser,
  AdminUserInfo,
  TokenPairResponse,
  TokenClaims,
  RefreshToken,
  APIKey,
  AuthLog,
} from './admin.types';

// Request types
export type {
  LoginRequest,
  RefreshTokenRequest,
  LogoutRequest,
  CreateAdminRequest,
} from './requests.types';

// Document types
export type {
  Document,
  CreateDocumentParams,
  UpdateDocumentParams,
} from './document.types';

// Chunk types
export type {
  Chunk,
  ChunkWithSimilarity,
  ChunkWithHybridSimilarity,
  SimilaritySearchParams,
  HybridSearchParams,
  BulkCreateChunksParams,
} from './chunk.types';

// Statistics types
export type {
  ChunkStatistics,
  TopChunkByUsage,
  UpdateChunkQualityMetricsParams,
  UpdateChunkStalenessParams,
} from './statistics.types';

// Parameter types
export type {
  Parameter,
  AddParameterParams,
  UpdateParameterParams,
} from './parameter.types';

// API Client
export {
  ApiClient,
  apiClient,
  generateBaseFields,
  ApiError,
} from './api-client';
export type { ApiClientConfig } from './api-client';
