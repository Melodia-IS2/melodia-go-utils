package app

import (
	"fmt"
	"net/http"

	"github.com/Melodia-IS2/melodia-go-utils/pkg/router"
)

type App struct {
	router *router.Router
	port   string
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *App) Run() error {
	if a.port == "" {
		a.port = "8080"
	}
	fmt.Println("Server is running on port", a.port)
	err := http.ListenAndServe(":"+a.port, a.router)
	if err != nil {
		return err
	}
	return nil
}
