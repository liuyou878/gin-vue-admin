<template>
  <div class="mock-page">
    <div class="mock-header">
      <div>
        <div class="mock-title">模拟生产工具提交</div>
        <div class="mock-subtitle">
          自动生成一组生产工具数据，提交后会进入生产订单和设备列表。
        </div>
      </div>
      <div class="header-actions">
        <el-button @click="generateMockData">重新生成</el-button>
        <el-input-number
          v-model="appendCount"
          :min="1"
          :max="100"
          controls-position="right"
          class="append-count"
        />
        <el-button type="success" plain @click="appendMockDevices(appendCount)">继续添加</el-button>
        <el-button type="primary" :loading="submitting" @click="submitMockData">
          提交模拟数据
        </el-button>
      </div>
    </div>

    <el-alert
      class="mock-alert"
      title="这里复用真实生产工具提交接口 /productionOrder/submitDeviceData，不新增后端模拟接口。"
      type="info"
      :closable="false"
      show-icon
    />

    <el-card class="mock-card" shadow="never">
      <template #header>
        <div class="card-title">本次模拟批次</div>
      </template>
      <el-descriptions :column="3" border size="small">
        <el-descriptions-item label="MO号">{{ form.moNumber }}</el-descriptions-item>
        <el-descriptions-item label="分批方式">
          {{ autoBindBatch ? '提交时绑定批次' : '只提交设备，后续手动分批' }}
        </el-descriptions-item>
        <el-descriptions-item label="批次号">
          {{ autoBindBatch ? form.batchNumber : '暂不绑定' }}
        </el-descriptions-item>
        <el-descriptions-item label="业务类型">{{ categoryLabel(form.instrumentCategory) }}</el-descriptions-item>
        <el-descriptions-item label="型号">{{ previewDevice.model }}</el-descriptions-item>
        <el-descriptions-item label="PN码">{{ previewDevice.pnCode }}</el-descriptions-item>
        <el-descriptions-item label="设备数量">{{ form.devices.length }}</el-descriptions-item>
        <el-descriptions-item label="固件版本">{{ previewDevice.firmwareVersion }}</el-descriptions-item>
        <el-descriptions-item label="主板固件">{{ previewDevice.mainboardFirmwareVersion }}</el-descriptions-item>
        <el-descriptions-item label="SN范围">
          {{ form.devices[0]?.sn || '-' }} ~ {{ form.devices[form.devices.length - 1]?.sn || '-' }}
        </el-descriptions-item>
      </el-descriptions>
      <div class="batch-options">
        <el-switch
          v-model="autoBindBatch"
          active-text="提交时绑定批次"
          inactive-text="后续手动分批"
        />
        <el-input
          v-if="autoBindBatch"
          v-model="form.batchNumber"
          class="batch-input"
          placeholder="请输入批次号"
        />
      </div>
    </el-card>

    <el-card class="mock-card" shadow="never">
      <template #header>
        <div class="card-title">设备数据预览</div>
      </template>
      <el-table :data="form.devices" border size="small" max-height="520">
        <el-table-column type="index" label="序号" width="70" />
        <el-table-column prop="sn" label="SN" min-width="150" />
        <el-table-column prop="model" label="型号" width="100" />
        <el-table-column prop="pnCode" label="PN码" min-width="150" />
        <el-table-column prop="firmwareVersion" label="固件版本" min-width="130" />
        <el-table-column prop="mainboardFirmwareVersion" label="主板固件" min-width="130" />
        <el-table-column prop="timeLicense" label="时间码" min-width="110" />
        <el-table-column prop="regionLicense" label="围栏码" min-width="110" />
        <el-table-column prop="ntripCode" label="Ntrip码" min-width="120" />
        <el-table-column label="操作" width="110" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              size="small"
              @click="regenerateDeviceSN(scope.$index)"
            >
              重生成SN
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
  import { computed, reactive, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { submitDeviceData } from '@/plugin/inspection/api/production_order'

  const submitting = ref(false)
  const appendCount = ref(5)
  const autoBindBatch = ref(false)
  const form = reactive({
    moNumber: '',
    batchNumber: '',
    instrumentCategory: 'online',
    devices: []
  })

  const previewDevice = computed(() => form.devices[0] || {})

  const categoryLabel = (value) =>
    ({
      online: '线上',
      offline: '线下',
      foreign_trade: '外贸',
      custom: '定制款'
    }[value] || value)

  const pad = (value, length = 2) => String(value).padStart(length, '0')

  const formatDate = (date = new Date()) => {
    const y = date.getFullYear()
    const m = pad(date.getMonth() + 1)
    const d = pad(date.getDate())
    return `${y}${m}${d}`
  }

  const formatTimeCode = (date = new Date()) => {
    const h = pad(date.getHours())
    const m = pad(date.getMinutes())
    const s = pad(date.getSeconds())
    return `${h}${m}${s}`
  }

  const buildDeviceInfo = (device) =>
    JSON.stringify({
      sn: device.sn,
      model: device.model,
      pnCode: device.pnCode,
      firmwareVersion: device.firmwareVersion,
      mainboardFirmwareVersion: device.mainboardFirmwareVersion,
      device: {
        fullType: device.model,
        model: device.model,
        pn: device.pnCode,
        firmwareVersion: device.firmwareVersion,
        mainboardFirmwareVersion: device.mainboardFirmwareVersion
      },
      mock: true,
      generatedAt: new Date().toISOString()
    })

  const buildMockSN = (snBase, index) =>
    `G3X${String(snBase + index).padStart(6, '0')}`

  const generateMockData = () => {
    const now = new Date()
    const dateText = formatDate(now)
    const timeText = formatTimeCode(now)
    const moNumber = `MO${dateText}${timeText}`
    const batchNumber = `${moNumber}-${dateText}-01`
    const model = 'G3X'
    const pnCode = 'G3X-MOCK-PN'
    const firmwareVersion = '0.01.250819'
    const mainboardFirmwareVersion = 'UM981-18383'
    const snBase = Number(`${dateText.slice(2)}${timeText}`.slice(-6))

    form.moNumber = moNumber
    form.batchNumber = batchNumber
    form.instrumentCategory = 'online'
    form.devices = buildMockDevices({
      count: 5,
      startIndex: 0,
      snBase,
      dateText,
      timeText,
      model,
      pnCode,
      firmwareVersion,
      mainboardFirmwareVersion
    })
  }

  const buildMockDevices = ({
    count,
    startIndex,
    snBase,
    dateText,
    timeText,
    model,
    pnCode,
    firmwareVersion,
    mainboardFirmwareVersion
  }) =>
    Array.from({ length: count }, (_, offset) => {
      const index = startIndex + offset
      const sn = buildMockSN(snBase, index)
      const device = {
        sn,
        model,
        pnCode,
        firmwareVersion,
        mainboardFirmwareVersion,
        timeLicense: `T${dateText.slice(2)}`,
        regionLicense: `R${timeText}`,
        ntripCode: `NTRIP${pad(index + 1, 3)}`
      }
      return {
        ...device,
        deviceInfo: buildDeviceInfo(device)
      }
    })

  const appendMockDevices = (count = 5) => {
    const first = form.devices[0]
    if (!first) {
      generateMockData()
      return
    }
    const now = new Date()
    const dateText = formatDate(now)
    const timeText = formatTimeCode(now)
    const firstSnNumber = Number(String(first.sn || '').match(/(\d+)$/)?.[1]) || 0
    const snBase = firstSnNumber
    const startIndex = form.devices.length
    form.devices.push(
      ...buildMockDevices({
        count,
        startIndex,
        snBase,
        dateText,
        timeText,
        model: first.model,
        pnCode: first.pnCode,
        firmwareVersion: first.firmwareVersion,
        mainboardFirmwareVersion: first.mainboardFirmwareVersion
      })
    )
  }

  const nextAvailableSN = () => {
    const used = new Set(form.devices.map((device) => device.sn))
    const now = new Date()
    const dateText = formatDate(now)
    const timeText = formatTimeCode(now)
    const snBase = Number(`${dateText.slice(2)}${timeText}`.slice(-6))
    for (let index = 0; index < 1000; index++) {
      const sn = buildMockSN(snBase, index)
      if (!used.has(sn)) return sn
    }
    return `G3X${Date.now().toString().slice(-9)}`
  }

  const regenerateDeviceSN = (index) => {
    const device = form.devices[index]
    if (!device) return
    device.sn = nextAvailableSN()
    device.deviceInfo = buildDeviceInfo(device)
    ElMessage.success('已重生成该设备SN')
  }

  const submitMockData = async () => {
    if (!form.devices.length) {
      ElMessage.warning('请先生成模拟设备')
      return
    }
    try {
      await ElMessageBox.confirm(
        `确定提交 ${form.devices.length} 台模拟设备到生产订单 ${form.moNumber}？${
          autoBindBatch.value ? `提交时会绑定批次 ${form.batchNumber}` : '提交后需要到生产订单里手动分批'
        }`,
        '提交模拟数据',
        {
          type: 'info',
          confirmButtonText: '提交'
        }
      )
    } catch {
      return
    }

    submitting.value = true
    try {
      const res = await submitDeviceData({
        moNumber: form.moNumber,
        batchNumber: autoBindBatch.value ? form.batchNumber : '',
        instrumentCategory: form.instrumentCategory,
        devices: form.devices
      })
      if (res.code !== 0) return
      ElMessage.success('模拟生产数据提交成功')
      generateMockData()
    } finally {
      submitting.value = false
    }
  }

  generateMockData()
</script>

<style scoped>
  .mock-page {
    padding: 16px;
    background: var(--el-bg-color, #fff);
  }

  .mock-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 12px;
  }

  .mock-title {
    color: var(--el-text-color-primary, #303133);
    font-size: 24px;
    font-weight: 800;
  }

  .mock-subtitle {
    margin-top: 6px;
    color: var(--el-text-color-secondary, #909399);
    font-size: 13px;
  }

  .header-actions {
    display: flex;
    gap: 8px;
    flex-shrink: 0;
  }

  .append-count {
    width: 110px;
  }

  .mock-alert,
  .mock-card {
    margin-bottom: 14px;
  }

  .batch-options {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 12px;
  }

  .batch-input {
    width: 360px;
  }

  .card-title {
    color: var(--el-text-color-primary, #303133);
    font-weight: 700;
  }

  @media (max-width: 768px) {
    .mock-header {
      flex-direction: column;
    }

    .header-actions {
      width: 100%;
    }

    .header-actions :deep(.el-button) {
      flex: 1;
      margin-left: 0;
    }

    .append-count {
      width: 100%;
    }

    .batch-options {
      align-items: stretch;
      flex-direction: column;
    }

    .batch-input {
      width: 100%;
    }
  }
</style>
