package verifier

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	*jwt.StandardClaims
	SourceId SourceId
	Uuid     string                 // source_id 与 随机 uuid 组合，允许同一账号多次登录
	Data     map[string]interface{} // 自定义数据
}

type SourceId interface{}

func (c *CustomClaims) IsParamsValid() bool {
	return c.SourceId != nil && len(c.Uuid) > 0
}
