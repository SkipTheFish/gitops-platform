import http from './http'
import type { DeploymentRecordItem } from '@/types'

export function getDeploymentsByApp(appId: number): Promise<{ data: DeploymentRecordItem[] }> {
  return http.get(`/apps/${appId}/deployments`)
}

export function getDeploymentsByEnv(envId: number): Promise<{ data: DeploymentRecordItem[] }> {
  return http.get(`/environments/${envId}/deployments`)
}

export function getDeploymentDetail(id: number): Promise<{ data: DeploymentRecordItem }> {
  return http.get(`/deployments/${id}`)
}