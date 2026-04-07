<template>
  <div>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="设备类别" name="category">
        <div class="gva-search-box">
          <el-form :inline="true" :model="categorySearch">
            <el-form-item label="类别名称">
              <el-input v-model="categorySearch.name" clearable placeholder="类别名称" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="loadCategories">查询</el-button>
              <el-button icon="refresh" @click="resetCategorySearch">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openCategoryDialog()">新增类别</el-button>
          </div>
          <el-table :data="categoryTableData" row-key="ID">
            <el-table-column label="创建时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column label="类别名称" prop="name" min-width="160" />
            <el-table-column label="排序" prop="sort" width="90" />
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">{{ scope.row.status === 1 ? '启用' : '禁用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="备注" prop="remark" min-width="180" show-overflow-tooltip />
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openCategoryDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link icon="delete" @click="deleteCategoryRow(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="设备型号" name="model">
        <div class="gva-search-box">
          <el-form :inline="true" :model="modelSearch">
            <el-form-item label="设备类别">
              <el-select v-model="modelSearch.categoryId" clearable placeholder="全部" style="width: 180px">
                <el-option v-for="item in categoryOptions" :key="item.ID" :label="item.name" :value="item.ID" />
              </el-select>
            </el-form-item>
            <el-form-item label="型号名称">
              <el-input v-model="modelSearch.modelName" clearable placeholder="型号名称" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="loadModels">查询</el-button>
              <el-button icon="refresh" @click="resetModelSearch">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openModelDialog()">新增型号</el-button>
          </div>
          <el-table :data="modelTableData" row-key="ID">
            <el-table-column label="设备类别" min-width="140">
              <template #default="scope">{{ scope.row.category?.name || '-' }}</template>
            </el-table-column>
            <el-table-column label="型号名称" prop="modelName" min-width="160" />
            <el-table-column label="代际" prop="generation" min-width="100" />
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">{{ scope.row.status === 1 ? '启用' : '禁用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openModelDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link icon="delete" @click="deleteModelRow(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="固件管理" name="firmware">
        <div class="gva-search-box">
          <el-form :inline="true" :model="firmwareSearch">
            <el-form-item label="版本号">
              <el-input v-model="firmwareSearch.versionCode" clearable placeholder="版本号" />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="firmwareSearch.status" clearable placeholder="全部" style="width: 160px">
                <el-option v-for="item in firmwareStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="loadFirmwares">查询</el-button>
              <el-button icon="refresh" @click="resetFirmwareSearch">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openFirmwareDialog()">新增固件</el-button>
          </div>
          <el-table :data="firmwareTableData" row-key="ID">
            <el-table-column label="版本号" prop="versionCode" min-width="130" />
            <el-table-column label="版本名称" prop="versionName" min-width="150" />
            <el-table-column label="状态" width="120">
              <template #default="scope">
                <el-tag :type="firmwareStatusTag(scope.row.status)">{{ firmwareStatusLabel(scope.row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="标签" min-width="180">
              <template #default="scope">
                <div class="tag-wrap">
                  <el-tag v-for="tagRel in scope.row.tags || []" :key="`${scope.row.ID}-${tagRel.tagId}`" size="small">
                    {{ tagRel.tag?.tagName || '-' }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="上传人" prop="uploadedBy" width="120" />
            <el-table-column label="时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.uploadedAt || scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="300" fixed="right">
              <template #default="scope">
                <el-button type="primary" link @click="focusFirmwareRelations(scope.row)">查看关联</el-button>
                <el-button type="primary" link icon="edit" @click="openFirmwareDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link @click="openStatusDialog(scope.row)">改状态</el-button>
                <el-button type="primary" link @click="openLogDrawer(scope.row)">日志</el-button>
                <el-button type="primary" link icon="delete" @click="deleteFirmwareRow(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="gva-search-box">
          <div class="relation-header">
            <span class="relation-title">型号固件关系</span>
            <span class="relation-current">
              {{ currentFirmwareLabel ? `当前固件：${currentFirmwareLabel}` : '当前显示全部固件关系' }}
            </span>
            <el-button v-if="relationSearch.firmwareId" link type="primary" @click="clearFirmwareRelationFilter">清除固件筛选</el-button>
          </div>
          <el-form :inline="true" :model="relationSearch">
            <el-form-item label="设备型号">
              <el-select v-model="relationSearch.modelId" clearable filterable placeholder="全部" style="width: 220px">
                <el-option v-for="item in modelOptions" :key="item.ID" :label="item.modelName" :value="item.ID" />
              </el-select>
            </el-form-item>
            <el-form-item label="固件版本">
              <el-select v-model="relationSearch.firmwareId" clearable filterable placeholder="全部" style="width: 220px">
                <el-option v-for="item in firmwareOptions" :key="item.ID" :label="item.versionCode" :value="item.ID" />
              </el-select>
            </el-form-item>
            <el-form-item label="测试结果">
              <el-select v-model="relationSearch.testResult" clearable placeholder="全部" style="width: 160px">
                <el-option v-for="item in relationTestResultOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="loadRelations">查询</el-button>
              <el-button icon="refresh" @click="resetRelationSearch">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openRelationDialog()">新增关系</el-button>
          </div>
          <el-table :data="relationTableData" row-key="ID">
            <el-table-column label="设备型号" min-width="160">
              <template #default="scope">{{ scope.row.model?.modelName || '-' }}</template>
            </el-table-column>
            <el-table-column label="固件版本" min-width="140">
              <template #default="scope">{{ scope.row.firmware?.versionCode || '-' }}</template>
            </el-table-column>
            <el-table-column label="推荐" width="100">
              <template #default="scope">
                <el-tag v-if="scope.row.isRecommended" type="warning">推荐中</el-tag>
                <span v-else>否</span>
              </template>
            </el-table-column>
            <el-table-column label="测试结果" width="110">
              <template #default="scope">
                <el-tag :type="relationTestTag(scope.row.testResult)">{{ relationTestLabel(scope.row.testResult) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="测试人" prop="tester" width="120" />
            <el-table-column label="操作" width="280" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openRelationDialog(scope.row)">编辑</el-button>
                <el-button type="primary" link @click="setRecommended(scope.row)">设推荐</el-button>
                <el-button type="primary" link @click="openTestDialog(scope.row)">测结果</el-button>
                <el-button type="primary" link icon="delete" @click="deleteRelationRow(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="categoryDialogVisible" :title="categoryDialogType === 'create' ? '新增设备类别' : '编辑设备类别'" width="520px">
      <el-form :model="categoryForm" label-width="90px">
        <el-form-item label="类别名称"><el-input v-model="categoryForm.name" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="categoryForm.sort" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="categoryForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="categoryForm.remark" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="modelDialogVisible" :title="modelDialogType === 'create' ? '新增设备型号' : '编辑设备型号'" width="620px">
      <el-form :model="modelForm" label-width="90px">
        <el-form-item label="设备类别">
          <el-select v-model="modelForm.categoryId" style="width: 100%">
            <el-option v-for="item in categoryOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="型号名称"><el-input v-model="modelForm.modelName" /></el-form-item>
        <el-form-item label="代际"><el-input v-model="modelForm.generation" /></el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="modelForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="modelForm.remark" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modelDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitModel">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="firmwareDialogVisible" :title="firmwareDialogType === 'create' ? '新增固件版本' : '编辑固件版本'" width="720px">
      <el-form :model="firmwareForm" label-width="100px">
        <el-form-item label="版本号"><el-input v-model="firmwareForm.versionCode" /></el-form-item>
        <el-form-item label="版本名称"><el-input v-model="firmwareForm.versionName" /></el-form-item>
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
              {{ firmwareUploading ? '正在上传到 MinIO...' : firmwareUploadName || '支持标准后缀和 update.bin_G3X... 这类自定义命名' }}
            </span>
          </div>
        </el-form-item>
        <el-form-item label="安装包地址"><el-input v-model="firmwareForm.packageUrl" readonly /></el-form-item>
        <el-form-item label="安装包名称"><el-input v-model="firmwareForm.packageName" readonly /></el-form-item>
        <el-form-item label="校验值"><el-input v-model="firmwareForm.checksum" /></el-form-item>
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
        <el-form-item label="上传人"><el-input v-model="firmwareForm.uploadedBy" /></el-form-item>
        <el-form-item label="版本说明"><el-input v-model="firmwareForm.releaseNote" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="测试总结"><el-input v-model="firmwareForm.testSummary" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="firmwareDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitFirmware">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="relationDialogVisible" :title="relationDialogType === 'create' ? '新增关系' : '编辑关系'" width="620px">
      <el-form :model="relationForm" label-width="100px">
        <el-form-item label="设备型号">
          <el-select v-model="relationForm.modelId" filterable style="width: 100%">
            <el-option v-for="item in modelOptions" :key="item.ID" :label="item.modelName" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="固件版本">
          <el-select v-model="relationForm.firmwareId" filterable style="width: 100%">
            <el-option v-for="item in firmwareOptions" :key="item.ID" :label="item.versionCode" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否支持"><el-switch v-model="relationForm.isSupported" /></el-form-item>
        <el-form-item label="测试结果">
          <el-select v-model="relationForm.testResult" style="width: 100%">
            <el-option v-for="item in relationTestResultOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="测试人"><el-input v-model="relationForm.tester" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="relationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRelation">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="statusDialogVisible" title="更新固件状态" width="520px">
      <el-form :model="statusForm" label-width="90px">
        <el-form-item label="目标状态">
          <el-select v-model="statusForm.status" style="width: 100%">
            <el-option v-for="item in firmwareStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作人"><el-input v-model="statusForm.operator" /></el-form-item>
        <el-form-item label="说明"><el-input v-model="statusForm.content" type="textarea" :rows="3" /></el-form-item>
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
        <el-form-item label="测试人"><el-input v-model="testForm.tester" /></el-form-item>
        <el-form-item label="操作人"><el-input v-model="testForm.operator" /></el-form-item>
        <el-form-item label="说明"><el-input v-model="testForm.content" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTestResult">确定</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="logDrawerVisible" title="固件日志" size="720px">
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
  deleteFirmwareVersion,
  updateFirmwareVersion,
  findFirmwareVersion,
  getFirmwareVersionList,
  changeFirmwareVersionStatus,
  createModelFirmwareRel,
  deleteModelFirmwareRel,
  updateModelFirmwareRel,
  findModelFirmwareRel,
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
const activeTab = ref('category')
const categoryOptions = ref([])
const modelOptions = ref([])
const firmwareOptions = ref([])
const tagOptions = ref([])
const categorySearch = ref({})
const modelSearch = ref({})
const firmwareSearch = ref({})
const relationSearch = ref({})
const categoryTableData = ref([])
const modelTableData = ref([])
const firmwareTableData = ref([])
const relationTableData = ref([])
const logTableData = ref([])
const logDrawerVisible = ref(false)
const currentFirmwareId = ref('')

const categoryDialogVisible = ref(false)
const modelDialogVisible = ref(false)
const firmwareDialogVisible = ref(false)
const relationDialogVisible = ref(false)
const statusDialogVisible = ref(false)
const testDialogVisible = ref(false)
const firmwareUploading = ref(false)
const firmwareUploadName = ref('')

const categoryDialogType = ref('create')
const modelDialogType = ref('create')
const firmwareDialogType = ref('create')
const relationDialogType = ref('create')

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
const firmwareForm = ref({ versionCode: '', versionName: '', packageUrl: '', packageName: '', checksum: '', status: 'draft', releaseNote: '', testSummary: '', uploadedBy: '', tagIds: [] })
const relationForm = ref({ modelId: '', firmwareId: '', isSupported: true, isRecommended: false, testResult: 'pending', tester: '', remark: '' })
const statusForm = ref({ id: '', status: 'draft', operator: '', content: '' })
const testForm = ref({ id: '', testResult: 'pending', tester: '', operator: '', content: '' })
const firmwareUploadAction = `${getBaseUrl()}/fileUploadAndDownload/upload`
const firmwareUploadHeaders = computed(() => ({
  'x-token': userStore.token,
  'x-user-id': userStore.userInfo?.ID || ''
}))

const firmwareStatusLabel = (status) => firmwareStatusOptions.find((item) => item.value === status)?.label || status
const firmwareStatusTag = (status) => status === 'stable' || status === 'tested_pass' ? 'success' : status === 'testing' ? 'warning' : status === 'deprecated' ? 'danger' : 'info'
const relationTestLabel = (status) => relationTestResultOptions.find((item) => item.value === status)?.label || status
const relationTestTag = (status) => status === 'passed' ? 'success' : status === 'testing' ? 'warning' : status === 'failed' ? 'danger' : 'info'
const defaultUploadedBy = () => userStore.userInfo?.nickName || userStore.userInfo?.userName || '系统用户'
const buildCategoryCode = () => `device-category-${Date.now()}`
const buildModelCode = () => `device-model-${Date.now()}`
const currentFirmwareLabel = computed(() => {
  const current = firmwareOptions.value.find((item) => item.ID === relationSearch.value.firmwareId || item.ID === currentFirmwareId.value)
  if (!current) {
    return ''
  }
  return `${current.versionCode}${current.versionName ? ` / ${current.versionName}` : ''}`
})

const resetCategorySearch = () => { categorySearch.value = {}; loadCategories() }
const resetModelSearch = () => { modelSearch.value = {}; loadModels() }
const resetFirmwareSearch = () => { firmwareSearch.value = {}; loadFirmwares() }
const resetRelationSearch = () => {
  currentFirmwareId.value = ''
  relationSearch.value = {}
  loadRelations()
}
const clearFirmwareRelationFilter = () => {
  currentFirmwareId.value = ''
  relationSearch.value = { ...relationSearch.value, firmwareId: '' }
  loadRelations()
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
  }
}

const loadFirmwares = async () => {
  const res = await getFirmwareVersionList({ page: 1, pageSize: 999, ...firmwareSearch.value })
  if (res.code === 0) {
    firmwareTableData.value = res.data.list || []
    firmwareOptions.value = res.data.list || []
  }
}

const loadRelations = async () => {
  const res = await getModelFirmwareRelList({ page: 1, pageSize: 999, ...relationSearch.value })
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

const openCategoryDialog = async (row) => {
  categoryDialogType.value = row?.ID ? 'update' : 'create'
  if (row?.ID) {
    const res = await findDeviceCategory({ ID: row.ID })
    if (res.code === 0) categoryForm.value = { ...res.data }
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
  const res = categoryDialogType.value === 'create' ? await createDeviceCategory(payload) : await updateDeviceCategory(payload)
  if (res.code === 0) {
    ElMessage.success(categoryDialogType.value === 'create' ? '创建成功' : '更新成功')
    categoryDialogVisible.value = false
    loadCategories()
  }
}

const deleteCategoryRow = (row) => {
  ElMessageBox.confirm(`确定删除类别 ${row.name} 吗？`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteDeviceCategory({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadCategories()
    }
  })
}

const openModelDialog = async (row) => {
  modelDialogType.value = row?.ID ? 'update' : 'create'
  if (row?.ID) {
    const res = await findDeviceModel({ ID: row.ID })
    if (res.code === 0) modelForm.value = { ...res.data }
  } else {
    modelForm.value = { categoryId: '', modelCode: '', modelName: '', seriesName: '', generation: '', status: 1, remark: '' }
  }
  modelDialogVisible.value = true
}

const submitModel = async () => {
  const payload = {
    ...modelForm.value,
    modelCode: modelForm.value.modelCode || buildModelCode()
  }
  const res = modelDialogType.value === 'create' ? await createDeviceModel(payload) : await updateDeviceModel(payload)
  if (res.code === 0) {
    ElMessage.success(modelDialogType.value === 'create' ? '创建成功' : '更新成功')
    modelDialogVisible.value = false
    loadModels()
  }
}

const deleteModelRow = (row) => {
  ElMessageBox.confirm(`确定删除型号 ${row.modelName} 吗？`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteDeviceModel({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadModels()
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
  firmwareDialogType.value = row?.ID ? 'update' : 'create'
  firmwareUploadName.value = ''
  firmwareUploading.value = false
  if (row?.ID) {
    const res = await findFirmwareVersion({ ID: row.ID })
    if (res.code === 0) {
      firmwareForm.value = {
        ...res.data,
        tagIds: (res.data.tags || []).map((item) => item.tagId)
      }
      firmwareUploadName.value = res.data.packageName || ''
    }
  } else {
    firmwareForm.value = { versionCode: '', versionName: '', packageUrl: '', packageName: '', checksum: '', status: 'draft', releaseNote: '', testSummary: '', uploadedBy: defaultUploadedBy(), tagIds: [] }
  }
  firmwareDialogVisible.value = true
}

const submitFirmware = async () => {
  const payload = { ...firmwareForm.value, uploadedAt: new Date() }
  const res = firmwareDialogType.value === 'create' ? await createFirmwareVersion(payload) : await updateFirmwareVersion(payload)
  if (res.code !== 0) return
  let firmwareId = firmwareForm.value.ID
  if (!firmwareId) {
    const listRes = await getFirmwareVersionList({ page: 1, pageSize: 1, versionCode: firmwareForm.value.versionCode })
    firmwareId = listRes.data?.list?.[0]?.ID
  }
  if (firmwareId) {
    await setFirmwareTags({ firmwareId, tagIds: firmwareForm.value.tagIds || [] })
  }
  ElMessage.success(firmwareDialogType.value === 'create' ? '创建成功' : '更新成功')
  firmwareDialogVisible.value = false
  firmwareUploadName.value = ''
  loadFirmwares()
}

const openStatusDialog = (row) => {
  statusForm.value = { id: row.ID, status: row.status, operator: '', content: '' }
  statusDialogVisible.value = true
}

const submitStatus = async () => {
  const res = await changeFirmwareVersionStatus(statusForm.value)
  if (res.code === 0) {
    ElMessage.success('状态更新成功')
    statusDialogVisible.value = false
    loadFirmwares()
  }
}

const openLogDrawer = async (row) => {
  const res = await getFirmwareVersionLogList({ page: 1, pageSize: 100, firmwareId: row.ID })
  if (res.code === 0) {
    logTableData.value = res.data.list || []
    logDrawerVisible.value = true
  }
}

const deleteFirmwareRow = (row) => {
  ElMessageBox.confirm(`确定删除固件 ${row.versionCode} 吗？`, '提示', { type: 'warning' }).then(async () => {
    const res = await deleteFirmwareVersion({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (currentFirmwareId.value === row.ID) {
        clearFirmwareRelationFilter()
      }
      loadFirmwares()
    }
  })
}

const focusFirmwareRelations = (row) => {
  activeTab.value = 'firmware'
  currentFirmwareId.value = row.ID
  relationSearch.value = {
    ...relationSearch.value,
    firmwareId: row.ID
  }
  loadRelations()
}

const openRelationDialog = async (row) => {
  relationDialogType.value = row?.ID ? 'update' : 'create'
  if (row?.ID) {
    const res = await findModelFirmwareRel({ ID: row.ID })
    if (res.code === 0) relationForm.value = { ...res.data }
  } else {
    relationForm.value = { modelId: '', firmwareId: relationSearch.value.firmwareId || currentFirmwareId.value || '', isSupported: true, isRecommended: false, testResult: 'pending', tester: '', remark: '' }
  }
  relationDialogVisible.value = true
}

const submitRelation = async () => {
  const duplicatedRelation = relationTableData.value.find((item) => item.modelId === relationForm.value.modelId && item.firmwareId === relationForm.value.firmwareId && item.ID !== relationForm.value.ID)
  if (duplicatedRelation) {
    ElMessage.warning('该型号已经关联了这个固件版本，请直接编辑现有关系')
    return
  }
  const res = relationDialogType.value === 'create' ? await createModelFirmwareRel(relationForm.value) : await updateModelFirmwareRel(relationForm.value)
  if (res.code === 0) {
    ElMessage.success(relationDialogType.value === 'create' ? '创建成功' : '更新成功')
    relationDialogVisible.value = false
    loadRelations()
  }
}

const deleteRelationRow = (row) => {
  ElMessageBox.confirm('确定删除这条关系吗？', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteModelFirmwareRel({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadRelations()
    }
  })
}

const setRecommended = (row) => {
  setModelFirmwareRecommended({ id: row.ID, content: '前端设置推荐版本' }).then((res) => {
    if (res.code === 0) {
      ElMessage.success('设置成功')
      loadRelations()
    }
  })
}

const openTestDialog = (row) => {
  testForm.value = { id: row.ID, testResult: row.testResult || 'pending', tester: row.tester || '', operator: '', content: '' }
  testDialogVisible.value = true
}

const submitTestResult = async () => {
  const res = await setModelFirmwareTestResult(testForm.value)
  if (res.code === 0) {
    ElMessage.success('更新成功')
    testDialogVisible.value = false
    loadRelations()
  }
}

onMounted(async () => {
  await Promise.all([loadCategories(), loadModels(), loadFirmwares(), loadRelations(), loadTags()])
})
</script>

<style scoped>
.relation-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  color: #606266;
  flex-wrap: wrap;
}

.relation-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.relation-current {
  font-size: 13px;
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
</style>
