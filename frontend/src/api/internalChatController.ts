// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 内部聊天接口 POST /internal/chat/completions */
export async function chatCompletions1(body: API.ChatRequest, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/internal/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
