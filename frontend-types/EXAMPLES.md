# API Client Usage Examples

Complete examples for using the API client with documents, chunks, parameters, and more.

## Table of Contents

- [Document Management](#document-management)
- [Chunk Management](#chunk-management)
- [Similarity & Hybrid Search](#similarity--hybrid-search)
- [Chunk Statistics & Analytics](#chunk-statistics--analytics)
- [Parameter Management](#parameter-management)
- [React Hooks Examples](#react-hooks-examples)
- [Complete Components](#complete-components)

## Document Management

### Fetch All Documents

```typescript
import { apiClient, Document } from '@api-chatbot/types';

// Get all documents with pagination
async function loadDocuments() {
  try {
    const result = await apiClient.getAllDocuments(100, 0);

    if (result.success) {
      const documents: Document[] = result.data;
      console.log('Documents:', documents);
      return documents;
    }
  } catch (error) {
    console.error('Failed to load documents:', error);
  }
}
```

### Search Documents by Category

```typescript
async function getDocumentsByCategory(category: string) {
  const result = await apiClient.getDocumentsByCategory(category, 50, 0);

  if (result.success) {
    return result.data;
  }

  throw new Error(result.info);
}
```

### Create a New Document

```typescript
async function createNewDocument() {
  const result = await apiClient.createDocument({
    category: 'programming',
    title: 'Introduction to TypeScript',
    summary: 'A comprehensive guide to TypeScript fundamentals',
    source: 'https://example.com/typescript-guide',
    publishedAt: new Date().toISOString(),
  });

  if (result.success) {
    console.log('Document created:', result.data);
  }
}
```

### Update Document

```typescript
async function updateDocument(docId: number) {
  const result = await apiClient.updateDocument({
    docId,
    category: 'programming',
    title: 'Advanced TypeScript Patterns',
    summary: 'Updated content about TypeScript',
    source: 'https://example.com/updated',
  });

  return result.success;
}
```

### Delete Document

```typescript
async function deleteDocument(docId: number) {
  const result = await apiClient.deleteDocument(docId);

  if (result.success) {
    console.log('Document deleted');
  }
}
```

## Chunk Management

### Get Chunks for a Document

```typescript
async function getDocumentChunks(docId: number) {
  const result = await apiClient.getChunksByDocument(docId);

  if (result.success) {
    const chunks = result.data;
    console.log(`Found ${chunks.length} chunks`);
    return chunks;
  }
}
```

### Create Single Chunk

```typescript
async function addChunk(documentId: number, content: string) {
  const result = await apiClient.createChunk(documentId, content);

  if (result.success) {
    console.log('Chunk created with auto-generated embedding');
  }
}
```

### Bulk Create Chunks

```typescript
async function bulkAddChunks(documentId: number) {
  const contents = [
    'First chunk of content about TypeScript types',
    'Second chunk explaining interfaces and type aliases',
    'Third chunk covering generics and utility types',
  ];

  const result = await apiClient.bulkCreateChunks({
    documentId,
    contents,
  });

  if (result.success) {
    console.log('All chunks created with embeddings');
  }
}
```

## Similarity & Hybrid Search

### Similarity Search (Vector Search Only)

```typescript
async function searchSimilarContent(query: string) {
  const result = await apiClient.similaritySearch({
    queryText: query,
    limit: 10,
    minSimilarity: 0.7,
  });

  if (result.success) {
    result.data.forEach((chunk) => {
      console.log(`Score: ${chunk.similarityScore.toFixed(3)}`);
      console.log(`Content: ${chunk.content}`);
      console.log(`From: ${chunk.docTitle} (${chunk.docCategory})`);
      console.log('---');
    });
  }
}
```

### Hybrid Search (Vector + Keyword)

```typescript
async function hybridSearch(query: string) {
  const result = await apiClient.hybridSearch({
    queryText: query,
    limit: 15,
    minSimilarity: 0.2,
    keywordWeight: 0.15, // 15% keyword, 85% vector similarity
  });

  if (result.success) {
    result.data.forEach((chunk) => {
      console.log(`Vector Score: ${chunk.similarityScore.toFixed(3)}`);
      console.log(`Keyword Score: ${chunk.keywordScore.toFixed(3)}`);
      console.log(`Combined Score: ${chunk.combinedScore.toFixed(3)}`);
      console.log(`Content: ${chunk.content.substring(0, 100)}...`);
      console.log('---');
    });
  }
}
```

## Chunk Statistics & Analytics

### Get Chunk Statistics

```typescript
async function getChunkStats(chunkId: number) {
  const result = await apiClient.getChunkStatistics(chunkId);

  if (result.success) {
    const stats = result.data;
    console.log('Usage Count:', stats.usageCount);
    console.log('F1 Score:', stats.f1AtK);
    console.log('Precision:', stats.precisionAtK);
    console.log('Recall:', stats.recallAtK);
    console.log('NDCG:', stats.ndcg);
  }
}
```

### Get Top Used Chunks

```typescript
async function getTopChunks() {
  const result = await apiClient.getTopChunksByUsage(20);

  if (result.success) {
    console.log('Most Used Chunks:');
    result.data.forEach((chunk, index) => {
      console.log(`${index + 1}. ${chunk.docTitle}`);
      console.log(`   Usage: ${chunk.usageCount} times`);
      console.log(`   F1 Score: ${chunk.f1Score || 'N/A'}`);
    });
  }
}
```

### Track Chunk Usage

```typescript
async function trackChunkUsage(chunkId: number) {
  // Call this when a chunk is shown to user
  await apiClient.incrementChunkUsage(chunkId);
}
```

### Update Quality Metrics

```typescript
async function updateMetrics(chunkId: number) {
  const result = await apiClient.updateChunkQualityMetrics({
    chunkId,
    precisionAtK: 0.85,
    recallAtK: 0.78,
    f1AtK: 0.81,
    mrr: 0.90,
    map: 0.88,
    ndcg: 0.92,
  });

  if (result.success) {
    console.log('Quality metrics updated');
  }
}
```

## Parameter Management

### Get All Parameters

```typescript
async function loadAllParameters() {
  const result = await apiClient.getAllParameters();

  if (result.success) {
    result.data.forEach(param => {
      console.log(`${param.code}: ${param.name}`);
      console.log('Data:', param.data);
    });
  }
}
```

### Get Specific Parameter

```typescript
async function getParameter(code: string) {
  const result = await apiClient.getParameterByCode(code);

  if (result.success) {
    const param = result.data;
    return param.data; // Returns the JSON data
  }
}
```

### Add/Update Parameters

```typescript
async function manageParameters() {
  // Add new parameter
  await apiClient.addParameter({
    name: 'Max Search Results',
    code: 'MAX_SEARCH_RESULTS',
    data: { value: 50, unit: 'items' },
    description: 'Maximum number of search results to return',
  });

  // Update existing parameter
  await apiClient.updateParameter({
    code: 'MAX_SEARCH_RESULTS',
    name: 'Max Search Results',
    data: { value: 100, unit: 'items' },
    description: 'Updated maximum results',
  });
}
```

## React Hooks Examples

### useDocuments Hook

```typescript
// hooks/useDocuments.ts
import { useState, useEffect } from 'react';
import { apiClient, Document } from '@api-chatbot/types';

export function useDocuments(category?: string) {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadDocs() {
      setLoading(true);
      setError(null);

      try {
        const result = category
          ? await apiClient.getDocumentsByCategory(category)
          : await apiClient.getAllDocuments();

        if (result.success) {
          setDocuments(result.data);
        } else {
          setError(result.info);
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    loadDocs();
  }, [category]);

  return { documents, loading, error };
}
```

### useSearch Hook

```typescript
// hooks/useSearch.ts
import { useState } from 'react';
import { apiClient, ChunkWithHybridSimilarity } from '@api-chatbot/types';

export function useSearch() {
  const [results, setResults] = useState<ChunkWithHybridSimilarity[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = async (query: string, useHybrid: boolean = true) => {
    if (!query.trim()) return;

    setLoading(true);
    setError(null);

    try {
      const result = useHybrid
        ? await apiClient.hybridSearch({ queryText: query, limit: 20 })
        : await apiClient.similaritySearch({ queryText: query, limit: 20 });

      if (result.success) {
        setResults(result.data);
      } else {
        setError(result.info);
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const clear = () => {
    setResults([]);
    setError(null);
  };

  return { results, loading, error, search, clear };
}
```

## Complete Components

### Document List Component

```typescript
// components/DocumentList.tsx
import React from 'react';
import { useDocuments } from '../hooks/useDocuments';

interface Props {
  category?: string;
}

export function DocumentList({ category }: Props) {
  const { documents, loading, error } = useDocuments(category);

  if (loading) return <div>Loading documents...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="document-list">
      <h2>{category ? `${category} Documents` : 'All Documents'}</h2>

      {documents.length === 0 ? (
        <p>No documents found</p>
      ) : (
        <ul>
          {documents.map((doc) => (
            <li key={doc.id}>
              <h3>{doc.title}</h3>
              <p className="category">{doc.category}</p>
              {doc.summary && <p>{doc.summary}</p>}
              {doc.source && (
                <a href={doc.source} target="_blank" rel="noopener noreferrer">
                  Source
                </a>
              )}
              <small>
                Created: {new Date(doc.createdAt).toLocaleDateString()}
              </small>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
```

### Search Component

```typescript
// components/SearchBox.tsx
import React, { useState } from 'react';
import { useSearch } from '../hooks/useSearch';

export function SearchBox() {
  const [query, setQuery] = useState('');
  const [searchType, setSearchType] = useState<'hybrid' | 'vector'>('hybrid');
  const { results, loading, error, search, clear } = useSearch();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    search(query, searchType === 'hybrid');
  };

  return (
    <div className="search-box">
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search knowledge base..."
          disabled={loading}
        />

        <select
          value={searchType}
          onChange={(e) => setSearchType(e.target.value as any)}
        >
          <option value="hybrid">Hybrid Search</option>
          <option value="vector">Vector Only</option>
        </select>

        <button type="submit" disabled={loading || !query.trim()}>
          {loading ? 'Searching...' : 'Search'}
        </button>

        {results.length > 0 && (
          <button type="button" onClick={clear}>Clear</button>
        )}
      </form>

      {error && <div className="error">{error}</div>}

      {results.length > 0 && (
        <div className="results">
          <h3>Found {results.length} results</h3>
          {results.map((chunk) => (
            <div key={chunk.id} className="result-item">
              <div className="scores">
                <span>Combined: {chunk.combinedScore.toFixed(3)}</span>
                <span>Vector: {chunk.similarityScore.toFixed(3)}</span>
                <span>Keyword: {chunk.keywordScore.toFixed(3)}</span>
              </div>
              <p className="content">{chunk.content}</p>
              <div className="meta">
                <span className="title">{chunk.docTitle}</span>
                <span className="category">{chunk.docCategory}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
```

### Document Upload Component

```typescript
// components/DocumentUpload.tsx
import React, { useState } from 'react';
import { apiClient } from '@api-chatbot/types';

export function DocumentUpload() {
  const [formData, setFormData] = useState({
    category: '',
    title: '',
    summary: '',
    source: '',
  });
  const [chunks, setChunks] = useState<string[]>(['']);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');

    try {
      // 1. Create document
      const docResult = await apiClient.createDocument({
        category: formData.category,
        title: formData.title,
        summary: formData.summary || undefined,
        source: formData.source || undefined,
      });

      if (!docResult.success) {
        throw new Error(docResult.info);
      }

      const docId = docResult.data.docId;

      // 2. Create chunks
      const validChunks = chunks.filter(c => c.trim());
      if (validChunks.length > 0) {
        const chunksResult = await apiClient.bulkCreateChunks({
          documentId: docId,
          contents: validChunks,
        });

        if (!chunksResult.success) {
          throw new Error(chunksResult.info);
        }
      }

      setMessage('Document and chunks created successfully!');

      // Reset form
      setFormData({ category: '', title: '', summary: '', source: '' });
      setChunks(['']);
    } catch (error: any) {
      setMessage(`Error: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const addChunk = () => {
    setChunks([...chunks, '']);
  };

  const updateChunk = (index: number, value: string) => {
    const newChunks = [...chunks];
    newChunks[index] = value;
    setChunks(newChunks);
  };

  return (
    <form onSubmit={handleSubmit} className="document-upload">
      <h2>Upload Document</h2>

      <input
        type="text"
        placeholder="Category"
        value={formData.category}
        onChange={(e) => setFormData({ ...formData, category: e.target.value })}
        required
      />

      <input
        type="text"
        placeholder="Title"
        value={formData.title}
        onChange={(e) => setFormData({ ...formData, title: e.target.value })}
        required
      />

      <textarea
        placeholder="Summary (optional)"
        value={formData.summary}
        onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
      />

      <input
        type="url"
        placeholder="Source URL (optional)"
        value={formData.source}
        onChange={(e) => setFormData({ ...formData, source: e.target.value })}
      />

      <h3>Chunks</h3>
      {chunks.map((chunk, index) => (
        <textarea
          key={index}
          placeholder={`Chunk ${index + 1}`}
          value={chunk}
          onChange={(e) => updateChunk(index, e.target.value)}
        />
      ))}

      <button type="button" onClick={addChunk}>
        Add Another Chunk
      </button>

      <button type="submit" disabled={loading}>
        {loading ? 'Uploading...' : 'Upload Document'}
      </button>

      {message && (
        <div className={message.startsWith('Error') ? 'error' : 'success'}>
          {message}
        </div>
      )}
    </form>
  );
}
```

### Analytics Dashboard

```typescript
// components/AnalyticsDashboard.tsx
import React, { useEffect, useState } from 'react';
import { apiClient, TopChunkByUsage } from '@api-chatbot/types';

export function AnalyticsDashboard() {
  const [topChunks, setTopChunks] = useState<TopChunkByUsage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadAnalytics() {
      const result = await apiClient.getTopChunksByUsage(10);

      if (result.success) {
        setTopChunks(result.data);
      }

      setLoading(false);
    }

    loadAnalytics();
  }, []);

  if (loading) return <div>Loading analytics...</div>;

  return (
    <div className="analytics-dashboard">
      <h2>Top Performing Chunks</h2>

      <table>
        <thead>
          <tr>
            <th>Rank</th>
            <th>Document</th>
            <th>Usage Count</th>
            <th>F1 Score</th>
            <th>Last Used</th>
          </tr>
        </thead>
        <tbody>
          {topChunks.map((chunk, index) => (
            <tr key={chunk.chunkId}>
              <td>{index + 1}</td>
              <td>{chunk.docTitle}</td>
              <td>{chunk.usageCount}</td>
              <td>{chunk.f1Score?.toFixed(3) || 'N/A'}</td>
              <td>
                {chunk.lastUsedAt
                  ? new Date(chunk.lastUsedAt).toLocaleDateString()
                  : 'Never'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

## Error Handling Best Practices

```typescript
import { ApiError } from '@api-chatbot/types';

async function handleApiCall() {
  try {
    const result = await apiClient.getAllDocuments();

    if (result.success) {
      // Handle success
      return result.data;
    } else {
      // Handle application error
      console.error('App error:', result.code, result.info);
    }
  } catch (error) {
    if (error instanceof ApiError) {
      // Handle API error
      switch (error.code) {
        case 'TIMEOUT':
          console.error('Request timed out');
          break;
        case 'NETWORK_ERROR':
          console.error('Network error');
          break;
        case 'ERR_UNAUTHORIZED':
          console.error('Unauthorized - redirect to login');
          break;
        default:
          console.error('API Error:', error.info);
      }
    } else {
      // Handle unknown error
      console.error('Unknown error:', error);
    }
  }
}
```

## TypeScript Type Guards

```typescript
import { ChunkWithSimilarity, ChunkWithHybridSimilarity } from '@api-chatbot/types';

function isHybridResult(
  chunk: ChunkWithSimilarity | ChunkWithHybridSimilarity
): chunk is ChunkWithHybridSimilarity {
  return 'keywordScore' in chunk && 'combinedScore' in chunk;
}

// Usage
function displayScore(chunk: ChunkWithSimilarity | ChunkWithHybridSimilarity) {
  if (isHybridResult(chunk)) {
    console.log(`Hybrid score: ${chunk.combinedScore}`);
  } else {
    console.log(`Similarity: ${chunk.similarityScore}`);
  }
}
```
