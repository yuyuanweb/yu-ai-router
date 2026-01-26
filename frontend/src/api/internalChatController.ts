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

/** 内部聊天接口（支持文件上传） POST /internal/chat/completions/upload */
export async function chatCompletionsWithFile(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.chatCompletionsWithFileParams,
  body: {},
  file?: File,
  options?: { [key: string]: any }
) {
  const formData = new FormData()

  if (file) {
    formData.append('file', file)
  }

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele]

    if (item !== undefined && item !== null) {
      if (typeof item === 'object' && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ''))
        } else {
          formData.append(ele, JSON.stringify(item))
        }
      } else {
        formData.append(ele, item)
      }
    }
  })

  return request<Record<string, any>>('/internal/chat/completions/upload', {
    method: 'POST',
    params: {
      ...params,
    },
    data: formData,
    requestType: 'form',
    ...(options || {}),
  })
}
