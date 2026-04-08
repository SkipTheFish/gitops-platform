-- 1. 应用基础信息表
-- 用于存储被管理应用的元数据，如仓库地址、Helm Chart 路径等

CREATE TABLE IF NOT EXISTs app (
    id BIGSERIAL PRIMARY KEY,                   -- 主键 ID，自增序列
    name VARCHAR(100) NOT NULL UNIQUE,          -- 应用名称，不可重复
    repo_url TEXT NOT NULL,                     -- 源代码仓库地址 (Git Repo)
    config_repo_url TEXT NOT NULL,              -- 配置仓库地址 (GitOps Config Repo)
    cluster_name VARCHAR(100) NOT NULL,         -- 目标集群名称
    namespace VARCHAR(100) NOT NULL,            -- 默认部署命名空间
    helm_chart_path TEXT NOT NULL,              -- Helm Chart 的相对路径
    values_file_path TEXT NOT NULL,             -- 默认 values 配置文件路径
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),    -- 记录创建时间
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()     -- 记录最后更新时间
);

-- 2. 环境配置表
-- 定义应用在不同环境（如开发、测试、生产）下的具体部署配置
CREATE TABLE IF NOT EXISTS environment (
    id BIGSERIAL PRIMARY KEY,              -- 主键 ID
    app_id BIGINT NOT NULL REFERENCES app(id) ON DELETE CASCADE, -- 关联应用ID，应用删除时级联删除环境
    env_name VARCHAR(50) NOT NULL,         -- 环境名称 (如: dev, staging, prod)
    cluster_name VARCHAR(100) NOT NULL,    -- 该环境所在的集群名称
    namespace VARCHAR(100) NOT NULL,       -- 该环境所在的命名空间
    auto_sync_enabled BOOLEAN NOT NULL DEFAULT FALSE, -- 是否开启自动同步 (GitOps)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- 更新时间
    UNIQUE(app_id, env_name)               -- 联合唯一约束：同一个应用下，环境名称不能重复
);

-- 3. 流水线运行记录表
-- 记录 CI/CD 流水线的执行情况，关联具体的应用和环境
CREATE TABLE IF NOT EXISTS pipeline_run (
    id BIGSERIAL PRIMARY KEY,              -- 主键 ID
    app_id BIGINT NOT NULL REFERENCES app(id) ON DELETE CASCADE, -- 关联应用ID
    env_id BIGINT NOT NULL REFERENCES environment(id) ON DELETE CASCADE, -- 关联环境ID
    git_commit VARCHAR(100),               -- 触发流水线的 Git Commit Hash
    branch VARCHAR(100),                   -- 代码分支名称
    image_tag VARCHAR(255),                -- 构建生成的镜像 Tag
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 运行状态 (pending, running, success, failed)
    trigger_type VARCHAR(50) NOT NULL,     -- 触发类型 (webhook, manual, schedule)
    started_at TIMESTAMP,                  -- 开始时间
    finished_at TIMESTAMP,                 -- 结束时间
    log_url TEXT                           -- 日志查看链接
);

-- 4. 部署历史/发布记录表
-- 记录每一次实际发生的生产/测试环境部署操作，用于版本回溯和审计
CREATE TABLE IF NOT EXISTS deployment_record (
    id BIGSERIAL PRIMARY KEY,              -- 主键 ID
    app_id BIGINT NOT NULL REFERENCES app(id) ON DELETE CASCADE, -- 关联应用ID
    env_id BIGINT NOT NULL REFERENCES environment(id) ON DELETE CASCADE, -- 关联环境ID
    version VARCHAR(100) NOT NULL,         -- 部署版本号 (如 Helm Chart Version)
    image_tag VARCHAR(255),                -- 实际部署的镜像 Tag
    git_commit VARCHAR(100),               -- 对应的代码版本
    argocd_app_name VARCHAR(255),          -- ArgoCD 中的应用名称 (如果是 GitOps 模式)
    sync_status VARCHAR(50),               -- 同步状态 (Synced, OutOfSync)
    health_status VARCHAR(50),             -- 健康状态 (Healthy, Degraded)
    operator VARCHAR(100),                 -- 操作人 (系统自动或人工)
    deployed_at TIMESTAMP NOT NULL DEFAULT NOW(), -- 部署时间
    rollback_from_version VARCHAR(100)     -- 如果是回滚操作，记录回滚来源的版本
);

-- 5. 操作审计日志表
-- 记录系统内的关键操作行为，用于安全审计和故障排查
CREATE TABLE IF NOT EXISTS operation_audit (
    id BIGSERIAL PRIMARY KEY,              -- 主键 ID
    operator VARCHAR(100) NOT NULL,        -- 操作人/账号
    action_type VARCHAR(50) NOT NULL,      -- 动作类型 (创建, 删除, 更新, 部署)
    target_id BIGINT NOT NULL,             -- 操作目标的 ID (如 app_id 或 env_id)
    detail TEXT,                           -- 操作详情或变更内容快照
    created_at TIMESTAMP NOT NULL DEFAULT NOW() -- 操作时间
);

















