package main

import (
	"canada-hires/container"
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/services"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	// Define command line flags
	var (
		dryRun    = flag.Bool("dry-run", false, "Preview the post without actually submitting to Reddit")
		force     = flag.Bool("force", false, "Post even if job is already marked as posted to Reddit")
		subreddit = flag.String("subreddit", "", "Override subreddit for posting (takes precedence over env var)")
	)
	flag.Parse()

	// Check for job posting ID argument
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <job_posting_id>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	jobPostingID := flag.Args()[0]

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn("Could not load .env file", "error", err)
	}

	// Create container and get required services
	cn, err := container.New()
	if err != nil {
		log.Fatal("Failed to create container", "error", err)
		os.Exit(1)
	}

	var (
		jobRepo       repos.JobBankRepository
		redditService services.RedditService
	)

	if err := cn.Invoke(func(jr repos.JobBankRepository, rs services.RedditService) {
		jobRepo = jr
		redditService = rs
	}); err != nil {
		log.Fatal("Failed to get services from container", "error", err)
		os.Exit(1)
	}

	// Retrieve the job posting
	log.Info("Retrieving job posting", "job_id", jobPostingID)
	jobPosting, err := jobRepo.GetJobPostingByID(jobPostingID)
	if err != nil {
		// Try by JobBankID if ID lookup fails
		log.Debug("Job not found by ID, trying by JobBankID", "job_id", jobPostingID)
		jobPosting, err = jobRepo.GetJobPostingByJobBankID(jobPostingID)
		if err != nil {
			log.Fatal("Job posting not found", "job_id", jobPostingID, "error", err)
			os.Exit(1)
		}
	}

	// Check if already posted (unless force flag is used)
	log.Debug("Checking Reddit posted status", 
		"job_id", jobPostingID, 
		"reddit_posted", jobPosting.RedditPosted, 
		"force_flag", *force)
		
	if jobPosting.RedditPosted && !*force {
		log.Warn("Job already posted to Reddit", "job_id", jobPostingID)
		fmt.Printf("Job '%s' is already marked as posted to Reddit.\n", jobPosting.Title)
		fmt.Println("Use --force flag to post anyway.")
		os.Exit(1)
	}
	
	if jobPosting.RedditPosted && *force {
		log.Info("Force flag enabled, proceeding with repost", "job_id", jobPostingID)
		fmt.Printf("⚠️  Force flag enabled - reposting job '%s'\n", jobPosting.Title)
	}

	// Handle subreddit override
	config := redditService.GetDefaultConfig()
	originalSubreddit := config.Subreddit

	// Priority: CLI flag > REDDIT_SUBREDDIT env var > default config
	if *subreddit != "" {
		config.Subreddit = *subreddit
		log.Info("Using subreddit from CLI flag", "subreddit", *subreddit)
	} else if redditSubreddit := os.Getenv("REDDIT_SUBREDDIT"); redditSubreddit != "" {
		config.Subreddit = redditSubreddit
		log.Info("Using subreddit from REDDIT_SUBREDDIT environment variable", "subreddit", redditSubreddit)
	}

	// Generate post data
	postData := config.GeneratePostData(jobPosting)
	if postData == nil {
		log.Fatal("Failed to generate Reddit post data")
		os.Exit(1)
	}

	// Display preview
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("REDDIT POST PREVIEW\n")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Subreddit: r/%s\n", postData.Subreddit)
	fmt.Printf("Title: %s\n", postData.Title)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Body:\n%s\n", postData.Body)
	fmt.Println(strings.Repeat("=", 60))

	if *dryRun {
		fmt.Println("\nDRY RUN MODE - Post preview only, not submitting to Reddit")
		fmt.Printf("Job Details: %s at %s (%s)\n", jobPosting.Title, jobPosting.Employer, jobPosting.Location)
		os.Exit(0)
	}

	// Confirm posting
	fmt.Print("\nProceed with posting to Reddit? (y/N): ")
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("Posting cancelled")
		os.Exit(0)
	}

	// Create temporary subreddit object for the posting
	tempSubreddit := &models.Subreddit{
		Name: postData.Subreddit,
	}

	// Post to Reddit
	log.Info("Posting job to Reddit", "job_id", jobPostingID, "subreddit", postData.Subreddit)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = redditService.PostJobWithConfig(ctx, jobPosting, config, tempSubreddit)
	if err != nil {
		log.Error("Failed to post job to Reddit", "error", err, "job_id", jobPostingID)
		fmt.Printf("Error posting to Reddit: %v\n", err)
		os.Exit(1)
	}

	// Restore original subreddit in config if we changed it
	config.Subreddit = originalSubreddit

	fmt.Printf("\n✅ Successfully posted job to r/%s!\n", postData.Subreddit)
	fmt.Printf("Job Title: %s\n", jobPosting.Title)
	fmt.Printf("Employer: %s\n", jobPosting.Employer)
	fmt.Printf("Job ID: %s\n", jobPostingID)

	log.Info("Reddit posting completed successfully", "job_id", jobPostingID)
}