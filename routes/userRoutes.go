package routes

import "github.com/gin-gonic/gin"

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/user/me")
}
