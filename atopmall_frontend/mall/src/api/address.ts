import request from '@/utils/request'

// 类型定义
export interface AddressItem {
  id: number
  province: string
  city: string
  district: string
  address: string
  signer_name: string
  signer_mobile: string
  user_id: number
}

export interface AddressListResponse {
  data: AddressItem[]
  total: number
}

export interface CreateAddressParams {
  province: string
  city: string
  district: string
  address: string
  signer_name: string
  signer_mobile: string
}

// 地址列表
export const getAddressList = () => {
  return request.get<AddressListResponse>('/op/v1/address')
}

// 添加地址
export const createAddress = (data: CreateAddressParams) => {
  return request.post<{ id: number }>('/op/v1/address', data)
}

// 更新地址
export const updateAddress = (id: number, data: CreateAddressParams) => {
  return request.put(`/op/v1/address/${id}`, data)
}

// 删除地址
export const deleteAddress = (id: number) => {
  return request.delete(`/op/v1/address/${id}`)
}
