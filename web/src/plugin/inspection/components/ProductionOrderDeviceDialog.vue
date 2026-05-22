<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="980px"
    destroy-on-close
    @open="loadDevices"
  >
    <div class="device-dialog-toolbar">
      <el-input
        v-model="searchInfo.sn"
        placeholder="按SN搜索"
        clearable
        size="small"
        style="width: 220px"
        @keyup.enter="reloadDevices"
        @clear="reloadDevices"
      />
      <el-button size="small" type="primary" @click="reloadDevices">查询</el-button>
      <el-button size="small" @click="resetSearch">重置</el-button>
      <span class="device-dialog-tip">共 {{ total }} 台</span>
    </div>

    <el-table :data="tableData" border size="small" v-loading="loading" max-height="520">
      <el-table-column prop="sn" label="SN" min-width="140" />
      <el-table-column prop="batchNumber" label="批次号" min-width="160">
        <template #default="scope">{{ scope.row.batchNumber || '-' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="deviceStatusTagType(scope.row.status)" size="small">
            {{ deviceStatusLabel(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="model" label="型号" width="110" />
      <el-table-column prop="pnCode" label="PN码" min-width="150" />
      <el-table-column prop="firmwareVersion" label="固件版本" min-width="130" />
      <el-table-column prop="mainboardFirmwareVersion" label="主板固件" min-width="130" />
      <el-table-column prop="returnReason" label="打回原因" min-width="160" show-overflow-tooltip>
        <template #default="scope">{{ scope.row.returnReason || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="170" fixed="right">
        <template #default="scope">
          <el-button
            v-if="scope.row.status === 'returned'"
            type="warning"
            link
            size="small"
            @click="handleConfirmReworkReceived(scope.row)"
          >
            接收返工
          </el-button>
          <el-button
            v-if="scope.row.status === 'rework'"
            type="warning"
            link
            size="small"
            @click="handleConfirmRework(scope.row)"
          >
            返工完成
          </el-button>
          <el-button type="primary" link size="small" @click="openFlowLogs(scope.row)">
            设备日志
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="gva-pagination device-dialog-pagination">
      <el-pagination
        v-model:current-page="searchInfo.page"
        v-model:page-size="searchInfo.pageSize"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        background
        small
        @size-change="loadDevices"
        @current-change="loadDevices"
      />
    </div>

    <el-drawer v-model="logVisible" title="设备日志" size="520px" destroy-on-close>
      <template v-if="logDevice">
        <div class="log-device-title">{{ logDevice.sn || '-' }}</div>
        <el-timeline v-if="flowLogs.length">
          <el-timeline-item
            v-for="log in flowLogs"
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
                <el-tag :type="flowStatusTagType(log)" size="small">
                  {{ flowStatusLabel(log, log.toStatus) }}
                </el-tag>
              </div>
              <div v-if="log.reason" class="log-reason">备注：{{ log.reason }}</div>
              <div v-if="log.operatorName" class="log-operator">操作人：{{ log.operatorName }}</div>
            </div>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-else description="暂无设备日志" />
      </template>
    </el-drawer>
  </el-dialog>
</template>

<script setup>
  import { computed, reactive, ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import {
    confirmReworkReceived,
    confirmReworkDone,
    getSubmittedDeviceList
  } from '@/plugin/inspection/api/production_order'
  import { getFlowLogs } from '@/plugin/inspection/api/work_order'

  const props = defineProps({
    modelValue: {
      type: Boolean,
      default: false
    },
    order: {
      type: Object,
      default: null
    },
    filterType: {
      type: String,
      default: 'all'
    },
    title: {
      type: String,
      default: ''
    }
  })

  const emit = defineEmits(['update:modelValue', 'changed'])

  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const logVisible = ref(false)
  const logDevice = ref(null)
  const flowLogs = ref([])

  const searchInfo = reactive({
    sn: '',
    page: 1,
    pageSize: 30
  })

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const filterConfig = computed(() => {
    const map = {
      all: { label: '全部设备', params: {} },
      pass: { label: '合格设备', params: { status: 'pass' } },
      fail: { label: '不合格设备', params: { status: 'fail' } },
      rework: { label: '返工设备', params: { statuses: 'returned,rework' } },
      recheck: { label: '待复检设备', params: { statuses: 'pending_recheck,rechecking' } },
      abnormal: { label: '异常设备', params: { statuses: 'fail,returned,rework,pending_recheck,rechecking' } }
    }
    return map[props.filterType] || map.all
  })

  const dialogTitle = computed(() => {
    if (props.title) return props.title
    const mo = props.order?.moNumber || '-'
    return `${mo} - ${filterConfig.value.label}`
  })

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

  const batchStatusLabel = (value) =>
    ({ 0: '未派检', 1: '待检测接收', 2: '检测中', 3: '已完成' }[Number(value)] || value)

  const batchStatusTagType = (value) =>
    ({ 0: 'info', 1: 'warning', 2: 'primary', 3: 'success' }[Number(value)] || 'info')

  const flowStatusLabel = (log, value) => {
    if (log.scope === 'batch') return batchStatusLabel(value)
    return deviceStatusLabel(value)
  }

  const flowStatusTagType = (log) => {
    if (log.scope === 'batch') return batchStatusTagType(log.toStatus)
    return deviceStatusTagType(log.toStatus)
  }

  const buildParams = () => ({
    productionOrderID: props.order?.ID,
    sn: searchInfo.sn,
    page: searchInfo.page,
    pageSize: searchInfo.pageSize,
    ...filterConfig.value.params
  })

  const loadDevices = async () => {
    if (!props.order?.ID) return
    loading.value = true
    try {
      const res = await getSubmittedDeviceList(buildParams())
      if (res.code === 0) {
        tableData.value = res.data.list || []
        total.value = res.data.total || 0
      }
    } finally {
      loading.value = false
    }
  }

  const reloadDevices = () => {
    searchInfo.page = 1
    loadDevices()
  }

  const resetSearch = () => {
    searchInfo.sn = ''
    searchInfo.page = 1
    loadDevices()
  }

  const handleConfirmRework = async (row) => {
    try {
      await ElMessageBox.confirm(
        `确认 ${row.sn} 已返工完成，并提交给检测复检？`,
        '确认返工完成',
        { type: 'warning', confirmButtonText: '确认完成' }
      )
    } catch {
      return
    }

    const res = await confirmReworkDone({ deviceID: row.ID })
    if (res.code !== 0) return
    ElMessage.success('已进入待复检')
    await loadDevices()
    emit('changed')
  }

  const handleConfirmReworkReceived = async (row) => {
    try {
      await ElMessageBox.confirm(
        `确认已接收 ${row.sn}，并开始返工？`,
        '确认接收返工',
        { type: 'warning', confirmButtonText: '确认接收' }
      )
    } catch {
      return
    }

    const res = await confirmReworkReceived({ deviceID: row.ID })
    if (res.code !== 0) return
    ElMessage.success('已进入返工中')
    await loadDevices()
    emit('changed')
  }

  const openFlowLogs = async (row) => {
    logDevice.value = row
    flowLogs.value = []
    logVisible.value = true
    const res = await getFlowLogs({ batchID: row.batchID, deviceID: row.ID })
    if (res.code === 0) {
      flowLogs.value = res.data || []
    }
  }

  watch(
    () => [props.order?.ID, props.filterType],
    () => {
      searchInfo.sn = ''
      searchInfo.page = 1
      if (visible.value) {
        loadDevices()
      }
    }
  )
</script>

<style scoped>
  .device-dialog-toolbar {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }

  .device-dialog-tip {
    color: var(--el-text-color-secondary, #909399);
    font-size: 12px;
  }

  .device-dialog-pagination {
    margin-top: 12px;
  }

  .log-device-title {
    margin-bottom: 12px;
    font-size: 16px;
    font-weight: 700;
  }

  .log-card {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .log-action {
    margin-left: 8px;
    font-weight: 600;
  }

  .log-current-status,
  .log-reason,
  .log-operator {
    color: var(--el-text-color-secondary, #909399);
    font-size: 12px;
  }
</style>
