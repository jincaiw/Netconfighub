import { get, post, del } from './index'
import type { ApiToken, TokenCreateRequest } from '@/types'

export function getTokens() {
  return get<ApiToken[]>('/tokens')
}

export function createToken(data: TokenCreateRequest) {
  return post<ApiToken>('/tokens', data as unknown as Record<string, any>)
}

export function deleteToken(id: number) {
  return del<null>(`/tokens/${id}`)
}
