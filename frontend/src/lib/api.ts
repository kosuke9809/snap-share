import type {
  APIError,
  ConfirmUploadRequest,
  CreateSessionRequest,
  Event,
  RefreshSessionRequest,
  RevokeSessionRequest,
  Session,
  UploadURLRequest,
  UploadURLResponse,
} from "@/types/api"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

class APIClient {
  private baseURL: string
  private authToken: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  setAuthToken(token: string | null) {
    this.authToken = token
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    }

    if (this.authToken) {
      headers.Authorization = `Bearer ${this.authToken}`
      console.log("Making request with token:", this.authToken.substring(0, 20) + "...")
    } else {
      console.log("Making request without auth token")
    }

    console.log(`API Request: ${options.method || 'GET'} ${url}`)

    const response = await fetch(url, {
      ...options,
      headers,
    })

    console.log(`API Response: ${response.status} ${response.statusText}`)

    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}`
      try {
        const error: APIError = await response.json()
        errorMessage = error.message || errorMessage
      } catch {
        errorMessage = response.statusText || errorMessage
      }
      throw new Error(errorMessage)
    }

    return response.json()
  }

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    return this.request("/health")
  }

  // Event endpoints
  async getEventByCode(code: string): Promise<Event> {
    return this.request(`/api/events/${code}`)
  }

  // Session endpoints
  async createSession(data: CreateSessionRequest): Promise<Session> {
    return this.request("/api/sessions", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async refreshSession(data: RefreshSessionRequest): Promise<Session> {
    return this.request("/api/sessions/refresh", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async revokeSession(data: RevokeSessionRequest): Promise<{ message: string }> {
    return this.request("/api/sessions", {
      method: "DELETE",
      body: JSON.stringify(data),
    })
  }

  // Photo endpoints
  async getUploadURL(data: UploadURLRequest): Promise<UploadURLResponse> {
    return this.request("/api/photos/upload-url", {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  async confirmUpload(
    photoId: string,
    data: ConfirmUploadRequest,
  ): Promise<{ message: string }> {
    return this.request(`/api/photos/confirm/${photoId}`, {
      method: "POST",
      body: JSON.stringify(data),
    })
  }

  // File upload to presigned URL
  async uploadFile(uploadURL: string, file: File): Promise<void> {
    console.log("Uploading file:", file.name, "Size:", file.size, "Type:", file.type)
    console.log("Upload URL:", uploadURL)

    try {
      const response = await fetch(uploadURL, {
        method: "PUT",
        body: file,
        headers: {
          "Content-Type": file.type,
        },
      })

      console.log("Upload response status:", response.status)
      console.log("Upload response headers:", Array.from(response.headers.entries()))

      if (!response.ok) {
        let errorText = ""
        try {
          errorText = await response.text()
        } catch (e) {
          errorText = response.statusText
        }
        console.error("Upload failed with response:", errorText)
        throw new Error(`Upload failed (${response.status}): ${errorText || response.statusText}`)
      }

      console.log("File uploaded successfully")
    } catch (error) {
      console.error("Upload error details:", error)
      if (error instanceof TypeError && error.message.includes("network")) {
        throw new Error("ネットワークエラー: アップロード先に接続できませんでした")
      }
      throw error
    }
  }
}

export const apiClient = new APIClient(API_BASE_URL)
export default apiClient