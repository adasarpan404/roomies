package routes

import (
	"github.com/adasarpan404/roomies-be/controllers"
	"github.com/gin-gonic/gin"
)

func RoomRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/room", controllers.GetRooms())
	incomingRoutes.POST("/room", controllers.CreateRoom())
	incomingRoutes.PUT("/room/:id", controllers.UpdateRoom())
	incomingRoutes.DELETE("/room/:id", controllers.DeleteRoom())
}
