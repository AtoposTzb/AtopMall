import request from "@/utils/request";

// 类型定义
export interface GoodsItem {
  id: number
  name: string
  goods_sn: string
  shop_price: number
  front_image: string
  goods_brief: string
  images: string[]
  desc_images: string[]
  is_hot: boolean
  is_new: boolean
  on_sale: boolean
  ship_free: boolean
  brand: {
    id: number
    name: string
    logo: string
  }
  ctegory: {
    id: number
    name: string
  }
  desc: string
}

export interface GoodsListParams {
  p?: number;
  pnum?: number;
  q?: string;
  c?: number;
  b?: number;
  pmin?: number;
  pmax?: number;
  ishot?: number;
  isnew?: number;
  istab?: number;
}

export interface GoodsListResponse {
  data: GoodsItem[];
  total: number;
}

export interface GoodsStockResponse {
  goodsId: number;
  num: number;
}

// 商品列表
export const getGoodsList = (params?: GoodsListParams) => {
  return request.get<GoodsListResponse>("/g/v1/goods", params);
};

// 商品详情
export const getGoodsDetail = (id: number) => {
  return request.get<GoodsItem>(`/g/v1/goods/${id}`);
};

// 商品库存
export const getGoodsStock = (id: number) => {
  return request.get<GoodsStockResponse>(`/g/v1/goods/${id}/stocks`);
};