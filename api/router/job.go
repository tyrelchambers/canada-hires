package router

import (
	"canada-hires/controllers"

	"github.com/go-chi/chi/v5"
)

func JobRoutes(r chi.Router, jobController *controllers.JobController) {
	r.Route("/api/jobs", func(r chi.Router) {
		// Job postings endpoints
		r.Get("/", jobController.GetJobPostings)
		r.Get("/stats", jobController.GetJobStats)
		
		// Scraping endpoints
		r.Post("/scraping-runs", jobController.CreateScrapingRun)
		r.Get("/scraping-runs", jobController.GetScrapingRuns)
		r.Post("/scraping-runs/{scraping_run_id}/jobs", jobController.SubmitScraperJobs)
		r.Post("/scraping-runs/{scraping_run_id}/complete", jobController.CompleteScrapingRun)
	})

}