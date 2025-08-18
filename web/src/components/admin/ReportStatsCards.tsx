import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useReports } from "@/hooks/useReports";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { 
  faFileAlt, 
  faCalendarWeek, 
  faCalendarDay,
  faBuilding
} from "@fortawesome/free-solid-svg-icons";

export function ReportStatsCards() {
  // Get total reports count
  const { data: allReports, isLoading: totalLoading } = useReports({ limit: 1000 });
  
  // Get this week's reports
  const weekAgo = new Date();
  weekAgo.setDate(weekAgo.getDate() - 7);
  const { data: weekReports, isLoading: weekLoading } = useReports({ 
    limit: 1000 
  });
  
  // Get today's reports  
  const today = new Date();
  const { data: todayReports, isLoading: todayLoading } = useReports({ 
    limit: 1000
  });

  // Calculate stats
  const totalReports = allReports?.reports?.length || 0;
  
  const thisWeekCount = weekReports?.reports?.filter(report => {
    const reportDate = new Date(report.created_at);
    return reportDate >= weekAgo;
  }).length || 0;
  
  const todayCount = todayReports?.reports?.filter(report => {
    const reportDate = new Date(report.created_at);
    return reportDate.toDateString() === today.toDateString();
  }).length || 0;

  // Get unique businesses count
  const uniqueBusinesses = new Set(
    allReports?.reports?.map(report => report.business_name.toLowerCase()) || []
  ).size;

  const isLoading = totalLoading || weekLoading || todayLoading;

  const stats = [
    {
      title: "Total Reports",
      value: totalReports,
      icon: faFileAlt,
      color: "text-blue-600",
      bgColor: "bg-blue-50",
    },
    {
      title: "This Week", 
      value: thisWeekCount,
      icon: faCalendarWeek,
      color: "text-green-600",
      bgColor: "bg-green-50",
    },
    {
      title: "Today",
      value: todayCount,
      icon: faCalendarDay,
      color: "text-orange-600", 
      bgColor: "bg-orange-50",
    },
    {
      title: "Unique Businesses",
      value: uniqueBusinesses,
      icon: faBuilding,
      color: "text-purple-600",
      bgColor: "bg-purple-50",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      {stats.map((stat, index) => (
        <Card key={index}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-600">
              {stat.title}
            </CardTitle>
            <div className={`p-2 rounded-lg ${stat.bgColor}`}>
              <FontAwesomeIcon 
                icon={stat.icon} 
                className={`h-4 w-4 ${stat.color}`} 
              />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {isLoading ? (
                <div className="animate-pulse bg-gray-200 h-8 w-16 rounded"></div>
              ) : (
                stat.value.toLocaleString()
              )}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}