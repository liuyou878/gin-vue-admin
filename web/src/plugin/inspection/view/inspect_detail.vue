<template>
  <div class="detail-page" v-if="detailLoaded">
    <div class="detail-toolbar">
      <el-button :icon="'ArrowLeft'" @click="goBack" />
      <span class="detail-info">{{ detailInfo }}</span>
    </div>

    <InspectResultOverview
      v-if="showResultOverview"
      :detail="detail"
      @refresh="loadDetail"
    />

    <template v-else>
    <el-tabs v-model="inspectMode" class="tab-bar">
      <el-tab-pane label="按项检测" name="byItem" />
      <el-tab-pane label="按台检测" name="byDevice" />
    </el-tabs>

    <div class="detail-scroll">
      <section class="workbench-header">
        <div class="subject-row">
          <el-button
            size="default"
            :disabled="currentSubjectIndex <= 0"
            @click="prevSubject"
          >
            {{ inspectMode === 'byItem' ? '上一项' : '上一台' }}
          </el-button>

          <el-select
            v-if="inspectMode === 'byItem'"
            v-model="currentItemIndex"
            filterable
            placeholder="选择检测项"
            class="subject-select"
          >
            <el-option
              v-for="(item, index) in detail.templateItems"
              :key="item.itemID"
              :label="`${index + 1}. ${item.itemName} (${doneCount(index)}/${detail.devices.length})`"
              :value="index"
            />
          </el-select>

          <el-select
            v-else
            v-model="currentDeviceIndex"
            filterable
            placeholder="选择设备"
            class="subject-select"
          >
            <el-option
              v-for="(device, index) in detail.devices"
              :key="device.ID"
              :label="`${device.lineNumber}. ${device.sn} (${deviceDoneCount(device)}/${detail.templateItems.length})`"
              :value="index"
            />
          </el-select>

          <el-button
            size="default"
            :disabled="currentSubjectIndex >= currentSubjectTotal - 1"
            @click="nextSubject"
          >
            {{ inspectMode === 'byItem' ? '下一项' : '下一台' }}
          </el-button>
        </div>

        <div class="subject-summary">
          <div>
            <div class="subject-title">{{ currentSubjectTitle }}</div>
            <div class="subject-subtitle">{{ currentSubjectSubtitle }}</div>
          </div>
          <el-tag :type="currentSubjectDone === currentSubjectTotalRows ? 'success' : 'warning'" size="large">
            {{ currentSubjectDone }} / {{ currentSubjectTotalRows }}
          </el-tag>
        </div>
      </section>

      <section class="inspection-list">
        <div
          v-for="row in inspectionRows"
          :key="row.key"
          class="inspection-row"
          :class="deviceRowClass(row.device)"
        >
          <div class="row-left">
            <div class="row-title">{{ row.title }}</div>
            <div class="row-meta">
              <el-tag :type="resultStatusTag(row.result)" size="small">
                {{ resultStatusLabel(row.result) }}
              </el-tag>
              <span v-if="row.meta">{{ row.meta }}</span>
              <span class="save-meta" :class="saveStateClass(row.result)">
                {{ resultSaveLabel(row.result) }}
              </span>
            </div>
          </div>

          <div class="row-controls">
            <div v-if="row.result.resultType !== 'number'" class="pass-actions">
              <button
                class="action-btn action-pass"
                :class="{ active: row.result._checked === true }"
                :disabled="isReadonly"
                @click="setPassAndSave(row.device, row.resultIndex, true)"
              >
                通过
              </button>
              <button
                class="action-btn action-fail"
                :class="{ active: row.result._checked === false }"
                :disabled="isReadonly"
                @click="setPassAndSave(row.device, row.resultIndex, false)"
              >
                不通过
              </button>
            </div>

            <el-input-number
              v-if="row.result.resultType !== 'pass_fail'"
              v-model="row.result._numVal"
              :disabled="isReadonly"
              controls-position="right"
              class="number-input"
              :class="getRangeClass(row.result)"
              @change="onNumChangeAndSave(row.device, row.resultIndex)"
            />

            <el-input
              v-model="row.result.remark"
              :disabled="isReadonly"
              placeholder="备注"
              class="remark-input"
              clearable
              @blur="saveSingleResult(row.device, row.resultIndex)"
            />
          </div>
        </div>

        <el-empty v-if="inspectionRows.length === 0" description="暂无检测数据" />
      </section>

      <div class="detail-footer" v-if="detail.order.status === 2 || hasRecheckingDevices">
        <el-button v-if="detail.order.status === 2" type="success" size="large" @click="onComplete">完成检测</el-button>
        <el-button v-if="hasRecheckingDevices" type="success" size="large" @click="onCompleteRecheck">完成复检</el-button>
      </div>
    </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import InspectResultOverview from '@/plugin/inspection/components/InspectResultOverview.vue'
import {
  completeInspection,
  completeRecheck,
  getInspectionDetail,
  saveResults as apiSaveResults,
  saveSingleResult as apiSaveSingleResult
} from '@/plugin/inspection/api/work_order'

const route = useRoute()

const detailLoaded = ref(false)
const inspectMode = ref('byItem')
const currentItemIndex = ref(0)
const currentDeviceIndex = ref(0)
const detail = ref({ order: {}, devices: [], templateItems: [] })

const currentDevice = computed(() => detail.value.devices[currentDeviceIndex.value] || null)
const currentItem = computed(() => detail.value.templateItems[currentItemIndex.value] || null)
const hasRecheckingDevices = computed(() => detail.value.devices.some((device) => device._startedRecheck || (device._status || device.status) === 'rechecking'))
const showResultOverview = computed(() => detail.value.order.status >= 3 && !hasRecheckingDevices.value)
const isReadonly = computed(() => detail.value.order.status >= 3 && !hasRecheckingDevices.value)

const detailInfo = computed(() => {
  const order = detail.value.order
  return [order.moNumber, order.batchNumber, order.model, order.inspectorName].filter(Boolean).join(' | ') || '-'
})

const currentSubjectIndex = computed(() => (inspectMode.value === 'byItem' ? currentItemIndex.value : currentDeviceIndex.value))
const currentSubjectTotal = computed(() => (inspectMode.value === 'byItem' ? detail.value.templateItems.length : detail.value.devices.length))
const currentSubjectTitle = computed(() => {
  if (inspectMode.value === 'byItem') {
    return currentItem.value?.itemName || '-'
  }
  return currentDevice.value?.sn || '-'
})
const currentSubjectSubtitle = computed(() => {
  if (inspectMode.value === 'byItem') {
    const item = currentItem.value
    if (!item) return '请选择检测项'
    const parts = []
    if (item.unit) parts.push(`单位：${item.unit}`)
    if (item.minValue != null || item.maxValue != null) {
      parts.push(`范围：${item.minValue ?? '-'} ~ ${item.maxValue ?? '-'}`)
    }
    return parts.join(' | ') || '按当前检测项批量检测所有设备'
  }
  return currentDevice.value ? `按当前设备完成 ${detail.value.templateItems.length} 个检测项` : '请选择设备'
})
const currentSubjectDone = computed(() => {
  if (inspectMode.value === 'byItem') return doneCount(currentItemIndex.value)
  return currentDevice.value ? deviceDoneCount(currentDevice.value) : 0
})
const currentSubjectTotalRows = computed(() => (inspectMode.value === 'byItem' ? detail.value.devices.length : detail.value.templateItems.length))

const inspectionRows = computed(() => {
  if (inspectMode.value === 'byItem') {
    const item = currentItem.value
    if (!item) return []
    return detail.value.devices.map((device, deviceIndex) => ({
      key: `device-${device.ID}-${item.itemID}`,
      title: `${device.lineNumber || deviceIndex + 1}. ${device.sn}`,
      meta: deviceStatusLabel(device),
      device,
      result: device.results[currentItemIndex.value],
      resultIndex: currentItemIndex.value
    })).filter((row) => row.result)
  }

  const device = currentDevice.value
  if (!device) return []
  return detail.value.templateItems.map((item, itemIndex) => ({
    key: `item-${device.ID}-${item.itemID}`,
    title: `${itemIndex + 1}. ${item.itemName}`,
    meta: item.unit ? `单位：${item.unit}` : '',
    device,
    result: device.results[itemIndex],
    resultIndex: itemIndex
  })).filter((row) => row.result)
})

const goBack = () => {
  window.location.hash = '/inspectWorkOrder'
}

const prevSubject = () => {
  if (inspectMode.value === 'byItem' && currentItemIndex.value > 0) currentItemIndex.value -= 1
  if (inspectMode.value === 'byDevice' && currentDeviceIndex.value > 0) currentDeviceIndex.value -= 1
}

const nextSubject = () => {
  if (inspectMode.value === 'byItem' && currentItemIndex.value < detail.value.templateItems.length - 1) currentItemIndex.value += 1
  if (inspectMode.value === 'byDevice' && currentDeviceIndex.value < detail.value.devices.length - 1) currentDeviceIndex.value += 1
}

const doneCount = (itemIndex) => {
  let count = 0
  detail.value.devices.forEach((device) => {
    if (resultCompletedForStatus(device.results[itemIndex])) count += 1
  })
  return count
}

const deviceDoneCount = (device) => {
  return (device?.results || []).filter(resultCompletedForStatus).length
}

const resultCompletedForStatus = (result) => {
  if (!result) return false
  if (result.resultType === 'number') {
    return result._numVal !== undefined && result._numVal !== null && result._numVal !== ''
  }
  if (result.resultType === 'pass_fail') {
    return result._checked === true || result._checked === false
  }
  return (result._checked === true || result._checked === false) &&
    result._numVal !== undefined &&
    result._numVal !== null &&
    result._numVal !== ''
}

const resultStatusLabel = (result) => {
  if (!resultCompletedForStatus(result)) return '未完成'
  if (result._checked === false) return '不通过'
  if (result._checked === true) return '通过'
  return '已填写'
}

const resultStatusTag = (result) => {
  if (!resultCompletedForStatus(result)) return 'info'
  if (result._checked === false) return 'danger'
  if (result._checked === true) return 'success'
  return 'warning'
}

const getRangeClass = (result) => {
  if (!result || result._numVal === undefined || result._numVal === null || result._numVal === '') return ''
  const min = result.minValue
  const max = result.maxValue
  if ((min != null && result._numVal < min) || (max != null && result._numVal > max)) return 'out-range'
  if (min != null || max != null) return 'in-range'
  return ''
}

const setPass = (device, resultIndex, value) => {
  device.results[resultIndex]._checked = value
  calcDeviceStatus(device)
}

const setPassAndSave = async (device, resultIndex, value) => {
  setPass(device, resultIndex, value)
  await saveSingleResult(device, resultIndex)
}

const onNumChange = (device, resultIndex) => {
  const result = device.results[resultIndex]
  if (result._numVal === undefined || result._numVal === null || result._numVal === '') {
    result._checked = null
  } else {
    const min = result.minValue
    const max = result.maxValue
    if ((min != null && result._numVal < min) || (max != null && result._numVal > max)) {
      result._checked = false
    } else if (min != null || max != null) {
      result._checked = true
    }
  }
  calcDeviceStatus(device)
}

const onNumChangeAndSave = async (device, resultIndex) => {
  onNumChange(device, resultIndex)
  await saveSingleResult(device, resultIndex)
}

const saveStateClass = (result) => {
  if (result._saveState === 'saving') return 'saving'
  if (result._saveState === 'error') return 'error'
  if (result.inspectorName || result.inspectedAt) return 'saved'
  return ''
}

const formatShortTime = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `${month}-${day} ${hour}:${minute}`
}

const resultSaveLabel = (result) => {
  if (result._saveState === 'saving') return '保存中...'
  if (result._saveState === 'error') return '保存失败'
  if (result.inspectorName || result.inspectedAt) {
    return `${result.inspectorName || '已保存'}${result.inspectedAt ? ` · ${formatShortTime(result.inspectedAt)}` : ''}`
  }
  return '未保存'
}

const resultHasValue = (result) => {
  if (!result) return false
  return result._checked === true ||
    result._checked === false ||
    (result._numVal !== undefined && result._numVal !== null && result._numVal !== '') ||
    Boolean((result.remark || '').trim())
}

const resultWasSaved = (result) => {
  if (!result) return false
  return Boolean(result.inspectorName || result.inspectedAt) ||
    result.passResult === true ||
    result.passResult === false ||
    result.numberResult !== undefined && result.numberResult !== null ||
    Boolean((result.remark || '').trim())
}

const saveSingleResult = async (device, resultIndex) => {
  if (isReadonly.value) return false
  const result = device?.results?.[resultIndex]
  if (!result) return false
  if (!resultHasValue(result) && !resultWasSaved(result)) return false
  result._saveState = 'saving'
  try {
    const numberResult = result._numVal !== undefined && result._numVal !== null && result._numVal !== '' ? Number(result._numVal) : null
    const status = detail.value.order.status >= 3 && (device._startedRecheck || device.status === 'rechecking')
      ? ''
      : (device._status || 'pending')
    const res = await apiSaveSingleResult({
      batchID: detail.value.order.ID,
      deviceID: device.ID,
      itemID: result.itemID,
      passResult: result._checked,
      numberResult,
      remark: result.remark || '',
      status
    })
    if (res.code !== 0) {
      result._saveState = 'error'
      return false
    }
    const data = res.data || {}
    result.passResult = data.passResult
    result.numberResult = data.numberResult
    result.remark = data.remark || ''
    result.inspectorID = data.inspectorID
    result.inspectorName = data.inspectorName
    result.inspectedAt = data.inspectedAt
    result._saveState = 'saved'
    return true
  } catch (error) {
    result._saveState = 'error'
    return false
  }
}

const calcDeviceStatus = (device) => {
  if (device._status === 'rework' || device._status === 'pending_recheck') return
  let failed = false
  let passed = false
  let unfinished = false
  device.results.forEach((result) => {
    if (!resultCompletedForStatus(result)) {
      unfinished = true
      return
    }
    if (result._checked === false) failed = true
    else passed = true
  })
  if (device._startedRecheck) {
    device._status = 'rechecking'
    return
  }
  device._status = failed ? 'fail' : (!unfinished && passed ? 'pass' : 'pending')
}

const deviceStatusLabel = (device) => {
  const status = device?._status || device?.status || 'pending'
  const progress = deviceProgressText(device)
  if (status === 'pending') return `待完成${progress}`
  return ({
    returned: '待生产接收',
    pending_recheck: '待复检',
    rechecking: '复检中',
    pass: '通过',
    fail: '不通过',
    rework: '返工中'
  }[status] || status)
}

const deviceProgressText = (device) => {
  const results = device?.results || []
  if (!results.length) return ''
  const done = results.filter(resultCompletedForStatus).length
  return done > 0 && done < results.length ? ` ${done}/${results.length}` : ''
}

const deviceStatusTag = (device) => {
  const status = device?._status || device?.status || 'pending'
  return status === 'pass' ? 'success' :
    status === 'fail' ? 'danger' :
      status === 'returned' || status === 'rework' ? 'warning' :
        status === 'pending_recheck' ? 'primary' :
          status === 'rechecking' ? 'warning' : 'info'
}

const deviceRowClass = (device) => {
  const status = device?._status || device?.status || 'pending'
  return status === 'fail' ? 'row-fail' :
    status === 'pass' ? 'row-pass' :
      status === 'returned' || status === 'rework' ? 'row-rework' :
        status === 'pending_recheck' || status === 'rechecking' ? 'row-recheck' : ''
}

const buildSavePayload = () => {
  const deviceStatuses = []
  const deviceResults = []
  detail.value.devices.forEach((device) => {
    if (isReadonly.value) return
    if (detail.value.order.status >= 3 && !device._startedRecheck) return
    deviceStatuses.push({ deviceID: device.ID, status: device._status || 'pending' })
    device.results.forEach((result) => {
      deviceResults.push({
        deviceID: device.ID,
        itemID: result.itemID,
        passResult: result._checked,
        numberResult: result._numVal !== undefined && result._numVal !== null && result._numVal !== '' ? Number(result._numVal) : null,
        remark: result.remark || ''
      })
    })
  })
  return { batchID: detail.value.order.ID, deviceStatuses, deviceResults }
}

const saveResults = async (silent = false) => {
  const res = await apiSaveResults(buildSavePayload())
  if (res.code === 0) {
    if (!silent) ElMessage.success('保存成功')
    await loadDetail()
    return true
  }
  return false
}

const shouldValidateDevice = (device) => {
  if (detail.value.order.status === 2) return (device._status || device.status) !== 'rework'
  if (hasRecheckingDevices.value) return (device._status || device.status) === 'rechecking' || device._startedRecheck
  return false
}

const getIncompleteSummary = () => {
  const missing = []
  detail.value.devices.forEach((device) => {
    if (!shouldValidateDevice(device)) return
    device.results.forEach((result) => {
      if (!resultCompletedForStatus(result)) missing.push(`${device.sn} / ${result.itemName}`)
    })
  })
  return missing
}

const onComplete = async () => {
  const saved = await saveResults(true)
  if (!saved) {
    ElMessage.error('完成前保存检测结果失败')
    return
  }
  const missing = getIncompleteSummary()
  if (missing.length > 0) {
    ElMessageBox.alert(
      `还有 ${missing.length} 个检测项未完成，请全部完成后再提交。\n\n${missing.slice(0, 8).join('\n')}${missing.length > 8 ? '\n...' : ''}`,
      '检测未完成',
      { type: 'warning' }
    )
    return
  }
  await ElMessageBox.confirm('确定完成该工单的检测？完成后才能处理不合格设备打回生产。', '提示', { type: 'info' })
  const res = await completeInspection({ ID: detail.value.order.ID })
  if (res.code === 0) {
    ElMessage.success('已提交待确认')
    await loadDetail()
  }
}

const onCompleteRecheck = async () => {
  const saved = await saveResults(true)
  if (!saved) {
    ElMessage.error('完成前保存复检结果失败')
    return
  }
  const missing = getIncompleteSummary()
  if (missing.length > 0) {
    ElMessageBox.alert(
      `还有 ${missing.length} 个复检项未完成，请全部完成后再提交。\n\n${missing.slice(0, 8).join('\n')}${missing.length > 8 ? '\n...' : ''}`,
      '复检未完成',
      { type: 'warning' }
    )
    return
  }
  await ElMessageBox.confirm('确定完成本次复检？完成后如果仍不合格，可以再次打回生产。', '提示', { type: 'info' })
  const res = await completeRecheck({ ID: detail.value.order.ID })
  if (res.code === 0) {
    ElMessage.success('复检完成')
    await loadDetail()
  }
}

const loadDetail = async () => {
  const batchId = route.query.batchId
  if (!batchId) return
  const res = await getInspectionDetail({ id: batchId })
  if (res.code !== 0) return
  const data = res.data
  data.devices.forEach((device) => {
    device._status = device.status || 'pending'
    device.results.forEach((result) => {
      result._checked = result.passResult
      result._numVal = result.numberResult
      result._saveState = ''
    })
    device._startedRecheck = device._status === 'rechecking'
    if (device._status !== 'returned' && device._status !== 'rework' && device._status !== 'pending_recheck' && device._status !== 'rechecking') {
      calcDeviceStatus(device)
    }
  })
  detail.value = data
  if (currentItemIndex.value >= detail.value.templateItems.length) currentItemIndex.value = 0
  if (currentDeviceIndex.value >= detail.value.devices.length) currentDeviceIndex.value = 0
  detailLoaded.value = true
}

loadDetail()
</script>

<style scoped>
.detail-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color, #fff);
  overflow: hidden;
}

.detail-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--el-border-color-light, #e4e7ed);
  background: var(--el-fill-color-light, #fafafa);
  flex-shrink: 0;
}

.detail-info {
  color: var(--el-text-color-secondary, #909399);
  font-size: 13px;
}

.tab-bar {
  flex-shrink: 0;
  padding: 0 16px;
  background: var(--el-bg-color, #fff);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.tab-bar :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.detail-scroll {
  flex: 1;
  overflow: auto;
  padding: 0 16px 16px;
  min-width: 0;
}

.workbench-header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 0 -16px 14px;
  padding: 14px;
  border: 0;
  border-bottom: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 0 0 14px 14px;
  background: var(--el-bg-color, #fff);
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.08);
}

.subject-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.subject-select {
  flex: 1;
  min-width: 220px;
}

.subject-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.subject-title {
  color: var(--el-text-color-primary, #303133);
  font-size: 22px;
  font-weight: 800;
}

.subject-subtitle {
  margin-top: 4px;
  color: var(--el-text-color-secondary, #909399);
  font-size: 13px;
}

.inspection-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.inspection-row {
  display: grid;
  grid-template-columns: minmax(220px, 320px) minmax(0, 1fr);
  align-items: center;
  gap: 14px;
  min-height: 72px;
  padding: 12px 14px;
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 12px;
  background: var(--el-bg-color, #fff);
}

.row-left {
  min-width: 0;
}

.row-title {
  overflow: hidden;
  color: var(--el-text-color-primary, #303133);
  font-size: 17px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 6px;
  color: var(--el-text-color-secondary, #909399);
  font-size: 12px;
}

.save-meta {
  color: var(--el-text-color-placeholder, #a8abb2);
}

.save-meta.saved {
  color: var(--el-color-success, #67c23a);
}

.save-meta.saving {
  color: var(--el-color-primary, #409eff);
}

.save-meta.error {
  color: var(--el-color-danger, #f56c6c);
}

.row-controls {
  display: grid;
  grid-template-columns: 320px 280px;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  min-width: 0;
}

.pass-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  align-items: center;
  gap: 10px;
  width: 320px;
}

.action-btn {
  width: 100%;
  height: 42px;
  padding: 0 18px;
  border: 1px solid var(--el-border-color, #dcdfe6);
  border-radius: 8px;
  background: var(--el-bg-color, #fff);
  cursor: pointer;
  font-size: 16px;
  font-weight: 800;
}

.action-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.action-pass {
  color: #15803d;
}

.action-fail {
  color: #b91c1c;
}

.action-pass.active {
  border-color: #67c23a;
  background: #67c23a;
  color: #fff;
}

.action-fail.active {
  border-color: #f56c6c;
  background: #f56c6c;
  color: #fff;
}

.number-input {
  width: 320px;
  flex-shrink: 0;
}

.remark-input {
  width: 280px;
}

.detail-footer {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: flex-end;
  padding: 16px 0;
  margin-top: 16px;
  border-top: 1px solid var(--el-border-color-light, #e4e7ed);
}

.row-pass {
  background: var(--el-color-success-light-9, #f0f9eb);
}

.row-fail {
  background: var(--el-color-danger-light-9, #fef0f0);
}

.row-rework {
  background: var(--el-color-warning-light-9, #fdf6ec);
}

.row-recheck {
  background: var(--el-color-primary-light-9, #ecf5ff);
}

.in-range :deep(.el-input__inner) {
  border-color: var(--el-color-success, #67c23a) !important;
  background: var(--el-color-success-light-9, #f0f9eb);
}

.out-range :deep(.el-input__inner) {
  border-color: var(--el-color-danger, #f56c6c) !important;
  background: var(--el-color-danger-light-9, #fef0f0);
}

@media (max-width: 768px) {
  .detail-scroll {
    padding: 0 10px 10px;
  }

  .workbench-header {
    margin: 0 -10px 10px;
    padding: 12px;
    border-radius: 0 0 12px 12px;
  }

  .subject-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
  }

  .subject-select {
    grid-column: 1 / -1;
    min-width: 0;
    order: -1;
  }

  .subject-title {
    font-size: 19px;
  }

  .inspection-row {
    grid-template-columns: 1fr;
    align-items: stretch;
  }

  .row-controls {
    flex-direction: column;
    align-items: stretch;
    display: flex;
  }

  .pass-actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
    width: 100%;
  }

  .action-btn {
    width: 100%;
  }

  .number-input,
  .remark-input {
    width: 100%;
  }

  .return-device-select,
  .return-reason-input {
    width: 100%;
  }
}
</style>
