package routes

import (
	"bootstrap/webrtc/api/controllers"
	"bootstrap/webrtc/infrastructure"
)

type StreamRoute struct {
	router           infrastructure.Router
	streamController controllers.StreamControllers
}

func NewStreamRoutes(
	router infrastructure.Router,
	streamController controllers.StreamControllers,
) StreamRoute {
	return StreamRoute{
		router:           router,
		streamController: streamController,
	}
}

func (rr StreamRoute) Setup() {
	stream := rr.router.Gin.Group("stream").Use()
	{
		stream.GET("", rr.streamController.GetAllStream)
		stream.GET("/create", rr.streamController.CreateStream)
		stream.GET("/:suuid", rr.streamController.JoinStream)
		// stream.DELETE("/delete/:id", rr.streamController.DeleteStream)
	}
	chatStream := rr.router.Gin.Group("stream/chat")
	{
		chatStream.GET("/:ssuid", rr.streamController.ChatStream)
	}
	streamWebsockets := rr.router.Gin.Group("stream/websockets")
	{
		streamWebsockets.GET("/:id", rr.streamController.StreamSocket)
		streamWebsockets.GET("/chat/:id", rr.streamController.ChatStreamSocket)
		streamWebsockets.GET("viewer/:id", rr.streamController.StreamViewerSocket)
	}
}
