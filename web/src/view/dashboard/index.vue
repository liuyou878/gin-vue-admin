<template>
  <div class="blank-dashboard">
    <div class="blank-dashboard__tip">
      <!-- 页面可能存在缓存，请刷新后重试 -->
    </div>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'Dashboard'
  })
</script>

<style lang="scss" scoped>
  .blank-dashboard {
    display: flex;
    min-height: 100%;
    align-items: center;
    justify-content: center;
    background: transparent;
  }

  .blank-dashboard__tip {
    color: #909399;
    font-size: 14px;
  }
</style>
