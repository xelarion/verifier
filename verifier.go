package verifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Verifier token 验证器
/**
token 生成和验证逻辑:

实现效果: 无操作后 3 小时自动退出，token 每 15 分钟刷新一次

token 生成逻辑:
	根据 sourceID 和 uuid 生成，token 过期时间为 15 分钟
	并以 sourceID 和 uuid 组合 为 key(keyA), token 为值 记录到存储器, 过期时间设置为 3 小时

验证 token 逻辑:
	大前提: keyA 存在于存储器, 不存在则不通过验证

	1. jwt 解析 token 没有 error
		=> token 有效 并且 存储器 中 (keyA) 对应的值 和 token 相同, 则通过验证
	2. jwt 解析 token 有 error, 且 error 为 token 过期
        => 1. token 与 存储器中 (keyA) 对应的值相同, 刷新 token, 通过验证并返回新 token
		=> 2. 当前 token 已被存储器记录(作为过渡token)， 通过验证
	3. 其他情况，不通过验证

PS: 刷新 token 时，将旧 token 记录到存储器 原因：
	并发请求时，只允许一个线程刷新 token，在新 token 未返回客户端前，短时间内允许旧 token 请求
*/
type Verifier struct {
	jwtKey      string           // jwt sign key
	jwtTimeFunc func() time.Time // jwt TimeFunc

	sourceName              string        // 指定 source, 用于区分不同资源登录, 影响 key, 如: user, admin, account, 默认为 user
	tokenExpireDuration     time.Duration // token 过期时长(在身份过期内，会自动刷新 token)
	authExpireDuration      time.Duration // 身份过期时长
	tempTokenExpireDuration time.Duration // 旧 token 临时保存过期时长

	tokenStorage TokenStorage // token 存储器
}

func New(jwtKey string, tokenStorage TokenStorage, options ...Option) *Verifier {
	verifier := Verifier{
		jwtKey:      jwtKey,
		jwtTimeFunc: time.Now,

		sourceName:              "user",
		tokenExpireDuration:     15 * time.Minute,
		authExpireDuration:      3 * time.Hour,
		tempTokenExpireDuration: 30 * time.Second,

		tokenStorage: tokenStorage,
	}

	for _, o := range options {
		o.Apply(&verifier)
	}

	return &verifier
}

// VerifyToken 验证 token (包含刷新 token 逻辑)
// 若刷新了 token 则返回新 token
// 若返回 error 则说明未通过验证
func (v *Verifier) VerifyToken(tokenStr string) (CustomClaims, string, error) {
	return v.verifyToken(tokenStr, true)
}

// IsTokenAuthorized token 是否通过身份验证(仅验证)
func (v *Verifier) IsTokenAuthorized(tokenStr string) (CustomClaims, bool) {
	claims, _, err := v.verifyToken(tokenStr, false)
	return claims, err == nil
}

// CreateToken 创建新 token
// data 为自定义数据
func (v *Verifier) CreateToken(sourceID any, data map[string]interface{}) (string, error) {
	uid := uuid.New().String()
	tokenStr, err := v.generateToken(sourceID, uid, data)
	if err != nil {
		return "", err
	}

	// sourceID 和 uuid 组合, 记录到存储器(允许同一账号多次登录)
	if err = v.tokenStorage.Set(
		v.tokenStorageKey(sourceID, uid),
		tokenStr,
		v.authExpireDuration,
	); err != nil {
		return "", err
	}

	return tokenStr, nil
}

// RefreshToken 刷新 token (根据原有的 sourceID 和 uuid)
// data 为自定义数据
func (v *Verifier) RefreshToken(sourceID any, uid string, data map[string]interface{}) (string, error) {
	tokenStr, err := v.generateToken(sourceID, uid, data)
	if err != nil {
		return "", err
	}

	// sourceID 和 uuid 组合, 记录到存储器
	if err = v.tokenStorage.Set(
		v.tokenStorageKey(sourceID, uid),
		tokenStr,
		v.authExpireDuration,
	); err != nil {
		return "", err
	}

	return tokenStr, nil
}

// DestroyToken 销毁 token
func (v *Verifier) DestroyToken(sourceID any, uid string) error {
	return v.tokenStorage.Del(v.tokenStorageKey(sourceID, uid))
}

// DestroyAllToken 销毁 sourceID 的所有 token
func (v *Verifier) DestroyAllToken(sourceID any) error {
	return v.tokenStorage.DelByKeyPrefix(v.tokenStorageKeySourceIdFilterPrefix(sourceID))
}

// 验证 token
// 若刷新了 token 则返回新 token
// 若返回 error 则说明未通过验证
func (v *Verifier) verifyToken(tokenStr string, isRefreshToken bool) (CustomClaims, string, error) {
	var (
		jwtToken        *jwt.Token   // jwt token
		jwtErr          error        // jwt parse err
		getStorageErr   error        // get storage token error
		claims          CustomClaims // jwt custom claims
		storageTokenStr string       // token of storage
	)

	// jwt 解析 tokenStr
	jwtToken, jwtErr = jwt.ParseWithClaims(tokenStr, &claims, v.jwtTokenKeyFunc())
	if !claims.IsParamsValid() {
		return CustomClaims{}, "", TokenInvalidError
	}

	// 从存储器中获取 token
	storageTokenStr, getStorageErr = v.tokenStorage.Get(v.tokenStorageKey(claims.SourceID, claims.UUID))
	if getStorageErr != nil {
		return CustomClaims{}, "", TokenInvalidError
	}

	if jwtErr != nil {
		// tokenStr 已过期
		if errors.Is(jwtErr, jwt.ErrTokenExpired) {
			// 存储器中 token 与 tokenStr 相同, 即身份未过期、仅 token 过期, 需要刷新 token
			if storageTokenStr == tokenStr {
				// 将旧的过期 token 保存到存储器中过渡(并发时可能多个请求携带旧 token)
				if isRefreshToken && v.tokenStorage.SetNX(tokenStr, "1", v.tempTokenExpireDuration) {
					// 刷新 token, 并返回
					newTokenStr, err := v.RefreshToken(claims.SourceID, claims.UUID, claims.Data)
					if err != nil {
						return CustomClaims{}, "", TokenInvalidError
					}

					return claims, newTokenStr, nil
				} else {
					// 返回 token 有效(其他线程刷新了 token)
					return claims, "", nil
				}
			} else {
				// 存储器中 token 与 tokenStr 不同, 可能 tokenStr 已被刷新
				// 如果 tokenStr 被临时保存, 则通过验证, 否则为弃用
				if v.tokenStorage.Exists(tokenStr) {
					return claims, "", nil
				}
			}
		}
	} else {
		// token 有效 且 与存储器中值相同
		if jwtToken.Valid && storageTokenStr == tokenStr {
			return claims, "", nil
		}
	}

	// 其他错误情况
	return CustomClaims{}, "", TokenInvalidError
}

// token 存储的 key
func (v *Verifier) tokenStorageKey(sourceID any, uid string) string {
	return fmt.Sprintf("verifier:%v:%v:uid:%s", v.sourceName, sourceID, uid)
}

// token 存储的 key 中的 sourceID 前缀, 用于过滤 sourceID 的所有 key
func (v *Verifier) tokenStorageKeySourceIdFilterPrefix(sourceID any) string {
	return fmt.Sprintf("verifier:%v:%v:uid:", v.sourceName, sourceID)
}

func (v *Verifier) jwtTokenKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(v.jwtKey), nil
	}
}

// 生成token
// data 为自定义数据
func (v *Verifier) generateToken(sourceID any, uid string, data map[string]interface{}) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(v.jwtTimeFunc().Add(v.tokenExpireDuration)),
		},
		SourceID: sourceID,
		UUID:     uid,
		Data:     data,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(v.jwtKey))
}
