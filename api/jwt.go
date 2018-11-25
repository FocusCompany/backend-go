package api

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdlatRjRjogo3WojgGHFHYLugdUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0/gCFQsHUfQrSDv+MuSUMAe8jzKE4qW+jK+xQU9a03GUnKHkkle+Q0pX/g6jXZ7r1/xAK5Do2kQ+X5xK9cipRgEKwIDAQAB
-----END PUBLIC KEY-----`)


// Check that a valid JWT exists in the Authorization header
func RequireBasicJwt(request *routing.Context) error {
	tokenString := string(request.Request.Header.Peek("Authorization"))

	userId, err := ValidateJwt(tokenString)
	if err != nil {
		fmt.Printf("RequireBasicJwt JWT: %v", err)
		fmt.Println("")
		request.Abort()
		return err
	}

	request.Set("userId", userId)
	return nil
}

// Validate JWT ensures the given token is valid and returns the user ID it contains.
// Returns an error and an empty user ID if the token is not valid.
func ValidateJwt(tokenString string) (uuid.UUID, error) {
	fmt.Printf("JWT = %s\n", tokenString)

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
