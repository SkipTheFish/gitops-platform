<template>
  <div>
    <PageHeader
      title="应用管理"
      description="管理应用元信息、代码仓库、配置仓库和 Helm 部署路径"
    >
      <el-button type="primary" @click="openCreateDialog">
        新建应用
      </el-button>
    </PageHeader>

    <el-card shadow="hover" class="card">
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="按应用名称搜索"
          clearable
          style="width: 280px"
        />
      </div>

      <el-table
        :data="filteredList"
        v-loading="loading"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="应用名称" min-width="180" />
        <el-table-column prop="cluster_name" label="集群" min-width="180" />
        <el-table-column prop="namespace" label="命名空间" min-width="140" />
        <el-table-column prop="helm_chart_path" label="Chart 路径" min-width="220" />
        <el-table-column prop="values_file_path" label="Values 路径" min-width="220" />
        <el-table-column label="创建时间" min-width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goDetail(row.id)">详情</el-button>
            <el-button link type="warning" @click="openEditDialog(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <AppFormDialog
      :visible="dialogVisible"
      :app="currentApp"
      @close="dialogVisible = false"
      @success="loadData"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import AppFormDialog from '@/components/AppFormDialog.vue'
import { getAppList } from '@/api/app'
import type { AppItem } from '@/types'

const router = useRouter()

const loading = ref(false)
const dialogVisible = ref(false)
const currentApp = ref<AppItem | null>(null)
const list = ref<AppItem[]>([])
const keyword = ref('')

// 计算属性：根据关键字做前端过滤
const filteredList = computed(() => {
  const key = keyword.value.trim().toLowerCase()
  if (!key) return list.value
  return list.value.filter((item) => item.name.toLowerCase().includes(key))
})

function formatTime(value: string) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

async function loadData() {
  try {
    loading.value = true
    const res = await getAppList()
    list.value = res.data || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载应用列表失败')
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  currentApp.value = null
  dialogVisible.value = true
}

function openEditDialog(app: AppItem) {
  currentApp.value = app
  dialogVisible.value = true
}

function goDetail(id: number) {
  router.push(`/apps/${id}`)
}

onMounted(() => {
  loadData()
})
</script>