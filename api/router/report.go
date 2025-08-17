package router

import (
	"canada-hires/container"
	"canada-hires/controllers"
	"canada-hires/middleware"
	"net/http"

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
	authMW           func(http.Handler) http.Handler
}

func NewReportRouter(cn *container.Container, reportController controllers.ReportController, authMW func(http.Handler) http.Handler) ReportRouter {
	return &reportRouter{
		cn:               cn,
		reportController: reportController,
		authMW:           authMW,
	}
}

func (rr *reportRouter) InjectReportRoutes(r chi.Router) {
	err := rr.cn.Invoke(func(reportController controllers.ReportController, authMW func(http.Handler) http.Handler) {
		rr := NewReportRouter(rr.cn, reportController, authMW)
		rr.Init(r)
	})

	if err != nil {
		log.Fatal("Failed to initialize report routes", "error", err)
	}
}

func (rr *reportRouter) Init(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		// Public routes - no authentication required
		r.Get("/", rr.reportController.GetReports)
		r.Get("/{id}", rr.reportController.GetReportByID)
		r.Get("/business/{businessName}", rr.reportController.GetReports)
		r.Get("/address", rr.reportController.GetReports)
		r.Get("/grouped-by-address", rr.reportController.GetReportsGrouped)
		
		// Protected routes - authentication required
		r.Group(func(r chi.Router) {
			// Apply auth middleware to extract user from cookie
			r.Use(rr.authMW)
			// Apply authentication requirement middleware
			r.Use(middleware.RequireAuth)
			r.Post("/", rr.reportController.CreateReport)
			r.Put("/{id}", rr.reportController.UpdateReport)
			r.Delete("/{id}", rr.reportController.DeleteReport)
			r.Get("/user/me", rr.reportController.GetUserReports)
		})
		
		// Admin-only routes
		r.Group(func(r chi.Router) {
			// Apply auth middleware to extract user from cookie
			r.Use(rr.authMW)
			// Apply admin requirement middleware
			r.Use(middleware.RequireAdmin)
			r.Get("/status/{status}", rr.reportController.GetReports)
			r.Post("/{id}/approve", rr.reportController.ApproveReport)
			r.Post("/{id}/reject", rr.reportController.RejectReport)
			r.Post("/{id}/flag", rr.reportController.FlagReport)
		})
	})
}
