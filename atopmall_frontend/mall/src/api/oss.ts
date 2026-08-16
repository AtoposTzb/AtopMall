import request from '@/utils/request'

// 类型定义
export interface OssTokenResponse {
  code: number
  data: {
    UploadUrl: string
    Url: string
  }
}

// 获取文件上传 URL（MinIO 预签名 PUT）
export const getOssToken = () => {
  return request.get<OssTokenResponse>('/oss/v1/oss/token')
}

// 上传文件到 MinIO
export const uploadFile = async (file: File): Promise<string> => {
  const res = await getOssToken()
  const uploadUrl = res.data.UploadUrl
  const fileUrl = res.data.Url

  await fetch(uploadUrl, {
    method: 'PUT',
    body: file,
    headers: {
      'Content-Type': file.type
    }
  })

  return fileUrl
}
