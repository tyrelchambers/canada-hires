package router

import (
	"canada-hires/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NonCompliantRoutes sets up routes for non-compliant employers API
func NonCompliantRoutes(controller *controllers.NonCompliantController, authMW func(http.Handler) http.Handler) func(chi.Router) {
	return func(r chi.Router) {
		// Public routes for non-compliant employers data
		r.Route("/non-compliant", func(r chi.Router) {
			r.Get("/employers", controller.GetNonCompliantEmployers)
			r.Get("/reasons", controller.GetNonCompliantReasons)
			r.Get("/locations", controller.GetNonCompliantLocations)
			r.Get("/employers/postal-code/{postal_code}", controller.GetNonCompliantEmployersByPostalCode)
			r.Get("/employers/coordinates/{lat}/{lng}", controller.GetNonCompliantEmployersByCoordinates)
		})

		// Admin routes for scraping operations
		r.Route("/admin/non-compliant", func(r chi.Router) {
			r.Use(authMW) // Apply authentication middleware
			r.Post("/scrape", controller.TriggerNonCompliantScraper)
			r.Get("/status", controller.GetNonCompliantScrapingStatus)
			r.Post("/geocode", controller.TriggerNonCompliantGeocoding)
		})
	}
}