<template>
  <div v-loading="loading">
    <PageHeader
      :title="appDetail?.name || '应用详情'"
      description="查看应用基本信息、管理环境以及查看部署记录"
    >
      <el-button @click="goBack">返回列表</el-button>
      <el-button type="primary" @click="openEditAppDialog" v-if="appDetail">
        编辑应用
      </el-button>
      <el-button type="success" @click="openPipelineDialog" v-if="appDetail">
        手动发布
      </el-button>
    </PageHeader>

    <el-row :gutter="16" v-if="appDetail">
      <el-col :span="24">
        <el-card shadow="hover" class="card">
          <template #header>
            <div class="card-header">应用基本信息</div>
          </template>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="应用名称">{{ appDetail.name }}</el-descriptions-item>
            <el-descriptions-item label="集群">{{ appDetail.cluster_name }}</el-descriptions-item>
            <el-descriptions-item label="命名空间">{{ appDetail.namespace }}</el-descriptions-item>
            <el-descriptions-item label="代码仓库">{{ appDetail.repo_url }}</el-descriptions-item>
            <el-descriptions-item label="配置仓库">{{ appDetail.config_repo_url }}</el-descriptions-item>
            <el-descriptions-item label="Chart 路径">{{ appDetail.helm_chart_path }}</el-descriptions-item>
            <el-descriptions-item label="Values 路径">{{ appDetail.values_file_path }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatTime(appDetail.updated_at) }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="hover" class="card" v-if="appDetail">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="环境管理" name="env">
          <div class="section-toolbar">
            <el-button type="primary" @click="openCreateEnvDialog">新增环境</el-button>
          </div>

          <el-row :gutter="16">
            <el-col :span="8" v-for="env in environments" :key="env.id">
              <el-card shadow="hover" class="env-card">
                <div class="env-title-row">
                  <div class="env-title">{{ env.env_name }}</div>
                  <el-switch
                    :model-value="env.auto_sync_enabled"
                    disabled
                  />
                </div>

                <div class="env-item"><span>集群：</span>{{ env.cluster_name }}</div>
                <div class="env-item"><span>命名空间：</span>{{ env.namespace }}</div>
                <div class="env-item"><span>自动同步：</span>{{ env.auto_sync_enabled ? '开启' : '关闭' }}</div>

                <div class="env-actions">
                  <el-button size="small" type="primary" plain @click="openEditEnvDialog(env)">
                    编辑环境
                  </el-button>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <el-tab-pane label="部署记录" name="deploy">
          <el-table :data="deployments" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="version" label="版本" min-width="120" />
            <el-table-column prop="image_tag" label="镜像 Tag" min-width="180" />
            <el-table-column prop="git_commit" label="Git Commit" min-width="140" />
            <el-table-column prop="argocd_app_name" label="ArgoCD App" min-width="180" />
            <el-table-column label="Sync 状态" min-width="120">
              <template #default="{ row }">
                <StatusTag :title="row.sync_status" />
              </template>
            </el-table-column>
            <el-table-column label="Health 状态" min-width="120">
              <template #default="{ row }">
                <StatusTag :title="row.health_status" />
              </template>
            </el-table-column>
            <el-table-column prop="operator" label="操作人" min-width="120" />
            <el-table-column label="部署时间" min-width="180">
              <template #default="{ row }">
                {{ formatTime(row.deployed_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="流水线记录" name="pipeline">
          <el-table :data="pipelines" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="branch" label="分支" min-width="120" />
            <el-table-column prop="git_commit" label="Git Commit" min-width="140" />
            <el-table-column prop="image_tag" label="镜像 Tag" min-width="180" />
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <StatusTag :title="row.status" />
              </template>
            </el-table-column>
            <el-table-column prop="trigger_type" label="触发方式" min-width="120" />
            <el-table-column label="开始时间" min-width="180">
              <template #default="{ row }">
                {{ formatTime(row.started_at) }}
              </template>
            </el-table-column>
            <el-table-column label="结束时间" min-width="180">
              <template #default="{ row }">
                {{ formatTime(row.finished_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>


      </el-tabs>
    </el-card>

    <AppFormDialog
      :visible="appDialogVisible"
      :app="appDetail"
      @close="appDialogVisible = false"
      @success="loadAll"
    />

    <EnvironmentFormDialog
      :visible="envDialogVisible"
      :app-id="appId"
      :environment="currentEnvironment"
      @close="envDialogVisible = false"
      @success="loadEnvironments"
    />

    <PipelineRunFormDialog
      :visible="pipelineDialogVisible"
      :app-id="appId"
      :environments="environments"
      @close="pipelineDialogVisible = false"
      @success="loadPipelines"
    />
  </div>
</template>

<script setup lang="ts">

import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import StatusTag from '@/components/StatusTag.vue'
import AppFormDialog from '@/components/AppFormDialog.vue'
import EnvironmentFormDialog from '@/components/EnvironmentFormDialog.vue'
import { getAppDetail } from '@/api/app'
import { getEnvironmentListByApp } from '@/api/environment'
import { getDeploymentsByApp } from '@/api/deployment'
import type { AppItem, DeploymentRecordItem, EnvironmentItem } from '@/types'

import { getPipelineRunsByApp } from '@/api/pipeline'
import PipelineRunFormDialog from '@/components/PipelineRunFormDialog.vue'
import type { PipelineRunItem } from '@/types'

const route = useRoute()
const router = useRouter()

const appId = Number(route.params.id)
const loading = ref(false)
const activeTab = ref('env')

const appDetail = ref<AppItem | null>(null)
const environments = ref<EnvironmentItem[]>([])
const deployments = ref<DeploymentRecordItem[]>([])

const appDialogVisible = ref(false)
const envDialogVisible = ref(false)
const currentEnvironment = ref<EnvironmentItem | null>(null)

const pipelines = ref<PipelineRunItem[]>([])
const pipelineDialogVisible = ref(false)
let refreshTimer: number | null = null

function formatTime(value: string) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

async function loadAppDetail() {
  const res = await getAppDetail(appId)
  appDetail.value = res.data
}

async function loadEnvironments() {
  const res = await getEnvironmentListByApp(appId)
  environments.value = res.data || []
}

async function loadDeployments() {
  const res = await getDeploymentsByApp(appId)
  deployments.value = res.data || []
}

async function loadAll() {
  try {
    loading.value = true
    await Promise.all([
      loadAppDetail(),
      loadEnvironments(),
      loadDeployments(),
      loadPipelines(),
    ])
  } catch (error: any) {
    ElMessage.error(error.message || '加载应用详情失败')
  } finally {
    loading.value = false
  }
}

async function loadPipelines() {
  const res = await getPipelineRunsByApp(appId)
  pipelines.value = res.data || []
}

function goBack() {
  router.push('/apps')
}

function openEditAppDialog() {
  appDialogVisible.value = true
}

function openCreateEnvDialog() {
  currentEnvironment.value = null
  envDialogVisible.value = true
}

function openEditEnvDialog(env: EnvironmentItem) {
  currentEnvironment.value = env
  envDialogVisible.value = true
}

function openPipelineDialog() {
  pipelineDialogVisible.value = true
}

function startPolling() {
  refreshTimer = window.setInterval(async () => {
    await Promise.all([loadPipelines(), loadDeployments()])
  }, 3000)
}

import { onMounted, onUnmounted, ref } from 'vue'

onMounted(() => {
  loadAll()
  startPolling()
})

onUnmounted(() => {
  if (refreshTimer) {
    window.clearInterval(refreshTimer)
  }
})
</script>