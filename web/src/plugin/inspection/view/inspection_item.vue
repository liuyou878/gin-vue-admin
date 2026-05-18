<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" ref="searchFormRef" inline>
        <el-form-item label="检测项名称">
          <el-input v-model="searchInfo.name" placeholder="请输入名称" clearable />
        </el-form-item>
        <el-form-item label="结果类型">
          <el-select v-model="searchInfo.resultType" placeholder="请选择类型" clearable style="width:160px">
            <el-option label="仅勾选" value="pass_fail" />
            <el-option label="仅数值" value="number" />
            <el-option label="勾选+数值" value="both" />
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
      <el-button :disabled="!multipleSelection.length" @click="onDeleteByIds">批量删除</el-button>
    </div>

    <el-table
      :data="tableData"
      @selection-change="handleSelectionChange"
      border
      v-loading="loading"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="name" label="检测项名称" min-width="160" />
      <el-table-column label="结果类型" width="120">
        <template #default="scope">
          <el-tag v-if="scope.row.resultType === 'pass_fail'" type="success">仅勾选</el-tag>
          <el-tag v-else-if="scope.row.resultType === 'number'" type="primary">仅数值</el-tag>
          <el-tag v-else type="warning">勾选+数值</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="unit" label="单位" width="80" />
      <el-table-column label="合格范围" width="180">
        <template #default="scope">
          <span v-if="scope.row.minValue != null || scope.row.maxValue != null">
            {{ scope.row.minValue ?? '-' }} ~ {{ scope.row.maxValue ?? '-' }}
          </span>
          <span v-else class="text-gray">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      <el-table-column prop="CreatedAt" label="创建时间" width="170">
        <template #default="scope">
          {{ formatDate(scope.row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
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
      :title="dialogType === 'create' ? '新增检测项' : '编辑检测项'"
      size="480px"
      destroy-on-close
    >
      <el-form :model="formData" ref="formRef" :rules="rules" label-width="100px">
        <el-form-item label="检测项名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入检测项名称" />
        </el-form-item>
        <el-form-item label="结果类型" prop="resultType">
          <el-select v-model="formData.resultType" placeholder="请选择结果类型" style="width:100%">
            <el-option label="仅勾选 (通过/未通过)" value="pass_fail" />
            <el-option label="仅数值" value="number" />
            <el-option label="勾选+数值 (兼有)" value="both" />
          </el-select>
        </el-form-item>
        <el-form-item label="单位" prop="unit" v-if="formData.resultType !== 'pass_fail'">
          <el-input v-model="formData.unit" placeholder="如 V、%'" />
        </el-form-item>
        <el-form-item label="合格范围" v-if="formData.resultType !== 'pass_fail'">
          <el-input-number v-model="formData.minValue" :precision="2" placeholder="下限" style="width:45%" controls-position="right" />
          <span class="mx-2">~</span>
          <el-input-number v-model="formData.maxValue" :precision="2" placeholder="上限" style="width:45%" controls-position="right" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" type="textarea" placeholder="备注说明" />
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
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import {
  getItemList,
  createItem,
  deleteItem,
  deleteItemByIds,
  updateItem,
  findItem
} from '@/plugin/inspection/api/inspection_item'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const multipleSelection = ref([])
const drawerVisible = ref(false)
const dialogType = ref('create')
const formRef = ref(null)
const searchFormRef = ref(null)

const searchInfo = reactive({
  name: '',
  resultType: '',
  page: 1,
  pageSize: 30
})

const formData = reactive({
  ID: 0,
  name: '',
  resultType: 'pass_fail',
  unit: '',
  minValue: null,
  maxValue: null,
  remark: ''
})

const rules = {
  name: [{ required: true, message: '请输入检测项名称', trigger: 'blur' }],
  resultType: [{ required: true, message: '请选择结果类型', trigger: 'change' }]
}

const getList = async () => {
  loading.value = true
  try {
    const res = await getItemList(searchInfo)
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchInfo.name = ''
  searchInfo.resultType = ''
  searchInfo.page = 1
  getList()
}

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const openDialog = (type, row) => {
  dialogType.value = type
  if (type === 'update' && row) {
    formData.ID = row.ID
    formData.name = row.name
    formData.resultType = row.resultType
    formData.unit = row.unit || ''
    formData.minValue = row.minValue
    formData.maxValue = row.maxValue
    formData.remark = row.remark || ''
  } else {
    resetForm()
  }
  drawerVisible.value = true
}

const resetForm = () => {
  formData.ID = 0
  formData.name = ''
  formData.resultType = 'pass_fail'
  formData.unit = ''
  formData.minValue = null
  formData.maxValue = null
  formData.remark = ''
  formRef.value?.resetFields()
}

const submitForm = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  const payload = {
    name: formData.name,
    resultType: formData.resultType,
    unit: formData.unit || '',
    minValue: formData.minValue,
    maxValue: formData.maxValue,
    remark: formData.remark || ''
  }

  try {
    if (dialogType.value === 'create') {
      const res = await createItem(payload)
      if (res.code === 0) {
        ElMessage.success('创建成功')
      }
    } else {
      payload.ID = formData.ID
      const res = await updateItem(payload)
      if (res.code === 0) {
        ElMessage.success('编辑成功')
      }
    }
    drawerVisible.value = false
    getList()
  } catch (e) {
    // error handled by request interceptor
  }
}

const onDelete = (id) => {
  ElMessageBox.confirm('确定删除该检测项吗？', '提示', { type: 'warning' })
    .then(async () => {
      const res = await deleteItem({ id })
      if (res.code === 0) {
        ElMessage.success('删除成功')
        getList()
      }
    })
    .catch(() => {})
}

const onDeleteByIds = () => {
  const ids = multipleSelection.value.map((row) => row.ID).join(',')
  ElMessageBox.confirm(`确定删除选中的 ${multipleSelection.value.length} 个检测项吗？`, '提示', { type: 'warning' })
    .then(async () => {
      const res = await deleteItemByIds({ ids })
      if (res.code === 0) {
        ElMessage.success('批量删除成功')
        getList()
      }
    })
    .catch(() => {})
}

getList()
</script>
