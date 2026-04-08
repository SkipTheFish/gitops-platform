-- 流水线运行记录索引增强

CREATE INDEX IF NOT EXISTS idx_pipeline_run_app_id
ON pipeline_run (app_id);

CREATE INDEX IF NOT EXISTS idx_pipeline_run_env_id
ON pipeline_run (env_id);

CREATE INDEX IF NOT EXISTS idx_pipeline_run_status
ON pipeline_run (status);

CREATE INDEX IF NOT EXISTS idx_pipeline_run_started_at
ON pipeline_run (started_at DESC);