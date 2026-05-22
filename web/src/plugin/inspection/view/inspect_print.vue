<template>
  <div class="print-page" v-if="loaded">
    <div class="no-print print-toolbar">
      <el-button @click="goBack">返回</el-button>
      <el-button type="primary" @click="printNow">打印</el-button>
    </div>

    <div class="sheet" :class="sheetDensityClass">
      <div class="title">GNSS接收机产品检测工单</div>

      <table class="info-table">
        <tbody>
          <tr>
            <td class="label">生产订单号</td>
            <td>{{ order.moNumber || '-' }}</td>
            <td class="label">批次号</td>
            <td>{{ order.batchNumber || '-' }}</td>
            <td class="label">业务类型</td>
            <td>{{ catLabel(order.instrumentCategory) || '-' }}</td>
          </tr>
          <tr>
            <td class="label">产品名称</td>
            <td>{{ order.productName || '-' }}</td>
            <td class="label">型号</td>
            <td>{{ order.model || '-' }}</td>
            <td class="label">PN码</td>
            <td>{{ order.pnCode || '-' }}</td>
          </tr>
          <tr>
            <td class="label">固件版本</td>
            <td>{{ order.firmwareVersion || '-' }}</td>
            <td class="label">主板固件版本</td>
            <td>{{ order.mainboardFirmwareVersion || '-' }}</td>
            <td class="label">检测员</td>
            <td>{{ order.inspectorName || '' }}</td>
          </tr>
        </tbody>
      </table>

      <table class="inspect-print-table">
        <thead>
          <tr>
            <th class="seq-col">序号</th>
            <th class="sn-col">机身码(SN)</th>
            <th class="result-col">检测结果</th>
            <th v-for="item in templateItems" :key="item.itemID" class="item-col">
              {{ item.itemName }}
            </th>
            <th class="remark-col">备注</th>
            <th class="sign-col">签名</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, index) in printRows" :key="row.ID || `empty-${index}`">
            <td>{{ index + 1 }}</td>
            <td class="sn-cell">{{ row.sn || '' }}</td>
            <td>{{ deviceResultLabel(row) }}</td>
            <td v-for="item in templateItems" :key="item.itemID">{{ resultText(row, item) }}</td>
            <td>{{ rowRemark(row) }}</td>
            <td></td>
          </tr>
        </tbody>
      </table>

      <div class="basis">
        仪器检测依据：企业标准：Q/440112014000AEFCHKJ 1-2025<br />
        计量检定规程参照：JJG 1200-2023<br />
        抽样检测参照：GB/T 2828.1-2021
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getInspectionDetail } from '@/plugin/inspection/api/work_order'

const route = useRoute()
const loaded = ref(false)
const detail = ref({ order: {}, devices: [], templateItems: [] })

const order = computed(() => detail.value.order || {})
const templateItems = computed(() => detail.value.templateItems || [])
const sheetDensityClass = computed(() => {
  const itemCount = templateItems.value.length
  if (itemCount >= 12) return 'sheet-compact'
  if (itemCount >= 8) return 'sheet-dense'
  return 'sheet-normal'
})
const printRows = computed(() => {
  const rows = [...(detail.value.devices || [])]
  while (rows.length < 8) {
    rows.push({})
  }
  return rows
})
const catLabel = (v) =>
  ({
    online: '线上',
    offline: '线下',
    foreign_trade: '外贸',
    custom: '定制款'
  }[v] || v)

const deviceResultLabel = (row) =>
  ({
    pass: '合格',
    fail: '不合格',
    pending: '未完成',
    returned: '待生产接收',
    rework: '返工中',
    pending_recheck: '待复检',
    rechecking: '复检中'
  }[row?.status || ''] || '')

const resultText = (row, item) => {
  if (!row?.results?.length || !item) return ''
  const result = row.results.find((r) => r.itemID === item.itemID)
  if (!result) return ''

  const parts = []
  if (result.resultType !== 'number') {
    if (result.passResult === true) parts.push('√')
    if (result.passResult === false) parts.push('×')
  }
  if (result.resultType !== 'pass_fail' && result.numberResult !== undefined && result.numberResult !== null && result.numberResult !== '') {
    parts.push(String(result.numberResult))
  }
  return parts.join(' ')
}

const rowRemark = (row) => {
  if (!row?.results?.length) return ''
  return row.results
    .filter((result) => result.remark)
    .map((result) => `${result.itemName}：${result.remark}`)
    .join('；')
}

const goBack = () => {
  window.close()
  if (!window.closed) {
    window.history.back()
  }
}

const printNow = () => {
  window.print()
}

const loadDetail = async () => {
  const batchId = route.query.batchId
  if (!batchId) return
  const res = await getInspectionDetail({ id: batchId })
  if (res.code === 0) {
    detail.value = res.data
    loaded.value = true
    await nextTick()
    setTimeout(() => window.print(), 300)
  }
}

loadDetail()
</script>

<style scoped>
.print-page {
  min-height: 100vh;
  background: #f3f4f6;
  padding: 16px;
  color: #111827;
}
.print-toolbar {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-bottom: 12px;
}
.sheet {
  width: 100%;
  max-width: 100%;
  background: #fff;
  padding: 12px;
  box-sizing: border-box;
  overflow: hidden;
}
.title {
  text-align: center;
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 10px;
}
.info-table,
.inspect-print-table {
  width: 100%;
  max-width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}
.info-table td,
.inspect-print-table th,
.inspect-print-table td {
  border: 1px solid #111;
  padding: 4px 5px;
  text-align: center;
  vertical-align: middle;
  font-size: 12px;
}
.info-table .label {
  width: 90px;
  font-weight: 700;
  background: #f3f4f6;
}
.inspect-print-table {
  margin-top: 8px;
}
.inspect-print-table th {
  height: 40px;
  font-weight: 700;
  background: #f3f4f6;
  word-break: break-all;
  overflow-wrap: anywhere;
}
.inspect-print-table td {
  height: 30px;
  overflow-wrap: anywhere;
}
.seq-col {
  width: 4.5%;
}
.sn-col {
  width: 15%;
}
.result-col {
  width: 8%;
}
.item-col {
  width: auto;
}
.remark-col,
.sign-col {
  width: 8%;
}
.sheet-dense .title {
  font-size: 18px;
}
.sheet-dense .info-table td,
.sheet-dense .inspect-print-table th,
.sheet-dense .inspect-print-table td {
  padding: 3px 4px;
  font-size: 11px;
}
.sheet-dense .seq-col {
  width: 4%;
}
.sheet-dense .sn-col {
  width: 13%;
}
.sheet-dense .result-col {
  width: 7%;
}
.sheet-dense .remark-col,
.sheet-dense .sign-col {
  width: 7%;
}
.sheet-compact .title {
  font-size: 16px;
  margin-bottom: 6px;
}
.sheet-compact .info-table td,
.sheet-compact .inspect-print-table th,
.sheet-compact .inspect-print-table td {
  padding: 2px 3px;
  font-size: 9px;
  line-height: 1.15;
}
.sheet-compact .info-table .label {
  width: 70px;
}
.sheet-compact .inspect-print-table th {
  height: 34px;
}
.sheet-compact .inspect-print-table td {
  height: 25px;
}
.sheet-compact .seq-col {
  width: 3.5%;
}
.sheet-compact .sn-col {
  width: 12%;
}
.sheet-compact .result-col {
  width: 6%;
}
.sheet-compact .remark-col,
.sheet-compact .sign-col {
  width: 6%;
}
.sheet-compact .basis {
  margin-top: 6px;
  font-size: 9px;
  line-height: 1.35;
}
.sn-cell {
  font-weight: 600;
}
.basis {
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.6;
}

@page {
  size: A4 landscape;
  margin: 8mm;
}

@media print {
  .no-print {
    display: none !important;
  }
  .print-page {
    padding: 0;
    background: #fff;
  }
  .sheet {
    padding: 0;
  }
  .sheet-compact {
    transform-origin: top left;
  }
  .inspect-print-table tr {
    break-inside: avoid;
  }
}
</style>
