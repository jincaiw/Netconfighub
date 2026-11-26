<template>
  <t-space break-line>
    <t-input
      v-model="keyword"
      placeholder="搜索设备名称/IP"
      clearable
      style="width: 180px"
      @clear="$emit('search')"
      @enter="$emit('search')"
    >
      <template #prefixIcon><SearchIcon /></template>
    </t-input>
    <t-select
      v-model="groupId"
      :options="groupOptions"
      placeholder="按分组筛选"
      clearable
      style="width: 160px"
      @change="$emit('group-change')"
    />
    <t-button variant="outline" :loading="exportLoading" @click="$emit('export')">导出</t-button>
    <t-button variant="outline" @click="$emit('import')">导入</t-button>
    <t-button theme="primary" @click="$emit('create')">
      <template #icon><AddIcon /></template>
      添加设备
    </t-button>
  </t-space>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { SearchIcon, AddIcon } from 'tdesign-icons-vue-next'

const props = defineProps<{
  searchKeyword: string
  selectedGroup: number | null
  groupOptions: Array<{ label: string; value: number }>
  exportLoading: boolean
}>()

const emit = defineEmits<{
  (e: 'update:searchKeyword', value: string): void
  (e: 'update:selectedGroup', value: number | null): void
  (e: 'search'): void
  (e: 'group-change'): void
  (e: 'export'): void
  (e: 'import'): void
  (e: 'create'): void
}>()

const keyword = computed({
  get: () => props.searchKeyword,
  set: (value: string) => emit('update:searchKeyword', value),
})

const groupId = computed({
  get: () => props.selectedGroup,
  set: (value: number | null) => emit('update:selectedGroup', value),
})
</script>
