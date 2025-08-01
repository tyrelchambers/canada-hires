package router

import (
	"canada-hires/controllers"

	"github.com/go-chi/chi/v5"
)

func LMIARoutes(lmiaController *controllers.LMIAController) func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/lmia", func(r chi.Router) {
			// Public endpoints for LMIA data
			r.Get("/employers/search", lmiaController.SearchEmployers)
			r.Get("/employers/location", lmiaController.GetEmployersByLocation)
			r.Get("/employers/resource/{resourceID}", lmiaController.GetEmployersByResource)
			r.Get("/resources", lmiaController.GetResources)
			r.Get("/stats", lmiaController.GetStats)
			r.Get("/status", lmiaController.GetUpdateStatus)
			
		})
	}
}