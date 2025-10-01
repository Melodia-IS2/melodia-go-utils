package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	router       *chi.Mux
	errorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

type RouterConfig struct {
	ErrorHandler            func(w http.ResponseWriter, r *http.Request, err error)
	NotFoundHandler         func(w http.ResponseWriter, r *http.Request)
	MethodNotAllowedHandler func(w http.ResponseWriter, r *http.Request)
}

func NewRouter(cfg *RouterConfig) *Router {
	rt := &Router{
		router: chi.NewRouter(),
	}

	if cfg != nil {
		if cfg.ErrorHandler != nil {
			rt.errorHandler = cfg.ErrorHandler
		}
		if cfg.NotFoundHandler != nil {
			rt.router.NotFound(cfg.NotFoundHandler)
		}
		if cfg.MethodNotAllowedHandler != nil {
			rt.router.MethodNotAllowed(cfg.MethodNotAllowedHandler)
		}
	}

	if rt.errorHandler == nil {
		rt.errorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	return rt
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.router.ServeHTTP(w, r)
}

func (rt *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	rt.router.Use(middlewares...)
}

func (rt *Router) UseErrorMiddleware(middlewares ...func(w http.ResponseWriter, r *http.Request) error) {
	for _, middleware := range middlewares {
		rt.router.Use(rt.executeMiddleware(middleware))
	}
}

func (rt *Router) Route(path string, fn func(r Router)) {
	rt.router.Route(path, func(r chi.Router) {
		subRouter := &Router{
			router:       r.(*chi.Mux),
			errorHandler: rt.errorHandler,
		}
		fn(*subRouter)
	})
}

func (rt *Router) executeHandler(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			rt.errorHandler(w, r, err)
		}
	}
}

func (rt *Router) executeMiddleware(middleware func(w http.ResponseWriter, r *http.Request) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := middleware(w, r); err != nil {
				rt.errorHandler(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (rt *Router) Get(path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	rt.router.Get(path, rt.executeHandler(handler))
}

func (rt *Router) Group(fn func(r *Router)) {
	rt.router.Group(func(r chi.Router) {
		subRouter := &Router{
			router:       r.(*chi.Mux),
			errorHandler: rt.errorHandler,
		}
		fn(subRouter)
	})
}

func (rt *Router) Post(path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	rt.router.Post(path, rt.executeHandler(handler))
}

func (rt *Router) Put(path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	rt.router.Put(path, rt.executeHandler(handler))
}

func (rt *Router) Delete(path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	rt.router.Delete(path, rt.executeHandler(handler))
}

type CanRegister interface {
	Register(rt *Router)
}
