<template>
  <el-dialog
    :model-value="visible"
    title="手动触发发布"
    width="640px"
    @close="$emit('close')"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="目标环境" prop="env_id">
        <el-select v-model="form.env_id" placeholder="请选择环境" style="width: 100%">
          <el-option
            v-for="env in environments"
            :key="env.id"
            :label="`${env.env_name} (${env.namespace})`"
            :value="env.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="版本号" prop="version">
        <el-input v-model="form.version" placeholder="例如：v1.0.1" />
      </el-form-item>

      <el-form-item label="分支" prop="branch">
        <el-input v-model="form.branch" placeholder="默认 main" />
      </el-form-item>

      <el-form-item label="Git Commit">
        <el-input v-model="form.git_commit" placeholder="例如：abc123def" />
      </el-form-item>

      <el-form-item label="镜像 Tag">
        <el-input v-model="form.image_tag" placeholder="例如：order-service:v1.0.1" />
      </el-form-item>

      <el-form-item label="操作人" prop="operator">
        <el-input v-model="form.operator" placeholder="例如：nofish" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('close')">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        触发发布
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createPipelineRun } from '@/api/pipeline'
import type { EnvironmentItem } from '@/types'

const props = defineProps<{
  visible: boolean
  appId: number
  environments: EnvironmentItem[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  env_id: undefined as number | undefined,
  version: '',
  branch: 'main',
  git_commit: '',
  image_tag: '',
  operator: '',
})

const rules: FormRules = {
  env_id: [{ required: true, message: '请选择目标环境', trigger: 'change' }],
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  operator: [{ required: true, message: '请输入操作人', trigger: 'blur' }],
}

watch(
  () => props.visible,
  (val) => {
    if (!val) return
    form.env_id = props.environments[0]?.id
    form.version = ''
    form.branch = 'main'
    form.git_commit = ''
    form.image_tag = ''
    form.operator = ''
  },
  { immediate: true },
)

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true

    await createPipelineRun({
      app_id: props.appId,
      env_id: Number(form.env_id),
      version: form.version,
      branch: form.branch,
      git_commit: form.git_commit,
      image_tag: form.image_tag,
      operator: form.operator,
      trigger_type: 'manual',
    })

    ElMessage.success('发布任务已创建')
    emit('success')
    emit('close')
  } catch (error: any) {
    ElMessage.error(error.message || '触发发布失败')
  } finally {
    submitting.value = false
  }
}
</script>