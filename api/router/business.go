package router

import (
	"canada-hires/container"
	"canada-hires/controllers"

	"github.com/go-chi/chi/v5"
)

type BusinessRouter interface {
	Init(r chi.Router)
}

type businessRouter struct {
	cn                 *container.Container
	businessController controllers.BusinessController
}

func NewBusinessRouter(cn *container.Container, businessController controllers.BusinessController) BusinessRouter {
	return &businessRouter{
		cn:                 cn,
		businessController: businessController,
	}
}

func (br *businessRouter) Init(r chi.Router) {
	r.Route("/businesses", func(r chi.Router) {
		r.Get("/", br.businessController.GetBusinesses)
		r.Post("/", br.businessController.CreateBusiness)
		r.Get("/{id}", br.businessController.GetBusiness)
		r.Put("/{id}", br.businessController.UpdateBusiness)
	})
}
