package services

import (
	"canada-hires/models"
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
}

type GeneratedRedditPost struct {
	JobID   string `json:"job_id"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

type BulkGenerationResponse struct {
	Posts []GeneratedRedditPost `json:"posts"`
}

// NewGeminiService creates a new Gemini service
func NewGeminiService() (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiService{
		client: client,
	}, nil
}

// Close closes the Gemini client (new SDK doesn't require explicit close)
func (g *GeminiService) Close() error {
	// The new Google Gen AI SDK doesn't require explicit closing
	return nil
}

// GenerateRedditPost generates sarcastic Reddit post content for a single job
func (g *GeminiService) GenerateRedditPost(ctx context.Context, job models.JobPosting) (string, error) {
	prompt := g.buildPrompt(job)

	// Set a reasonable timeout for generation
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := g.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	content := strings.TrimSpace(result.Text())
	if content == "" {
		return "", fmt.Errorf("generated content is empty")
	}

	return content, nil
}

// GenerateBulkRedditPosts generates Reddit post content for multiple jobs
func (g *GeminiService) GenerateBulkRedditPosts(ctx context.Context, jobs []models.JobPosting) BulkGenerationResponse {
	response := BulkGenerationResponse{
		Posts: make([]GeneratedRedditPost, 0, len(jobs)),
	}

	// Process jobs concurrently for better performance
	resultChan := make(chan GeneratedRedditPost, len(jobs))

	for _, job := range jobs {
		go func(j models.JobPosting) {
			content, err := g.GenerateRedditPost(ctx, j)
			post := GeneratedRedditPost{
				JobID:   j.ID,
				Content: content,
			}
			if err != nil {
				post.Error = err.Error()
			}
			resultChan <- post
		}(job)
	}

	// Collect results
	for i := 0; i < len(jobs); i++ {
		post := <-resultChan
		response.Posts = append(response.Posts, post)
	}

	return response
}

// buildPrompt creates the prompt for Gemini to generate sarcastic Reddit content
func (g *GeminiService) buildPrompt(job models.JobPosting) string {
	// Format salary information
	salaryInfo := "Salary not specified"
	if job.SalaryMin != nil && job.SalaryMax != nil {
		if *job.SalaryMin == *job.SalaryMax {
			salaryInfo = fmt.Sprintf("$%.2f %s", *job.SalaryMin, g.getSalaryType(job.SalaryType))
		} else {
			salaryInfo = fmt.Sprintf("$%.2f - $%.2f %s", *job.SalaryMin, *job.SalaryMax, g.getSalaryType(job.SalaryType))
		}
	} else if job.SalaryRaw != nil {
		salaryInfo = *job.SalaryRaw
	}

	// Format posting date
	postingDate := "Recently posted"
	if job.PostingDate != nil {
		postingDate = job.PostingDate.Format("January 2, 2006")
	}

	prompt := fmt.Sprintf(`You are a witty, sarcastic Canadian Reddit user who cares about fair employment practices and transparency around the Temporary Foreign Worker (TFW) program. Your job is to create engaging Reddit posts that highlight TFW job postings in a way that informs the community while maintaining a slightly sarcastic but not mean-spirited tone.

Create a Reddit post about this TFW job posting:

**Job Details:**
- Title: %s
- Employer: %s
- Location: %s
- Salary: %s
- Posted: %s
- Job URL: %s

**Instructions:**
1. Write in a conversational, slightly sarcastic tone that a typical Canadian Redditor would use
2. Highlight the fact that this is a TFW position requiring an LMIA
3. Include the key job details (title, employer, location, salary)
4. Keep it informative but engaging - people should learn something
5. Add some Canadian humor or references where appropriate
6. Stay factual - don't make assumptions beyond what's provided
7. Keep the post between 100-200 words
8. End with exactly this footer (no changes): "---\nLearn more at https://jobwatchcanada.com | Join the discussion at r/jobwatchcanada"
9. Don't make the post sound stupid and unnatural. This is an example of a stupid sentence - "Well, butter my bannock and call me Sheila, another day, another LMIA job posting has popped up on the ol' Job Bank!"
10. Keep the post very low-key and natural, business casual. You are not pretending to be a user, don't try to speak like a human. But just create a reddit post to convey the provided information. Imagine you're a CEO releasing a very professional and succinct newsletter.
11. Do NOT include an preamble such as "Alright team, another one for the Job Bank files, just landed in the queue." OR any other wording.
12. Bold the job details

Don't sound like you are actual a human. For example:

"As with other TFW positions, this role necessitates an approved Labour Market Impact Assessment (LMIA). It consistently raises questions when a position like this, at that compensation, appears to be unfillable by the local workforce, especially within a well-populated area like the Fraser Valley. It makes one wonder about the nature of the "labour shortage" being addressed. Just something to consider as these postings continue to surface."

This is too much. It is good to factual point out how companies are exploiting this in a sentence or two, but don't sound like you're trying to relate to people.

Make it sound natural, like something a real person would post on Reddit, not like corporate marketing copy or cheesy AI slop.`,
		job.Title,
		job.Employer,
		job.Location,
		salaryInfo,
		postingDate,
		job.URL,
	)

	return prompt
}

// getSalaryType returns a readable salary type
func (g *GeminiService) getSalaryType(salaryType *string) string {
	if salaryType == nil {
		return "per hour"
	}

	switch *salaryType {
	case "yearly":
		return "per year"
	case "monthly":
		return "per month"
	case "weekly":
		return "per week"
	case "biweekly":
		return "biweekly"
	case "hourly":
		return "per hour"
	default:
		return "per hour"
	}
}
