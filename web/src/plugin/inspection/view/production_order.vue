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
        <el-form-item label="批次号">
          <el-input
            v-model="searchInfo.batchNumber"
            placeholder="请输入"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="SN">
          <el-input
            v-model="searchInfo.sn"
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
        <el-form-item label="完成状态">
          <el-select
            v-model="searchInfo.status"
            placeholder="请选择"
            clearable
            size="small"
            style="width: 130px"
          >
            <el-option
              v-for="item in orderStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="分批状态">
          <el-select
            v-model="searchInfo.batchComplete"
            placeholder="请选择"
            clearable
            size="small"
            style="width: 130px"
          >
            <el-option label="分批完成" :value="1" />
            <el-option label="未分批完成" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="检测异常">
          <el-select
            v-model="searchInfo.hasAbnormal"
            placeholder="请选择"
            clearable
            size="small"
            style="width: 130px"
          >
            <el-option label="有异常" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="提交时间">
          <el-date-picker
            v-model="submitDateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            size="small"
            style="width: 240px"
            clearable
          />
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
      <el-table-column label="未分批" width="90">
        <template #default="scope">
          <el-tag
            size="small"
            :type="scope.row.unbatchedCount > 0 ? 'warning' : 'info'"
          >
            {{ scope.row.unbatchedCount || 0 }}
          </el-tag>
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
      <el-table-column label="完成状态" width="110">
        <template #default="scope">
          <el-tag
            :type="productionStatusTagType(scope.row.status)"
            size="small"
          >
            {{ productionStatusLabel(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="submitterName" label="提交人" width="100" />
      <el-table-column label="提交时间" width="170">
        <template #default="scope">
          {{ formatDate(scope.row.submitDate) || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="设备数" width="80">
        <template #default="scope">
          <DeviceStatusCount
            :row="scope.row"
            type="all"
            :count="scope.row.deviceCount"
            allow-rework-actions
            @changed="getList"
          />
        </template>
      </el-table-column>
      <el-table-column label="合格数" width="90">
        <template #default="scope">
          <DeviceStatusCount
            :row="scope.row"
            type="pass"
            :count="scope.row.passCount"
            allow-rework-actions
            @changed="getList"
          />
        </template>
      </el-table-column>
      <el-table-column label="异常数" width="90">
        <template #default="scope">
          <DeviceStatusCount
            :row="scope.row"
            type="abnormal"
            :count="scope.row.abnormalCount"
            allow-rework-actions
            @changed="getList"
          />
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
            v-auth="btnAuth.delete"
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
      class="batch-scan-dialog"
      destroy-on-close
      @opened="focusScanInput"
    >
      <div v-if="batchScanOrder" class="batch-scan-layout">
        <div class="batch-scan-fixed">
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
              <el-select
                v-model="batchScanTarget"
                filterable
                style="width: 320px"
              >
                <el-option
                  v-for="b in batchScanExistingBatches"
                  :key="b.ID"
                  :label="b.batchNumber"
                  :value="b.batchNumber"
                />
                <el-option
                  :label="newBatchOption"
                  :value="newBatchOption"
                />
              </el-select>
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
        </div>

        <div class="scan-board">
          <div class="scan-basket">
            <div class="scan-title">批次</div>
            <div class="scan-list">
              <el-empty
                v-if="!existingBatchDevices.length && !scanBasket.length"
                description="还没有加入设备"
              />
              <template v-else>
                <div
                  v-for="(item, index) in existingBatchDevices"
                  :key="'exist-' + item.sn"
                  class="scan-existing-item"
                >
                  <span>{{ index + 1 }}. {{ item.sn }}</span>
                </div>
                <div
                  v-for="(item, index) in scanBasket"
                  :key="item.sn"
                  class="scan-item"
                >
                  <span>{{ existingBatchDevices.length + index + 1 }}. {{ item.sn }}</span>
                  <el-button
                    type="danger"
                    link
                    size="small"
                    @click="removeScanItem(item.sn)"
                    >移除</el-button
                  >
                </div>
              </template>
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
      class="production-detail-dialog"
      destroy-on-close
    >
      <div v-if="detailOrder" class="production-detail-layout">
        <div class="production-detail-fixed">
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
            <el-descriptions-item label="未分批">{{
              detailOrder.unbatchedCount ?? unbatchedDevices.length
            }}</el-descriptions-item>
            <el-descriptions-item label="设备数">{{
              detailOrder.deviceCount ?? detailOrder.devices?.length ?? 0
            }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="production-detail-scroll">
          <el-collapse v-model="detailActivePanels">
            <el-collapse-item
              v-for="batch in sortedDetailBatches"
              :key="batch.ID"
              :name="`batch-${batch.ID}`"
            >
              <template #title>
                <div class="detail-collapse-title">
                  <span>批次: {{ batch.batchNumber || '-' }}</span>
                  <el-tag size="small" type="info"
                    >{{ batch.devices?.length || 0 }} 台</el-tag
                  >
                  <el-tag v-if="batch.template" size="small" type="info">
                    模板: {{ batch.template.name }}
                  </el-tag>
                  <el-tag size="small" :type="orderStatusTagType(batch.status)">
                    {{ batchStatusLabel(batch.status) }}
                  </el-tag>
                  <span
                    v-if="batch.lastOperatorName"
                    class="batch-operator"
                  >操作人: {{ batch.lastOperatorName }}</span
                  >
                </div>
              </template>
              <div class="detail-section-actions">
                <el-button
                  size="small"
                  type="success"
                  link
                  @click.stop="onExportBatchExcel(batch)"
                >
                  导出Excel
                </el-button>
                <el-button
                  size="small"
                  type="primary"
                  link
                  @click.stop="openBatchPrint(batch)"
                >
                  打印
                </el-button>
                <el-button
                  size="small"
                  type="primary"
                  link
                  @click.stop="openFlowLogs({ batch })"
                >
                  流转日志
                </el-button>
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
                <el-table-column
                  prop="firmwareVersion"
                  label="固件"
                  width="110"
                />
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
                      v-if="batch.status !== 4"
                      type="danger"
                      link
                      size="small"
                      @click="handleRemoveFromBatch(scope.row)"
                    >
                      移出批次
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
            </el-collapse-item>

            <el-collapse-item v-if="unbatchedDevices.length" name="unbatched">
              <template #title>
                <div class="detail-collapse-title">
                  <span>无批次设备</span>
                  <el-tag size="small" type="warning"
                    >{{ unbatchedDevices.length }} 台</el-tag
                  >
                </div>
              </template>
              <el-table
                :data="unbatchedDevices"
                border
                size="small"
                class="mt-1"
              >
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
                <el-table-column
                  prop="firmwareVersion"
                  label="固件"
                  width="110"
                />
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
                      type="success"
                      link
                      size="small"
                      @click="openBatchPicker(scope.row)"
                    >
                      加入批次
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
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </el-dialog>

    <el-dialog
      v-model="batchPickerVisible"
      title="选择目标批次"
      width="450px"
      destroy-on-close
    >
      <el-radio-group v-model="pickedBatchID" class="batch-picker-group">
        <el-radio
          v-for="b in availableBatches"
          :key="b.ID"
          :value="b.ID"
          class="batch-picker-item"
        >
          {{ b.batchNumber }}
          <el-tag size="small" :type="batchStatusTag(b.status)" class="ml-1">
            {{ batchStatusLabel(b.status) }}
          </el-tag>
          <span class="batch-device-count">{{ b.deviceCount || b.devices?.length || 0 }} 台</span>
        </el-radio>
      </el-radio-group>
      <div v-if="availableBatches.length === 0" class="text-muted">没有可用的批次（所有批次均已完成）</div>
      <template #footer>
        <el-button @click="batchPickerVisible = false">取消</el-button>
        <el-button type="primary" :disabled="!pickedBatchID" @click="confirmAddToBatch">确认加入</el-button>
      </template>
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
  import { useBtnAuth } from '@/utils/btnAuth'
  import {
    getProductionOrderList,
    forceDeleteProductionOrder,
    updateProductionOrder,
    findProductionOrder,
    confirmReworkReceived,
    confirmReworkDone,
    scanAssignBatch,
    addDevicesToBatch,
    removeDeviceFromBatch
  } from '@/plugin/inspection/api/production_order'
  import { getTemplateList } from '@/plugin/inspection/api/template'
  import {
    assignOrderTemplate,
    exportInspectionExcel,
    getFlowLogs
  } from '@/plugin/inspection/api/work_order'
  import DeviceStatusCount from '@/plugin/inspection/components/DeviceStatusCount.vue'
  import FlowLogDrawer from '@/plugin/inspection/components/FlowLogDrawer.vue'

  const btnAuth = useBtnAuth()
  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const drawerVisible = ref(false)
  const batchScanVisible = ref(false)
  const detailVisible = ref(false)
  const dispatchVisible = ref(false)
  const batchPickerVisible = ref(false)
  const pickedBatchID = ref(null)
  const pendingDevice = ref(null)
  const formRef = ref(null)
  const scanInputRef = ref(null)
  const detailOrder = ref(null)
  const detailActivePanels = ref([])
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
    scanSN: ''
  })
  const batchScanTarget = ref('')
  const scanBasket = ref([])
  const submitDateRange = ref([])

  const searchInfo = reactive({
    moNumber: '',
    model: '',
    batchNumber: '',
    sn: '',
    instrumentCategory: '',
    status: undefined,
    batchComplete: undefined,
    hasAbnormal: undefined,
    startSubmitDate: '',
    endSubmitDate: '',
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
    ({ 0: '未派检', 1: '待检测接收', 2: '检测中', 3: '检测中', 4: '已完成' }[
      value
    ] || value)
  const batchStatusTag = (value) =>
    ({ 0: 'info', 1: 'warning', 2: 'primary', 3: 'warning', 4: 'success' }[
      value
    ] || 'info')
  const orderStatusTagType = (value) =>
    ({ 0: 'info', 1: 'warning', 2: 'primary', 3: 'warning', 4: 'success' }[
      value
    ] || 'info')
  const productionStatusLabel = (value) =>
    Number(value) === 4 ? '已完成' : '未完成'
  const productionStatusTagType = (value) =>
    Number(value) === 4 ? 'success' : 'warning'
  const orderStatusOptions = [
    { label: '未完成', value: 0 },
    { label: '已完成', value: 4 }
  ]
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

  const sortedDetailBatches = computed(() =>
    [...(detailOrder.value?.batches || [])].sort((a, b) => {
      const bTime = new Date(b.CreatedAt || b.createdAt || 0).getTime()
      const aTime = new Date(a.CreatedAt || a.createdAt || 0).getTime()
      if (bTime !== aTime) return bTime - aTime
      return Number(b.ID || 0) - Number(a.ID || 0)
    })
  )

  const availableBatches = computed(() =>
    (detailOrder.value?.batches || []).filter((b) => b.status !== 4)
  )

  const openBatchPicker = (device) => {
    pendingDevice.value = device
    pickedBatchID.value = null
    batchPickerVisible.value = true
  }

  const confirmAddToBatch = async () => {
    if (!pickedBatchID.value || !pendingDevice.value) return
    try {
      const res = await addDevicesToBatch({
        batchID: pickedBatchID.value,
        sns: [pendingDevice.value.sn]
      })
      if (res.code !== 0) return
      ElMessage.success(`已将 ${pendingDevice.value.sn} 加入批次`)
      batchPickerVisible.value = false
      await refreshDetail()
    } catch (e) {
      ElMessage.error('操作失败')
    }
  }

  const handleRemoveFromBatch = async (device) => {
    try {
      await ElMessageBox.confirm(
        `确定将 ${device.sn} 移出当前批次吗？`,
        '移出批次',
        { type: 'warning', confirmButtonText: '确定' }
      )
    } catch {
      return
    }
    try {
      const res = await removeDeviceFromBatch({ deviceID: device.ID })
      if (res.code !== 0) return
      ElMessage.success(`已将 ${device.sn} 移出批次`)
      await refreshDetail()
    } catch (e) {
      ElMessage.error('操作失败')
    }
  }

  const refreshDetail = async () => {
    if (!detailOrder.value) return
    const res = await findProductionOrder({ id: detailOrder.value.ID })
    if (res.code === 0) detailOrder.value = res.data
  }

  const batchScanExistingBatches = computed(() =>
    (batchScanOrder.value?.batches || [])
      .filter((b) => b.status !== 4)
      .sort((a, b) => new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime())
  )

  const newBatchOption = computed(() => {
    if (!batchScanOrder.value) return ''
    return previewNextBatchNumber(batchScanOrder.value)
  })

  const selectedExistingBatch = computed(() =>
    (batchScanOrder.value?.batches || []).find(
      (b) => b.batchNumber === batchScanTarget.value
    )
  )

  const existingBatchDevices = computed(() =>
    selectedExistingBatch.value?.devices || []
  )

  const existingBatchSNs = computed(() =>
    new Set(existingBatchDevices.value.map((d) => d.sn))
  )

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
    const filename = `${detailOrder.value?.moNumber || 'MO'}-${
      batch.batchNumber || batch.ID
    }-检测工单.xlsx`
    downloadBlob(res.data || res, filename)
  }

  const getList = async () => {
    loading.value = true
    try {
      searchInfo.startSubmitDate = submitDateRange.value?.[0] || ''
      searchInfo.endSubmitDate = submitDateRange.value?.[1] || ''
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
    searchInfo.batchNumber = ''
    searchInfo.sn = ''
    searchInfo.instrumentCategory = ''
    searchInfo.status = undefined
    searchInfo.batchComplete = undefined
    searchInfo.hasAbnormal = undefined
    searchInfo.startSubmitDate = ''
    searchInfo.endSubmitDate = ''
    submitDateRange.value = []
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
    detailActivePanels.value = []
    detailVisible.value = true
  }

  const openBatchScan = async (row) => {
    const res = await findProductionOrder({ id: row.ID })
    if (res.code !== 0) return
    batchScanOrder.value = res.data
    batchScanTarget.value = previewNextBatchNumber(res.data)
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
    if (existingBatchSNs.value.has(sn)) {
      ElMessage.warning(`SN ${sn} 已在此批次中`)
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
      batchNumber: batchScanTarget.value,
      sns: scanBasket.value.map((item) => item.sn)
    })
    if (res.code !== 0) return
    ElMessage.success('分批成功')
    const refresh = await findProductionOrder({ id: batchScanOrder.value.ID })
    if (refresh.code === 0) {
      batchScanOrder.value = refresh.data
      batchScanTarget.value = previewNextBatchNumber(refresh.data)
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

  :global(.batch-scan-dialog.el-dialog) {
    height: min(860px, calc(100vh - 48px));
    margin-top: 24px !important;
    margin-bottom: 24px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  :global(.batch-scan-dialog .el-dialog__header),
  :global(.batch-scan-dialog .el-dialog__footer) {
    flex: none;
  }

  :global(.batch-scan-dialog .el-dialog__body) {
    flex: 1;
    min-height: 0;
    overflow: hidden;
    padding-top: 12px;
    padding-bottom: 0;
  }

  .batch-scan-layout {
    height: 100%;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .batch-scan-fixed {
    flex: none;
  }

  :global(.production-detail-dialog.el-dialog) {
    height: min(860px, calc(100vh - 48px));
    margin-top: 24px !important;
    margin-bottom: 24px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  :global(.production-detail-dialog .el-dialog__header),
  :global(.production-detail-dialog .el-dialog__footer) {
    flex: none;
  }

  :global(.production-detail-dialog .el-dialog__body) {
    flex: 1;
    min-height: 0;
    overflow: hidden;
    padding-top: 12px;
  }

  .production-detail-layout {
    height: 100%;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .production-detail-fixed {
    flex: none;
  }

  .production-detail-scroll {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding-right: 4px;
  }

  .detail-collapse-title {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
    flex-wrap: wrap;
    font-weight: 600;
  }

  .batch-operator {
    color: var(--el-text-color-secondary, #909399);
    font-size: 12px;
    font-weight: 400;
  }

  .detail-section-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-bottom: 8px;
  }

  .scan-board {
    flex: 1;
    min-height: 0;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 14px;
    overflow: hidden;
  }

  .scan-basket,
  .scan-waiting {
    display: flex;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;
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
    flex: none;
    margin-bottom: 8px;
    font-size: 12px;
    color: var(--el-text-color-secondary, #909399);
  }

  .scan-list {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
    overflow-y: auto;
  }

  .scan-existing-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 7px 10px;
    border-radius: 6px;
    background: var(--el-fill-color-light, #f5f7fa);
    border: 1px solid var(--el-border-color-lighter, #ebeef5);
    color: var(--el-text-color-secondary, #909399);
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

  .batch-picker-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 100%;
  }

  .batch-picker-item {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 8px 12px;
    border: 1px solid var(--el-border-color-lighter, #ebeef5);
    border-radius: 6px;
    margin-right: 0;
  }

  .batch-device-count {
    margin-left: auto;
    font-size: 12px;
    color: var(--el-text-color-secondary, #909399);
  }

  .text-muted {
    color: var(--el-text-color-secondary, #909399);
    text-align: center;
    padding: 20px;
  }
</style>
