<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="模板名称">
          <el-input
            v-model="searchInfo.name"
            placeholder="请输入名称"
            clearable
          />
        </el-form-item>
        <el-form-item label="型号">
          <el-input
            v-model="searchInfo.model"
            placeholder="请输入型号"
            clearable
          />
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
      <el-table-column prop="model" label="模板型号" width="120" />
      <el-table-column label="状态" width="80">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
            {{ scope.row.status === 1 ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="itemCount" label="检测项数" width="100" />
      <el-table-column prop="CreatedAt" label="创建时间" width="170">
        <template #default="scope">{{
          formatDate(scope.row.CreatedAt)
        }}</template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="scope">
          <el-button
            size="small"
            type="primary"
            link
            @click="openDialog('update', scope.row)"
            >编辑</el-button
          >
          <el-button size="small" type="success" link @click="onCopy(scope.row)"
            >复制</el-button
          >
          <el-button
            v-auth="btnAuth.delete"
            size="small"
            type="danger"
            link
            @click="onDelete(scope.row.ID)"
            >删除</el-button
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
      <el-form
        :model="formData"
        ref="formRef"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="formData.name" placeholder="如 G3X标准检测" />
        </el-form-item>
        <el-form-item label="产品名称" prop="productName">
          <el-input
            v-model="formData.productName"
            placeholder="如 GNSS接收机（RTK）"
          />
        </el-form-item>
        <el-form-item label="模板型号" prop="model">
          <el-input v-model="formData.model" placeholder="如 G3X" />
        </el-form-item>
        <el-form-item v-if="dialogType === 'update'" label="状态">
          <el-switch
            v-model="formData.status"
            :active-value="1"
            :inactive-value="2"
            active-text="启用"
            inactive-text="停用"
          />
        </el-form-item>

        <el-divider>检测项分配 & 排序</el-divider>

        <div class="mb-2">
          <el-button size="small" @click="showItemPicker = true"
            >从池中选择检测项</el-button
          >
        </div>

        <div class="drag-header">
          <span class="drag-col-sort">排序</span>
          <span class="drag-col-name">检测项名称</span>
          <span class="drag-col-type">类型</span>
          <span class="drag-col-unit">单位</span>
          <span class="drag-col-action">操作</span>
        </div>
        <draggable
          v-model="selectedItems"
          item-key="itemID"
          handle=".drag-handle"
          @end="renumberItems"
        >
          <template #item="{ element, index }">
            <div class="drag-row">
              <span class="drag-col-sort">
                <span class="drag-handle">⠿</span>
                <span class="ml-1">{{ index + 1 }}</span>
              </span>
              <span class="drag-col-name">{{ element.name }}</span>
              <span class="drag-col-type">
                <el-tag
                  v-if="element.resultType === 'pass_fail'"
                  type="success"
                  size="small"
                  >仅勾选</el-tag
                >
                <el-tag
                  v-else-if="element.resultType === 'number'"
                  type="primary"
                  size="small"
                  >仅数值</el-tag
                >
                <el-tag v-else type="warning" size="small">勾选+数值</el-tag>
              </span>
              <span class="drag-col-unit">{{ element.unit || '-' }}</span>
              <span class="drag-col-action">
                <el-button
                  size="small"
                  type="danger"
                  link
                  @click="removeItem(index)"
                  >移除</el-button
                >
              </span>
            </div>
          </template>
        </draggable>
      </el-form>

      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-drawer>

    <!-- Item picker dialog -->
    <el-dialog
      v-model="showItemPicker"
      title="选择检测项"
      width="500px"
      destroy-on-close
      @opened="onPickerOpened"
    >
      <el-table
        :data="allItems"
        border
        @selection-change="onPickerSelection"
        ref="pickerTableRef"
        max-height="400"
        row-key="ID"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="类型" width="100">
          <template #default="scope">
            <el-tag
              v-if="scope.row.resultType === 'pass_fail'"
              type="success"
              size="small"
              >仅勾选</el-tag
            >
            <el-tag
              v-else-if="scope.row.resultType === 'number'"
              type="primary"
              size="small"
              >仅数值</el-tag
            >
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
  import { useBtnAuth } from '@/utils/btnAuth'
  import draggable from 'vuedraggable'
  import {
    getTemplateList,
    createTemplate,
    copyTemplate,
    deleteTemplate,
    updateTemplate,
    findTemplate
  } from '@/plugin/inspection/api/template'
  import { getItemList } from '@/plugin/inspection/api/inspection_item'

  const btnAuth = useBtnAuth()
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
  const pickerTableRef = ref(null)

  const searchInfo = reactive({ name: '', model: '', page: 1, pageSize: 30 })
  const formData = reactive({
    ID: 0,
    name: '',
    productName: '',
    model: '',
    status: 1
  })

  const rules = {
    name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }]
  }

  const getList = async () => {
    loading.value = true
    try {
      const res = await getTemplateList(searchInfo)
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
    searchInfo.model = ''
    searchInfo.page = 1
    getList()
  }

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
      formData.ID = 0
      formData.name = ''
      formData.productName = ''
      formData.model = ''
      formData.status = 1
      selectedItems.value = []
    }
    drawerVisible.value = true
    await nextTick()
    formRef.value?.clearValidate()
  }

  const onPickerSelection = (val) => {
    pickerSelected.value = val
  }

  const onPickerOpened = () => {
    // Pre-select rows already in selectedItems
    const existingIDs = new Set(selectedItems.value.map((s) => s.itemID))
    allItems.value.forEach((row) => {
      if (existingIDs.has(row.ID)) {
        pickerTableRef.value?.toggleRowSelection(row, true)
      }
    })
  }

  const confirmPicker = () => {
    // Full sync: checked items are in, unchecked are out
    const newSelected = []
    pickerSelected.value.forEach((row) => {
      const existing = selectedItems.value.find((s) => s.itemID === row.ID)
      newSelected.push({
        itemID: row.ID,
        name: row.name,
        resultType: row.resultType,
        unit: row.unit || '',
        sort: existing ? existing.sort : 0
      })
    })
    selectedItems.value = newSelected
    renumberItems()
    showItemPicker.value = false
  }

  const removeItem = (index) => {
    selectedItems.value.splice(index, 1)
    renumberItems()
  }

  const renumberItems = () => {
    selectedItems.value.forEach((item, i) => {
      item.sort = i + 1
    })
  }

  const submitForm = async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return
    if (selectedItems.value.length === 0) {
      ElMessage.warning('请至少选择一个检测项')
      return
    }

    const payload = {
      name: formData.name,
      productName: formData.productName,
      model: formData.model,
      items: selectedItems.value.map((s) => ({
        itemID: s.itemID,
        sort: s.sort
      }))
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
    } catch (e) {
      /* handled by interceptor */
    }
  }

  const onDelete = (id) => {
    ElMessageBox.confirm(
      '删除模板会同时删除其关联的检测项配置，确定？',
      '提示',
      { type: 'warning' }
    )
      .then(async () => {
        const res = await deleteTemplate({ id })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          getList()
        }
      })
      .catch(() => {})
  }

  const onCopy = async (row) => {
    try {
      const { value } = await ElMessageBox.prompt(
        `请输入新模板名称`,
        `复制模板：${row.name}`,
      {
        confirmButtonText: '复制',
        cancelButtonText: '取消',
        inputPlaceholder: '请输入新模板名称',
        inputPattern: /\S+/,
        inputErrorMessage: '请输入模板名称'
      }
      )
      const res = await copyTemplate({ ID: row.ID, name: value.trim() })
      if (res.code === 0) {
        ElMessage.success('复制成功')
        getList()
      }
    } catch {
      // cancelled
    }
  }

  getList()
</script>

<style scoped>
  .drag-header,
  .drag-row {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid #eee;
  }
  .drag-header {
    font-weight: 600;
    color: #606266;
    background: #f5f7fa;
    border-radius: 4px 4px 0 0;
  }
  .drag-row {
    background: #fff;
  }
  .drag-row:hover {
    background: #f5f7fa;
  }
  .drag-handle {
    cursor: grab;
    color: #909399;
    font-size: 18px;
    user-select: none;
  }
  .drag-handle:active {
    cursor: grabbing;
  }
  .drag-col-sort {
    width: 80px;
    display: flex;
    align-items: center;
  }
  .drag-col-name {
    flex: 1;
  }
  .drag-col-type {
    width: 100px;
  }
  .drag-col-unit {
    width: 60px;
  }
  .drag-col-action {
    width: 70px;
  }
</style>
