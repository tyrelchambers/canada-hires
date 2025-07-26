package router

import (
	"canada-hires/container"
	"canada-hires/controllers"

	"github.com/go-chi/chi/v5"
)

type ReportRouter interface {
	Init(r chi.Router)
}

type reportRouter struct {
	cn               *container.Container
	reportController controllers.ReportController
}

func NewReportRouter(cn *container.Container, reportController controllers.ReportController) ReportRouter {
	return &reportRouter{
		cn:               cn,
		reportController: reportController,
	}
}

func (rr *reportRouter) Init(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		r.Post("/", rr.reportController.CreateReport)
		r.Get("/business/{id}", rr.reportController.GetReportsByBusiness)
		r.Put("/{id}", rr.reportController.UpdateReport)
		r.Delete("/{id}", rr.reportController.DeleteReport)
	})
}
