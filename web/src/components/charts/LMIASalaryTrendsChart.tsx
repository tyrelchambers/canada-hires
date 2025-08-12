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
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart'
import { LMIAStatistics } from '@/hooks/useLMIATrends'
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts'

interface LMIASalaryTrendsChartProps {
  data: LMIAStatistics[]
  title?: string
  description?: string
  className?: string
}

const chartConfig = {
  avg_salary_min: {
    label: 'Min Salary',
    color: 'hsl(var(--chart-1))',
  },
  avg_salary_max: {
    label: 'Max Salary', 
    color: 'hsl(var(--chart-2))',
  },
} satisfies ChartConfig

export function LMIASalaryTrendsChart({ 
  data, 
  title = "Salary Trends", 
  description,
  className 
}: LMIASalaryTrendsChartProps) {
  // Filter out entries without salary data and transform for chart
  const chartData = data
    .filter(item => item.avg_salary_min || item.avg_salary_max)
    .map((item) => ({
      date: new Date(item.date).toLocaleDateString('en-US', { 
        month: 'short', 
        day: 'numeric',
        year: item.period_type === 'monthly' ? 'numeric' : undefined
      }),
      avg_salary_min: item.avg_salary_min ? Math.round(item.avg_salary_min) : null,
      avg_salary_max: item.avg_salary_max ? Math.round(item.avg_salary_max) : null,
    }))

  // Show message if no salary data available
  if (chartData.length === 0) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
          {description && <CardDescription>{description}</CardDescription>}
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-[200px] text-muted-foreground">
            No salary data available for the selected period
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <AreaChart
            accessibilityLayer
            data={chartData}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => value}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => `$${value.toLocaleString()}`}
            />
            <ChartTooltip 
              cursor={false} 
              content={<ChartTooltipContent 
                formatter={(value, name) => [
                  `$${Number(value).toLocaleString()}`,
                  name === 'avg_salary_min' ? 'Min Salary' : 'Max Salary'
                ]}
              />} 
            />
            <Area
              dataKey="avg_salary_min"
              type="monotone"
              fill="var(--color-avg_salary_min)"
              fillOpacity={0.4}
              stroke="var(--color-avg_salary_min)"
              strokeWidth={2}
              stackId="a"
            />
            <Area
              dataKey="avg_salary_max"
              type="monotone"
              fill="var(--color-avg_salary_max)"
              fillOpacity={0.4}
              stroke="var(--color-avg_salary_max)"
              strokeWidth={2}
              stackId="a"
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}