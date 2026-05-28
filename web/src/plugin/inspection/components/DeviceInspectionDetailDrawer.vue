<template>
  <el-drawer
    v-model="visible"
    title="设备检测明细"
    size="720px"
    append-to-body
    destroy-on-close
  >
    <div v-loading="loading">
      <template v-if="device">
        <div class="device-title">{{ device.sn || '-' }}</div>
        <el-table :data="device.results || []" border size="small">
          <el-table-column prop="itemName" label="检测项" min-width="150" />
          <el-table-column label="标准" min-width="130">
            <template #default="scope">{{
              resultStandard(scope.row)
            }}</template>
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
          <el-table-column
            prop="remark"
            label="备注"
            min-width="140"
            show-overflow-tooltip
          >
            <template #default="scope">{{ scope.row.remark || '-' }}</template>
          </el-table-column>
          <el-table-column prop="inspectorName" label="检测人" width="110">
            <template #default="scope">{{
              scope.row.inspectorName || '-'
            }}</template>
          </el-table-column>
          <el-table-column prop="inspectedAt" label="检测时间" width="160">
            <template #default="scope">{{
              formatDate(scope.row.inspectedAt) || '-'
            }}</template>
          </el-table-column>
        </el-table>
      </template>
      <el-empty v-else description="暂无检测明细" />
    </div>
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
    device: {
      type: Object,
      default: null
    },
    loading: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['update:modelValue'])

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const resultCompleted = (result) => {
    if (!result) return false
    if (result.resultType === 'number') {
      return result.numberResult !== undefined && result.numberResult !== null
    }
    if (result.resultType === 'pass_fail') {
      return result.passResult === true || result.passResult === false
    }
    return (
      (result.passResult === true || result.passResult === false) &&
      result.numberResult !== undefined &&
      result.numberResult !== null
    )
  }

  const resultFailed = (result) => {
    if (!resultCompleted(result)) return false
    const passFailed = result.passResult === false
    const numberFailed =
      result.numberResult !== undefined &&
      result.numberResult !== null &&
      ((result.minValue != null && result.numberResult < result.minValue) ||
        (result.maxValue != null && result.numberResult > result.maxValue))
    return passFailed || numberFailed
  }

  const resultLabel = (result) => {
    if (!resultCompleted(result)) return '未完成'
    if (resultFailed(result)) return '未通过'
    if (
      result.passResult === true ||
      result.minValue != null ||
      result.maxValue != null
    )
      return '通过'
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
    if (result.passResult === false) values.push('未通过')
    if (result.numberResult !== undefined && result.numberResult !== null) {
      values.push(`${result.numberResult}${result.unit || ''}`)
    }
    return values.join(' / ') || '-'
  }
</script>

<style scoped>
  .device-title {
    margin-bottom: 12px;
    color: var(--el-text-color-primary, #303133);
    font-size: 18px;
    font-weight: 700;
  }
</style>
