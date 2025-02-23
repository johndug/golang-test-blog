import { defineStore } from 'pinia'
import axios from 'axios'
import type { User, LoginCredentials, RegisterData, AuthResponse } from '../types'

interface AuthState {
  user: User | null
  token: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: localStorage.getItem('token')
  }),

  getters: {
    isAuthenticated: (state): boolean => !!state.token,
    currentUser: (state): User | null => state.user
  },

  actions: {
    async login(credentials: LoginCredentials): Promise<void> {
      try {
        const response = await axios.post<AuthResponse>('/api/login', credentials)
        this.token = response.data.token
        this.user = response.data.user
        localStorage.setItem('token', this.token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      } catch (error: any) {
        throw new Error(error.response?.data?.error || 'Login failed')
      }
    },

    async register(data: RegisterData): Promise<void> {
      try {
        const response = await axios.post<AuthResponse>('/api/register', data, {
          headers: {
            'Content-Type': 'application/json'
          }
        })

        this.token = response.data.token
        this.user = response.data.user
        localStorage.setItem('token', this.token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      } catch (error: any) {
        throw new Error(error.response?.data?.error || 'Registration failed')
      }
    },

    logout(): void {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
    },

    async fetchCurrentUser() {
      try {
        if (!this.token) return

        const response = await axios.get<User>('/api/me', {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        })
        this.user = response.data
      } catch (error) {
        this.logout()
      }
    },

    async initializeAuth() {
      const token = localStorage.getItem('token')
      if (token) {
        this.token = token
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        await this.fetchCurrentUser()
      }
    }
  }
})
