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
	db, _ := utils.DBConnect()
	godotenv.Load()
	defer db.Close()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Static("/uploads/profiles", "./uploads/profiles")
	r.Static("/uploads/posters", "./uploads/posters")
	r.Static("/uploads/backdrops", "./uploads/backdrops")
	routers.CombineRouter(r)
	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
