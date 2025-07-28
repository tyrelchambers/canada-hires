# Reddit Post CLI Tool

A command-line tool for manually posting job postings to Reddit for testing purposes.

## Usage

```bash
# Basic usage - post a job to Reddit
go run cmd/reddit_post/main.go <job_posting_id>

# Preview the post without actually submitting (dry run)
go run cmd/reddit_post/main.go --dry-run <job_posting_id>

# Force post even if job is already marked as posted
go run cmd/reddit_post/main.go --force <job_posting_id>

# Override the subreddit for testing
go run cmd/reddit_post/main.go --subreddit testsubreddit <job_posting_id>
```

## Environment Variables

### Required
- `REDDIT_ID`: Reddit app client ID (required)
- `REDDIT_SECRET`: Reddit app client secret (required) 
- `REDDIT_USERNAME`: Reddit username (must be registered as app developer) (required)
- `REDDIT_PASSWORD`: Reddit password (required)

### Optional
- `REDDIT_TEST_SUBREDDIT`: Override the default subreddit for testing
- `REDDIT_USER_AGENT`: Reddit API user agent (defaults to "JobWatchCanada/1.0")

## Features

- **Job Lookup**: Accepts either job posting UUID or JobBankID
- **Preview Mode**: `--dry-run` shows what would be posted without actually submitting
- **Force Mode**: `--force` allows reposting jobs already marked as posted
- **Subreddit Override**: CLI flag or environment variable to test in different subreddits
- **Automatic Authentication**: Fetches and caches Reddit access tokens automatically
- **Token Management**: Automatically refreshes expired tokens
- **Automatic Status Update**: Marks jobs as posted in database when successful
- **Interactive Confirmation**: Asks for confirmation before posting (unless dry-run)

## Example Output

```
============================================================
REDDIT POST PREVIEW
============================================================
Subreddit: r/testjobs
Title: 🇨🇦 New TFW Job: Software Developer at Tech Corp - Toronto, ON
------------------------------------------------------------
Body:
**New Temporary Foreign Worker (TFW) Job Posting**

**Position:** Software Developer
**Employer:** Tech Corp
**Location:** Toronto, ON
**Salary:** $70,000 - $80,000 yearly
**Posted:** January 15, 2025

**Job Details:** [View on Job Bank](https://www.jobbank.gc.ca/jobposting/12345)

---
*This posting was automatically detected from Government of Canada Job Bank TFW listings. Data provided by JobWatch Canada for transparency in hiring practices.*
============================================================

Proceed with posting to Reddit? (y/N): 
```

## Error Handling

- Job not found: Returns exit code 1
- Already posted (without --force): Returns exit code 1  
- Reddit API errors: Returns exit code 1
- Missing Reddit credentials: Returns exit code 1

## Reddit App Setup

Before using this tool, you need to create a Reddit app:

1. Go to https://www.reddit.com/prefs/apps/
2. Click "Create App" or "Create Another App"
3. Choose "script" as the app type
4. Fill in name and description
5. Set redirect URI (not used for script apps, can be http://localhost)
6. Note the client ID (under the app name) and client secret
7. Ensure your Reddit username is registered as a developer for this app

## Notes

- The tool automatically tries to find jobs by UUID first, then by JobBankID
- Posts are marked as `reddit_posted = true` in the database upon successful submission
- Reddit API rate limits apply - use responsibly for testing
- Access tokens are automatically fetched and cached (expire after ~1 hour)
- Ensure your Reddit app has proper permissions for posting to target subreddits