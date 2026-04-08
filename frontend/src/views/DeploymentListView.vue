<template>
  <div>
    <PageHeader
      title="部署记录"
      description="集中查看各应用的部署历史、版本信息和同步状态"
    />

    <el-card shadow="hover" class="card">
      <el-table :data="rows" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="app_id" label="应用 ID" width="100" />
        <el-table-column prop="env_id" label="环境 ID" width="100" />
        <el-table-column prop="version" label="版本" min-width="120" />
        <el-table-column prop="image_tag" label="镜像 Tag" min-width="180" />
        <el-table-column prop="git_commit" label="Git Commit" min-width="140" />
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
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import StatusTag from '@/components/StatusTag.vue'
import { getAppList } from '@/api/app'
import { getDeploymentsByApp } from '@/api/deployment'
import type { AppItem, DeploymentRecordItem } from '@/types'

const loading = ref(false)
const rows = ref<DeploymentRecordItem[]>([])

function formatTime(value: string) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

// 这里因为当前后端还没有“全局部署记录列表接口”
// 所以前端做一个聚合：先拿所有 app，再逐个拉部署记录。
async function loadData() {
  try {
    loading.value = true

    const appRes = await getAppList()
    const apps: AppItem[] = appRes.data || []

    const deploymentResults = await Promise.all(
      apps.map((app) =>
        getDeploymentsByApp(app.id)
          .then((res) => res.data || [])
          .catch(() => []),
      ),
    )

    rows.value = deploymentResults
      .flat()
      .sort((a, b) => new Date(b.deployed_at).getTime() - new Date(a.deployed_at).getTime())
  } catch (error: any) {
    ElMessage.error(error.message || '加载部署记录失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>