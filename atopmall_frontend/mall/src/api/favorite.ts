import request from '@/utils/request'

// 类型定义
export interface FavItem {
  id: number
  goods_id?: number
  name: string
  shop_price: number
  front_image?: string
  goods_brief?: string
}

export interface FavListResponse {
  data: FavItem[]
  total: number
}

// 收藏列表
export const getFavList = () => {
  return request.get<FavListResponse>('/op/v1/userfavs/')
}

// 添加收藏
export const addFav = (goodsId: number) => {
  return request.post('/op/v1/userfavs', { goods: goodsId })
}

// 取消收藏
export const deleteFav = (goodsId: number) => {
  return request.delete(`/op/v1/userfavs/${goodsId}`)
}

export const getFavStatus = (goodsId: number) => {
  return request.get(`/op/v1/userfavs/${goodsId}`, undefined, { skipErrorToast: true } as any)
}
