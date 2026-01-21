// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 此处后端没有提供注释 POST /provider/add */
export async function addProvider(body: API.ProviderAddRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseLong>('/provider/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 POST /provider/delete */
export async function deleteProvider(body: API.DeleteRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/provider/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /provider/get/vo */
export async function getProviderVoById(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getProviderVOByIdParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseProviderVO>('/provider/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /provider/list/healthy */
export async function listHealthyProviders(options?: { [key: string]: any }) {
  return request<API.BaseResponseListProviderVO>('/provider/list/healthy', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 POST /provider/list/page/vo */
export async function listProviderVoByPage(
  body: API.ProviderQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageProviderVO>('/provider/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /provider/list/vo */
export async function listProviderVo(options?: { [key: string]: any }) {
  return request<API.BaseResponseListProviderVO>('/provider/list/vo', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 POST /provider/update */
export async function updateProvider(
  body: API.ProviderUpdateRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/provider/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
