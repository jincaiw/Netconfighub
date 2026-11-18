<template>
  <div class="baseline-manage">
    <t-card title="基线列表">
      <template #actions>
        <t-space>
          <t-button :loading="exportLoading" @click="handleExportDeviations">导出报告</t-button>
          <t-button theme="primary" @click="openCreateModal">创建基线</t-button>
        </t-space>
      </template>
      <t-table
        :columns="baselineColumns"
        :data="baselines"
        :loading="loading"
        :pagination="tablePagination"
        @page-change="handleBaselinePageChange"
      />
    </t-card>

    <t-dialog
      v-model:visible="showCreateModal"
      :header="editingBaseline ? '编辑基线' : '创建基线'"
      :width="560"
    >
      <t-form :data="baselineForm" label-align="left" label-width="90px">
        <t-form-item v-if="!editingBaseline" label="适用范围">
          <t-radio-group v-model="baselineForm.scope">
            <t-radio value="device">设备</t-radio>
            <t-radio value="group">分组</t-radio>
          </t-radio-group>
        </t-form-item>
        <t-form-item v-if="!editingBaseline && baselineForm.scope === 'device'" label="选择设备">
          <t-select
            v-model="baselineForm.device_id"
            :options="deviceOptions"
            placeholder="请选择设备"
            filterable
          />
        </t-form-item>
        <t-form-item v-if="!editingBaseline && baselineForm.scope === 'group'" label="选择分组">
          <t-select
            v-model="baselineForm.group_id"
            :options="groupSelectOptions"
            placeholder="请选择分组"
          />
        </t-form-item>
        <t-form-item label="基线内容">
          <t-textarea
            v-model="baselineForm.content"
            :autosize="{ minRows: 8 }"
            placeholder="请输入基线配置内容"
          />
        </t-form-item>
      </t-form>
      <template #footer>
        <t-button variant="outline" @click="closeCreateModal">取消</t-button>
        <t-button theme="primary" :loading="submitting" @click="handleSubmitBaseline">确定</t-button>
      </template>
    </t-dialog>

    <t-dialog v-model:visible="showDeviationModal" header="偏差详情" :width="640">
      <t-empty v-if="currentDeviations.length === 0" description="暂无偏差" />
      <div v-else>
        <t-card
          v-for="deviation in currentDeviations"
          :key="deviation.id"
          size="small"
          style="margin-bottom: 8px"
        >
          <pre class="diff-content">{{ deviation.diff_content }}</pre>
          <template #footer>
            <span class="diff-time">检测时间: {{ deviation.detected_at }}</span>
          </template>
        </t-card>
      </div>
      <template #footer>
        <t-button variant="outline" @click="showDeviationModal = false">关闭</t-button>
      </template>
    </t-dialog>

    <t-dialog v-model:visible="showExportModal" header="导出偏差报告" :width="480">
      <t-form label-align="left" label-width="90px">
        <t-form-item label="设备筛选">
          <t-select
            v-model="exportDeviceId"
            :options="deviceOptions"
            placeholder="全部设备"
            clearable
            filterable
          />
        </t-form-item>
        <t-form-item label="基线筛选">
          <t-select
            v-model="exportBaselineId"
            :options="baselineOptions"
            placeholder="全部基线"
            clearable
          />
        </t-form-item>
      </t-form>
      <template #footer>
        <t-button variant="outline" @click="showExportModal = false">取消</t-button>
        <t-button theme="primary" :loading="exportLoading" @click="doExportDeviations">导出</t-button>
      </template>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Button, Space, DialogPlugin, MessagePlugin, type PrimaryTableCol } from 'tdesign-vue-next'
import { getDevices, getGroups } from '@/api/device'
import {
  getBaselines,
  createBaseline,
  updateBaseline,
  deleteBaseline,
  getDeviations,
  exportDeviations,
} from '@/api/baseline'
import type { Baseline, Deviation, Device, Group } from '@/types'

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
  }),
)

const baselinePagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
})

const tablePagination = computed(() => ({
  current: baselinePagination.page,
  pageSize: baselinePagination.pageSize,
  total: baselinePagination.itemCount,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50],
}))

const baselineForm = reactive({
  scope: 'device' as 'device' | 'group',
  device_id: undefined as number | undefined,
  group_id: undefined as number | undefined,
  content: '',
})

const baselineColumns: PrimaryTableCol<Baseline>[] = [
  {
    title: '设备/分组名称',
    colKey: 'name',
    width: 160,
    cell: (_h, { row }) => {
      if (row.scope === 'device') {
        return deviceMap.value[row.device_id!] || `设备#${row.device_id}`
      }
      return groupMap.value[row.group_id!] || `分组#${row.group_id}`
    },
  },
  {
    title: '适用范围',
    colKey: 'scope',
    width: 100,
    cell: (_h, { row }) => (row.scope === 'device' ? '设备' : '分组'),
  },
  { title: '创建时间', colKey: 'created_at', width: 170 },
  { title: '更新时间', colKey: 'updated_at', width: 170 },
  {
    title: '操作',
    colKey: 'actions',
    width: 220,
    cell: (h, { row }) =>
      h(
        Space,
        { size: 'small' },
        {
          default: () => [
            h(
              Button,
              { size: 'small', onClick: () => handleViewDeviation(row.id) },
              { default: () => '偏差' },
            ),
            h(
              Button,
              { size: 'small', theme: 'primary', onClick: () => openEditModal(row) },
              { default: () => '编辑' },
            ),
            h(
              Button,
              { size: 'small', theme: 'danger', onClick: () => handleDelete(row) },
              { default: () => '删除' },
            ),
          ],
        },
      ),
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
    MessagePlugin.error('获取基线列表失败')
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
    MessagePlugin.error('获取设备/分组选项失败')
  }
}

function handleBaselinePageChange(pageInfo: { current: number; pageSize: number }) {
  if (pageInfo.pageSize && pageInfo.pageSize !== baselinePagination.pageSize) {
    baselinePagination.pageSize = pageInfo.pageSize
    baselinePagination.page = 1
  } else {
    baselinePagination.page = pageInfo.current
  }
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
    MessagePlugin.warning('请输入基线内容')
    return
  }
  if (!editingBaseline.value) {
    if (baselineForm.scope === 'device' && !baselineForm.device_id) {
      MessagePlugin.warning('请选择设备')
      return
    }
    if (baselineForm.scope === 'group' && !baselineForm.group_id) {
      MessagePlugin.warning('请选择分组')
      return
    }
  }

  submitting.value = true
  try {
    if (editingBaseline.value) {
      await updateBaseline(editingBaseline.value.id, { content: baselineForm.content })
      MessagePlugin.success('基线更新成功')
    } else {
      await createBaseline(baselineForm)
      MessagePlugin.success('基线创建成功')
    }
    showCreateModal.value = false
    editingBaseline.value = null
    fetchBaselines()
  } catch {
    MessagePlugin.error(editingBaseline.value ? '基线更新失败' : '基线创建失败')
  } finally {
    submitting.value = false
  }
}

function handleDelete(baseline: Baseline) {
  DialogPlugin.confirm({
    header: '确认删除',
    body: '确定要删除该基线吗？',
    theme: 'warning',
    confirmBtn: '删除',
    cancelBtn: '取消',
    onConfirm: async () => {
      try {
        await deleteBaseline(baseline.id)
        MessagePlugin.success('基线删除成功')
        fetchBaselines()
      } catch {
        MessagePlugin.error('基线删除失败')
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
    MessagePlugin.error('获取偏差信息失败')
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
    MessagePlugin.success('导出成功')
    showExportModal.value = false
  } catch {
    MessagePlugin.error('导出偏差报告失败')
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

.diff-time {
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}
</style>
