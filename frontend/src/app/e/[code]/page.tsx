"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { useForm } from "react-hook-form"
import { Camera, Heart, Users } from "lucide-react"
import { apiClient } from "@/lib/api"
import { useAuthStore } from "@/stores/auth"
import type { Event } from "@/types/api"

interface GuestFormData {
  guest_name: string
}

interface PageProps {
  params: Promise<{ code: string }>
}

export default function EventLandingPage({ params }: PageProps) {
  const router = useRouter()
  const [event, setEvent] = useState<Event | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { setSession } = useAuthStore()
  const [code, setCode] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<GuestFormData>()

  // Unwrap params
  useEffect(() => {
    params.then((resolvedParams) => {
      setCode(resolvedParams.code)
    })
  }, [params])

  // Load event data
  useEffect(() => {
    if (!code) return

    const loadEvent = async () => {
      try {
        const eventData = await apiClient.getEventByCode(code)
        setEvent(eventData)
      } catch (err) {
        setError(err instanceof Error ? err.message : "イベントが見つかりません")
      } finally {
        setLoading(false)
      }
    }

    loadEvent()
  }, [code])

  const onSubmit = async (data: GuestFormData) => {
    if (!code) return

    setSubmitting(true)
    setError(null)

    try {
      const session = await apiClient.createSession({
        event_code: code,
        guest_name: data.guest_name
      })

      setSession(session)
      router.push("/upload")
    } catch (err) {
      setError(err instanceof Error ? err.message : "セッションの作成に失敗しました")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-pink-50">
        <div className="text-center">
          <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </div>
    )
  }

  if (error && !event) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-pink-50">
        <div className="text-center p-8">
          <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-red-600 text-2xl">✕</span>
          </div>
          <h1 className="text-xl font-bold text-gray-900 mb-2">
            エラーが発生しました
          </h1>
          <p className="text-gray-600">{error}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-soft relative overflow-hidden py-6 px-4 flex flex-col">
      {/* Background Elements */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute top-10 left-10 w-40 h-40 bg-white rounded-full mix-blend-multiply filter blur-xl animate-float"></div>
        <div className="absolute bottom-10 right-10 w-40 h-40 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl animate-float" style={{ animationDelay: '3s' }}></div>
      </div>

      <div className="relative z-10 max-w-md mx-auto w-full flex-1 flex flex-col justify-center">
        {/* Event Header */}
        <div className="text-center mb-12 animate-fade-in-up">
          <div className="inline-flex items-center justify-center w-28 h-28 mb-8 glass rounded-3xl animate-float">
            <Heart className="w-14 h-14 text-white" />
          </div>
          <h1 className="text-3xl sm:text-4xl font-bold text-white mb-4 px-2 text-elegant">
            {event?.name}
          </h1>
          {event?.description && (
            <p className="text-lg text-white/90 px-2 leading-relaxed">
              {event.description}
            </p>
          )}
          {event?.event_date && (
            <p className="text-sm text-white/80 mt-4 font-medium">
              {new Date(event.event_date).toLocaleDateString("ja-JP", {
                year: "numeric",
                month: "long",
                day: "numeric"
              })}
            </p>
          )}
        </div>

        {/* Feature Icons */}
        <div className="grid grid-cols-2 gap-4 mb-8">
          <div className="card-elevated p-6 text-center group hover:scale-105 transition-all duration-300">
            <div className="w-12 h-12 bg-gradient-primary rounded-xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform">
              <Camera className="w-6 h-6 text-white" />
            </div>
            <p className="text-sm font-semibold text-gray-900">写真をシェア</p>
          </div>
          <div className="card-elevated p-6 text-center group hover:scale-105 transition-all duration-300">
            <div className="w-12 h-12 bg-gradient-primary rounded-xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform">
              <Users className="w-6 h-6 text-white" />
            </div>
            <p className="text-sm font-semibold text-gray-900">みんなで共有</p>
          </div>
        </div>

        {/* Guest Form */}
        <div className="card-premium p-8 animate-fade-in-up" style={{ animationDelay: '0.3s' }}>
          <h2 className="text-2xl font-bold text-gray-900 mb-8 text-center">
            お名前を入力してください
          </h2>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div>
              <input
                {...register("guest_name", {
                  required: "お名前を入力してください",
                  minLength: {
                    value: 1,
                    message: "お名前を入力してください"
                  },
                  maxLength: {
                    value: 100,
                    message: "お名前は100文字以内で入力してください"
                  }
                })}
                type="text"
                placeholder="田中 太郎"
                className="input-modern w-full text-center text-lg font-medium"
                disabled={submitting}
              />
              {errors.guest_name && (
                <p className="text-sm text-red-600 mt-3 text-center font-medium">
                  {errors.guest_name.message}
                </p>
              )}
            </div>

            {error && (
              <div className="bg-red-50 border-2 border-red-200 rounded-2xl p-4">
                <p className="text-sm text-red-600 text-center font-medium">{error}</p>
              </div>
            )}

            <button
              type="submit"
              disabled={submitting}
              className="btn-primary w-full text-lg font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {submitting ? (
                <div className="flex items-center justify-center">
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-3" />
                  参加中...
                </div>
              ) : (
                "写真をシェアする"
              )}
            </button>
          </form>
        </div>

        {/* Footer */}
        <div className="text-center mt-8 pb-4 animate-fade-in-up" style={{ animationDelay: '0.6s' }}>
          <p className="text-white/80 px-4 leading-relaxed">
            このページから写真をアップロードして、<br />
            みんなで素敵な思い出を共有しましょう！
          </p>
        </div>
      </div>
    </div>
  )
}