package middlewares

import (
	"be-tickitz/utils"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("APP_SECRET")
		token := strings.Split(c.GetHeader("Authorization"), "Bearer ")

		if len(token) < 2 {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Unauthorized",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		redisClient := utils.RedisConnect()
		defer redisClient.Close()

		blacklisted, _ := redisClient.Get(context.Background(), "blacklist:"+token[1]).Result()
		if blacklisted == "true" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Unauthorized",
				Errors:  "Token has been blacklisted",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		rawToken, err := jwt.Parse(token[1], func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Invalid token",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userId := rawToken.Claims.(jwt.MapClaims)["userId"]
		role := rawToken.Claims.(jwt.MapClaims)["role"]
		c.Set("userId", userId)
		c.Set("role", role)
		c.Next()
	}
}
