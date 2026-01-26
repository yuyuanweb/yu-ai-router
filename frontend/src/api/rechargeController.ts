// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 获取我的充值记录 GET /recharge/list/my */
export async function getMyRechargeRecords(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getMyRechargeRecordsParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponsePageRechargeRecord>('/recharge/list/my', {
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

/** Stripe充值取消回调 GET /recharge/stripe/cancel */
export async function stripeCancel(options?: { [key: string]: any }) {
  return request<API.BaseResponseString>('/recharge/stripe/cancel', {
    method: 'GET',
    ...(options || {}),
  })
}

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

/** Stripe充值成功回调 GET /recharge/stripe/success */
export async function stripeSuccess(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.stripeSuccessParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponseString>('/recharge/stripe/success', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}
