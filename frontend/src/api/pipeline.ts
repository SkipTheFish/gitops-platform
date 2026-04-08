import http from './http'
import type { PipelineRunItem } from '@/types'

export interface CreatePipelineRunPayload {
  app_id: number
  env_id: number
  git_commit?: string
  branch?: string
  image_tag?: string
  trigger_type?: string
  operator: string
  version: string
}

export function createPipelineRun(
  payload: CreatePipelineRunPayload,
): Promise<{ message: string; data: PipelineRunItem }> {
  return http.post('/pipeline-runs', payload)
}

export function getPipelineRunDetail(id: number): Promise<{ data: PipelineRunItem }> {
  return http.get(`/pipeline-runs/${id}`)
}

export function getPipelineRunsByApp(appId: number): Promise<{ data: PipelineRunItem[] }> {
  return http.get(`/apps/${appId}/pipeline-runs`)
}

export function getPipelineRunsByEnv(envId: number): Promise<{ data: PipelineRunItem[] }> {
  return http.get(`/environments/${envId}/pipeline-runs`)
}