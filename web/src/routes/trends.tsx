import { createFileRoute } from "@tanstack/react-router";
import {
  useDailyTrends,
  useMonthlyTrends,
  useTrendsSummary,
} from "@/hooks/useLMIATrends";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
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
    meta: [
      {
        title: "LMIA Job Trends - JobWatch Canada",
        description:
          "Track Temporary Foreign Worker (TFW) job posting trends across Canada with interactive charts and statistics.",
      },
    ],
  }),
});

function TrendsPage() {
  const [timeRange, setTimeRange] = useState<
    "week" | "month" | "quarter" | "year"
  >("month");
  const [chartType, setChartType] = useState<"daily" | "monthly">("daily");

  const trendsSummary = useTrendsSummary();

  // Get date ranges based on selected time range
  const getDateRange = () => {
    const now = new Date();
    const endDate = now.toISOString().split("T")[0];
    let startDate: Date;

    switch (timeRange) {
      case "week":
        startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      case "month":
        startDate = new Date(
          now.getFullYear(),
          now.getMonth() - 1,
          now.getDate(),
        );
        break;
      case "quarter":
        startDate = new Date(
          now.getFullYear(),
          now.getMonth() - 3,
          now.getDate(),
        );
        break;
      case "year":
        startDate = new Date(
          now.getFullYear() - 1,
          now.getMonth(),
          now.getDate(),
        );
        break;
    }

    return {
      start_date: startDate.toISOString().split("T")[0],
      end_date: endDate,
    };
  };

  const dateRange = getDateRange();
  const dailyTrends = useDailyTrends(
    chartType === "daily" ? dateRange : undefined,
  );
  const monthlyTrends = useMonthlyTrends(
    chartType === "monthly" ? dateRange : undefined,
  );

  const currentTrends = chartType === "daily" ? dailyTrends : monthlyTrends;

  return (
    <section>
      <AuthNav />
      <div className="container mx-auto py-6 space-y-6">
        {/* Header */}
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">LMIA Job Trends</h1>
          <p className="text-muted-foreground">
            Track Temporary Foreign Worker (TFW) job posting trends across
            Canada
          </p>
        </div>

        {/* Summary Stats */}
        {trendsSummary.data && (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
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
          <div className="flex gap-2">
            <Select
              value={timeRange}
              onValueChange={(value: "week" | "month" | "quarter" | "year") =>
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

            <Select
              value={chartType}
              onValueChange={(value: "daily" | "monthly") =>
                setChartType(value)
              }
            >
              <SelectTrigger className="w-[100px]">
                <SelectValue placeholder="Select period" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="daily">Daily</SelectItem>
                <SelectItem value="monthly">Monthly</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Button
            variant="outline"
            onClick={() => {
              void dailyTrends.refetch();
              void monthlyTrends.refetch();
              void trendsSummary.refetch();
            }}
          >
            Refresh Data
          </Button>
        </div>

        {/* Job Trends Section */}
        <div className="space-y-4">
          <div>
            <h2 className="text-2xl font-semibold tracking-tight">
              Job Trends -{" "}
              {timeRange === "week"
                ? "Past Week"
                : timeRange === "month"
                  ? "Past Month"
                  : timeRange === "quarter"
                    ? "Past Quarter"
                    : "Past Year"}
            </h2>
            <p className="text-sm text-muted-foreground">
              {chartType === "daily" ? "Daily" : "Monthly"} LMIA job postings
              and employer trends over time
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
              title={`${chartType === "daily" ? "Daily" : "Monthly"} Job Trends`}
              description={`${chartType === "daily" ? "Daily" : "Monthly"} LMIA job postings and employer trends over time`}
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
              Top provinces and cities with the most LMIA job postings today
            </p>
          </div>

          {trendsSummary.isLoading ? (
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
          ) : trendsSummary.error ? (
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-center h-[400px] text-red-600">
                  Error loading regional data: {trendsSummary.error.message}
                </div>
              </CardContent>
            </Card>
          ) : (
            <div className="grid md:grid-cols-2 gap-4">
              <LMIARegionalChart
                data={trendsSummary.data?.top_provinces_today || []}
                title="Top Provinces Today"
                description="Provinces with the most LMIA job postings today"
                type="province"
              />
              <LMIARegionalChart
                data={trendsSummary.data?.top_cities_today || []}
                title="Top Cities Today"
                description="Cities with the most LMIA job postings today"
                type="city"
              />
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
