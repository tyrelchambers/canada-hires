package router

import (
	"canada-hires/controllers"

	"github.com/go-chi/chi/v5"
)

func SearchRoutes(searchController controllers.SearchController) func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/search", func(r chi.Router) {
			r.Get("/", searchController.Search)
		})
	}
}