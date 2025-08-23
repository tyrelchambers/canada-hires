import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  useLMIAOperations,
  useScraperOperations,
} from "@/hooks/useLMIAOperations";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlay,
  faChartLine,
  faSpinner,
  faDownload,
  faCogs,
  faMapMarkerAlt,
} from "@fortawesome/free-solid-svg-icons";

export function ScraperManager() {
  const lmiaOperations = useLMIAOperations();
  const scraperOperations = useScraperOperations();
  const [scraperStatus, setScraperStatus] = useState<string | null>(null);
  const [statisticsStatus, setStatisticsStatus] = useState<string | null>(null);
  const [lmiaScraperStatus, setLmiaScraperStatus] = useState<string | null>(
    null,
  );
  const [lmiaProcessorStatus, setLmiaProcessorStatus] = useState<string | null>(
    null,
  );
  const [geocodingStatus, setGeocodingStatus] = useState<string | null>(null);

  const handleTriggerScraper = () => {
    setScraperStatus(null);
    scraperOperations.scraper.mutate(undefined, {
      onSuccess: (data) => {
        setScraperStatus("Scraper job started successfully!");
        console.log("Scraper triggered:", data);
      },
      onError: (error: unknown) => {
        const errorMessage =
          error instanceof Error ? error.message : "Unknown error";
        setScraperStatus(`Error: ${errorMessage}`);
        console.error("Scraper error:", error);
      },
    });
  };

  const handleTriggerStatistics = () => {
    setStatisticsStatus(null);
    scraperOperations.statistics.mutate(undefined, {
      onSuccess: (data) => {
        setStatisticsStatus("Statistics aggregation started successfully!");
        console.log("Statistics triggered:", data);
      },
      onError: (error: unknown) => {
        const errorMessage =
          error instanceof Error ? error.message : "Unknown error";
        setStatisticsStatus(`Error: ${errorMessage}`);
        console.error("Statistics error:", error);
      },
    });
  };

  const handleTriggerLMIABackfill = () => {
    setStatisticsStatus(null);
    lmiaOperations.backfill.mutate(undefined, {
      onSuccess: (data) => {
        setStatisticsStatus("LMIA statistics backfill completed successfully!");
        console.log("LMIA backfill triggered:", data);
      },
      onError: (error: unknown) => {
        const errorMessage =
          error instanceof Error ? error.message : "Unknown error";
        setStatisticsStatus(`Error: ${errorMessage}`);
        console.error("LMIA backfill error:", error);
      },
    });
  };

  const handleTriggerLMIAScraper = () => {
    setLmiaScraperStatus(null);
    lmiaOperations.fullUpdate.mutate(undefined, {
      onSuccess: (data) => {
        setLmiaScraperStatus("LMIA data scraper started successfully!");
        console.log("LMIA scraper triggered:", data);
      },
      onError: (error: unknown) => {
        const errorMessage =
          error instanceof Error ? error.message : "Unknown error";
        setLmiaScraperStatus(`Error: ${errorMessage}`);
        console.error("LMIA scraper error:", error);
      },
    });
  };

  const handleTriggerLMIAProcessor = () => {
    setLmiaProcessorStatus(null);
    lmiaOperations.processor.mutate(undefined, {
      onSuccess: (data) => {
        setLmiaProcessorStatus(
          "LMIA resource processing started successfully!",
        );
        console.log("LMIA processor triggered:", data);
      },
      onError: (error: unknown) => {
        const errorMessage =
          error instanceof Error ? error.message : "Unknown error";
        setLmiaProcessorStatus(`Error: ${errorMessage}`);
        console.error("LMIA processor error:", error);
      },
    });
  };

  const handleTriggerGeocoding = () => {
    setGeocodingStatus(null);
    lmiaOperations.geocoding.mutate(undefined, {
      onSuccess: (data) => {
        setGeocodingStatus("LMIA geocoding started successfully!");
        console.log("LMIA geocoding triggered:", data);
      },
      onError: (error: unknown) => {
        const errorMessage =
          error instanceof Error ? error.message : "Unknown error";
        setGeocodingStatus(`Error: ${errorMessage}`);
        console.error("LMIA geocoding error:", error);
      },
    });
  };

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h3 className="font-bold text-xl">Scraper Management</h3>
        <p className="text-muted-foreground">
          Manually trigger job scraping and statistics aggregation
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {/* Job Scraper */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon
                icon={faPlay}
                className="h-5 w-5 text-blue-600"
              />
              Job Scraper
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Triggers the job scraper to collect new LMIA job postings from
              government sources. This process also automatically runs
              statistics aggregation when complete.
            </p>

            <Button
              onClick={handleTriggerScraper}
              disabled={scraperOperations.scraper.isPending}
              className="w-full"
            >
              {scraperOperations.scraper.isPending ? (
                <>
                  <FontAwesomeIcon
                    icon={faSpinner}
                    className="mr-2 h-4 w-4 animate-spin"
                  />
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
              <div
                className={`p-3 rounded-md text-sm ${
                  scraperStatus.startsWith("Error")
                    ? "bg-red-50 text-red-700 border border-red-200"
                    : "bg-green-50 text-green-700 border border-green-200"
                }`}
              >
                {scraperStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* LMIA Data Scraper */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon
                icon={faDownload}
                className="h-5 w-5 text-orange-600"
              />
              LMIA Data Scraper
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Downloads and processes LMIA employer data from Open Canada API.
              This fetches new CSV/Excel files and populates employer records.
            </p>

            <Button
              onClick={handleTriggerLMIAScraper}
              disabled={lmiaOperations.fullUpdate.isPending}
              variant="outline"
              className="w-full border-orange-200 text-orange-700 hover:bg-orange-50"
            >
              {lmiaOperations.fullUpdate.isPending ? (
                <>
                  <FontAwesomeIcon
                    icon={faSpinner}
                    className="mr-2 h-4 w-4 animate-spin"
                  />
                  Downloading Data...
                </>
              ) : (
                <>
                  <FontAwesomeIcon icon={faDownload} className="mr-2 h-4 w-4" />
                  Download LMIA Data
                </>
              )}
            </Button>

            {lmiaScraperStatus && (
              <div
                className={`p-3 rounded-md text-sm ${
                  lmiaScraperStatus.startsWith("Error")
                    ? "bg-red-50 text-red-700 border border-red-200"
                    : "bg-green-50 text-green-700 border border-green-200"
                }`}
              >
                {lmiaScraperStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* LMIA Resource Processor */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon
                icon={faCogs}
                className="h-5 w-5 text-teal-600"
              />
              LMIA Processor
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Processes existing LMIA resources to populate employer records.
              Use this to reprocess data without downloading new files.
            </p>

            <Button
              onClick={handleTriggerLMIAProcessor}
              disabled={lmiaOperations.processor.isPending}
              variant="outline"
              className="w-full border-teal-200 text-teal-700 hover:bg-teal-50"
            >
              {lmiaOperations.processor.isPending ? (
                <>
                  <FontAwesomeIcon
                    icon={faSpinner}
                    className="mr-2 h-4 w-4 animate-spin"
                  />
                  Processing Data...
                </>
              ) : (
                <>
                  <FontAwesomeIcon icon={faCogs} className="mr-2 h-4 w-4" />
                  Process LMIA Resources
                </>
              )}
            </Button>

            {lmiaProcessorStatus && (
              <div
                className={`p-3 rounded-md text-sm ${
                  lmiaProcessorStatus.startsWith("Error")
                    ? "bg-red-50 text-red-700 border border-red-200"
                    : "bg-green-50 text-green-700 border border-green-200"
                }`}
              >
                {lmiaProcessorStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* LMIA Geocoding */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon
                icon={faMapMarkerAlt}
                className="h-5 w-5 text-indigo-600"
              />
              LMIA Geocoding
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Geocodes postal codes for LMIA employers using concurrent Mapbox v6 batch processing. 
              16 parallel workers each process 1000 postal codes per batch for maximum efficiency.
            </p>

            <Button
              onClick={handleTriggerGeocoding}
              disabled={lmiaOperations.geocoding.isPending}
              variant="outline"
              className="w-full border-indigo-200 text-indigo-700 hover:bg-indigo-50"
            >
              {lmiaOperations.geocoding.isPending ? (
                <>
                  <FontAwesomeIcon
                    icon={faSpinner}
                    className="mr-2 h-4 w-4 animate-spin"
                  />
                  Geocoding Postal Codes...
                </>
              ) : (
                <>
                  <FontAwesomeIcon
                    icon={faMapMarkerAlt}
                    className="mr-2 h-4 w-4"
                  />
                  Geocode Postal Codes
                </>
              )}
            </Button>

            {geocodingStatus && (
              <div
                className={`p-3 rounded-md text-sm ${
                  geocodingStatus.startsWith("Error")
                    ? "bg-red-50 text-red-700 border border-red-200"
                    : "bg-green-50 text-green-700 border border-green-200"
                }`}
              >
                {geocodingStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Statistics Aggregation */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon
                icon={faChartLine}
                className="h-5 w-5 text-green-600"
              />
              Statistics Aggregation
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Triggers the statistics aggregation to generate daily and monthly
              LMIA job trends from existing job posting data.
            </p>

            <Button
              onClick={handleTriggerStatistics}
              disabled={scraperOperations.statistics.isPending}
              variant="outline"
              className="w-full"
            >
              {scraperOperations.statistics.isPending ? (
                <>
                  <FontAwesomeIcon
                    icon={faSpinner}
                    className="mr-2 h-4 w-4 animate-spin"
                  />
                  Running Aggregation...
                </>
              ) : (
                <>
                  <FontAwesomeIcon
                    icon={faChartLine}
                    className="mr-2 h-4 w-4"
                  />
                  Run Statistics Aggregation
                </>
              )}
            </Button>

            {statisticsStatus && (
              <div
                className={`p-3 rounded-md text-sm ${
                  statisticsStatus.startsWith("Error")
                    ? "bg-red-50 text-red-700 border border-red-200"
                    : "bg-green-50 text-green-700 border border-green-200"
                }`}
              >
                {statisticsStatus}
              </div>
            )}
          </CardContent>
        </Card>

        {/* LMIA Statistics Backfill */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon
                icon={faChartLine}
                className="h-5 w-5 text-purple-600"
              />
              LMIA Statistics Backfill
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Generates all historical LMIA statistics from existing job posting
              data. This populates province/city charts and trend data for the
              trends page.
            </p>

            <Button
              onClick={handleTriggerLMIABackfill}
              disabled={lmiaOperations.backfill.isPending}
              variant="secondary"
              className="w-full"
            >
              {lmiaOperations.backfill.isPending ? (
                <>
                  <FontAwesomeIcon
                    icon={faSpinner}
                    className="mr-2 h-4 w-4 animate-spin"
                  />
                  Running Backfill...
                </>
              ) : (
                <>
                  <FontAwesomeIcon
                    icon={faChartLine}
                    className="mr-2 h-4 w-4"
                  />
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
          <CardTitle>Process Overview</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3 text-sm text-muted-foreground">
            <p>
              <strong>Recommended Workflow:</strong>
            </p>
            <ol className="list-decimal list-inside space-y-1 ml-4">
              <li>Download LMIA Data - Fetches latest employer data files</li>
              <li>Process LMIA Resources - Parses files and extracts postal codes</li>
              <li>Geocode Postal Codes - Adds latitude/longitude for map features</li>
              <li>Backfill LMIA Statistics - Generates trend data and charts</li>
            </ol>
            <p>
              <strong>Daily Schedule:</strong> The job scraper runs
              automatically every day at midnight UTC.
            </p>
            <p>
              <strong>Auto-Aggregation:</strong> Statistics are automatically
              generated after each successful scraper run.
            </p>
            <p>
              <strong>Concurrent Batch Processing:</strong> Geocoding uses 16 parallel workers 
              with Mapbox v6 batch API, processing up to 16,000 postal codes simultaneously.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
