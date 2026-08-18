import request from "@/utils/request";

export interface LoginParams {
  mobile: string;
  password: string;
  captcha: string;
  captcha_id: string;
}

export interface LoginResponse {
  id: number;
  nick_name: string;
  token: string;
  expired_at?: number;
}

export interface CaptchaData {
  id: string;
  picBase64: string;
}

export interface UserInfo {
  id: number;
  name: string;
  nick_name?: string;
  mobile: string;
  gender: string;
  birthday: string;
  role?: number;
}

export const login = (data: LoginParams) => {
  return request.post<LoginResponse>("/u/v1/user/pwd_login", data);
};

export const getCaptcha = () => {
  return request.get<CaptchaData>("/u/v1/base/captcha");
};

export const getUserDetail = () => {
  return request.get<UserInfo>("/u/v1/user/detail");
};

export const getUserList = (params?: { pn?: number; psize?: number }) => {
  return request.get<UserInfo[]>("/u/v1/user/list", params);
};
