// API Response Types
export interface Event {
  id: string
  name: string
  code: string
  description?: string
  event_date?: string
  status: "active" | "inactive" | "closed"
  owner_email: string
  created_at: string
  updated_at: string
}

export interface Session {
  id: string
  event_id: string
  guest_name: string
  session_token: string
  expires_at: string
  created_at: string
  event?: Event
}

export interface Photo {
  id: string
  event_id: string
  uploader_name: string
  object_key: string
  file_size: number
  mime_type: string
  created_at: string
  updated_at: string
}

export interface UploadURLResponse {
  upload_url: string
  object_key: string
  photo_id: string
}

// API Request Types
export interface CreateSessionRequest {
  event_code: string
  guest_name: string
}

export interface RefreshSessionRequest {
  session_token: string
}

export interface RevokeSessionRequest {
  session_token: string
}

export interface UploadURLRequest {
  event_id: string
  content_type: string
}

export interface ConfirmUploadRequest {
  file_size: number
}

// Error Response
export interface APIError {
  message: string
  error?: string
}