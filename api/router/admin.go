package router

import (
	"canada-hires/container"
	"canada-hires/controllers"
	"canada-hires/middleware"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type AdminRouter interface {
	Init(r chi.Router)
	InjectAdminRoutes(r chi.Router)
}

type adminRouter struct {
	cn            *container.Container
	jobController *controllers.JobController
	authMW        func(http.Handler) http.Handler
}

func NewAdminRouter(cn *container.Container, jobController *controllers.JobController, authMW func(http.Handler) http.Handler) AdminRouter {
	return &adminRouter{
		cn:            cn,
		jobController: jobController,
		authMW:        authMW,
	}
}

func (ar *adminRouter) InjectAdminRoutes(r chi.Router) {
	err := ar.cn.Invoke(func(jobController *controllers.JobController, authMW func(http.Handler) http.Handler) {
		ar := NewAdminRouter(ar.cn, jobController, authMW)
		ar.Init(r)
	})

	if err != nil {
		log.Fatal("Failed to initialize admin routes", "error", err)
	}
}

func (ar *adminRouter) Init(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		// Apply auth middleware to extract user from cookie
		r.Use(ar.authMW)
		// Apply admin middleware to require admin role
		r.Use(middleware.RequireAdmin)

		// Reddit approval endpoints for jobs
		r.Route("/jobs/reddit", func(r chi.Router) {
			r.Get("/pending", ar.jobController.GetPendingJobsForReddit)
			r.Get("/posted", ar.jobController.GetPostedJobsForReddit)
			r.Post("/approve/{job_id}", ar.jobController.ApproveJobForReddit)
			r.Post("/reject/{job_id}", ar.jobController.RejectJobForReddit)
			r.Post("/bulk-approve", ar.jobController.BulkApproveJobsForReddit)
			r.Post("/bulk-reject", ar.jobController.BulkRejectJobsForReddit)
			
			// Content generation endpoints
			r.Post("/generate-content/{job_id}", ar.jobController.GenerateRedditPostContent)
			r.Post("/generate-content/bulk", ar.jobController.GenerateBulkRedditPostContent)
			
			// Preview endpoint to see what will be posted
			r.Get("/preview/{job_id}", ar.jobController.PreviewRedditPost)
		})

		// Scraper and statistics endpoints
		r.Route("/scraper", func(r chi.Router) {
			r.Post("/run", ar.jobController.TriggerScraper)
			r.Post("/statistics", ar.jobController.TriggerStatisticsAggregation)
		})
	})
}