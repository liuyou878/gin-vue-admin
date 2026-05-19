<template>
  <div />
</template>

<script setup>
import { onMounted } from 'vue'
onMounted(() => {
  const url = `${window.location.origin}${window.location.pathname}#/inspectWorkOrder`
  window.open(url, '_self')
})
</script>
