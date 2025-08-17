import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useApiClient } from "@/hooks/useApiClient";
import { useMutation } from "@tanstack/react-query";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlay, faChartLine, faSpinner } from "@fortawesome/free-solid-svg-icons";

export function ScraperManager() {
  const apiClient = useApiClient();
  const [scraperStatus, setScraperStatus] = useState<string | null>(null);
  const [statisticsStatus, setStatisticsStatus] = useState<string | null>(null);

  const scraperMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.post("/admin/scraper/run");
      return response.data;
    },
    onSuccess: (data) => {
      setScraperStatus("Scraper job started successfully!");
      console.log("Scraper triggered:", data);
    },
    onError: (error: any) => {
      setScraperStatus(`Error: ${error.response?.data?.message || error.message}`);
      console.error("Scraper error:", error);
    },
  });

  const statisticsMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.post("/admin/scraper/statistics");
      return response.data;
    },
    onSuccess: (data) => {
      setStatisticsStatus("Statistics aggregation started successfully!");
      console.log("Statistics triggered:", data);
    },
    onError: (error: any) => {
      setStatisticsStatus(`Error: ${error.response?.data?.message || error.message}`);
      console.error("Statistics error:", error);
    },
  });

  const lmiaBackfillMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.post("/lmia/statistics/backfill");
      return response.data;
    },
    onSuccess: (data) => {
      setStatisticsStatus("LMIA statistics backfill completed successfully!");
      console.log("LMIA backfill triggered:", data);
    },
    onError: (error: any) => {
      setStatisticsStatus(`Error: ${error.response?.data?.message || error.message}`);
      console.error("LMIA backfill error:", error);
    },
  });

  const handleTriggerScraper = () => {
    setScraperStatus(null);
    scraperMutation.mutate();
  };

  const handleTriggerStatistics = () => {
    setStatisticsStatus(null);
    statisticsMutation.mutate();
  };

  const handleTriggerLMIABackfill = () => {
    setStatisticsStatus(null);
    lmiaBackfillMutation.mutate();
  };

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h3 className="font-bold text-xl">Scraper Management</h3>
        <p className="text-muted-foreground">
          Manually trigger job scraping and statistics aggregation
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {/* Job Scraper */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon icon={faPlay} className="h-5 w-5 text-blue-600" />
              Job Scraper
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Triggers the job scraper to collect new LMIA job postings from government sources.
              This process also automatically runs statistics aggregation when complete.
            </p>
            
            <Button
              onClick={handleTriggerScraper}
              disabled={scraperMutation.isPending}
              className="w-full"
            >
              {scraperMutation.isPending ? (
                <>
                  <FontAwesomeIcon icon={faSpinner} className="mr-2 h-4 w-4 animate-spin" />
                  Starting Scraper...
                </>
              ) : (
                <>
                  <FontAwesomeIcon icon={faPlay} className="mr-2 h-4 w-4" />
                  Run Job Scraper
                </>
              )}
            </Button>

            {scraperStatus && (
              <div className={`p-3 rounded-md text-sm ${
                scraperStatus.startsWith("Error") 
                  ? "bg-red-50 text-red-700 border border-red-200"
                  : "bg-green-50 text-green-700 border border-green-200"
              }`}>
                {scraperStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Statistics Aggregation */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon icon={faChartLine} className="h-5 w-5 text-green-600" />
              Statistics Aggregation
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Triggers the statistics aggregation to generate daily and monthly LMIA job trends
              from existing job posting data.
            </p>
            
            <Button
              onClick={handleTriggerStatistics}
              disabled={statisticsMutation.isPending}
              variant="outline"
              className="w-full"
            >
              {statisticsMutation.isPending ? (
                <>
                  <FontAwesomeIcon icon={faSpinner} className="mr-2 h-4 w-4 animate-spin" />
                  Running Aggregation...
                </>
              ) : (
                <>
                  <FontAwesomeIcon icon={faChartLine} className="mr-2 h-4 w-4" />
                  Run Statistics Aggregation
                </>
              )}
            </Button>

            {statisticsStatus && (
              <div className={`p-3 rounded-md text-sm ${
                statisticsStatus.startsWith("Error") 
                  ? "bg-red-50 text-red-700 border border-red-200"
                  : "bg-green-50 text-green-700 border border-green-200"
              }`}>
                {statisticsStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* LMIA Statistics Backfill */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon icon={faChartLine} className="h-5 w-5 text-purple-600" />
              LMIA Statistics Backfill
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Generates all historical LMIA statistics from existing job posting data. 
              This populates province/city charts and trend data for the trends page.
            </p>
            
            <Button
              onClick={handleTriggerLMIABackfill}
              disabled={lmiaBackfillMutation.isPending}
              variant="secondary"
              className="w-full"
            >
              {lmiaBackfillMutation.isPending ? (
                <>
                  <FontAwesomeIcon icon={faSpinner} className="mr-2 h-4 w-4 animate-spin" />
                  Running Backfill...
                </>
              ) : (
                <>
                  <FontAwesomeIcon icon={faChartLine} className="mr-2 h-4 w-4" />
                  Backfill LMIA Statistics
                </>
              )}
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Information Section */}
      <Card>
        <CardHeader>
          <CardTitle>Automated Schedule</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3 text-sm text-muted-foreground">
            <p>
              <strong>Daily Schedule:</strong> The job scraper runs automatically every day at midnight UTC.
            </p>
            <p>
              <strong>Auto-Aggregation:</strong> Statistics are automatically generated after each successful scraper run.
            </p>
            <p>
              <strong>Manual Triggers:</strong> Use the buttons above for immediate execution when needed for testing or updates.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}