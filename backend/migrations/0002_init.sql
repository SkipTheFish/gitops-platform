-- 环境管理 + 发布记录模型增强

-- 1. 确保同一个应用下，环境名称唯一
CREATE UNIQUE INDEX IF NOT EXISTS uk_environment_app_env
ON environment (app_id, env_name);

-- 2. deployment_record 常用查询索引
CREATE INDEX IF NOT EXISTS idx_deployment_record_app_id
ON deployment_record (app_id);

CREATE INDEX IF NOT EXISTS idx_deployment_record_env_id
ON deployment_record (env_id);

CREATE INDEX IF NOT EXISTS idx_deployment_record_deployed_at
ON deployment_record (deployed_at DESC);

-- 3. operation_audit 常用查询索引
CREATE INDEX IF NOT EXISTS idx_operation_audit_target_id
ON operation_audit (target_id);

CREATE INDEX IF NOT EXISTS idx_operation_audit_created_at
ON operation_audit (created_at DESC);