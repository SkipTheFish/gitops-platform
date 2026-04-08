import http from './http'
import type { AppItem } from '@/types'

export interface CreateAppPayload {
  name: string
  repo_url: string
  config_repo_url: string
  cluster_name: string
  namespace: string
  helm_chart_path: string
  values_file_path: string
}

export interface UpdateAppPayload {
  name?: string
  repo_url?: string
  config_repo_url?: string
  cluster_name?: string
  namespace?: string
  helm_chart_path?: string
  values_file_path?: string
}

export function getAppList(): Promise<{ data: AppItem[] }> {
  return http.get('/apps')
}

export function getAppDetail(id: number): Promise<{ data: AppItem }> {
  return http.get(`/apps/${id}`)
}

export function createApp(payload: CreateAppPayload): Promise<{ message: string; data: AppItem }> {
  return http.post('/apps', payload)
}

export function updateApp(id: number, payload: UpdateAppPayload): Promise<{ message: string; data: AppItem }> {
  return http.put(`/apps/${id}`, payload)
}