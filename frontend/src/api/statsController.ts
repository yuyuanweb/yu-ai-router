// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 获取我的请求日志 GET /stats/my/logs */
export async function getMyLogs(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getMyLogsParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseListRequestLog>('/stats/my/logs', {
    method: 'GET',
    params: {
      // limit has a default value: 100
      limit: '100',
      ...params,
    },
    ...(options || {}),
  })
}

/** 获取我的 Token 消耗统计 GET /stats/my/tokens */
export async function getMyTokenStats(options?: { [key: string]: any }) {
  return request<API.BaseResponseTokenStatsVO>('/stats/my/tokens', {
    method: 'GET',
    ...(options || {}),
  })
}
