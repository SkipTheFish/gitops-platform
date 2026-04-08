<template>
  <el-dialog
    :model-value="visible"
    :title="environment ? '编辑环境' : '新增环境'"
    width="620px"
    @close="$emit('close')"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="环境名称" prop="env_name">
        <el-input v-model="form.env_name" placeholder="例如：dev / test / prod" />
      </el-form-item>

      <el-form-item label="集群名称" prop="cluster_name">
        <el-input v-model="form.cluster_name" placeholder="例如：kind-gitops-platform" />
      </el-form-item>

      <el-form-item label="命名空间" prop="namespace">
        <el-input v-model="form.namespace" placeholder="例如：order-dev" />
      </el-form-item>

      <el-form-item label="自动同步">
        <el-switch v-model="form.auto_sync_enabled" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('close')">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        保存
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createEnvironment, updateEnvironment } from '@/api/environment'
import type { EnvironmentItem } from '@/types'

const props = defineProps<{
  visible: boolean
  appId: number
  environment?: EnvironmentItem | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  env_name: '',
  cluster_name: '',
  namespace: '',
  auto_sync_enabled: false,
})

const rules: FormRules = {
  env_name: [{ required: true, message: '请输入环境名称', trigger: 'blur' }],
  cluster_name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
  namespace: [{ required: true, message: '请输入命名空间', trigger: 'blur' }],
}

watch(
  () => props.visible,
  (val) => {
    if (!val) return

    if (props.environment) {
      form.env_name = props.environment.env_name
      form.cluster_name = props.environment.cluster_name
      form.namespace = props.environment.namespace
      form.auto_sync_enabled = props.environment.auto_sync_enabled
    } else {
      form.env_name = ''
      form.cluster_name = ''
      form.namespace = ''
      form.auto_sync_enabled = false
    }
  },
  { immediate: true },
)

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true

    if (props.environment) {
      await updateEnvironment(props.environment.id, { ...form })
      ElMessage.success('环境更新成功')
    } else {
      await createEnvironment(props.appId, { ...form })
      ElMessage.success('环境创建成功')
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