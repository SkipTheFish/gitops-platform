package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/handler"
	"github.com/SkipTheFish/gitops-platform/backend/internal/service"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
	"github.com/gin-contrib/cors"
)

func New(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	// ===== 跨域配置 =====
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // 允许所有域名跨域（开发环境为了图省事可以这么写）
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// ===== store =====
	appStore := store.NewAppStore(db)
	envStore := store.NewEnvironmentStore(db)
	deploymentStore := store.NewDeploymentRecordStore(db)
	auditStore := store.NewOperationAuditStore(db)
	pipelineStore := store.NewPipelineRunStore(db)

	// infra services
	gitopsService := service.NewGitOpsService(
		getEnv("GITOPS_REPO_LOCAL_PATH", "/tmp/config-repo"),
		getEnv("GITOPS_REPO_BRANCH", "main"),
	)

	argocdService := service.NewArgoCDService(
		getEnv("ARGOCD_SERVER", "http://127.0.0.1:8081"),
		getEnv("ARGOCD_TOKEN", ""),
	)

	pipelineExecutor := service.NewPipelineExecutor(
		pipelineStore,
		appStore,
		envStore,
		deploymentStore,
		auditStore,
		gitopsService,
		argocdService,
	)

	// biz services
	appService := service.NewAppService(appStore)
	envService := service.NewEnvironmentService(envStore, appStore)
	deploymentService := service.NewDeploymentRecordService(deploymentStore, appStore, envStore, auditStore)
	auditService := service.NewOperationAuditService(auditStore)

	pipelineService := service.NewPipelineRunService(
		pipelineStore,
		appStore,
		envStore,
		deploymentStore,
		auditStore,
		pipelineExecutor,
	)

	// handlers
	healthHandler := handler.NewHealthHandler(db, rdb)
	appHandler := handler.NewAppHandler(appService)
	envHandler := handler.NewEnvironmentHandler(envService)
	deploymentHandler := handler.NewDeploymentRecordHandler(deploymentService)
	auditHandler := handler.NewOperationAuditHandler(auditService)
	pipelineHandler := handler.NewPipelineRunHandler(pipelineService)

	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Health)

		api.POST("/apps", appHandler.CreateApp)
		api.GET("/apps", appHandler.ListApps)
		api.GET("/apps/:id", appHandler.GetApp)
		api.PUT("/apps/:id", appHandler.UpdateApp)

		api.POST("/apps/:id/environments", envHandler.CreateEnvironment)
		api.GET("/apps/:id/environments", envHandler.ListEnvironments)
		api.GET("/environments/:id", envHandler.GetEnvironment)
		api.PUT("/environments/:id", envHandler.UpdateEnvironment)

		api.POST("/deployments", deploymentHandler.CreateDeploymentRecord)
		api.GET("/deployments/:id", deploymentHandler.GetDeploymentRecord)
		api.GET("/apps/:id/deployments", deploymentHandler.ListDeploymentsByApp)
		api.GET("/environments/:id/deployments", deploymentHandler.ListDeploymentsByEnv)

		api.POST("/audits", auditHandler.CreateAudit)
		api.GET("/audits/target/:id", auditHandler.ListAuditsByTargetID)

		api.POST("/pipeline-runs", pipelineHandler.CreateManualRun)
		api.GET("/pipeline-runs/:id", pipelineHandler.GetPipelineRun)
		api.GET("/apps/:id/pipeline-runs", pipelineHandler.ListPipelineRunsByApp)
		api.GET("/environments/:id/pipeline-runs", pipelineHandler.ListPipelineRunsByEnv)
	}

	return r
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
