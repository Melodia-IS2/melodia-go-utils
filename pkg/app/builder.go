package app

import (
	"github.com/Melodia-IS2/melodia-go-utils/pkg/http"
	"github.com/Melodia-IS2/melodia-go-utils/pkg/router"
)

type Builder struct {
	app    App
	port   string
	router *router.Router
}

func NewBuilder(rtcfg *router.RouterConfig, port string) (*Builder, error) {
	if rtcfg == nil {
		rtcfg = &router.RouterConfig{
			ErrorHandler:            http.ErrorHandler,
			NotFoundHandler:         http.NotFoundHandler,
			MethodNotAllowedHandler: http.MethodNotAllowedHandler,
		}
	}

	return &Builder{
		router: router.NewRouter(rtcfg),
		port:   port,
	}, nil
}

func (b *Builder) RegisterHandler(handler router.CanRegister) *Builder {
	handler.Register(b.router)
	return b
}

func (b *Builder) Build() *App {
	return &App{
		router: b.router,
		port:   b.port,
	}
}
