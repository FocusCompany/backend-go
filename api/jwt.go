package api

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdlatRjRjogo3WojgGHFHYLugdUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0/gCFQsHUfQrSDv+MuSUMAe8jzKE4qW+jK+xQU9a03GUnKHkkle+Q0pX/g6jXZ7r1/xAK5Do2kQ+X5xK9cipRgEKwIDAQAB
-----END PUBLIC KEY-----`)

func ValidateJwt(tokenString string) (uuid.UUID, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected JWT signing method : "+token.Header["alg"].(string))
		}
		key, _ := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		return key, nil
	})

	if err != nil {
		return uuid.UUID{}, errors.New("failed to get claims : "+err.Error())
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID, err := uuid.FromString(claims["uuid"].(string))
	if err != nil {
		return uuid.UUID{}, errors.New("missing userID in JWT")
	}

	return userID, nil
}

func ValidateJwtFromRequest(request *routing.Context) (uuid.UUID, error) {
	// Decode and validate JWT

	tokenString := string(request.Request.Header.Peek("Authorization"))

	if tokenString == "" {
		return uuid.UUID{}, errors.New("missing JWT token")
	}

	userID, err := ValidateJwt(tokenString)
	if err != nil {
		request.Response.SetBody([]byte(err.Error()))
		request.Response.SetStatusCode(401)
		return uuid.UUID{}, err
	}
	return userID, nil
}
