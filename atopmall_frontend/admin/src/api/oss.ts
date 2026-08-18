import request from '@/utils/request'

export interface OssTokenResp {
  upload_url: string
  url: string
}

export const getOssToken = () => {
  return request.get<{ code: number; data: OssTokenResp }>('/oss/v1/oss/token')
}

export const uploadFile = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<{ code: number; data: { url: string } }>('/oss/v1/oss/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 60000
  })
}