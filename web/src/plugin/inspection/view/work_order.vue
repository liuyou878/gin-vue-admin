<template>
  <div class="list-page">
    <div class="sticky-header">
      <el-tabs v-model="activeTab" @tab-change="onTabChange">
        <el-tab-pane label="待接收" name="pending" />
        <el-tab-pane label="检测中" name="inspecting" />
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
          <el-form-item label="批次号">
            <el-input
              v-model="searchInfo.batchNumber"
              placeholder="请输入"
              clearable
              size="small"
            />
          </el-form-item>
          <el-form-item label="SN">
            <el-input
              v-model="searchInfo.sn"
              placeholder="请输入"
              clearable
              size="small"
            />
          </el-form-item>
          <!-- <el-form-item label="设备状态">
            <el-select
              v-model="searchInfo.deviceStatus"
              placeholder="请选择"
              clearable
              size="small"
              style="width: 140px"
            >
              <el-option
                v-for="item in deviceStatusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item> -->
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
      <div class="desktop-table">
        <el-table
          :data="tableData"
          border
          stripe
          v-loading="loading"
          size="small"
        >
          <el-table-column prop="moNumber" label="MO号" min-width="140" />
          <el-table-column prop="batchNumber" label="批次号" min-width="160" />
          <el-table-column
            prop="productName"
            label="产品名称"
            min-width="150"
          />
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
            <template #default="s2">{{
              s2.row.template?.name || '-'
            }}</template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="s3">
              <el-tag size="small" :type="activeTabTagType">
                {{ activeTabLabel }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="总数" width="80" align="center">
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
          <el-table-column label="合格" width="80" align="center">
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
          <el-table-column label="待测" width="80" align="center">
            <template #default="s3">
              <DeviceStatusCount
                :row="s3.row"
                type="pending"
                :count="pendingCount(s3.row)"
                :batch-id="s3.row.ID"
                @changed="getList"
              />
            </template>
          </el-table-column>
          <el-table-column label="异常" width="80" align="center">
            <template #default="s3">
              <DeviceStatusCount
                :row="s3.row"
                type="abnormal"
                :count="abnormalCount(s3.row)"
                :batch-id="s3.row.ID"
                allow-recheck-actions
                @changed="getList"
              />
            </template>
          </el-table-column>
          <el-table-column label="合格率" width="90" align="center">
            <template #default="s3">
              {{ passRateLabel(s3.row.passCount, s3.row.deviceCount) }}
            </template>
          </el-table-column>
          <el-table-column label="检测人" width="100">
            <template #default="s4">{{ s4.row.inspectorName || '-' }}</template>
          </el-table-column>
          <el-table-column label="更新时间" width="160">
            <template #default="s5">{{
              formatDate(s5.row.UpdatedAt || s5.row.CreatedAt)
            }}</template>
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
                v-if="activeTab !== 'pending'"
                size="small"
                type="primary"
                @click="openDetail(s6.row)"
              >
                {{
                  activeTab === 'completed'
                    ? '查看'
                    : '检测'
                }}
              </el-button>
              <el-button
                size="small"
                type="success"
                link
                @click="onExportExcel(s6.row)"
              >
                导出Excel
              </el-button>
              <el-button
                size="small"
                type="primary"
                link
                @click="openPrint(s6.row)"
              >
                打印
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div
        v-loading="loading && !tableData.length && !mobileRefreshing"
        class="mobile-cards"
        :class="{ 'is-refreshing': mobileRefreshing }"
      >
        <div v-for="row in tableData" :key="row.ID" class="work-card">
          <div class="card-head">
            <div>
              <div class="card-title">{{ row.moNumber || '-' }}</div>
              <div class="card-subtitle">{{ row.batchNumber || '-' }}</div>
            </div>
            <el-tag
              size="small"
              :type="
                activeTab === 'completed'
                  ? 'success'
                  : activeTab === 'inspecting'
                  ? 'primary'
                  : 'info'
              "
            >
              {{ activeTabLabel }}
            </el-tag>
          </div>

          <div class="card-meta">
            <span>型号：{{ row.model || '-' }}</span>
            <span>模板：{{ row.template?.name || '-' }}</span>
            <span>业务：{{ catLabel(row.instrumentCategory) || '-' }}</span>
            <span>检测人：{{ row.inspectorName || '-' }}</span>
          </div>

          <div class="card-counts">
            <div class="count-box">
              <span>总数</span>
              <DeviceStatusCount
                :row="row"
                type="all"
                :count="row.deviceCount"
                :batch-id="row.ID"
                @changed="getList"
              />
            </div>
            <div class="count-box">
              <span>合格</span>
              <DeviceStatusCount
                :row="row"
                type="pass"
                :count="row.passCount"
                :batch-id="row.ID"
                @changed="getList"
              />
            </div>
            <div class="count-box">
              <span>待测</span>
              <DeviceStatusCount
                :row="row"
                type="pending"
                :count="pendingCount(row)"
                :batch-id="row.ID"
                @changed="getList"
              />
            </div>
            <div class="count-box">
              <span>异常</span>
              <DeviceStatusCount
                :row="row"
                type="abnormal"
                :count="abnormalCount(row)"
                :batch-id="row.ID"
                allow-recheck-actions
                @changed="getList"
              />
            </div>
            <div class="count-box rate">
              <span>合格率</span>
              <strong>{{
                passRateLabel(row.passCount, row.deviceCount)
              }}</strong>
            </div>
          </div>

          <div class="card-actions">
            <el-button
              v-if="activeTab === 'pending'"
              type="primary"
              size="small"
              @click="onStartInspect(row)"
            >
              接收并开始检测
            </el-button>
            <el-button
              v-if="activeTab !== 'pending'"
              type="primary"
              size="small"
              @click="openDetail(row)"
            >
              {{
                activeTab === 'completed'
                  ? '查看'
                  : '检测'
              }}
            </el-button>
          </div>
        </div>
        <div
          v-if="!loading && tableData.length === 0"
          class="mobile-empty-state"
        >
          <el-empty description="暂无检测工单" />
        </div>
      </div>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="searchInfo.page"
          v-model:page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :pager-count="isMobile ? 5 : 7"
          :total="total"
          :layout="paginationLayout"
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
  import { computed, onMounted, onUnmounted, ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import DeviceStatusCount from '@/plugin/inspection/components/DeviceStatusCount.vue'
  import {
    getInspectionBatchList,
    exportInspectionExcel,
    startInspection
  } from '@/plugin/inspection/api/work_order'

  const loading = ref(false)
  const mobileRefreshing = ref(false)
  const isMobile = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const requestSeq = ref(0)
  const savedTab = sessionStorage.getItem('inspectTab')
  const activeTab = ref(
    savedTab === 'recheck' || savedTab === 'confirming'
      ? 'inspecting'
      : savedTab || 'pending'
  )
  const searchInfo = reactive({
    moNumber: '',
    model: '',
    batchNumber: '',
    sn: '',
    deviceStatus: '',
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
  const activeTabLabelMap = {
    pending: '待接收',
    inspecting: '检测中',
    completed: '已完成'
  }
  const deviceStatusOptions = [
    { label: '待检测设备', value: 'pending' },
    { label: '合格', value: 'pass' },
    { label: '不合格', value: 'fail' },
    { label: '待生产接收', value: 'returned' },
    { label: '返工中', value: 'rework' },
    { label: '待复检', value: 'pending_recheck' },
    { label: '复检中', value: 'rechecking' }
  ]
  const activeTabLabel = computed(
    () => activeTabLabelMap[activeTab.value] || '-'
  )
  const activeTabTagType = computed(
    () =>
      ({
        pending: 'info',
        inspecting: 'primary',
        completed: 'success'
      }[activeTab.value] || 'info')
  )
  const paginationLayout = computed(() =>
    isMobile.value
      ? 'prev, pager, next'
      : 'total, sizes, prev, pager, next, jumper'
  )
  const passRateLabel = (passCount, deviceCount) => {
    const total = Number(deviceCount || 0)
    if (!total) return '-'
    return `${((Number(passCount || 0) / total) * 100).toFixed(1)}%`
  }
  const abnormalCount = (row) =>
    Number(row.abnormalCount ?? 0) ||
    Number(row.failCount || 0) +
      Number(row.reworkCount || 0) +
      Number(row.recheckCount || 0)
  const pendingCount = (row) => {
    const total = Number(row.deviceCount || 0)
    const count = total - Number(row.passCount || 0) - abnormalCount(row)
    return count > 0 ? count : 0
  }
  const getList = async (options = {}) => {
    const smooth = options.smooth === true && isMobile.value
    const seq = requestSeq.value + 1
    requestSeq.value = seq
    if (smooth) {
      mobileRefreshing.value = true
    } else {
      loading.value = true
    }
    const statusMap = { pending: 1, inspecting: 2, completed: 4 }
    try {
      const res = await getInspectionBatchList({
        ...searchInfo,
        status: statusMap[activeTab.value]
      })
      if (seq !== requestSeq.value) return
      if (res.code === 0) {
        tableData.value = Array.isArray(res.data?.list) ? res.data.list : []
        total.value = Number(res.data?.total || 0)
      }
    } finally {
      if (seq === requestSeq.value) {
        loading.value = false
        mobileRefreshing.value = false
      }
    }
  }
  const resetSearch = () => {
    searchInfo.moNumber = ''
    searchInfo.model = ''
    searchInfo.batchNumber = ''
    searchInfo.sn = ''
    searchInfo.deviceStatus = ''
    searchInfo.page = 1
    getList()
  }
  const onTabChange = () => {
    sessionStorage.setItem('inspectTab', activeTab.value)
    searchInfo.page = 1
    getList({ smooth: true })
  }
  const refreshOnVisible = () => {
    if (document.visibilityState === 'visible') {
      getList()
    }
  }
  const refreshList = () => {
    getList()
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
    const filename = `${row.moNumber || 'MO'}-${
      row.batchNumber || row.ID
    }-检测工单.xlsx`
    downloadBlob(res.data || res, filename)
  }
  const onStartInspect = async (row) => {
    const countText = `本批次共 ${Number(row.deviceCount || 0)} 台设备`
    await ElMessageBox.confirm(
      `确定接收该批次并开始检测？\n\n批次号：${row.batchNumber || '-'}\n${countText}`,
      '提示',
      {
      type: 'info'
      }
    )
    const res = await startInspection({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success({ message: '已接收并开始检测', duration: 1000 })
      getList()
      window.setTimeout(() => {
        ElMessage.closeAll()
        openDetail(row)
      }, 300)
    }
  }
  const updateIsMobile = () => {
    isMobile.value = window.matchMedia('(max-width: 768px)').matches
  }
  updateIsMobile()
  getList()
  onMounted(() => {
    updateIsMobile()
    document.addEventListener('visibilitychange', refreshOnVisible)
    window.addEventListener('focus', refreshList)
    window.addEventListener('resize', updateIsMobile)
  })
  onUnmounted(() => {
    document.removeEventListener('visibilitychange', refreshOnVisible)
    window.removeEventListener('focus', refreshList)
    window.removeEventListener('resize', updateIsMobile)
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
  .mobile-cards {
    display: none;
  }
  .mobile-cards.is-refreshing {
    opacity: 0.72;
    transition: opacity 0.16s ease;
    pointer-events: none;
  }
  .work-card {
    padding: 14px;
    border: 1px solid var(--el-border-color-light, #e4e7ed);
    border-radius: 14px;
    background: var(--el-bg-color, #fff);
    box-shadow: 0 6px 18px rgba(20, 36, 64, 0.08);
  }
  .work-card + .work-card {
    margin-top: 12px;
  }
  .mobile-empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 260px;
    border: 1px dashed var(--el-border-color-light, #e4e7ed);
    border-radius: 14px;
    background: var(--el-bg-color, #fff);
  }
  .card-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 10px;
  }
  .card-title {
    color: var(--el-text-color-primary, #303133);
    font-size: 17px;
    font-weight: 800;
    line-height: 1.25;
  }
  .card-subtitle {
    margin-top: 4px;
    color: var(--el-text-color-secondary, #909399);
    font-size: 13px;
  }
  .card-meta {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 6px 10px;
    margin-bottom: 12px;
    color: var(--el-text-color-regular, #606266);
    font-size: 12px;
  }
  .card-meta span {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .card-counts {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 8px;
    margin-bottom: 12px;
  }
  .count-box {
    min-width: 0;
    padding: 8px;
    border-radius: 10px;
    background: var(--el-fill-color-lighter, #fafafa);
    text-align: center;
  }
  .count-box span {
    display: block;
    margin-bottom: 2px;
    color: var(--el-text-color-secondary, #909399);
    font-size: 12px;
  }
  .count-box.rate strong {
    color: var(--el-text-color-primary, #303133);
    font-size: 14px;
  }
  .card-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .card-actions :deep(.el-button) {
    margin-left: 0;
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
  @media (max-width: 768px) {
    .sticky-header {
      padding: 6px 10px 0;
    }
    .search-bar :deep(.el-form) {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 6px;
    }
    .search-bar :deep(.el-form-item) {
      margin-right: 0;
      margin-bottom: 0;
    }
    .search-bar :deep(.el-form-item:last-child) {
      grid-column: 1 / -1;
    }
    .search-bar :deep(.el-form-item__content) {
      width: 100%;
    }
    .search-bar :deep(.el-input) {
      width: 100%;
    }
    .scroll-content {
      padding: 10px;
    }
    .desktop-table {
      display: none;
    }
    .mobile-cards {
      display: block;
      min-height: 320px;
    }
    .mobile-empty-state {
      min-height: calc(100vh - 260px);
    }
    .card-counts {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
    .pagination-wrap {
      justify-content: center;
      min-height: 44px;
      padding: 10px 0 calc(10px + env(safe-area-inset-bottom));
    }
    .pagination-wrap :deep(.el-pagination) {
      width: 100%;
      justify-content: center;
      flex-wrap: nowrap;
    }
    .pagination-wrap :deep(.el-pager) {
      flex-wrap: nowrap;
    }
    .pagination-wrap :deep(.btn-prev),
    .pagination-wrap :deep(.btn-next),
    .pagination-wrap :deep(.el-pager li) {
      min-width: 30px;
    }
  }
</style>
