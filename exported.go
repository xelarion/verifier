package verifier

var (
	// defaultVerifier is default verifier
	defaultVerifier *Verifier
)

func InitDefaultVerifier(jwtKey string, tokenStorage TokenStorage, options ...Option) {
	defaultVerifier = New(jwtKey, tokenStorage, options...)
}

func DefaultVerifierVerifier() *Verifier {
	return defaultVerifier
}

// VerifyToken 验证 token (包含刷新 token 逻辑)
func VerifyToken(tokenStr string) (CustomClaims, string, error) {
	return defaultVerifier.VerifyToken(tokenStr)
}

// IsTokenAuthorized token 是否通过身份验证(仅验证)
func IsTokenAuthorized(tokenStr string) (CustomClaims, bool) {
	return defaultVerifier.IsTokenAuthorized(tokenStr)
}

// CreateToken 创建新 token
// data 为自定义数据
func CreateToken(sourceID any, data map[string]interface{}) (string, error) {
	return defaultVerifier.CreateToken(sourceID, data)
}

// RefreshToken 刷新 token (根据原有的 sourceID 和 uuid)
// data 为自定义数据
func RefreshToken(sourceID any, uid string, data map[string]interface{}) (string, error) {
	return defaultVerifier.RefreshToken(sourceID, uid, data)
}

// DestroyToken 销毁 token
func DestroyToken(sourceID any, uid string) error {
	return defaultVerifier.DestroyToken(sourceID, uid)
}

// DestroyAllToken 销毁 sourceID 的所有 token
func DestroyAllToken(sourceID any) error {
	return defaultVerifier.DestroyAllToken(sourceID)
}
