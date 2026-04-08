<template>
  <div>
    <PageHeader
      title="审计日志"
      description="查看系统中的关键操作记录，便于追踪发布与变更"
    />

    <el-card shadow="hover" class="card">
      <div class="toolbar">
        <el-input
          v-model="targetId"
          placeholder="输入 target_id，例如 deployment_record.id"
          style="width: 280px"
          clearable
        />
        <el-button type="primary" @click="loadData">查询</el-button>
      </div>

      <el-table :data="rows" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="operator" label="操作人" min-width="120" />
        <el-table-column prop="action_type" label="动作类型" min-width="140" />
        <el-table-column prop="target_id" label="目标 ID" min-width="100" />
        <el-table-column prop="detail" label="详情" min-width="320" />
        <el-table-column label="时间" min-width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { getAuditListByTarget } from '@/api/audit'
import type { AuditItem } from '@/types'

const targetId = ref('')
const loading = ref(false)
const rows = ref<AuditItem[]>([])

function formatTime(value: string) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

async function loadData() {
  try {
    if (!targetId.value) {
      ElMessage.warning('请先输入 target_id')
      return
    }

    loading.value = true
    const res = await getAuditListByTarget(Number(targetId.value))
    rows.value = res.data || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载审计日志失败')
  } finally {
    loading.value = false
  }
}
</script>