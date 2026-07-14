package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// 自定义JWT claims结构体加密解密的时候需要使用这个结构体
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityID uint
	jwt.RegisteredClaims
}
