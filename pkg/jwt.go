package pkg

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningMethon = errors.New("signing method is invalid")
	ErrClaimsFailed         = errors.New("claims failed")
	ErrTokenExpired         = errors.New("token is expired")
)

type JWTpkg interface {
	GenerateToken(Name string) (string, error)
	ParseToken(token string) (string, error)
}

type JWTgo struct {
}

func NewJWTgo() JWTgo {
	return JWTgo{}
}

func (j *JWTgo) GenerateToken(Name string) (string, error) {
	payload := jwt.MapClaims{
		"name": Name,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))
	if err != nil {
		log.Println("GenerateToken err")
		return "", err
	}
	return t, nil
}

func (j *JWTgo) ParseToken(token string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		key := os.Getenv("JWTSECRET")
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("Keyfunc err: mismatch signing method")
			return nil, ErrInvalidSigningMethon
		}
		return []byte(key), nil
	}
	tkn, err := jwt.Parse(token, keyFunc)
	if err != nil {
		log.Println("JWT parse err", err)
		return "", err
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrClaimsFailed
	}
	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return "", ErrTokenExpired
	}
	name, ok := claims["name"].(string)
	if !ok {
		return "", ErrClaimsFailed
	}
	return name, nil
}
