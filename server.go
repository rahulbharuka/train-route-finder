package main

import (
	"log"
	"os"

	"github.com/rahulbharuka/train-route-finder/logic"
	"github.com/rahulbharuka/train-route-finder/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// set release mode logging.
	gin.SetMode(gin.ReleaseMode)

	// create default Gin router
	router := gin.New()

	// init log middleware
	router.Use(gin.Logger())

	// init recovery middleware
	router.Use(gin.Recovery())

	// init rail network
	repository.RailNetworkInit()

	// get logic handler
	h := logic.GetHandler()

	// API handlers.
	router.GET("/routes", h.Routes)

	// run app on the specified port
	router.Run(":" + port)
}
