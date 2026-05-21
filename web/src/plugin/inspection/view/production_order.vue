<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="MO号">
          <el-input
            v-model="searchInfo.moNumber"
            placeholder="请输入"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="型号">
          <el-input
            v-model="searchInfo.model"
            placeholder="请输入"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="业务类型">
          <el-select
            v-model="searchInfo.instrumentCategory"
            placeholder="请选择"
            clearable
            size="small"
            style="width: 120px"
          >
            <el-option label="线上" value="online" />
            <el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" />
            <el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="getList"
            >查询</el-button
          >
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData" border v-loading="loading" size="small">
      <el-table-column prop="moNumber" label="MO号" min-width="140" />
      <el-table-column label="批次数" width="90">
        <template #default="scope">
          <el-tag size="small" type="info">{{
            scope.row.batchCount || 0
          }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="model" label="型号" width="100" />
      <el-table-column
        prop="firmwareVersion"
        label="固件版本"
        min-width="130"
      />
      <el-table-column
        prop="mainboardFirmwareVersion"
        label="主板固件版本"
        min-width="150"
      />
      <el-table-column prop="pnCode" label="PN码" min-width="150" />
      <el-table-column label="业务类型" width="100">
        <template #default="scope">{{
          catLabel(scope.row.instrumentCategory)
        }}</template>
      </el-table-column>
      <el-table-column prop="submitterName" label="提交人" width="100" />
      <el-table-column prop="deviceCount" label="设备数" width="80" />
      <el-table-column label="合格数" width="90">
        <template #default="scope">
          <span class="count-pass">{{ scope.row.passCount || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="不合格数" width="100">
        <template #default="scope">
          <span class="count-fail">{{ scope.row.failCount || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="合格率" width="100">
        <template #default="scope">
          {{ passRateLabel(scope.row.passCount, scope.row.deviceCount) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="scope">
          <el-button
            size="small"
            type="primary"
            link
            @click="viewDetail(scope.row)"
            >详情</el-button
          >
          <el-button
            size="small"
            type="primary"
            link
            @click="editOrder(scope.row)"
            >编辑</el-button
          >
          <el-button
            v-if="scope.row.status === 0"
            size="small"
            type="danger"
            link
            @click="onDelete(scope.row.ID)"
            >删除</el-button
          >
          <el-button
            v-else
            size="small"
            type="danger"
            link
            @click="onForceDelete(scope.row)"
            >强制删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <div class="gva-pagination">
      <el-pagination
        v-model:current-page="searchInfo.page"
        v-model:page-size="searchInfo.pageSize"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        background
        small
        @size-change="getList"
        @current-change="getList"
      />
    </div>

    <el-drawer
      v-model="drawerVisible"
      title="编辑生产订单"
      size="520px"
      destroy-on-close
    >
      <el-form
        :model="formData"
        ref="formRef"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item label="MO号" prop="moNumber"
          ><el-input v-model="formData.moNumber"
        /></el-form-item>
        <el-form-item label="内部型号"
          ><el-input v-model="formData.model"
        /></el-form-item>
        <el-form-item label="固件版本"
          ><el-input v-model="formData.firmwareVersion"
        /></el-form-item>
        <el-form-item label="主板固件版本"
          ><el-input v-model="formData.mainboardFirmwareVersion"
        /></el-form-item>
        <el-form-item label="PN码"
          ><el-input v-model="formData.pnCode"
        /></el-form-item>
        <el-form-item label="业务类型">
          <el-select v-model="formData.instrumentCategory" style="width: 100%">
            <el-option label="线上" value="online" />
            <el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" />
            <el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"
          ><el-input v-model="formData.remark" type="textarea"
        /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit">确定</el-button>
      </template>
    </el-drawer>

    <el-dialog
      v-model="detailVisible"
      title="生产订单详情"
      width="1100px"
      destroy-on-close
    >
      <div v-if="detailOrder">
        <el-descriptions :column="2" border size="small" class="mb-4">
          <el-descriptions-item label="MO号">{{
            detailOrder.moNumber || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="内部型号">{{
            detailOrder.model || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="固件版本">{{
            detailOrder.firmwareVersion || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="主板固件版本">{{
            detailOrder.mainboardFirmwareVersion || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="PN码">{{
            detailOrder.pnCode || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="业务类型">{{
            catLabel(detailOrder.instrumentCategory) || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="批次数">{{
            detailOrder.batchCount ?? detailOrder.batches?.length ?? 0
          }}</el-descriptions-item>
          <el-descriptions-item label="设备数">{{
            detailOrder.deviceCount ?? detailOrder.devices?.length ?? 0
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="mb-3">
          <el-alert
            title="当前页面按生产订单号汇总显示；模板和检测工单生成发生在批次层。"
            type="info"
            :closable="false"
          />
        </div>

        <div v-if="unbatchedDevices.length" class="mb-4">
          <el-text type="info" size="small"
            >无批次 ({{ unbatchedDevices.length }}台)</el-text
          >
          <el-table :data="unbatchedDevices" border size="small" class="mt-1">
            <el-table-column prop="sn" label="SN" min-width="140" />
            <el-table-column prop="model" label="型号" width="90" />
            <el-table-column prop="pnCode" label="PN码" width="130" />
            <el-table-column prop="firmwareVersion" label="固件" width="110" />
            <el-table-column
              prop="mainboardFirmwareVersion"
              label="主板固件"
              width="130"
            />
          </el-table>
        </div>

        <div v-for="batch in detailOrder.batches" :key="batch.ID" class="mb-4">
          <div class="flex items-center justify-between mb-2">
            <el-text type="primary" size="small"
              >批次: {{ batch.batchNumber || '-' }}</el-text
            >
            <div class="flex items-center gap-2">
              <div class="text-sm text-gray-500">
                设备数: {{ batch.devices?.length || 0 }}
              </div>
              <el-select
                v-model="batchTemplateMap[batch.ID]"
                size="small"
                placeholder="选择检测模板"
                clearable
                style="width: 220px"
              >
                <el-option
                  v-for="template in templateList"
                  :key="template.ID"
                  :label="template.name"
                  :value="template.ID"
                />
              </el-select>
              <el-button
                v-if="batch.status === 0"
                type="primary"
                size="small"
                :disabled="!batchTemplateMap[batch.ID]"
                @click="assignTemplateToBatch(batch)"
              >
                生成检测工单
              </el-button>
              <el-tag
                v-else
                size="small"
                :type="
                  batch.status === 1
                    ? 'warning'
                    : batch.status === 2
                    ? 'primary'
                    : 'success'
                "
              >
                {{ batchStatusLabel(batch.status) }}
              </el-tag>
            </div>
          </div>
          <el-table :data="batch.devices" border size="small" class="mt-1">
            <el-table-column prop="sn" label="SN" min-width="140" />
            <el-table-column prop="model" label="型号" width="90" />
            <el-table-column prop="pnCode" label="PN码" width="130" />
            <el-table-column prop="firmwareVersion" label="固件" width="110" />
            <el-table-column
              prop="mainboardFirmwareVersion"
              label="主板固件"
              width="130"
            />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, reactive, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    getProductionOrderList,
    deleteProductionOrder,
    forceDeleteProductionOrder,
    updateProductionOrder,
    findProductionOrder
  } from '@/plugin/inspection/api/production_order'
  import { getTemplateList } from '@/plugin/inspection/api/template'
  import { assignBatchTemplate } from '@/plugin/inspection/api/work_order'

  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const drawerVisible = ref(false)
  const detailVisible = ref(false)
  const formRef = ref(null)
  const detailOrder = ref(null)
  const templateList = ref([])
  const batchTemplateMap = reactive({})

  const searchInfo = reactive({
    moNumber: '',
    model: '',
    instrumentCategory: '',
    page: 1,
    pageSize: 30
  })

  const formData = reactive({
    ID: 0,
    moNumber: '',
    model: '',
    firmwareVersion: '',
    mainboardFirmwareVersion: '',
    pnCode: '',
    instrumentCategory: '',
    remark: ''
  })

  const rules = {
    moNumber: [{ required: true, message: '请输入MO号', trigger: 'blur' }]
  }

  const catLabel = (value) =>
    ({
      online: '线上',
      offline: '线下',
      foreign_trade: '外贸',
      custom: '定制款'
    }[value] || value)
  const batchStatusLabel = (value) =>
    ({ 0: '未生成', 1: '待检测', 2: '检测中', 3: '已完成' }[value] || value)
  const passRateLabel = (passCount, deviceCount) => {
    const total = Number(deviceCount || 0)
    if (!total) return '-'
    return `${((Number(passCount || 0) / total) * 100).toFixed(1)}%`
  }

  const unbatchedDevices = computed(() => {
    if (!detailOrder.value) return []
    const allDevices = detailOrder.value.devices || []
    const batchedSet = new Set()
    detailOrder.value.batches?.forEach((batch) =>
      batch.devices?.forEach((device) => batchedSet.add(device.ID))
    )
    return allDevices.filter((device) => !batchedSet.has(device.ID))
  })

  const getList = async () => {
    loading.value = true
    try {
      const res = await getProductionOrderList(searchInfo)
      if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
      }
    } finally {
      loading.value = false
    }
  }

  const resetSearch = () => {
    searchInfo.moNumber = ''
    searchInfo.model = ''
    searchInfo.instrumentCategory = ''
    searchInfo.page = 1
    getList()
  }

  const editOrder = async (row) => {
    formData.ID = row.ID
    formData.moNumber = row.moNumber || ''
    formData.model = row.model || ''
    formData.firmwareVersion = row.firmwareVersion || ''
    formData.mainboardFirmwareVersion = row.mainboardFirmwareVersion || ''
    formData.pnCode = row.pnCode || ''
    formData.instrumentCategory = row.instrumentCategory || ''
    formData.remark = row.remark || ''
    drawerVisible.value = true
  }

  const submitEdit = async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return
    const res = await updateProductionOrder({ ...formData })
    if (res.code !== 0) return
    ElMessage.success('更新成功')
    drawerVisible.value = false
    getList()
  }

  const loadTemplates = async () => {
    if (templateList.value.length > 0) return
    const res = await getTemplateList({ page: 1, pageSize: 100 })
    if (res.code === 0) {
      templateList.value = res.data.list
    }
  }

  const viewDetail = async (row) => {
    await loadTemplates()
    const res = await findProductionOrder({ id: row.ID })
    if (res.code !== 0) return
    detailOrder.value = res.data
    Object.keys(batchTemplateMap).forEach((key) => delete batchTemplateMap[key])
    ;(res.data.batches || []).forEach((batch) => {
      batchTemplateMap[batch.ID] = batch.templateID || null
    })
    detailVisible.value = true
  }

  const assignTemplateToBatch = async (batch) => {
    const templateID = batchTemplateMap[batch.ID]
    if (!templateID) {
      ElMessage.warning('请先选择检测模板')
      return
    }
    const res = await assignBatchTemplate({
      ID: batch.ID,
      templateID: Number(templateID)
    })
    if (res.code !== 0) return
    ElMessage.success('已生成待检测工单')
    if (detailOrder.value) {
      await viewDetail({ ID: detailOrder.value.ID })
    }
    getList()
  }

  const onDelete = (id) => {
    ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
      .then(async () => {
        const res = await deleteProductionOrder({ id })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          getList()
        }
      })
      .catch(() => {})
  }

  const onForceDelete = (row) => {
    ElMessageBox.confirm(
      `确定强制删除生产订单 ${row.moNumber} 吗？这会一并删除设备、批次和检测结果。`,
      '强制删除',
      { type: 'warning', confirmButtonText: '确认删除' }
    )
      .then(async () => {
        const res = await forceDeleteProductionOrder({ id: row.ID })
        if (res.code === 0) {
          ElMessage.success('强制删除成功')
          getList()
        }
      })
      .catch(() => {})
  }

  getList()
</script>

<style scoped>
  .count-pass {
    color: #16a34a;
    font-weight: 600;
  }

  .count-fail {
    color: #dc2626;
    font-weight: 600;
  }
</style>
