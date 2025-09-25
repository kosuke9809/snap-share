import { create } from "zustand"
import { persist } from "zustand/middleware"
import type { Session } from "@/types/api"
import { apiClient } from "@/lib/api"

interface AuthState {
  session: Session | null
  isAuthenticated: boolean
  setSession: (session: Session | null) => void
  clearSession: () => void
  refreshSession: () => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      session: null,
      isAuthenticated: false,

      setSession: (session: Session | null) => {
        if (session) {
          apiClient.setAuthToken(session.session_token)
        } else {
          apiClient.setAuthToken(null)
        }

        set({
          session,
          isAuthenticated: !!session
        })
      },

      clearSession: () => {
        apiClient.setAuthToken(null)
        set({
          session: null,
          isAuthenticated: false
        })
      },

      refreshSession: async () => {
        const { session } = get()
        if (!session) return

        try {
          const refreshedSession = await apiClient.refreshSession({
            session_token: session.session_token,
          })
          get().setSession(refreshedSession)
        } catch (error) {
          console.error("Failed to refresh session:", error)
          get().clearSession()
        }
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        session: state.session,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
)