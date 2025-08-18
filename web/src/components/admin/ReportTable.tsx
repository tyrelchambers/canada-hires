import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useReports, useDeleteReport, type Report } from "@/hooks/useReports";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { 
  faSearch, 
  faEdit, 
  faTrash, 
  faCalendar,
  faUser,
  faMapMarkerAlt,
  faEye
} from "@fortawesome/free-solid-svg-icons";

interface ReportTableProps {
  onEditReport: (report: Report) => void;
  onViewReport: (report: Report) => void;
}

export function ReportTable({ onEditReport, onViewReport }: ReportTableProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [cityFilter, setCityFilter] = useState("");
  const [provinceFilter, setProvinceFilter] = useState("");
  const [yearFilter, setYearFilter] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);

  const filters = {
    query: searchQuery,
    city: cityFilter,
    province: provinceFilter,
    year: yearFilter,
    limit: pageSize,
    offset: (currentPage - 1) * pageSize,
  };

  const { data: reportsData, isLoading, error } = useReports(filters);
  const deleteReportMutation = useDeleteReport();

  const reports = reportsData?.reports || [];
  const totalPages = Math.ceil((reportsData?.pagination.total || 0) / pageSize);

  const handleDeleteReport = async (reportId: string) => {
    if (window.confirm("Are you sure you want to delete this report? This action cannot be undone.")) {
      try {
        await deleteReportMutation.mutateAsync(reportId);
      } catch (error) {
        console.error("Failed to delete report:", error);
        alert("Failed to delete report. Please try again.");
      }
    }
  };

  const clearFilters = () => {
    setSearchQuery("");
    setCityFilter("");
    setProvinceFilter("");
    setYearFilter("");
    setCurrentPage(1);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-CA', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getReportSourceBadge = (source: string) => {
    const variants = {
      employment: "bg-blue-100 text-blue-800",
      observation: "bg-green-100 text-green-800",
      public_record: "bg-purple-100 text-purple-800"
    };
    
    const labels = {
      employment: "Employment",
      observation: "Observation", 
      public_record: "Public Record"
    };

    return (
      <Badge className={variants[source as keyof typeof variants] || "bg-gray-100 text-gray-800"}>
        {labels[source as keyof typeof labels] || source}
      </Badge>
    );
  };

  // Generate year options (current year back to 2015)
  const currentYear = new Date().getFullYear();
  const yearOptions = Array.from(
    { length: currentYear - 2014 },
    (_, i) => currentYear - i
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Report Management</CardTitle>
        
        {/* Search and Filter Bar */}
        <div className="flex flex-col lg:flex-row gap-4 mt-4">
          <div className="flex-1 relative">
            <FontAwesomeIcon
              icon={faSearch}
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
            />
            <Input
              placeholder="Search business name or address..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
          
          <div className="flex gap-2">
            <Input
              placeholder="City"
              value={cityFilter}
              onChange={(e) => setCityFilter(e.target.value)}
              className="w-32"
            />
            <Input
              placeholder="Province"
              value={provinceFilter}
              onChange={(e) => setProvinceFilter(e.target.value)}
              className="w-32"
            />
            <select
              value={yearFilter}
              onChange={(e) => setYearFilter(e.target.value)}
              className="flex h-10 w-32 rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="">All Years</option>
              {yearOptions.map((year) => (
                <option key={year} value={year.toString()}>
                  {year}
                </option>
              ))}
            </select>
            <Button onClick={clearFilters} variant="outline">
              Clear
            </Button>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        {/* Loading State */}
        {isLoading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="mt-2 text-gray-600">Loading reports...</p>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="text-center py-12">
            <p className="text-red-600">
              Error loading reports: {error instanceof Error ? error.message : "Unknown error"}
            </p>
          </div>
        )}

        {/* Reports Table */}
        {!isLoading && !error && (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left p-4 font-semibold">Business</th>
                    <th className="text-left p-4 font-semibold">Address</th>
                    <th className="text-left p-4 font-semibold">Source</th>
                    <th className="text-left p-4 font-semibold">Confidence</th>
                    <th className="text-left p-4 font-semibold">Submitted</th>
                    <th className="text-left p-4 font-semibold">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {reports.map((report) => (
                    <tr key={report.id} className="border-b hover:bg-gray-50">
                      <td className="p-4">
                        <div className="font-medium">{report.business_name}</div>
                        <div className="text-sm text-gray-500 flex items-center">
                          <FontAwesomeIcon icon={faUser} className="w-3 h-3 mr-1" />
                          User ID: {report.user_id.slice(0, 8)}...
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex items-center text-gray-600">
                          <FontAwesomeIcon icon={faMapMarkerAlt} className="w-3 h-3 mr-1" />
                          {report.business_address}
                        </div>
                      </td>
                      <td className="p-4">
                        {getReportSourceBadge(report.report_source)}
                      </td>
                      <td className="p-4">
                        {report.confidence_level ? (
                          <span className="font-medium">{report.confidence_level}/10</span>
                        ) : (
                          <span className="text-gray-400">-</span>
                        )}
                      </td>
                      <td className="p-4">
                        <div className="flex items-center text-gray-600">
                          <FontAwesomeIcon icon={faCalendar} className="w-3 h-3 mr-1" />
                          {formatDate(report.created_at)}
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="flex gap-2">
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => onViewReport(report)}
                          >
                            <FontAwesomeIcon icon={faEye} className="w-3 h-3" />
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => onEditReport(report)}
                          >
                            <FontAwesomeIcon icon={faEdit} className="w-3 h-3" />
                          </Button>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDeleteReport(report.id)}
                            disabled={deleteReportMutation.isPending}
                          >
                            <FontAwesomeIcon icon={faTrash} className="w-3 h-3" />
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Empty State */}
            {reports.length === 0 && (
              <div className="text-center py-12">
                <FontAwesomeIcon
                  icon={faSearch}
                  className="text-gray-300 text-6xl mb-4"
                />
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  No reports found
                </h3>
                <p className="text-gray-500 mb-4">
                  Try adjusting your search criteria
                </p>
                <Button onClick={clearFilters} variant="outline">
                  Clear all filters
                </Button>
              </div>
            )}

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center space-x-4 mt-6 pt-6 border-t">
                <Button
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                  variant="outline"
                >
                  Previous
                </Button>

                <div className="flex items-center space-x-2">
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    const pageNum = i + 1;
                    return (
                      <Button
                        key={pageNum}
                        onClick={() => setCurrentPage(pageNum)}
                        variant={currentPage === pageNum ? "default" : "outline"}
                        size="sm"
                      >
                        {pageNum}
                      </Button>
                    );
                  })}
                  {totalPages > 5 && <span className="text-gray-500">...</span>}
                </div>

                <Button
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                  variant="outline"
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}