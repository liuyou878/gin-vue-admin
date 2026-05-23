<template>
  <div>
    <div class="gva-search-box">
      <el-form :model="searchInfo" inline>
        <el-form-item label="生产号">
          <el-input
            v-model="searchInfo.moNumber"
            placeholder="请输入MO号"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="批次号">
          <el-input
            v-model="searchInfo.batchNumber"
            placeholder="请输入批次号"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="SN">
          <el-input
            v-model="searchInfo.sn"
            placeholder="请输入SN"
            clearable
            size="small"
          />
        </el-form-item>
        <el-form-item label="型号">
          <el-input
            v-model="searchInfo.model"
            placeholder="请输入型号"
            clearable
            size="small"
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
      <el-table-column prop="moNumber" label="生产号" min-width="140" />
      <el-table-column prop="batchNumber" label="批次号" min-width="180">
        <template #default="scope">{{ scope.row.batchNumber || '-' }}</template>
      </el-table-column>
      <el-table-column prop="sn" label="SN" min-width="140" />
      <el-table-column prop="model" label="型号" width="110" />
      <el-table-column prop="pnCode" label="PN码" min-width="150" />
      <el-table-column
        prop="firmwareVersion"
        label="固件版本"
        min-width="130"
      />
      <el-table-column
        prop="timeLicense"
        label="时间码"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column
        prop="regionLicense"
        label="围栏码"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column
        prop="ntripCode"
        label="Ntrip状态"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="statusTagType(scope.row.status)" size="small">
            {{ statusLabel(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="submitterName" label="提交人" width="100" />
      <el-table-column prop="CreatedAt" label="提交时间" width="170">
        <template #default="scope">{{
          formatDate(scope.row.CreatedAt)
        }}</template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="scope">
          <el-button
            type="primary"
            link
            size="small"
            @click="openDetail(scope.row)"
            >GETALL详情</el-button
          >
          <el-button
            v-auth="btnAuth.delete"
            type="danger"
            link
            size="small"
            @click="handleDelete(scope.row)"
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
        background
        small
        @size-change="getList"
        @current-change="getList"
      />
    </div>

    <el-drawer
      v-model="detailVisible"
      title="生产工具提交详情"
      size="60%"
      destroy-on-close
    >
      <template v-if="detailData">
        <div class="mb-3">
          <!-- <el-button v-auth="btnAuth.delete" type="danger" size="small" @click="handleDelete(detailData, true)">删除这条提交数据</el-button> -->
        </div>
        <el-descriptions :column="2" border size="small" class="mb-4">
          <el-descriptions-item label="生产号">{{
            detailData.productionOrder?.moNumber || detailData.moNumber || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="批次号">{{
            detailData.batch?.batchNumber || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="SN">{{
            detailData.sn || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="型号">{{
            detailData.model || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="PN码">{{
            detailData.pnCode || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="固件版本">{{
            detailData.firmwareVersion || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="主板固件版本">{{
            detailMainboardFirmwareVersion || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="时间码">{{
            detailData.timeLicense || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="围栏码">{{
            detailData.regionLicense || '-'
          }}</el-descriptions-item>
          <el-descriptions-item label="Ntrip状态">{{
            detailData.ntripCode || '-'
          }}</el-descriptions-item>
          <el-descriptions-item v-if="detailNetworkModel" label="网络型号">{{
            detailNetworkModel
          }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTagType(detailData.status)" size="small">
              {{ statusLabel(detailData.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item
            v-if="detailData.returnReason"
            label="打回原因"
            >{{ detailData.returnReason }}</el-descriptions-item
          >
          <el-descriptions-item v-if="detailData.returnByName" label="打回人">{{
            detailData.returnByName
          }}</el-descriptions-item>
          <el-descriptions-item v-if="detailData.returnAt" label="打回时间">{{
            formatDate(detailData.returnAt)
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="mb-2 font-600">GETALL 原始内容</div>
        <el-input
          :model-value="prettyDeviceInfo"
          type="textarea"
          :rows="20"
          readonly
        />
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
  import { computed, reactive, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import { useBtnAuth } from '@/utils/btnAuth'
  import {
    getSubmittedDeviceList,
    findSubmittedDevice,
    deleteSubmittedDevice
  } from '@/plugin/inspection/api/production_order'

  const btnAuth = useBtnAuth()
  const loading = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const detailVisible = ref(false)
  const detailData = ref(null)

  const searchInfo = reactive({
    moNumber: '',
    batchNumber: '',
    sn: '',
    model: '',
    page: 1,
    pageSize: 30
  })

  const prettyDeviceInfo = computed(() => {
    const raw = detailData.value?.deviceInfo
    if (!raw) return ''
    try {
      return JSON.stringify(JSON.parse(raw), null, 2)
    } catch {
      return raw
    }
  })

  const parsedDeviceInfo = computed(() => {
    const raw = detailData.value?.deviceInfo
    if (!raw) return {}
    try {
      return JSON.parse(raw)
    } catch {
      return {}
    }
  })

  const detailNetworkModel = computed(() => {
    const value =
      parsedDeviceInfo.value?.network?.model ||
      parsedDeviceInfo.value?.networkModel ||
      ''
    return String(value || '').trim()
  })

  const detailMainboardFirmwareVersion = computed(() => {
    const value =
      detailData.value?.mainboardFirmwareVersion ||
      parsedDeviceInfo.value?.device?.mainboardFirmwareVersion ||
      parsedDeviceInfo.value?.mainboardFirmwareVersion ||
      ''
    return String(value || '').trim()
  })

  const statusLabel = (status) =>
    ({
      pending: '待检测设备',
      pass: '合格',
      fail: '不合格',
      returned: '待生产接收',
      rework: '返工中',
      pending_recheck: '待复检'
    }[status] ||
    status ||
    '-')

  const statusTagType = (status) =>
    ({
      pass: 'success',
      fail: 'danger',
      returned: 'warning',
      rework: 'warning',
      pending_recheck: 'primary',
      pending: 'info'
    }[status] || 'info')

  const getList = async () => {
    loading.value = true
    try {
      const res = await getSubmittedDeviceList(searchInfo)
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
    searchInfo.batchNumber = ''
    searchInfo.sn = ''
    searchInfo.model = ''
    searchInfo.page = 1
    searchInfo.pageSize = 30
    getList()
  }

  const openDetail = async (row) => {
    const res = await findSubmittedDevice({ id: row.ID })
    if (res.code === 0) {
      detailData.value = res.data
      detailVisible.value = true
    }
  }

  const handleDelete = async (row, closeDetail = false) => {
    const label = [
      row.moNumber || row.productionOrder?.moNumber,
      row.batchNumber || row.batch?.batchNumber,
      row.sn
    ]
      .filter(Boolean)
      .join(' / ')
    try {
      await ElMessageBox.confirm(
        `确定删除这条生产提交数据吗？${label ? `\n${label}` : ''}`,
        '提示',
        { type: 'warning', confirmButtonText: '删除' }
      )
    } catch {
      return
    }

    const res = await deleteSubmittedDevice({ id: row.ID })
    if (res.code !== 0) return

    ElMessage.success('删除成功')
    if (closeDetail) {
      detailVisible.value = false
      detailData.value = null
    }
    if (tableData.value.length === 1 && searchInfo.page > 1) {
      searchInfo.page -= 1
    }
    getList()
  }

  getList()
</script>
