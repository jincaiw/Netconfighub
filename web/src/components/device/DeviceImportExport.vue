<template>
  <t-dialog v-model:visible="importVisible" header="导入设备" :width="480">
    <t-upload
      v-model="uploadFiles"
      :auto-upload="false"
      accept=".csv"
      theme="file"
      @change="onFileChange"
    />
    <p class="import-hint">
      CSV 格式: name, ip, vendor, model, protocol, port, username, password, group_name
    </p>
    <template #footer>
      <t-button variant="outline" @click="importVisible = false">取消</t-button>
      <t-button theme="primary" :loading="importLoading" :disabled="!hasFile" @click="emit('import')">导入</t-button>
    </template>
  </t-dialog>

  <t-dialog v-model:visible="resultVisible" header="导入结果" :width="480">
    <div v-if="importResult">
      <t-descriptions bordered :column="1">
        <t-descriptions-item label="总数">{{ importResult.total }}</t-descriptions-item>
        <t-descriptions-item label="成功">
          <span class="result-success">{{ importResult.success }}</span>
        </t-descriptions-item>
        <t-descriptions-item label="失败">
          <span class="result-error">{{ importResult.failed }}</span>
        </t-descriptions-item>
      </t-descriptions>
      <div v-if="importResult.errors && importResult.errors.length > 0" style="margin-top: 12px">
        <span class="import-hint">失败详情:</span>
        <div v-for="(err, idx) in importResult.errors" :key="idx" style="margin-top: 4px">
          <span class="result-error">第 {{ err.row }} 行 ({{ err.name }}): {{ err.reason }}</span>
        </div>
      </div>
    </div>
    <template #footer>
      <t-button theme="primary" @click="resultVisible = false">确定</t-button>
    </template>
  </t-dialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { UploadFile } from 'tdesign-vue-next'
import type { ImportResult } from '@/types'

const props = defineProps<{
  visible: boolean
  resultVisible: boolean
  importResult: ImportResult | null
  importLoading: boolean
  hasFile: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'update:resultVisible', value: boolean): void
  (e: 'file-change', files: UploadFile[]): void
  (e: 'import'): void
}>()

const uploadFiles = ref<UploadFile[]>([])

const importVisible = computed({
  get: () => props.visible,
  set: (value: boolean) => emit('update:visible', value),
})
const resultVisible = computed({
  get: () => props.resultVisible,
  set: (value: boolean) => emit('update:resultVisible', value),
})

function onFileChange(files: UploadFile[]) {
  uploadFiles.value = files
  emit('file-change', files)
}
</script>

<style scoped>
.import-hint {
  margin-top: 12px;
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}

.result-success {
  color: var(--td-success-color-5);
}

.result-error {
  color: var(--td-error-color-5);
}
</style>
