/**
 * Type-safe API client for React frontend
 * Automatically handles Base fields and API communication
 */

import { Base, Result, Data } from './base.types';
import {
  LoginRequest,
  RefreshTokenRequest,
  LogoutRequest,
  CreateAdminRequest,
} from './requests.types';
import {
  TokenPairResponse,
  AdminUser,
} from './admin.types';
import {
  Document,
  CreateDocumentParams,
  UpdateDocumentParams,
} from './document.types';
import {
  Chunk,
  ChunkWithSimilarity,
  ChunkWithHybridSimilarity,
  SimilaritySearchParams,
  HybridSearchParams,
  BulkCreateChunksParams,
} from './chunk.types';
import {
  ChunkStatistics,
  TopChunkByUsage,
  UpdateChunkQualityMetricsParams,
  UpdateChunkStalenessParams,
} from './statistics.types';
import {
  Parameter,
  AddParameterParams,
  UpdateParameterParams,
} from './parameter.types';

/**
 * API Client configuration
 */
export interface ApiClientConfig {
  baseUrl: string;
  timeout?: number;
  getAccessToken?: () => string | null;
  onUnauthorized?: () => void;
  onError?: (error: ApiError) => void;
}

/**
 * API Error class
 */
export class ApiError extends Error {
  constructor(
    public code: string,
    public info: string,
    public statusCode?: number,
    public data?: any
  ) {
    super(info);
    this.name = 'ApiError';
  }
}

/**
 * Generate Base fields for requests
 */
export function generateBaseFields(): Base {
  return {
    idSession: sessionStorage.getItem('sessionId') || crypto.randomUUID(),
    idRequest: crypto.randomUUID(),
    process: 'web-frontend',
    idDevice: getDeviceId(),
    deviceAddress: '0.0.0.0', // Will be set by backend from actual IP
    dateProcess: new Date().toISOString(),
  };
}

/**
 * Get or create device ID
 */
function getDeviceId(): string {
  let deviceId = localStorage.getItem('deviceId');
  if (!deviceId) {
    deviceId = crypto.randomUUID();
    localStorage.setItem('deviceId', deviceId);
  }
  return deviceId;
}

/**
 * Main API Client class
 */
export class ApiClient {
  private config: Required<ApiClientConfig>;

  constructor(config: ApiClientConfig) {
    this.config = {
      timeout: 30000,
      getAccessToken: () => null,
      onUnauthorized: () => {},
      onError: () => {},
      ...config,
    };
  }

  /**
   * Generic request method
   */
  private async request<T>(
    endpoint: string,
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' = 'POST',
    body?: any,
    requiresAuth: boolean = false
  ): Promise<Result<T>> {
    const url = `${this.config.baseUrl}${endpoint}`;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (requiresAuth) {
      const token = this.config.getAccessToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.config.timeout);

    try {
      const response = await fetch(url, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      const result: Result<T> = await response.json();

      // Handle unsuccessful responses
      if (!response.ok) {
        if (response.status === 401) {
          this.config.onUnauthorized();
        }

        const error = new ApiError(
          result.code || 'UNKNOWN_ERROR',
          result.info || 'An error occurred',
          response.status,
          result.data
        );

        this.config.onError(error);
        throw error;
      }

      // Handle application-level errors
      if (!result.success) {
        const error = new ApiError(
          result.code,
          result.info,
          response.status,
          result.data
        );

        this.config.onError(error);
        throw error;
      }

      return result;
    } catch (error) {
      clearTimeout(timeoutId);

      if (error instanceof ApiError) {
        throw error;
      }

      if (error instanceof Error) {
        if (error.name === 'AbortError') {
          throw new ApiError('TIMEOUT', 'Request timeout');
        }
        throw new ApiError('NETWORK_ERROR', error.message);
      }

      throw new ApiError('UNKNOWN_ERROR', 'An unknown error occurred');
    }
  }

  /**
   * Admin Authentication Endpoints
   */

  async login(username: string, password: string): Promise<Result<TokenPairResponse>> {
    const request: LoginRequest = {
      ...generateBaseFields(),
      username,
      password,
    };

    return this.request<TokenPairResponse>('/admin/auth/login', 'POST', request);
  }

  async refreshToken(refreshToken: string): Promise<Result<TokenPairResponse>> {
    const request: RefreshTokenRequest = {
      ...generateBaseFields(),
      refreshToken,
    };

    return this.request<TokenPairResponse>('/admin/auth/refresh', 'POST', request);
  }

  async logout(refreshToken: string): Promise<Result<Data>> {
    const request: LogoutRequest = {
      ...generateBaseFields(),
      refreshToken,
    };

    return this.request<Data>('/admin/auth/logout', 'POST', request);
  }

  async createAdmin(params: Omit<CreateAdminRequest, keyof Base>): Promise<Result<AdminUser>> {
    const request: CreateAdminRequest = {
      ...generateBaseFields(),
      ...params,
    };

    return this.request<AdminUser>('/admin/users/create', 'POST', request, true);
  }

  /**
   * Document Endpoints
   */

  async getAllDocuments(limit: number = 100, offset: number = 0): Promise<Result<Document[]>> {
    return this.request<Document[]>('/api/v1/documents/get-all', 'POST', {
      ...generateBaseFields(),
      limit,
      offset,
    }, true);
  }

  async getDocumentById(docId: number): Promise<Result<Document>> {
    return this.request<Document>('/api/v1/documents/get-by-id', 'POST', {
      ...generateBaseFields(),
      docId,
    }, true);
  }

  async getDocumentsByCategory(category: string, limit: number = 100, offset: number = 0): Promise<Result<Document[]>> {
    return this.request<Document[]>('/api/v1/documents/get-by-category', 'POST', {
      ...generateBaseFields(),
      category,
      limit,
      offset,
    }, true);
  }

  async searchDocumentsByTitle(titlePattern: string, limit: number = 100): Promise<Result<Document[]>> {
    return this.request<Document[]>('/api/v1/documents/search-by-title', 'POST', {
      ...generateBaseFields(),
      titlePattern,
      limit,
    }, true);
  }

  async createDocument(params: CreateDocumentParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/documents/create', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  async updateDocument(params: UpdateDocumentParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/documents/update', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  async deleteDocument(docId: number): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/documents/delete', 'POST', {
      ...generateBaseFields(),
      docId,
    }, true);
  }

  /**
   * Chunk Endpoints
   */

  async getChunksByDocument(docId: number): Promise<Result<Chunk[]>> {
    return this.request<Chunk[]>('/api/v1/chunks/get-by-document', 'POST', {
      ...generateBaseFields(),
      docId,
    }, true);
  }

  async getChunkById(chunkId: number): Promise<Result<Chunk>> {
    return this.request<Chunk>('/api/v1/chunks/get-by-id', 'POST', {
      ...generateBaseFields(),
      chunkId,
    }, true);
  }

  async similaritySearch(params: SimilaritySearchParams): Promise<Result<ChunkWithSimilarity[]>> {
    return this.request<ChunkWithSimilarity[]>('/api/v1/chunks/similarity-search', 'POST', {
      ...generateBaseFields(),
      queryText: params.queryText,
      limit: params.limit || 10,
      minSimilarity: params.minSimilarity || 0.7,
    }, true);
  }

  async hybridSearch(params: HybridSearchParams): Promise<Result<ChunkWithHybridSimilarity[]>> {
    return this.request<ChunkWithHybridSimilarity[]>('/api/v1/chunks/hybrid-search', 'POST', {
      ...generateBaseFields(),
      queryText: params.queryText,
      limit: params.limit || 10,
      minSimilarity: params.minSimilarity || 0.2,
      keywordWeight: params.keywordWeight || 0.15,
    }, true);
  }

  async createChunk(documentId: number, content: string): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunks/create', 'POST', {
      ...generateBaseFields(),
      documentId,
      content,
    }, true);
  }

  async updateChunkContent(chunkId: number, content: string): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunks/update-content', 'POST', {
      ...generateBaseFields(),
      chunkId,
      content,
    }, true);
  }

  async deleteChunk(chunkId: number): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunks/delete', 'POST', {
      ...generateBaseFields(),
      chunkId,
    }, true);
  }

  async bulkCreateChunks(params: BulkCreateChunksParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunks/bulk-create', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  /**
   * Chunk Statistics Endpoints
   */

  async getChunkStatistics(chunkId: number): Promise<Result<ChunkStatistics>> {
    return this.request<ChunkStatistics>('/api/v1/chunk-statistics/get-by-chunk', 'POST', {
      ...generateBaseFields(),
      chunkId,
    }, true);
  }

  async getTopChunksByUsage(limit: number = 10): Promise<Result<TopChunkByUsage[]>> {
    return this.request<TopChunkByUsage[]>('/api/v1/chunk-statistics/get-top-by-usage', 'POST', {
      ...generateBaseFields(),
      limit,
    }, true);
  }

  async incrementChunkUsage(chunkId: number): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunk-statistics/increment-usage', 'POST', {
      ...generateBaseFields(),
      chunkId,
    }, true);
  }

  async updateChunkQualityMetrics(params: UpdateChunkQualityMetricsParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunk-statistics/update-quality-metrics', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  async updateChunkStaleness(params: UpdateChunkStalenessParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/chunk-statistics/update-staleness', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  /**
   * Parameter Endpoints
   */

  async getAllParameters(): Promise<Result<Parameter[]>> {
    return this.request<Parameter[]>('/api/v1/parameters/get-all', 'POST', {
      ...generateBaseFields(),
    }, true);
  }

  async getParameterByCode(code: string): Promise<Result<Parameter>> {
    return this.request<Parameter>('/api/v1/parameters/get-by-code', 'POST', {
      ...generateBaseFields(),
      code,
    }, true);
  }

  async addParameter(params: AddParameterParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/parameters/add', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  async updateParameter(params: UpdateParameterParams): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/parameters/update', 'POST', {
      ...generateBaseFields(),
      ...params,
    }, true);
  }

  async deleteParameter(code: string): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/parameters/delete', 'POST', {
      ...generateBaseFields(),
      code,
    }, true);
  }

  async reloadParameterCache(): Promise<Result<Data>> {
    return this.request<Data>('/api/v1/parameters/reload-cache', 'POST', {
      ...generateBaseFields(),
    }, true);
  }
}

/**
 * Default API client instance
 * You can create multiple instances with different configurations
 */
export const apiClient = new ApiClient({
  baseUrl: process.env.REACT_APP_API_URL || 'http://localhost:8080',
});
