import request from '@/utils/request'

// 类型定义
export interface OssTokenResponse {
  code: number
  data: {
    upload_url: string
    url: string
  }
}

// 获取文件上传 URL（MinIO 预签名 PUT）
export const getOssToken = () => {
  return request.get<OssTokenResponse>('/oss/v1/oss/token')
}

// 上传文件到 MinIO（通过后端代理，避免浏览器直连 MinIO 的兼容性问题）
export const uploadFile = async (file: File): Promise<string> => {
  const formData = new FormData()
  formData.append('file', file)
  const res = await request.post<{ code: number; data: { url: string } }>(
    '/oss/v1/oss/upload',
    formData,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 60000
    }
  )
  return (res as any).data.url
}