<template>
  <div class="inspect-page">
    <div class="page-header">
      <el-button @click="goBack">← 返回列表</el-button>
      <span class="text-lg font-bold ml-4" v-if="detailLoaded">检测工单</span>
    </div>

    <div v-if="detailLoaded">
      <div class="mb-4 flex gap-4 flex-wrap text-sm info-bar">
        <span><b>MO号:</b> {{ detail.order.moNumber }}</span>
        <span><b>批次号:</b> {{ detail.order.batchNumber }}</span>
        <span><b>型号:</b> {{ detail.order.model }}</span>
        <span><b>类别:</b> {{ catLabel(detail.order.instrumentCategory) }}</span>
        <span><b>模板:</b> {{ detail.order.template?.name || '-' }}</span>
        <span><b>检测人:</b> {{ detail.order.inspectorName || '-' }}</span>
      </div>

      <el-tabs v-model="inspectMode">
        <el-tab-pane label="逐台检测" name="byDevice" />
        <el-tab-pane label="逐项检测" name="byItem" />
      </el-tabs>

      <!-- Mode 1: Desktop table -->
      <div v-if="inspectMode === 'byDevice'" style="overflow-x:auto" class="desktop-only">
        <table class="inspect-table">
          <thead>
            <tr>
              <th class="fixed-col">序号</th>
              <th class="fixed-col sn-col">SN</th>
              <th class="fixed-col">判定</th>
              <th v-for="ti in detail.templateItems" :key="ti.itemID">{{ ti.itemName }}<br /><small>{{ ti.unit || '' }}</small></th>
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
              <td><el-input v-model="dev._remark" size="small" style="width:100px" placeholder="设备备注" /></td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mode 1 Mobile -->
      <div v-if="inspectMode === 'byDevice'" class="mobile-only">
        <div class="device-swiper">
          <div class="swiper-controls">
            <el-button :disabled="currentDeviceIndex === 0" @click="currentDeviceIndex--">◀ 上一台</el-button>
            <span class="text-lg font-bold">{{ currentDeviceIndex + 1 }} / {{ detail.devices.length }}</span>
            <el-button :disabled="currentDeviceIndex >= detail.devices.length - 1" @click="currentDeviceIndex++">下一台 ▶</el-button>
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
                <div class="item-label">
                  <span class="font-bold">{{ ti.itemName }}</span>
                  <small v-if="ti.unit" class="text-gray ml-1">({{ ti.unit }})</small>
                </div>
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
            <div class="card-remark mt-2">
              <el-input v-model="currentDevice._remark" placeholder="设备备注" size="large" />
            </div>
          </div>
        </div>
      </div>

      <!-- Mode 2: By Item -->
      <div v-if="inspectMode === 'byItem'">
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
            <el-input-number v-if="currentItem.resultType !== 'pass_fail'" v-model="dev.results[currentItemIndex]._numVal"
              :precision="2" size="small" controls-position="right" style="width:120px"
              :class="getRangeClass(dev.results[currentItemIndex])" @change="onNumChange(dev, currentItemIndex)"
              @keyup.enter="focusNextDevice(dev)" />
            <el-input v-model="dev.results[currentItemIndex].remark" size="small" placeholder="备注" style="width:120px" class="ml-2" />
          </template>
        </div>
      </div>

      <div class="mt-4 page-footer">
        <el-button type="primary" size="large" @click="saveResults">保存进度</el-button>
        <el-button v-if="detail.order.status === 2" type="success" size="large" @click="onComplete">完成检测</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { startInspection, saveResults as apiSaveResults, completeInspection, getInspectionDetail } from '@/plugin/inspection/api/work_order'

const router = useRouter()
const route = useRoute()

const detailLoaded = ref(false)
const inspectMode = ref('byDevice')
const currentItemIndex = ref(0)
const currentDeviceIndex = ref(0)
const detail = ref({ order: {}, devices: [], templateItems: [] })

const currentDevice = computed(() => detail.value.devices[currentDeviceIndex.value] || null)
const currentItem = computed(() => detail.value.templateItems[currentItemIndex.value] || null)

const catLabel = (v) => ({ online: '线上', offline: '线下', foreign_trade: '外贸', custom: '定制款' }[v] || v)

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

const goBack = () => router.back()

const loadDetail = async () => {
  const orderId = route.query.orderId
  if (!orderId) return
  const res = await getInspectionDetail({ id: orderId })
  if (res.code === 0) {
    const data = res.data
    data.devices.forEach(dev => {
      dev._status = dev.status || 'pending'
      dev._remark = ''
      dev.results.forEach(r => { r._checked = r.passResult; r._numVal = r.numberResult })
      calcDeviceStatus(dev)
    })
    detail.value = data
    detailLoaded.value = true
  }
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
  if (res.code === 0) { ElMessage.success('检测完成'); router.back() }
}

const onItemChange = () => {}
const prevItem = () => { if (currentItemIndex.value > 0) currentItemIndex.value-- }
const nextItem = () => { if (currentItemIndex.value < detail.value.templateItems.length - 1) currentItemIndex.value++ }

loadDetail()
</script>

<style scoped>
.inspect-page { padding: 16px; max-width: 100%; }
.page-header { display: flex; align-items: center; padding-bottom: 12px; border-bottom: 1px solid #e4e7ed; margin-bottom: 16px; }
.info-bar { background: #f5f7fa; padding: 10px 16px; border-radius: 4px; }
.page-footer { position: sticky; bottom: 0; background: #fff; padding: 12px 0; border-top: 1px solid #e4e7ed; z-index: 10; }
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
.pass-toggle button { width: 26px; height: 26px; border: 1px solid #dcdfe6; border-radius: 4px; background: #fff; cursor: pointer; font-size: 13px; padding: 0; line-height: 1; }
.pass-toggle button:first-child.active { background: #67c23a; color: #fff; border-color: #67c23a; }
.pass-toggle button:last-child.active { background: #f56c6c; color: #fff; border-color: #f56c6c; }
.row-pass td { background: #f0f9eb !important; }
.row-fail td { background: #fef0f0 !important; }
.cursor-pointer { cursor: pointer; }
.desktop-only { display: block; }
.mobile-only { display: none; }
@media (max-width: 768px) {
  .desktop-only { display: none !important; }
  .mobile-only { display: block !important; }
  .device-swiper { padding: 8px; }
  .swiper-controls { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; position: sticky; top: 0; background: #fff; z-index: 5; border-bottom: 1px solid #eee; margin-bottom: 12px; }
  .device-card { border: 1px solid #e4e7ed; border-radius: 8px; padding: 12px; }
  .card-header { display: flex; align-items: center; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #eee; }
  .card-item-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #f5f5f5; min-height: 48px; }
  .item-label { flex: 1; font-size: 15px; }
  .item-controls { display: flex; align-items: center; }
  .pass-toggle.large button { width: auto; height: 38px; padding: 0 14px; font-size: 15px; border-radius: 6px; }
}
</style>
