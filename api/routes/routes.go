package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewStreamRoutes),
)

type Routes []Route

type Route interface {
	Setup()
}

func NewRoutes(
	streamRoutes StreamRoute,
) Routes {
	return Routes{
		streamRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
