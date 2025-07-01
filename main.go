package main

import (
	"be-tickitz/routers"
	"be-tickitz/utils"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	db, _ := utils.DBConnect()
	godotenv.Load()
	defer db.Close()
	r := gin.Default()
	routers.CombineRouter(r)
	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
