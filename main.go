package main

import (
	"be-tickitz/routers"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title CRUD
// @version         1.0
// @description     CRUD Swagger
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	r := gin.Default()
	godotenv.Load()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})
	r.Static("/uploads/profiles", "./uploads/profiles")
	r.Static("/uploads/poster", "./uploads/poster")
	r.Static("/uploads/backdrop", "./uploads/backdrop")

	routers.CombineRouter(r)
	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
