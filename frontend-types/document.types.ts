/**
 * Document types
 * Generated from Go backend domain models
 */

/**
 * Document represents a knowledge base document
 */
export interface Document {
  id: number;
  category: string;
  title: string;
  summary?: string;
  source?: string;
  publishedAt?: string;      // ISO 8601 timestamp
  active: boolean;
  createdAt: string;          // ISO 8601 timestamp
  updatedAt: string;          // ISO 8601 timestamp
}

/**
 * CreateDocumentParams for creating a new document
 */
export interface CreateDocumentParams {
  category: string;
  title: string;              // Min 1, max 200 characters
  summary?: string;           // Max 5000 characters
  source?: string;            // Max 500 characters
  publishedAt?: string;       // ISO 8601 timestamp
}

/**
 * UpdateDocumentParams for updating an existing document
 */
export interface UpdateDocumentParams {
  docId: number;              // Document ID (must be >= 1)
  category: string;
  title: string;              // Min 1, max 200 characters
  summary?: string;           // Max 5000 characters
  source?: string;            // Max 500 characters
  publishedAt?: string;       // ISO 8601 timestamp
}
