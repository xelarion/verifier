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
func CreateToken(sourceId SourceId, data map[string]interface{}) (string, error) {
	return defaultVerifier.CreateToken(sourceId, data)
}

// RefreshToken 刷新 token (根据原有的 sourceId 和 uuid)
// data 为自定义数据
func RefreshToken(sourceId SourceId, uid string, data map[string]interface{}) (string, error) {
	return defaultVerifier.RefreshToken(sourceId, uid, data)
}

// DestroyToken 销毁 token
func DestroyToken(sourceId SourceId, uid string) error {
	return defaultVerifier.DestroyToken(sourceId, uid)
}

// DestroyAllToken 销毁 sourceId 的所有 token
func DestroyAllToken(sourceId SourceId) error {
	return defaultVerifier.DestroyAllToken(sourceId)
}
