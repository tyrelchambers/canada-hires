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
import { RegionData } from '@/hooks/useLMIATrends'
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from 'recharts'

interface LMIARegionalChartProps {
  data: RegionData[]
  title?: string
  description?: string
  className?: string
  type?: 'province' | 'city'
}

const chartConfig = {
  count: {
    label: 'Jobs',
    color: 'var(--chart-1)',
  },
} satisfies ChartConfig

export function LMIARegionalChart({ 
  data, 
  title = "Regional Distribution", 
  description,
  className,
  type = 'province'
}: LMIARegionalChartProps) {
  // Limit to top 10 and sort by count
  const chartData = data
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
    .map((item) => ({
      name: item.name,
      count: item.count,
    }))

  // Show message if no data available
  if (!data || data.length === 0) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
          {description && <CardDescription>{description}</CardDescription>}
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-[300px] text-muted-foreground">
            No {type} data available
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
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer 
          config={chartConfig}
          className="aspect-auto h-[300px] w-full"
        >
          <BarChart
            accessibilityLayer
            data={chartData}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="name"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              angle={-45}
              textAnchor="end"
              height={80}
              minTickGap={32}
              tickFormatter={(value) => {
                // Truncate long city names
                if (type === 'city' && value.length > 12) {
                  return value.substring(0, 12) + '...'
                }
                return value
              }}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel indicator="dot" />}
            />
            <Bar
              dataKey="count"
              fill="var(--color-count)"
              radius={[4, 4, 0, 0]}
            />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}