package Routes

import (
	"RMS/Handler"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {

	route := chi.NewRouter()

	route.Route("/v1", func(r chi.Router) {
		r.Get("/home", Handler.Home)
	})

}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
