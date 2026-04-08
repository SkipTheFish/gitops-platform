<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="logo">GitOps Platform</div>

      <el-menu
        :default-active="activeMenu"
        class="menu"
        router
      >
        <el-menu-item index="/apps">
          <el-icon><Grid /></el-icon>
          <span>应用管理</span>
        </el-menu-item>

        <el-menu-item index="/deployments">
          <el-icon><Histogram /></el-icon>
          <span>部署记录</span>
        </el-menu-item>

        <el-menu-item index="/audits">
          <el-icon><Document /></el-icon>
          <span>审计日志</span>
        </el-menu-item>
      </el-menu>
    </aside>

    <div class="main">
      <header class="topbar">
        <div class="topbar-left">
          <div class="title">Kubernetes + GitOps 发布平台</div>
          <div class="subtitle">企业风格的最小可用管理后台</div>
        </div>
      </header>

      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const activeMenu = computed(() => {
  if (route.path.startsWith('/apps')) return '/apps'
  if (route.path.startsWith('/deployments')) return '/deployments'
  if (route.path.startsWith('/audits')) return '/audits'
  return '/apps'
})
</script>