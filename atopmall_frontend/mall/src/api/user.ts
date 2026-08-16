import request from "@/utils/request";

// 类型定义
export interface LoginParams {
  mobile: string;
  password: string;
  captcha: string;
  captcha_id: string;
}

export interface RegisterParams {
  mobile: string;
  password: string;
  email: string;
  code: string;
}

export interface LoginResponse {
  id: number;
  nick_name: string;
  token: string;
  expired_at?: number;
}

export interface UserInfo {
  name: string;
  mobile: string;
  gender: string;
  birthday: string;
  email: string;
}

export interface UpdateUserParams {
  name: string;
  gender: string;
  birthday: string;
}

export interface CaptchaResponse {
  id: string;
  picBase64: string;
}

// 图片验证码
export const getCaptcha = () => {
  return request.get<CaptchaResponse>("/u/v1/base/captcha");
};

// 发送邮箱验证码
export const sendEmailCode = (data: { email: string; type: number }) => {
  return request.post("/u/v1/base/send-code", data);
};

// 登录
export const login = (data: LoginParams) => {
  return request.post<LoginResponse>("/u/v1/user/pwd_login", data);
};

// 注册
export const register = (data: RegisterParams) => {
  return request.post<LoginResponse>("/u/v1/user/register", data);
};

// 用户详情
export const getUserDetail = () => {
  return request.get<UserInfo>("/u/v1/user/detail");
};

// 更新用户信息
export const updateUser = (data: UpdateUserParams) => {
  return request.put("/u/v1/user/update", data);
};

// 用户列表（管理员）
export const getUserList = (params?: { pn?: number; psize?: number }) => {
  return request.get<UserInfo[]>("/u/v1/user/list", params);
};