// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 获取我的余额信息 GET /balance/my */
export async function getMyBalance(options?: { [key: string]: any }) {
  return request<API.BaseResponseBalanceVO>('/balance/my', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 获取我的消费账单 GET /balance/billing/my */
export async function getMyBillingRecords(
  params: {
    pageNum?: number
    pageSize?: number
  },
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageBillingRecord>('/balance/billing/my', {
    method: 'GET',
    params: {
      pageNum: 1,
      pageSize: 10,
      ...params,
    },
    ...(options || {}),
  })
}
