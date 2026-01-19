// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** Chat Completions POST /v1/chat/completions */
export async function chatCompletions(body: API.ChatRequest, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
