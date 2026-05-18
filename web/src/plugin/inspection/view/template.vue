<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="模板名称">
          <el-input v-model="searchInfo.name" placeholder="请输入名称" clearable />
        </el-form-item>
        <el-form-item label="型号">
          <el-input v-model="searchInfo.model" placeholder="请输入型号" clearable />
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
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="name" label="模板名称" min-width="140" />
      <el-table-column prop="productName" label="产品名称" min-width="140" />
      <el-table-column prop="model" label="型号" width="100" />
      <el-table-column prop="firmwareVersion" label="固件版本" width="130" />
      <el-table-column label="状态" width="80">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
            {{ scope.row.status === 1 ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="itemCount" label="检测项数" width="100" />
      <el-table-column prop="CreatedAt" label="创建时间" width="170">
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
      :title="dialogType === 'create' ? '新增检测模板' : '编辑检测模板'"
      size="640px"
      destroy-on-close
    >
      <el-form :model="formData" ref="formRef" :rules="rules" label-width="100px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="formData.name" placeholder="如 G3X标准检测" />
        </el-form-item>
        <el-form-item label="产品名称" prop="productName">
          <el-input v-model="formData.productName" placeholder="如 GNSS接收机（RTK）" />
        </el-form-item>
        <el-form-item label="型号" prop="model">
          <el-input v-model="formData.model" placeholder="如 G3X" />
        </el-form-item>
        <el-form-item label="固件版本" prop="firmwareVersion">
          <el-input v-model="formData.firmwareVersion" placeholder="如 UM980-11833" />
        </el-form-item>
        <el-form-item v-if="dialogType === 'update'" label="状态">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="2" active-text="启用" inactive-text="停用" />
        </el-form-item>

        <el-divider>检测项分配 & 排序</el-divider>

        <div class="mb-2">
          <el-button size="small" @click="showItemPicker = true">从池中选择检测项</el-button>
        </div>

        <el-table :data="selectedItems" border size="small" max-height="400">
          <el-table-column label="排序" width="70">
            <template #default="scope">
              <el-input-number v-model="scope.row.sort" :min="1" size="small" controls-position="right" style="width:60px" />
            </template>
          </el-table-column>
          <el-table-column prop="name" label="检测项名称" min-width="140" />
          <el-table-column label="类型" width="100">
            <template #default="scope">
              <el-tag v-if="scope.row.resultType === 'pass_fail'" type="success" size="small">仅勾选</el-tag>
              <el-tag v-else-if="scope.row.resultType === 'number'" type="primary" size="small">仅数值</el-tag>
              <el-tag v-else type="warning" size="small">勾选+数值</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="unit" label="单位" width="60" />
          <el-table-column label="操作" width="70">
            <template #default="scope">
              <el-button size="small" type="danger" link @click="removeItem(scope.$index)">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form>

      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-drawer>

    <!-- Item picker dialog -->
    <el-dialog v-model="showItemPicker" title="选择检测项" width="500px" destroy-on-close>
      <el-table :data="allItems" border @selection-change="onPickerSelection" ref="pickerTable" max-height="400">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="类型" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.resultType === 'pass_fail'" type="success" size="small">仅勾选</el-tag>
            <el-tag v-else-if="scope.row.resultType === 'number'" type="primary" size="small">仅数值</el-tag>
            <el-tag v-else type="warning" size="small">勾选+数值</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="showItemPicker = false">取消</el-button>
        <el-button type="primary" @click="confirmPicker">确认选择</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import { getTemplateList, createTemplate, deleteTemplate, updateTemplate, findTemplate } from '@/plugin/inspection/api/template'
import { getItemList } from '@/plugin/inspection/api/inspection_item'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const drawerVisible = ref(false)
const dialogType = ref('create')
const formRef = ref(null)
const allItems = ref([])
const selectedItems = ref([])
const showItemPicker = ref(false)
const pickerSelected = ref([])

const searchInfo = reactive({ name: '', model: '', page: 1, pageSize: 30 })
const formData = reactive({ ID: 0, name: '', productName: '', model: '', firmwareVersion: '', status: 1 })

const rules = { name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }] }

const getList = async () => {
  loading.value = true
  try {
    const res = await getTemplateList(searchInfo)
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  } finally { loading.value = false }
}

const resetSearch = () => { searchInfo.name = ''; searchInfo.model = ''; searchInfo.page = 1; getList() }

const fetchAllItems = async () => {
  const res = await getItemList({ page: 1, pageSize: 1000 })
  if (res.code === 0) allItems.value = res.data.list
}

const openDialog = async (type, row) => {
  dialogType.value = type
  await fetchAllItems()
  if (type === 'update' && row) {
    formData.ID = row.ID
    formData.name = row.name
    formData.productName = row.productName || ''
    formData.model = row.model || ''
    formData.firmwareVersion = row.firmwareVersion || ''
    formData.status = row.status
    const res = await findTemplate({ id: row.ID })
    if (res.code === 0 && res.data.templateItems) {
      selectedItems.value = res.data.templateItems.map((ti) => ({
        itemID: ti.itemID,
        name: ti.item?.name || '',
        resultType: ti.item?.resultType || 'pass_fail',
        unit: ti.item?.unit || '',
        sort: ti.sort
      }))
    }
  } else {
    formData.ID = 0; formData.name = ''; formData.productName = ''; formData.model = ''
    formData.firmwareVersion = ''; formData.status = 1
    selectedItems.value = []
  }
  drawerVisible.value = true
  await nextTick(); formRef.value?.clearValidate()
}

const onPickerSelection = (val) => { pickerSelected.value = val }

const confirmPicker = () => {
  for (const row of pickerSelected.value) {
    if (!selectedItems.value.find((s) => s.itemID === row.ID)) {
      const maxSort = selectedItems.value.length > 0 ? Math.max(...selectedItems.value.map((s) => s.sort)) : 0
      selectedItems.value.push({
        itemID: row.ID,
        name: row.name,
        resultType: row.resultType,
        unit: row.unit || '',
        sort: maxSort + 1
      })
    }
  }
  showItemPicker.value = false
}

const removeItem = (index) => { selectedItems.value.splice(index, 1) }

const submitForm = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  if (selectedItems.value.length === 0) { ElMessage.warning('请至少选择一个检测项'); return }

  const payload = {
    name: formData.name,
    productName: formData.productName,
    model: formData.model,
    firmwareVersion: formData.firmwareVersion,
    items: selectedItems.value.map((s) => ({ itemID: s.itemID, sort: s.sort }))
  }

  try {
    if (dialogType.value === 'create') {
      const res = await createTemplate(payload)
      if (res.code === 0) ElMessage.success('创建成功')
    } else {
      payload.ID = formData.ID
      payload.status = formData.status
      const res = await updateTemplate(payload)
      if (res.code === 0) ElMessage.success('编辑成功')
    }
    drawerVisible.value = false
    getList()
  } catch (e) { /* handled by interceptor */ }
}

const onDelete = (id) => {
  ElMessageBox.confirm('删除模板会同时删除其关联的检测项配置，确定？', '提示', { type: 'warning' })
    .then(async () => {
      const res = await deleteTemplate({ id })
      if (res.code === 0) { ElMessage.success('删除成功'); getList() }
    })
    .catch(() => {})
}

getList()
</script>
