package Routes

import (
	"RMS/Handler"
	"RMS/Middleware"
	"RMS/Models"
	"RMS/Utils"
	"context"
	"github.com/go-chi/chi/v5"

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
		r.Get("/health", Utils.Health)
		r.Get("/home", Handler.Home)

		r.Post("/login", Handler.LoginUser)

		r.Group(func(r chi.Router) {

			r.Use(Middleware.Authenticate)

			r.Get("/dishesByRestaurant", Handler.DishesByRestaurant)
			r.Post("/logout", Handler.LogoutUser)

			r.Route("/admin", func(admin chi.Router) {
				admin.Use(Middleware.ShouldHaveRole(Models.RoleAdmin))

				admin.Post("/createUser", Handler.CreateUser)
				admin.Get("/getAllUsers", Handler.GetAllUsersByAdmin)
				admin.Post("/createRestaurants", Handler.CreateRestaurants)
				admin.Get("/getAllRestaurants", Handler.GetallRestaurants)
				admin.Post("/subAdminCreation", Handler.SubAdminCreation)
				admin.Get("/getAllSubadmins", Handler.GetAllSubAdmins)
				admin.Route("/createDish/{restaurantId}", func(restId chi.Router) {
					restId.Post("/", Handler.CreateDish)
				})
				admin.Get("/getAllDishes", Handler.GetAllDishes)
			})
			r.Route("/sub-admin", func(subAdmin chi.Router) {
				subAdmin.Use(Middleware.ShouldHaveRole(Models.RoleSubAdmin))
				subAdmin.Post("/createUser", Handler.CreateUser)
				subAdmin.Get("/userBySubAdmin", Handler.GetAllUsersBySubAdmin)
				subAdmin.Post("/createRestaurant", Handler.CreateRestaurants)
				subAdmin.Get("/getAllRetaurants", Handler.GetAllRestaurantsBySubAdmin)
				subAdmin.Route("/createDish/{restaurantId}", func(restId chi.Router) {
					restId.Post("/", Handler.CreateDish)
				})
				subAdmin.Get("/getDishesBySubadmin", Handler.GetAllDishesBySubAdmin)
			})
			r.Route("/user", func(user chi.Router) {
				user.Use(Middleware.ShouldHaveRole(Models.RoleUser))
				user.Get("/getAllRestaurants", Handler.GetallRestaurants)
				user.Get("/getAllDishes", Handler.GetAllDishes)
			})
		})
	})

	return &Server{
		Router: route,
	}

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
