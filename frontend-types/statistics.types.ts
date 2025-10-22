/**
 * Chunk statistics types
 * Generated from Go backend domain models
 */

/**
 * ChunkStatistics represents analytics and quality metrics for a chunk
 */
export interface ChunkStatistics {
  id: number;
  chunkId: number;
  usageCount: number;
  lastUsedAt?: string;              // ISO 8601 timestamp
  precisionAtK?: number;            // Precision@K metric
  recallAtK?: number;               // Recall@K metric
  f1AtK?: number;                   // F1@K metric
  mrr?: number;                     // Mean Reciprocal Rank
  map?: number;                     // Mean Average Precision
  ndcg?: number;                    // Normalized Discounted Cumulative Gain
  stalenessDays?: number;           // Days since last refresh
  lastRefreshAt?: string;           // ISO 8601 timestamp
  curriculumCoveragePct?: number;   // Percentage of curriculum coverage
  createdAt: string;                // ISO 8601 timestamp
  updatedAt: string;                // ISO 8601 timestamp
}

/**
 * TopChunkByUsage for analytics queries
 */
export interface TopChunkByUsage {
  chunkId: number;
  content: string;
  docTitle: string;
  usageCount: number;
  lastUsedAt?: string;              // ISO 8601 timestamp
  f1Score?: number;
}

/**
 * UpdateChunkQualityMetricsParams for updating quality metrics
 */
export interface UpdateChunkQualityMetricsParams {
  chunkId: number;
  precisionAtK?: number;
  recallAtK?: number;
  f1AtK?: number;
  mrr?: number;
  map?: number;
  ndcg?: number;
}

/**
 * UpdateChunkStalenessParams for updating staleness
 */
export interface UpdateChunkStalenessParams {
  chunkId: number;
  stalenessDays: number;
}
