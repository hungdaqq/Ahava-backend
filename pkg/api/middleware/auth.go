package middleware

import (
	"fmt"
	"net/http"
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

func UserAuthMiddleware(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		ctx.Abort()
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return []byte("ahava"), nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		ctx.Abort()
		return
	}

	fmt.Println("claims", claims)

	role, ok := claims["role"].(string)
	if !ok || role != "client" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		ctx.Abort()
		return
	}

	id, ok := claims["id"].(float64)
	if !ok || id == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "error in retrieving id"})
		ctx.Abort()
		return
	}

	ctx.Set("role", role)
	ctx.Set("id", int(id))

	ctx.Next()
}
