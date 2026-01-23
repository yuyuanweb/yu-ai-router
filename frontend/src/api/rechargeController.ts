// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 创建Stripe充值订单 POST /recharge/stripe/create */
export async function createStripeRecharge(
  body: API.CreateRechargeRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseCreateRechargeResponse>('/recharge/stripe/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取我的充值记录 GET /recharge/list/my */
export async function getMyRechargeRecords(
  params: {
    pageNum?: number
    pageSize?: number
  },
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageRechargeRecord>('/recharge/list/my', {
    method: 'GET',
    params: {
      pageNum: 1,
      pageSize: 10,
      ...params,
    },
    ...(options || {}),
  })
}
