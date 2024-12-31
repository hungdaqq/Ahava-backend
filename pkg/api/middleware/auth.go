package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AdminAuthMiddleware(ctx *gin.Context) {

	accessToken := ctx.Request.Header.Get("Authorization")

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("accesssecret"), nil
	})
	if err != nil {
		// The access token is invalid.
		fmt.Println("error catches here")
		ctx.AbortWithStatus(401)
		return
	}

	ctx.Next()
}
