<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="MO号">
          <el-input v-model="searchInfo.moNumber" placeholder="请输入" clearable size="small" />
        </el-form-item>
        <el-form-item label="型号">
          <el-input v-model="searchInfo.model" placeholder="请输入" clearable size="small" />
        </el-form-item>
        <el-form-item label="仪器类别">
          <el-select v-model="searchInfo.instrumentCategory" placeholder="请选择" clearable size="small" style="width:120px">
            <el-option label="线上" value="online" /><el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" /><el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择" clearable size="small" style="width:110px">
            <el-option label="待确认" :value="0" /><el-option label="待检测" :value="1" />
            <el-option label="检测中" :value="2" /><el-option label="已完成" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="getList">查询</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData" border v-loading="loading" size="small">
      <el-table-column prop="moNumber" label="MO号" min-width="140" />
      <el-table-column label="模板" min-width="120">
        <template #default="s1">{{ s1.row.template?.name || '-' }}</template>
      </el-table-column>
      <el-table-column prop="productName" label="产品名称" min-width="120" />
      <el-table-column prop="model" label="型号" width="80" />
      <el-table-column label="类别" width="70">
        <template #default="s2">{{ catLabel(s2.row.instrumentCategory) }}</template>
      </el-table-column>
      <el-table-column prop="submitterName" label="提交人" width="100" />
      <el-table-column label="状态" width="80">
        <template #default="s3">
          <el-tag :type="statusTag(s3.row.status)" size="small">{{ statusLabel(s3.row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="deviceCount" label="设备数" width="70" />
      <el-table-column label="合格/不合格" width="100">
        <template #default="s4">{{ s4.row.passCount||0 }}/{{ s4.row.failCount||0 }}</template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="s5">
          <el-button v-if="s5.row.status===0" size="small" type="primary" @click="confirmOrder(s5.row)">确认</el-button>
          <el-button size="small" type="primary" link @click="viewDetail(s5.row)">设备</el-button>
          <el-button v-if="s5.row.status===0" size="small" type="primary" link @click="editOrder(s5.row)">编辑</el-button>
          <el-button v-if="s5.row.status===0" size="small" type="danger" link @click="onDelete(s5.row.ID)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="gva-pagination">
      <el-pagination
        v-model:current-page="searchInfo.page" v-model:page-size="searchInfo.pageSize"
        :page-sizes="[10,30,50,100]" :total="total" layout="total, sizes, prev, pager, next, jumper" background
        @size-change="getList" @current-change="getList" small
      />
    </div>

    <!-- Edit drawer -->
    <el-drawer v-model="drawerVisible" title="编辑生产订单" size="480px" destroy-on-close>
      <el-form :model="formData" ref="formRef" :rules="rules" label-width="100px">
        <el-form-item label="MO号" prop="moNumber"><el-input v-model="formData.moNumber" /></el-form-item>
        <el-form-item label="模板">
          <el-select v-model="formData.templateID" placeholder="选择检测模板" style="width:100%" clearable>
            <el-option v-for="t in templateList" :key="t.ID" :label="t.name" :value="t.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="产品名称"><el-input v-model="formData.productName" /></el-form-item>
        <el-form-item label="型号"><el-input v-model="formData.model" /></el-form-item>
        <el-form-item label="固件版本"><el-input v-model="formData.firmwareVersion" /></el-form-item>
        <el-form-item label="仪器类别">
          <el-select v-model="formData.instrumentCategory" style="width:100%">
            <el-option label="线上" value="online" /><el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" /><el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible=false">取消</el-button>
        <el-button type="primary" @click="submitEdit">确定</el-button>
      </template>
    </el-drawer>

    <!-- Device detail dialog -->
    <el-dialog v-model="detailVisible" title="设备列表" width="800px" destroy-on-close>
      <div v-if="detailOrder">
        <div class="mb-2 text-sm">
          <b>MO号:</b> {{ detailOrder.moNumber }} | <b>型号:</b> {{ detailOrder.model }}
        </div>
        <!-- Unbatched devices -->
        <div v-if="unbatchedDevices.length" class="mb-4">
          <el-text type="info" size="small">无批次 ({{ unbatchedDevices.length }}台)</el-text>
          <el-table :data="unbatchedDevices" border size="small" class="mt-1">
            <el-table-column prop="sn" label="SN" min-width="140" />
            <el-table-column prop="model" label="型号" width="90" />
            <el-table-column prop="pnCode" label="PN码" width="130" />
            <el-table-column prop="firmwareVersion" label="固件" width="100" />
          </el-table>
        </div>
        <!-- Batches -->
        <div v-for="b in detailOrder.batches" :key="b.ID" class="mb-4">
          <el-text type="primary" size="small">批次: {{ b.batchNumber }} ({{ b.devices?.length||0 }}台)</el-text>
          <el-table :data="b.devices" border size="small" class="mt-1">
            <el-table-column prop="sn" label="SN" min-width="140" />
            <el-table-column prop="model" label="型号" width="90" />
            <el-table-column prop="pnCode" label="PN码" width="130" />
            <el-table-column prop="firmwareVersion" label="固件" width="100" />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { getProductionOrderList, deleteProductionOrder, updateProductionOrder, findProductionOrder } from '@/plugin/inspection/api/production_order'
import { getTemplateList } from '@/plugin/inspection/api/template'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const drawerVisible = ref(false)
const detailVisible = ref(false)
const formRef = ref(null)
const templateList = ref([])
const detailOrder = ref(null)

const searchInfo = reactive({ moNumber: '', model: '', instrumentCategory: '', status: null, page: 1, pageSize: 30 })
const formData = reactive({ ID: 0, moNumber: '', templateID: null, productName: '', model: '', firmwareVersion: '', instrumentCategory: '', remark: '' })
const rules = { moNumber: [{ required: true, message: '请输入MO号', trigger: 'blur' }] }

const catLabel = (v) => ({ online: '线上', offline: '线下', foreign_trade: '外贸', custom: '定制款' }[v] || v)
const statusLabel = (v) => ({ 0: '待确认', 1: '待检测', 2: '检测中', 3: '已完成' }[v] || v)
const statusTag = (v) => ({ 0: 'info', 1: 'warning', 2: 'primary', 3: 'success' }[v] || '')

const unbatchedDevices = computed(() => {
  if (!detailOrder.value) return []
  const allDevices = detailOrder.value.devices || []
  const batchedSet = new Set()
  detailOrder.value.batches?.forEach(b => b.devices?.forEach(d => batchedSet.add(d.ID)))
  return allDevices.filter(d => !batchedSet.has(d.ID))
})

const getList = async () => {
  loading.value = true
  try {
    const res = await getProductionOrderList(searchInfo)
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  } finally { loading.value = false }
}
const resetSearch = () => { searchInfo.moNumber = ''; searchInfo.model = ''; searchInfo.instrumentCategory = ''; searchInfo.status = null; searchInfo.page = 1; getList() }

const confirmOrder = async (row) => {
  await ElMessageBox.confirm('确认该生产订单，设为待检测？', '提示', { type: 'info' })
  await updateProductionOrder({ ID: row.ID, moNumber: row.moNumber, status: 1 })
  ElMessage.success('已确认，进入待检测'); getList()
}

const editOrder = async (row) => {
  const tRes = await getTemplateList({ page: 1, pageSize: 100 })
  if (tRes.code === 0) templateList.value = tRes.data.list
  formData.ID = row.ID; formData.moNumber = row.moNumber; formData.templateID = row.templateID || null
  formData.productName = row.productName||''; formData.model = row.model||''
  formData.firmwareVersion = row.firmwareVersion||''; formData.instrumentCategory = row.instrumentCategory||''
  formData.remark = row.remark||''
  drawerVisible.value = true
}

const submitEdit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  await updateProductionOrder({ ...formData })
  ElMessage.success('更新成功'); drawerVisible.value = false; getList()
}

const viewDetail = async (row) => {
  const res = await findProductionOrder({ id: row.ID })
  if (res.code === 0) { detailOrder.value = res.data; detailVisible.value = true }
}

const onDelete = (id) => {
  ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteProductionOrder({ id })
    if (res.code === 0) { ElMessage.success('删除成功'); getList() }
  }).catch(() => {})
}

getList()
</script>
