<template>
  <div class="inspect-page">
    <!-- Header -->
    <div class="page-header">
      <h2 class="m-0 text-lg">检测工单</h2>
      <span class="text-sm text-gray ml-auto" v-if="!showDetail">{{ tableData.length }} 条记录</span>
    </div>

    <!-- List View -->
    <div v-if="!showDetail" class="page-body">
      <el-tabs v-model="activeTab" @tab-change="onTabChange" class="tab-section">
        <el-tab-pane label="待检测" name="pending" />
        <el-tab-pane label="检测中" name="inspecting" />
        <el-tab-pane label="已完成" name="completed" />
      </el-tabs>

      <div class="search-bar">
        <el-form :model="searchInfo" inline>
          <el-form-item label="MO号">
            <el-input v-model="searchInfo.moNumber" placeholder="请输入" clearable size="small" />
          </el-form-item>
          <el-form-item label="型号">
            <el-input v-model="searchInfo.model" placeholder="请输入" clearable size="small" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="small" @click="getList">查询</el-button>
            <el-button size="small" @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="table-wrap">
        <el-table :data="tableData" border stripe v-loading="loading" size="small">
          <el-table-column prop="moNumber" label="MO号" min-width="140" />
          <el-table-column prop="batchNumber" label="批次号" min-width="160" />
          <el-table-column prop="model" label="型号" width="90" />
          <el-table-column label="类别" width="80">
            <template #default="a">{{ catLabel(a.row.instrumentCategory) }}</template>
          </el-table-column>
          <el-table-column label="模板" min-width="120">
            <template #default="b">{{ b.row.template?.name || '-' }}</template>
          </el-table-column>
          <el-table-column prop="deviceCount" label="总数" width="60" />
          <el-table-column label="合格/不合格" width="100">
            <template #default="f">{{ f.row.passCount || 0 }} / {{ f.row.failCount || 0 }}</template>
          </el-table-column>
          <el-table-column label="检测人" width="100">
            <template #default="c">{{ c.row.inspectorName || '-' }}</template>
          </el-table-column>
          <el-table-column prop="CreatedAt" label="创建时间" width="160">
            <template #default="d">{{ formatDate(d.row.CreatedAt) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="e">
              <el-button v-if="activeTab === 'pending'" size="small" type="primary" @click="onStartInspect(e.row)">开始检测</el-button>
              <el-button v-else size="small" type="primary" @click="openDetail(e.row)">
                {{ activeTab === 'completed' ? '查看' : '检测' }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="searchInfo.page"
          v-model:page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total"
          layout="total, sizes, prev, pager, next, jumper" background
          @size-change="getList" @current-change="getList" small
        />
      </div>
    </div>

    <!-- Detail View -->
    <div v-if="showDetail" class="page-body">
      <div class="detail-toolbar">
        <el-button @click="showDetail = false" :icon="'ArrowLeft'" />
        <span class="ml-4 text-sm" v-if="detailLoaded">
          <b>MO号:</b> {{ detail.order.moNumber }} &nbsp;
          <b>批次号:</b> {{ detail.order.batchNumber }} &nbsp;
          <b>型号:</b> {{ detail.order.model }} &nbsp;
          <b>检测人:</b> {{ detail.order.inspectorName || '-' }}
        </span>
      </div>

      <div v-if="detailLoaded">
        <el-tabs v-model="inspectMode" class="tab-section">
          <el-tab-pane label="逐台检测" name="byDevice" />
          <el-tab-pane label="逐项检测" name="byItem" />
        </el-tabs>

        <!-- By Device: Desktop table -->
        <div v-if="inspectMode === 'byDevice'" class="desktop-only" style="overflow-x:auto">
          <table class="inspect-table">
            <thead>
              <tr>
                <th class="fixed-col">序号</th>
                <th class="fixed-col sn-col">SN</th>
                <th class="fixed-col">判定</th>
                <th v-for="ti in detail.templateItems" :key="ti.itemID">{{ ti.itemName }}<br /><small>{{ ti.unit || '' }}</small></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="dev in detail.devices" :key="dev.ID" :class="deviceRowClass(dev)">
                <td class="fixed-col">{{ dev.lineNumber }}</td>
                <td class="fixed-col sn-col">{{ dev.sn }}</td>
                <td class="fixed-col">
                  <el-tag :type="deviceStatusTag(dev)" size="small" @click="toggleDeviceStatus(dev)" class="cursor-pointer">
                    {{ deviceStatusLabel(dev) }}
                  </el-tag>
                </td>
                <td v-for="(ti, ri) in detail.templateItems" :key="ti.itemID">
                  <template v-if="ti.resultType === 'pass_fail'">
                    <span class="pass-toggle">
                      <button :class="{active: dev.results[ri]._checked === true}" @click="setPass(dev, ri, true)">✓</button>
                      <button :class="{active: dev.results[ri]._checked === false}" @click="setPass(dev, ri, false)">✗</button>
                    </span>
                  </template>
                  <template v-else-if="ti.resultType === 'number'">
                    <el-input-number v-model="dev.results[ri]._numVal" :precision="2" size="small"
                      controls-position="right" style="width:100px" :class="getRangeClass(dev.results[ri])"
                      @change="onNumChange(dev, ri)" />
                  </template>
                  <template v-else>
                    <div class="flex gap-1 items-center">
                      <span class="pass-toggle">
                        <button :class="{active: dev.results[ri]._checked === true}" @click="setPass(dev, ri, true)">✓</button>
                        <button :class="{active: dev.results[ri]._checked === false}" @click="setPass(dev, ri, false)">✗</button>
                      </span>
                      <el-input-number v-model="dev.results[ri]._numVal" :precision="2" size="small"
                        controls-position="right" style="width:90px" :class="getRangeClass(dev.results[ri])"
                        @change="onNumChange(dev, ri)" />
                    </div>
                  </template>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- By Device: Mobile -->
        <div v-if="inspectMode === 'byDevice'" class="mobile-only">
          <div class="swiper-controls">
            <el-button size="small" :disabled="currentDeviceIndex === 0" @click="currentDeviceIndex--">◀ 上一台</el-button>
            <span class="text-lg font-bold">{{ currentDeviceIndex + 1 }} / {{ detail.devices.length }}</span>
            <el-button size="small" :disabled="currentDeviceIndex >= detail.devices.length - 1" @click="currentDeviceIndex++">下一台 ▶</el-button>
          </div>
          <div class="device-card" v-if="currentDevice">
            <div class="card-header">
              <span class="text-lg font-bold">{{ currentDevice.sn }}</span>
              <el-tag :type="deviceStatusTag(currentDevice)" size="large" @click="toggleDeviceStatus(currentDevice)" class="cursor-pointer ml-2">
                {{ deviceStatusLabel(currentDevice) }}
              </el-tag>
            </div>
            <div class="card-items">
              <div v-for="(ti, ri) in detail.templateItems" :key="ti.itemID" class="card-item-row">
                <div class="item-label"><span class="font-bold">{{ ti.itemName }}</span><small v-if="ti.unit" class="text-gray ml-1">({{ ti.unit }})</small></div>
                <div class="item-controls">
                  <template v-if="ti.resultType === 'pass_fail'">
                    <span class="pass-toggle large">
                      <button :class="{active: currentDevice.results[ri]._checked === true}" @click="setPass(currentDevice, ri, true)">✓ 通过</button>
                      <button :class="{active: currentDevice.results[ri]._checked === false}" @click="setPass(currentDevice, ri, false)">✗ 不通过</button>
                    </span>
                  </template>
                  <template v-else-if="ti.resultType === 'number'">
                    <el-input-number v-model="currentDevice.results[ri]._numVal" :precision="2" size="large"
                      controls-position="right" style="width:160px" :class="getRangeClass(currentDevice.results[ri])"
                      @change="onNumChange(currentDevice, ri)" />
                  </template>
                  <template v-else>
                    <span class="pass-toggle large mr-2">
                      <button :class="{active: currentDevice.results[ri]._checked === true}" @click="setPass(currentDevice, ri, true)">✓ 通过</button>
                      <button :class="{active: currentDevice.results[ri]._checked === false}" @click="setPass(currentDevice, ri, false)">✗ 不通过</button>
                    </span>
                    <el-input-number v-model="currentDevice.results[ri]._numVal" :precision="2" size="large"
                      controls-position="right" style="width:140px" :class="getRangeClass(currentDevice.results[ri])"
                      @change="onNumChange(currentDevice, ri)" />
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- By Item -->
        <div v-if="inspectMode === 'byItem'">
          <div class="mb-3">
            <el-select v-model="currentItemIndex" placeholder="选择检测项" style="width:280px">
              <el-option v-for="(ti, idx) in detail.templateItems" :key="ti.itemID"
                :label="`${idx+1}. ${ti.itemName} (${doneCount(idx)}/${detail.devices.length})`" :value="idx" />
            </el-select>
            <el-button class="ml-2" size="small" @click="prevItem" :disabled="currentItemIndex === 0">上一项</el-button>
            <el-button size="small" @click="nextItem" :disabled="currentItemIndex >= detail.templateItems.length - 1">下一项</el-button>
          </div>
          <div v-for="dev in detail.devices" :key="dev.ID" class="inspect-card mb-2">
            <span class="text-sm font-bold mr-4 inline-block w-150px">{{ dev.lineNumber }}. {{ dev.sn }}</span>
            <template v-if="currentItem">
              <span class="pass-toggle mr-2" v-if="currentItem.resultType !== 'number'">
                <button :class="{active: dev.results[currentItemIndex]._checked === true}" @click="setPass(dev, currentItemIndex, true)">✓</button>
                <button :class="{active: dev.results[currentItemIndex]._checked === false}" @click="setPass(dev, currentItemIndex, false)">✗</button>
              </span>
              <el-input-number v-if="currentItem.resultType !== 'pass_fail'" v-model="dev.results[currentItemIndex]._numVal"
                :precision="2" size="small" controls-position="right" style="width:120px"
                :class="getRangeClass(dev.results[currentItemIndex])" @change="onNumChange(dev, currentItemIndex)" />
              <el-input v-model="dev.results[currentItemIndex].remark" size="small" placeholder="备注" style="width:120px" class="ml-2" />
            </template>
          </div>
        </div>

        <div class="detail-footer">
          <el-button type="primary" @click="saveResults">保存进度</el-button>
          <el-button v-if="detail.order.status === 2" type="success" @click="onComplete">完成检测</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { getProductionOrderList } from '@/plugin/inspection/api/production_order'
import { startInspection, saveResults as apiSaveResults, completeInspection, getInspectionDetail } from '@/plugin/inspection/api/work_order'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const activeTab = ref('pending')
const showDetail = ref(false)
const detailLoaded = ref(false)
const inspectMode = ref('byDevice')
const currentItemIndex = ref(0)
const currentDeviceIndex = ref(0)
const detail = ref({ order: {}, devices: [], templateItems: [] })

const searchInfo = reactive({ moNumber: '', model: '', page: 1, pageSize: 30 })

const currentDevice = computed(() => detail.value.devices[currentDeviceIndex.value] || null)
const currentItem = computed(() => detail.value.templateItems[currentItemIndex.value] || null)

const catLabel = (v) => ({ online: '线上', offline: '线下', foreign_trade: '外贸', custom: '定制款' }[v] || v)

const getList = async () => {
  loading.value = true
  const statusMap = { pending: 1, inspecting: 2, completed: 3 }
  try {
    const res = await getProductionOrderList({ ...searchInfo, status: statusMap[activeTab.value] })
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  } finally { loading.value = false }
}

const resetSearch = () => { searchInfo.moNumber = ''; searchInfo.model = ''; searchInfo.page = 1; getList() }
const onTabChange = () => { searchInfo.page = 1; getList() }

const openDetail = async (row) => {
  const res = await getInspectionDetail({ id: row.ID })
  if (res.code === 0) {
    const data = res.data
    data.devices.forEach(dev => {
      dev._status = dev.status || 'pending'
      dev.results.forEach(r => { r._checked = r.passResult; r._numVal = r.numberResult })
      calcDeviceStatus(dev)
    })
    detail.value = data
    detailLoaded.value = true
    showDetail.value = true
    currentDeviceIndex.value = 0
    currentItemIndex.value = 0
  }
}

const onStartInspect = async (row) => {
  await ElMessageBox.confirm('确定开始检测该订单？', '提示', { type: 'info' })
  const res = await startInspection({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('已开始检测')
    getList()
    await openDetail(row)
  }
}

// --- Detail helpers ---
const doneCount = (idx) => {
  let c = 0
  detail.value.devices.forEach(d => { if (d.results[idx]?._checked !== undefined && d.results[idx]?._checked !== null) c++ })
  return c
}

const getRangeClass = (r) => {
  if (!r || r._numVal === undefined || r._numVal === null || r._numVal === '') return ''
  const min = r.minValue, max = r.maxValue
  if ((min != null && r._numVal < min) || (max != null && r._numVal > max)) return 'out-range'
  if (min != null || max != null) return 'in-range'
  return ''
}

const deviceStatusLabel = (dev) => ({ pending: '待判定', pass: '通过', fail: '不通过' }[dev._status || 'pending'])
const deviceStatusTag = (dev) => { const s = dev._status || 'pending'; return s === 'pass' ? 'success' : s === 'fail' ? 'danger' : 'warning' }
const deviceRowClass = (dev) => { const s = dev._status || 'pending'; return s === 'fail' ? 'row-fail' : s === 'pass' ? 'row-pass' : '' }
const toggleDeviceStatus = (dev) => { const order = ['pending', 'pass', 'fail']; const idx = order.indexOf(dev._status || 'pending'); dev._status = order[(idx + 1) % 3] }

const setPass = (dev, ri, val) => { dev.results[ri]._checked = val; calcDeviceStatus(dev) }

const onNumChange = (dev, ri) => {
  const r = dev.results[ri]
  if (r._numVal === undefined || r._numVal === null || r._numVal === '') { r._checked = null }
  else {
    const min = r.minValue, max = r.maxValue
    if ((min != null && r._numVal < min) || (max != null && r._numVal > max)) { if (r._checked === undefined || r._checked === null) r._checked = false }
    else if (min != null || max != null) { if (r._checked === undefined || r._checked === null) r._checked = true }
  }
  calcDeviceStatus(dev)
}

const calcDeviceStatus = (dev) => {
  let hasFail = false, hasPass = false, hasUntested = false
  dev.results.forEach(r => {
    if (r._checked === true) hasPass = true
    else if (r._checked === false) hasFail = true
    else hasUntested = true
  })
  if (hasFail) dev._status = 'fail'
  else if (!hasUntested && hasPass) dev._status = 'pass'
  else dev._status = 'pending'
}

const saveResults = async () => {
  const deviceStatuses = [], deviceResults = []
  detail.value.devices.forEach(dev => {
    deviceStatuses.push({ deviceID: dev.ID, status: dev._status || 'pending' })
    dev.results.forEach(r => {
      deviceResults.push({
        deviceID: dev.ID, itemID: r.itemID,
        passResult: r._checked === true ? true : (r._checked === false ? false : null),
        numberResult: (r._numVal !== undefined && r._numVal !== null && r._numVal !== '') ? Number(r._numVal) : null,
        remark: r.remark || ''
      })
    })
  })
  const res = await apiSaveResults({ productionOrderID: detail.value.order.ID, deviceStatuses, deviceResults })
  if (res.code === 0) ElMessage.success('保存成功')
}

const onComplete = async () => {
  await ElMessageBox.confirm('确定完成该工单的检测？', '提示', { type: 'info' })
  const res = await completeInspection({ ID: detail.value.order.ID })
  if (res.code === 0) { ElMessage.success('检测完成'); showDetail.value = false; getList() }
}

const prevItem = () => { if (currentItemIndex.value > 0) currentItemIndex.value-- }
const nextItem = () => { if (currentItemIndex.value < detail.value.templateItems.length - 1) currentItemIndex.value++ }

getList()
</script>

<style scoped>
.inspect-page {
  max-width: 1400px; margin: 0 auto; padding: 16px;
  padding-bottom: 80px;
  background: var(--el-bg-color, #fff);
}
.page-header {
  display: flex; align-items: center; padding: 12px 16px;
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 6px 6px 0 0;
  background: var(--el-fill-color-light, #fafafa);
  margin-bottom: 0;
}
.page-body {
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-top: 0; border-radius: 0 0 6px 6px;
  padding: 16px;
  background: var(--el-bg-color, #fff);
}
.tab-section { margin-top: 0; }
.tab-section :deep(.el-tabs__header) { margin-bottom: 12px; }
.search-bar { margin-bottom: 12px; }
.table-wrap { margin-bottom: 12px; }
.pagination-wrap { display: flex; justify-content: flex-end; flex-wrap: wrap; }
.pagination-wrap :deep(.el-pagination) { flex-wrap: wrap; justify-content: flex-end; }
.pagination-wrap :deep(.el-pager) { flex-wrap: wrap; justify-content: center; }

.detail-toolbar {
  display: flex; align-items: center; padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-light, #e4e7ed);
  margin-bottom: 12px;
}

.inspect-table { border-collapse: collapse; width: 100%; font-size: 13px; }
.inspect-table th, .inspect-table td { border: 1px solid var(--el-border-color-light, #e4e7ed); padding: 6px 8px; text-align: center; white-space: nowrap; }
.inspect-table th { background: var(--el-fill-color-light, #f5f7fa); font-weight: 600; }
.inspect-table .fixed-col { position: sticky; left: 0; background: var(--el-bg-color, #fff); z-index: 2; }
.inspect-table th.fixed-col { background: var(--el-fill-color-light, #f5f7fa); z-index: 3; }
.inspect-table .sn-col { min-width: 130px; }

.inspect-card {
  display: flex; align-items: center; padding: 8px 12px;
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 4px; background: var(--el-bg-color, #fafafa);
}

.detail-footer {
  display: flex; gap: 8px; padding: 16px 0; margin-top: 16px;
  border-top: 1px solid var(--el-border-color-light, #e4e7ed);
}

.pass-toggle { display: inline-flex; gap: 2px; }
.pass-toggle button {
  width: 26px; height: 26px; border: 1px solid var(--el-border-color, #dcdfe6);
  border-radius: 4px; background: var(--el-bg-color, #fff);
  cursor: pointer; font-size: 13px; padding: 0; line-height: 1;
  color: var(--el-text-color-regular, #606266);
}
.pass-toggle button:first-child.active { background: var(--el-color-success, #67c23a); color: #fff; border-color: var(--el-color-success, #67c23a); }
.pass-toggle button:last-child.active { background: var(--el-color-danger, #f56c6c); color: #fff; border-color: var(--el-color-danger, #f56c6c); }
.row-pass td { background: var(--el-color-success-light-9, #f0f9eb) !important; }
.row-fail td { background: var(--el-color-danger-light-9, #fef0f0) !important; }
.cursor-pointer { cursor: pointer; }

.in-range :deep(.el-input__inner) { border-color: var(--el-color-success, #67c23a) !important; background: var(--el-color-success-light-9, #f0f9eb); }
.out-range :deep(.el-input__inner) { border-color: var(--el-color-danger, #f56c6c) !important; background: var(--el-color-danger-light-9, #fef0f0); }

.desktop-only { display: block; }
.mobile-only { display: none; }

@media (max-width: 768px) {
  .inspect-page { padding: 8px; }
  .page-body { padding: 8px; }
  .desktop-only { display: none !important; }
  .mobile-only { display: block !important; }
  .swiper-controls {
    display: flex; justify-content: space-between; align-items: center;
    padding: 12px 0; position: sticky; top: 0;
    background: var(--el-bg-color, #fff); z-index: 5;
    border-bottom: 1px solid var(--el-border-color-light, #eee);
    margin-bottom: 12px;
  }
  .device-card { border: 1px solid var(--el-border-color-light, #e4e7ed); border-radius: 8px; padding: 12px; }
  .card-header { display: flex; align-items: center; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid var(--el-border-color-light, #eee); }
  .card-item-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid var(--el-border-color-light, #f5f5f5); min-height: 48px; }
  .item-label { flex: 1; font-size: 15px; }
  .item-controls { display: flex; align-items: center; }
  .pass-toggle.large button { width: auto; height: 38px; padding: 0 14px; font-size: 15px; border-radius: 6px; }
}
</style>

<style>
html, body, #app {
  overflow-y: auto !important;
  height: auto !important;
  min-height: 100vh;
}
</style>
