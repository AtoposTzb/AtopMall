import request from '@/utils/request'

export interface CategoryItem {
  id: number
  name: string
  level: number
  isTab: boolean
  parent: number | null
  sub_category?: CategoryItem[]
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

export interface CategoryBrandItem {
  id: number
  category: {
    id: number
    name: string
  }
  brand: {
    id: number
    name: string
    logo: string
  }
}

export interface CategoryBrandListResponse {
  data: CategoryBrandItem[]
  total: number
}

export const getCategoryList = () => {
  return request.get<CategoryItem[]>('/g/v1/categorys', { _t: Date.now() })
}

export const createCategory = (data: { name: string; parent: number; level: number; is_tab: boolean }) => {
  return request.post('/g/v1/categorys', data)
}

export const updateCategory = (id: number, data: { name: string; is_tab: boolean }) => {
  return request.put(`/g/v1/categorys/${id}`, data)
}

export const deleteCategory = (id: number) => {
  return request.delete(`/g/v1/categorys/${id}`)
}

export const getBrandList = (pn = 0, psize = 10) => {
  return request.get<BrandListResponse>('/g/v1/brands', { pn, psize })
}

export const createBrand = (data: { name: string; logo: string }) => {
  return request.post('/g/v1/brands', data)
}

export const updateBrand = (id: number, data: { name: string; logo: string }) => {
  return request.put(`/g/v1/brands/${id}`, data)
}

export const deleteBrand = (id: number) => {
  return request.delete(`/g/v1/brands/${id}`)
}

export const getBannerList = () => {
  return request.get<BannerItem[]>('/g/v1/banners')
}

export const createBanner = (data: { image: string; url: string; index: number }) => {
  return request.post('/g/v1/banners', data)
}

export const updateBanner = (id: number, data: { image: string; url: string; index: number }) => {
  return request.put(`/g/v1/banners/${id}`, data)
}

export const deleteBanner = (id: number) => {
  return request.delete(`/g/v1/banners/${id}`)
}

export const getCategoryBrandList = () => {
  return request.get<CategoryBrandListResponse>('/g/v1/categorybrands')
}

export const createCategoryBrand = (data: { category_id: number; brand_id: number }) => {
  return request.post('/g/v1/categorybrands', data)
}

export const updateCategoryBrand = (id: number, data: { category_id: number; brand_id: number }) => {
  return request.put(`/g/v1/categorybrands/${id}`, data)
}

export const deleteCategoryBrand = (id: number) => {
  return request.delete(`/g/v1/categorybrands/${id}`)
}

export const getBrandsByCategory = (categoryId: number) => {
  return request.get<BrandItem[]>(`/g/v1/categorybrands/${categoryId}`)
}