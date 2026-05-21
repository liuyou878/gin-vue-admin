<template>
  <div class="detail-page" v-if="detailLoaded">
    <div class="detail-toolbar">
      <el-button :icon="'ArrowLeft'" @click="goBack" />
      <span class="ml-2 text-sm">{{ detailInfo }}</span>
    </div>

    <el-tabs v-model="inspectMode" class="tab-bar">
      <el-tab-pane label="逐台检测" name="byDevice" />
      <el-tab-pane label="逐项检测" name="byItem" />
    </el-tabs>

    <div class="detail-scroll">
      <div v-if="inspectMode === 'byDevice'" class="desktop-only desktop-focus-toolbar">
        <div class="desktop-focus-mode-row">
          <el-radio-group v-model="desktopViewMode" size="default">
            <el-radio-button label="all">全表模式</el-radio-button>
            <el-radio-button label="single">单项模式</el-radio-button>
          </el-radio-group>
        </div>
        <div v-if="desktopViewMode === 'single'" class="desktop-focus-main">
          <el-button size="default" @click="prevItem" :disabled="currentItemIndex === 0">上一项</el-button>
          <el-select v-model="currentItemIndex" placeholder="选择检测项" size="default" class="desktop-focus-select">
            <el-option
              v-for="(ti, idx) in detail.templateItems"
              :key="ti.itemID"
              :label="`${idx + 1}. ${ti.itemName}`"
              :value="idx"
            />
          </el-select>
          <el-button size="default" @click="nextItem" :disabled="currentItemIndex >= detail.templateItems.length - 1">下一项</el-button>
        </div>
      </div>

      <!-- By Device: Desktop table -->
      <div v-if="inspectMode === 'byDevice' && desktopViewMode === 'all'" class="desktop-only" style="overflow-x:auto">
        <table class="inspect-table">
          <thead>
            <tr>
              <th class="fixed-col">序号</th>
              <th class="fixed-col sn-col">SN</th>
              <th class="fixed-col">判定</th>
              <th v-for="entry in visibleDesktopItems" :key="entry.item.itemID">
                {{ entry.item.itemName }}<br /><small>{{ entry.item.unit || '' }}</small>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="dev in detail.devices" :key="dev.ID" :class="deviceRowClass(dev)">
              <td class="fixed-col">{{ dev.lineNumber }}</td>
              <td class="fixed-col sn-col">{{ dev.sn }}</td>
              <td class="fixed-col">
                <el-tag :type="deviceStatusTag(dev)" size="small" @click="!isReadonly && toggleDeviceStatus(dev)" :class="!isReadonly && 'cursor-pointer'">
                  {{ deviceStatusLabel(dev) }}
                </el-tag>
              </td>
              <td v-for="entry in visibleDesktopItems" :key="entry.item.itemID">
                <template v-if="entry.item.resultType === 'pass_fail'">
                  <span class="pass-toggle">
                    <button :class="{active: dev.results[entry.index]._checked === true}" :disabled="isReadonly" @click="setPass(dev, entry.index, true)">✓</button>
                    <button :class="{active: dev.results[entry.index]._checked === false}" :disabled="isReadonly" @click="setPass(dev, entry.index, false)">✗</button>
                  </span>
                </template>
                <template v-else-if="entry.item.resultType === 'number'">
                  <el-input-number :disabled="isReadonly" v-model="dev.results[entry.index]._numVal" :precision="2" size="small"
                    controls-position="right" style="width:100px" :class="getRangeClass(dev.results[entry.index])"
                    @change="onNumChange(dev, entry.index)" />
                </template>
                <template v-else>
                  <div class="flex gap-1 items-center">
                    <span class="pass-toggle">
                      <button :class="{active: dev.results[entry.index]._checked === true}" :disabled="isReadonly" @click="setPass(dev, entry.index, true)">✓</button>
                      <button :class="{active: dev.results[entry.index]._checked === false}" :disabled="isReadonly" @click="setPass(dev, entry.index, false)">✗</button>
                    </span>
                    <el-input-number :disabled="isReadonly" v-model="dev.results[entry.index]._numVal" :precision="2" size="small"
                      controls-position="right" style="width:90px" :class="getRangeClass(dev.results[entry.index])"
                      @change="onNumChange(dev, entry.index)" />
                  </div>
                </template>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="inspectMode === 'byDevice' && desktopViewMode === 'single'" class="desktop-only">
        <div v-if="currentItem" class="single-mode-header">
          <div class="single-mode-title">
            <span>{{ currentItem.itemName }}</span>
            <small v-if="currentItem.unit">({{ currentItem.unit }})</small>
          </div>
          <div class="single-mode-progress">
            已完成 {{ doneCount(currentItemIndex) }} / {{ detail.devices.length }}
          </div>
        </div>

        <div class="single-mode-list">
          <div
            v-for="(dev, devIndex) in detail.devices"
            :key="dev.ID"
            :ref="el => setSingleModeRowRef(el, devIndex)"
            class="single-mode-row"
            :class="[deviceRowClass(dev), { 'single-mode-active': currentSingleDeviceIndex === devIndex }]"
          >
            <div class="single-mode-meta">
              <div class="single-mode-sn">{{ dev.lineNumber }}. {{ dev.sn }}</div>
              <el-tag :type="deviceStatusTag(dev)" size="large">
                {{ deviceStatusLabel(dev) }}
              </el-tag>
            </div>

            <div v-if="currentItem" class="single-mode-actions">
              <template v-if="currentItem.resultType === 'pass_fail'">
                <div class="single-mode-pass-toggle">
                  <button
                    class="single-action-btn single-action-pass"
                    :class="{ active: dev.results[currentItemIndex]._checked === true }"
                    :disabled="isReadonly"
                    @click="setPassAndAdvance(dev, devIndex, currentItemIndex, true)"
                  >
                    通过
                  </button>
                  <button
                    class="single-action-btn single-action-fail"
                    :class="{ active: dev.results[currentItemIndex]._checked === false }"
                    :disabled="isReadonly"
                    @click="setPassAndAdvance(dev, devIndex, currentItemIndex, false)"
                  >
                    不通过
                  </button>
                </div>
              </template>

              <template v-else-if="currentItem.resultType === 'number'">
                <el-input-number
                  :disabled="isReadonly"
                  v-model="dev.results[currentItemIndex]._numVal"
                  :precision="2"
                  size="large"
                  controls-position="right"
                  style="width: 180px"
                  :class="getRangeClass(dev.results[currentItemIndex])"
                  @change="onNumChange(dev, currentItemIndex)"
                />
              </template>

              <template v-else>
                <div class="single-mode-pass-toggle">
                  <button
                    class="single-action-btn single-action-pass"
                    :class="{ active: dev.results[currentItemIndex]._checked === true }"
                    :disabled="isReadonly"
                    @click="setPassAndAdvance(dev, devIndex, currentItemIndex, true)"
                  >
                    通过
                  </button>
                  <button
                    class="single-action-btn single-action-fail"
                    :class="{ active: dev.results[currentItemIndex]._checked === false }"
                    :disabled="isReadonly"
                    @click="setPassAndAdvance(dev, devIndex, currentItemIndex, false)"
                  >
                    不通过
                  </button>
                </div>
                <el-input-number
                  :disabled="isReadonly"
                  v-model="dev.results[currentItemIndex]._numVal"
                  :precision="2"
                  size="large"
                  controls-position="right"
                  style="width: 180px"
                  :class="getRangeClass(dev.results[currentItemIndex])"
                  @change="onNumChange(dev, currentItemIndex)"
                />
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- By Device: Mobile card -->
      <div v-if="inspectMode === 'byDevice'" class="mobile-only">
        <div class="swiper-controls">
          <el-button size="small" :disabled="currentDeviceIndex === 0" @click="currentDeviceIndex--">◀ 上一台</el-button>
          <span class="text-lg font-bold">{{ currentDeviceIndex + 1 }} / {{ detail.devices.length }}</span>
          <el-button size="small" :disabled="currentDeviceIndex >= detail.devices.length - 1" @click="currentDeviceIndex++">下一台 ▶</el-button>
        </div>
        <div class="device-card" v-if="currentDevice">
          <div class="card-header">
            <span class="text-lg font-bold">{{ currentDevice.sn }}</span>
            <el-tag :type="deviceStatusTag(currentDevice)" size="large" @click="toggleDeviceStatus(currentDevice)" class="cursor-pointer ml-2">{{ deviceStatusLabel(currentDevice) }}</el-tag>
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

      <div class="detail-footer" v-if="detail.order.status === 2">
        <el-button type="primary" size="large" @click="saveResults">保存进度</el-button>
        <el-button type="success" size="large" @click="onComplete">完成检测</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { saveResults as apiSaveResults, completeInspection, getInspectionDetail } from '@/plugin/inspection/api/work_order'

const route = useRoute()

const detailLoaded = ref(false)
const inspectMode = ref('byDevice')
const desktopViewMode = ref('all')
const currentItemIndex = ref(0)
const currentDeviceIndex = ref(0)
const currentSingleDeviceIndex = ref(0)
const singleModeRowRefs = ref([])
const detail = ref({ order: {}, devices: [], templateItems: [] })

const isReadonly = computed(() => detail.value.order.status === 3)
const currentDevice = computed(() => detail.value.devices[currentDeviceIndex.value] || null)
const currentItem = computed(() => detail.value.templateItems[currentItemIndex.value] || null)
const visibleDesktopItems = computed(() => {
  if (desktopViewMode.value === 'single') {
    const item = detail.value.templateItems[currentItemIndex.value]
    return item ? [{ item, index: currentItemIndex.value }] : []
  }
  return detail.value.templateItems.map((item, index) => ({ item, index }))
})
const detailInfo = computed(() => {
  const o = detail.value.order
  return [o.moNumber, o.batchNumber, o.model, o.inspectorName].filter(Boolean).join(' | ') || '-'
})

const goBack = () => { window.location.hash = '/inspectWorkOrder' }

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
const deviceStatusLabel = (d) => ({ pending: '待判定', pass: '通过', fail: '不通过' }[d._status || 'pending'])
const deviceStatusTag = (d) => { const s = d._status || 'pending'; return s === 'pass' ? 'success' : s === 'fail' ? 'danger' : 'warning' }
const deviceRowClass = (d) => { const s = d._status || 'pending'; return s === 'fail' ? 'row-fail' : s === 'pass' ? 'row-pass' : '' }
const toggleDeviceStatus = (d) => { const o = ['pending','pass','fail']; d._status = o[(o.indexOf(d._status||'pending') + 1) % 3] }
const setPass = (dev, ri, val) => { dev.results[ri]._checked = val; calcDeviceStatus(dev) }
const setSingleModeRowRef = (el, index) => {
  if (!el) return
  singleModeRowRefs.value[index] = el
}
const scrollToSingleModeRow = async (index) => {
  await nextTick()
  const row = singleModeRowRefs.value[index]
  if (row?.scrollIntoView) {
    row.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
  }
}
const moveToNextSingleModeTarget = async (currentIndex) => {
  const lastDeviceIndex = Math.max(detail.value.devices.length - 1, 0)
  if (currentIndex < lastDeviceIndex) {
    const nextIndex = currentIndex + 1
    currentSingleDeviceIndex.value = nextIndex
    await scrollToSingleModeRow(nextIndex)
    return
  }

  if (currentItemIndex.value < detail.value.templateItems.length - 1) {
    currentItemIndex.value += 1
    currentSingleDeviceIndex.value = 0
    await scrollToSingleModeRow(0)
  }
}
const setPassAndAdvance = async (dev, devIndex, resultIndex, val) => {
  setPass(dev, resultIndex, val)
  if (desktopViewMode.value === 'single') {
    await moveToNextSingleModeTarget(devIndex)
  }
}
const onNumChange = (dev, ri) => {
  const r = dev.results[ri]
  if (r._numVal === undefined || r._numVal === null || r._numVal === '') r._checked = null
  else {
    const min = r.minValue, max = r.maxValue
    if ((min != null && r._numVal < min) || (max != null && r._numVal > max)) { if (r._checked == null) r._checked = false }
    else if (min != null || max != null) { if (r._checked == null) r._checked = true }
  }
  calcDeviceStatus(dev)
}
const calcDeviceStatus = (dev) => {
  let f = false, p = false, u = false
  dev.results.forEach(r => { if (r._checked === true) p = true; else if (r._checked === false) f = true; else u = true })
  dev._status = f ? 'fail' : (!u && p ? 'pass' : 'pending')
}

const buildSavePayload = () => {
  const ds = [], dr = []
  detail.value.devices.forEach(dev => {
    ds.push({ deviceID: dev.ID, status: dev._status || 'pending' })
    dev.results.forEach(r => dr.push({ deviceID: dev.ID, itemID: r.itemID, passResult: r._checked, numberResult: (r._numVal !== undefined && r._numVal !== null && r._numVal !== '') ? Number(r._numVal) : null, remark: r.remark || '' }))
  })
  return { batchID: detail.value.order.ID, deviceStatuses: ds, deviceResults: dr }
}

const saveResults = async (silent = false) => {
  const res = await apiSaveResults(buildSavePayload())
  if (res.code === 0) {
    if (!silent) ElMessage.success('保存成功')
    return true
  }
  return false
}
const onComplete = async () => {
  await ElMessageBox.confirm('确定完成该工单的检测？', '提示', { type: 'info' })
  const saved = await saveResults(true)
  if (!saved) {
    ElMessage.error('完成前保存检测结果失败')
    return
  }
  const res = await completeInspection({ ID: detail.value.order.ID })
  if (res.code === 0) {
    ElMessage.success('检测完成')
    window.location.hash = '/inspectWorkOrder'
  }
}
const prevItem = () => { if (currentItemIndex.value > 0) currentItemIndex.value-- }
const nextItem = () => { if (currentItemIndex.value < detail.value.templateItems.length - 1) currentItemIndex.value++ }

watch(
  () => [desktopViewMode.value, currentItemIndex.value],
  () => {
    currentSingleDeviceIndex.value = 0
    singleModeRowRefs.value = []
  }
)

const loadDetail = async () => {
  const batchId = route.query.batchId
  if (!batchId) return
  const res = await getInspectionDetail({ id: batchId })
  if (res.code === 0) {
    const data = res.data
    data.devices.forEach(dev => {
      dev._status = dev.status || 'pending'
      dev.results.forEach(r => { r._checked = r.passResult; r._numVal = r.numberResult })
      calcDeviceStatus(dev)
    })
    detail.value = data; detailLoaded.value = true
  }
}
loadDetail()
</script>

<style scoped>
.detail-page { height: 100vh; display: flex; flex-direction: column; background: var(--el-bg-color, #fff); overflow: hidden; }
.detail-toolbar { display: flex; align-items: center; padding: 10px 16px; border-bottom: 1px solid var(--el-border-color-light, #e4e7ed); background: var(--el-fill-color-light, #fafafa); flex-shrink: 0; }
.tab-bar { flex-shrink: 0; padding: 0 16px; background: var(--el-bg-color, #fff); box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.tab-bar :deep(.el-tabs__header) { margin-bottom: 0; }
.detail-scroll { flex: 1; overflow-y: auto; padding: 16px; }
.detail-footer { display: flex; gap: 8px; padding: 16px 0; margin-top: 16px; border-top: 1px solid var(--el-border-color-light, #e4e7ed); }
.desktop-focus-toolbar { display: flex; flex-direction: column; align-items: flex-start; gap: 10px; margin-bottom: 12px; padding: 12px; border: 1px solid var(--el-border-color-light, #e4e7ed); border-radius: 8px; background: var(--el-fill-color-lighter, #fafcff); }
.desktop-focus-mode-row { display: flex; align-items: center; gap: 8px; }
.desktop-focus-main { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.desktop-focus-select { width: 320px; }
.single-mode-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 12px; padding: 10px 12px; border-radius: 8px; background: var(--el-fill-color-light, #f5f7fa); }
.single-mode-title { font-size: 18px; font-weight: 600; color: var(--el-text-color-primary, #303133); }
.single-mode-title small { margin-left: 4px; color: var(--el-text-color-secondary, #909399); font-size: 13px; }
.single-mode-progress { color: var(--el-text-color-secondary, #909399); font-size: 13px; }
.single-mode-list { display: flex; flex-direction: column; gap: 10px; }
.single-mode-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 14px 16px; border: 1px solid var(--el-border-color-light, #e4e7ed); border-radius: 10px; background: var(--el-bg-color, #fff); }
.single-mode-active { box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.18); border-color: rgba(64, 158, 255, 0.45); }
.single-mode-meta { display: flex; align-items: center; gap: 12px; min-width: 260px; }
.single-mode-sn { font-size: 16px; font-weight: 600; color: var(--el-text-color-primary, #303133); }
.single-mode-actions { display: flex; align-items: center; gap: 12px; justify-content: flex-start; flex: 1; }
.single-mode-pass-toggle { display: flex; align-items: center; gap: 10px; }
.single-action-btn { min-width: 92px; height: 42px; border: 1px solid var(--el-border-color, #dcdfe6); border-radius: 8px; background: var(--el-bg-color, #fff); cursor: pointer; font-size: 15px; font-weight: 600; transition: all 0.15s ease; }
.single-action-btn:disabled { cursor: not-allowed; opacity: 0.6; }
.single-action-pass { color: #15803d; }
.single-action-fail { color: #b91c1c; }
.single-action-pass.active { background: #67c23a; border-color: #67c23a; color: #fff; }
.single-action-fail.active { background: #f56c6c; border-color: #f56c6c; color: #fff; }

.inspect-table { border-collapse: collapse; width: 100%; font-size: 13px; }
.inspect-table th, .inspect-table td { border: 1px solid var(--el-border-color-light, #e4e7ed); padding: 6px 8px; text-align: center; white-space: nowrap; }
.inspect-table th { background: var(--el-fill-color-light, #f5f7fa); font-weight: 600; }
.inspect-table .fixed-col { position: sticky; left: 0; background: var(--el-bg-color, #fff); z-index: 2; }
.inspect-table th.fixed-col { background: var(--el-fill-color-light, #f5f7fa); z-index: 3; }
.inspect-table .sn-col { min-width: 130px; }
.inspect-card { display: flex; align-items: center; padding: 8px 12px; border: 1px solid var(--el-border-color-light, #e4e7ed); border-radius: 4px; background: var(--el-bg-color, #fafafa); }

.pass-toggle { display: inline-flex; gap: 2px; }
.pass-toggle button { width: 26px; height: 26px; border: 1px solid var(--el-border-color, #dcdfe6); border-radius: 4px; background: var(--el-bg-color, #fff); cursor: pointer; font-size: 13px; padding: 0; line-height: 1; color: var(--el-text-color-regular, #606266); }
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
  .desktop-only { display: none !important; }
  .mobile-only { display: block !important; }
  .swiper-controls { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; position: sticky; top: 0; background: var(--el-bg-color, #fff); z-index: 11; border-bottom: 1px solid var(--el-border-color-light, #eee); margin-bottom: 12px; box-shadow: 0 2px 6px rgba(0,0,0,0.08); }
  .device-card { border: 1px solid var(--el-border-color-light, #e4e7ed); border-radius: 8px; padding: 12px; }
  .card-header { display: flex; align-items: center; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid var(--el-border-color-light, #eee); }
  .card-item-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid var(--el-border-color-light, #f5f5f5); min-height: 48px; }
  .item-label { flex: 1; font-size: 15px; }
  .item-controls { display: flex; align-items: center; }
  .pass-toggle.large button { width: auto; height: 38px; padding: 0 14px; font-size: 15px; border-radius: 6px; }
}
</style>
