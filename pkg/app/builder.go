package app

import (
	"context"
	"net/http"

	"github.com/Melodia-IS2/melodia-events/pkg/suscriber/kafka"
	httpUtils "github.com/Melodia-IS2/melodia-go-utils/pkg/http"
	"github.com/Melodia-IS2/melodia-go-utils/pkg/router"
)

type Builder struct {
	app       App
	port      string
	router    *router.Router
	workers   []Worker
	consumers []kafka.Consumer
}

func NewBuilder(rtcfg *router.RouterConfig, port string) (*Builder, error) {
	if rtcfg == nil {
		rtcfg = &router.RouterConfig{
			ErrorHandler:            httpUtils.ErrorHandler,
			NotFoundHandler:         httpUtils.NotFoundHandler,
			MethodNotAllowedHandler: httpUtils.MethodNotAllowedHandler,
		}
	}

	return &Builder{
		router: router.NewRouter(rtcfg),
		port:   port,
	}, nil
}

func (b *Builder) RegisterWorker(worker Worker) *Builder {
	b.workers = append(b.workers, worker)
	return b
}

func (b *Builder) RegisterHandler(handler router.CanRegister) *Builder {
	handler.Register(b.router)
	return b
}

func (b *Builder) RegisterMiddleware(middleware func(http.Handler) http.Handler) *Builder {
	b.router.Use(middleware)
	return b
}

func (b *Builder) RegisterConsumer(consumer kafka.Consumer) *Builder {
	b.consumers = append(b.consumers, consumer)
	return b
}

func (b *Builder) Build() *App {
	for _, worker := range b.workers {
		go worker.Start()
	}
	for _, consumer := range b.consumers {
		go consumer.Start(context.Background())
	}
	return &App{
		router: b.router,
		port:   b.port,
	}
}
