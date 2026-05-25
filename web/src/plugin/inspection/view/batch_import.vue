<template>
  <div class="batch-page">
    <div class="batch-header">
      <div>
        <div class="batch-title">批量导入生产设备</div>
        <div class="batch-subtitle">
          粘贴序列号列表，一键批量提交到同一个生产号下。适合同一型号仅SN不同的设备。
        </div>
      </div>
      <div class="header-actions">
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          批量提交（{{ snList.length }} 台）
        </el-button>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :md="12">
        <el-card class="batch-card" shadow="never">
          <template #header>
            <div class="card-title">基本信息（所有设备共用）</div>
          </template>
          <el-form
            ref="formRef"
            :model="form"
            :rules="formRules"
            label-width="100px"
            size="default"
          >
            <el-form-item label="MO号" prop="moNumber">
              <el-input v-model="form.moNumber" placeholder="生产号，如 MO20250525001" clearable />
            </el-form-item>
            <el-form-item label="批次号">
              <el-input v-model="form.batchNumber" placeholder="可选，留空则不绑定批次" clearable />
            </el-form-item>
            <el-form-item label="业务类型">
              <el-select v-model="form.instrumentCategory" class="full-width" clearable placeholder="可选">
                <el-option label="线上" value="online" />
                <el-option label="线下" value="offline" />
                <el-option label="外贸" value="foreign_trade" />
                <el-option label="定制款" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item label="型号" prop="model">
              <el-input v-model="form.model" placeholder="如 G3X" clearable />
            </el-form-item>
            <el-form-item label="PN码">
              <el-input v-model="form.pnCode" placeholder="产品PN码" clearable />
            </el-form-item>
            <el-form-item label="固件版本">
              <el-input v-model="form.firmwareVersion" placeholder="如 0.01.250819" clearable />
            </el-form-item>
            <el-form-item label="主板固件">
              <el-input v-model="form.mainboardFirmwareVersion" placeholder="如 UM981-18383" clearable />
            </el-form-item>
            <el-form-item label="时间注册码">
              <el-input v-model="form.timeLicense" placeholder="时间码" clearable />
            </el-form-item>
            <el-form-item label="围栏注册码">
              <el-input v-model="form.regionLicense" placeholder="围栏码" clearable />
            </el-form-item>
            <el-form-item label="Ntrip码">
              <el-input v-model="form.ntripCode" placeholder="Ntrip码" clearable />
            </el-form-item>
            <el-form-item label="GETALL模板">
              <el-input
                v-model="form.deviceInfoTemplate"
                type="textarea"
                :rows="5"
                placeholder="可选，粘贴一台设备的GETALL JSON作为模板，其中的SN会自动替换为每台设备的序列号"
              />
              <div class="form-tip">SN占位符用 {SN} 标记，提交时自动替换。留空则自动生成。</div>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="batch-card" shadow="never">
          <template #header>
            <div class="card-title">
              序列号列表
              <el-tag size="small" type="info" class="sn-count">{{ snList.length }} 台</el-tag>
            </div>
          </template>
          <el-input
            v-model="snText"
            type="textarea"
            :rows="18"
            placeholder="每行一个序列号，粘贴到这里&#10;&#10;例如：&#10;AZ16651001&#10;AZ16651002&#10;AZ16651003&#10;..."
            class="sn-textarea"
          />
          <div class="sn-actions">
            <el-button size="small" @click="handleDeduplicate">去重</el-button>
            <el-button size="small" @click="handleSort">排序</el-button>
            <el-button size="small" @click="handleClear">清空</el-button>
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              accept=".txt,.csv"
              @change="handleFileImport"
              class="inline-upload"
            >
              <el-button size="small" type="primary" plain>导入文件</el-button>
            </el-upload>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card v-if="snList.length > 0" class="batch-card preview-card" shadow="never">
      <template #header>
        <div class="card-title">提交预览</div>
      </template>
      <el-table :data="previewData" border size="small" max-height="400">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="sn" label="SN" min-width="160" />
        <el-table-column prop="model" label="型号" width="100" />
        <el-table-column prop="pnCode" label="PN码" min-width="150" />
        <el-table-column prop="firmwareVersion" label="固件版本" width="120" />
        <el-table-column prop="mainboardFirmwareVersion" label="主板固件" width="120" />
        <el-table-column label="时间码" width="100">
          <template #default="{ row }">
            <span v-if="row.timeLicense">{{ row.timeLicense }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="围栏码" width="100">
          <template #default="{ row }">
            <span v-if="row.regionLicense">{{ row.regionLicense }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="Ntrip码" width="110">
          <template #default="{ row }">
            <span v-if="row.ntripCode">{{ row.ntripCode }}</span>
            <span v-else class="text-muted">-</span>
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
const snText = ref('')
const formRef = ref(null)

const form = reactive({
  moNumber: '',
  batchNumber: '',
  instrumentCategory: 'online',
  model: '',
  pnCode: '',
  firmwareVersion: '',
  mainboardFirmwareVersion: '',
  timeLicense: '',
  regionLicense: '',
  ntripCode: '',
  deviceInfoTemplate: ''
})

const formRules = {
  moNumber: [{ required: true, message: '请输入MO号（生产号）', trigger: 'blur' }],
  model: [{ required: true, message: '请输入型号', trigger: 'blur' }],
  instrumentCategory: []
}

const snList = computed(() => {
  return snText.value
    .split(/[\n,;，；]+/)
    .map((s) => s.trim())
    .filter(Boolean)
})

const buildDeviceInfo = (sn) => {
  const tpl = form.deviceInfoTemplate.trim()
  if (tpl) {
    try {
      const replaced = tpl.replace(/\{SN\}/g, sn)
      JSON.parse(replaced)
      return replaced
    } catch {
      return tpl.replace(/\{SN\}/g, sn)
    }
  }
  return JSON.stringify({
    sn,
    model: form.model,
    pnCode: form.pnCode,
    firmwareVersion: form.firmwareVersion,
    mainboardFirmwareVersion: form.mainboardFirmwareVersion,
    device: {
      fullType: form.model,
      model: form.model,
      pn: form.pnCode,
      firmwareVersion: form.firmwareVersion,
      mainboardFirmwareVersion: form.mainboardFirmwareVersion
    }
  })
}

const previewData = computed(() =>
  snList.value.map((sn) => ({
    sn,
    model: form.model,
    pnCode: form.pnCode,
    firmwareVersion: form.firmwareVersion,
    mainboardFirmwareVersion: form.mainboardFirmwareVersion,
    timeLicense: form.timeLicense,
    regionLicense: form.regionLicense,
    ntripCode: form.ntripCode
  }))
)

const handleDeduplicate = () => {
  const unique = [...new Set(snList.value)]
  snText.value = unique.join('\n')
  ElMessage.success(`已去重，剩余 ${unique.length} 条`)
}

const handleSort = () => {
  const sorted = [...snList.value].sort((a, b) => a.localeCompare(b, undefined, { numeric: true }))
  snText.value = sorted.join('\n')
  ElMessage.success('已排序')
}

const handleClear = () => {
  snText.value = ''
}

const handleFileImport = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target.result
    snText.value = snText.value ? snText.value + '\n' + content : content
    ElMessage.success(`已导入文件: ${file.name}`)
  }
  reader.readAsText(file.raw)
}

const handleSubmit = async () => {
  if (snList.value.length === 0) {
    ElMessage.warning('请输入至少一个序列号')
    return
  }

  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定提交 ${snList.value.length} 台设备到生产号 ${form.moNumber}？`,
      '批量提交',
      { type: 'info', confirmButtonText: '提交' }
    )
  } catch {
    return
  }

  submitting.value = true
  try {
    const devices = snList.value.map((sn) => ({
      sn,
      model: form.model,
      pnCode: form.pnCode,
      firmwareVersion: form.firmwareVersion,
      mainboardFirmwareVersion: form.mainboardFirmwareVersion,
      timeLicense: form.timeLicense,
      regionLicense: form.regionLicense,
      ntripCode: form.ntripCode,
      deviceInfo: buildDeviceInfo(sn)
    }))

    const res = await submitDeviceData({
      moNumber: form.moNumber,
      batchNumber: form.batchNumber,
      instrumentCategory: form.instrumentCategory,
      devices
    })

    if (res.code === 0) {
      ElMessage.success(`成功提交 ${snList.value.length} 台设备`)
      snText.value = ''
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.batch-page {
  padding: 16px;
  background: var(--el-bg-color, #fff);
}

.batch-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.batch-title {
  font-size: 24px;
  font-weight: 800;
  color: var(--el-text-color-primary, #303133);
}

.batch-subtitle {
  margin-top: 6px;
  font-size: 13px;
  color: var(--el-text-color-secondary, #909399);
}

.header-actions {
  flex-shrink: 0;
}

.batch-card {
  margin-bottom: 16px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  color: var(--el-text-color-primary, #303133);
}

.sn-count {
  margin-left: 4px;
}

.sn-textarea :deep(.el-textarea__inner) {
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.8;
}

.sn-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
  align-items: center;
}

.inline-upload {
  display: inline-block;
}

.full-width {
  width: 100%;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  margin-top: 4px;
}

.text-muted {
  color: var(--el-text-color-disabled, #c0c4cc);
}

.preview-card {
  margin-top: 0;
}
</style>
