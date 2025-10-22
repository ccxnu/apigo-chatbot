/**
 * Chunk types
 * Generated from Go backend domain models
 */

/**
 * Chunk represents a text chunk with optional embedding
 */
export interface Chunk {
  id: number;
  documentId: number;
  content: string;
  embedding?: number[];       // Vector embedding (optional)
  createdAt: string;          // ISO 8601 timestamp
  updatedAt: string;          // ISO 8601 timestamp
}

/**
 * ChunkWithSimilarity extends Chunk for similarity search results
 */
export interface ChunkWithSimilarity {
  id: number;
  documentId: number;
  content: string;
  similarityScore: number;    // 0-1 cosine similarity
  docTitle: string;
  docCategory: string;
}

/**
 * ChunkWithHybridSimilarity extends ChunkWithSimilarity for hybrid search results
 */
export interface ChunkWithHybridSimilarity {
  id: number;
  documentId: number;
  content: string;
  similarityScore: number;    // 0-1 cosine similarity (vector search)
  keywordScore: number;       // 0-1 full-text search score
  combinedScore: number;      // Weighted combination of both scores
  docTitle: string;
  docCategory: string;
}

/**
 * SimilaritySearchParams for vector similarity search
 */
export interface SimilaritySearchParams {
  queryText: string;          // Text query (automatically converted to embedding)
  limit?: number;             // Max results (default: 10, max: 100)
  minSimilarity?: number;     // Min similarity 0-1 (default: 0.7)
}

/**
 * HybridSearchParams for hybrid search (vector + keyword)
 */
export interface HybridSearchParams {
  queryText: string;          // Text query (converted to embedding + full-text search)
  limit?: number;             // Max results (default: 10, max: 100)
  minSimilarity?: number;     // Min similarity 0-1 (default: 0.2)
  keywordWeight?: number;     // Weight for keyword score 0-1 (default: 0.15)
}

/**
 * BulkCreateChunksParams for creating multiple chunks at once
 */
export interface BulkCreateChunksParams {
  documentId: number;
  contents: string[];         // Array of chunk contents (embeddings generated automatically)
}
