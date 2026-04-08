import http from './http'
import type { AuditItem } from '@/types'

export function getAuditListByTarget(targetId: number): Promise<{ data: AuditItem[] }> {
  return http.get(`/audits/target/${targetId}`)
}