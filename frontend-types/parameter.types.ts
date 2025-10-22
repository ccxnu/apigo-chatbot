/**
 * Parameter types
 * Generated from Go backend domain models
 */

import { Data } from './base.types';

/**
 * Parameter represents a configuration parameter
 */
export interface Parameter {
  id: number;
  name: string;
  code: string;               // Unique code identifier
  data: any;                  // JSON data (can be object or primitive)
  description: string;
  active: boolean;
  createdAt: string;          // ISO 8601 timestamp
  updatedAt: string;          // ISO 8601 timestamp
}

/**
 * AddParameterParams for creating a new parameter
 */
export interface AddParameterParams {
  name: string;               // Min 3, max 100 characters
  code: string;               // Min 2, max 100 characters (unique)
  data: any;                  // JSON data
  description?: string;       // Max 500 characters
}

/**
 * UpdateParameterParams for updating an existing parameter
 */
export interface UpdateParameterParams {
  code: string;               // Parameter code to update
  name: string;               // Min 3, max 100 characters
  data: any;                  // JSON data
  description?: string;       // Max 500 characters
}
