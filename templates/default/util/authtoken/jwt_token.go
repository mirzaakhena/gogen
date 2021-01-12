package authtoken

import (
	"fmt"
	"log"
	"time"

	gonanoid "github.com/matoous/go-nanoid"

	"github.com/dgrijalva/jwt-go"
)

type JwtToken interface {
	Generate(subject, audience string, extendInfo interface{}, expired time.Duration) string
	Validate(subject, tokenString string) (*CustomClaims, bool)
}

type jwtToken struct {
	secretKey []byte
	issuer    string
}

type CustomClaims struct {
	jwt.StandardClaims
	ExtendInfo interface{} `json:"extendInfo,omitempty"`
}

func GenerateSecretKey() string {
	key, _ := gonanoid.Generate("abcdef1234567890", 128)
	return key
}

func New(issuer string, secretKey []byte) JwtToken {
	return &jwtToken{
		issuer:    issuer,
		secretKey: secretKey,
	}
}

func (j *jwtToken) Generate(subject, audience string, extendInfo interface{}, expired time.Duration) string {
	now := time.Now().UTC()
	claims := &CustomClaims{
		ExtendInfo: extendInfo,
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			Issuer:    j.issuer,
			Audience:  audience,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(expired).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(j.secretKey)
	return tokenString
}

func (j *jwtToken) Validate(subject, tokenString string) (*CustomClaims, bool) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		log.Printf(err.Error())
		return nil, false
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, false
	}

	if claims.Subject != subject {
		return nil, false
	}

	return claims, token.Valid
}
