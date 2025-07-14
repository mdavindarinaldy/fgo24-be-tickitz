package main

import (
	"be-tickitz/routers"
	"be-tickitz/utils"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://146.190.102.54:9602", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	db, _ := utils.DBConnect()
	godotenv.Load()
	defer db.Close()
	r.Static("/uploads/profiles", "./uploads/profiles")
	r.Static("/uploads/poster", "./uploads/poster")
	r.Static("/uploads/backdrop", "./uploads/backdrop")
	routers.CombineRouter(r)
	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
