<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="MO号">
          <el-input v-model="searchInfo.moNumber" placeholder="请输入MO号" clearable />
        </el-form-item>
        <el-form-item label="型号">
          <el-input v-model="searchInfo.model" placeholder="请输入型号" clearable />
        </el-form-item>
        <el-form-item label="仪器类别">
          <el-select v-model="searchInfo.instrumentCategory" placeholder="请选择" clearable style="width:130px">
            <el-option label="线上" value="online" />
            <el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" />
            <el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择" clearable style="width:120px">
            <el-option label="待检测" :value="1" />
            <el-option label="检测中" :value="2" />
            <el-option label="已完成" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-btn-list">
      <el-button type="primary" @click="openDialog('create')">新增</el-button>
    </div>

    <el-table :data="tableData" border v-loading="loading">
      <el-table-column prop="ID" label="ID" width="70" />
      <el-table-column prop="moNumber" label="MO号" min-width="140" />
      <el-table-column label="模板" min-width="120">
        <template #default="scope">{{ scope.row.template?.name || '-' }}</template>
      </el-table-column>
      <el-table-column prop="productName" label="产品名称" min-width="130" />
      <el-table-column prop="model" label="型号" width="90" />
      <el-table-column prop="firmwareVersion" label="固件版本" width="120" />
      <el-table-column label="仪器类别" width="100">
        <template #default="scope">
          <el-tag :type="categoryTagType(scope.row.instrumentCategory)">
            {{ categoryLabel(scope.row.instrumentCategory) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="batchNumber" label="批次号" width="180" />
      <el-table-column label="状态" width="90">
        <template #default="scope">
          <el-tag :type="statusTagType(scope.row.status)">{{ statusLabel(scope.row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="deviceCount" label="SN数量" width="80" />
      <el-table-column prop="CreatedAt" label="提交时间" width="170">
        <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="scope">
          <el-button size="small" type="primary" link @click="openDialog('update', scope.row)">编辑</el-button>
          <el-button size="small" type="danger" link @click="onDelete(scope.row.ID)">删除</el-button>
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
        @size-change="getList"
        @current-change="getList"
      />
    </div>

    <el-drawer
      v-model="drawerVisible"
      :title="dialogType === 'create' ? '新增生产订单' : '编辑生产订单'"
      size="640px"
      destroy-on-close
    >
      <el-form :model="formData" ref="formRef" :rules="rules" label-width="100px">
        <el-form-item label="MO号" prop="moNumber">
          <el-input v-model="formData.moNumber" placeholder="如 MO2026030023" />
        </el-form-item>
        <el-form-item label="检测模板">
          <el-select v-model="formData.templateID" placeholder="选择检测模板（定检测标准）" style="width:100%" clearable>
            <el-option v-for="t in templateList" :key="t.ID" :label="t.name" :value="t.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="产品名称">
          <el-input v-model="formData.productName" placeholder="如 GNSS接收机（RTK）" />
        </el-form-item>
        <el-form-item label="型号">
          <el-input v-model="formData.model" placeholder="如 G3X" />
        </el-form-item>
        <el-form-item label="固件版本">
          <el-input v-model="formData.firmwareVersion" placeholder="如 UM980-11833" />
        </el-form-item>
        <el-form-item label="仪器类别">
          <el-select v-model="formData.instrumentCategory" placeholder="请选择仪器类别" style="width:100%">
            <el-option label="线上" value="online" />
            <el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" />
            <el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="批次号">
          <el-input v-model="formData.batchNumber" placeholder="如 Alpha G3X 260324-06" />
        </el-form-item>
        <el-form-item v-if="dialogType === 'update'" label="状态">
          <el-select v-model="formData.status" style="width:100%">
            <el-option label="待检测" :value="1" />
            <el-option label="检测中" :value="2" />
            <el-option label="已完成" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" placeholder="备注" />
        </el-form-item>

        <el-divider>SN 列表（一个 SN 一行）</el-divider>
        <el-form-item label="SN列表" prop="snText">
          <el-input
            v-model="snText"
            type="textarea"
            :rows="10"
            placeholder="逐行输入或扫描枪录入 SN，一行一个"
          />
        </el-form-item>
        <el-form-item>
          <span class="text-gray">当前 {{ snCount }} 个 SN</span>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import {
  getProductionOrderList, createProductionOrder, deleteProductionOrder,
  updateProductionOrder, findProductionOrder
} from '@/plugin/inspection/api/production_order'
import { getTemplateList } from '@/plugin/inspection/api/template'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const drawerVisible = ref(false)
const dialogType = ref('create')
const formRef = ref(null)
const snText = ref('')

const searchInfo = reactive({ moNumber: '', model: '', instrumentCategory: '', status: null, page: 1, pageSize: 30 })
const templateList = ref([])
const formData = reactive({ ID: 0, moNumber: '', templateID: null, productName: '', model: '', firmwareVersion: '', instrumentCategory: '', batchNumber: '', status: 1, remark: '' })
const rules = { moNumber: [{ required: true, message: '请输入MO号', trigger: 'blur' }] }

const snCount = computed(() => snText.value.split('\n').filter(s => s.trim()).length)

const categoryLabel = (v) => ({ online: '线上', offline: '线下', foreign_trade: '外贸', custom: '定制款' }[v] || v)
const categoryTagType = (v) => ({ online: 'success', offline: '', foreign_trade: 'warning', custom: 'danger' }[v] || '')
const statusLabel = (v) => ({ 1: '待检测', 2: '检测中', 3: '已完成' }[v] || v)
const statusTagType = (v) => ({ 1: 'warning', 2: 'primary', 3: 'success' }[v] || '')

const getList = async () => {
  loading.value = true
  try {
    const res = await getProductionOrderList(searchInfo)
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  } finally { loading.value = false }
}

const resetSearch = () => {
  searchInfo.moNumber = ''; searchInfo.model = ''; searchInfo.instrumentCategory = ''
  searchInfo.status = null; searchInfo.page = 1
  getList()
}

const openDialog = async (type, row) => {
  dialogType.value = type
  const tRes = await getTemplateList({ page: 1, pageSize: 100 })
  if (tRes.code === 0) templateList.value = tRes.data.list
  if (type === 'update' && row) {
    formData.ID = row.ID
    formData.moNumber = row.moNumber
    formData.templateID = row.templateID || null
    formData.productName = row.productName || ''
    formData.model = row.model || ''
    formData.firmwareVersion = row.firmwareVersion || ''
    formData.instrumentCategory = row.instrumentCategory || ''
    formData.batchNumber = row.batchNumber || ''
    formData.status = row.status
    formData.remark = row.remark || ''
    const res = await findProductionOrder({ id: row.ID })
    if (res.code === 0 && res.data.devices) {
      snText.value = res.data.devices.map(d => d.sn).join('\n')
    }
  } else {
    formData.ID = 0; formData.moNumber = ''; formData.templateID = null
    formData.productName = ''; formData.model = ''; formData.firmwareVersion = ''
    formData.instrumentCategory = ''; formData.batchNumber = ''
    formData.status = 1; formData.remark = ''
    snText.value = ''
  }
  drawerVisible.value = true
}

const submitForm = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  const sns = snText.value.split('\n').map(s => s.trim()).filter(s => s)
  if (sns.length === 0) { ElMessage.warning('请至少输入一个 SN'); return }

  const payload = {
    moNumber: formData.moNumber,
    templateID: formData.templateID || null,
    productName: formData.productName,
    model: formData.model,
    firmwareVersion: formData.firmwareVersion,
    instrumentCategory: formData.instrumentCategory,
    batchNumber: formData.batchNumber,
    remark: formData.remark,
    sns
  }

  try {
    if (dialogType.value === 'create') {
      const res = await createProductionOrder(payload)
      if (res.code === 0) ElMessage.success('创建成功')
    } else {
      payload.ID = formData.ID; payload.status = formData.status
      const res = await updateProductionOrder(payload)
      if (res.code === 0) ElMessage.success('编辑成功')
    }
    drawerVisible.value = false
    getList()
  } catch (e) { /* handled by interceptor */ }
}

const onDelete = (id) => {
  ElMessageBox.confirm('删除订单会同时删除其 SN 列表，确定？', '提示', { type: 'warning' })
    .then(async () => {
      const res = await deleteProductionOrder({ id })
      if (res.code === 0) { ElMessage.success('删除成功'); getList() }
    })
    .catch(() => {})
}

getList()
</script>
