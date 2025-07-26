package router

import (
	"canada-hires/container"
	"canada-hires/controllers"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(cn *container.Container, r *chi.Mux) {
	// Initialize routers
	ar := &authRouter{}
	br := &businessRouter{}
	rr := &reportRouter{}
	ur := &userRouter{}

	// Invoke the router initializers
	err := cn.Invoke(func(authController controllers.AuthController, businessController controllers.BusinessController, reportController controllers.ReportController, userController controllers.UserController, jobController *controllers.JobController, authMW func(http.Handler) http.Handler, requireMW func(http.Handler) http.Handler) {
		*ar = *NewAuthRouter(cn, authController).(*authRouter)
		*br = *NewBusinessRouter(cn, businessController).(*businessRouter)
		*rr = *NewReportRouter(cn, reportController).(*reportRouter)
		*ur = *NewUserRouter(cn, userController, authMW, requireMW).(*userRouter)
		
		// Initialize job routes
		JobRoutes(r, jobController)
	})

	if err != nil {
		log.Fatal("Failed to initialize routes", "error", err)
	}

	r.Route("/api", func(r chi.Router) {
		// Inject routes
		ar.InjectAuthRoutes(r)
		br.Init(r)
		rr.Init(r)
		ur.Init(r)
		
		// Add LMIA routes
		err := cn.Invoke(func(lmiaController *controllers.LMIAController) {
			LMIARoutes(lmiaController)(r)
		})
		if err != nil {
			log.Error("Failed to initialize LMIA routes", "error", err)
		}
	})
}
