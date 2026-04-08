<template>
  <div class="page-header">
    <div>
      <div class="page-title">{{ title }}</div>
      <div v-if="description" class="page-description">{{ description }}</div>
    </div>
    <div class="page-extra">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>