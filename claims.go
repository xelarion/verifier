package verifier

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	*jwt.StandardClaims
	SourceId uint
	Uuid     string                 // source_id 与 随机 uuid 组合，允许同一账号多次登录
	Data     map[string]interface{} // 自定义数据
}

func (c *CustomClaims) IsParamsValid() bool {
	return c.SourceId > 0 && len(c.Uuid) > 0
}
