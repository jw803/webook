package ginx

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Data Staff  `json:"data"`
	Iss  string `json:"iss"`
	Iat  int    `json:"iat"`
	Exp  int    `json:"exp"`
	jwt.RegisteredClaims
}

type Staff struct {
	MerchantId string `json:"merchant_id"`
	StaffId    string `json:"staff_id"`
	ClientId   string `json:"client_id"`
	SessionId  string `json:"session_id"`
}
