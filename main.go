package main

import (
	"bootstrap/webrtc/bootstrap"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	// infrastructure.InitializeLogger()
	// infrastructure.Logger.Info("started Server")
	godotenv.Load()
	fx.New(bootstrap.Module).Run()

}
