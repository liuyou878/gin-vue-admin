import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/view/init/index.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/view/login/index.vue')
  },
  {
    path: '/scanUpload',
    name: 'ScanUpload',
    meta: {
      title: '扫码上传',
      client: true
    },
    component: () => import('@/view/example/upload/scanUpload.vue')
  },
  {
    path: '/inspectDetail',
    name: 'InspectDetail',
    meta: { title: '检测详情', client: true },
    component: () => import('@/plugin/inspection/view/inspect_detail.vue')
  },
  {
    path: '/inspectPrint',
    name: 'InspectPrint',
    meta: { title: '检测工单打印', client: true },
    component: () => import('@/plugin/inspection/view/inspect_print.vue')
  },
  {
    path: '/inspectWorkOrder',
    name: 'InspectWorkOrder',
    meta: {
      title: '检测工单',
      client: true
    },
    component: () => import('@/plugin/inspection/view/work_order.vue')
  },
  {
    path: '/productionMockSubmit',
    name: 'ProductionMockSubmit',
    meta: {
      title: '模拟生产工具提交',
      client: true
    },
    component: () => import('@/plugin/inspection/view/production_mock_submit.vue')
  },
  {
    path: '/batchImportDevice',
    name: 'BatchImportDevice',
    meta: {
      title: '批量导入设备',
      client: true
    },
    component: () => import('@/plugin/inspection/view/batch_import.vue')
  },
  {
    path: '/publicFirmwareDownload',
    name: 'PublicFirmwareDownload',
    meta: {
      title: '版本下载',
      client: true
    },
    component: () => import('@/view/publicFirmwareDownload/index.vue')
  },
  {
    path: '/:catchAll(.*)',
    meta: {
      closeTab: true
    },
    component: () => import('@/view/error/index.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router

