package forms

type ShoppingCartIntemForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShoppingCartUpdateForm struct {
	Nums  int32 `json:"nums" binding:"required,min=1"`
	Check *bool `json:"checked"`
}
