import request from '@/utils/request'

export interface GoodsItem {
  id: number
  name: string
  goods_sn: string
  market_price: number
  shop_price: number
  front_image: string
  goods_brief: string
  images: string[]
  desc_images: string[]
  is_hot: boolean
  is_new: boolean
  on_sale: boolean
  ship_free: boolean
  stocks: number
  brand: {
    id: number
    name: string
    logo: string
  }
  ctegory: {
    id: number
    name: string
  }
}

export interface GoodsListParams {
  p?: number
  pnum?: number
  q?: string
  c?: number
  b?: number
  pmin?: number
  pmax?: number
  ishot?: number
  isnew?: number
  istab?: number
  onsale?: number
}

export interface GoodsListResponse {
  data: GoodsItem[]
  total: number
}

export const getGoodsList = (params?: GoodsListParams) => {
  return request.get<GoodsListResponse>('/g/v1/goods', params)
}

export const getGoodsDetail = (id: number) => {
  return request.get(`/g/v1/goods/${id}`)
}

export const createGoods = (data: any) => {
  return request.post('/g/v1/goods', data)
}

export const updateGoods = (id: number, data: any) => {
  return request.put(`/g/v1/goods/${id}`, data)
}

export const updateGoodsStatus = (id: number, data: { sale: boolean; new: boolean; hot: boolean }) => {
  return request.patch(`/g/v1/goods/${id}`, data)
}

export const deleteGoods = (id: number) => {
  return request.delete(`/g/v1/goods/${id}`)
}

export const getGoodsStock = (id: number) => {
  return request.get<{ goodsId: number; num: number }>(`/g/v1/goods/${id}/stocks`)
}