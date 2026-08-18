import request from '@/utils/request'

export interface OrderItem {
  id: number
  order_sn: string
  user: number
  name: string
  mobile: string
  address: string
  post: string
  pay_type: string
  status: string
  total: number
  add_time: string
  goods?: Array<{
    id: number
    name: string
    image: string
    price: number
    nums: number
  }>
}

export interface OrderListResponse {
  data: OrderItem[]
  total: number
}

export const getOrderList = (params?: { p?: number; pnum?: number }) => {
  return request.get<OrderListResponse>('/o/v1/order', params)
}

export const getOrderDetail = (id: number) => {
  return request.get<OrderItem>(`/o/v1/order/${id}`)
}
