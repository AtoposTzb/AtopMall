import request from "@/utils/request";

export interface MessageItem {
  id: number;
  user_id: number;
  type: number;
  subject: string;
  message: string;
  file: string;
}

export interface MessageListResponse {
  total: number;
  data: MessageItem[];
}

export const getMessageList = () => {
  return request.get<MessageListResponse>("/op/v1/message", {});
};