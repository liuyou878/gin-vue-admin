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
                    : scope.row.nodeType === 'model'
                    ? scope.row.modelName
                    : resolveFirmware(scope.row).versionCode ||
                      resolveFirmware(scope.row).versionName ||
                      '未命名固件'
                }}
              </span>
              <el-tag
                v-if="scope.row.nodeType === 'category'"
                type="info"
                size="small"
                >类别</el-tag
              >
              <el-tag
                v-else-if="scope.row.nodeType === 'model'"
                type="success"
                size="small"
                >型号</el-tag
              >
              <el-tag v-else type="warning" size="small">固件</el-tag>
              <el-tag
                v-if="
                  scope.row.nodeType === 'firmware' && scope.row.isRecommended
                "
                type="danger"
                size="small"
              >
                当前发布
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="260">
          <template #default="scope">
            <div v-if="scope.row.nodeType === 'firmware'" class="tag-wrap">
              <el-tag :type="firmwareStatusTag(resolveFirmware(scope.row).status)">
                {{ firmwareStatusLabel(resolveFirmware(scope.row).status) }}
              </el-tag>
              <el-tag
                :type="
                  firmwarePublishStatusTag(resolveFirmware(scope.row).publishStatus)
                "
              >
                {{
                  firmwarePublishStatusLabel(
                    resolveFirmware(scope.row).publishStatus
                  )
                }}
              </el-tag>
              <el-tag
                v-if="
                  resolveFirmware(scope.row).publishStatus === 'published' &&
                  resolveFirmware(scope.row).isLatest
                "
                type="danger"
              >
                最新版本
              </el-tag>
              <el-tag
                v-if="
                  resolveFirmware(scope.row).publishStatus === 'published' &&
                  resolveFirmware(scope.row).isStable
                "
                type="success"
              >
                稳定版本
              </el-tag>
              <el-tag
                v-if="isHistoryVersion(resolveFirmware(scope.row))"
                type="info"
              >
                历史版本
              </el-tag>
            </div>
            <el-tag v-else :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- <el-table-column label="排序/代际" min-width="120">
          <template #default="scope">
            {{
              scope.row.nodeType === 'category'
                ? scope.row.sort
                : scope.row.nodeType === 'model'
                  ? (scope.row.generation || '-')
                  : (scope.row.isRecommended ? '当前发布' : '-')
            }}
          </template>
        </el-table-column> -->
        <el-table-column label="备注" min-width="220" show-overflow-tooltip>
          <template #default="scope">
            {{
              scope.row.nodeType === 'firmware'
                ? (resolveFirmware(scope.row).status === 'test_failed'
                    ? latestFailureLogMap[resolveFirmware(scope.row).ID]?.content ||
                      resolveFirmware(scope.row).releaseNote
                    : resolveFirmware(scope.row).releaseNote) ||
                  resolveFirmware(scope.row).packageName ||
                  '-'
                : scope.row.remark || '-'
            }}
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
                icon="plus"
                @click="openFirmwareDialog(null, scope.row)"
                >添加固件</el-button
              >
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
            <template v-else>
              <el-button
                type="primary"
                link
                icon="edit"
                @click="openFirmwareDialog(scope.row)"
                >{{
                  resolveFirmware(scope.row).status === 'test_failed'
                    ? '上传修复包'
                    : '编辑固件'
                }}</el-button
              >
              <el-button
                v-if="canStartTesting(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="changeFirmwareStage(scope.row, 'testing', '开始测试')"
                >开始测试</el-button
              >
              <el-button
                v-if="canSubmitTestResult(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="openTestResultDialog(scope.row)"
                >测试结果</el-button
              >
              <el-button
                v-if="canSubmitRelease(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="
                  changeFirmwareStage(scope.row, 'pending_release', '提交发布')
                "
                >提交发布</el-button
              >
              <el-button
                v-if="canRejectRelease(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="changeFirmwareStage(scope.row, 'testing', '驳回到测试中')"
                >驳回</el-button
              >
              <el-button
                v-if="canPublish(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="publishFirmware(scope.row)"
                >发布</el-button
              >
              <el-button
                v-if="canToggleStable(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="toggleFirmwareStable(scope.row)"
              >
                {{ resolveFirmware(scope.row).isStable ? '取消稳定' : '标记稳定' }}
              </el-button>
              <el-button
                v-if="canSetCurrentRelease(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="setCurrentRelease(scope.row)"
                >设为当前发布</el-button
              >
              <el-button
                v-if="canVoid(resolveFirmware(scope.row))"
                type="primary"
                link
                @click="voidFirmware(scope.row)"
                >作废</el-button
              >
              <el-button type="primary" link @click="openLogDrawer(scope.row)"
                >日志</el-button
              >
              <el-button
                v-if="resolveFirmware(scope.row).packageUrl"
                type="primary"
                link
                @click="downloadPackage(scope.row)"
              >
                下载
              </el-button>
              <el-button
                v-if="canDeleteFirmwareRelation(resolveFirmware(scope.row))"
                type="primary"
                link
                icon="delete"
                @click="deleteRelationRow(scope.row)"
                >移除</el-button
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

    <el-dialog
      v-model="modelDialogVisible"
      :title="modelDialogType === 'create' ? '新增设备型号' : '编辑设备型号'"
      width="620px"
    >
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

    <el-dialog
      v-model="firmwareDialogVisible"
      :title="firmwareDialogTitle"
      width="720px"
    >
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
              :disabled="!canEditFirmwarePackage(firmwareForm)"
              :before-upload="beforeFirmwareUpload"
              :on-success="handleFirmwareUploadSuccess"
              :on-error="handleFirmwareUploadError"
            >
              <el-button type="primary" :loading="firmwareUploading"
                :disabled="!canEditFirmwarePackage(firmwareForm)"
                >上传固件包</el-button
              >
            </el-upload>
            <span class="upload-tip">
              {{
                firmwareUploading
                  ? '正在上传到 MinIO...'
                  : firmwareUploadName ||
                    '支持 update.bin_G3X... 这类自定义命名文件'
              }}
            </span>
          </div>
        </el-form-item>
        <el-form-item label="安装包地址">
          <el-input v-model="firmwareForm.packageUrl" readonly />
        </el-form-item>
        <!-- <el-form-item label="安装包名称">
          <el-input v-model="firmwareForm.packageName" readonly />
        </el-form-item>
        <el-form-item label="校验值">
          <el-input v-model="firmwareForm.checksum" />
        </el-form-item> -->
        <el-form-item label="开发状态">
          <el-input :model-value="firmwareStatusLabel(firmwareForm.status || 'pending_test')" readonly />
        </el-form-item>
        <el-form-item v-if="firmwareDialogType === 'create'" label="流程说明">
          <el-input model-value="新建后自动进入待测试，后续通过按钮推进流程" readonly />
        </el-form-item>
        <el-form-item label="上传人">
          <el-input v-model="firmwareForm.uploadedBy" />
        </el-form-item>
        <el-form-item v-if="firmwareDialogType === 'update'" label="发布状态">
          <el-input
            :model-value="firmwarePublishStatusLabel(firmwareForm.publishStatus)"
            readonly
          />
        </el-form-item>
        <el-form-item label="版本说明">
          <el-input
            v-model="firmwareForm.releaseNote"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="firmwareDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitFirmware">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="testResultDialogVisible" title="提交测试结果" width="560px">
      <el-form :model="testResultForm" label-width="90px">
        <el-form-item label="测试结果">
          <el-radio-group v-model="testResultForm.result">
            <el-radio value="tested_pass">通过</el-radio>
            <el-radio value="test_failed">不通过</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="testResultForm.result === 'test_failed'" label="原因分类">
          <el-select
            v-model="testResultForm.reasonTypes"
            multiple
            collapse-tags
            collapse-tags-tooltip
            style="width: 100%"
            placeholder="请选择未通过原因"
          >
            <el-option
              v-for="item in failReasonOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="说明">
          <el-input
            v-model="testResultForm.description"
            type="textarea"
            :rows="4"
            :placeholder="
              testResultForm.result === 'test_failed'
                ? '请填写未通过原因或先贴腾讯文档链接'
                : '可选，记录测试通过说明'
            "
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testResultDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTestResult">确定</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="logDrawerVisible" :title="logDrawerTitle" size="720px">
      <el-table :data="logTableData">
        <el-table-column label="时间" width="180">
          <template #default="scope">{{
            formatDate(scope.row.operateAt || scope.row.CreatedAt)
          }}</template>
        </el-table-column>
        <el-table-column label="动作" min-width="180">
          <template #default="scope">{{ firmwareLogActionLabel(scope.row) }}</template>
        </el-table-column>
        <el-table-column label="原状态" width="120">
          <template #default="scope">{{ firmwareStatusLabel(scope.row.fromStatus) }}</template>
        </el-table-column>
        <el-table-column label="目标状态" width="120">
          <template #default="scope">{{ firmwareStatusLabel(scope.row.toStatus) }}</template>
        </el-table-column>
        <el-table-column label="操作人" prop="operator" width="120" />
        <el-table-column
          label="说明"
          prop="content"
          min-width="180"
          show-overflow-tooltip
        />
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
    publishFirmwareVersion,
    setFirmwareStable,
    voidFirmwareVersion,
    createModelFirmwareRel,
    deleteModelFirmwareRel,
    getModelFirmwareRelList,
    setModelFirmwareRecommended,
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

  const categorySearch = ref({})
  const modelSearch = ref({})
  const firmwareContext = ref({ categoryId: '', modelId: '' })

  const categoryTableData = ref([])
  const modelTableData = ref([])
  const relationTableData = ref([])
  const logTableData = ref([])
  const firmwareLogList = ref([])

  const categoryDialogVisible = ref(false)
  const modelDialogVisible = ref(false)
  const firmwareDialogVisible = ref(false)
  const testResultDialogVisible = ref(false)
  const logDrawerVisible = ref(false)

  const categoryDialogType = ref('create')
  const modelDialogType = ref('create')
  const firmwareDialogType = ref('create')

  const firmwareUploading = ref(false)
  const firmwareUploadName = ref('')
  const logDrawerTitle = ref('固件日志')
  const currentTestResultRow = ref(null)
  const failReasonOptions = ['有Bug', '少功能', '需优化']

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
    generation: '',
    status: 1,
    remark: ''
  })
  const firmwareForm = ref({
    versionCode: '',
    versionName: '',
    packageUrl: '',
    packageName: '',
    checksum: '',
    status: 'pending_test',
    publishStatus: 'unpublished',
    releaseNote: '',
    uploadedBy: ''
  })
  const testResultForm = ref({
    result: 'tested_pass',
    reasonTypes: [],
    description: ''
  })

  const firmwareUploadAction = `${getBaseUrl()}/fileUploadAndDownload/upload`
  const firmwareUploadHeaders = computed(() => ({
    'x-token': userStore.token,
    'x-user-id': userStore.userInfo?.ID || ''
  }))

  const defaultUploadedBy = () =>
    userStore.userInfo?.nickName || userStore.userInfo?.userName || '系统用户'
  const buildCategoryCode = () => `device-category-${Date.now()}`
  const buildModelCode = () => `device-model-${Date.now()}`

  const currentModel = computed(() =>
    modelOptions.value.find((item) => item.ID === firmwareContext.value.modelId)
  )
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
            displayOrder: `${modelIndex + 1}`,
            children: relationTableData.value
              .filter((relation) => relation.modelId === model.ID)
              .map((relation, relationIndex) => ({
                ...relation,
                nodeType: 'firmware',
                rowKey: `firmware-${relation.ID}`,
                displayOrder: `${relationIndex + 1}`
              }))
          }))
      }))
      .filter(
        (category) =>
          !modelSearch.value.modelName || category.children.length > 0
      )
  )
  const firmwareDialogTitle = computed(() => {
    if (firmwareDialogType.value === 'update') {
      return firmwareForm.value.status === 'test_failed' ? '上传修复包' : '编辑固件版本'
    }
    return currentModel.value
      ? `为 ${currentModel.value.modelName} 上传新固件`
      : '新增固件版本'
  })

  const firmwareStatusLabel = (status) =>
    ({
      pending_test: '待测试',
      pending: '待测试',
      draft: '待测试',
      testing: '测试中',
      passed: '测试通过',
      tested_pass: '测试通过',
      test_failed: '测试不通过',
      failed: '测试不通过',
      pending_release: '待发布',
      unpublished: '未发布',
      published: '已发布',
      voided: '已作废'
    }[status] ||
    status ||
    '-')
  const firmwarePublishStatusLabel = (status) =>
    ({
      unpublished: '未发布',
      published: '已发布',
      voided: '已作废'
    }[status] || '-')
  const firmwareStatusTag = (status) => {
    if (status === 'tested_pass') {
      return 'success'
    }
    if (status === 'testing') {
      return 'warning'
    }
    if (status === 'test_failed') {
      return 'danger'
    }
    return 'info'
  }
  const firmwarePublishStatusTag = (status) => {
    if (status === 'published') {
      return 'success'
    }
    if (status === 'voided') {
      return 'danger'
    }
    return 'info'
  }
  const firmwareLogActionLabel = (log) => {
    const modelName = log?.model?.modelName
    const withModel = (prefix, fallback) => (modelName ? `${prefix}${modelName}` : fallback)
    return ({
      upload: '上传固件',
      bind_model: withModel('绑定到 ', '绑定型号'),
      start_testing: '开始测试',
      test_pass: '测试通过',
      test_fail: '测试未通过',
      fix_upload: '已修复并重新上传',
      submit_release: '提交发布',
      reject_release: '驳回到测试中',
      publish: '发布版本',
      mark_stable: '标记稳定版本',
      unmark_stable: '取消稳定版本',
      void_release: '作废发布版本',
      set_recommended: modelName ? `设为 ${modelName} 当前发布` : '设为当前发布'
    }[log?.action] || log?.action || '-')
  }
  const isHistoryVersion = (firmware) =>
    firmware?.publishStatus === 'published' &&
    !firmware?.isLatest &&
    !firmware?.isStable
  const canStartTesting = (firmware) =>
    ['pending_test', 'test_failed'].includes(firmware?.status) &&
    !['published', 'voided'].includes(firmware?.publishStatus)
  const canSubmitTestResult = (firmware) => firmware?.status === 'testing'
  const canSubmitRelease = (firmware) =>
    firmware?.status === 'tested_pass' && firmware?.publishStatus === 'unpublished'
  const canRejectRelease = (firmware) =>
    firmware?.status === 'pending_release' &&
    firmware?.publishStatus === 'unpublished'
  const canPublish = (firmware) =>
    firmware?.status === 'pending_release' &&
    firmware?.publishStatus === 'unpublished'
  const canToggleStable = (firmware) => firmware?.publishStatus === 'published'
  const canVoid = (firmware) => firmware?.publishStatus === 'published'
  const canSetCurrentRelease = (firmware) =>
    firmware?.publishStatus === 'published'
  const canDeleteFirmwareRelation = (firmware) =>
    !['published', 'voided'].includes(firmware?.publishStatus)
  const canEditFirmwarePackage = (firmware) =>
    !firmware?.ID ||
    (firmware?.publishStatus === 'unpublished' &&
      ['pending_test', 'test_failed'].includes(firmware?.status))
  const latestFailureLogMap = computed(() => {
    const map = {}
    ;(firmwareLogList.value || []).forEach((log) => {
      if (log.action !== 'test_fail' || !log.firmwareId) {
        return
      }
      const current = map[log.firmwareId]
      const currentTime = current
        ? new Date(current.operateAt || current.CreatedAt || 0).getTime()
        : 0
      const nextTime = new Date(log.operateAt || log.CreatedAt || 0).getTime()
      if (!current || nextTime >= currentTime) {
        map[log.firmwareId] = log
      }
    })
    return map
  })

  const resolveFirmware = (relation) => {
    if (!relation) {
      return {}
    }
    return (
      firmwareOptions.value.find((item) => item.ID === relation.firmwareId) ||
      relation.firmware ||
      {}
    )
  }

  const loadDeviceTree = async () => {
    await Promise.all([
      loadCategories(),
      loadModels(),
      loadFirmwareOptions(),
      loadAllRelations(),
      loadFirmwareLogs()
    ])
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
    if (row.nodeType === 'model') {
      return 'device-tree-model-row'
    }
    return 'device-tree-firmware-row'
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
      if (firmwareContext.value.modelId) {
        const stillExists = (res.data.list || []).some(
          (item) => item.ID === firmwareContext.value.modelId
        )
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

  const loadAllRelations = async () => {
    const res = await getModelFirmwareRelList({ page: 1, pageSize: 999 })
    if (res.code === 0) {
      relationTableData.value = res.data.list || []
    }
  }

  const loadFirmwareLogs = async () => {
    const res = await getFirmwareVersionLogList({ page: 1, pageSize: 999 })
    if (res.code === 0) {
      firmwareLogList.value = res.data.list || []
    }
  }

  const setFirmwareContextByModel = (row) => {
    firmwareContext.value = {
      categoryId: row.categoryId,
      modelId: row.ID
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
        if (firmwareContext.value.modelId === row.ID) {
          firmwareContext.value = { categoryId: '', modelId: '' }
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

  const openFirmwareDialog = async (row, parentModel) => {
    if (parentModel) {
      setFirmwareContextByModel(parentModel)
    } else if (row?.nodeType === 'firmware') {
      const relationModel = modelOptions.value.find(
        (item) => item.ID === row.modelId
      )
      if (relationModel) {
        setFirmwareContextByModel(relationModel)
      }
    }
    if (!row && !currentModel.value) {
      ElMessage.warning('请先选择设备型号，再上传固件')
      return
    }
    firmwareDialogType.value =
      row?.nodeType === 'firmware' ? 'update' : 'create'
    firmwareUploadName.value = ''
    firmwareUploading.value = false
    if (row?.nodeType === 'firmware') {
      const firmware = resolveFirmware(row)
      const res = await findFirmwareVersion({ ID: firmware.ID })
      if (res.code === 0) {
        firmwareForm.value = {
          ...res.data
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
        status: 'pending_test',
        publishStatus: 'unpublished',
        releaseNote: '',
        uploadedBy: defaultUploadedBy()
      }
    }
    firmwareDialogVisible.value = true
  }

  const submitFirmware = async () => {
    const isFixUpload =
      firmwareDialogType.value === 'update' && firmwareForm.value.status === 'test_failed'
    const payload = {
      ...firmwareForm.value,
      uploadedAt: new Date()
    }
    const res =
      firmwareDialogType.value === 'create'
        ? await createFirmwareVersion(payload)
        : await updateFirmwareVersion(payload)
    if (res.code !== 0) {
      return
    }

    let firmwareId = firmwareForm.value.ID
    if (!firmwareId) {
      const listRes = await getFirmwareVersionList({
        page: 1,
        pageSize: 1,
        versionCode: firmwareForm.value.versionCode
      })
      firmwareId = listRes.data?.list?.[0]?.ID
    }

    if (
      firmwareDialogType.value === 'create' &&
      firmwareId &&
      currentModel.value
    ) {
      const bindRes = await createModelFirmwareRel({
        modelId: currentModel.value.ID,
        firmwareId,
        isSupported: true,
        isRecommended: false,
        testResult: '',
        tester: '',
        remark: ''
      })
      if (bindRes.code !== 0) {
        return
      }
    }

    ElMessage.success(
      firmwareDialogType.value === 'create'
        ? '上传并关联成功'
        : isFixUpload
          ? '修复包已更新，已进入待测试'
          : '更新成功'
    )
    firmwareDialogVisible.value = false
    firmwareUploadName.value = ''
    await loadDeviceTree()
  }

  const changeFirmwareStage = async (row, status, content) => {
    const firmware = resolveFirmware(row)
    const res = await changeFirmwareVersionStatus({
      id: firmware.ID,
      status,
      operator: defaultUploadedBy(),
      content
    })
    if (res.code === 0) {
      ElMessage.success('状态更新成功')
      await loadDeviceTree()
      return true
    }
    return false
  }

  const openTestResultDialog = (row) => {
    currentTestResultRow.value = row
    testResultForm.value = {
      result: 'tested_pass',
      reasonTypes: [],
      description: ''
    }
    testResultDialogVisible.value = true
  }

  const submitTestResult = async () => {
    const row = currentTestResultRow.value
    if (!row) {
      return
    }
    const result = testResultForm.value.result
    if (
      result === 'test_failed' &&
      !testResultForm.value.reasonTypes.length &&
      !testResultForm.value.description?.trim()
    ) {
      ElMessage.warning('测试不通过时，请至少选择一个原因或填写说明')
      return
    }
    const reasonText =
      result === 'test_failed' && testResultForm.value.reasonTypes.length
        ? `原因分类：${testResultForm.value.reasonTypes.join('、')}`
        : ''
    const descriptionText = testResultForm.value.description?.trim() || ''
    const content = [reasonText, descriptionText].filter(Boolean).join('；') ||
      (result === 'tested_pass' ? '测试通过' : '测试不通过')
    const res = await changeFirmwareStage(row, result, content)
    if (res !== false) {
      testResultDialogVisible.value = false
      currentTestResultRow.value = null
    }
  }

  const publishFirmware = async (row) => {
    const firmware = resolveFirmware(row)
    const res = await publishFirmwareVersion({
      id: firmware.ID,
      operator: defaultUploadedBy(),
      content: '发布进入发布版本'
    })
    if (res.code === 0) {
      ElMessage.success('发布成功，已进入发布版本')
      await loadDeviceTree()
    }
  }

  const toggleFirmwareStable = async (row) => {
    const firmware = resolveFirmware(row)
    const stable = !firmware.isStable
    const res = await setFirmwareStable({
      id: firmware.ID,
      stable,
      operator: defaultUploadedBy(),
      content: stable ? '手动标记为稳定版本' : '取消稳定版本标记'
    })
    if (res.code === 0) {
      ElMessage.success(stable ? '已标记为稳定版本' : '已取消稳定版本')
      await loadDeviceTree()
    }
  }

  const openLogDrawer = async (row) => {
    const firmware = resolveFirmware(row)
    const res = await getFirmwareVersionLogList({
      page: 1,
      pageSize: 100,
      firmwareId: firmware.ID
    })
    if (res.code === 0) {
      logDrawerTitle.value = `固件日志 - ${firmware.versionCode || ''}`
      logTableData.value = res.data.list || []
      logDrawerVisible.value = true
    }
  }

  const setCurrentRelease = (row) => {
    setModelFirmwareRecommended({
      id: row.ID,
      operator: defaultUploadedBy(),
      content: '设为当前发布版'
    }).then(async (res) => {
      if (res.code === 0) {
        ElMessage.success('已设为当前发布版')
        await loadDeviceTree()
      }
    })
  }

  const voidFirmware = (row) => {
    const firmware = resolveFirmware(row)
    ElMessageBox.prompt('请输入作废原因', '作废版本', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '作废原因',
      inputValue: ''
    }).then(async ({ value }) => {
      const res = await voidFirmwareVersion({
        id: firmware.ID,
        operator: defaultUploadedBy(),
        voidReason: value,
        content: value || '作废已发布版本'
      })
      if (res.code === 0) {
        ElMessage.success('版本已作废')
        await loadDeviceTree()
      }
    })
  }

  const deleteRelationRow = (row) => {
    const firmware = resolveFirmware(row)
    ElMessageBox.confirm(
      `确定把 ${firmware.versionCode || '该固件'} 从当前型号中移除吗？`,
      '提示',
      { type: 'warning' }
    ).then(async () => {
      const res = await deleteModelFirmwareRel({ ID: row.ID })
      if (res.code === 0) {
        ElMessage.success('移除成功')
        await loadDeviceTree()
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
    await Promise.all([loadDeviceTree()])
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
