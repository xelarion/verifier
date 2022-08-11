package verifier

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// An Option configures a verifier.
type Option interface {
	Apply(*Verifier)
}

// OptionFunc is a function that configures a verifier.
type OptionFunc func(*Verifier)

// Apply calls f(verifier)
func (f OptionFunc) Apply(verifier *Verifier) {
	f(verifier)
}

// WithJwtTimeFunc can be used to set jwtTimeFunc.
func WithJwtTimeFunc(timeFunc func() time.Time) Option {
	return OptionFunc(func(v *Verifier) {
		v.jwtTimeFunc = timeFunc
		jwt.TimeFunc = timeFunc
	})
}

// WithSourceName can be used to set sourceName.
func WithSourceName(sourceName string) Option {
	return OptionFunc(func(v *Verifier) {
		v.sourceName = sourceName
	})
}

// WithTokenExpireDuration can be used to set tokenExpireDuration.
func WithTokenExpireDuration(duration time.Duration) Option {
	return OptionFunc(func(v *Verifier) {
		v.tokenExpireDuration = duration
	})
}

// WithAuthExpireDuration can be used to set authExpireDuration.
func WithAuthExpireDuration(duration time.Duration) Option {
	return OptionFunc(func(v *Verifier) {
		v.authExpireDuration = duration
	})
}

// WithTempTokenExpireDuration can be used to set tempTokenExpireDuration.
func WithTempTokenExpireDuration(duration time.Duration) Option {
	return OptionFunc(func(v *Verifier) {
		v.tempTokenExpireDuration = duration
	})
}
