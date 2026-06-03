<template>
  <n-tag :type="tagType" size="small">{{ label }}</n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NTag } from 'naive-ui'

const props = defineProps<{
  status: 'success' | 'failed' | 'pending' | 'running' | 'never' | 'retrying' | 'timeout' | 'cancelled' | 'skipped' | ''
}>()

const tagType = computed(() => {
  const map: Record<string, 'success' | 'error' | 'warning' | 'info' | 'default'> = {
    success: 'success',
    failed: 'error',
    pending: 'info',
    running: 'warning',
    never: 'default',
    retrying: 'warning',
    timeout: 'error',
    cancelled: 'default',
    skipped: 'default',
    '': 'default',
  }
  return map[props.status] || 'default'
})

const label = computed(() => {
  const map: Record<string, string> = {
    success: '成功',
    failed: '失败',
    pending: '等待中',
    running: '运行中',
    never: '未备份',
    retrying: '重试中',
    timeout: '超时',
    cancelled: '已取消',
    skipped: '已跳过',
    '': '未备份',
  }
  return map[props.status] || props.status
})
</script>
