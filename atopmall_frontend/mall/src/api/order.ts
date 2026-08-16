import request from '@/utils/request'

// 类型定义
export interface OrderGoods {
  id: number
  name: string
  image: string
  price: number
  nums: number
}

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
  goods?: OrderGoods[]
  alipay_url?: string
}

export interface OrderListResponse {
  data: OrderItem[]
  total: number
}

export interface CreateOrderParams {
  address: string
  name: string
  mobile: string
  post: string
}

// 订单列表
export const getOrderList = (params?: { p?: number; pnum?: number }) => {
  return request.get<OrderListResponse>('/o/v1/order', params)
}

// 订单详情
export const getOrderDetail = (id: number) => {
  return request.get<OrderItem>(`/o/v1/order/${id}`)
}

// 创建订单
export const createOrder = (data: CreateOrderParams) => {
  return request.post<{ id: number }>('/o/v1/order', data)
}
