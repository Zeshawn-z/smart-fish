import { httpPost, httpGet, httpPut } from '@/network/http'
import type { User, LoginInput, RegisterInput, TokenPair } from '@/types'
import CryptoJS from 'crypto-js'

// SHA256 加密密码（与旧项目保持一致）
function hashPassword(password: string): string {
  return CryptoJS.SHA256(password).toString()
}

export const AuthService = {
  async login(input: LoginInput): Promise<TokenPair> {
    return httpPost<TokenPair>('/api/auth/login', {
      username: input.username,
      password: hashPassword(input.password)
    })
  },

  async register(input: RegisterInput): Promise<{ message: string; user: User }> {
    return httpPost('/api/auth/register', {
      ...input,
      password: hashPassword(input.password)
    })
  },

  async refreshToken(refreshToken: string): Promise<{ access_token: string }> {
    return httpPost('/api/auth/refresh', { refresh_token: refreshToken })
  },

  async getMe(): Promise<User> {
    return httpGet<User>('/api/auth/me')
  },

  async updateMe(data: { phone?: string; email?: string }): Promise<User> {
    return httpPut<User>('/api/auth/me', data)
  },

  async updatePassword(oldPassword: string, newPassword: string): Promise<{ message: string }> {
    return httpPut('/api/auth/password', {
      old_password: hashPassword(oldPassword),
      new_password: hashPassword(newPassword)
    })
  }
}
