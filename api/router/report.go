package router

import (
	"canada-hires/container"
	"canada-hires/controllers"
	"canada-hires/middleware"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type ReportRouter interface {
	Init(r chi.Router)
	InjectReportRoutes(r chi.Router)
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

func (rr *reportRouter) InjectReportRoutes(r chi.Router) {
	err := rr.cn.Invoke(func(reportController controllers.ReportController) {
		rr := NewReportRouter(rr.cn, reportController)
		rr.Init(r)
	})

	if err != nil {
		log.Fatal("Failed to initialize report routes", "error", err)
	}
}

func (rr *reportRouter) Init(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		// Public routes - no authentication required
		r.Get("/", rr.reportController.GetAllReports)
		r.Get("/{id}", rr.reportController.GetReportByID)
		r.Get("/business/{businessName}", rr.reportController.GetBusinessReports)
		
		// Protected routes - authentication required
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth)
			r.Post("/", rr.reportController.CreateReport)
			r.Put("/{id}", rr.reportController.UpdateReport)
			r.Delete("/{id}", rr.reportController.DeleteReport)
			r.Get("/user/me", rr.reportController.GetUserReports)
		})
		
		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdmin)
			r.Get("/status/{status}", rr.reportController.GetReportsByStatus)
			r.Post("/{id}/approve", rr.reportController.ApproveReport)
			r.Post("/{id}/reject", rr.reportController.RejectReport)
			r.Post("/{id}/flag", rr.reportController.FlagReport)
		})
	})
}
