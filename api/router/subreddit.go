package router

import (
	"canada-hires/controllers"
	"canada-hires/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SubredditRoutes(subredditController *controllers.SubredditController, authMW func(http.Handler) http.Handler) func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/subreddits", func(r chi.Router) {
			// Public routes
			r.Get("/active", subredditController.GetActiveSubreddits)
			
			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(authMW)
				r.Use(middleware.RequireAdmin)
				
				r.Get("/", subredditController.GetSubreddits)
				r.Post("/", subredditController.CreateSubreddit)
				r.Put("/{subreddit_id}", subredditController.UpdateSubreddit)
				r.Delete("/{subreddit_id}", subredditController.DeleteSubreddit)
			})
		})
	}
}