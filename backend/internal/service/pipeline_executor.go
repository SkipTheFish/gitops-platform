package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
)

type PipelineExecutor struct {
	pipelineStore   *store.PipelineRunStore
	appStore        *store.AppStore
	envStore        *store.EnvironmentStore
	deploymentStore *store.DeploymentRecordStore
	auditStore      *store.OperationAuditStore

	gitopsService *GitOpsService
	argocdService *ArgoCDService
}

func NewPipelineExecutor(
	pipelineStore *store.PipelineRunStore,
	appStore *store.AppStore,
	envStore *store.EnvironmentStore,
	deploymentStore *store.DeploymentRecordStore,
	auditStore *store.OperationAuditStore,
	gitopsService *GitOpsService,
	argocdService *ArgoCDService,
) *PipelineExecutor {
	return &PipelineExecutor{
		pipelineStore:   pipelineStore,
		appStore:        appStore,
		envStore:        envStore,
		deploymentStore: deploymentStore,
		auditStore:      auditStore,
		gitopsService:   gitopsService,
		argocdService:   argocdService,
	}
}

func (e *PipelineExecutor) Execute(runID int64, input CreatePipelineRunInput) {
	ctx := context.Background()

	run, err := e.pipelineStore.GetByID(ctx, runID)
	if err != nil {
		return
	}

	app, err := e.appStore.GetByID(ctx, input.AppID)
	if err != nil {
		e.failRun(ctx, run, input.Operator, "app not found")
		return
	}

	env, err := e.envStore.GetByID(ctx, input.EnvID)
	if err != nil {
		e.failRun(ctx, run, input.Operator, "environment not found")
		return
	}

	if strings.TrimSpace(env.ValuesFilePath) == "" {
		e.failRun(ctx, run, input.Operator, "environment values_file_path is empty")
		return
	}

	if strings.TrimSpace(env.ArgoCDAppName) == "" {
		e.failRun(ctx, run, input.Operator, "environment argocd_app_name is empty")
		return
	}

	now := time.Now()
	run.Status = "running"
	run.StartedAt = &now
	run.LogURL = fmt.Sprintf("/fake-logs/pipeline/%d", run.ID)
	_ = e.pipelineStore.Update(ctx, run)

	updateResult, err := e.gitopsService.UpdateImageTagAndPush(UpdateValuesInput{
		ValuesFilePath: env.ValuesFilePath,
		ImageTag:       input.ImageTag,
	})
	if err != nil {
		e.failRun(ctx, run, input.Operator, "update gitops repo failed: "+err.Error())
		return
	}

	if !env.AutoSyncEnabled {
		if err := e.argocdService.SyncApplication(env.ArgoCDAppName); err != nil {
			e.failRun(ctx, run, input.Operator, "argocd sync failed: "+err.Error())
			return
		}
	}

	status, err := e.waitForArgoCDHealthy(env.ArgoCDAppName, 2*time.Minute)
	if err != nil {
		e.failRun(ctx, run, input.Operator, "argocd wait failed: "+err.Error())
		return
	}

	finishedAt := time.Now()
	run.Status = "success"
	run.FinishedAt = &finishedAt
	_ = e.pipelineStore.Update(ctx, run)

	record := &model.DeploymentRecord{
		AppID:           app.ID,
		EnvID:           env.ID,
		Version:         input.Version,
		ImageTag:        input.ImageTag,
		GitCommit:       input.GitCommit,
		ConfigCommitSHA: updateResult.CommitSHA,
		ArgoCDAppName:   env.ArgoCDAppName,
		SyncStatus:      status.SyncStatus,
		HealthStatus:    status.HealthStatus,
		Operator:        input.Operator,
	}

	if err := e.deploymentStore.Create(ctx, record); err != nil {
		e.failRun(ctx, run, input.Operator, "create deployment record failed: "+err.Error())
		return
	}

	_ = e.auditStore.Create(ctx, &model.OperationAudit{
		Operator:   input.Operator,
		ActionType: "deploy",
		TargetID:   record.ID,
		Detail: fmt.Sprintf(
			"deploy success, app=%s env=%s version=%s config_commit=%s",
			app.Name,
			env.EnvName,
			input.Version,
			updateResult.CommitSHA,
		),
	})
}

func (e *PipelineExecutor) waitForArgoCDHealthy(appName string, timeout time.Duration) (*ArgoCDAppStatus, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		status, err := e.argocdService.GetApplicationStatus(appName)
		if err == nil && status.SyncStatus == "Synced" && status.HealthStatus == "Healthy" {
			return status, nil
		}
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("timeout waiting for argocd app to become synced and healthy")
}

func (e *PipelineExecutor) failRun(ctx context.Context, run *model.PipelineRun, operator string, detail string) {
	finishedAt := time.Now()
	run.Status = "failed"
	run.FinishedAt = &finishedAt
	_ = e.pipelineStore.Update(ctx, run)

	_ = e.auditStore.Create(ctx, &model.OperationAudit{
		Operator:   operator,
		ActionType: "deploy_failed",
		TargetID:   run.ID,
		Detail:     detail,
	})
}
