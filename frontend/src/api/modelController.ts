// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 此处后端没有提供注释 POST /model/add */
export async function addModel(body: API.ModelAddRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseLong>('/model/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 POST /model/delete */
export async function deleteModel(body: API.DeleteRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/model/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /model/get/vo */
export async function getModelVoById(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getModelVOByIdParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseModelVO>('/model/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /model/list/active */
export async function listActiveModels(options?: { [key: string]: any }) {
  return request<API.BaseResponseListModelVO>('/model/list/active', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /model/list/active/provider/${param0} */
export async function listActiveModelsByProvider(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.listActiveModelsByProviderParams,
  options?: { [key: string]: any }
) {
  const { providerId: param0, ...queryParams } = params
  return request<API.BaseResponseListModelVO>(`/model/list/active/provider/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /model/list/active/type/${param0} */
export async function listActiveModelsByType(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.listActiveModelsByTypeParams,
  options?: { [key: string]: any }
) {
  const { modelType: param0, ...queryParams } = params
  return request<API.BaseResponseListModelVO>(`/model/list/active/type/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 POST /model/list/page/vo */
export async function listModelVoByPage(
  body: API.ModelQueryRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageModelVO>('/model/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 GET /model/list/vo */
export async function listModelVo(options?: { [key: string]: any }) {
  return request<API.BaseResponseListModelVO>('/model/list/vo', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 此处后端没有提供注释 POST /model/update */
export async function updateModel(body: API.ModelUpdateRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/model/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
