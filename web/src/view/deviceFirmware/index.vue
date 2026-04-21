<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item label="类别名称">
          <el-input
            v-model="categorySearch.name"
            clearable
            placeholder="类别名称"
          />
        </el-form-item>
        <el-form-item label="型号名称">
          <el-input
            v-model="modelSearch.modelName"
            clearable
            placeholder="型号名称"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="loadDeviceTree"
            >查询</el-button
          >
          <el-button icon="refresh" @click="resetDeviceTreeSearch"
            >重置</el-button
          >
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="section-head">
        <div>
          <div class="section-title">设备管理</div>
          <div class="section-subtitle">
            先维护设备类别，再在类别下面添加型号，形成树形结构。
          </div>
        </div>
        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openCategoryDialog()"
            >新增类别</el-button
          >
        </div>
      </div>

      <el-table
        :data="deviceTreeData"
        row-key="rowKey"
        default-expand-all
        :tree-props="{ children: 'children' }"
        :row-class-name="getDeviceTreeRowClassName"
      >
        <el-table-column label="序号" width="90">
          <template #default="scope">
            <span class="tree-seq-col">{{ scope.row.displayOrder }}</span>
          </template>
        </el-table-column>
        <el-table-column label="名称" min-width="220">
          <template #default="scope">
            <div
              class="tree-name-cell"
              :class="{
                'tree-name-cell--child': scope.row.nodeType === 'model'
              }"
            >
              <span class="tree-name">
                {{
                  scope.row.nodeType === 'category'
                    ? scope.row.name
                    : scope.row.modelName || '-'
                }}
              </span>
              <el-tag
                v-if="scope.row.nodeType === 'category'"
                type="info"
                size="small"
                >类别</el-tag
              >
              <el-tag v-else type="success" size="small">型号</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="260">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="备注" min-width="220" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.remark || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="420" fixed="right">
          <template #default="scope">
            <template v-if="scope.row.nodeType === 'category'">
              <el-button
                type="primary"
                link
                icon="plus"
                @click="openModelDialog(null, scope.row)"
                >添加型号</el-button
              >
              <el-button
                type="primary"
                link
                icon="edit"
                @click="openCategoryDialog(scope.row)"
                >编辑类别</el-button
              >
              <el-button
                type="primary"
                link
                icon="delete"
                @click="deleteCategoryRow(scope.row)"
                >删除类别</el-button
              >
            </template>
            <template v-else-if="scope.row.nodeType === 'model'">
              <el-button
                type="primary"
                link
                icon="edit"
                @click="openModelDialog(scope.row)"
                >编辑型号</el-button
              >
              <el-button
                type="primary"
                link
                icon="delete"
                @click="deleteModelRow(scope.row)"
                >删除型号</el-button
              >
            </template>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="categoryDialogVisible"
      :title="categoryDialogType === 'create' ? '新增设备类别' : '编辑设备类别'"
      width="520px"
    >
      <el-form
        ref="categoryFormRef"
        :model="categoryForm"
        :rules="categoryRules"
        label-width="90px"
      >
        <el-form-item label="类别名称" prop="name">
          <el-input v-model="categoryForm.name" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="categoryForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="categoryForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="categoryForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="modelDialogVisible"
      :title="modelDialogType === 'create' ? '新增设备型号' : '编辑设备型号'"
      width="620px"
    >
      <el-form
        ref="modelFormRef"
        :model="modelForm"
        :rules="modelRules"
        label-width="90px"
      >
        <el-form-item label="设备类别" prop="categoryId">
          <el-input :model-value="selectedModelCategoryName" readonly />
        </el-form-item>
        <el-form-item label="型号名称" prop="modelName">
          <el-input v-model="modelForm.modelName" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="modelForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="modelForm.remark"
            type="textarea"
            :rows="3"
            placeholder="可选"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modelDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitModel">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    createDeviceCategory,
    deleteDeviceCategory,
    updateDeviceCategory,
    findDeviceCategory,
    getDeviceCategoryList,
    createDeviceModel,
    deleteDeviceModel,
    updateDeviceModel,
    findDeviceModel,
    getDeviceModelList
  } from '@/api/deviceFirmware'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { computed, onMounted, ref } from 'vue'

  defineOptions({ name: 'DeviceFirmwareCenter' })

  const categoryOptions = ref([])
  const modelOptions = ref([])
  const categorySearch = ref({})
  const modelSearch = ref({})

  const categoryTableData = ref([])
  const modelTableData = ref([])

  const categoryDialogVisible = ref(false)
  const modelDialogVisible = ref(false)
  const categoryDialogType = ref('create')
  const modelDialogType = ref('create')
  const categoryFormRef = ref()
  const modelFormRef = ref()

  const categoryForm = ref({
    name: '',
    code: '',
    sort: 0,
    status: 1,
    remark: ''
  })
  const modelForm = ref({
    categoryId: '',
    modelCode: '',
    modelName: '',
    seriesName: '',
    status: 1,
    remark: ''
  })

  const buildCategoryCode = () => `device-category-${Date.now()}`
  const buildModelCode = () => `device-model-${Date.now()}`
  const requiredRule = (message) => ({
    required: true,
    message,
    trigger: ['blur', 'change']
  })

  const categoryRules = {
    name: [requiredRule('请输入类别名称')],
    sort: [requiredRule('请填写排序')],
    status: [requiredRule('请选择状态')],
    remark: [requiredRule('请填写备注')]
  }
  const modelRules = {
    categoryId: [requiredRule('请选择设备类别')],
    modelName: [requiredRule('请输入型号名称')],
    status: [requiredRule('请选择状态')]
  }

  const validateForm = async (formRef) => {
    if (!formRef.value) {
      return false
    }
    try {
      await formRef.value.validate()
      return true
    } catch (error) {
      return false
    }
  }

  const selectedModelCategoryName = computed(() => {
    if (!modelForm.value.categoryId) {
      return ''
    }
    return (
      categoryOptions.value.find(
        (item) => item.ID === modelForm.value.categoryId
      )?.name || ''
    )
  })

  const deviceTreeData = computed(() =>
    categoryTableData.value
      .map((category, categoryIndex) => ({
        ...category,
        nodeType: 'category',
        rowKey: `category-${category.ID}`,
        displayOrder: `${categoryIndex + 1}`,
        children: modelTableData.value
          .filter((model) => model.categoryId === category.ID)
          .map((model, modelIndex) => ({
            ...model,
            nodeType: 'model',
            rowKey: `model-${model.ID}`,
            displayOrder: `${modelIndex + 1}`
          }))
      }))
      .filter(
        (category) =>
          !modelSearch.value.modelName || category.children.length > 0
      )
  )

  const loadDeviceTree = async () => {
    await Promise.all([loadCategories(), loadModels()])
  }

  const resetDeviceTreeSearch = async () => {
    categorySearch.value = {}
    modelSearch.value = {}
    await loadDeviceTree()
  }

  const getDeviceTreeRowClassName = ({ row }) => {
    if (row.nodeType === 'category') {
      return 'device-tree-category-row'
    }
    return 'device-tree-model-row'
  }

  const loadCategories = async () => {
    const res = await getDeviceCategoryList({
      page: 1,
      pageSize: 999,
      ...categorySearch.value
    })
    if (res.code === 0) {
      categoryTableData.value = res.data.list || []
      categoryOptions.value = res.data.list || []
    }
  }

  const loadModels = async () => {
    const res = await getDeviceModelList({
      page: 1,
      pageSize: 999,
      ...modelSearch.value
    })
    if (res.code === 0) {
      modelTableData.value = res.data.list || []
      modelOptions.value = res.data.list || []
    }
  }

  const openCategoryDialog = async (row) => {
    categoryDialogType.value = row?.ID ? 'update' : 'create'
    if (row?.ID) {
      const res = await findDeviceCategory({ ID: row.ID })
      if (res.code === 0) {
        categoryForm.value = { ...res.data }
      }
    } else {
      categoryForm.value = {
        name: '',
        code: '',
        sort: 0,
        status: 1,
        remark: ''
      }
    }
    categoryDialogVisible.value = true
  }

  const submitCategory = async () => {
    if (!(await validateForm(categoryFormRef))) {
      return
    }
    const payload = {
      ...categoryForm.value,
      code: categoryForm.value.code || buildCategoryCode()
    }
    const res =
      categoryDialogType.value === 'create'
        ? await createDeviceCategory(payload)
        : await updateDeviceCategory(payload)
    if (res.code === 0) {
      ElMessage.success(
        categoryDialogType.value === 'create' ? '创建成功' : '更新成功'
      )
      categoryDialogVisible.value = false
      await loadDeviceTree()
    }
  }

  const deleteCategoryRow = (row) => {
    ElMessageBox.confirm(`确定删除类别 ${row.name} 吗？`, '提示', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteDeviceCategory({ ID: row.ID })
      if (res.code === 0) {
        ElMessage.success('删除成功')
        await loadDeviceTree()
      }
    })
  }

  const openModelDialog = async (row, parentCategory) => {
    modelDialogType.value = row?.ID ? 'update' : 'create'
    if (row?.ID) {
      const res = await findDeviceModel({ ID: row.ID })
      if (res.code === 0) {
        modelForm.value = { ...res.data }
      }
    } else {
      modelForm.value = {
        categoryId: parentCategory?.ID || '',
        modelCode: '',
        modelName: '',
        seriesName: '',
        status: 1,
        remark: ''
      }
    }
    modelDialogVisible.value = true
  }

  const submitModel = async () => {
    if (!(await validateForm(modelFormRef))) {
      return
    }
    const payload = {
      ...modelForm.value,
      modelCode: modelForm.value.modelCode || buildModelCode()
    }
    const res =
      modelDialogType.value === 'create'
        ? await createDeviceModel(payload)
        : await updateDeviceModel(payload)
    if (res.code === 0) {
      ElMessage.success(
        modelDialogType.value === 'create' ? '创建成功' : '更新成功'
      )
      modelDialogVisible.value = false
      await loadDeviceTree()
    }
  }

  const deleteModelRow = (row) => {
    ElMessageBox.confirm(`确定删除型号 ${row.modelName} 吗？`, '提示', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteDeviceModel({ ID: row.ID })
      if (res.code === 0) {
        ElMessage.success('删除成功')
        await loadDeviceTree()
      }
    })
  }

  onMounted(async () => {
    await loadDeviceTree()
  })
</script>

<style scoped>
  .section-head {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }

  .section-title {
    color: #303133;
    font-size: 18px;
    font-weight: 600;
  }

  .section-subtitle {
    margin-top: 6px;
    color: #909399;
    font-size: 13px;
  }

  .tree-name-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .tree-seq-col {
    color: #909399;
    font-size: 13px;
    font-variant-numeric: tabular-nums;
    font-weight: 500;
  }

  .tree-name {
    color: #303133;
    font-weight: 500;
  }

  :deep(.device-tree-category-row .el-table__cell) {
    background: #f7fbff;
  }

  :deep(.device-tree-category-row .tree-name) {
    font-weight: 600;
  }

  :deep(.el-table__placeholder) {
    display: inline-block;
    width: 20px;
  }

  :deep(.el-table__indent) {
    width: 0 !important;
  }

  :deep(.el-table__expand-icon) {
    color: #606266;
    font-size: 14px;
  }
</style>
