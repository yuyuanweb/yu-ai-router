// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 获取我的消费账单 GET /balance/billing/my */
export async function getMyBillingRecords(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getMyBillingRecordsParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageBillingRecord>('/balance/billing/my', {
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

/** 获取我的余额信息 GET /balance/my */
export async function getMyBalance(options?: { [key: string]: any }) {
  return request<API.BaseResponseBalanceVO>('/balance/my', {
    method: 'GET',
    ...(options || {}),
  })
}
