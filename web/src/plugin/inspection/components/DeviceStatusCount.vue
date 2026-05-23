<template>
  <button
    class="device-count-link"
    :class="countClass"
    type="button"
    @click="openDialog"
  >
    {{ displayCount }}
  </button>

  <ProductionOrderDeviceDialog
    v-model="visible"
    :order="dialogOrder"
    :batch-id="batchID"
    :filter-type="type"
    :title="dialogTitle"
    :allow-rework-actions="allowReworkActions"
    :allow-recheck-actions="allowRecheckActions"
    @changed="emit('changed')"
  />
</template>

<script setup>
  import { computed, ref } from 'vue'
  import ProductionOrderDeviceDialog from '@/plugin/inspection/components/ProductionOrderDeviceDialog.vue'

  const props = defineProps({
    row: {
      type: Object,
      required: true
    },
    type: {
      type: String,
      default: 'all'
    },
    count: {
      type: [Number, String],
      default: 0
    },
    batchId: {
      type: [Number, String],
      default: ''
    },
    allowReworkActions: {
      type: Boolean,
      default: false
    },
    allowRecheckActions: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['changed'])

  const visible = ref(false)

  const displayCount = computed(() => Number(props.count || 0))
  const batchID = computed(() => props.batchId || '')
  const dialogOrder = computed(() => ({
    ...props.row,
    ID:
      props.row?.productionOrderID || props.row?.productionOrderID === 0
        ? props.row.productionOrderID
        : props.row?.ID
  }))

  const typeConfig = computed(() => {
    const map = {
      all: { label: '全部设备', className: '' },
      pending: { label: '待测设备', className: 'count-pending' },
      pass: { label: '合格设备', className: 'count-pass' },
      fail: { label: '不合格设备', className: 'count-fail' },
      rework: { label: '返工设备', className: 'count-return' },
      recheck: { label: '待复检设备', className: 'count-recheck' },
      abnormal: { label: '异常设备', className: 'count-abnormal' }
    }
    return map[props.type] || map.all
  })

  const countClass = computed(() => typeConfig.value.className)
  const dialogTitle = computed(() => {
    const mo = props.row?.moNumber || '-'
    const batch = props.row?.batchNumber ? ` / ${props.row.batchNumber}` : ''
    return `${mo}${batch} - ${typeConfig.value.label}`
  })

  const openDialog = () => {
    visible.value = true
  }
</script>

<style scoped>
  .device-count-link {
    appearance: none;
    min-width: 28px;
    padding: 2px 6px;
    border: 0;
    border-radius: 6px;
    background: transparent;
    color: var(--el-color-primary, #409eff);
    cursor: pointer;
    font: inherit;
    font-weight: 700;
    line-height: 1.4;
  }

  .device-count-link:hover {
    background: var(--el-fill-color-light, #f5f7fa);
    text-decoration: underline;
  }

  .count-pass {
    color: #16a34a;
  }

  .count-pending {
    color: #d97706;
  }

  .count-fail {
    color: #dc2626;
  }

  .count-return {
    color: #d97706;
  }

  .count-recheck {
    color: #2563eb;
  }

  .count-abnormal {
    color: #9333ea;
  }
</style>
