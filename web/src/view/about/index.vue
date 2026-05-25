<template>
  <div class="blank-about">
    <div class="blank-about__tip">
      <!-- 页面可能存在缓存，请刷新后重试 -->
    </div>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'About'
  })
</script>

<style lang="scss" scoped>
  .blank-about {
    display: flex;
    min-height: 100%;
    align-items: center;
    justify-content: center;
    background: transparent;
  }

  .blank-about__tip {
    color: #909399;
    font-size: 14px;
  }
</style>
