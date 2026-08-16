import request from '@/utils/request'

// 类型定义（字段名与后端 Go 代码返回的 JSON key 一致）
export interface CartItem {
  id: number
  goods_id: number
  goods_name: string
  goods_price: number
  goods_image: string
  nums: number
  checked: boolean
}

export interface AddCartParams {
  goods: number
  nums: number
}

export interface UpdateCartParams {
  nums?: number
  checked?: boolean
}

// 购物车列表
export const getCartList = () => {
  return request.get<{ data: CartItem[]; total: number }>('/o/v1/shoppingcart')
}

// 添加商品到购物车
export const addToCart = (data: AddCartParams) => {
  return request.post<{ id: number }>('/o/v1/shoppingcart', data)
}

// 修改购物车记录（后端将 :id 当作 goodsId 使用）
export const updateCartItem = (goodsId: number, data: UpdateCartParams) => {
  return request.patch(`/o/v1/shoppingcart/${goodsId}`, data)
}

// 删除购物车记录（后端将 :id 当作 goodsId 使用）
export const deleteCartItem = (goodsId: number) => {
  return request.delete(`/o/v1/shoppingcart/${goodsId}`)
}