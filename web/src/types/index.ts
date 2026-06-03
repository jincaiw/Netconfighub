export interface Device {
  id: number
  name: string
  ip: string
  port: number
  protocol: 'ssh' | 'telnet'
  vendor: 'cisco' | 'h3c' | 'huawei' | 'ruijie'
  model: string
  username: string
  group_id?: number
  group_name?: string
  enabled: boolean
  last_backup_status: 'success' | 'failed' | 'pending' | 'running' | 'never' | 'retrying' | 'timeout' | 'cancelled' | 'skipped'
  last_backup_at?: string
  last_failure_reason?: string
  retry_count: number
  backup_interval?: string
  created_at: string
  updated_at: string
}

export interface DeviceCreateRequest {
  name: string
  ip: string
  port: number
  protocol: 'ssh' | 'telnet'
  vendor: 'cisco' | 'h3c' | 'huawei' | 'ruijie'
  model: string
  username: string
  password?: string
  ssh_key?: string
  enable_password?: string
  backup_interval?: string
  group_id?: number
  enabled?: boolean
}

export interface Group {
  id: number
  name: string
  description?: string
  device_count?: number
  created_at: string
  updated_at: string
}

export interface GroupCreateRequest {
  name: string
  description?: string
}

export interface BackupTask {
  id: number
  device_id: number
  device_name?: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'retrying' | 'timeout' | 'cancelled' | 'skipped'
  started_at?: string
  finished_at?: string
  failure_reason?: string
  retry_count: number
  error_type?: string
  config_hash?: string
  created_at: string
}

export interface Baseline {
  id: number
  device_id?: number
  group_id?: number
  scope: 'device' | 'group'
  content: string
  created_at: string
  updated_at: string
}

export interface Deviation {
  id: number
  backup_task_id: number
  baseline_id: number
  diff_content: string
  detected_at: string
}

export interface ConfigVersion {
  hash: string
  message: string
  author: string
  date: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string | null
}

export interface RefreshResponse {
  token: string
  expires_at: string
}

export interface ApiToken {
  id: number
  name: string
  token?: string
  created_at: string
  expires_at: string | null
}

export interface TokenCreateRequest {
  name: string
  expires_at: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface Alert {
  id: number
  type: 'backup_failed' | 'baseline_deviation'
  device_id?: number
  device_name: string
  title: string
  message: string
  severity: 'error' | 'warning' | 'info'
  is_read: boolean
  created_at: string
}

export interface DiffResult {
  from_hash: string
  to_hash: string
  diff_content: string
  added_lines: number
  removed_lines: number
}

export interface ImportResult {
  total: number
  success: number
  failed: number
  errors?: ImportError[]
}

export interface ImportError {
  row: number
  name: string
  reason: string
}

export interface AuditLog {
  id: number
  user_id: number
  username: string
  action: string
  target_type: string
  target_id: number
  detail: string
  client_ip: string
  created_at: string
}
