# PDF Upload Endpoint Documentation

## Overview
The PDF upload endpoint allows you to upload a PDF file (encoded as base64), extract text from it, create a document in the knowledge base, and automatically generate chunks with embeddings for semantic search.

## Endpoint Details

**URL:** `/api/v1/documents/upload-pdf`
**Method:** `POST`
**Content-Type:** `application/json`

## Request Body

```json
{
  "idSession": "string (required)",
  "idRequest": "string (required, UUID format)",
  "process": "string (required)",
  "idDevice": "string (required)",
  "publicIp": "string (required)",
  "dateProcess": "string (required, ISO 8601 timestamp)",
  "category": "string (required)",
  "title": "string (required, 1-200 characters)",
  "source": "string (optional, max 500 characters)",
  "fileBase64": "string (required, base64-encoded PDF)",
  "chunkSize": "integer (optional, default: 1000, range: 100-5000)",
  "chunkOverlap": "integer (optional, default: 200, range: 0-500)"
}
```

### Field Descriptions

- **idSession**: Session identifier for tracking
- **idRequest**: Unique request ID (UUID format)
- **process**: Process name or identifier
- **idDevice**: Device identifier
- **publicIp**: Client's public IP address
- **dateProcess**: Timestamp of the request
- **category**: Document category (e.g., "Technology", "Science", "Business")
- **title**: Document title
- **source**: Optional source or reference URL
- **fileBase64**: Base64-encoded PDF file content
- **chunkSize**: Size of each text chunk in characters (default: 1000)
- **chunkOverlap**: Number of overlapping characters between chunks (default: 200)

## Response

### Success Response (200 OK)

```json
{
  "success": true,
  "code": "SUCCESS",
  "message": "OK",
  "data": {
    "docId": 123,
    "chunksCreated": 15,
    "message": "PDF uploaded and processed successfully"
  }
}
```

### Error Responses

#### PDF Processing Error
```json
{
  "success": false,
  "code": "ERR_PDF_PROCESSING",
  "message": "Error al procesar el archivo PDF",
  "data": null
}
```

#### Chunk Creation Error
```json
{
  "success": false,
  "code": "ERR_CHUNK_CREATION",
  "message": "Error al crear los fragmentos del documento",
  "data": null
}
```

#### Validation Error
```json
{
  "success": false,
  "code": "ERR_VALIDATION",
  "message": "Validation error message",
  "data": null
}
```

## Processing Flow

1. **PDF Text Extraction**: The base64-encoded PDF is decoded and text is extracted from all pages
2. **Document Creation**: A new document record is created in the database with the extracted text summary
3. **Text Chunking**: The extracted text is split into chunks based on sentence boundaries with optional overlap
4. **Embedding Generation**: Embeddings are generated for each chunk using the configured embedding service
5. **Chunk Storage**: All chunks with their embeddings are stored in the database

## Features

- **Automatic Text Extraction**: Extracts text from text-based PDFs
- **Smart Chunking**: Splits text at sentence boundaries for better semantic coherence
- **Overlap Support**: Maintains context between chunks with configurable overlap
- **Bulk Operations**: Efficiently processes multiple chunks in a single operation
- **Embedding Generation**: Automatically generates embeddings for semantic search

## Example Usage

### Using cURL

```bash
# First, encode your PDF file to base64
BASE64_PDF=$(base64 -w 0 your_document.pdf)

# Make the API request
curl -X POST http://localhost:8080/api/v1/documents/upload-pdf \
  -H 'Content-Type: application/json' \
  -d '{
    "idSession": "user-session-123",
    "idRequest": "550e8400-e29b-41d4-a716-446655440000",
    "process": "document-upload",
    "idDevice": "device-001",
    "publicIp": "192.168.1.1",
    "dateProcess": "2025-10-30T10:00:00Z",
    "category": "Technology",
    "title": "Introduction to Machine Learning",
    "source": "https://example.com/ml-guide",
    "fileBase64": "'"$BASE64_PDF"'",
    "chunkSize": 800,
    "chunkOverlap": 150
  }'
```

### Using Python

```python
import requests
import base64
from datetime import datetime

# Read and encode PDF file
with open('your_document.pdf', 'rb') as f:
    pdf_content = f.read()
    base64_pdf = base64.b64encode(pdf_content).decode('utf-8')

# Prepare request
url = 'http://localhost:8080/api/v1/documents/upload-pdf'
headers = {'Content-Type': 'application/json'}
payload = {
    'idSession': 'user-session-123',
    'idRequest': '550e8400-e29b-41d4-a716-446655440000',
    'process': 'document-upload',
    'idDevice': 'device-001',
    'publicIp': '192.168.1.1',
    'dateProcess': datetime.utcnow().isoformat() + 'Z',
    'category': 'Technology',
    'title': 'Introduction to Machine Learning',
    'source': 'https://example.com/ml-guide',
    'fileBase64': base64_pdf,
    'chunkSize': 800,
    'chunkOverlap': 150
}

# Send request
response = requests.post(url, json=payload, headers=headers)
print(response.json())
```

## Notes

- **PDF Type Support**: Currently supports text-based PDFs. For scanned/image-based PDFs, OCR functionality can be added by installing Tesseract OCR libraries
- **File Size**: Large PDFs may take longer to process due to text extraction and embedding generation
- **Timeout**: The endpoint has a 2-minute timeout for processing
- **Chunking Strategy**: Uses sentence-based chunking to maintain semantic coherence
- **Embedding Model**: Uses the configured embedding service (OpenAI by default)

## OCR Support (Future Enhancement)

To add OCR support for scanned PDFs:

1. Install Tesseract OCR and dependencies:
   ```bash
   apt-get install tesseract-ocr libtesseract-dev libleptonica-dev
   ```

2. Add the gosseract dependency:
   ```bash
   go get github.com/otiai10/gosseract/v2
   ```

3. Update the PDF processor to use OCR when text extraction yields minimal results

## Related Endpoints

- `POST /api/v1/documents/create` - Create document manually
- `POST /api/v1/chunks/bulk-create` - Bulk create chunks
- `POST /api/v1/chunks/similarity-search` - Search chunks by similarity
- `POST /api/v1/documents/get-all` - List all documents

## Error Codes Reference

| Code | Description |
|------|-------------|
| `SUCCESS` | Operation completed successfully |
| `ERR_PDF_PROCESSING` | Failed to extract text from PDF |
| `ERR_CHUNK_CREATION` | Failed to create chunks |
| `ERR_VALIDATION` | Request validation failed |
| `ERR_INTERNAL_DB` | Database operation failed |
| `ERR_EMBEDDING_GENERATION` | Failed to generate embeddings |
