package verifier

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	*jwt.StandardClaims
	UserId uint
	Uuid   string                 // user_id 与 随机 uuid 组合，允许同一账号多次登录
	Data   map[string]interface{} // 自定义数据
}

func (c *CustomClaims) IsParamsValid() bool {
	return c.UserId > 0 && len(c.Uuid) > 0
}
