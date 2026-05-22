<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="MO号">
          <el-input
            v-model="searchInfo.moNumber"
            placeholder="请输入"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="型号">
          <el-input
            v-model="searchInfo.model"
            placeholder="请输入"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="业务类型">
          <el-select
            v-model="searchInfo.instrumentCategory"
            placeholder="请选择"
            clearable
            size="small"
            style="width: 120px"
          >
            <el-option label="线上" value="online" />
            <el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" />
            <el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="getList"
            >查询</el-button
          >
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData" border v-loading="loading" size="small">
      <el-table-column prop="moNumber" label="MO号" min-width="140" />
      <el-table-column label="批次数" width="90">
        <template #default="scope">
          <el-tag size="small" type="info">{{
            scope.row.batchCount || 0
          }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="model" label="型号" width="100" />
      <el-table-column
        prop="firmwareVersion"
        label="固件版本"
        min-width="130"
      />
      <el-table-column
        prop="mainboardFirmwareVersion"
        label="主板固件版本"
        min-width="150"
      />
      <el-table-column prop="pnCode" label="PN码" min-width="150" />
      <el-table-column label="业务类型" width="100">
        <template #default="scope">{{
          catLabel(scope.row.instrumentCategory)
        }}</template>
      </el-table-column>
      <el-table-column label="订单状态" width="110">
        <template #default="scope">
          <el-tag :type="orderStatusTagType(scope.row.status)" size="small">
            {{ batchStatusLabel(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="submitterName" label="提交人" width="100" />
      <el-table-column label="设备数" width="80">
        <template #default="scope">
          <DeviceStatusCount :row="scope.row" type="all" :count="scope.row.deviceCount" allow-rework-actions @changed="getList" />
        </template>
      </el-table-column>
      <el-table-column label="合格数" width="90">
        <template #default="scope">
          <DeviceStatusCount :row="scope.row" type="pass" :count="scope.row.passCount" allow-rework-actions @changed="getList" />
        </template>
      </el-table-column>
      <el-table-column label="不合格数" width="100">
        <template #default="scope">
          <DeviceStatusCount :row="scope.row" type="fail" :count="scope.row.failCount" allow-rework-actions @changed="getList" />
        </template>
      </el-table-column>
      <el-table-column label="返工数" width="90">
        <template #default="scope">
          <DeviceStatusCount :row="scope.row" type="rework" :count="scope.row.reworkCount" allow-rework-actions @changed="getList" />
        </template>
      </el-table-column>
      <el-table-column label="待复检" width="90">
        <template #default="scope">
          <DeviceStatusCount :row="scope.row" type="recheck" :count="scope.row.recheckCount" allow-rework-actions @changed="getList" />
        </template>
      </el-table-column>
      <el-table-column label="异常数" width="90">
        <template #default="scope">
          <DeviceStatusCount :row="scope.row" type="abnormal" :count="scope.row.abnormalCount" allow-rework-actions @changed="getList" />
        </template>
      </el-table-column>
      <el-table-column label="合格率" width="100">
        <template #default="scope">
          {{ passRateLabel(scope.row.passCount, scope.row.deviceCount) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="scope">
          <el-button
            size="small"
            type="warning"
            link
            @click="openBatchScan(scope.row)"
          >
            分批
          </el-button>
          <el-button
            size="small"
            type="success"
            link
            @click="openDispatch(scope.row)"
          >
            派检
          </el-button>
          <el-button
            size="small"
            type="primary"
            link
            @click="viewDetail(scope.row)"
            >详情</el-button
          >
          <!-- <el-button
            size="small"
            type="primary"
            link
            @click="editOrder(scope.row)"
            >编辑</el-button
          > -->
          <el-button
            v-if="scope.row.status === 0"
            size="small"
            type="danger"
            link
            @click="onDelete(scope.row.ID)"
            >删除</el-button
          >
          <el-button
            v-else
            size="small"
            type="danger"
            link
            @click="onForceDelete(scope.row)"
            >强制删除</el-button
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
        background
        small
        @size-change="getList"
        @current-change="getList"
      />
    </div>

    <el-drawer
      v-model="drawerVisible"
      title="编辑生产订单"
      size="520px"
      destroy-on-close
    >
      <el-form
        :model="formData"
        ref="formRef"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item label="MO号" prop="moNumber"
          ><el-input v-model="formData.moNumber"
        /></el-form-item>
        <el-form-item label="内部型号"
          ><el-input v-model="formData.model"
        /></el-form-item>
        <el-form-item label="固件版本"
          ><el-input v-model="formData.firmwareVersion"
        /></el-form-item>
        <el-form-item label="主板固件版本"
          ><el-input v-model="formData.mainboardFirmwareVersion"
        /></el-form-item>
        <el-form-item label="PN码"
          ><el-input v-model="formData.pnCode"
        /></el-form-item>
        <el-form-item label="业务类型">
          <el-select v-model="formData.instrumentCategory" style="width: 100%">
            <el-option label="线上" value="online" />
            <el-option label="线下" value="offline" />
            <el-option label="外贸" value="foreign_trade" />
            <el-option label="定制款" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"
          ><el-input v-model="formData.remark" type="textarea"
        /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit">确定</el-button>
      </template>
    </el-drawer>

    <el-dialog
      v-model="batchScanVisible"
      title="生产分批"
      width="820px"
      destroy-on-close
      @opened="focusScanInput"
    >
      <div v-if="batchScanOrder">
        <el-descriptions :column="3" border size="small" class="mb-4">
          <el-descriptions-item label="MO号">{{
            batchScanOrder.moNumber || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="未分批">{{
            unbatchedForScan.length
          }}</el-descriptions-item>
          <el-descriptions-item label="已入篮">{{
            scanBasket.length
          }}</el-descriptions-item>
        </el-descriptions>

        <el-form label-width="90px">
          <el-form-item label="批次号">
            <el-input
              v-model="batchScanForm.batchNumber"
              readonly
              style="width: 320px"
            />
          </el-form-item>
          <el-form-item label="扫码SN">
            <el-input
              ref="scanInputRef"
              v-model="batchScanForm.scanSN"
              placeholder="扫描设备条码后回车加入批次"
              class="scan-input"
              @keyup.enter="addScannedSN"
            />
          </el-form-item>
        </el-form>

        <div class="scan-board">
          <div class="scan-basket">
            <div class="scan-title">批次</div>
            <el-empty v-if="!scanBasket.length" description="还没有加入设备" />
            <div v-else class="scan-list">
              <div
                v-for="(item, index) in scanBasket"
                :key="item.sn"
                class="scan-item"
              >
                <span>{{ index + 1 }}. {{ item.sn }}</span>
                <el-button
                  type="danger"
                  link
                  size="small"
                  @click="removeScanItem(item.sn)"
                  >移除</el-button
                >
              </div>
            </div>
          </div>
          <div class="scan-waiting">
            <div class="scan-title">未分批设备</div>
            <div class="scan-tip">只允许以下序列号加入</div>
            <div class="scan-list">
              <div
                v-for="device in unbatchedForScan.slice(0, 80)"
                :key="device.ID"
                class="scan-waiting-item"
              >
                {{ device.sn }}
              </div>
            </div>
            <div v-if="unbatchedForScan.length > 80" class="scan-tip">
              还有 {{ unbatchedForScan.length - 80 }} 台未显示，可直接扫码加入。
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="batchScanVisible = false">取消</el-button>
        <el-button
          type="primary"
          :disabled="!scanBasket.length"
          @click="submitBatchScan"
        >
          确认绑定批次
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="dispatchVisible"
      title="生产派检"
      width="760px"
      destroy-on-close
    >
      <div v-if="dispatchOrder">
        <el-descriptions :column="2" border size="small" class="mb-4">
          <el-descriptions-item label="MO号">{{
            dispatchOrder.moNumber || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="内部型号">{{
            dispatchOrder.model || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="PN码">{{
            dispatchOrder.pnCode || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="设备数">{{
            dispatchOrder.deviceCount || 0
          }}</el-descriptions-item>
        </el-descriptions>

        <el-form label-width="100px" class="dispatch-form">
          <el-form-item label="业务类型">
            <el-select
              v-model="dispatchForm.instrumentCategory"
              placeholder="请选择业务类型"
              style="width: 260px"
            >
              <el-option label="线上" value="online" />
              <el-option label="线下" value="offline" />
              <el-option label="外贸" value="foreign_trade" />
              <el-option label="定制款" value="custom" />
            </el-select>
          </el-form-item>
          <el-form-item label="检测模板">
            <el-select
              v-model="dispatchForm.templateID"
              placeholder="请选择检测模板"
              filterable
              style="width: 320px"
            >
              <el-option
                v-for="template in templateList"
                :key="template.ID"
                :label="template.name"
                :value="template.ID"
              />
            </el-select>
          </el-form-item>
        </el-form>

        <div class="dispatch-summary">
          <div>
            未派检批次：
            <span class="dispatch-strong">{{
              dispatchPendingBatches.length
            }}</span>
            个
          </div>
          <div>
            将统一使用模板：
            <span class="dispatch-strong">{{
              selectedDispatchTemplate?.name || '-'
            }}</span>
          </div>
        </div>

        <el-table
          :data="dispatchOrder.batches || []"
          border
          size="small"
          class="mt-3"
        >
          <el-table-column prop="batchNumber" label="批次号" min-width="160" />
          <el-table-column label="设备数" width="90">
            <template #default="scope">{{
              scope.row.devices?.length || 0
            }}</template>
          </el-table-column>
          <el-table-column label="当前模板" min-width="160">
            <template #default="scope">{{
              scope.row.template?.name || '-'
            }}</template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="scope">
              <el-tag
                size="small"
                :type="
                  scope.row.status === 0
                    ? 'info'
                    : scope.row.status === 1
                    ? 'warning'
                    : scope.row.status === 2
                    ? 'primary'
                    : scope.row.status === 3
                    ? 'warning'
                    : 'success'
                "
              >
                {{ batchStatusLabel(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>

        <el-alert
          v-if="!dispatchPendingBatches.length"
          class="mt-3"
          title="当前生产订单没有未派检批次。已派检、检测中或已完成的批次不会重复提交。"
          type="info"
          :closable="false"
        />
      </div>
      <template #footer>
        <el-button @click="dispatchVisible = false">取消</el-button>
        <el-button
          type="success"
          :disabled="!dispatchForm.templateID || !dispatchPendingBatches.length"
          @click="submitDispatch"
        >
          提交检测接收
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="detailVisible"
      title="生产订单详情"
      width="1100px"
      destroy-on-close
    >
      <div v-if="detailOrder">
        <el-descriptions :column="2" border size="small" class="mb-4">
          <el-descriptions-item label="MO号">{{
            detailOrder.moNumber || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="内部型号">{{
            detailOrder.model || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="固件版本">{{
            detailOrder.firmwareVersion || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="主板固件版本">{{
            detailOrder.mainboardFirmwareVersion || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="PN码">{{
            detailOrder.pnCode || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="业务类型">{{
            catLabel(detailOrder.instrumentCategory) || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="批次数">{{
            detailOrder.batchCount ?? detailOrder.batches?.length ?? 0
          }}</el-descriptions-item>
          <el-descriptions-item label="设备数">{{
            detailOrder.deviceCount ?? detailOrder.devices?.length ?? 0
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="mb-3">
          <el-alert
            title="当前页面用于查看设备、批次和日志；派检请回到列表点击“派检”，一个生产号只选一次检测模板。"
            type="info"
            :closable="false"
          />
        </div>

        <div v-if="unbatchedDevices.length" class="mb-4">
          <el-text type="info" size="small"
            >无批次 ({{ unbatchedDevices.length }}台)</el-text
          >
          <el-table :data="unbatchedDevices" border size="small" class="mt-1">
            <el-table-column prop="sn" label="SN" min-width="140" />
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag
                  :type="deviceStatusTagType(scope.row.status)"
                  size="small"
                >
                  {{ deviceStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="model" label="型号" width="90" />
            <el-table-column prop="pnCode" label="PN码" width="130" />
            <el-table-column prop="firmwareVersion" label="固件" width="110" />
            <el-table-column
              prop="mainboardFirmwareVersion"
              label="主板固件"
              width="130"
            />
            <el-table-column label="操作" width="170" fixed="right">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 'returned'"
                  type="warning"
                  link
                  size="small"
                  @click="handleConfirmReworkReceived(scope.row)"
                >
                  确认接收返工
                </el-button>
                <el-button
                  v-if="scope.row.status === 'rework'"
                  type="warning"
                  link
                  size="small"
                  @click="handleConfirmRework(scope.row)"
                >
                  确认返工完成
                </el-button>
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="openFlowLogs({ device: scope.row })"
                >
                  设备日志
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-for="batch in detailOrder.batches" :key="batch.ID" class="mb-4">
          <div class="flex items-center justify-between mb-2">
            <el-text type="primary" size="small"
              >批次: {{ batch.batchNumber || '-' }}</el-text
            >
            <div class="flex items-center gap-2">
              <div class="text-sm text-gray-500">
                设备数: {{ batch.devices?.length || 0 }}
              </div>
              <el-tag v-if="batch.template" size="small" type="info">
                模板: {{ batch.template.name }}
              </el-tag>
              <el-tag
                size="small"
                :type="
                  batch.status === 1
                    ? 'warning'
                    : batch.status === 2
                    ? 'primary'
                    : batch.status === 3
                    ? 'warning'
                    : 'success'
                "
              >
                {{ batchStatusLabel(batch.status) }}
              </el-tag>
              <el-button
                size="small"
                type="success"
                link
                @click="onExportBatchExcel(batch)"
              >
                导出Excel
              </el-button>
              <el-button
                size="small"
                type="primary"
                link
                @click="openBatchPrint(batch)"
              >
                打印
              </el-button>
              <el-button
                size="small"
                type="primary"
                link
                @click="openFlowLogs({ batch })"
              >
                流转日志
              </el-button>
            </div>
          </div>
          <el-table :data="batch.devices" border size="small" class="mt-1">
            <el-table-column prop="sn" label="SN" min-width="140" />
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag
                  :type="deviceStatusTagType(scope.row.status)"
                  size="small"
                >
                  {{ deviceStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="model" label="型号" width="90" />
            <el-table-column prop="pnCode" label="PN码" width="130" />
            <el-table-column prop="firmwareVersion" label="固件" width="110" />
            <el-table-column
              prop="mainboardFirmwareVersion"
              label="主板固件"
              width="130"
            />
            <el-table-column label="操作" width="170" fixed="right">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 'returned'"
                  type="warning"
                  link
                  size="small"
                  @click="handleConfirmReworkReceived(scope.row)"
                >
                  确认接收返工
                </el-button>
                <el-button
                  v-if="scope.row.status === 'rework'"
                  type="warning"
                  link
                  size="small"
                  @click="handleConfirmRework(scope.row)"
                >
                  确认返工完成
                </el-button>
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="openFlowLogs({ batch, device: scope.row })"
                >
                  设备日志
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-dialog>

    <FlowLogDrawer
      v-model="flowLogVisible"
      :title="flowLogDrawerTitle"
      :subject="flowLogTitle"
      :logs="flowLogs"
      :mode="flowLogMode"
    />

  </div>
</template>

<script setup>
  import { computed, reactive, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import {
    getProductionOrderList,
    deleteProductionOrder,
    forceDeleteProductionOrder,
    updateProductionOrder,
    findProductionOrder,
    confirmReworkReceived,
    confirmReworkDone,
    scanAssignBatch,
  } from '@/plugin/inspection/api/production_order'
  import { getTemplateList } from '@/plugin/inspection/api/template'
  import {
    assignOrderTemplate,
    exportInspectionExcel,
    getFlowLogs
  } from '@/plugin/inspection/api/work_order'
  import DeviceStatusCount from '@/plugin/inspection/components/DeviceStatusCount.vue'
  import FlowLogDrawer from '@/plugin/inspection/components/FlowLogDrawer.vue'

  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const drawerVisible = ref(false)
  const batchScanVisible = ref(false)
  const detailVisible = ref(false)
  const dispatchVisible = ref(false)
  const formRef = ref(null)
  const scanInputRef = ref(null)
  const detailOrder = ref(null)
  const batchScanOrder = ref(null)
  const dispatchOrder = ref(null)
  const templateList = ref([])
  const flowLogVisible = ref(false)
  const flowLogTitle = ref('')
  const flowLogDrawerTitle = ref('流转日志')
  const flowLogMode = ref('flow')
  const flowLogs = ref([])
  const dispatchForm = reactive({
    templateID: null,
    instrumentCategory: ''
  })
  const batchScanForm = reactive({
    batchNumber: '',
    scanSN: ''
  })
  const scanBasket = ref([])

  const searchInfo = reactive({
    moNumber: '',
    model: '',
    instrumentCategory: '',
    page: 1,
    pageSize: 30
  })

  const formData = reactive({
    ID: 0,
    moNumber: '',
    model: '',
    firmwareVersion: '',
    mainboardFirmwareVersion: '',
    pnCode: '',
    instrumentCategory: '',
    remark: ''
  })

  const rules = {
    moNumber: [{ required: true, message: '请输入MO号', trigger: 'blur' }]
  }

  const catLabel = (value) =>
    ({
      online: '线上',
      offline: '线下',
      foreign_trade: '外贸',
      custom: '定制款'
    }[value] || value)
  const batchStatusLabel = (value) =>
    ({ 0: '未派检', 1: '待检测接收', 2: '检测中', 3: '待确认', 4: '已完成' }[value] || value)
  const orderStatusTagType = (value) =>
    ({ 0: 'info', 1: 'warning', 2: 'primary', 3: 'warning', 4: 'success' }[value] || 'info')
  const deviceStatusLabel = (value) =>
    ({
      pending: '待检测设备',
      pass: '合格',
      fail: '不合格',
      returned: '待生产接收',
      rework: '返工中',
      pending_recheck: '待复检',
      rechecking: '复检中'
    }[value] ||
    value ||
    '-')
  const deviceStatusTagType = (value) =>
    ({
      pending: 'info',
      pass: 'success',
      fail: 'danger',
      returned: 'warning',
      rework: 'warning',
      pending_recheck: 'primary',
      rechecking: 'warning'
    }[value] || 'info')
  const passRateLabel = (passCount, deviceCount) => {
    const total = Number(deviceCount || 0)
    if (!total) return '-'
    return `${((Number(passCount || 0) / total) * 100).toFixed(1)}%`
  }

  const unbatchedDevices = computed(() => {
    if (!detailOrder.value) return []
    const allDevices = detailOrder.value.devices || []
    const batchedSet = new Set()
    detailOrder.value.batches?.forEach((batch) =>
      batch.devices?.forEach((device) => batchedSet.add(device.ID))
    )
    return allDevices.filter((device) => !batchedSet.has(device.ID))
  })

  const unbatchedForScan = computed(() => {
    if (!batchScanOrder.value) return []
    const basketSet = new Set(scanBasket.value.map((item) => item.sn))
    const allDevices = batchScanOrder.value.devices || []
    return allDevices.filter(
      (device) => !device.batchID && !basketSet.has(device.sn)
    )
  })

  const dispatchPendingBatches = computed(() =>
    (dispatchOrder.value?.batches || []).filter((batch) => batch.status === 0)
  )

  const selectedDispatchTemplate = computed(() =>
    templateList.value.find(
      (template) => template.ID === dispatchForm.templateID
    )
  )

  const ensureBatchPrintable = (batch) => {
    if (!batch?.template && !batch?.templateID) {
      ElMessage.warning('请先派检并绑定检测模板，再打印或导出')
      return false
    }
    return true
  }

  const openBatchPrint = (batch) => {
    if (!ensureBatchPrintable(batch)) return
    const url = `${window.location.origin}${window.location.pathname}#/inspectPrint?batchId=${batch.ID}`
    window.open(url, '_blank')
  }

  const downloadBlob = (blob, filename) => {
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  }

  const onExportBatchExcel = async (batch) => {
    if (!ensureBatchPrintable(batch)) return
    const res = await exportInspectionExcel({ id: batch.ID })
    const filename = `${detailOrder.value?.moNumber || 'MO'}-${batch.batchNumber || batch.ID}-检测工单.xlsx`
    downloadBlob(res.data || res, filename)
  }

  const getList = async () => {
    loading.value = true
    try {
      const res = await getProductionOrderList(searchInfo)
      if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
      }
    } finally {
      loading.value = false
    }
  }

  const resetSearch = () => {
    searchInfo.moNumber = ''
    searchInfo.model = ''
    searchInfo.instrumentCategory = ''
    searchInfo.page = 1
    getList()
  }

  const editOrder = async (row) => {
    formData.ID = row.ID
    formData.moNumber = row.moNumber || ''
    formData.model = row.model || ''
    formData.firmwareVersion = row.firmwareVersion || ''
    formData.mainboardFirmwareVersion = row.mainboardFirmwareVersion || ''
    formData.pnCode = row.pnCode || ''
    formData.instrumentCategory = row.instrumentCategory || ''
    formData.remark = row.remark || ''
    drawerVisible.value = true
  }

  const submitEdit = async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return
    const res = await updateProductionOrder({ ...formData })
    if (res.code !== 0) return
    ElMessage.success('更新成功')
    drawerVisible.value = false
    getList()
  }

  const loadTemplates = async () => {
    if (templateList.value.length > 0) return
    const res = await getTemplateList({ page: 1, pageSize: 100 })
    if (res.code === 0) {
      templateList.value = res.data.list
    }
  }

  const viewDetail = async (row) => {
    const res = await findProductionOrder({ id: row.ID })
    if (res.code !== 0) return
    detailOrder.value = res.data
    detailVisible.value = true
  }

  const openBatchScan = async (row) => {
    const res = await findProductionOrder({ id: row.ID })
    if (res.code !== 0) return
    batchScanOrder.value = res.data
    batchScanForm.batchNumber = previewNextBatchNumber(res.data)
    batchScanForm.scanSN = ''
    scanBasket.value = []
    batchScanVisible.value = true
  }

  const formatDateCompact = (date = new Date()) => {
    const y = date.getFullYear()
    const m = String(date.getMonth() + 1).padStart(2, '0')
    const d = String(date.getDate()).padStart(2, '0')
    return `${y}${m}${d}`
  }

  const previewNextBatchNumber = (order) => {
    const dateText = formatDateCompact()
    const prefix = `${order.moNumber || ''}-${dateText}-`
    const sameDay = (order.batches || []).filter((batch) =>
      String(batch.batchNumber || '').startsWith(prefix)
    )
    return `${prefix}${String(sameDay.length + 1).padStart(2, '0')}`
  }

  const focusScanInput = () => {
    setTimeout(() => {
      scanInputRef.value?.focus?.()
    }, 50)
  }

  const addScannedSN = () => {
    const sn = String(batchScanForm.scanSN || '').trim()
    batchScanForm.scanSN = ''
    if (!sn) {
      focusScanInput()
      return
    }
    if (scanBasket.value.some((item) => item.sn === sn)) {
      ElMessage.warning(`SN ${sn} 已在篮子里`)
      focusScanInput()
      return
    }
    const device = (batchScanOrder.value?.devices || []).find(
      (item) => item.sn === sn
    )
    if (!device) {
      ElMessage.error(`SN ${sn} 不属于当前生产订单`)
      focusScanInput()
      return
    }
    if (device.batchID) {
      const batch = (batchScanOrder.value?.batches || []).find(
        (item) => item.ID === device.batchID
      )
      ElMessage.error(
        `SN ${sn} 已在批次 ${batch?.batchNumber || device.batchID} 中`
      )
      focusScanInput()
      return
    }
    if (device.status !== 'pending') {
      ElMessage.error(`SN ${sn} 当前状态不是待检测设备，不能分批`)
      focusScanInput()
      return
    }
    scanBasket.value.unshift(device)
    ElMessage.success(`已加入 ${sn}`)
    focusScanInput()
  }

  const removeScanItem = (sn) => {
    scanBasket.value = scanBasket.value.filter((item) => item.sn !== sn)
    focusScanInput()
  }

  const submitBatchScan = async () => {
    if (!batchScanOrder.value) return
    if (!scanBasket.value.length) {
      ElMessage.warning('请先扫码加入设备')
      focusScanInput()
      return
    }
    const res = await scanAssignBatch({
      productionOrderID: batchScanOrder.value.ID,
      batchNumber: batchScanForm.batchNumber,
      sns: scanBasket.value.map((item) => item.sn)
    })
    if (res.code !== 0) return
    ElMessage.success('分批成功')
    const refresh = await findProductionOrder({ id: batchScanOrder.value.ID })
    if (refresh.code === 0) {
      batchScanOrder.value = refresh.data
      batchScanForm.batchNumber = previewNextBatchNumber(refresh.data)
    }
    scanBasket.value = []
    getList()
    focusScanInput()
  }

  const openDispatch = async (row) => {
    await loadTemplates()
    const res = await findProductionOrder({ id: row.ID })
    if (res.code !== 0) return
    dispatchOrder.value = res.data
    dispatchForm.templateID = res.data.templateID || null
    dispatchForm.instrumentCategory = res.data.instrumentCategory || ''
    dispatchVisible.value = true
  }

  const submitDispatch = async () => {
    if (!dispatchOrder.value) return
    if (!dispatchForm.templateID) {
      ElMessage.warning('请先选择检测模板')
      return
    }
    if (!dispatchPendingBatches.value.length) {
      ElMessage.warning('没有未派检批次')
      return
    }
    const res = await assignOrderTemplate({
      productionOrderID: dispatchOrder.value.ID,
      templateID: Number(dispatchForm.templateID),
      instrumentCategory: dispatchForm.instrumentCategory
    })
    if (res.code !== 0) return
    ElMessage.success('已提交检测接收')
    dispatchVisible.value = false
    dispatchOrder.value = null
    getList()
  }

  const handleConfirmReworkReceived = async (row) => {
    try {
      await ElMessageBox.confirm(
        `确认已接收 ${row.sn}，并开始返工？`,
        '确认接收返工',
        { type: 'warning', confirmButtonText: '确认接收' }
      )
    } catch {
      return
    }

    const res = await confirmReworkReceived({ deviceID: row.ID })
    if (res.code !== 0) return
    ElMessage.success('已进入返工中')
    if (detailOrder.value) {
      await viewDetail({ ID: detailOrder.value.ID })
    }
    getList()
  }

  const handleConfirmRework = async (row) => {
    try {
      await ElMessageBox.confirm(
        `确认 ${row.sn} 已返工完成，并提交给检测复检？`,
        '确认返工完成',
        { type: 'warning', confirmButtonText: '确认完成' }
      )
    } catch {
      return
    }

    const res = await confirmReworkDone({ deviceID: row.ID })
    if (res.code !== 0) return
    ElMessage.success('已进入待复检')
    if (detailOrder.value) {
      await viewDetail({ ID: detailOrder.value.ID })
    }
    getList()
  }

  const openFlowLogs = async ({ batch, device }) => {
    flowLogDrawerTitle.value = device ? '设备日志' : '流转日志'
    flowLogMode.value = device ? 'device' : 'flow'
    flowLogTitle.value = device?.sn || batch?.batchNumber || '-'
    flowLogs.value = []
    flowLogVisible.value = true
    const res = await getFlowLogs({
      batchID: batch?.ID || device?.batchID,
      deviceID: device?.ID
    })
    if (res.code === 0) {
      flowLogs.value = res.data || []
    }
  }

  const onDelete = (id) => {
    ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
      .then(async () => {
        const res = await deleteProductionOrder({ id })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          getList()
        }
      })
      .catch(() => {})
  }

  const onForceDelete = (row) => {
    ElMessageBox.confirm(
      `确定强制删除生产订单 ${row.moNumber} 吗？这会一并删除设备、批次和检测结果。`,
      '强制删除',
      { type: 'warning', confirmButtonText: '确认删除' }
    )
      .then(async () => {
        const res = await forceDeleteProductionOrder({ id: row.ID })
        if (res.code === 0) {
          ElMessage.success('强制删除成功')
          getList()
        }
      })
      .catch(() => {})
  }

  getList()
</script>

<style scoped>
  .dispatch-form {
    margin-top: 12px;
  }

  .dispatch-summary {
    display: flex;
    gap: 24px;
    padding: 12px 14px;
    border-radius: 8px;
    background: var(--el-fill-color-lighter, #fafafa);
    color: var(--el-text-color-regular, #606266);
  }

  .dispatch-strong {
    color: var(--el-color-primary, #409eff);
    font-weight: 700;
  }

  .scan-input {
    width: 420px;
  }

  .scan-board {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 14px;
  }

  .scan-basket,
  .scan-waiting {
    min-height: 280px;
    padding: 12px;
    border: 1px solid var(--el-border-color-light, #e4e7ed);
    border-radius: 10px;
    background: var(--el-fill-color-lighter, #fafafa);
  }

  .scan-title {
    font-weight: 700;
    margin-bottom: 8px;
  }

  .scan-tip {
    margin-bottom: 8px;
    font-size: 12px;
    color: var(--el-text-color-secondary, #909399);
  }

  .scan-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    max-height: 360px;
    overflow-y: auto;
  }

  .scan-item,
  .scan-waiting-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 7px 10px;
    border-radius: 6px;
    background: var(--el-bg-color, #fff);
    border: 1px solid var(--el-border-color-lighter, #ebeef5);
  }

</style>
