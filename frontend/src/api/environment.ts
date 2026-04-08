import http from './http'
import type { EnvironmentItem } from '@/types'

export interface CreateEnvironmentPayload {
  env_name: string
  cluster_name: string
  namespace: string
  auto_sync_enabled: boolean
}

export interface UpdateEnvironmentPayload {
  env_name?: string
  cluster_name?: string
  namespace?: string
  auto_sync_enabled?: boolean
}

export function getEnvironmentListByApp(appId: number): Promise<{ data: EnvironmentItem[] }> {
  return http.get(`/apps/${appId}/environments`)
}

export function getEnvironmentDetail(id: number): Promise<{ data: EnvironmentItem }> {
  return http.get(`/environments/${id}`)
}

export function createEnvironment(
  appId: number,
  payload: CreateEnvironmentPayload,
): Promise<{ message: string; data: EnvironmentItem }> {
  return http.post(`/apps/${appId}/environments`, payload)
}

export function updateEnvironment(
  id: number,
  payload: UpdateEnvironmentPayload,
): Promise<{ message: string; data: EnvironmentItem }> {
  return http.put(`/environments/${id}`, payload)
}