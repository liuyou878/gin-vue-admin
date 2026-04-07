<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item label="类别名称">
          <el-input v-model="categorySearch.name" clearable placeholder="类别名称" />
        </el-form-item>
        <el-form-item label="型号名称">
          <el-input v-model="modelSearch.modelName" clearable placeholder="型号名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="loadDeviceTree">查询</el-button>
          <el-button icon="refresh" @click="resetDeviceTreeSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="section-head">
        <div>
          <div class="section-title">设备管理</div>
          <div class="section-subtitle">先维护设备类别，再在类别下面添加型号，形成树形结构。</div>
        </div>
        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openCategoryDialog()">新增类别</el-button>
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
            <div class="tree-name-cell" :class="{ 'tree-name-cell--child': scope.row.nodeType === 'model' }">
              <span class="tree-name">{{ scope.row.nodeType === 'category' ? scope.row.name : scope.row.modelName }}</span>
              <el-tag v-if="scope.row.nodeType === 'category'" type="info" size="small">类别</el-tag>
              <el-tag v-else type="success" size="small">型号</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="排序/代际" min-width="120">
          <template #default="scope">
            {{ scope.row.nodeType === 'category' ? scope.row.sort : (scope.row.generation || '-') }}
          </template>
        </el-table-column>
        <el-table-column label="备注" min-width="220" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.remark || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="300" fixed="right">
          <template #default="scope">
            <template v-if="scope.row.nodeType === 'category'">
              <el-button type="primary" link icon="plus" @click="openModelDialog(null, scope.row)">添加型号</el-button>
              <el-button type="primary" link icon="edit" @click="openCategoryDialog(scope.row)">编辑类别</el-button>
              <el-button type="primary" link icon="delete" @click="deleteCategoryRow(scope.row)">删除类别</el-button>
            </template>
            <template v-else>
              <el-button type="primary" link @click="enterFirmwareCenter(scope.row)">固件管理</el-button>
              <el-button type="primary" link icon="edit" @click="openModelDialog(scope.row)">编辑型号</el-button>
              <el-button type="primary" link icon="delete" @click="deleteModelRow(scope.row)">删除型号</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="gva-table-box">
      <div class="section-head section-head--compact">
        <div>
          <div class="section-title">型号固件管理</div>
          <div class="section-subtitle">
            {{ currentModel ? `当前型号：${currentModel.modelName}` : '从上面的树形设备管理里，点某个型号的“固件管理”进入。' }}
          </div>
        </div>
        <div v-if="currentModel" class="gva-btn-list">
          <el-button icon="refresh" @click="clearCurrentModel">清除当前型号</el-button>
        </div>
      </div>

      <div v-if="currentModel">
        <div class="firmware-context-card">
          <div class="context-item">
            <span class="context-label">当前型号</span>
            <span class="context-value">{{ currentModel.modelName }}</span>
          </div>
          <div class="context-item">
            <span class="context-label">设备类别</span>
            <span class="context-value">{{ currentModel.category?.name || currentCategoryName || '-' }}</span>
          </div>
          <div class="context-item">
            <span class="context-label">推荐版本</span>
            <span class="context-value">{{ recommendedFirmwareLabel || '暂未设置' }}</span>
          </div>
          <div class="context-item">
            <span class="context-label">已关联固件</span>
            <span class="context-value">{{ relationTableData.length }} 个</span>
          </div>
        </div>

        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openFirmwareDialog()">上传新固件</el-button>
          <el-button icon="link" @click="openBindExistingDialog">关联已有固件</el-button>
        </div>

        <el-table :data="relationTableData" row-key="ID">
            <el-table-column label="版本号" min-width="150">
              <template #default="scope">{{ resolveFirmware(scope.row).versionCode || '-' }}</template>
            </el-table-column>
            <el-table-column label="版本名称" min-width="180">
              <template #default="scope">{{ resolveFirmware(scope.row).versionName || '-' }}</template>
            </el-table-column>
            <el-table-column label="状态" width="110">
              <template #default="scope">
                <el-tag :type="firmwareStatusTag(resolveFirmware(scope.row).status)">
                  {{ firmwareStatusLabel(resolveFirmware(scope.row).status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="测试结果" width="110">
              <template #default="scope">
                <el-tag :type="relationTestTag(scope.row.testResult)">
                  {{ relationTestLabel(scope.row.testResult) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="推荐" width="90">
              <template #default="scope">
                <el-tag v-if="scope.row.isRecommended" type="warning">推荐中</el-tag>
                <span v-else>否</span>
              </template>
            </el-table-column>
            <el-table-column label="标签" min-width="180">
              <template #default="scope">
                <div class="tag-wrap">
                  <el-tag
                    v-for="tagRel in resolveFirmware(scope.row).tags || []"
                    :key="`${scope.row.ID}-${tagRel.tagId}`"
                    size="small"
                  >
                    {{ tagRel.tag?.tagName || '-' }}
                  </el-tag>
                  <span v-if="!(resolveFirmware(scope.row).tags || []).length">-</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="上传人" width="120">
              <template #default="scope">{{ resolveFirmware(scope.row).uploadedBy || '-' }}</template>
            </el-table-column>
            <el-table-column label="上传时间" width="180">
              <template #default="scope">{{ formatDate(resolveFirmware(scope.row).uploadedAt || resolveFirmware(scope.row).CreatedAt) }}</template>
            </el-table-column>
            <el-table-column label="操作" min-width="420" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openFirmwareDialog(scope.row)">编辑固件</el-button>
                <el-button type="primary" link @click="openStatusDialog(scope.row)">改状态</el-button>
                <el-button type="primary" link @click="openTestDialog(scope.row)">测结果</el-button>
                <el-button type="primary" link @click="setRecommended(scope.row)">设推荐</el-button>
                <el-button type="primary" link @click="openLogDrawer(scope.row)">日志</el-button>
                <el-button
                  v-if="resolveFirmware(scope.row).packageUrl"
                  type="primary"
                  link
                  @click="downloadPackage(scope.row)"
                >
                  下载
                </el-button>
                <el-button type="primary" link icon="delete" @click="deleteRelationRow(scope.row)">移除关联</el-button>
              </template>
            </el-table-column>
        </el-table>
      </div>

      <div v-else class="empty-panel">
        <el-empty description="请先在上面的设备树里添加类别和型号，然后从型号行进入固件管理" />
      </div>
    </div>

    <el-dialog v-model="categoryDialogVisible" :title="categoryDialogType === 'create' ? '新增设备类别' : '编辑设备类别'" width="520px">
      <el-form :model="categoryForm" label-width="90px">
        <el-form-item label="类别名称">
          <el-input v-model="categoryForm.name" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="categoryForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="categoryForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="categoryForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="modelDialogVisible" :title="modelDialogType === 'create' ? '新增设备型号' : '编辑设备型号'" width="620px">
      <el-form :model="modelForm" label-width="90px">
        <el-form-item label="设备类别">
          <el-input :model-value="selectedModelCategoryName" readonly />
        </el-form-item>
        <el-form-item label="型号名称">
          <el-input v-model="modelForm.modelName" />
        </el-form-item>
        <el-form-item label="代际">
          <el-input v-model="modelForm.generation" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="modelForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="modelForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modelDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitModel">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="firmwareDialogVisible" :title="firmwareDialogTitle" width="720px">
      <el-form :model="firmwareForm" label-width="100px">
        <el-form-item label="当前型号">
          <el-input :model-value="currentModel?.modelName || '-'" readonly />
        </el-form-item>
        <el-form-item label="版本号">
          <el-input v-model="firmwareForm.versionCode" />
        </el-form-item>
        <el-form-item label="版本名称">
          <el-input v-model="firmwareForm.versionName" />
        </el-form-item>
        <el-form-item label="固件包上传">
          <div class="upload-row">
            <el-upload
              :action="firmwareUploadAction"
              :headers="firmwareUploadHeaders"
              :show-file-list="false"
              :before-upload="beforeFirmwareUpload"
              :on-success="handleFirmwareUploadSuccess"
              :on-error="handleFirmwareUploadError"
            >
              <el-button type="primary" :loading="firmwareUploading">上传固件包</el-button>
            </el-upload>
            <span class="upload-tip">
              {{ firmwareUploading ? '正在上传到 MinIO...' : firmwareUploadName || '支持 update.bin_G3X... 这类自定义命名文件' }}
            </span>
          </div>
        </el-form-item>
        <el-form-item label="安装包地址">
          <el-input v-model="firmwareForm.packageUrl" readonly />
        </el-form-item>
        <el-form-item label="安装包名称">
          <el-input v-model="firmwareForm.packageName" readonly />
        </el-form-item>
        <el-form-item label="校验值">
          <el-input v-model="firmwareForm.checksum" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="firmwareForm.status" style="width: 100%">
            <el-option v-for="item in firmwareStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本标签">
          <el-select v-model="firmwareForm.tagIds" multiple clearable filterable style="width: 100%">
            <el-option v-for="item in tagOptions" :key="item.ID" :label="item.tagName" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="上传人">
          <el-input v-model="firmwareForm.uploadedBy" />
        </el-form-item>
        <el-form-item label="版本说明">
          <el-input v-model="firmwareForm.releaseNote" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="测试总结">
          <el-input v-model="firmwareForm.testSummary" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="firmwareDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitFirmware">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="bindDialogVisible" title="关联已有固件" width="560px">
      <el-form :model="bindForm" label-width="90px">
        <el-form-item label="当前型号">
          <el-input :model-value="currentModel?.modelName || '-'" readonly />
        </el-form-item>
        <el-form-item label="固件版本">
          <el-select v-model="bindForm.firmwareId" filterable style="width: 100%" placeholder="选择已有固件版本">
            <el-option
              v-for="item in availableBindFirmwareOptions"
              :key="item.ID"
              :label="`${item.versionCode}${item.versionName ? ` / ${item.versionName}` : ''}`"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBindExisting">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="statusDialogVisible" title="更新固件状态" width="520px">
      <el-form :model="statusForm" label-width="90px">
        <el-form-item label="目标状态">
          <el-select v-model="statusForm.status" style="width: 100%">
            <el-option v-for="item in firmwareStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作人">
          <el-input v-model="statusForm.operator" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="statusForm.content" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitStatus">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="testDialogVisible" title="更新测试结果" width="520px">
      <el-form :model="testForm" label-width="90px">
        <el-form-item label="测试结果">
          <el-select v-model="testForm.testResult" style="width: 100%">
            <el-option v-for="item in relationTestResultOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="测试人">
          <el-input v-model="testForm.tester" />
        </el-form-item>
        <el-form-item label="操作人">
          <el-input v-model="testForm.operator" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="testForm.content" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTestResult">确定</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="logDrawerVisible" :title="logDrawerTitle" size="720px">
      <el-table :data="logTableData">
        <el-table-column label="时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.operateAt || scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="动作" prop="action" width="140" />
        <el-table-column label="原状态" prop="fromStatus" width="120" />
        <el-table-column label="目标状态" prop="toStatus" width="120" />
        <el-table-column label="操作人" prop="operator" width="120" />
        <el-table-column label="说明" prop="content" min-width="180" show-overflow-tooltip />
      </el-table>
    </el-drawer>
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
  getDeviceModelList,
  createFirmwareVersion,
  updateFirmwareVersion,
  findFirmwareVersion,
  getFirmwareVersionList,
  changeFirmwareVersionStatus,
  createModelFirmwareRel,
  deleteModelFirmwareRel,
  getModelFirmwareRelList,
  setModelFirmwareRecommended,
  setModelFirmwareTestResult,
  getFirmwareTagList,
  setFirmwareTags,
  getFirmwareVersionLogList
} from '@/api/deviceFirmware'
import { formatDate, getBaseUrl } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({ name: 'DeviceFirmwareCenter' })

const userStore = useUserStore()

const categoryOptions = ref([])
const modelOptions = ref([])
const firmwareOptions = ref([])
const tagOptions = ref([])

const categorySearch = ref({})
const modelSearch = ref({})
const firmwareContext = ref({ categoryId: '', modelId: '' })

const categoryTableData = ref([])
const modelTableData = ref([])
const relationTableData = ref([])
const logTableData = ref([])

const categoryDialogVisible = ref(false)
const modelDialogVisible = ref(false)
const firmwareDialogVisible = ref(false)
const bindDialogVisible = ref(false)
const statusDialogVisible = ref(false)
const testDialogVisible = ref(false)
const logDrawerVisible = ref(false)

const categoryDialogType = ref('create')
const modelDialogType = ref('create')
const firmwareDialogType = ref('create')

const firmwareUploading = ref(false)
const firmwareUploadName = ref('')
const logDrawerTitle = ref('固件日志')

const firmwareStatusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '测试中', value: 'testing' },
  { label: '测试通过', value: 'tested_pass' },
  { label: '稳定版', value: 'stable' },
  { label: '已废弃', value: 'deprecated' }
]

const relationTestResultOptions = [
  { label: '待测试', value: 'pending' },
  { label: '测试中', value: 'testing' },
  { label: '已通过', value: 'passed' },
  { label: '未通过', value: 'failed' }
]

const categoryForm = ref({ name: '', code: '', sort: 0, status: 1, remark: '' })
const modelForm = ref({ categoryId: '', modelCode: '', modelName: '', seriesName: '', generation: '', status: 1, remark: '' })
const firmwareForm = ref({
  versionCode: '',
  versionName: '',
  packageUrl: '',
  packageName: '',
  checksum: '',
  status: 'draft',
  releaseNote: '',
  testSummary: '',
  uploadedBy: '',
  tagIds: []
})
const bindForm = ref({ firmwareId: '' })
const statusForm = ref({ id: '', status: 'draft', operator: '', content: '' })
const testForm = ref({ id: '', testResult: 'pending', tester: '', operator: '', content: '' })

const firmwareUploadAction = `${getBaseUrl()}/fileUploadAndDownload/upload`
const firmwareUploadHeaders = computed(() => ({
  'x-token': userStore.token,
  'x-user-id': userStore.userInfo?.ID || ''
}))

const defaultUploadedBy = () => userStore.userInfo?.nickName || userStore.userInfo?.userName || '系统用户'
const buildCategoryCode = () => `device-category-${Date.now()}`
const buildModelCode = () => `device-model-${Date.now()}`

const currentModel = computed(() => modelOptions.value.find((item) => item.ID === firmwareContext.value.modelId))
const currentCategoryName = computed(() => categoryOptions.value.find((item) => item.ID === firmwareContext.value.categoryId)?.name || '')
const selectedModelCategoryName = computed(() => {
  if (!modelForm.value.categoryId) {
    return ''
  }
  return categoryOptions.value.find((item) => item.ID === modelForm.value.categoryId)?.name || ''
})
const recommendedRelation = computed(() => relationTableData.value.find((item) => item.isRecommended))
const recommendedFirmwareLabel = computed(() => {
  if (!recommendedRelation.value) {
    return ''
  }
  const firmware = resolveFirmware(recommendedRelation.value)
  if (!firmware.versionCode) {
    return '-'
  }
  return `${firmware.versionCode}${firmware.versionName ? ` / ${firmware.versionName}` : ''}`
})
const selectedModelFirmwareIds = computed(() => relationTableData.value.map((item) => item.firmwareId))
const availableBindFirmwareOptions = computed(() => firmwareOptions.value.filter((item) => !selectedModelFirmwareIds.value.includes(item.ID)))
const deviceTreeData = computed(() => categoryTableData.value
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
  .filter((category) => !modelSearch.value.modelName || category.children.length > 0))
const firmwareDialogTitle = computed(() => {
  if (firmwareDialogType.value === 'update') {
    return '编辑固件版本'
  }
  return currentModel.value ? `为 ${currentModel.value.modelName} 上传新固件` : '新增固件版本'
})

const firmwareStatusLabel = (status) => firmwareStatusOptions.find((item) => item.value === status)?.label || status || '-'
const firmwareStatusTag = (status) => {
  if (status === 'stable' || status === 'tested_pass') {
    return 'success'
  }
  if (status === 'testing') {
    return 'warning'
  }
  if (status === 'deprecated') {
    return 'danger'
  }
  return 'info'
}
const relationTestLabel = (status) => relationTestResultOptions.find((item) => item.value === status)?.label || status || '-'
const relationTestTag = (status) => {
  if (status === 'passed') {
    return 'success'
  }
  if (status === 'testing') {
    return 'warning'
  }
  if (status === 'failed') {
    return 'danger'
  }
  return 'info'
}

const resolveFirmware = (relation) => {
  if (!relation) {
    return {}
  }
  return firmwareOptions.value.find((item) => item.ID === relation.firmwareId) || relation.firmware || {}
}

const loadDeviceTree = async () => {
  await Promise.all([loadCategories(), loadModels()])
}

const resetDeviceTreeSearch = async () => {
  categorySearch.value = {}
  modelSearch.value = {}
  await loadDeviceTree()
}

const resetFirmwareContext = () => {
  firmwareContext.value = { categoryId: '', modelId: '' }
  relationTableData.value = []
}

const clearCurrentModel = () => {
  resetFirmwareContext()
}

const getDeviceTreeRowClassName = ({ row }) => {
  return row.nodeType === 'category' ? 'device-tree-category-row' : 'device-tree-model-row'
}

const loadCategories = async () => {
  const res = await getDeviceCategoryList({ page: 1, pageSize: 999, ...categorySearch.value })
  if (res.code === 0) {
    categoryTableData.value = res.data.list || []
    categoryOptions.value = res.data.list || []
  }
}

const loadModels = async () => {
  const res = await getDeviceModelList({ page: 1, pageSize: 999, ...modelSearch.value })
  if (res.code === 0) {
    modelTableData.value = res.data.list || []
    modelOptions.value = res.data.list || []
    if (firmwareContext.value.modelId) {
      const stillExists = (res.data.list || []).some((item) => item.ID === firmwareContext.value.modelId)
      if (!stillExists) {
        firmwareContext.value.modelId = ''
        relationTableData.value = []
      }
    }
  }
}

const loadFirmwareOptions = async () => {
  const res = await getFirmwareVersionList({ page: 1, pageSize: 999 })
  if (res.code === 0) {
    firmwareOptions.value = res.data.list || []
  }
}

const loadRelationsForCurrentModel = async () => {
  if (!firmwareContext.value.modelId) {
    relationTableData.value = []
    return
  }
  const res = await getModelFirmwareRelList({ page: 1, pageSize: 999, modelId: firmwareContext.value.modelId })
  if (res.code === 0) {
    relationTableData.value = res.data.list || []
  }
}

const loadTags = async () => {
  const res = await getFirmwareTagList({ page: 1, pageSize: 999, status: 1 })
  if (res.code === 0) {
    tagOptions.value = res.data.list || []
  }
}

const enterFirmwareCenter = async (row) => {
  firmwareContext.value = {
    categoryId: row.categoryId,
    modelId: row.ID
  }
  await loadRelationsForCurrentModel()
}

const openCategoryDialog = async (row) => {
  categoryDialogType.value = row?.ID ? 'update' : 'create'
  if (row?.ID) {
    const res = await findDeviceCategory({ ID: row.ID })
    if (res.code === 0) {
      categoryForm.value = { ...res.data }
    }
  } else {
    categoryForm.value = { name: '', code: '', sort: 0, status: 1, remark: '' }
  }
  categoryDialogVisible.value = true
}

const submitCategory = async () => {
  const payload = {
    ...categoryForm.value,
    code: categoryForm.value.code || buildCategoryCode()
  }
  const res = categoryDialogType.value === 'create'
    ? await createDeviceCategory(payload)
    : await updateDeviceCategory(payload)
  if (res.code === 0) {
    ElMessage.success(categoryDialogType.value === 'create' ? '创建成功' : '更新成功')
    categoryDialogVisible.value = false
    await loadCategories()
  }
}

const deleteCategoryRow = (row) => {
  ElMessageBox.confirm(`确定删除类别 ${row.name} 吗？`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteDeviceCategory({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await loadCategories()
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
      generation: '',
      status: 1,
      remark: ''
    }
  }
  modelDialogVisible.value = true
}

const submitModel = async () => {
  const payload = {
    ...modelForm.value,
    modelCode: modelForm.value.modelCode || buildModelCode()
  }
  const res = modelDialogType.value === 'create'
    ? await createDeviceModel(payload)
    : await updateDeviceModel(payload)
  if (res.code === 0) {
    ElMessage.success(modelDialogType.value === 'create' ? '创建成功' : '更新成功')
    modelDialogVisible.value = false
    await loadModels()
  }
}

const deleteModelRow = (row) => {
  ElMessageBox.confirm(`确定删除型号 ${row.modelName} 吗？`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteDeviceModel({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await loadModels()
      if (firmwareContext.value.modelId === row.ID) {
        resetFirmwareContext()
      }
    }
  })
}

const beforeFirmwareUpload = (file) => {
  if (!file?.name) {
    ElMessage.error('未读取到文件名，请重新选择文件')
    return false
  }
  firmwareUploading.value = true
  firmwareUploadName.value = file.name
  return true
}

const handleFirmwareUploadSuccess = (res) => {
  firmwareUploading.value = false
  const file = res?.data?.file
  if (!file?.url) {
    ElMessage.error('上传成功，但未返回文件地址')
    return
  }
  firmwareForm.value.packageUrl = file.url
  firmwareForm.value.packageName = file.name || firmwareUploadName.value
  if (!firmwareForm.value.uploadedBy) {
    firmwareForm.value.uploadedBy = defaultUploadedBy()
  }
  ElMessage.success('固件包上传成功，已自动回填地址')
}

const handleFirmwareUploadError = () => {
  firmwareUploading.value = false
  ElMessage.error('固件包上传失败')
}

const openFirmwareDialog = async (row) => {
  if (!row && !currentModel.value) {
    ElMessage.warning('请先选择设备型号，再上传固件')
    return
  }
  firmwareDialogType.value = row ? 'update' : 'create'
  firmwareUploadName.value = ''
  firmwareUploading.value = false
  if (row) {
    const firmware = resolveFirmware(row)
    const res = await findFirmwareVersion({ ID: firmware.ID })
    if (res.code === 0) {
      firmwareForm.value = {
        ...res.data,
        tagIds: (res.data.tags || []).map((item) => item.tagId)
      }
      firmwareUploadName.value = res.data.packageName || ''
    }
  } else {
    firmwareForm.value = {
      versionCode: '',
      versionName: '',
      packageUrl: '',
      packageName: '',
      checksum: '',
      status: 'draft',
      releaseNote: '',
      testSummary: '',
      uploadedBy: defaultUploadedBy(),
      tagIds: []
    }
  }
  firmwareDialogVisible.value = true
}

const submitFirmware = async () => {
  const payload = {
    ...firmwareForm.value,
    uploadedAt: new Date()
  }
  const res = firmwareDialogType.value === 'create'
    ? await createFirmwareVersion(payload)
    : await updateFirmwareVersion(payload)
  if (res.code !== 0) {
    return
  }

  let firmwareId = firmwareForm.value.ID
  if (!firmwareId) {
    const listRes = await getFirmwareVersionList({ page: 1, pageSize: 1, versionCode: firmwareForm.value.versionCode })
    firmwareId = listRes.data?.list?.[0]?.ID
  }

  if (firmwareId) {
    await setFirmwareTags({ firmwareId, tagIds: firmwareForm.value.tagIds || [] })
  }

  if (firmwareDialogType.value === 'create' && firmwareId && currentModel.value) {
    const bindRes = await createModelFirmwareRel({
      modelId: currentModel.value.ID,
      firmwareId,
      isSupported: true,
      isRecommended: false,
      testResult: 'pending',
      tester: '',
      remark: ''
    })
    if (bindRes.code !== 0) {
      return
    }
  }

  ElMessage.success(firmwareDialogType.value === 'create' ? '上传并关联成功' : '更新成功')
  firmwareDialogVisible.value = false
  firmwareUploadName.value = ''
  await Promise.all([loadFirmwareOptions(), loadRelationsForCurrentModel()])
}

const openBindExistingDialog = () => {
  if (!currentModel.value) {
    ElMessage.warning('请先选择设备型号')
    return
  }
  if (!availableBindFirmwareOptions.value.length) {
    ElMessage.info('当前没有可关联的已有固件')
    return
  }
  bindForm.value = { firmwareId: '' }
  bindDialogVisible.value = true
}

const submitBindExisting = async () => {
  if (!bindForm.value.firmwareId || !currentModel.value) {
    ElMessage.warning('请选择要关联的固件版本')
    return
  }
  const res = await createModelFirmwareRel({
    modelId: currentModel.value.ID,
    firmwareId: bindForm.value.firmwareId,
    isSupported: true,
    isRecommended: false,
    testResult: 'pending',
    tester: '',
    remark: ''
  })
  if (res.code === 0) {
    ElMessage.success('关联成功')
    bindDialogVisible.value = false
    await loadRelationsForCurrentModel()
  }
}

const openStatusDialog = (row) => {
  const firmware = resolveFirmware(row)
  statusForm.value = {
    id: firmware.ID,
    status: firmware.status || 'draft',
    operator: defaultUploadedBy(),
    content: ''
  }
  statusDialogVisible.value = true
}

const submitStatus = async () => {
  const res = await changeFirmwareVersionStatus(statusForm.value)
  if (res.code === 0) {
    ElMessage.success('状态更新成功')
    statusDialogVisible.value = false
    await Promise.all([loadFirmwareOptions(), loadRelationsForCurrentModel()])
  }
}

const openLogDrawer = async (row) => {
  const firmware = resolveFirmware(row)
  const res = await getFirmwareVersionLogList({ page: 1, pageSize: 100, firmwareId: firmware.ID })
  if (res.code === 0) {
    logDrawerTitle.value = `固件日志 - ${firmware.versionCode || ''}`
    logTableData.value = res.data.list || []
    logDrawerVisible.value = true
  }
}

const openTestDialog = (row) => {
  testForm.value = {
    id: row.ID,
    testResult: row.testResult || 'pending',
    tester: row.tester || '',
    operator: defaultUploadedBy(),
    content: ''
  }
  testDialogVisible.value = true
}

const submitTestResult = async () => {
  const res = await setModelFirmwareTestResult(testForm.value)
  if (res.code === 0) {
    ElMessage.success('测试结果已更新')
    testDialogVisible.value = false
    await loadRelationsForCurrentModel()
  }
}

const setRecommended = (row) => {
  setModelFirmwareRecommended({
    id: row.ID,
    operator: defaultUploadedBy(),
    content: '前端设置推荐版本'
  }).then(async (res) => {
    if (res.code === 0) {
      ElMessage.success('设置成功')
      await loadRelationsForCurrentModel()
    }
  })
}

const deleteRelationRow = (row) => {
  const firmware = resolveFirmware(row)
  ElMessageBox.confirm(`确定把 ${firmware.versionCode || '该固件'} 从当前型号中移除吗？`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteModelFirmwareRel({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('移除成功')
      await loadRelationsForCurrentModel()
    }
  })
}

const downloadPackage = (row) => {
  const firmware = resolveFirmware(row)
  if (!firmware.packageUrl) {
    ElMessage.warning('当前固件还没有安装包地址')
    return
  }
  window.open(firmware.packageUrl, '_blank')
}

onMounted(async () => {
  await Promise.all([
    loadDeviceTree(),
    loadFirmwareOptions(),
    loadTags()
  ])
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

.section-head--compact {
  margin-bottom: 20px;
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

.context-tip {
  color: #909399;
  font-size: 13px;
  font-variant-numeric: tabular-nums;
}

.firmware-context-card {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.context-item {
  padding: 14px 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background: linear-gradient(180deg, #fbfdff 0%, #f5f9ff 100%);
}

.context-label {
  display: block;
  margin-bottom: 6px;
  color: #909399;
  font-size: 12px;
}

.context-value {
  color: #303133;
  font-size: 16px;
  font-weight: 600;
}

.tag-wrap {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.upload-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.upload-tip {
  color: #909399;
  font-size: 13px;
}

.empty-panel {
  padding: 48px 0;
  background: #fff;
  border-radius: 4px;
}
</style>
