<template>
  <el-drawer v-model="visible" :title="title" size="520px" destroy-on-close>
    <template v-if="subject">
      <div class="log-subject">{{ subject }}</div>
      <el-timeline v-if="visibleLogs.length">
        <el-timeline-item
          v-for="log in visibleLogs"
          :key="`${log.scope}-${log.ID}`"
          :timestamp="formatDate(log.CreatedAt)"
          placement="top"
        >
          <div class="log-card">
            <div>
              <el-tag :type="log.scope === 'batch' ? 'primary' : 'success'" size="small">
                {{ log.scopeLabel }}
              </el-tag>
              <span class="log-action">{{ log.title || log.action || '-' }}</span>
            </div>
            <div>
              <span class="log-current-status">当前状态：</span>
              <el-tag :type="statusTagType(log)" size="small">
                {{ statusLabel(log, log.toStatus) }}
              </el-tag>
            </div>
            <div v-if="log.deviceSN" class="log-line">设备：{{ log.deviceSN }}</div>
            <div v-if="displayReason(log)" class="log-line">备注：{{ displayReason(log) }}</div>
            <div v-if="log.operatorName" class="log-operator">操作人：{{ log.operatorName }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else :description="`暂无${title}`" />
    </template>
  </el-drawer>
</template>

<script setup>
import { computed } from 'vue'
import { formatDate } from '@/utils/format'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '流转日志'
  },
  subject: {
    type: String,
    default: ''
  },
  logs: {
    type: Array,
    default: () => []
  },
  mode: {
    type: String,
    default: 'flow'
  }
})

const emit = defineEmits(['update:modelValue'])

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const batchStatusLabel = (value) =>
  ({ 0: '未派检', 1: '待检测接收', 2: '检测中', 3: '待确认', 4: '已完成' }[Number(value)] || value)

const batchStatusTagType = (value) =>
  ({ 0: 'info', 1: 'warning', 2: 'primary', 3: 'warning', 4: 'success' }[Number(value)] || 'info')

const deviceStatusLabel = (value) =>
  ({
    pending: '待检测设备',
    pass: '合格',
    fail: '不合格',
    returned: '待生产接收',
    rework: '返工中',
    pending_recheck: '待复检',
    rechecking: '复检中'
  }[value] ||
  value ||
  '-')

const deviceStatusTagType = (value) =>
  ({
    pending: 'info',
    pass: 'success',
    fail: 'danger',
    returned: 'warning',
    rework: 'warning',
    pending_recheck: 'primary',
    rechecking: 'warning'
  }[value] || 'info')

const statusLabel = (log, value) => {
  if (log.scope === 'batch') return batchStatusLabel(value)
  return deviceStatusLabel(value)
}

const statusTagType = (log) => {
  if (log.scope === 'batch') return batchStatusTagType(log.toStatus)
  return deviceStatusTagType(log.toStatus)
}

const deviceBatchPrefixActions = new Set([
  '生产分批',
  '生产提交检测接收',
  '检测接收并开始检测'
])

const visibleLogs = computed(() => {
  if (props.mode !== 'device') return props.logs
  return props.logs.filter((log) => {
    if (log.scope !== 'batch') return true
    const action = log.action || log.title || ''
    return deviceBatchPrefixActions.has(action)
  })
})

const systemReasons = new Set([
  '',
  '保存检测结果',
  '保存单项检测结果',
  '开始复检',
  '完成复检',
  '设备状态变更',
  '生产确认接收返工',
  '生产确认返工完成',
  '批次流转',
  '历史数据无原始派检日志，系统自动补显',
  '历史数据无原始接收日志，系统自动补显',
  '历史数据无原始待确认日志，系统自动补显',
  '历史数据无原始完成日志，系统自动补显'
])

const displayReason = (log) => {
  const reason = String(log.reason || '').trim()
  if (!reason || systemReasons.has(reason)) return ''
  if (reason === log.title || reason === log.action) return ''
  return reason
}
</script>

<style scoped>
.log-subject {
  margin-bottom: 16px;
  color: var(--el-text-color-primary, #303133);
  font-size: 16px;
  font-weight: 800;
}

.log-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px 12px;
  border: 1px solid var(--el-border-color-light, #e4e7ed);
  border-radius: 8px;
  background: var(--el-fill-color-lighter, #fafafa);
}

.log-action {
  margin-left: 8px;
  color: var(--el-text-color-primary, #303133);
  font-weight: 700;
}

.log-current-status,
.log-line {
  color: var(--el-text-color-primary, #303133);
  font-size: 13px;
}

.log-operator {
  color: var(--el-text-color-secondary, #909399);
  font-size: 12px;
}
</style>
