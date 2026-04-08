<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item label="设备类别">
          <el-select
            v-model="searchForm.categoryId"
            clearable
            filterable
            placeholder="设备类别"
            style="width: 180px"
          >
            <el-option
              v-for="item in categoryOptions"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="设备型号">
          <el-select
            v-model="searchForm.modelId"
            clearable
            filterable
            placeholder="设备型号"
            style="width: 220px"
          >
            <el-option
              v-for="item in filteredSearchModels"
              :key="item.ID"
              :label="item.modelName"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="开发状态">
          <el-select
            v-model="searchForm.status"
            clearable
            placeholder="开发状态"
            style="width: 160px"
          >
            <el-option
              v-for="item in devStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发布状态">
          <el-select
            v-model="searchForm.publishStatus"
            clearable
            placeholder="发布状态"
            style="width: 160px"
          >
            <el-option
              v-for="item in publishStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="版本关键字">
          <el-input
            v-model="searchForm.keyword"
            clearable
            placeholder="版本号/版本名称"
            style="width: 220px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="loadPage"
            >查询</el-button
          >
          <el-button icon="refresh" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="section-head">
        <div>
          <div class="section-title">固件流程管理</div>
          <div class="section-subtitle">
            独立管理开发版本与发布版本，新增固件时直接选择设备类别和型号。
          </div>
        </div>
        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openFirmwareDialog()"
            >新增固件</el-button
          >
        </div>
      </div>

      <el-table :data="filteredRows" row-key="ID">
        <el-table-column label="设备类别" min-width="120">
          <template #default="scope">{{
            scope.row.model?.category?.name || '-'
          }}</template>
        </el-table-column>
        <el-table-column label="设备型号" min-width="140">
          <template #default="scope">{{
            scope.row.model?.modelName || '-'
          }}</template>
        </el-table-column>
        <el-table-column label="版本号" min-width="130">
          <template #default="scope">{{
            firmwareOf(scope.row).versionCode || '-'
          }}</template>
        </el-table-column>
        <el-table-column label="版本名称" min-width="170" show-overflow-tooltip>
          <template #default="scope">{{
            firmwareOf(scope.row).versionName || '-'
          }}</template>
        </el-table-column>
        <el-table-column label="状态" min-width="300">
          <template #default="scope">
            <div class="tag-wrap">
              <el-tag :type="devStatusTag(firmwareOf(scope.row).status)">{{
                devStatusLabel(firmwareOf(scope.row).status)
              }}</el-tag>
              <el-tag
                :type="publishStatusTag(firmwareOf(scope.row).publishStatus)"
                >{{
                  publishStatusLabel(firmwareOf(scope.row).publishStatus)
                }}</el-tag
              >
              <el-tag
                v-if="
                  firmwareOf(scope.row).publishStatus === 'published' &&
                  firmwareOf(scope.row).isLatest
                "
                type="danger"
                >最新版本</el-tag
              >
              <el-tag
                v-if="
                  firmwareOf(scope.row).publishStatus === 'published' &&
                  firmwareOf(scope.row).isStable
                "
                type="success"
                >稳定版本</el-tag
              >
              <el-tag v-if="isHistoryVersion(firmwareOf(scope.row))" type="info"
                >历史版本</el-tag
              >
              <el-tag v-if="scope.row.isRecommended" type="warning"
                >当前发布</el-tag
              >
            </div>
          </template>
        </el-table-column>
        <el-table-column label="说明" min-width="260" show-overflow-tooltip>
          <template #default="scope">{{ rowNote(scope.row) }}</template>
        </el-table-column>
        <el-table-column label="操作" min-width="300" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              @click="openFirmwareDialog(scope.row)"
              >{{
                firmwareOf(scope.row).status === 'test_failed'
                  ? '上传修复包'
                  : '编辑固件'
              }}</el-button
            >
            <el-button
              v-if="canStartTesting(firmwareOf(scope.row))"
              type="primary"
              link
              @click="changeStage(scope.row, 'testing', '开始测试')"
              >开始测试</el-button
            >
            <el-button
              v-if="canSubmitTestResult(firmwareOf(scope.row))"
              type="primary"
              link
              @click="openTestResultDialog(scope.row)"
              >测试结果</el-button
            >
            <el-button
              v-if="canSubmitRelease(firmwareOf(scope.row))"
              type="primary"
              link
              @click="changeStage(scope.row, 'pending_release', '提交发布')"
              >提交发布</el-button
            >
            <el-button
              v-if="canRejectRelease(firmwareOf(scope.row))"
              type="primary"
              link
              @click="changeStage(scope.row, 'testing', '驳回到测试中')"
              >驳回</el-button
            >
            <el-button
              v-if="canPublish(firmwareOf(scope.row))"
              type="primary"
              link
              @click="publishRow(scope.row)"
              >发布</el-button
            >
            <el-button
              v-if="canToggleStable(firmwareOf(scope.row))"
              type="primary"
              link
              @click="toggleStable(scope.row)"
              >{{
                firmwareOf(scope.row).isStable ? '取消稳定' : '标记稳定'
              }}</el-button
            >
            <el-button
              v-if="canSetCurrentRelease(firmwareOf(scope.row))"
              type="primary"
              link
              @click="setCurrentRelease(scope.row)"
              >设为当前发布</el-button
            >
            <el-button
              v-if="canVoid(firmwareOf(scope.row))"
              type="primary"
              link
              @click="voidRow(scope.row)"
              >作废</el-button
            >
            <el-button type="primary" link @click="openLogDrawer(scope.row)"
              >日志</el-button
            >
            <el-button
              v-if="firmwareOf(scope.row).packageUrl"
              type="primary"
              link
              @click="downloadPackage(scope.row)"
              >下载</el-button
            >
            <el-button
              v-if="canDeleteRelation(firmwareOf(scope.row))"
              type="primary"
              link
              @click="deleteRelation(scope.row)"
              >移除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="firmwareDialogVisible"
      :title="firmwareDialogTitle"
      width="720px"
    >
      <el-form :model="firmwareForm" label-width="100px">
        <el-form-item label="设备类别">
          <el-select
            v-model="firmwareForm.categoryId"
            filterable
            clearable
            placeholder="请选择设备类别"
            style="width: 100%"
            :disabled="firmwareDialogType === 'update'"
          >
            <el-option
              v-for="item in categoryOptions"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="设备型号">
          <el-select
            v-model="firmwareForm.modelId"
            filterable
            clearable
            placeholder="请选择设备型号"
            style="width: 100%"
            :disabled="firmwareDialogType === 'update'"
          >
            <el-option
              v-for="item in filteredDialogModels"
              :key="item.ID"
              :label="item.modelName"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="版本号"
          ><el-input v-model="firmwareForm.versionCode"
        /></el-form-item>
        <el-form-item label="版本名称"
          ><el-input v-model="firmwareForm.versionName"
        /></el-form-item>
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
              <el-button
                type="primary"
                :loading="firmwareUploading"
                :disabled="!canEditFirmwarePackage(firmwareForm)"
                >上传固件包</el-button
              >
            </el-upload>
            <el-button
              v-if="firmwareForm.packageUrl"
              type="danger"
              plain
              :disabled="!firmwareForm.packageUrl"
              @click="deleteFirmwarePackage"
              >删除固件包</el-button
            >
            <span class="upload-file-name">
              已选文件：{{ firmwareUploadName || firmwareForm.packageName || '未选择' }}
            </span>
          </div>
        </el-form-item>
        <!-- <el-form-item label="安装包地址"
          ><el-input v-model="firmwareForm.packageUrl" readonly
        /></el-form-item> -->
        <el-form-item label="开发状态"
          ><el-input
            :model-value="devStatusLabel(firmwareForm.status || 'pending_test')"
            readonly
        /></el-form-item>
        <el-form-item v-if="firmwareDialogType === 'update'" label="发布状态"
          ><el-input
            :model-value="publishStatusLabel(firmwareForm.publishStatus)"
            readonly
        /></el-form-item>
        <el-form-item label="上传人"
          ><el-input v-model="firmwareForm.uploadedBy"
        /></el-form-item>
        <el-form-item label="版本说明"
          ><el-input
            v-model="firmwareForm.releaseNote"
            type="textarea"
            :rows="3"
        /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="firmwareDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitFirmware">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="testResultDialogVisible"
      title="提交测试结果"
      width="560px"
    >
      <el-form :model="testResultForm" label-width="90px">
        <el-form-item label="测试结果">
          <el-radio-group v-model="testResultForm.result">
            <el-radio value="tested_pass">通过</el-radio>
            <el-radio value="test_failed">不通过</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item
          v-if="testResultForm.result === 'test_failed'"
          label="原因分类"
        >
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
      <el-table :data="logRows">
        <el-table-column label="时间" width="180"
          ><template #default="scope">{{
            formatDate(scope.row.operateAt || scope.row.CreatedAt)
          }}</template></el-table-column
        >
        <el-table-column label="动作" min-width="180"
          ><template #default="scope">{{
            logActionLabel(scope.row)
          }}</template></el-table-column
        >
        <el-table-column label="目标状态" width="120"
          ><template #default="scope">{{
            devStatusLabel(scope.row.toStatus)
          }}</template></el-table-column
        >
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
    getFirmwareVersionLogList,
    getDeviceCategoryList,
    getDeviceModelList
  } from '@/api/deviceFirmware'
  import { deleteFirmwarePackage as deleteFirmwarePackageApi } from '@/api/deviceFirmware'
  import { deleteFile } from '@/api/fileUploadAndDownload'
  import { formatDate, getBaseUrl } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { computed, onMounted, ref, watch } from 'vue'
  import { useUserStore } from '@/pinia/modules/user'

  defineOptions({ name: 'DeviceFirmwareWorkflow' })

  const userStore = useUserStore()
  const categoryOptions = ref([])
  const modelOptions = ref([])
  const firmwareOptions = ref([])
  const relationRows = ref([])
  const firmwareLogs = ref([])
  const logRows = ref([])
  const firmwareDialogVisible = ref(false)
  const testResultDialogVisible = ref(false)
  const logDrawerVisible = ref(false)
  const firmwareDialogType = ref('create')
  const currentTestRow = ref(null)
  const firmwareUploading = ref(false)
  const firmwareUploadName = ref('')
  const logDrawerTitle = ref('固件日志')
  const failReasonOptions = ['有Bug', '少功能', '需优化']
  const devStatusOptions = [
    { label: '待测试', value: 'pending_test' },
    { label: '测试中', value: 'testing' },
    { label: '测试通过', value: 'tested_pass' },
    { label: '测试不通过', value: 'test_failed' },
    { label: '待发布', value: 'pending_release' }
  ]
  const publishStatusOptions = [
    { label: '未发布', value: 'unpublished' },
    { label: '已发布', value: 'published' },
    { label: '已作废', value: 'voided' }
  ]
  const searchForm = ref({
    categoryId: '',
    modelId: '',
    status: '',
    publishStatus: '',
    keyword: ''
  })
  const firmwareForm = ref({
    ID: undefined,
    categoryId: '',
    modelId: '',
    versionCode: '',
    versionName: '',
    packageUrl: '',
    packageName: '',
    packageFileId: undefined,
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
  const filteredSearchModels = computed(() =>
    !searchForm.value.categoryId
      ? modelOptions.value
      : modelOptions.value.filter(
          (item) => item.categoryId === searchForm.value.categoryId
        )
  )
  const filteredDialogModels = computed(() =>
    !firmwareForm.value.categoryId
      ? modelOptions.value
      : modelOptions.value.filter(
          (item) => item.categoryId === firmwareForm.value.categoryId
        )
  )
  const filteredRows = computed(() =>
    relationRows.value.filter((row) => {
      const firmware = firmwareOf(row)
      const categoryId = row.model?.categoryId || row.model?.category?.ID
      const keyword = (searchForm.value.keyword || '').trim().toLowerCase()
      if (
        searchForm.value.categoryId &&
        categoryId !== searchForm.value.categoryId
      )
        return false
      if (searchForm.value.modelId && row.modelId !== searchForm.value.modelId)
        return false
      if (
        searchForm.value.status &&
        firmware.status !== searchForm.value.status
      )
        return false
      if (
        searchForm.value.publishStatus &&
        firmware.publishStatus !== searchForm.value.publishStatus
      )
        return false
      if (!keyword) return true
      const haystack = [
        firmware.versionCode,
        firmware.versionName,
        row.model?.modelName,
        row.model?.category?.name
      ]
        .filter(Boolean)
        .join(' ')
        .toLowerCase()
      return haystack.includes(keyword)
    })
  )
  const currentDialogModel = computed(() =>
    modelOptions.value.find((item) => item.ID === firmwareForm.value.modelId)
  )
  const firmwareDialogTitle = computed(() => {
    if (firmwareDialogType.value === 'update')
      return firmwareForm.value.status === 'test_failed'
        ? '上传修复包'
        : '编辑固件'
    return currentDialogModel.value
      ? `为 ${currentDialogModel.value.modelName} 上传新固件`
      : '新增固件'
  })
  const latestFailureLogMap = computed(() => {
    const map = {}
    ;(firmwareLogs.value || []).forEach((log) => {
      if (log.action !== 'test_fail' || !log.firmwareId) return
      const current = map[log.firmwareId]
      const currentTime = current
        ? new Date(current.operateAt || current.CreatedAt || 0).getTime()
        : 0
      const nextTime = new Date(log.operateAt || log.CreatedAt || 0).getTime()
      if (!current || nextTime >= currentTime) map[log.firmwareId] = log
    })
    return map
  })
  const firmwareOf = (row) => {
    if (!row) return {}
    return (
      firmwareOptions.value.find((item) => item.ID === row.firmwareId) ||
      row.firmware ||
      {}
    )
  }
  const devStatusLabel = (status) =>
    ({
      pending_test: '待测试',
      pending: '待测试',
      draft: '待测试',
      testing: '测试中',
      tested_pass: '测试通过',
      passed: '测试通过',
      test_failed: '测试不通过',
      failed: '测试不通过',
      pending_release: '待发布'
    }[status] ||
    status ||
    '-')
  const publishStatusLabel = (status) =>
    ({
      unpublished: '未发布',
      published: '已发布',
      voided: '已作废'
    }[status] || '-')
  const devStatusTag = (status) => {
    if (status === 'tested_pass') return 'success'
    if (status === 'testing' || status === 'pending_release') return 'warning'
    if (status === 'test_failed') return 'danger'
    return 'info'
  }
  const publishStatusTag = (status) => {
    if (status === 'published') return 'success'
    if (status === 'voided') return 'danger'
    return 'info'
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
    firmware?.status === 'tested_pass' &&
    firmware?.publishStatus === 'unpublished'
  const canRejectRelease = (firmware) =>
    firmware?.status === 'pending_release' &&
    firmware?.publishStatus === 'unpublished'
  const canPublish = (firmware) =>
    firmware?.status === 'pending_release' &&
    firmware?.publishStatus === 'unpublished'
  const canToggleStable = (firmware) => firmware?.publishStatus === 'published'
  const canSetCurrentRelease = (firmware) =>
    firmware?.publishStatus === 'published'
  const canVoid = (firmware) => firmware?.publishStatus === 'published'
  const canDeleteRelation = (firmware) =>
    !['published', 'voided'].includes(firmware?.publishStatus)
  const canEditFirmwarePackage = (firmware) =>
    !firmware?.ID ||
    (firmware?.publishStatus === 'unpublished' &&
      ['pending_test', 'test_failed'].includes(firmware?.status))
  const rowNote = (row) => {
    const firmware = firmwareOf(row)
    if (firmware.status === 'test_failed') {
      return (
        latestFailureLogMap.value[firmware.ID]?.content ||
        firmware.releaseNote ||
        firmware.packageName ||
        '-'
      )
    }
    return firmware.releaseNote || firmware.packageName || '-'
  }
  const logActionLabel = (log) => {
    const row = relationRows.value.find(
      (item) =>
        item.firmwareId === log?.firmwareId &&
        (!log?.modelId || item.modelId === log.modelId)
    )
    const modelName = log?.model?.modelName || row?.model?.modelName
    return (
      {
        upload: '上传固件',
        bind_model: modelName ? `绑定到 ${modelName}` : '绑定型号',
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
        delete_package: '删除固件包',
        set_recommended: modelName
          ? `设为 ${modelName} 当前发布`
          : '设为当前发布'
      }[log?.action] ||
      log?.action ||
      '-'
    )
  }
  const resetSearch = async () => {
    searchForm.value = {
      categoryId: '',
      modelId: '',
      status: '',
      publishStatus: '',
      keyword: ''
    }
    await loadPage()
  }
  const loadCategories = async () => {
    const res = await getDeviceCategoryList({ page: 1, pageSize: 999 })
    if (res.code === 0) categoryOptions.value = res.data.list || []
  }
  const loadModels = async () => {
    const res = await getDeviceModelList({ page: 1, pageSize: 999 })
    if (res.code === 0) modelOptions.value = res.data.list || []
  }
  const loadFirmwareOptions = async () => {
    const res = await getFirmwareVersionList({ page: 1, pageSize: 999 })
    if (res.code === 0) firmwareOptions.value = res.data.list || []
  }
  const loadRelations = async () => {
    const res = await getModelFirmwareRelList({ page: 1, pageSize: 999 })
    if (res.code === 0) relationRows.value = res.data.list || []
  }
  const loadLogs = async () => {
    const res = await getFirmwareVersionLogList({ page: 1, pageSize: 999 })
    if (res.code === 0) firmwareLogs.value = res.data.list || []
  }
  const loadPage = async () => {
    await Promise.all([
      loadCategories(),
      loadModels(),
      loadFirmwareOptions(),
      loadRelations(),
      loadLogs()
    ])
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
    firmwareForm.value.packageFileId = file.ID || file.id || undefined
    if (!firmwareForm.value.uploadedBy)
      firmwareForm.value.uploadedBy = defaultUploadedBy()
    ElMessage.success('固件包上传成功，已自动回填地址')
  }
  const handleFirmwareUploadError = () => {
    firmwareUploading.value = false
    ElMessage.error('固件包上传失败')
  }
  const deleteFirmwarePackage = async () => {
    if (!firmwareForm.value.packageUrl) {
      ElMessage.warning('当前没有可删除的安装包')
      return
    }
    try {
      await ElMessageBox.confirm(
        '确定删除当前固件包吗？该操作会同时删除 MinIO 里的文件。',
        '提示',
        { type: 'warning' }
      )
    } catch (error) {
      return
    }
    if (firmwareDialogType.value === 'create' && !firmwareForm.value.packageFileId) {
      ElMessage.warning('未找到可删除的安装包记录')
      return
    }
    const res = firmwareDialogType.value === 'update' && firmwareForm.value.ID
      ? await deleteFirmwarePackageApi({
          id: firmwareForm.value.ID,
          operator: defaultUploadedBy(),
          content: '删除安装包'
        })
      : await deleteFile({
          ID: firmwareForm.value.packageFileId
        })
    if (res.code !== 0) {
      return
    }
    firmwareForm.value.packageUrl = ''
    firmwareForm.value.packageName = ''
    firmwareForm.value.packageFileId = undefined
    firmwareUploadName.value = ''
    ElMessage.success('固件包已删除')
  }
  const openFirmwareDialog = async (row) => {
    firmwareDialogType.value = row?.ID ? 'update' : 'create'
    firmwareUploadName.value = ''
    firmwareUploading.value = false
    if (row?.ID) {
      const firmware = firmwareOf(row)
      const res = await findFirmwareVersion({ ID: firmware.ID })
      if (res.code === 0) {
        firmwareForm.value = {
          categoryId: row.model?.categoryId || row.model?.category?.ID || '',
          modelId: row.modelId,
          publishStatus: 'unpublished',
          ...res.data
        }
        firmwareUploadName.value = res.data.packageName || ''
      }
    } else {
      firmwareForm.value = {
        ID: undefined,
        categoryId: searchForm.value.categoryId || '',
        modelId: searchForm.value.modelId || '',
        versionCode: '',
        versionName: '',
        packageUrl: '',
        packageName: '',
        packageFileId: undefined,
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
    if (!firmwareForm.value.categoryId || !firmwareForm.value.modelId) {
      ElMessage.warning('请先选择设备类别和设备型号')
      return
    }
    if (!firmwareForm.value.versionCode) {
      ElMessage.warning('请填写版本号')
      return
    }
    const isFixUpload =
      firmwareDialogType.value === 'update' &&
      firmwareForm.value.status === 'test_failed'
    const payload = { ...firmwareForm.value, uploadedAt: new Date() }
    if (firmwareDialogType.value === 'create') {
      delete payload.ID
    }
    const res =
      firmwareDialogType.value === 'create'
        ? await createFirmwareVersion(payload)
        : await updateFirmwareVersion(payload)
    if (res.code !== 0) return

    let firmwareId = firmwareForm.value.ID || res.data?.ID
    if (!firmwareId && firmwareDialogType.value === 'create') {
      const listRes = await getFirmwareVersionList({
        page: 1,
        pageSize: 1,
        versionCode: firmwareForm.value.versionCode
      })
      firmwareId = listRes.data?.list?.[0]?.ID
    }

    if (firmwareDialogType.value === 'create' && firmwareId) {
      const bindRes = await createModelFirmwareRel({
        modelId: firmwareForm.value.modelId,
        firmwareId,
        isSupported: true,
        isRecommended: false,
        testResult: '',
        tester: '',
        remark: ''
      })
      if (bindRes.code !== 0) return
    }

    ElMessage.success(
      firmwareDialogType.value === 'create'
        ? '新增固件成功'
        : isFixUpload
        ? '修复包已更新，已进入待测试'
        : '固件更新成功'
    )
    firmwareDialogVisible.value = false
    await loadPage()
  }
  const changeStage = async (row, status, content) => {
    const firmware = firmwareOf(row)
    const res = await changeFirmwareVersionStatus({
      id: firmware.ID,
      status,
      operator: defaultUploadedBy(),
      content
    })
    if (res.code === 0) {
      ElMessage.success('状态更新成功')
      await loadPage()
      return true
    }
    return false
  }
  const openTestResultDialog = (row) => {
    currentTestRow.value = row
    testResultForm.value = {
      result: 'tested_pass',
      reasonTypes: [],
      description: ''
    }
    testResultDialogVisible.value = true
  }
  const submitTestResult = async () => {
    if (!currentTestRow.value) return
    if (
      testResultForm.value.result === 'test_failed' &&
      !testResultForm.value.reasonTypes.length &&
      !testResultForm.value.description?.trim()
    ) {
      ElMessage.warning('测试不通过时，请至少选择一个原因或填写说明')
      return
    }
    const reasonText =
      testResultForm.value.result === 'test_failed' &&
      testResultForm.value.reasonTypes.length
        ? `原因分类：${testResultForm.value.reasonTypes.join('、')}`
        : ''
    const descriptionText = testResultForm.value.description?.trim() || ''
    const content =
      [reasonText, descriptionText].filter(Boolean).join('；') ||
      (testResultForm.value.result === 'tested_pass'
        ? '测试通过'
        : '测试不通过')
    const success = await changeStage(
      currentTestRow.value,
      testResultForm.value.result,
      content
    )
    if (success) {
      testResultDialogVisible.value = false
      currentTestRow.value = null
    }
  }
  const publishRow = async (row) => {
    const firmware = firmwareOf(row)
    const res = await publishFirmwareVersion({
      id: firmware.ID,
      operator: defaultUploadedBy(),
      content: '发布进入发布版本'
    })
    if (res.code === 0) {
      ElMessage.success('发布成功，已进入发布版本')
      await loadPage()
    }
  }
  const toggleStable = async (row) => {
    const firmware = firmwareOf(row)
    const stable = !firmware.isStable
    const res = await setFirmwareStable({
      id: firmware.ID,
      stable,
      operator: defaultUploadedBy(),
      content: stable ? '手动标记为稳定版本' : '取消稳定版本标记'
    })
    if (res.code === 0) {
      ElMessage.success(stable ? '已标记为稳定版本' : '已取消稳定版本')
      await loadPage()
    }
  }
  const setCurrentRelease = async (row) => {
    const res = await setModelFirmwareRecommended({
      id: row.ID,
      operator: defaultUploadedBy(),
      content: '设为当前发布版'
    })
    if (res.code === 0) {
      ElMessage.success('已设为当前发布版')
      await loadPage()
    }
  }
  const voidRow = (row) => {
    const firmware = firmwareOf(row)
    ElMessageBox.prompt('请输入作废原因', '作废版本', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '作废原因'
    })
      .then(async ({ value }) => {
        const res = await voidFirmwareVersion({
          id: firmware.ID,
          operator: defaultUploadedBy(),
          voidReason: value,
          content: value || '作废已发布版本'
        })
        if (res.code === 0) {
          ElMessage.success('版本已作废')
          await loadPage()
        }
      })
      .catch(() => {})
  }
  const openLogDrawer = async (row) => {
    const firmware = firmwareOf(row)
    const res = await getFirmwareVersionLogList({
      page: 1,
      pageSize: 100,
      firmwareId: firmware.ID
    })
    if (res.code === 0) {
      logRows.value = res.data.list || []
      logDrawerTitle.value = `固件日志 - ${
        firmware.versionCode || firmware.versionName || ''
      }`
      logDrawerVisible.value = true
    }
  }
  const deleteRelation = (row) => {
    const firmware = firmwareOf(row)
    ElMessageBox.confirm(
      `确定把 ${firmware.versionCode || '该固件'} 从当前型号中移除吗？`,
      '提示',
      {
        type: 'warning'
      }
    )
      .then(async () => {
        const res = await deleteModelFirmwareRel({ ID: row.ID })
        if (res.code === 0) {
          ElMessage.success('移除成功')
          await loadPage()
        }
      })
      .catch(() => {})
  }
  const downloadPackage = (row) => {
    const firmware = firmwareOf(row)
    if (!firmware.packageUrl) {
      ElMessage.warning('当前固件还没有安装包地址')
      return
    }
    window.open(firmware.packageUrl, '_blank')
  }

  watch(
    () => searchForm.value.categoryId,
    (value) => {
      if (
        value &&
        !filteredSearchModels.value.some(
          (item) => item.ID === searchForm.value.modelId
        )
      ) {
        searchForm.value.modelId = ''
      }
    }
  )
  watch(
    () => firmwareForm.value.categoryId,
    (value) => {
      if (
        value &&
        !filteredDialogModels.value.some(
          (item) => item.ID === firmwareForm.value.modelId
        )
      ) {
        firmwareForm.value.modelId = ''
      }
    }
  )

  onMounted(async () => {
    await loadPage()
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

  .upload-file-name {
    color: #909399;
    font-size: 13px;
  }

  .upload-tip {
    color: #909399;
    font-size: 13px;
  }
</style>
