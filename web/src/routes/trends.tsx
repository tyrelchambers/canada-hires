import { createFileRoute } from "@tanstack/react-router";
import {
  useTrendsSummary,
  useRegionalStats,
  useTrendsForPeriod,
} from "@/hooks/useLMIATrends";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { LMIATrendsChart } from "@/components/charts/LMIATrendsChart";
import { LMIARegionalChart } from "@/components/charts/LMIARegionalChart";
import { useState } from "react";
import Stat from "@/components/Stat";
import {
  faArrowDown,
  faArrowUp,
  faBuilding,
  faUsers,
} from "@fortawesome/free-solid-svg-icons";
import { AuthNav } from "@/components/AuthNav";

export const Route = createFileRoute("/trends")({
  component: TrendsPage,
  head: () => ({
    title: "LMIA Job Trends - JobWatch Canada",
    meta: [
      {
        name: "description",
        content:
          "Track Temporary Foreign Worker (TFW) job posting trends across Canada with interactive charts and statistics.",
      },
    ],
  }),
});

type TimeRange = "week" | "month" | "quarter" | "year";

const timeRangeMap: Record<TimeRange, string> = {
  week: "Past Week",
  month: "Past Month",
  quarter: "Past Quarter",
  year: "Past Year",
};

function TrendsPage() {
  const [timeRange, setTimeRange] = useState<TimeRange>("month");

  const trendsSummary = useTrendsSummary();
  const regionalStats = useRegionalStats(timeRange);
  const currentTrends = useTrendsForPeriod(timeRange);

  const chartTitle = timeRangeMap[timeRange];

  return (
    <section>
      <AuthNav />
      <div className="container mx-auto py-6 space-y-6 px-4">
        {/* Header */}
        <div className="space-y-2 my-10">
          <h1 className="text-3xl font-bold tracking-tight">LMIA Job Trends</h1>
          <p className="text-muted-foreground">
            Track Temporary Foreign Worker (TFW) job posting trends across
            Canada
          </p>
        </div>

        {/* Summary Stats */}
        {trendsSummary.data && (
          <div className="grid gap-8 grid-cols-2 lg:grid-cols-4 my-10">
            <Stat
              label="Jobs Today"
              value={String(trendsSummary.data.total_jobs_today)}
              icon={faBuilding}
            />
            <Stat
              label="Jobs This Month"
              value={String(trendsSummary.data.total_jobs_this_month)}
              icon={faUsers}
            />
            <Stat
              label="Month-over-Month"
              value={`${trendsSummary.data.percentage_change > 0 ? "+" : ""}${trendsSummary.data.percentage_change.toFixed(1)}%`}
              icon={
                trendsSummary.data.percentage_change >= 0
                  ? faArrowUp
                  : faArrowDown
              }
            />
            <Stat
              label="Last Month"
              value={String(trendsSummary.data.total_jobs_last_month)}
              icon={faBuilding}
            />
          </div>
        )}

        {/* Controls */}
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
          <Select
            value={timeRange}
            onValueChange={(value: TimeRange) =>
              setTimeRange(value)
            }
          >
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="Select time range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="week">Past Week</SelectItem>
              <SelectItem value="month">Past Month</SelectItem>
              <SelectItem value="quarter">Past Quarter</SelectItem>
              <SelectItem value="year">Past Year</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Job Trends Section */}
        <div className="space-y-4">
          <div>
            <h2 className="text-2xl font-semibold tracking-tight">
              Job Trends - {chartTitle}
            </h2>
            <p className="text-sm text-muted-foreground">
              Daily LMIA job postings and employer trends over time
            </p>
          </div>

          {currentTrends.isLoading ? (
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-center h-[400px]">
                  Loading trends data...
                </div>
              </CardContent>
            </Card>
          ) : currentTrends.error ? (
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-center h-[400px] text-red-600">
                  Error loading trends data: {currentTrends.error.message}
                </div>
              </CardContent>
            </Card>
          ) : (
            <LMIATrendsChart
              data={currentTrends.data?.data || []}
              title="Daily Job Trends"
              description="Daily LMIA job postings and employer trends over time"
            />
          )}
        </div>

        {/* Regional Breakdown Section */}
        <div className="space-y-4">
          <div>
            <h2 className="text-2xl font-semibold tracking-tight">
              Regional Breakdown
            </h2>
            <p className="text-sm text-muted-foreground">
              Top provinces and cities with the most LMIA job postings in the{" "}
              {timeRangeMap[timeRange].toLowerCase()}
            </p>
          </div>

          {regionalStats.isLoading ? (
            <div className="grid md:grid-cols-2 gap-4">
              <Card>
                <CardContent className="p-6">
                  <div className="flex items-center justify-center h-[300px]">
                    Loading regional data...
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardContent className="p-6">
                  <div className="flex items-center justify-center h-[300px]">
                    Loading regional data...
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : regionalStats.error ? (
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-center h-[400px] text-red-600">
                  Error loading regional data: {regionalStats.error.message}
                </div>
              </CardContent>
            </Card>
          ) : (
            <div className="grid md:grid-cols-2 gap-4">
              <LMIARegionalChart
                data={regionalStats.data?.top_provinces || []}
                title={`Top Provinces - ${chartTitle}`}
                description={`Provinces with the most LMIA job postings in the ${timeRangeMap[timeRange].toLowerCase()}`}
                type="province"
              />
              <LMIARegionalChart
                data={regionalStats.data?.top_cities || []}
                title={`Top Cities - ${chartTitle}`}
                description={`Cities with the most LMIA job postings in the ${timeRangeMap[timeRange].toLowerCase()}`}
                type="city"
              />
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
