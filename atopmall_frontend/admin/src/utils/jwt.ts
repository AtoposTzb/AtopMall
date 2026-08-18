interface JwtPayload {
  ID: number
  NickName: string
  AuthorityID: number
  iss: string
  exp: number
  nbf: number
}

export function decodeJwt(token: string): JwtPayload | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const payload = parts[1]
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(decoded) as JwtPayload
  } catch {
    return null
  }
}

export function isAdminUser(token: string): boolean {
  const payload = decodeJwt(token)
  return payload?.AuthorityID === 2
}