package route

import (
	"github.com/kouhei-github/book-advertisement-site/controller"
	"net/http"
)

type Router struct {
	Mutex *http.ServeMux
}

func (router *Router) GetRouter() {
	router.Mutex.HandleFunc("/api/v1/hello-world", controller.HelloWorldHandler)
}
