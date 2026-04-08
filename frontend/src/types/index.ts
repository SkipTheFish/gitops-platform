// 这里集中定义前端会用到的类型。
// 好处是：接口返回的数据结构一眼能看清，后面写页面时也有类型提示。

export interface AppItem {
  id: number
  name: string
  repo_url: string
  config_repo_url: string
  cluster_name: string
  namespace: string
  helm_chart_path: string
  values_file_path: string
  created_at: string
  updated_at: string
}

export interface EnvironmentItem {
  id: number
  app_id: number
  env_name: string
  cluster_name: string
  namespace: string
  auto_sync_enabled: boolean
  created_at: string
  updated_at: string
}

export interface DeploymentRecordItem {
  id: number
  app_id: number
  env_id: number
  version: string
  image_tag: string
  git_commit: string
  argocd_app_name: string
  sync_status: string
  health_status: string
  operator: string
  deployed_at: string
  rollback_from_version?: string | null
}

export interface AuditItem {
  id: number
  operator: string
  action_type: string
  target_id: number
  detail: string
  created_at: string
}

export interface PipelineRunItem {
  id: number
  app_id: number
  env_id: number
  git_commit: string
  branch: string
  image_tag: string
  status: string
  trigger_type: string
  started_at?: string
  finished_at?: string
  log_url: string
}