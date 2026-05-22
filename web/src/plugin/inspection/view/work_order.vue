<template>
  <div class="list-page">
    <div class="sticky-header">
      <el-tabs v-model="activeTab" @tab-change="onTabChange">
        <el-tab-pane label="待接收" name="pending" />
        <el-tab-pane label="检测中" name="inspecting" />
        <el-tab-pane label="待确认" name="confirming" />
        <el-tab-pane label="已完成" name="completed" />
      </el-tabs>
      <div class="search-bar">
        <el-form :model="searchInfo" inline>
          <el-form-item label="MO号">
            <el-input
              v-model="searchInfo.moNumber"
              placeholder="请输入"
              clearable
              size="small"
            />
          </el-form-item>
          <el-form-item label="型号">
            <el-input
              v-model="searchInfo.model"
              placeholder="请输入"
              clearable
              size="small"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="small" @click="getList"
              >查询</el-button
            >
            <el-button size="small" @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <div class="scroll-content">
      <el-table
        :data="tableData"
        border
        stripe
        v-loading="loading"
        size="small"
      >
        <el-table-column prop="moNumber" label="MO号" min-width="140" />
        <el-table-column prop="batchNumber" label="批次号" min-width="160" />
        <el-table-column prop="productName" label="产品名称" min-width="150" />
        <el-table-column prop="model" label="型号" width="100" />
        <el-table-column
          prop="firmwareVersion"
          label="固件版本"
          min-width="130"
        />
        <el-table-column label="业务类型" width="90">
          <template #default="s1">{{
            catLabel(s1.row.instrumentCategory)
          }}</template>
        </el-table-column>
        <el-table-column label="模板" min-width="120">
          <template #default="s2">{{ s2.row.template?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="总数" width="70">
          <template #default="s3">
            <DeviceStatusCount
              :row="s3.row"
              type="all"
              :count="s3.row.deviceCount"
              :batch-id="s3.row.ID"
              @changed="getList"
            />
          </template>
        </el-table-column>
        <el-table-column label="合格数" width="90">
          <template #default="s3">
            <DeviceStatusCount
              :row="s3.row"
              type="pass"
              :count="s3.row.passCount"
              :batch-id="s3.row.ID"
              @changed="getList"
            />
          </template>
        </el-table-column>
        <el-table-column label="不合格数" width="100">
          <template #default="s3">
            <DeviceStatusCount
              :row="s3.row"
              type="fail"
              :count="s3.row.failCount"
              :batch-id="s3.row.ID"
              @changed="getList"
            />
          </template>
        </el-table-column>
        <el-table-column label="返工数" width="90">
          <template #default="s3">
            <DeviceStatusCount
              :row="s3.row"
              type="rework"
              :count="s3.row.reworkCount"
              :batch-id="s3.row.ID"
              @changed="getList"
            />
          </template>
        </el-table-column>
        <el-table-column label="待复检" width="90">
          <template #default="s3">
            <DeviceStatusCount
              :row="s3.row"
              type="recheck"
              :count="s3.row.recheckCount"
              :batch-id="s3.row.ID"
              allow-recheck-actions
              @changed="getList"
            />
          </template>
        </el-table-column>
        <el-table-column label="合格率" width="100">
          <template #default="s3">
            {{ passRateLabel(s3.row.passCount, s3.row.deviceCount) }}
          </template>
        </el-table-column>
        <el-table-column label="检测人" width="100">
          <template #default="s4">{{ s4.row.inspectorName || '-' }}</template>
        </el-table-column>
        <el-table-column prop="CreatedAt" label="创建时间" width="160">
          <template #default="s5">{{ formatDate(s5.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="s6">
            <el-button
              v-if="activeTab === 'pending'"
              size="small"
              type="primary"
              @click="onStartInspect(s6.row)"
              >接收并开始检测</el-button
            >
            <el-button
              v-if="activeTab === 'confirming' && canConfirmComplete(s6.row)"
              size="small"
              type="success"
              @click="onConfirmComplete(s6.row)"
            >
              确认完成
            </el-button>
            <el-button
              v-if="activeTab !== 'pending'"
              size="small"
              type="primary"
              @click="openDetail(s6.row)"
            >
              {{ activeTab === 'completed' || activeTab === 'confirming' ? '查看' : '检测' }}
            </el-button>
            <el-button size="small" type="success" link @click="onExportExcel(s6.row)">
              导出Excel
            </el-button>
            <el-button size="small" type="primary" link @click="openPrint(s6.row)">
              打印
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="searchInfo.page"
          v-model:page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="getList"
          @current-change="getList"
          small
        />
      </div>
    </div>
  </div>
</template>

<script setup>
  import { onMounted, onUnmounted, ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import DeviceStatusCount from '@/plugin/inspection/components/DeviceStatusCount.vue'
  import {
    getInspectionBatchList,
    exportInspectionExcel,
    confirmInspectionComplete,
    startInspection
  } from '@/plugin/inspection/api/work_order'

  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const savedTab = sessionStorage.getItem('inspectTab')
  const activeTab = ref(savedTab === 'recheck' ? 'confirming' : (savedTab || 'pending'))
  const searchInfo = reactive({
    moNumber: '',
    model: '',
    page: 1,
    pageSize: 30
  })
  const catLabel = (v) =>
    ({
      online: '线上',
      offline: '线下',
      foreign_trade: '外贸',
      custom: '定制款'
    }[v] || v)
  const passRateLabel = (passCount, deviceCount) => {
    const total = Number(deviceCount || 0)
    if (!total) return '-'
    return `${((Number(passCount || 0) / total) * 100).toFixed(1)}%`
  }
  const canConfirmComplete = (row) => {
    const deviceCount = Number(row.deviceCount || 0)
    if (!deviceCount) return false
    const passCount = Number(row.passCount || 0)
    const pendingCount = Number(row.failCount || 0) +
      Number(row.reworkCount || 0) +
      Number(row.recheckCount || 0)
    return passCount === deviceCount && pendingCount === 0
  }

  const getList = async () => {
    loading.value = true
    const statusMap = { pending: 1, inspecting: 2, confirming: 3, completed: 4 }
    try {
      const res = await getInspectionBatchList({
        ...searchInfo,
        status: statusMap[activeTab.value]
      })
      if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
      }
    } finally {
      loading.value = false
    }
  }
  const resetSearch = () => {
    searchInfo.moNumber = ''
    searchInfo.model = ''
    searchInfo.page = 1
    getList()
  }
  const onTabChange = () => {
    sessionStorage.setItem('inspectTab', activeTab.value)
    searchInfo.page = 1
    getList()
  }
  const refreshOnVisible = () => {
    if (document.visibilityState === 'visible') {
      getList()
    }
  }

  const openDetail = (row) => {
    window.location.hash = `/inspectDetail?batchId=${row.ID}`
  }
  const openPrint = (row) => {
    const url = `${window.location.origin}${window.location.pathname}#/inspectPrint?batchId=${row.ID}`
    window.open(url, '_blank')
  }
  const downloadBlob = (blob, filename) => {
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  }
  const onExportExcel = async (row) => {
    const res = await exportInspectionExcel({ id: row.ID })
    const filename = `${row.moNumber || 'MO'}-${row.batchNumber || row.ID}-检测工单.xlsx`
    downloadBlob(res.data || res, filename)
  }
  const onStartInspect = async (row) => {
    await ElMessageBox.confirm('确定接收该批次并开始检测？', '提示', { type: 'info' })
    const res = await startInspection({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('已接收并开始检测')
      getList()
      openDetail(row)
    }
  }
  const onConfirmComplete = async (row) => {
    await ElMessageBox.confirm('确认该批次全部闭环并完成检测？确认后只能查看、打印和导出。', '确认完成', {
      type: 'success',
      confirmButtonText: '确认完成'
    })
    const res = await confirmInspectionComplete({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('已确认完成')
      getList()
    }
  }
  getList()
  onMounted(() => {
    document.addEventListener('visibilitychange', refreshOnVisible)
    window.addEventListener('focus', getList)
  })
  onUnmounted(() => {
    document.removeEventListener('visibilitychange', refreshOnVisible)
    window.removeEventListener('focus', getList)
  })
</script>

<style scoped>
  .list-page {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--el-bg-color, #fff);
    overflow: hidden;
  }
  .sticky-header {
    flex-shrink: 0;
    padding: 8px 16px 0 16px;
    background: var(--el-bg-color, #fff);
    z-index: 10;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
  .sticky-header :deep(.el-tabs__header) {
    margin-bottom: 0;
  }
  .search-bar {
    padding: 4px 0 8px;
    background: var(--el-bg-color, #fff);
  }
  .scroll-content {
    flex: 1;
    overflow-y: auto;
    padding: 0 16px 16px;
  }
  .pagination-wrap {
    display: flex;
    justify-content: flex-end;
    flex-wrap: wrap;
    padding: 12px 0 8px;
  }
  .pagination-wrap :deep(.el-pagination) {
    flex-wrap: wrap;
    justify-content: flex-end;
  }
</style>
