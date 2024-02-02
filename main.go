package main

import (
	"github.com/adasarpan404/roomies-be/controllers"
	"github.com/adasarpan404/roomies-be/environment"
	"github.com/adasarpan404/roomies-be/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	router.Use(controllers.Authenticate())
	routes.UserRoutes(router)
	router.Run(":" + environment.PORT)
}
