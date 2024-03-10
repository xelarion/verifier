package verifier

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	*jwt.RegisteredClaims
	SourceID any
	UUID     string                 // source_id 与 随机 uuid 组合，允许同一账号多次登录
	Data     map[string]interface{} // 自定义数据
}

func (c *CustomClaims) IsParamsValid() bool {
	return c.SourceID != nil && len(c.UUID) > 0
}
