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

func fail(request *routing.Context, errMsg string) error {
	request.Response.SetBody([]byte(errMsg))
	request.Response.SetStatusCode(401)
	return errors.New(errMsg)
}

func validateJwt(request *routing.Context) (uuid.UUID, error) {
	tokenString := string(request.Request.Header.Peek("Authorization"))

	if tokenString == "" {
		return uuid.UUID{}, errors.New("missing JWT token")
	}

	// Decode and validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fail(request, "unexpected JWT signing method : "+token.Header["alg"].(string))
		}
		key, _ := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		return key, nil
	})

	if err != nil {
		return uuid.UUID{}, fail(request, "failed to get claims : "+err.Error())
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID, err := uuid.FromString(claims["uuid"].(string))
	if err != nil {
		return uuid.UUID{}, fail(request, "missing userID in JWT")
	}
	return userID, nil
}
