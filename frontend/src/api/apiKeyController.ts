// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 创建 API Key POST /api/key/create */
export async function createApiKey(
  body: API.ApiKeyCreateRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseApiKeyVO>('/api/key/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取我的 API Key 列表 GET /api/key/list/my */
export async function listMyApiKeys(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.listMyApiKeysParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageApiKeyVO>('/api/key/list/my', {
    method: 'GET',
    params: {
      // pageNum has a default value: 1
      pageNum: '1',
      // pageSize has a default value: 10
      pageSize: '10',
      ...params,
    },
    ...(options || {}),
  })
}

/** 撤销 API Key POST /api/key/revoke */
export async function revokeApiKey(body: API.DeleteRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/api/key/revoke', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
