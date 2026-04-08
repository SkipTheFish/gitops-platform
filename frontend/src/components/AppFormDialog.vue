<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? '编辑应用' : '新建应用'"
    width="720px"
    @close="$emit('close')"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="应用名称" prop="name">
        <el-input v-model="form.name" placeholder="例如：order-service" />
      </el-form-item>

      <el-form-item label="代码仓库" prop="repo_url">
        <el-input v-model="form.repo_url" placeholder="GitHub 仓库地址" />
      </el-form-item>

      <el-form-item label="配置仓库" prop="config_repo_url">
        <el-input v-model="form.config_repo_url" placeholder="Helm / 部署配置仓库地址" />
      </el-form-item>

      <el-form-item label="集群名称" prop="cluster_name">
        <el-input v-model="form.cluster_name" placeholder="例如：kind-gitops-platform" />
      </el-form-item>

      <el-form-item label="命名空间" prop="namespace">
        <el-input v-model="form.namespace" placeholder="例如：order" />
      </el-form-item>

      <el-form-item label="Chart 路径" prop="helm_chart_path">
        <el-input v-model="form.helm_chart_path" placeholder="例如：charts/order-service" />
      </el-form-item>

      <el-form-item label="Values 路径" prop="values_file_path">
        <el-input v-model="form.values_file_path" placeholder="例如：values-dev.yaml" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('close')">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        {{ isEdit ? '保存修改' : '创建应用' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createApp, updateApp } from '@/api/app'
import type { AppItem } from '@/types'

const props = defineProps<{
  visible: boolean
  app?: AppItem | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  name: '',
  repo_url: '',
  config_repo_url: '',
  cluster_name: '',
  namespace: '',
  helm_chart_path: '',
  values_file_path: '',
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  repo_url: [{ required: true, message: '请输入代码仓库地址', trigger: 'blur' }],
  config_repo_url: [{ required: true, message: '请输入配置仓库地址', trigger: 'blur' }],
  cluster_name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
  namespace: [{ required: true, message: '请输入命名空间', trigger: 'blur' }],
  helm_chart_path: [{ required: true, message: '请输入 Chart 路径', trigger: 'blur' }],
  values_file_path: [{ required: true, message: '请输入 Values 路径', trigger: 'blur' }],
}

const isEdit = ref(false)

watch(
  () => props.visible,
  (val) => {
    if (!val) return

    if (props.app) {
      isEdit.value = true
      form.name = props.app.name
      form.repo_url = props.app.repo_url
      form.config_repo_url = props.app.config_repo_url
      form.cluster_name = props.app.cluster_name
      form.namespace = props.app.namespace
      form.helm_chart_path = props.app.helm_chart_path
      form.values_file_path = props.app.values_file_path
    } else {
      isEdit.value = false
      form.name = ''
      form.repo_url = ''
      form.config_repo_url = ''
      form.cluster_name = ''
      form.namespace = ''
      form.helm_chart_path = ''
      form.values_file_path = ''
    }
  },
  { immediate: true },
)

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true

    if (isEdit.value && props.app) {
      await updateApp(props.app.id, { ...form })
      ElMessage.success('应用更新成功')
    } else {
      await createApp({ ...form })
      ElMessage.success('应用创建成功')
    }

    emit('success')
    emit('close')
  } catch (error: any) {
    if (error?.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}
</script>