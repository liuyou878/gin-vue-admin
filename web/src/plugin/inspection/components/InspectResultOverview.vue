<template>
  <div class="result-overview">
    <section class="overview-header">
      <div>
        <div class="overview-title">检测结果总览</div>
        <div class="overview-subtitle">
          {{ detail.order.moNumber || '-' }} / {{ detail.order.batchNumber || '-' }}
        </div>
      </div>
      <div class="overview-actions">
        <el-tag :type="isFinalCompleted ? 'success' : 'primary'" size="large">
          {{ isFinalCompleted ? '已完成' : '检测中' }}
        </el-tag>
      </div>
    </section>

    <section class="info-grid">
      <div class="info-card">
        <span>模板</span>
        <strong>{{ detail.order.templateName || detail.order.productName || '-' }}</strong>
      </div>
      <div class="info-card">
        <span>型号</span>
        <strong>{{ detail.order.model || '-' }}</strong>
      </div>
      <div class="info-card">
        <span>固件版本</span>
        <strong>{{ detail.order.firmwareVersion || '-' }}</strong>
      </div>
      <div class="info-card">
        <span>主板固件</span>
        <strong>{{ detail.order.mainboardFirmwareVersion || '-' }}</strong>
      </div>
      <div class="info-card">
        <span>检测人</span>
        <strong>{{ detail.order.inspectorName || '-' }}</strong>
      </div>
      <div class="info-card">
        <span>检测时间</span>
        <strong>{{ formatDate(detail.order.inspectionDate) || '-' }}</strong>
      </div>
    </section>

    <section class="summary-grid">
      <div class="summary-card">
        <span>总数</span>
        <strong>{{ summary.total }}</strong>
      </div>
      <div class="summary-card pass">
        <span>合格</span>
        <strong>{{ summary.pass }}</strong>
      </div>
      <div class="summary-card fail">
        <span>不合格</span>
        <strong>{{ summary.fail }}</strong>
      </div>
      <div class="summary-card rework">
        <span>返工/待生产</span>
        <strong>{{ summary.rework }}</strong>
      </div>
      <div class="summary-card recheck">
        <span>待复检</span>
        <strong>{{ summary.recheck }}</strong>
      </div>
      <div class="summary-card rate">
        <span>合格率</span>
        <strong>{{ passRate }}</strong>
      </div>
    </section>

    <section v-if="canReturn && failDevices.length" class="return-panel">
      <div>
        <div class="return-title">不合格处理</div>
        <div class="return-desc">选择需要返回生产处理的设备，填写原因后打回。</div>
      </div>
      <el-select
        v-model="returnDeviceIDs"
        multiple
        collapse-tags
        collapse-tags-tooltip
        filterable
        placeholder="选择不合格设备"
        class="return-select"
      >
        <el-option
          v-for="device in failDevices"
          :key="device.ID"
          :label="`${device.lineNumber || '-'} . ${device.sn}`"
          :value="device.ID"
        />
      </el-select>
      <el-input v-model="returnReason" placeholder="打回原因" class="return-reason" />
      <el-button type="warning" :disabled="returnDeviceIDs.length === 0" @click="onReturnDevices">
        打回生产
      </el-button>
    </section>

    <el-table :data="deviceRows" border stripe size="small" class="result-table">
      <el-table-column prop="lineNumber" label="序号" width="70" />
      <el-table-column prop="sn" label="SN" min-width="150" />
      <el-table-column label="结果" width="110">
        <template #default="scope">
          <el-tag :type="deviceStatusTag(scope.row)">
            {{ deviceStatusLabel(scope.row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="完成项" width="90">
        <template #default="scope">{{ scope.row.doneCount }} / {{ scope.row.totalCount }}</template>
      </el-table-column>
      <el-table-column label="失败项" width="90">
        <template #default="scope">
          <span :class="{ 'fail-text': scope.row.failCount > 0 }">{{ scope.row.failCount }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="inspectorName" label="最后检测人" width="120">
        <template #default="scope">{{ scope.row.inspectorName || '-' }}</template>
      </el-table-column>
      <el-table-column prop="inspectedAt" label="最后检测时间" width="160">
        <template #default="scope">{{ formatDate(scope.row.inspectedAt) || '-' }}</template>
      </el-table-column>
      <el-table-column prop="returnReason" label="打回原因" min-width="160" show-overflow-tooltip>
        <template #default="scope">{{ scope.row.returnReason || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="210" fixed="right">
        <template #default="scope">
          <el-button type="primary" link size="small" @click="openDeviceDetail(scope.row)">
            查看明细
          </el-button>
          <el-button
            v-if="canStartRecheck(scope.row)"
            type="warning"
            link
            size="small"
            @click="handleStartRecheck(scope.row)"
          >
            开始复检
          </el-button>
          <el-button
            v-if="canContinueRecheck(scope.row)"
            type="warning"
            link
            size="small"
            @click="openInspectDetail(scope.row)"
          >
            继续复检
          </el-button>
          <el-button
            v-if="canReturn && scope.row.status === 'fail'"
            type="warning"
            link
            size="small"
            @click="returnOne(scope.row)"
          >
            打回
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-drawer v-model="detailVisible" title="设备检测明细" size="720px" destroy-on-close>
      <template v-if="currentDevice">
        <div class="device-title">{{ currentDevice.sn }}</div>
        <el-table :data="currentDevice.results" border size="small">
          <el-table-column prop="itemName" label="检测项" min-width="150" />
          <el-table-column label="标准" min-width="130">
            <template #default="scope">{{ resultStandard(scope.row) }}</template>
          </el-table-column>
          <el-table-column label="检测值" min-width="120">
            <template #default="scope">{{ resultValue(scope.row) }}</template>
          </el-table-column>
          <el-table-column label="结果" width="100">
            <template #default="scope">
              <el-tag :type="resultTag(scope.row)" size="small">
                {{ resultLabel(scope.row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="140" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.remark || '-' }}</template>
          </el-table-column>
          <el-table-column prop="inspectorName" label="检测人" width="110">
            <template #default="scope">{{ scope.row.inspectorName || '-' }}</template>
          </el-table-column>
          <el-table-column prop="inspectedAt" label="检测时间" width="160">
            <template #default="scope">{{ formatDate(scope.row.inspectedAt) || '-' }}</template>
          </el-table-column>
        </el-table>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { returnDevices as apiReturnDevices, startRecheck } from '@/plugin/inspection/api/work_order'

const props = defineProps({
  detail: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['refresh'])

const returnDeviceIDs = ref([])
const returnReason = ref('')
const detailVisible = ref(false)
const currentDevice = ref(null)

const normalizedDevices = computed(() => props.detail.devices || [])

const summary = computed(() => {
  const value = {
    total: normalizedDevices.value.length,
    pass: 0,
    fail: 0,
    rework: 0,
    recheck: 0
  }
  normalizedDevices.value.forEach((device) => {
    const status = device.status || device._status || 'pending'
    if (status === 'pass') value.pass += 1
    if (status === 'fail') value.fail += 1
    if (['returned', 'rework'].includes(status)) value.rework += 1
    if (['pending_recheck', 'rechecking'].includes(status)) value.recheck += 1
  })
  return value
})

const passRate = computed(() => {
  if (!summary.value.total) return '-'
  return `${((summary.value.pass / summary.value.total) * 100).toFixed(1)}%`
})

const isFinalCompleted = computed(() => props.detail.order?.status === 4)
const pendingStatusKeys = ['pending', 'fail', 'returned', 'rework', 'pending_recheck', 'rechecking']

const hasRecheckDevices = computed(() =>
  normalizedDevices.value.some((device) => ['pending_recheck', 'rechecking'].includes(device.status || device._status))
)

const canReturn = computed(() => props.detail.order?.status === 2 && !hasRecheckDevices.value)

const failDevices = computed(() =>
  normalizedDevices.value.filter((device) => (device.status || device._status) === 'fail')
)

const deviceRows = computed(() =>
  normalizedDevices.value.map((device, index) => {
    const doneCount = (device.results || []).filter(resultCompleted).length
    const failCount = (device.results || []).filter(resultFailed).length
    const latest = latestResult(device.results || [])
    return {
      ...device,
      lineNumber: device.lineNumber || index + 1,
      doneCount,
      failCount,
      totalCount: (device.results || []).length,
      inspectorName: latest?.inspectorName || '',
      inspectedAt: latest?.inspectedAt || null
    }
  })
)

const latestResult = (results) => {
  return results
    .filter((item) => item.inspectedAt || item.inspectorName)
    .sort((a, b) => new Date(b.inspectedAt || 0) - new Date(a.inspectedAt || 0))[0]
}

const resultCompleted = (result) => {
  if (!result) return false
  if (result.resultType === 'number') {
    return result.numberResult !== undefined && result.numberResult !== null
  }
  if (result.resultType === 'pass_fail') {
    return result.passResult === true || result.passResult === false
  }
  return (result.passResult === true || result.passResult === false) &&
    result.numberResult !== undefined &&
    result.numberResult !== null
}

const resultFailed = (result) => {
  if (!resultCompleted(result)) return false
  const passFailed = result.passResult === false
  const numberFailed = result.numberResult !== undefined &&
    result.numberResult !== null &&
    ((result.minValue != null && result.numberResult < result.minValue) ||
      (result.maxValue != null && result.numberResult > result.maxValue))
  return passFailed || numberFailed
}

const resultLabel = (result) => {
  if (!resultCompleted(result)) return '未完成'
  if (resultFailed(result)) return '不通过'
  if (result.passResult === true || result.minValue != null || result.maxValue != null) return '通过'
  return '已填写'
}

const resultTag = (result) => {
  if (!resultCompleted(result)) return 'info'
  return resultFailed(result) ? 'danger' : 'success'
}

const resultStandard = (result) => {
  const parts = []
  if (result.minValue != null || result.maxValue != null) {
    parts.push(`${result.minValue ?? '-'} ~ ${result.maxValue ?? '-'}`)
  }
  if (result.unit) parts.push(result.unit)
  return parts.join(' ') || '-'
}

const resultValue = (result) => {
  const values = []
  if (result.passResult === true) values.push('通过')
  if (result.passResult === false) values.push('不通过')
  if (result.numberResult !== undefined && result.numberResult !== null) {
    values.push(`${result.numberResult}${result.unit || ''}`)
  }
  return values.join(' / ') || '-'
}

const deviceStatusLabel = (device) =>
  ({
    pending: '待完成',
    pass: '合格',
    fail: '不合格',
    returned: '待生产接收',
    rework: '返工中',
    pending_recheck: '待复检',
    rechecking: '复检中'
  }[device.status || device._status] || device.status || '-')

const deviceStatusTag = (device) =>
  ({
    pending: 'info',
    pass: 'success',
    fail: 'danger',
    returned: 'warning',
    rework: 'warning',
    pending_recheck: 'primary',
    rechecking: 'warning'
  }[device.status || device._status] || 'info')

const openDeviceDetail = (device) => {
  currentDevice.value = device
  detailVisible.value = true
}

const canStartRecheck = (device) =>
  !isFinalCompleted.value && (device.status || device._status) === 'pending_recheck'

const canContinueRecheck = (device) =>
  !isFinalCompleted.value && (device.status || device._status) === 'rechecking'

const openInspectDetail = () => {
  window.location.hash = `/inspectDetail?batchId=${props.detail.order.ID}`
}

const handleStartRecheck = async (device) => {
  try {
    await ElMessageBox.confirm(
      `确定开始复检 ${device.sn}？开始后检测端才能录入复检结果。`,
      '开始复检',
      {
        type: 'warning',
        confirmButtonText: '开始复检'
      }
    )
  } catch {
    return
  }

  const res = await startRecheck({ ID: props.detail.order.ID, deviceID: device.ID })
  if (res.code !== 0) return
  ElMessage.success({ message: '已开始复检', duration: 1000 })
  emit('refresh')
  window.setTimeout(() => {
    ElMessage.closeAll()
    openInspectDetail()
  }, 300)
}

const returnOne = (device) => {
  returnDeviceIDs.value = [device.ID]
  onReturnDevices()
}

const onReturnDevices = async () => {
  if (!returnDeviceIDs.value.length) {
    ElMessage.warning('请先选择要打回的设备')
    return
  }
  try {
    await ElMessageBox.confirm('确定将选中的不合格设备打回生产？', '打回生产', {
      type: 'warning',
      confirmButtonText: '确认打回'
    })
  } catch {
    return
  }

  const res = await apiReturnDevices({
    batchID: props.detail.order.ID,
    deviceIDs: returnDeviceIDs.value,
    reason: returnReason.value || ''
  })
  if (res.code !== 0) return
  ElMessage.success('已打回生产')
  returnDeviceIDs.value = []
  returnReason.value = ''
  emit('refresh')
}

</script>

<style scoped>
.result-overview {
  flex: 1;
  overflow: auto;
  padding: 16px;
  background: var(--el-bg-color, #fff);
}

.overview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.overview-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.confirm-alert {
  margin-bottom: 14px;
}

.overview-title {
  color: var(--el-text-color-primary, #303133);
  font-size: 24px;
  font-weight: 800;
}

.overview-subtitle {
  margin-top: 4px;
  color: var(--el-text-color-secondary, #909399);
  font-size: 13px;
}

.info-grid,
.summary-grid {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 14px;
}

.info-card,
.summary-card {
  min-width: 0;
  padding: 12px;
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 12px;
  background: var(--el-fill-color-lighter, #fafafa);
}

.info-card span,
.summary-card span {
  display: block;
  color: var(--el-text-color-secondary, #909399);
  font-size: 12px;
}

.info-card strong,
.summary-card strong {
  display: block;
  margin-top: 6px;
  overflow: hidden;
  color: var(--el-text-color-primary, #303133);
  font-size: 18px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.summary-card.pass strong {
  color: #16a34a;
}

.summary-card.fail strong,
.fail-text {
  color: #dc2626;
}

.summary-card.rework strong {
  color: #d97706;
}

.summary-card.recheck strong,
.summary-card.rate strong {
  color: #2563eb;
}

.return-panel {
  display: grid;
  grid-template-columns: minmax(180px, 260px) minmax(220px, 1fr) minmax(220px, 1fr) auto;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
  padding: 12px;
  border: 1px solid var(--el-color-warning-light-7, #f3d19e);
  border-radius: 12px;
  background: var(--el-color-warning-light-9, #fdf6ec);
}

.return-title {
  color: var(--el-text-color-primary, #303133);
  font-weight: 700;
}

.return-desc {
  margin-top: 4px;
  color: var(--el-text-color-secondary, #909399);
  font-size: 12px;
}

.return-select,
.return-reason {
  width: 100%;
}

.result-table {
  width: 100%;
}

.device-title {
  margin-bottom: 12px;
  color: var(--el-text-color-primary, #303133);
  font-size: 18px;
  font-weight: 800;
}

@media (max-width: 1100px) {
  .info-grid,
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .return-panel {
    grid-template-columns: 1fr;
  }
}
</style>
