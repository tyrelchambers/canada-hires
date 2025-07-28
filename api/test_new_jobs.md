# Testing New vs Existing Job Detection

## How it Works

The system now tracks which jobs are truly "new" vs updates to existing jobs:

### 1. **Pre-insertion Check**
- Before inserting jobs, we query the database to see which URLs already exist
- This gives us a map of existing URLs: `existingUrls[url] = true`

### 2. **New Job Identification**
- For each job being processed, if `!existingUrls[posting.URL]`, it's a new job
- We track the IDs of these new jobs in `newJobIds[]`

### 3. **Database Operation**
- Insert/update all jobs using `ON CONFLICT (url) DO UPDATE`
- Existing jobs get updated, new jobs get inserted

### 4. **Return Only New Jobs**
- Query the database to fetch only the jobs with IDs in `newJobIds[]`
- Return these new jobs to the service layer

### 5. **Reddit Posting**
- Only the returned new jobs get posted to Reddit
- Existing job updates are ignored for Reddit posting

## Testing Scenarios

### First Scrape (All New)
```
Input: 100 jobs from scraper
Database: Empty
Result: 100 new jobs → 100 Reddit posts
```

### Second Scrape (Some Duplicates)
```
Input: 100 jobs from scraper (80 duplicates, 20 new)
Database: 80 existing jobs
Result: 20 new jobs → 20 Reddit posts
```

### Third Scrape (All Duplicates)
```
Input: 100 jobs from scraper (all duplicates)
Database: 100 existing jobs  
Result: 0 new jobs → 0 Reddit posts
Log: "No new jobs to post to Reddit - all jobs were updates to existing postings"
```

## Benefits

1. **No Spam**: Reddit only gets posts for genuinely new job listings
2. **Accurate Tracking**: Clear distinction between new jobs and job updates
3. **Efficient**: Single query to check existing URLs before processing
4. **Logging**: Clear visibility into new vs updated job counts
5. **Scalable**: Works with batched processing for large scrape runs