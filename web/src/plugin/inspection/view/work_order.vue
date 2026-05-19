<template>
  <div>
    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="待检测" name="pending" />
      <el-tab-pane label="检测中" name="inspecting" />
      <el-tab-pane label="已完成" name="completed" />
    </el-tabs>

    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="MO号">
          <el-input v-model="searchInfo.moNumber" placeholder="请输入" clearable />
        </el-form-item>
        <el-form-item label="型号">
          <el-input v-model="searchInfo.model" placeholder="请输入" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData" border v-loading="loading">
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
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="e">
          <el-button v-if="activeTab === 'pending'" size="small" type="primary" @click="onStartInspect(e.row)">开始检测</el-button>
          <el-button v-else size="small" type="primary" @click="openInspectDrawer(e.row)">{{ activeTab === 'completed' ? '查看' : '检测' }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="gva-pagination">
      <el-pagination
        v-model:current-page="searchInfo.page"
        v-model:page-size="searchInfo.pageSize"
        :page-sizes="[10, 30, 50, 100]" :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="getList" @current-change="getList"
      />
    </div>

    <!-- Inspection Drawer -->
    <el-drawer v-model="drawerVisible" title="检测工单" size="95%" destroy-on-close>
      <template v-if="detailLoaded">
        <div class="mb-4 flex gap-4 flex-wrap text-sm">
          <span><b>MO号:</b> {{ detail.order.moNumber }}</span>
          <span><b>批次号:</b> {{ detail.order.batchNumber }}</span>
          <span><b>型号:</b> {{ detail.order.model }}</span>
          <span><b>类别:</b> {{ catLabel(detail.order.instrumentCategory) }}</span>
          <span><b>模板:</b> {{ detail.order.template?.name || '-' }}</span>
          <span><b>检测人:</b> {{ detail.order.inspectorID || '-' }}</span>
        </div>

        <el-tabs v-model="inspectMode">
          <el-tab-pane label="逐台检测" name="byDevice" />
          <el-tab-pane label="逐项检测" name="byItem" />
        </el-tabs>

        <!-- Mode 1: By Device (row=device, col=item) -->
        <div v-if="inspectMode === 'byDevice'" style="overflow-x:auto">
          <table class="inspect-table">
            <thead>
              <tr>
                <th class="fixed-col">序号</th>
                <th class="fixed-col sn-col">SN</th>
                <th class="fixed-col">判定</th>
                <th v-for="ti in detail.templateItems" :key="ti.itemID">
                  {{ ti.itemName }}<br /><small>{{ ti.unit || '' }}</small>
                </th>
                <th>备注</th>
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
                    <el-input-number
                      v-model="dev.results[ri]._numVal"
                      :precision="2" size="small" controls-position="right" style="width:100px"
                      :class="getRangeClass(dev.results[ri])"
                      @change="onNumChange(dev, ri)"
                    />
                  </template>
                  <template v-else>
                    <div class="flex gap-1 items-center">
                      <span class="pass-toggle">
                        <button :class="{active: dev.results[ri]._checked === true}" @click="setPass(dev, ri, true)">✓</button>
                        <button :class="{active: dev.results[ri]._checked === false}" @click="setPass(dev, ri, false)">✗</button>
                      </span>
                      <el-input-number
                        v-model="dev.results[ri]._numVal"
                        :precision="2" size="small" controls-position="right" style="width:90px"
                        :class="getRangeClass(dev.results[ri])"
                        @change="onNumChange(dev, ri)"
                      />
                    </div>
                  </template>
                </td>
                <td><el-input v-model="dev._remark" size="small" style="width:100px" placeholder="设备备注" /></td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Mode 2: By Item (pick one item, list all devices) -->
        <div v-else>
          <div class="mb-3">
            <el-select v-model="currentItemIndex" placeholder="选择检测项" style="width:280px" @change="onItemChange">
              <el-option v-for="(ti, idx) in detail.templateItems" :key="ti.itemID"
                :label="`${idx+1}. ${ti.itemName} (${doneCount(idx)}/${detail.devices.length})`" :value="idx" />
            </el-select>
            <el-button class="ml-2" @click="prevItem" :disabled="currentItemIndex === 0">上一项</el-button>
            <el-button @click="nextItem" :disabled="currentItemIndex >= detail.templateItems.length - 1">下一项</el-button>
          </div>
          <div v-for="dev in detail.devices" :key="dev.ID" class="inspect-card mb-2">
            <span class="text-sm font-bold mr-4 w-150px inline-block">{{ dev.lineNumber }}. {{ dev.sn }}</span>
            <template v-if="currentItem">
              <span class="pass-toggle mr-2" v-if="currentItem.resultType !== 'number'">
                <button :class="{active: dev.results[currentItemIndex]._checked === true}" @click="setPass(dev, currentItemIndex, true)">✓</button>
                <button :class="{active: dev.results[currentItemIndex]._checked === false}" @click="setPass(dev, currentItemIndex, false)">✗</button>
              </span>
              <el-input-number
                v-if="currentItem.resultType !== 'pass_fail'"
                v-model="dev.results[currentItemIndex]._numVal"
                :precision="2" size="small" controls-position="right" style="width:120px"
                :class="getRangeClass(dev.results[currentItemIndex])"
                @change="onNumChange(dev, currentItemIndex)"
                @keyup.enter="focusNextDevice(dev)"
              />
              <el-input v-model="dev.results[currentItemIndex].remark" size="small" placeholder="备注" style="width:120px" class="ml-2" />
            </template>
          </div>
        </div>

        <div class="mt-4">
          <el-button type="primary" @click="saveResults">保存进度</el-button>
          <el-button v-if="detail.order.status === 2" type="success" @click="onComplete">完成检测</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { getProductionOrderList } from '@/plugin/inspection/api/production_order'
import { startInspection, saveResults as apiSaveResults, completeInspection, getInspectionDetail } from '@/plugin/inspection/api/work_order'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const activeTab = ref('pending')
const drawerVisible = ref(false)
const inspectMode = ref('byDevice')
const detailLoaded = ref(false)
const detail = ref({ order: {}, devices: [], templateItems: [] })
const currentItemIndex = ref(0)

const searchInfo = reactive({ moNumber: '', model: '', page: 1, pageSize: 30 })

const currentItem = computed(() => detail.value.templateItems[currentItemIndex.value] || null)
const catLabel = (v) => ({ online: '线上', offline: '线下', foreign_trade: '外贸', custom: '定制款' }[v] || v)

const doneCount = (idx) => {
  let c = 0
  detail.value.devices.forEach(d => {
    const r = d.results[idx]
    if (r && (r._checked || r._numVal !== undefined)) c++
  })
  return c
}

const getRangeClass = (r) => {
  if (!r || r._numVal === undefined || r._numVal === null || r._numVal === '') return ''
  const min = r.minValue, max = r.maxValue
  if ((min != null && r._numVal < min) || (max != null && r._numVal > max)) return 'out-range'
  if (min != null || max != null) return 'in-range'
  return ''
}

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

const onStartInspect = async (row) => {
  await ElMessageBox.confirm('确定开始检测该订单？', '提示', { type: 'info' })
  const res = await startInspection({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('已开始检测')
    getList()
    // 直接打开检测抽屉
    row.status = 2
    openInspectDrawer(row)
  }
}

const openInspectDrawer = async (row) => {
  const res = await getInspectionDetail({ id: row.ID })
  if (res.code === 0) {
    const data = res.data
    // Add reactive helper fields to each result
    data.devices.forEach(dev => {
      dev._status = dev.status || 'pending'
      dev._remark = ''
      dev.results.forEach(r => {
        r._checked = r.passResult
        r._numVal = r.numberResult
      })
      calcDeviceStatus(dev)
    })
    detail.value = data
    detailLoaded.value = true
    drawerVisible.value = true
  }
}

const saveResults = async () => {
  const deviceStatuses = []
  const deviceResults = []
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
  if (res.code === 0) { ElMessage.success('检测完成'); drawerVisible.value = false; getList() }
}

const onItemChange = () => {}
const prevItem = () => { if (currentItemIndex.value > 0) currentItemIndex.value-- }
const nextItem = () => { if (currentItemIndex.value < detail.value.templateItems.length - 1) currentItemIndex.value++ }

// Auto-check/uncheck based on number range for "both" type
watch(() => detail.value.devices, { deep: true })

// Device status helpers
const deviceStatusLabel = (dev) => ({ pending: '待判定', pass: '通过', fail: '不通过' }[dev._status || 'pending'])
const deviceStatusTag = (dev) => {
  const s = dev._status || 'pending'
  return s === 'pass' ? 'success' : s === 'fail' ? 'danger' : 'warning'
}
const deviceRowClass = (dev) => {
  const s = dev._status || 'pending'
  return s === 'fail' ? 'row-fail' : s === 'pass' ? 'row-pass' : ''
}
const toggleDeviceStatus = (dev) => {
  const order = ['pending', 'pass', 'fail']
  const idx = order.indexOf(dev._status || 'pending')
  dev._status = order[(idx + 1) % 3]
}

const setPass = (dev, ri, val) => {
  dev.results[ri]._checked = val
  calcDeviceStatus(dev)
}

const onNumChange = (dev, ri) => {
  const r = dev.results[ri]
  if (r._numVal === undefined || r._numVal === null || r._numVal === '') {
    r._checked = null
  } else {
    const min = r.minValue, max = r.maxValue
    if ((min != null && r._numVal < min) || (max != null && r._numVal > max)) {
      if (r._checked === undefined || r._checked === null) r._checked = false
    } else if (min != null || max != null) {
      if (r._checked === undefined || r._checked === null) r._checked = true
    }
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

getList()
</script>

<style scoped>
.inspect-table { border-collapse: collapse; width: 100%; font-size: 13px; }
.inspect-table th, .inspect-table td { border: 1px solid #e4e7ed; padding: 6px 8px; text-align: center; white-space: nowrap; }
.inspect-table th { background: #f5f7fa; font-weight: 600; }
.inspect-table .fixed-col { position: sticky; left: 0; background: #fff; z-index: 2; }
.inspect-table th.fixed-col { background: #f5f7fa; z-index: 3; }
.inspect-table .sn-col { min-width: 130px; }
.inspect-card { display: flex; align-items: center; padding: 8px 12px; border: 1px solid #e4e7ed; border-radius: 4px; background: #fafafa; }
.in-range :deep(.el-input__inner) { border-color: #67c23a !important; background: #f0f9eb; }
.out-range :deep(.el-input__inner) { border-color: #f56c6c !important; background: #fef0f0; }
.pass-toggle { display: inline-flex; gap: 2px; }
.pass-toggle button {
  width: 26px; height: 26px; border: 1px solid #dcdfe6; border-radius: 4px;
  background: #fff; cursor: pointer; font-size: 13px; padding: 0; line-height: 1;
}
.pass-toggle button:first-child.active { background: #67c23a; color: #fff; border-color: #67c23a; }
.pass-toggle button:last-child.active { background: #f56c6c; color: #fff; border-color: #f56c6c; }
.row-pass td { background: #f0f9eb !important; }
.row-fail td { background: #fef0f0 !important; }
.cursor-pointer { cursor: pointer; }
</style>
