import request from '@/utils/request'

// 类型定义
export interface CategoryItem {
  id: number
  name: string
  level: number
  isTab: boolean
  parent: number | null
  sub_category: CategoryItem[]
}

export interface BrandItem {
  id: number
  name: string
  logo: string
}

export interface BannerItem {
  id: number
  image: string
  url: string
  index: number
}

export interface BrandListResponse {
  data: BrandItem[]
  total: number
}

// 商品分类列表（树形结构）
export const getCategoryList = () => {
  return request.get<CategoryItem[]>('/g/v1/categorys')
}

// 分类详情
export const getCategoryDetail = (id: number) => {
  return request.get<CategoryItem>(`/g/v1/categorys/${id}`)
}

// 品牌列表
export const getBrandList = (params?: { pn?: number; psize?: number }) => {
  return request.get<BrandListResponse>('/g/v1/brands', params)
}

// 轮播图列表
export const getBannerList = () => {
  return request.get<BannerItem[]>('/g/v1/banners')
}

// 通过分类获取品牌
export const getBrandsByCategory = (categoryId: number) => {
  return request.get<BrandItem[]>(`/g/v1/categorybrands/${categoryId}`)
}