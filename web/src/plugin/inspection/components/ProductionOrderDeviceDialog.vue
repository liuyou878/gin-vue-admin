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
            v-if="scope.row.status === 'rework'"
            type="warning"
            link
            size="small"
            @click="handleConfirmRework(scope.row)"
          >
            返工完成
          </el-button>
          <el-button type="primary" link size="small" @click="openStatusLogs(scope.row)">
            日志
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

    <el-drawer v-model="logVisible" title="设备状态日志" size="520px" destroy-on-close>
      <template v-if="logDevice">
        <div class="log-device-title">{{ logDevice.sn || '-' }}</div>
        <el-timeline v-if="statusLogs.length">
          <el-timeline-item
            v-for="log in statusLogs"
            :key="log.ID"
            :timestamp="formatDate(log.CreatedAt)"
            placement="top"
          >
            <div class="log-card">
              <div>
                <el-tag :type="deviceStatusTagType(log.fromStatus)" size="small">
                  {{ deviceStatusLabel(log.fromStatus) }}
                </el-tag>
                <span class="log-arrow">→</span>
                <el-tag :type="deviceStatusTagType(log.toStatus)" size="small">
                  {{ deviceStatusLabel(log.toStatus) }}
                </el-tag>
              </div>
              <div class="log-reason">{{ log.reason || '-' }}</div>
              <div class="log-operator">操作人：{{ log.operatorName || '-' }}</div>
            </div>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-else description="暂无状态日志" />
      </template>
    </el-drawer>
  </el-dialog>
</template>

<script setup>
  import { computed, reactive, ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import {
    confirmReworkDone,
    getDeviceStatusLogs,
    getSubmittedDeviceList
  } from '@/plugin/inspection/api/production_order'

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
  const statusLogs = ref([])

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
      rework: { label: '返工中设备', params: { status: 'rework' } },
      recheck: { label: '待复检设备', params: { statuses: 'pending_recheck,rechecking' } },
      abnormal: { label: '异常设备', params: { statuses: 'fail,rework,pending_recheck,rechecking' } }
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
      pending: '待检测',
      pass: '合格',
      fail: '不合格',
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
      rework: 'warning',
      pending_recheck: 'primary',
      rechecking: 'warning'
    }[value] || 'info')

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

  const openStatusLogs = async (row) => {
    logDevice.value = row
    statusLogs.value = []
    logVisible.value = true
    const res = await getDeviceStatusLogs({ deviceID: row.ID })
    if (res.code === 0) {
      statusLogs.value = res.data || []
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

  .log-arrow {
    margin: 0 8px;
    color: var(--el-text-color-secondary, #909399);
  }

  .log-reason,
  .log-operator {
    color: var(--el-text-color-secondary, #909399);
    font-size: 12px;
  }
</style>
