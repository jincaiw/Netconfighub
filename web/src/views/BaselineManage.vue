<template>
  <div class="baseline-manage">
    <n-card title="基线列表">
      <template #header-extra>
        <n-space>
          <n-button @click="handleExportDeviations" :loading="exportLoading">导出报告</n-button>
          <n-button type="primary" @click="openCreateModal">创建基线</n-button>
        </n-space>
      </template>
      <n-data-table
        :columns="baselineColumns"
        :data="baselines"
        :loading="loading"
        :pagination="baselinePagination"
        remote
        @update:page="handleBaselinePageChange"
        @update:page-size="handleBaselinePageSizeChange"
      />
    </n-card>

    <n-modal v-model:show="showCreateModal" preset="dialog" :title="editingBaseline ? '编辑基线' : '创建基线'" style="width: 560px">
      <n-form :model="baselineForm" label-placement="left" label-width="90">
        <n-form-item v-if="!editingBaseline" label="适用范围">
          <n-radio-group v-model:value="baselineForm.scope">
            <n-radio value="device">设备</n-radio>
            <n-radio value="group">分组</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="!editingBaseline && baselineForm.scope === 'device'" label="选择设备">
          <n-select
            v-model:value="baselineForm.device_id"
            :options="deviceOptions"
            placeholder="请选择设备"
            filterable
          />
        </n-form-item>
        <n-form-item v-if="!editingBaseline && baselineForm.scope === 'group'" label="选择分组">
          <n-select
            v-model:value="baselineForm.group_id"
            :options="groupSelectOptions"
            placeholder="请选择分组"
          />
        </n-form-item>
        <n-form-item label="基线内容">
          <n-input
            v-model:value="baselineForm.content"
            type="textarea"
            placeholder="请输入基线配置内容"
            :rows="8"
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="closeCreateModal">取消</n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmitBaseline">确定</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showDeviationModal" preset="dialog" title="偏差详情" style="width: 640px">
      <n-empty v-if="currentDeviations.length === 0" description="暂无偏差" />
      <div v-else>
        <n-card
          v-for="deviation in currentDeviations"
          :key="deviation.id"
          size="small"
          style="margin-bottom: 8px"
        >
          <pre class="diff-content">{{ deviation.diff_content }}</pre>
          <template #footer>
            <n-text depth="3">检测时间: {{ deviation.detected_at }}</n-text>
          </template>
        </n-card>
      </div>
      <template #action>
        <n-button @click="showDeviationModal = false">关闭</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showExportModal" preset="dialog" title="导出偏差报告" style="width: 480px">
      <n-form label-placement="left" label-width="90">
        <n-form-item label="设备筛选">
          <n-select
            v-model:value="exportDeviceId"
            :options="deviceOptions"
            placeholder="全部设备"
            clearable
            filterable
          />
        </n-form-item>
        <n-form-item label="基线筛选">
          <n-select
            v-model:value="exportBaselineId"
            :options="baselineOptions"
            placeholder="全部基线"
            clearable
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showExportModal = false">取消</n-button>
        <n-button type="primary" :loading="exportLoading" @click="doExportDeviations">导出</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import {
  NCard,
  NDataTable,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NRadioGroup,
  NRadio,
  NButton,
  NEmpty,
  NText,
  NSpace,
  useMessage,
  useDialog,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getDevices, getGroups } from '@/api/device'
import { getBaselines, createBaseline, updateBaseline, deleteBaseline, getDeviations, exportDeviations } from '@/api/baseline'
import type { Baseline, Deviation, Device, Group } from '@/types'

const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const submitting = ref(false)
const exportLoading = ref(false)
const baselines = ref<Baseline[]>([])
const showCreateModal = ref(false)
const showDeviationModal = ref(false)
const showExportModal = ref(false)
const editingBaseline = ref<Baseline | null>(null)
const currentDeviations = ref<Deviation[]>([])
const exportDeviceId = ref<number | null>(null)
const exportBaselineId = ref<number | null>(null)

const deviceOptions = ref<Array<{ label: string; value: number }>>([])
const groupSelectOptions = ref<Array<{ label: string; value: number }>>([])
const deviceMap = ref<Record<number, string>>({})
const groupMap = ref<Record<number, string>>({})

const baselineOptions = computed(() =>
  baselines.value.map((b) => {
    let label = ''
    if (b.scope === 'device') {
      label = deviceMap.value[b.device_id!] || `设备#${b.device_id}`
    } else {
      label = groupMap.value[b.group_id!] || `分组#${b.group_id}`
    }
    return { label: `${label} (ID: ${b.id})`, value: b.id }
  })
)

const baselinePagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

const baselineForm = reactive({
  scope: 'device' as 'device' | 'group',
  device_id: undefined as number | undefined,
  group_id: undefined as number | undefined,
  content: '',
})

const baselineColumns: DataTableColumns<Baseline> = [
  {
    title: '设备/分组名称',
    key: 'name',
    width: 160,
    render: (row) => {
      if (row.scope === 'device') {
        return deviceMap.value[row.device_id!] || `设备#${row.device_id}`
      }
      return groupMap.value[row.group_id!] || `分组#${row.group_id}`
    },
  },
  {
    title: '适用范围',
    key: 'scope',
    width: 100,
    render: (row) => (row.scope === 'device' ? '设备' : '分组'),
  },
  { title: '创建时间', key: 'created_at', width: 170 },
  { title: '更新时间', key: 'updated_at', width: 170 },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleViewDeviation(row.id) },
            { default: () => '偏差' },
          ),
          h(
            NButton,
            { size: 'small', type: 'primary', onClick: () => openEditModal(row) },
            { default: () => '编辑' },
          ),
          h(
            NButton,
            { size: 'small', type: 'error', onClick: () => handleDelete(row) },
            { default: () => '删除' },
          ),
        ],
      }),
  },
]

async function fetchBaselines() {
  loading.value = true
  try {
    const res = await getBaselines({
      page: baselinePagination.page,
      page_size: baselinePagination.pageSize,
    })
    baselines.value = res.data.items
    baselinePagination.itemCount = res.data.total
  } catch {
    message.error('获取基线列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchOptions() {
  try {
    const [devRes, grpRes] = await Promise.all([
      getDevices({ page: 1, page_size: 100 }),
      getGroups({ page: 1, page_size: 100 }),
    ])
    const devItems = devRes.data.items
    const grpItems = grpRes.data.items
    deviceOptions.value = devItems.map((d: Device) => ({
      label: `${d.name} (${d.ip})`,
      value: d.id,
    }))
    groupSelectOptions.value = grpItems.map((g: Group) => ({ label: g.name, value: g.id }))
    deviceMap.value = {}
    for (const d of devItems) {
      deviceMap.value[d.id] = d.name
    }
    groupMap.value = {}
    for (const g of grpItems) {
      groupMap.value[g.id] = g.name
    }
  } catch {
    message.error('获取设备/分组选项失败')
  }
}

function handleBaselinePageChange(page: number) {
  baselinePagination.page = page
  fetchBaselines()
}

function handleBaselinePageSizeChange(pageSize: number) {
  baselinePagination.pageSize = pageSize
  baselinePagination.page = 1
  fetchBaselines()
}

function openCreateModal() {
  editingBaseline.value = null
  baselineForm.scope = 'device'
  baselineForm.device_id = undefined
  baselineForm.group_id = undefined
  baselineForm.content = ''
  showCreateModal.value = true
}

function openEditModal(baseline: Baseline) {
  editingBaseline.value = baseline
  baselineForm.scope = baseline.scope
  baselineForm.device_id = baseline.device_id
  baselineForm.group_id = baseline.group_id
  baselineForm.content = baseline.content
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  editingBaseline.value = null
}

async function handleSubmitBaseline() {
  if (!baselineForm.content.trim()) {
    message.warning('请输入基线内容')
    return
  }
  if (!editingBaseline.value) {
    if (baselineForm.scope === 'device' && !baselineForm.device_id) {
      message.warning('请选择设备')
      return
    }
    if (baselineForm.scope === 'group' && !baselineForm.group_id) {
      message.warning('请选择分组')
      return
    }
  }

  submitting.value = true
  try {
    if (editingBaseline.value) {
      await updateBaseline(editingBaseline.value.id, { content: baselineForm.content })
      message.success('基线更新成功')
    } else {
      await createBaseline(baselineForm)
      message.success('基线创建成功')
    }
    showCreateModal.value = false
    editingBaseline.value = null
    fetchBaselines()
  } catch {
    message.error(editingBaseline.value ? '基线更新失败' : '基线创建失败')
  } finally {
    submitting.value = false
  }
}

function handleDelete(baseline: Baseline) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除该基线吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteBaseline(baseline.id)
        message.success('基线删除成功')
        fetchBaselines()
      } catch {
        message.error('基线删除失败')
      }
    },
  })
}

async function handleViewDeviation(baselineId: number) {
  try {
    const res = await getDeviations({ baseline_id: baselineId, page: 1, page_size: 20 })
    currentDeviations.value = res.data.items
    showDeviationModal.value = true
  } catch {
    message.error('获取偏差信息失败')
  }
}

function handleExportDeviations() {
  exportDeviceId.value = null
  exportBaselineId.value = null
  showExportModal.value = true
}

async function doExportDeviations() {
  exportLoading.value = true
  try {
    const params: Record<string, any> = {}
    if (exportDeviceId.value) {
      params.device_id = exportDeviceId.value
    }
    if (exportBaselineId.value) {
      params.baseline_id = exportBaselineId.value
    }
    const blob = await exportDeviations(params)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'deviations.csv'
    a.click()
    window.URL.revokeObjectURL(url)
    message.success('导出成功')
    showExportModal.value = false
  } catch {
    message.error('导出偏差报告失败')
  } finally {
    exportLoading.value = false
  }
}

onMounted(() => {
  fetchBaselines()
  fetchOptions()
})
</script>

<style scoped>
.baseline-manage {
  width: 100%;
}

.diff-content {
  margin: 0;
  font-family: 'SF Mono', 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
