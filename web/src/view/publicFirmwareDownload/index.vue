<template>
  <div class="public-firmware-page">
    <div class="page-shell">
      <section class="hero-panel">
        <div class="hero-copy">
          <div class="hero-eyebrow">公开版本下载</div>
          <h1>Alpha 版本下载中心</h1>
          <p>请先选择设备类别和型号，查看对应的上传包列表</p>
        </div>
      </section>

      <section class="filter-panel">
        <div class="panel-title-block">
          <div class="panel-title">选择设备</div>
        </div>
        <div class="filter-row">
          <el-select
            v-model="selectedCategoryId"
            filterable
            placeholder="选择设备类别"
            class="filter-select"
            @change="handleCategoryChange"
          >
            <el-option
              v-for="item in categories"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
          <el-select
            v-model="selectedModelId"
            filterable
            placeholder="选择设备型号"
            class="filter-select"
            @change="handleModelChange"
          >
            <el-option
              v-for="item in models"
              :key="item.ID"
              :label="item.modelName"
              :value="item.ID"
            />
          </el-select>
          <!-- <el-button :loading="loading" @click="loadPage">刷新</el-button> -->
        </div>
      </section>

      <Transition name="panel-swap" mode="out-in">
        <section v-if="loading" key="loading" class="loading-panel">
          <el-skeleton :rows="6" animated />
        </section>

        <div v-else key="content" class="content-stage">
          <section v-if="!hasSelectedDeviceModels" class="empty-panel">
            <el-empty description="请先选择设备和型号" />
          </section>

          <section v-else-if="packageCards.length" class="version-list-panel">
            <div class="panel-head version-list-head">
              <div>
                <div class="panel-kicker">上传包列表</div>
              </div>
              <el-tag type="info">{{ packageCards.length }} 个上传包</el-tag>
            </div>

            <div class="version-list">
              <article
                v-for="(item, index) in packageCards"
                :key="item.logId"
                class="download-card"
                :style="{ '--card-index': index }"
              >
                <div class="download-card-head">
                  <div class="download-card-title">
                    <span class="version-code-text">
                      {{ item.firmware?.versionCode || '-' }}
                    </span>
                    <span class="version-name version-name-inline">
                      {{ item.firmware?.versionName || '-' }}
                    </span>
                    <span
                      v-if="item.packageSize > 0"
                      class="version-size version-size-inline"
                    >
                      {{ formatPackageSize(item.packageSize) }}
                    </span>
                    <el-tag v-if="item.isRecommended" type="success">推荐版</el-tag>
                    <el-tag v-if="item.firmware?.isLatest" type="danger">最新版</el-tag>
                    <el-tag
                      v-if="!item.isRecommended && !item.firmware?.isLatest"
                      type="info"
                    >
                      历史版本
                    </el-tag>
                  </div>
                  <el-button
                    type="primary"
                    size="small"
                    :disabled="!downloadable(item)"
                    @click="downloadVersion(item)"
                  >
                    下载
                  </el-button>
                </div>

                <div class="download-card-body">
                  <div class="download-card-meta">
                    <div class="meta-item compact-meta">
                      <span>设备类别</span>
                      <strong>{{ item.category?.name || '-' }}</strong>
                    </div>
                    <div class="meta-item compact-meta">
                      <span>设备型号</span>
                      <strong>{{
                        (item.modelNames && item.modelNames.length
                          ? item.modelNames.join('、')
                          : item.model?.modelName) || '-'
                      }}</strong>
                    </div>
                    <div class="meta-item compact-meta">
                      <span>上传时间</span>
                      <strong>{{
                        formatDate(item.operateAt || item.firmware?.uploadedAt)
                      }}</strong>
                    </div>
                    <div
                      class="meta-item compact-meta"
                      v-if="item.firmware?.checksum"
                    >
                      <span>校验值</span>
                      <strong class="mono">{{ item.firmware.checksum }}</strong>
                    </div>
                  </div>

                  <div class="release-note">
                    <div class="release-note-title">包说明</div>
                    <div class="release-note-body">
                      {{ item.firmware?.releaseNote || '暂无版本说明' }}
                    </div>
                  </div>
                </div>
              </article>
            </div>
          </section>

          <section v-else class="empty-panel">
            <el-empty description="暂无可下载包" />
          </section>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import { getUrl } from '@/utils/image'
  import { getPublicFirmwareDownloadPage } from '@/api/publicFirmware'

  defineOptions({
    name: 'PublicFirmwareDownload'
  })

  const route = useRoute()
  const router = useRouter()

  const loading = ref(false)
  const MIN_LOADING_DURATION = 240
  let loadSequence = 0
  const pageData = ref({
    categories: [],
    models: [],
    packages: []
  })
  const selectedCategoryId = ref('')
  const selectedModelId = ref('')

  const categories = computed(() => pageData.value.categories || [])
  const models = computed(() => pageData.value.models || [])
  const hasSelectedDeviceModels = computed(
    () => !!selectedCategoryId.value && !!selectedModelId.value
  )

  const packageCards = computed(() =>
    [...(pageData.value.packages || [])].sort((a, b) => {
      const recommendedA = !!a?.isRecommended
      const recommendedB = !!b?.isRecommended
      if (recommendedA !== recommendedB) {
        return recommendedA ? -1 : 1
      }
      const timeA = new Date(a?.operateAt || a?.firmware?.uploadedAt || 0).getTime()
      const timeB = new Date(b?.operateAt || b?.firmware?.uploadedAt || 0).getTime()
      if (timeA !== timeB) return timeB - timeA
      return (b?.logId || 0) - (a?.logId || 0)
    })
  )
  const normalizeSelectId = (value) => {
    if (Array.isArray(value)) {
      return normalizeSelectId(value[0])
    }
    if (value === undefined || value === null || value === '') {
      return ''
    }
    const id = Number(value)
    return Number.isFinite(id) && id > 0 ? id : ''
  }

  const loadPage = async (params = {}) => {
    const sequence = ++loadSequence
    const startedAt = Date.now()
    loading.value = true
    try {
      const categoryId = normalizeSelectId(
        params.categoryId ?? selectedCategoryId.value ?? ''
      )
      const modelId = normalizeSelectId(
        params.modelId ?? params.modelIds ?? selectedModelId.value ?? ''
      )
      const query = {
        categoryId,
        modelId
      }
      const res = await getPublicFirmwareDownloadPage(query)
      if (sequence !== loadSequence) {
        return
      }
      const data = res?.data || {}
      pageData.value = {
        ...data,
        packages: data.packages || []
      }
      selectedCategoryId.value = normalizeSelectId(
        data.selectedCategoryId || categoryId || ''
      )
      selectedModelId.value = normalizeSelectId(
        data.selectedModelId || modelId || ''
      )

      const nextQuery = {
        categoryId: selectedCategoryId.value || undefined,
        modelId: selectedModelId.value || undefined
      }
      const currentQuery = {
        categoryId: route.query.categoryId || undefined,
        modelId: route.query.modelId || route.query.modelIds || undefined
      }
      if (
        String(currentQuery.categoryId || '') !==
          String(nextQuery.categoryId || '') ||
        String(currentQuery.modelId || '') !== String(nextQuery.modelId || '')
      ) {
        await router.replace({ name: route.name, query: nextQuery })
      }
    } catch (error) {
      if (sequence !== loadSequence) {
        return
      }
      ElMessage.error(error?.message || '获取公开版本数据失败')
    } finally {
      if (sequence !== loadSequence) {
        return
      }
      const elapsed = Date.now() - startedAt
      if (elapsed < MIN_LOADING_DURATION) {
        await new Promise((resolve) =>
          setTimeout(resolve, MIN_LOADING_DURATION - elapsed)
        )
      }
      if (sequence === loadSequence) {
        loading.value = false
      }
    }
  }

  const handleCategoryChange = async () => {
    selectedModelId.value = ''
    await loadPage({
      categoryId: selectedCategoryId.value,
      modelId: ''
    })
  }

  const handleModelChange = async () => {
    await loadPage({
      categoryId: selectedCategoryId.value,
      modelId: selectedModelId.value
    })
  }

  const downloadable = (item) => !!item?.firmware?.packageUrl

  const extractFileNameFromDisposition = (disposition) => {
    if (!disposition) return ''
    const utf8Match = disposition.match(/filename\*=UTF-8''([^;]+)/i)
    if (utf8Match?.[1]) {
      try {
        return decodeURIComponent(utf8Match[1])
      } catch (error) {
        return utf8Match[1]
      }
    }
    const normalMatch = disposition.match(/filename="?([^";]+)"?/i)
    return normalMatch?.[1] || ''
  }

  const downloadVersion = async (item) => {
    const url = item?.firmware?.packageUrl
    if (!url) {
      ElMessage.warning('当前上传包没有可下载的安装包')
      return
    }
    const loadingMessage = ElMessage({
      message: '下载中...',
      type: 'info',
      duration: 0,
      showClose: false
    })

    try {
      const response = await fetch(getUrl(url), {
        credentials: 'include'
      })
      if (!response.ok) {
        let errorMessage = '下载失败'
        try {
          const text = await response.text()
          if (text) {
            try {
              const json = JSON.parse(text)
              errorMessage = json?.msg || errorMessage
            } catch (error) {
              errorMessage = text
            }
          }
        } catch (error) {
          // ignore parse error
        }
        throw new Error(errorMessage)
      }

      const blob = await response.blob()
      const contentType = String(
        response.headers.get('content-type') || blob.type || ''
      ).toLowerCase()
      const disposition = String(response.headers.get('content-disposition') || '')
      const isErrorBlob =
        contentType.includes('application/json') ||
        contentType.includes('text/plain')
      if (!blob.size || isErrorBlob) {
        let errorMessage = '下载失败'
        try {
          const text = await blob.text()
          if (text) {
            try {
              const json = JSON.parse(text)
              errorMessage = json?.msg || errorMessage
            } catch (error) {
              errorMessage = text
            }
          }
        } catch (error) {
          // ignore parse error
        }
        throw new Error(errorMessage)
      }

      const downloadName =
        extractFileNameFromDisposition(disposition) ||
        item?.firmware?.packageName ||
        `${item?.firmware?.versionCode || item?.firmware?.versionName || 'firmware'}.bin`
      const objectUrl = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = objectUrl
      link.download = downloadName
      link.rel = 'noopener'
      link.style.display = 'none'
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.setTimeout(() => {
        window.URL.revokeObjectURL(objectUrl)
      }, 1000)
      ElMessage.success('下载成功')
    } catch (error) {
      ElMessage.error(error?.message || '下载失败')
    } finally {
      loadingMessage?.close?.()
    }
  }

  const formatPackageSize = (size) => {
    const value = Number(size || 0)
    if (!value) return '-'
    if (value < 1024) return `${value} B`
    if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
    if (value < 1024 * 1024 * 1024) {
      return `${(value / (1024 * 1024)).toFixed(1)} MB`
    }
    return `${(value / (1024 * 1024 * 1024)).toFixed(1)} GB`
  }

  onMounted(async () => {
    const query = {
      categoryId: normalizeSelectId(route.query.categoryId || ''),
      modelId: normalizeSelectId(route.query.modelId || route.query.modelIds)
    }
    await loadPage(query)
  })
</script>

<style scoped>
  .public-firmware-page {
    height: 100%;
    min-height: 100%;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    background: radial-gradient(
        circle at top right,
        rgba(59, 130, 246, 0.16),
        transparent 30%
      ),
      radial-gradient(
        circle at top left,
        rgba(16, 185, 129, 0.12),
        transparent 28%
      ),
      linear-gradient(
        180deg,
        #0b1220 0%,
        #132238 26%,
        #f5f7fb 26%,
        #f5f7fb 100%
      );
    padding: 24px 18px 56px;
  }

  .page-shell {
    max-width: 1240px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 18px;
    min-height: 100%;
  }

  .hero-panel {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 28px;
    border-radius: 24px;
    color: #f8fafc;
    background: linear-gradient(
        135deg,
        rgba(15, 23, 42, 0.92),
        rgba(30, 41, 59, 0.88)
      ),
      linear-gradient(135deg, #0f172a, #1d4ed8);
    box-shadow: 0 24px 80px rgba(15, 23, 42, 0.28);
  }

  .hero-copy h1 {
    margin: 10px 0 12px;
    font-size: clamp(28px, 4vw, 48px);
    line-height: 1.05;
  }

  .hero-copy p {
    margin: 0;
    max-width: 640px;
    color: rgba(226, 232, 240, 0.88);
    font-size: 14px;
    line-height: 1.8;
  }

  .hero-eyebrow {
    display: inline-flex;
    align-items: center;
    padding: 6px 12px;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.1);
    color: #bfdbfe;
    font-size: 12px;
    letter-spacing: 0.08em;
  }

  .filter-panel,
  .primary-panel,
  .secondary-panel,
  .loading-panel,
  .empty-panel {
    background: #fff;
    border-radius: 22px;
    box-shadow: 0 16px 40px rgba(15, 23, 42, 0.08);
    border: 1px solid rgba(148, 163, 184, 0.16);
    padding: 22px;
  }

  .content-stage {
    display: flex;
    flex-direction: column;
    gap: 18px;
    min-width: 0;
  }

  .panel-title-block {
    margin-bottom: 16px;
  }

  .panel-title {
    color: #0f172a;
    font-size: 18px;
    font-weight: 700;
  }

  .panel-kicker {
    color: #0f172a;
    font-size: 16px;
    font-weight: 700;
  }

  .panel-subtitle {
    margin-top: 6px;
    color: #64748b;
    font-size: 13px;
    line-height: 1.6;
  }

  .filter-row {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    align-items: center;
  }

  .filter-select {
    min-width: 240px;
    flex: 1 1 240px;
  }

  .panel-head {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    align-items: flex-start;
    margin-bottom: 18px;
    flex-wrap: wrap;
  }

  .panel-head.compact {
    margin-bottom: 0;
  }

  .version-list-head {
    align-items: center;
  }

  .panel-swap-enter-active,
  .panel-swap-leave-active {
    transition: opacity 0.24s ease, transform 0.24s ease;
  }

  .panel-swap-enter-from,
  .panel-swap-leave-to {
    opacity: 0;
    transform: translateY(10px);
  }

  .panel-swap-enter-to,
  .panel-swap-leave-from {
    opacity: 1;
    transform: translateY(0);
  }

  @keyframes card-rise {
    from {
      opacity: 0;
      transform: translateY(10px);
    }

    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .version-title {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
  }

  .version-code {
    display: inline-flex;
    align-items: baseline;
    gap: 8px;
    color: #0f172a;
    font-size: clamp(24px, 3vw, 36px);
    font-weight: 800;
  }

  .version-name {
    color: #334155;
    font-size: 16px;
    font-weight: 600;
  }

  .tag-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin: 16px 0 18px;
  }

  .meta-grid {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 12px;
  }

  .meta-item {
    padding: 14px 15px;
    border-radius: 14px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
  }

  .meta-item span {
    display: block;
    color: #64748b;
    font-size: 12px;
    margin-bottom: 6px;
  }

  .meta-item strong {
    display: block;
    color: #0f172a;
    font-size: 13px;
    line-height: 1.6;
    word-break: break-word;
  }

  .mono {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
      monospace;
  }

  .release-note {
    margin-top: 18px;
    padding: 16px 18px;
    border-radius: 16px;
    background: linear-gradient(180deg, #f8fafc, #eef6ff);
    border: 1px solid #dbeafe;
  }

  .release-note-title {
    color: #0f172a;
    font-size: 13px;
    font-weight: 700;
    margin-bottom: 8px;
  }

  .release-note-body {
    color: #334155;
    line-height: 1.8;
    white-space: pre-wrap;
  }

  .version-list-panel {
    background: #fff;
    border-radius: 22px;
    box-shadow: 0 16px 40px rgba(15, 23, 42, 0.08);
    border: 1px solid rgba(148, 163, 184, 0.16);
    padding: 22px;
  }

  .version-list {
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .download-card {
    padding: 20px 20px 18px;
    border-radius: 20px;
    border: 1px solid #e2e8f0;
    background: linear-gradient(180deg, #fcfdff 0%, #f8fbff 100%);
    box-shadow: none;
    animation: card-rise 0.38s ease both;
    animation-delay: calc(var(--card-index, 0) * 60ms);
  }

  .download-card.is-recommended {
    border-color: #f4d6a8;
    box-shadow: 0 18px 48px rgba(245, 158, 11, 0.08);
  }

  .download-card.is-latest {
    border-color: #f6c4c4;
  }

  .download-card.is-history {
    border-color: #dbe3f0;
  }

  .download-card-head {
    display: flex;
    justify-content: space-between;
    gap: 14px;
    align-items: flex-start;
    margin-bottom: 14px;
  }

  .download-card-title {
    display: inline-flex;
    align-items: baseline;
    gap: 10px;
    flex-wrap: wrap;
    line-height: 1.1;
    flex: 1 1 auto;
    min-width: 0;
  }

  .version-code-text {
    font-size: clamp(26px, 4vw, 40px);
    font-weight: 800;
    color: #0f172a;
    line-height: 1;
  }

  .version-name-inline {
    font-size: 15px;
    font-weight: 600;
    color: #475569;
  }

  .version-size-inline {
    font-size: 12px;
  }

  .download-card-title :deep(.el-tag) {
    height: 24px;
    line-height: 22px;
  }

  .download-card-body {
    display: grid;
    gap: 14px;
  }

  .download-card-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 10px 12px;
    align-items: center;
  }

  .compact-meta {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    border-radius: 999px;
    background: #f8fafc;
  }

  .compact-meta span {
    margin-bottom: 0;
    display: inline;
    font-size: 12px;
    white-space: nowrap;
  }

  .compact-meta strong {
    font-size: 12px;
    display: inline;
    white-space: nowrap;
  }

  .download-card .release-note {
    margin-top: 0;
  }

  .secondary-grid,
  .history-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px;
  }

  .version-card,
  .history-card {
    padding: 18px;
    border-radius: 18px;
    border: 1px solid #e2e8f0;
    background: linear-gradient(180deg, #ffffff, #f8fafc);
  }

  .version-card-head {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    align-items: flex-start;
    margin-bottom: 14px;
  }

  .version-card-title {
    color: #0f172a;
    font-size: 16px;
    font-weight: 700;
  }

  .version-card-subtitle {
    margin-top: 4px;
    color: #64748b;
    font-size: 12px;
  }

  .version-code-line {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .version-size {
    color: #0f172a;
    font-size: 12px;
    font-weight: 600;
    padding: 2px 8px;
    border-radius: 999px;
    background: #e2e8f0;
  }

  .version-card-body {
    display: grid;
    gap: 10px;
  }

  .card-line {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    font-size: 13px;
    line-height: 1.6;
    color: #334155;
  }

  .card-line span {
    color: #64748b;
    flex: 0 0 auto;
  }

  .card-line strong {
    text-align: right;
    word-break: break-word;
  }

  .version-card-footer {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    align-items: center;
    margin-top: 16px;
  }

  .history-list {
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .history-note {
    color: #334155;
    font-weight: 400;
    text-align: right;
  }

  @media (max-width: 960px) {
    .hero-panel,
    .meta-grid,
    .secondary-grid,
    .history-grid {
      grid-template-columns: 1fr;
    }

    .filter-select {
      min-width: 100%;
    }

    .panel-head {
      align-items: stretch;
    }
  }

  @media (max-width: 640px) {
    .public-firmware-page {
      padding: 14px 12px 24px;
    }

    .page-shell {
      gap: 12px;
    }

    .hero-panel,
    .filter-panel,
    .version-list-panel,
    .loading-panel,
    .empty-panel {
      padding: 16px;
      border-radius: 18px;
    }

    .download-card {
      padding: 16px;
      border-radius: 16px;
    }

    .download-card-head {
      align-items: flex-start;
    }

    .download-card-head .el-button {
      margin-top: 2px;
      margin-left: auto;
      flex: 0 0 auto;
    }

    .download-card-meta {
      flex-direction: row;
      flex-wrap: wrap;
      align-items: stretch;
      gap: 8px;
    }

    .compact-meta {
      flex: 1 1 140px;
      min-width: 0;
    }

    .card-line {
      flex-direction: column;
      align-items: flex-start;
    }

    .card-line strong {
      text-align: left;
    }
  }
</style>

