// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 禁用插件 POST /plugin/disable */
export async function disablePlugin(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.disablePluginParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/plugin/disable', {
    method: 'POST',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 启用插件 POST /plugin/enable */
export async function enablePlugin(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.enablePluginParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/plugin/enable', {
    method: 'POST',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 执行插件 POST /plugin/execute */
export async function executePlugin(
  body: API.PluginExecuteRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePluginExecuteVO>('/plugin/execute', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取插件详情 GET /plugin/get */
export async function getPlugin(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPluginParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePluginConfigVO>('/plugin/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 获取所有插件列表 GET /plugin/list */
export async function listPlugins(options?: { [key: string]: any }) {
  return request<API.BaseResponseListPluginConfigVO>('/plugin/list', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 获取启用的插件列表 GET /plugin/list/enabled */
export async function listEnabledPlugins(options?: { [key: string]: any }) {
  return request<API.BaseResponseListPluginConfigVO>('/plugin/list/enabled', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 重新加载插件 POST /plugin/reload */
export async function reloadPlugin(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.reloadPluginParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/plugin/reload', {
    method: 'POST',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 重新加载所有插件 POST /plugin/reload/all */
export async function reloadAllPlugins(options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/plugin/reload/all', {
    method: 'POST',
    ...(options || {}),
  })
}

/** 更新插件配置 POST /plugin/update */
export async function updatePlugin(
  body: API.PluginUpdateRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseBoolean>('/plugin/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
