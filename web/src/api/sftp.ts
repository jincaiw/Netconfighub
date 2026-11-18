import { get, post, put, del } from './index'
import type {
  SftpAccount,
  SftpAccountCredentials,
  SftpStatus,
  SftpUpload,
  PaginatedResponse,
} from '@/types'

/** SFTP 服务运行状态（监听地址、主机密钥指纹、连接数） */
export function getSftpStatus() {
  return get<SftpStatus>('/sftp/status')
}

/** 查询设备的上传账号元数据（不含口令） */
export function getSftpAccount(deviceId: number) {
  return get<SftpAccount>(`/devices/${deviceId}/sftp-account`)
}

/** 创建账号，响应中的 password 只出现这一次 */
export function createSftpAccount(deviceId: number) {
  return post<SftpAccountCredentials>(`/devices/${deviceId}/sftp-account`)
}

/** 重置口令，响应中的 password 只出现这一次 */
export function resetSftpAccount(deviceId: number) {
  return post<SftpAccountCredentials>(`/devices/${deviceId}/sftp-account/reset`)
}

export function enableSftpAccount(deviceId: number) {
  return put<{ message: string }>(`/devices/${deviceId}/sftp-account/enable`)
}

export function disableSftpAccount(deviceId: number) {
  return put<{ message: string }>(`/devices/${deviceId}/sftp-account/disable`)
}

/** 撤销账号：之后该设备无法再通过 SFTP 认证 */
export function revokeSftpAccount(deviceId: number) {
  return del<{ message: string }>(`/devices/${deviceId}/sftp-account`)
}

/** 设备的 SFTP 上传记录（接收与归档两阶段状态） */
export function getSftpUploads(deviceId: number, params?: { page?: number; page_size?: number }) {
  return get<PaginatedResponse<SftpUpload>>(`/devices/${deviceId}/sftp-uploads`, params)
}

export function getSftpUpload(id: number) {
  return get<SftpUpload>(`/sftp-uploads/${id}`)
}

/** 重新排队一次失败的归档 */
export function retrySftpUpload(id: number) {
  return post<{ message: string }>(`/sftp-uploads/${id}/retry`)
}
