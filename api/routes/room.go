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
	room := rr.router.Gin.Group("rooms").Use()
	{
		room.GET("", rr.roomController.GetAllRoom)
		room.GET("/create", rr.roomController.CreateRoom)
		room.GET("/:uid", rr.roomController.JoinRoom)
	}
	chatRoom := rr.router.Gin.Group("rooms/chat")
	{
		chatRoom.GET("/:id", rr.roomController.ChatRoom)
	}
	roomWebsockets := rr.router.Gin.Group("rooms/websockets")
	{
		roomWebsockets.GET("/:uid", rr.roomController.RoomSocket)
		roomWebsockets.GET("/chat/:uid", rr.roomController.ChatRoomSocket)
		roomWebsockets.GET("viewer/:id", rr.roomController.RoomViewerSocket)
	}
}
