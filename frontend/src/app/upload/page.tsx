"use client"

import { useCallback, useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Camera, Upload, X, Check, AlertCircle, LogOut } from "lucide-react"
import { useAuthStore } from "@/stores/auth"
import { apiClient } from "@/lib/api"

interface UploadFile {
  id: string
  file: File
  preview: string
  status: "pending" | "uploading" | "success" | "error"
  progress: number
  error?: string
}

export default function UploadPage() {
  const router = useRouter()
  const { session, isAuthenticated, clearSession } = useAuthStore()
  const [files, setFiles] = useState<UploadFile[]>([])
  const [dragActive, setDragActive] = useState(false)

  // Redirect if not authenticated
  useEffect(() => {
    if (!isAuthenticated || !session) {
      router.push("/")
    }
  }, [isAuthenticated, session, router])

  // Loading state while checking authentication
  if (!isAuthenticated || !session) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-ocean">
        <div className="text-center text-white">
          <div className="w-8 h-8 border-4 border-white border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p>認証確認中...</p>
        </div>
      </div>
    )
  }

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true)
    } else if (e.type === "dragleave") {
      setDragActive(false)
    }
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)

    const droppedFiles = Array.from(e.dataTransfer.files)
    handleFiles(droppedFiles)
  }, [])

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const selectedFiles = Array.from(e.target.files)
      handleFiles(selectedFiles)
    }
  }

  const handleFiles = (newFiles: File[]) => {
    const imageFiles = newFiles.filter(file => file.type.startsWith("image/"))

    const uploadFiles: UploadFile[] = imageFiles.map(file => ({
      id: Math.random().toString(36).substr(2, 9),
      file,
      preview: URL.createObjectURL(file),
      status: "pending",
      progress: 0
    }))

    setFiles(prev => [...prev, ...uploadFiles])
  }

  const uploadFile = async (uploadFile: UploadFile) => {
    console.log("Starting upload for:", uploadFile.file.name)
    console.log("Session:", session)

    if (!session?.event_id) {
      console.error("No event_id in session")
      setFiles(prev =>
        prev.map(f =>
          f.id === uploadFile.id
            ? { ...f, status: "error", error: "セッション情報がありません" }
            : f
        )
      )
      return
    }

    setFiles(prev =>
      prev.map(f =>
        f.id === uploadFile.id
          ? { ...f, status: "uploading", progress: 0 }
          : f
      )
    )

    try {
      console.log("Getting upload URL...")
      // Get upload URL
      const uploadResponse = await apiClient.getUploadURL({
        event_id: session.event_id,
        content_type: uploadFile.file.type
      })
      console.log("Upload response:", uploadResponse)

      // Update progress
      setFiles(prev =>
        prev.map(f =>
          f.id === uploadFile.id ? { ...f, progress: 30 } : f
        )
      )

      console.log("Uploading file to:", uploadResponse.upload_url)
      // Upload file to MinIO
      await apiClient.uploadFile(uploadResponse.upload_url, uploadFile.file)

      // Update progress
      setFiles(prev =>
        prev.map(f =>
          f.id === uploadFile.id ? { ...f, progress: 70 } : f
        )
      )

      console.log("Confirming upload...")
      // Confirm upload
      await apiClient.confirmUpload(uploadResponse.photo_id, {
        file_size: uploadFile.file.size
      })

      console.log("Upload completed successfully")
      // Success
      setFiles(prev =>
        prev.map(f =>
          f.id === uploadFile.id
            ? { ...f, status: "success", progress: 100 }
            : f
        )
      )
    } catch (error) {
      console.error("Upload error:", error)
      const errorMessage = error instanceof Error ? error.message : "アップロードに失敗しました"
      setFiles(prev =>
        prev.map(f =>
          f.id === uploadFile.id
            ? {
                ...f,
                status: "error",
                progress: 0,
                error: errorMessage
              }
            : f
        )
      )
    }
  }

  const removeFile = (id: string) => {
    setFiles(prev => {
      const file = prev.find(f => f.id === id)
      if (file) {
        URL.revokeObjectURL(file.preview)
      }
      return prev.filter(f => f.id !== id)
    })
  }

  const uploadAll = () => {
    const pendingFiles = files.filter(f => f.status === "pending")
    pendingFiles.forEach(uploadFile)
  }

  const retryUpload = (uploadFileData: UploadFile) => {
    uploadFile(uploadFileData)
  }

  const handleLogout = () => {
    clearSession()
    router.push("/")
  }

  const pendingCount = files.filter(f => f.status === "pending").length
  const uploadingCount = files.filter(f => f.status === "uploading").length
  const successCount = files.filter(f => f.status === "success").length

  return (
    <div className="min-h-screen bg-gradient-ocean relative overflow-hidden">
      {/* Background Elements */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-20 left-20 w-64 h-64 bg-white rounded-full mix-blend-multiply filter blur-xl animate-float"></div>
        <div className="absolute bottom-20 right-20 w-64 h-64 bg-blue-300 rounded-full mix-blend-multiply filter blur-xl animate-float" style={{ animationDelay: '2s' }}></div>
      </div>

      {/* Header */}
      <div className="relative z-10 glass-dark sticky top-0 backdrop-blur-md">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="min-w-0 flex-1">
            <h1 className="text-xl sm:text-2xl font-bold text-white truncate text-elegant">
              {session.event?.name}
            </h1>
            <p className="text-white/80 text-sm">
              ようこそ、{session.guest_name}さん
            </p>
          </div>
          <button
            onClick={handleLogout}
            className="glass px-4 py-2 text-white hover:bg-white/20 focus:ring-2 focus:ring-white/50 rounded-xl transition-all duration-300 flex items-center gap-2"
          >
            <LogOut className="w-4 h-4" />
            <span className="text-sm font-medium">ログアウト</span>
          </button>
        </div>
      </div>

      <div className="relative z-10 max-w-6xl mx-auto px-4 py-6">
        {/* Upload Area */}
        <div
          className={`
            card-elevated p-8 sm:p-12 text-center transition-all touch-none mb-8
            ${dragActive
              ? "border-4 border-blue-500 bg-blue-50 scale-105"
              : "hover:scale-102 hover:shadow-2xl"
            }
          `}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
          onTouchStart={(e) => e.preventDefault()}
        >
          <input
            type="file"
            multiple
            accept="image/*"
            capture="environment"
            onChange={handleFileInput}
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
          />

          <div className="space-y-6 animate-fade-in-up">
            <div className="w-20 h-20 sm:w-24 sm:h-24 bg-gradient-primary rounded-3xl flex items-center justify-center mx-auto shadow-lg animate-float">
              <Camera className="w-10 h-10 sm:w-12 sm:h-12 text-white" />
            </div>
            <div>
              <h3 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-4 text-elegant">
                写真をアップロード
              </h3>
              <p className="text-lg text-gray-600 px-4 leading-relaxed">
                <span className="hidden sm:inline">ファイルをドラッグ&ドロップするか、</span>
                <span className="sm:hidden">写真を選択するには</span>
                ここをタップしてください
              </p>
              <p className="text-sm text-gray-500 mt-4 font-medium">
                対応形式: JPG, PNG, GIF (最大10MB)
              </p>
            </div>
          </div>
        </div>

        {/* File List */}
        {files.length > 0 && (
          <div className="animate-fade-in-up">
            <div className="flex items-center justify-between mb-8">
              <h3 className="text-2xl font-bold text-white text-elegant">
                選択された写真 ({files.length})
              </h3>
              {pendingCount > 0 && (
                <button
                  onClick={uploadAll}
                  disabled={uploadingCount > 0}
                  className="btn-secondary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <Upload className="w-4 h-4" />
                  すべてアップロード ({pendingCount})
                </button>
              )}
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6">
              {files.map((uploadFile) => (
                <div
                  key={uploadFile.id}
                  className="card-elevated overflow-hidden group hover:scale-105 transition-all duration-300"
                >
                  <div className="aspect-square relative">
                    <img
                      src={uploadFile.preview}
                      alt="Preview"
                      className="w-full h-full object-cover"
                    />

                    {/* Status Overlay */}
                    <div className="absolute top-3 right-3">
                      {uploadFile.status === "pending" && (
                        <button
                          onClick={() => removeFile(uploadFile.id)}
                          className="w-8 h-8 bg-black/60 backdrop-blur-sm rounded-full flex items-center justify-center text-white hover:bg-black/80 focus:ring-2 focus:ring-white/50 transition-all"
                        >
                          <X className="w-4 h-4" />
                        </button>
                      )}
                      {uploadFile.status === "success" && (
                        <div className="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center text-white shadow-lg">
                          <Check className="w-4 h-4" />
                        </div>
                      )}
                      {uploadFile.status === "error" && (
                        <div className="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center text-white shadow-lg">
                          <AlertCircle className="w-4 h-4" />
                        </div>
                      )}
                    </div>

                    {/* Progress Overlay */}
                    {uploadFile.status === "uploading" && (
                      <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                        <div className="text-white text-center">
                          <div className="w-12 h-12 border-4 border-white border-t-transparent rounded-full animate-spin mb-2 mx-auto"></div>
                          <p className="text-sm font-medium">{uploadFile.progress}%</p>
                        </div>
                      </div>
                    )}
                  </div>

                  <div className="p-4">
                    <p className="text-sm text-gray-600 truncate mb-2 font-medium">
                      {uploadFile.file.name}
                    </p>

                    {/* Status */}
                    {uploadFile.status === "success" && (
                      <p className="text-sm text-green-600 font-semibold">
                        アップロード完了
                      </p>
                    )}

                    {uploadFile.status === "error" && (
                      <div className="space-y-2">
                        <p className="text-xs text-red-600">
                          {uploadFile.error}
                        </p>
                        <button
                          onClick={() => retryUpload(uploadFile)}
                          className="text-xs text-blue-600 hover:text-blue-700 font-semibold bg-blue-50 px-2 py-1 rounded transition-colors"
                        >
                          再試行
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>

            {/* Summary */}
            {successCount > 0 && (
              <div className="mt-8 card-elevated p-6 bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center">
                    <Check className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <p className="text-lg font-bold text-green-800">
                      {successCount}枚の写真をアップロードしました！
                    </p>
                    <p className="text-green-700 mt-1">
                      他の参加者もあなたの写真を見ることができます。
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}