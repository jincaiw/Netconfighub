<template>
  <t-tag :theme="tagTheme" variant="light" size="small">{{ label }}</t-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type TagTheme = 'default' | 'primary' | 'success' | 'warning' | 'danger'

const props = defineProps<{
  status:
    | 'success'
    | 'failed'
    | 'pending'
    | 'running'
    | 'never'
    | 'retrying'
    | 'timeout'
    | 'cancelled'
    | 'skipped'
    | ''
}>()

const tagTheme = computed<TagTheme>(() => {
  const map: Record<string, TagTheme> = {
    success: 'success',
    failed: 'danger',
    pending: 'primary',
    running: 'warning',
    never: 'default',
    retrying: 'warning',
    timeout: 'danger',
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
