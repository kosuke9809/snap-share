"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Camera, Heart, Share2 } from "lucide-react"

export default function HomePage() {
  const router = useRouter()
  const [eventCode, setEventCode] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (eventCode.trim()) {
      router.push(`/e/${eventCode.trim().toUpperCase()}`)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-elegant relative overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-20 left-20 w-72 h-72 bg-white rounded-full mix-blend-multiply filter blur-xl animate-float"></div>
        <div className="absolute top-40 right-20 w-72 h-72 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl animate-float" style={{ animationDelay: '2s' }}></div>
        <div className="absolute -bottom-8 left-1/2 w-72 h-72 bg-pink-300 rounded-full mix-blend-multiply filter blur-xl animate-float" style={{ animationDelay: '4s' }}></div>
      </div>

      <div className="relative z-10 max-w-6xl mx-auto px-4 py-8 sm:py-12">
        {/* Hero Section */}
        <div className="text-center mb-16 animate-fade-in-up">
          <div className="inline-flex items-center justify-center w-24 h-24 mb-8 glass rounded-3xl animate-float">
            <Heart className="w-12 h-12 text-white" />
          </div>
          <h1 className="text-5xl sm:text-7xl font-bold text-white mb-6 text-elegant">
            SnapShare
          </h1>
          <p className="text-xl sm:text-2xl text-white/90 max-w-3xl mx-auto leading-relaxed">
            結婚式や特別なイベントで、みんなで写真を共有しよう
          </p>
        </div>

        <div className="grid lg:grid-cols-3 gap-8 mb-16">
          {/* Features */}
          <div className="lg:col-span-2 space-y-6">
            <div className="grid sm:grid-cols-2 gap-6">
              <div className="card-elevated p-8 text-center group hover:scale-105 transition-all duration-300">
                <div className="w-16 h-16 bg-gradient-primary rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform">
                  <Camera className="w-8 h-8 text-white" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-4">
                  簡単アップロード
                </h3>
                <p className="text-gray-600 leading-relaxed">
                  QRコードを読み取って、名前を入力するだけで写真をアップロード
                </p>
              </div>

              <div className="card-elevated p-8 text-center group hover:scale-105 transition-all duration-300">
                <div className="w-16 h-16 bg-gradient-primary rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform">
                  <Share2 className="w-8 h-8 text-white" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-4">
                  リアルタイム共有
                </h3>
                <p className="text-gray-600 leading-relaxed">
                  アップロードした写真は即座に他の参加者と共有
                </p>
              </div>
            </div>

            <div className="card-elevated p-8 text-center group hover:scale-105 transition-all duration-300">
              <div className="w-16 h-16 bg-gradient-primary rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform">
                <Heart className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-4">
                思い出を保存
              </h3>
              <p className="text-gray-600 leading-relaxed">
                すべての写真をまとめてダウンロードして、永遠の思い出として保存
              </p>
            </div>
          </div>

          {/* CTA */}
          <div className="lg:col-span-1">
            <div className="card-premium p-8 h-full flex flex-col justify-center">
              <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
                イベントに参加
              </h2>
              <p className="text-gray-600 mb-8 text-center leading-relaxed">
                イベントコードを入力して参加してください
              </p>

              <form onSubmit={handleSubmit} className="space-y-6" role="form" aria-label="イベント参加フォーム">
                <div>
                  <input
                    type="text"
                    value={eventCode}
                    onChange={(e) => setEventCode(e.target.value)}
                    placeholder="WEDDING1"
                    className="input-modern w-full text-center font-semibold text-lg tracking-widest uppercase"
                    maxLength={8}
                    aria-label="イベントコード"
                    aria-describedby="event-code-help"
                    autoComplete="off"
                  />
                  <p id="event-code-help" className="text-sm text-gray-500 mt-2 text-center">
                    8文字のコードを入力
                  </p>
                </div>

                <button
                  type="submit"
                  className="btn-primary w-full text-lg font-semibold"
                  aria-label="イベントに参加する"
                >
                  参加する
                </button>
              </form>

              {/* Sample Codes */}
              <div className="mt-8 pt-6 border-t border-gray-200">
                <p className="text-sm text-gray-500 mb-4 text-center">開発用サンプル</p>
                <div className="flex flex-wrap gap-2 justify-center">
                  {["WEDDING1", "TRAVEL02", "REUNION2"].map((code) => (
                    <button
                      key={code}
                      onClick={() => router.push(`/e/${code}`)}
                      className="px-3 py-2 text-xs bg-gray-100 text-gray-600 rounded-lg hover:bg-gray-200 transition-colors font-mono"
                    >
                      {code}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="text-center">
          <p className="text-white/70 text-lg">
            SnapShare - みんなで作る、特別な思い出
          </p>
        </div>
      </div>
    </div>
  )
}
