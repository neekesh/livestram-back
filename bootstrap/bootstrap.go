package bootstrap

import (
	"bootstrap/webrtc/api/controllers"
	"bootstrap/webrtc/api/routes"
	infrastructure "bootstrap/webrtc/infrastructure"
	"os"

	"go.uber.org/fx"
)

var Module = fx.Options(
	controllers.Module,
	infrastructure.Module,
	routes.Module,
	// services.Module,
	// repository.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	handler infrastructure.Router,
	routes routes.Routes,

) {
	routes.Setup()
	handler.Gin.Run(":" + os.Getenv("ServerPort"))

}
