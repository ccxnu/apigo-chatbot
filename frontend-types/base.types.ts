/**
 * Base types for API requests and responses
 * Generated from Go backend domain models
 */

/**
 * Base contains common fields required in all API requests
 * This ensures traceability and audit trail for all operations
 */
export interface Base {
  idSession: string;      // Session identifier
  idRequest: string;      // Unique request ID (UUID v4)
  process: string;        // Process name initiating the request
  idDevice: string;       // Device identifier
  deviceAddress: string;  // Device IP address (validation: must be valid IP)
  dateProcess: string;    // Timestamp when request was initiated (ISO 8601)
}

/**
 * Generic Result wrapper for all API responses
 */
export interface Result<T> {
  success: boolean;  // Indicates if the operation was successful
  code: string;      // Error/success code (e.g., "OK", "ERR_INTERNAL_DB")
  info: string;      // Human-readable message
  data: T;          // Response payload
}

/**
 * Generic data object for dynamic responses
 */
export type Data = Record<string, any>;
