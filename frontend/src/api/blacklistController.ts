// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 添加 IP 到黑名单 POST /admin/blacklist/add */
export async function addToBlacklist(body: API.BlacklistRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/admin/blacklist/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 检查 IP 是否在黑名单中 GET /admin/blacklist/check */
export async function checkBlacklist(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.checkBlacklistParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/admin/blacklist/check', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 获取黑名单数量 GET /admin/blacklist/count */
export async function getBlacklistCount(options?: { [key: string]: any }) {
  return request<API.BaseResponseLong>('/admin/blacklist/count', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 获取黑名单列表 GET /admin/blacklist/list */
export async function getBlacklist(options?: { [key: string]: any }) {
  return request<API.BaseResponseSetString>('/admin/blacklist/list', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 从黑名单移除 IP POST /admin/blacklist/remove */
export async function removeFromBlacklist(
  body: API.BlacklistRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/admin/blacklist/remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
