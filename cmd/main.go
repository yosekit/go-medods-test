package main

import (
	"fmt"
	"log"
	"net/http"

	"goauth/config"
	"goauth/internal/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	appConfig = config.App()
	router = gin.Default()
	err error
)

func main() {
	mapRoutes()

	log.Println("Starting server on port", appConfig.Port, "...")
	err = http.ListenAndServe(fmt.Sprintf(":%s", appConfig.Port), router)
	if err != nil {
		log.Fatal("Could not start server:", err)
	}
}

func mapRoutes() {
	router.POST("/auth/tokens", handlers.GenerateTokens)
	router.POST("/auth/refresh", handlers.RefreshTokens)
}