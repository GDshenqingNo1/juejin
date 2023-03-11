package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	Config *Config
}

type Config struct {
	SecretKey  string
	ExpireTime int64
	BufferTime int64
	Issuer     string
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("toke not active yet")
	TokenMalformed   = errors.New("not a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT(config *Config) *JWT {
	return &JWT{Config: config}
}

func (j *JWT) CreateClaims(baseClaims *BaseClaims) CustomClaims {
	claims := CustomClaims{
		BufferTime: j.Config.BufferTime,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.Config.ExpireTime) * time.Second)),
			Issuer:    j.Config.Issuer,
		},
		BaseClaims: *baseClaims,
	}
	return claims
}

func (j *JWT) GenerateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	signingKey := []byte(j.Config.SecretKey)
	return token.SignedString(signingKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	signingKey := []byte(j.Config.SecretKey)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return signingKey, nil
	})
	if err != nil {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

func (j *JWT) GenerateIdToken(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, *claims)
	signingKey := []byte(j.Config.SecretKey)
	return token.SignedString(signingKey)
}

func (j *JWT) CreateOIDCClaims(sub string, aud []string) jwt.RegisteredClaims {
	claims := jwt.RegisteredClaims{
		Issuer:    j.Config.Issuer,
		Subject:   sub,
		Audience:  aud,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.Config.ExpireTime) * time.Second)),
		NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)),
	}
	return claims
}

func (j *JWT) ParseIdToken(tokenString string) (*jwt.RegisteredClaims, error) {
	signingKey := []byte(j.Config.SecretKey)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return signingKey, nil
	})
	if err != nil {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
