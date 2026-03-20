/**
 * v1 API 认证 header 工具
 * v1 接口使用 Bearer token（Go 后端的 access_token 同时被 FlaskAuthRequired 识别）
 */

import type { AxiosRequestConfig } from 'axios'

export function v1AuthHeaders(): AxiosRequestConfig {
  const token = localStorage.getItem('access_token')
  return token ? { headers: { Authorization: `Bearer ${token}` } } : {}
}
