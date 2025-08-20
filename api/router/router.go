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
	adr := &adminRouter{}

	// Invoke the router initializers
	err := cn.Invoke(func(authController controllers.AuthController, businessController controllers.BusinessController, reportController controllers.ReportController, userController controllers.UserController, jobController *controllers.JobController, authMW func(http.Handler) http.Handler, requireMW func(http.Handler) http.Handler) {
		*ar = *NewAuthRouter(cn, authController).(*authRouter)
		*br = *NewBusinessRouter(cn, businessController).(*businessRouter)
		*rr = *NewReportRouter(cn, reportController, authMW).(*reportRouter)
		*ur = *NewUserRouter(cn, userController, authMW, requireMW).(*userRouter)
		*adr = *NewAdminRouter(cn, jobController, authMW).(*adminRouter)
		
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
		rr.InjectReportRoutes(r)
		ur.Init(r)
		adr.InjectAdminRoutes(r)
		
		// Add LMIA routes
		err := cn.Invoke(func(lmiaController *controllers.LMIAController) {
			LMIARoutes(lmiaController)(r)
		})
		if err != nil {
			log.Error("Failed to initialize LMIA routes", "error", err)
		}
		
		// Add subreddit routes
		err = cn.Invoke(func(subredditController *controllers.SubredditController, authMW func(http.Handler) http.Handler) {
			SubredditRoutes(subredditController, authMW)(r)
		})
		if err != nil {
			log.Error("Failed to initialize subreddit routes", "error", err)
		}
		
		// Add LMIA statistics routes
		err = cn.Invoke(func(statisticsController controllers.LMIAStatisticsController, authMW func(http.Handler) http.Handler) {
			LMIAStatisticsRoutes(statisticsController, authMW)(r)
		})
		if err != nil {
			log.Error("Failed to initialize LMIA statistics routes", "error", err)
		}
		
		// Add boycott routes
		err = cn.Invoke(func(boycottController controllers.BoycottController, authMW func(http.Handler) http.Handler) {
			BoycottRoutes(boycottController, authMW)(r)
		})
		if err != nil {
			log.Error("Failed to initialize boycott routes", "error", err)
		}
	})
}
