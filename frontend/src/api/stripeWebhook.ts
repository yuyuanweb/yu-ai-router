// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** Stripe Webhook回调 POST /webhook/stripe */
export async function handleStripeWebhook(body: string, options?: { [key: string]: any }) {
  return request<string>('/webhook/stripe', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
