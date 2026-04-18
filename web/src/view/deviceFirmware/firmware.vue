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
          <div class="section-title">版本流程管理</div>
          <!-- <div class="section-subtitle">
            按测试版和正式版分层管理，正式版对应已发布版本，测试版保留所有未发布流程。
          </div> -->
        </div>
        <div class="gva-btn-list">
          <el-button plain @click="openPublicFirmwareDownloadPage"
            >公开下载页</el-button
          >
          <el-button type="primary" icon="plus" @click="openFirmwareDialog()"
            >新增固件</el-button
          >
        </div>
      </div>

      <div class="firmware-tabs">
        <el-tabs v-model="activeTab">
          <el-tab-pane
            :label="`测试版 (${testRows.length})`"
            name="test"
          ></el-tab-pane>
          <el-tab-pane
            :label="`正式版 (${officialRows.length})`"
            name="official"
          ></el-tab-pane>
          <el-tab-pane
            :label="`已下架 (${voidRows.length})`"
            name="voided"
          ></el-tab-pane>
        </el-tabs>
      </div>

      <el-table :data="displayRows" row-key="ID">
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
              <el-tag :type="firmwareStatusTag(firmwareOf(scope.row))">{{
                firmwareStatusLabel(firmwareOf(scope.row))
              }}</el-tag>
              <el-tag
                v-if="
                  firmwareOf(scope.row).publishStatus === 'published' &&
                  firmwareOf(scope.row).isLatest
                "
                type="danger"
                >最新版本</el-tag
              >
              <el-tag v-if="isHistoryVersion(firmwareOf(scope.row))" type="info"
                >历史版本</el-tag
              >
              <el-tag v-if="scope.row.isRecommended" type="warning"
                >当前推荐</el-tag
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
              v-if="firmwareOf(scope.row).publishStatus !== 'voided'"
              type="primary"
              link
              @click="openFirmwareDialog(scope.row)"
              >编辑信息</el-button
            >
            <el-button
              v-if="canEditFirmwarePackage(firmwareOf(scope.row))"
              type="primary"
              link
              @click="openPackageUpdateDialog(scope.row)"
              >更新包</el-button
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
              v-if="canDirectPublish(firmwareOf(scope.row))"
              type="primary"
              link
              @click="publishRow(scope.row, true)"
              >直接发布</el-button
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
              v-if="canSetCurrentRelease(firmwareOf(scope.row))"
              type="primary"
              link
              @click="setCurrentRelease(scope.row)"
              >设为当前推荐</el-button
            >
            <el-button
              v-if="canVoid(firmwareOf(scope.row))"
              type="primary"
              link
              @click="voidRow(scope.row)"
              >下架</el-button
            >
            <el-button
              v-if="canOnShelf(firmwareOf(scope.row))"
              type="primary"
              link
              @click="onShelfRow(scope.row)"
              >上架</el-button
            >
            <el-button type="primary" link @click="openLogDrawer(scope.row)"
              >日志</el-button
            >
            <el-button type="primary" link @click="openPackageDialog(scope.row)"
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
      <el-form
        ref="firmwareFormRef"
        :model="firmwareForm"
        :rules="firmwareRules"
        label-width="100px"
      >
        <el-form-item label="设备类别" prop="categoryId">
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
        <el-form-item label="设备型号" prop="modelId">
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
        <el-form-item label="版本号" prop="versionCode"
          ><el-input
            v-model="firmwareForm.versionCode"
            :disabled="
              firmwareDialogType === 'update' &&
              firmwareForm.publishStatus === 'published'
            "
        /></el-form-item>
        <el-form-item label="版本名称" prop="versionName"
          ><el-input
            v-model="firmwareForm.versionName"
            :disabled="
              firmwareDialogType === 'update' &&
              firmwareForm.publishStatus === 'published'
            "
        /></el-form-item>
        <el-form-item label="固件包上传" prop="packageUrl">
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
              v-if="firmwareDialogType === 'create' && firmwareForm.packageUrl"
              type="danger"
              plain
              :disabled="!firmwareForm.packageUrl"
              @click="deleteFirmwarePackage"
              >删除固件包</el-button
            >
            <span class="upload-file-name">
              {{ firmwareUploadName || firmwareForm.packageName || '未选择' }}
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
        <!-- <el-form-item label="上传人"
          ><el-input
            v-model="firmwareForm.uploadedBy"
            :disabled="
              firmwareDialogType === 'update' &&
              firmwareForm.publishStatus === 'published'
            "
        /></el-form-item> -->
        <el-form-item label="版本说明" prop="releaseNote"
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
      v-model="packageUpdateDialogVisible"
      :title="packageUpdateDialogTitle"
      width="520px"
    >
      <el-form
        ref="packageUpdateFormRef"
        class="package-update-form"
        :model="packageUpdateForm"
        :rules="packageUpdateRules"
        label-width="90px"
      >
        <el-form-item label="上传包" prop="packageUrl">
          <div class="upload-row">
            <el-upload
              :action="firmwareUploadAction"
              :headers="firmwareUploadHeaders"
              :show-file-list="false"
              :disabled="packageUpdateUploading"
              :before-upload="beforePackageUpdateUpload"
              :on-success="handlePackageUpdateUploadSuccess"
              :on-error="handlePackageUpdateUploadError"
            >
              <el-button
                type="primary"
                :loading="packageUpdateUploading"
                :disabled="packageUpdateUploading"
                >上传安装包</el-button
              >
            </el-upload>
            <span class="upload-file-name">
              {{
                packageUpdateUploadName ||
                packageUpdateForm.packageName ||
                '未选择'
              }}
            </span>
          </div>
        </el-form-item>
        <el-form-item label="版本说明" prop="releaseNote">
          <el-input
            v-model="packageUpdateForm.releaseNote"
            type="textarea"
            :rows="4"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closePackageUpdateDialog">取消</el-button>
        <el-button
          type="primary"
          :disabled="!packageUpdateForm.packageUrl"
          @click="submitPackageUpdate"
          >确定</el-button
        >
      </template>
    </el-dialog>

    <el-dialog
      v-model="testResultDialogVisible"
      title="提交测试结果"
      width="720px"
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
        <el-form-item label="通知邮箱">
          <el-input
            v-model="testResultForm.notifyTo"
            placeholder="留空则使用系统默认收件人，支持多个邮箱用英文逗号分隔"
          />
        </el-form-item>
        <el-form-item label="邮件附加">
          <el-input
            v-model="testResultForm.emailContent"
            type="textarea"
            :rows="4"
            placeholder="会追加到固定邮件模板后面"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testResultDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTestResult">确定</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="logDrawerVisible" :title="logDrawerTitle" size="720px">
      <el-empty v-if="!logTimelineRows.length" description="暂无日志" />
      <el-timeline v-else class="firmware-timeline">
        <el-timeline-item
          v-for="item in logTimelineRows"
          :key="item.ID"
          :timestamp="formatDate(item.operateAt || item.CreatedAt)"
          placement="top"
        >
          <div class="timeline-card">
            <div class="timeline-card-head">
              <div class="timeline-card-title">{{ logActionLabel(item) }}</div>
              <el-tag size="small">{{ logStatusLabel(item) }}</el-tag>
            </div>
            <div class="timeline-card-meta">
              <span>操作人：{{ item.operator || '-' }}</span>
              <span>当前状态：{{ logStatusLabel(item) }}</span>
            </div>
            <div class="timeline-card-content">{{ item.content || '-' }}</div>
            <div v-if="timelinePackageName(item)" class="timeline-card-footer">
              <span>包名：{{ timelinePackageName(item) }}</span>
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-drawer>

    <el-dialog
      v-model="packageDialogVisible"
      :title="packageDialogTitle"
      width="820px"
      @closed="closePackageDialog"
    >
      <el-empty v-if="!packageRows.length" description="暂无可下载安装包" />
      <el-table v-else :data="packageRows" row-key="ID">
        <el-table-column label="上传日期" width="180">
          <template #default="scope">{{
            formatDate(scope.row.operateAt || scope.row.CreatedAt)
          }}</template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="scope">
            <el-tag size="small">{{ logStatusLabel(scope.row) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="包名" min-width="220" show-overflow-tooltip>
          <template #default="scope">{{
            timelinePackageName(scope.row) || '-'
          }}</template>
        </el-table-column>
        <el-table-column v-if="showPackageSizeColumn" label="大小" width="120">
          <template #default="scope">{{
            formatPackageSize(timelinePackageSize(scope.row))
          }}</template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              :disabled="!timelinePackageUrl(scope.row)"
              @click="downloadPackageByRow(scope.row)"
              >下载</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="packageDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
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
    voidFirmwareVersion,
    onShelfFirmwareVersion,
    createModelFirmwareRel,
    deleteModelFirmwareRel,
    getModelFirmwareRelList,
    setModelFirmwareTestResult,
    setModelFirmwareRecommended,
    getFirmwareVersionLogList,
    getDeviceCategoryList,
    getDeviceModelList
  } from '@/api/deviceFirmware'
  import { deleteFile } from '@/api/fileUploadAndDownload'
  import { formatDate, getBaseUrl } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { computed, onMounted, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { useUserStore } from '@/pinia/modules/user'

  defineOptions({ name: 'DeviceFirmwareWorkflow' })

  const router = useRouter()
  const userStore = useUserStore()
  const categoryOptions = ref([])
  const modelOptions = ref([])
  const firmwareOptions = ref([])
  const relationRows = ref([])
  const firmwareLogs = ref([])
  const logRows = ref([])
  const packageRows = ref([])
  const firmwareDialogVisible = ref(false)
  const testResultDialogVisible = ref(false)
  const logDrawerVisible = ref(false)
  const packageDialogVisible = ref(false)
  const packageUpdateDialogVisible = ref(false)
  const firmwareDialogType = ref('create')
  const firmwareFormRef = ref()
  const packageUpdateFormRef = ref()
  const currentTestRow = ref(null)
  const selectedPackageFirmware = ref(null)
  const firmwareUploading = ref(false)
  const firmwareUploadName = ref('')
  const packageUpdateUploading = ref(false)
  const packageUpdateUploadName = ref('')
  const logDrawerTitle = ref('固件日志')
  const activeTab = ref('test')
  const failReasonOptions = ['有Bug', '少功能', '需优化']
  const devStatusOptions = [
    { label: '待测试', value: 'pending_test' },
    { label: '测试中', value: 'testing' },
    { label: '测试通过', value: 'tested_pass' },
    { label: '测试不通过', value: 'test_failed' }
  ]
  const publishStatusOptions = [
    { label: '未发布', value: 'unpublished' },
    { label: '已发布', value: 'published' },
    { label: '已下架', value: 'voided' }
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
  const packageUpdateForm = ref({
    ID: undefined,
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
    description: '',
    notifyTo: '',
    emailContent: ''
  })
  const firmwareUploadAction = `${getBaseUrl()}/fileUploadAndDownload/upload`
  const firmwareUploadHeaders = computed(() => ({
    'x-token': userStore.token,
    'x-user-id': userStore.userInfo?.ID || ''
  }))
  const defaultUploadedBy = () =>
    userStore.userInfo?.nickName || userStore.userInfo?.userName || '系统用户'
  const requiredRule = (message) => ({
    required: true,
    message,
    trigger: ['blur', 'change']
  })
  const firmwareRules = {
    categoryId: [requiredRule('请选择设备类别')],
    modelId: [requiredRule('请选择设备型号')],
    versionCode: [requiredRule('请输入版本号')],
    versionName: [requiredRule('请输入版本名称')],
    packageUrl: [requiredRule('请先上传安装包')],
    releaseNote: [requiredRule('请填写版本说明')]
  }
  const packageUpdateRules = {
    packageUrl: [requiredRule('请先上传安装包')],
    releaseNote: [requiredRule('请填写版本说明')]
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
  const testRows = computed(() =>
    filteredRows.value.filter(
      (row) => !['published', 'voided'].includes(firmwareOf(row)?.publishStatus)
    )
  )
  const officialRows = computed(() =>
    filteredRows.value.filter(
      (row) => firmwareOf(row)?.publishStatus === 'published'
    )
  )
  const voidRows = computed(() =>
    filteredRows.value.filter(
      (row) => firmwareOf(row)?.publishStatus === 'voided'
    )
  )
  const displayRows = computed(() =>
    activeTab.value === 'official'
      ? officialRows.value
      : activeTab.value === 'voided'
      ? voidRows.value
      : testRows.value
  )
  const currentDialogModel = computed(() =>
    modelOptions.value.find((item) => item.ID === firmwareForm.value.modelId)
  )
  const firmwareDialogTitle = computed(() => {
    if (firmwareDialogType.value === 'update') return '编辑信息'
    return currentDialogModel.value
      ? `为 ${currentDialogModel.value.modelName} 上传新固件`
      : '新增固件'
  })
  const packageUpdateDialogTitle = computed(() => {
    const firmware = packageUpdateForm.value
    return firmware.versionCode || firmware.versionName
      ? `更新包 - ${firmware.versionCode || firmware.versionName}`
      : '更新包'
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
  const logTimelineRows = computed(() =>
    [...logRows.value].sort((a, b) => {
      const bTime = new Date(b.operateAt || b.CreatedAt || 0).getTime()
      const aTime = new Date(a.operateAt || a.CreatedAt || 0).getTime()
      return bTime - aTime
    })
  )
  const packageDialogTitle = computed(() => {
    const firmware = selectedPackageFirmware.value
    const name = firmware?.versionCode || firmware?.versionName || ''
    return name ? `安装包列表 - ${name}` : '安装包列表'
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
      voided: '已下架'
    }[status] || '-')
  const firmwareStatusLabel = (firmware) => {
    if (firmware?.publishStatus === 'published') return '已发布'
    if (firmware?.publishStatus === 'voided') return '已下架'
    return devStatusLabel(firmware?.status)
  }
  const firmwareStatusTag = (firmware) => {
    if (firmware?.publishStatus === 'published') return 'success'
    if (firmware?.publishStatus === 'voided') return 'danger'
    return devStatusTag(firmware?.status)
  }
  const logStatusLabel = (log) => {
    if (log?.toStatus === 'published') return '已发布'
    if (log?.toStatus === 'voided') return '已下架'
    return devStatusLabel(log?.toStatus || log?.fromStatus)
  }
  const devStatusTag = (status) => {
    if (status === 'tested_pass') return 'success'
    if (status === 'testing' || status === 'pending_release') return 'warning'
    if (status === 'test_failed') return 'danger'
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
  const canDirectPublish = (firmware) =>
    firmware?.publishStatus === 'unpublished' &&
    !['published', 'voided'].includes(firmware?.publishStatus)
  const canRejectRelease = (firmware) =>
    firmware?.status === 'tested_pass' &&
    firmware?.publishStatus === 'unpublished'
  const canPublish = (firmware) =>
    firmware?.status === 'tested_pass' &&
    firmware?.publishStatus === 'unpublished'
  const canSetCurrentRelease = (firmware) =>
    firmware?.publishStatus === 'published'
  const canVoid = (firmware) => firmware?.publishStatus === 'published'
  const canOnShelf = (firmware) => firmware?.publishStatus === 'voided'
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
        mark_stable: '版本状态变更',
        unmark_stable: '版本状态变更',
        void_release: '下架发布版本',
        on_shelf_release: '上架发布版本',
        delete_package: '删除固件包',
        set_recommended: modelName
          ? `设为 ${modelName} 推荐版本`
          : '设为当前推荐'
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
  const beforePackageUpdateUpload = (file) => {
    if (!file?.name) {
      ElMessage.error('未读取到文件名，请重新选择文件')
      return false
    }
    packageUpdateUploading.value = true
    packageUpdateUploadName.value = file.name
    return true
  }
  const handlePackageUpdateUploadSuccess = (res) => {
    packageUpdateUploading.value = false
    const file = res?.data?.file
    if (!file?.url) {
      ElMessage.error('上传成功，但未返回文件地址')
      return
    }
    packageUpdateForm.value.packageUrl = file.url
    packageUpdateForm.value.packageName =
      file.name || packageUpdateUploadName.value
    packageUpdateForm.value.packageFileId = file.ID || file.id || undefined
    if (!packageUpdateForm.value.uploadedBy) {
      packageUpdateForm.value.uploadedBy = defaultUploadedBy()
    }
    ElMessage.success('安装包上传成功，已自动回填地址')
  }
  const handlePackageUpdateUploadError = () => {
    packageUpdateUploading.value = false
    ElMessage.error('安装包上传失败')
  }
  const deleteFirmwarePackage = async () => {
    if (!firmwareForm.value.packageUrl) {
      ElMessage.warning('当前没有可删除的安装包')
      return
    }
    if (firmwareDialogType.value !== 'create') {
      ElMessage.warning('已创建版本不能删除安装包，只能重新上传替换')
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
    if (
      firmwareDialogType.value === 'create' &&
      !firmwareForm.value.packageFileId
    ) {
      ElMessage.warning('未找到可删除的安装包记录')
      return
    }
    const res = await deleteFile({
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
  const openPackageUpdateDialog = async (row) => {
    const firmware = firmwareOf(row)
    const res = await findFirmwareVersion({ ID: firmware.ID })
    if (res.code === 0) {
      packageUpdateForm.value = {
        ID: res.data.ID,
        versionCode: res.data.versionCode || '',
        versionName: res.data.versionName || '',
        packageUrl: '',
        packageName: '',
        packageFileId: undefined,
        checksum: '',
        status: res.data.status || 'pending_test',
        publishStatus: res.data.publishStatus || 'unpublished',
        releaseNote: res.data.releaseNote || ''
      }
      packageUpdateUploadName.value = ''
      packageUpdateDialogVisible.value = true
    }
  }
  const submitPackageUpdate = async () => {
    if (!packageUpdateForm.value.packageUrl) {
      ElMessage.warning('请先上传安装包')
      return
    }
    if (!(await validateForm(packageUpdateFormRef))) {
      return
    }
    const payload = {
      ...packageUpdateForm.value,
      uploadedAt: new Date()
    }
    const res = await updateFirmwareVersion(payload)
    if (res.code !== 0) {
      return
    }
    ElMessage.success('安装包更新成功')
    closePackageUpdateDialog()
    await loadDeviceTree()
  }
  const closePackageUpdateDialog = () => {
    packageUpdateDialogVisible.value = false
    packageUpdateUploading.value = false
    packageUpdateUploadName.value = ''
    packageUpdateForm.value = {
      ID: undefined,
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
    }
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
    if (!(await validateForm(firmwareFormRef))) {
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
      description: '',
      notifyTo: '',
      emailContent: ''
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
    const res = await setModelFirmwareTestResult({
      id: currentTestRow.value.ID,
      testResult: testResultForm.value.result,
      tester: defaultUploadedBy(),
      operator: defaultUploadedBy(),
      content,
      notifyTo: testResultForm.value.notifyTo?.trim(),
      emailContent: testResultForm.value.emailContent?.trim()
    })
    if (res.code === 0) {
      await loadPage()
      testResultDialogVisible.value = false
      currentTestRow.value = null
    }
  }
  const publishRow = async (row, direct = false) => {
    const firmware = firmwareOf(row)
    try {
      await ElMessageBox.confirm(
        direct
          ? `确定直接发布 ${
              firmware.versionCode || firmware.versionName || '该版本'
            } 吗？`
          : `确定发布 ${
              firmware.versionCode || firmware.versionName || '该版本'
            } 吗？`,
        direct ? '直接发布' : '发布',
        {
          type: 'warning',
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }
      )
    } catch (error) {
      return
    }
    const res = await publishFirmwareVersion({
      id: firmware.ID,
      direct,
      operator: defaultUploadedBy(),
      content: direct ? '上传后直接发布' : '发布进入发布版本'
    })
    if (res.code === 0) {
      ElMessage.success(direct ? '已直接发布' : '发布成功，已进入发布版本')
      await loadPage()
    }
  }
  const setCurrentRelease = async (row) => {
    const res = await setModelFirmwareRecommended({
      id: row.ID,
      operator: defaultUploadedBy(),
      content: '设为当前推荐'
    })
    if (res.code === 0) {
      ElMessage.success('已设为当前推荐')
      await loadPage()
    }
  }
  const voidRow = (row) => {
    const firmware = firmwareOf(row)
    ElMessageBox.prompt('请输入下架原因', '下架版本', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '下架原因'
    })
      .then(async ({ value }) => {
        const res = await voidFirmwareVersion({
          id: firmware.ID,
          operator: defaultUploadedBy(),
          voidReason: value,
          content: value || '下架已发布版本'
        })
        if (res.code === 0) {
          ElMessage.success('版本已下架')
          await loadPage()
        }
      })
      .catch(() => {})
  }
  const onShelfRow = async (row) => {
    const firmware = firmwareOf(row)
    try {
      await ElMessageBox.confirm(
        `确定上架 ${
          firmware.versionCode || firmware.versionName || '该版本'
        } 吗？`,
        '上架版本',
        {
          type: 'warning',
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }
      )
    } catch (error) {
      return
    }
    const res = await onShelfFirmwareVersion({
      id: firmware.ID,
      operator: defaultUploadedBy(),
      content: '上架已下架版本'
    })
    if (res.code === 0) {
      ElMessage.success('版本已上架')
      await loadPage()
    }
  }
  const openLogDrawer = async (row) => {
    const firmware = firmwareOf(row)
    const res = await getFirmwareVersionLogList({
      page: 1,
      pageSize: 999,
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
  const openPublicFirmwareDownloadPage = () => {
    const { href } = router.resolve({ name: 'PublicFirmwareDownload' })
    window.open(href, '_blank', 'noopener,noreferrer')
  }
  const downloadPackage = (row) => {
    const firmware = firmwareOf(row)
    if (!firmware.packageUrl) {
      ElMessage.warning('当前固件还没有安装包地址')
      return
    }
    window.open(firmware.packageUrl, '_blank')
  }
  const packageHistoryActions = new Set(['upload', 'fix_upload'])
  const timelineCanViewPackage = (log) =>
    packageHistoryActions.has(log?.action) && !!log?.packageUrl
  const timelinePackageUrl = (log) => log?.packageUrl || ''
  const timelinePackageName = (log) => log?.packageName || ''
  const timelinePackageSize = (log) => Number(log?.packageSize || 0)
  const showPackageSizeColumn = computed(() =>
    packageRows.value.some((item) => timelinePackageSize(item) > 0)
  )
  const formatPackageSize = (size) => {
    if (!size) return '-'
    if (size < 1024) return `${size} B`
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
    if (size < 1024 * 1024 * 1024)
      return `${(size / (1024 * 1024)).toFixed(1)} MB`
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`
  }
  const openPackageDialog = async (log) => {
    const firmware = firmwareOf(log)
    if (!firmware?.ID) {
      ElMessage.warning('未找到可下载的固件记录')
      return
    }
    selectedPackageFirmware.value = firmware
    packageRows.value = []
    await loadPackageRows(firmware.ID)
    packageDialogVisible.value = true
  }
  const loadPackageRows = async (firmwareId) => {
    const res = await getFirmwareVersionLogList({
      page: 1,
      pageSize: 999,
      firmwareId
    })
    if (res.code !== 0) {
      packageRows.value = []
      return
    }
    packageRows.value = (res.data.list || [])
      .filter((item) => timelineCanViewPackage(item))
      .sort((a, b) => {
        const bTime = new Date(b.operateAt || b.CreatedAt || 0).getTime()
        const aTime = new Date(a.operateAt || a.CreatedAt || 0).getTime()
        return bTime - aTime
      })
  }
  const downloadPackageByRow = (row) => {
    const url = timelinePackageUrl(row)
    if (!url) {
      ElMessage.warning('当前安装包没有可下载地址')
      return
    }
    window.open(url, '_blank')
  }
  const closePackageDialog = () => {
    packageDialogVisible.value = false
    selectedPackageFirmware.value = null
    packageRows.value = []
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

  .firmware-tabs {
    margin-bottom: 12px;
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

  .firmware-timeline {
    padding: 8px 8px 0 6px;
  }

  .timeline-card {
    padding: 12px 14px;
    border: 1px solid #ebeef5;
    border-radius: 10px;
    background: #fff;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  }

  .timeline-card-head,
  .timeline-card-meta,
  .timeline-card-footer {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }

  .timeline-card-head {
    justify-content: space-between;
    margin-bottom: 8px;
  }

  .timeline-card-title {
    color: #303133;
    font-size: 14px;
    font-weight: 600;
  }

  .timeline-card-meta {
    color: #909399;
    font-size: 12px;
    margin-bottom: 8px;
  }

  .timeline-card-content {
    color: #606266;
    font-size: 13px;
    line-height: 1.6;
    word-break: break-word;
  }

  .timeline-card-footer {
    justify-content: space-between;
    margin-top: 10px;
    color: #909399;
    font-size: 12px;
  }
</style>
