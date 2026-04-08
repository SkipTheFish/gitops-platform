import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import AppListView from '@/views/AppListView.vue'
import AppDetailView from '@/views/AppDetailView.vue'
import DeploymentListView from '@/views/DeploymentListView.vue'
import AuditListView from '@/views/AuditListView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: MainLayout,
      redirect: '/apps',
      children: [
        {
          path: '/apps',
          name: 'apps',
          component: AppListView,
        },
        {
          path: '/apps/:id',
          name: 'app-detail',
          component: AppDetailView,
        },
        {
          path: '/deployments',
          name: 'deployments',
          component: DeploymentListView,
        },
        {
          path: '/audits',
          name: 'audits',
          component: AuditListView,
        },
      ],
    },
  ],
})

export default router