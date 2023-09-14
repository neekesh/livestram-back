package routes

import (
	"bootstrap/webrtc/api/controllers"
	"bootstrap/webrtc/infrastructure"
)

type RoomRoute struct {
	router         infrastructure.Router
	roomController controllers.RoomControllers
}

func NewRoomRoutes(
	router infrastructure.Router,
	roomController controllers.RoomControllers,
) RoomRoute {
	return RoomRoute{
		router:         router,
		roomController: roomController,
	}
}

func (rr RoomRoute) Setup() {
	Room := rr.router.Gin.Group("rooms").Use()
	{
		Room.GET("", rr.roomController.GetAllRoom)
		Room.POST("/create", rr.roomController.PostRoom)
		Room.PUT("/update/:id", rr.roomController.UpdateRoom)
		Room.DELETE("/delete/:id", rr.roomController.DeleteRoom)
	}
}
