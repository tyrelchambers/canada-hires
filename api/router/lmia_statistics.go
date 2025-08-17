package router

import (
	"canada-hires/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// LMIAStatisticsRoutes sets up routes for LMIA statistics endpoints
func LMIAStatisticsRoutes(controller controllers.LMIAStatisticsController, authMW func(http.Handler) http.Handler) func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/lmia/statistics", func(r chi.Router) {
			// Public routes
			r.Get("/daily", controller.GetDailyTrends)
			r.Get("/monthly", controller.GetMonthlyTrends)
			r.Get("/summary", controller.GetTrendsSummary)
			r.Get("/regional", controller.GetRegionalStats)

			// Admin routes (require authentication)
			r.Group(func(r chi.Router) {
				r.Use(authMW)
				r.Post("/backfill", controller.BackfillHistoricalStatistics)
				r.Post("/generate", controller.GenerateStatisticsForDateRange)
				r.Post("/aggregate", controller.RunDailyAggregation)
			})
		})
	}
}