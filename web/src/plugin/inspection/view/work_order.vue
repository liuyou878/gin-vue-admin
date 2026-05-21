<template>
  <div class="list-page">
    <div class="sticky-header">
      <el-tabs v-model="activeTab" @tab-change="onTabChange">
        <el-tab-pane label="待检测" name="pending" />
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
        <el-table-column prop="deviceCount" label="总数" width="60" />
        <el-table-column label="合格数" width="90">
          <template #default="s3">
            <span class="count-pass">{{ s3.row.passCount || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="不合格数" width="100">
          <template #default="s3">
            <span class="count-fail">{{ s3.row.failCount || 0 }}</span>
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
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="s6">
            <el-button
              v-if="activeTab === 'pending'"
              size="small"
              type="primary"
              @click="onStartInspect(s6.row)"
              >开始检测</el-button
            >
            <el-button
              v-else
              size="small"
              type="primary"
              @click="openDetail(s6.row)"
            >
              {{ activeTab === 'completed' ? '查看' : '检测' }}
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
  import { ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import {
    getInspectionBatchList,
    startInspection
  } from '@/plugin/inspection/api/work_order'

  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const activeTab = ref(sessionStorage.getItem('inspectTab') || 'pending')
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

  const getList = async () => {
    loading.value = true
    const statusMap = { pending: 1, inspecting: 2, completed: 3 }
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

  const openDetail = (row) => {
    window.location.hash = `/inspectDetail?batchId=${row.ID}`
  }
  const onStartInspect = async (row) => {
    await ElMessageBox.confirm('确定开始检测？', '提示', { type: 'info' })
    const res = await startInspection({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('已开始检测')
      getList()
      openDetail(row)
    }
  }
  getList()
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
  .count-pass {
    color: #16a34a;
    font-weight: 600;
  }
  .count-fail {
    color: #dc2626;
    font-weight: 600;
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
