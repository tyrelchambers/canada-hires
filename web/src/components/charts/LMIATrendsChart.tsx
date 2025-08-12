import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart'
import { LMIAStatistics } from '@/hooks/useLMIATrends'
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts'

interface LMIATrendsChartProps {
  data: LMIAStatistics[]
  title?: string
  description?: string
  className?: string
}

const chartConfig = {
  total_jobs: {
    label: 'Total Jobs',
    color: 'var(--chart-1)',
  },
  unique_employers: {
    label: 'Unique Employers',
    color: 'var(--chart-2)',
  },
} satisfies ChartConfig

export function LMIATrendsChart({ data, title = "Job Trends", description, className }: LMIATrendsChartProps) {
  // Transform data for chart
  const chartData = data.map((item) => ({
    date: new Date(item.date).toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric',
      year: item.period_type === 'monthly' ? 'numeric' : undefined
    }),
    total_jobs: item.total_jobs,
    unique_employers: item.unique_employers,
    avg_salary_min: item.avg_salary_min,
    avg_salary_max: item.avg_salary_max,
  }))

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          <AreaChart data={chartData}>
            <defs>
              <linearGradient id="fillTotalJobs" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-total_jobs)"
                  stopOpacity={0.8}
                />
                <stop
                  offset="95%"
                  stopColor="var(--color-total_jobs)"
                  stopOpacity={0.1}
                />
              </linearGradient>
              <linearGradient id="fillUniqueEmployers" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-unique_employers)"
                  stopOpacity={0.8}
                />
                <stop
                  offset="95%"
                  stopColor="var(--color-unique_employers)"
                  stopOpacity={0.1}
                />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={(value) => value}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              domain={[0, 'dataMax']}
            />
            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  labelFormatter={(value) => value}
                  indicator="dot"
                />
              }
            />
            <Area
              dataKey="unique_employers"
              type="natural"
              fill="url(#fillUniqueEmployers)"
              stroke="var(--color-unique_employers)"
              stackId="a"
            />
            <Area
              dataKey="total_jobs"
              type="natural"
              fill="url(#fillTotalJobs)"
              stroke="var(--color-total_jobs)"
              stackId="a"
            />
            <ChartLegend content={<ChartLegendContent />} />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}