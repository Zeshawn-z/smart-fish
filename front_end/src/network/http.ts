import apiClient from './axios'
import type { AxiosRequestConfig } from 'axios'

export async function httpGet<T = any>(url: string, params?: Record<string, any>, config?: AxiosRequestConfig): Promise<T> {
  const response = await apiClient.get<T>(url, { params, ...config })
  return response.data
}

export async function httpPost<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  const response = await apiClient.post<T>(url, data, config)
  return response.data
}

export async function httpPut<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  const response = await apiClient.put<T>(url, data, config)
  return response.data
}

export async function httpPatch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  const response = await apiClient.patch<T>(url, data, config)
  return response.data
}

export async function httpDelete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  const response = await apiClient.delete<T>(url, config)
  return response.data
}
