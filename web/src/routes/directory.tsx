import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSearch,
  faMapMarkerAlt,
  faFilter,
} from "@fortawesome/free-solid-svg-icons";
import { useReports } from "@/hooks/useReports";
import { AuthNav } from "@/components/AuthNav";

const getStatusColor = (status: string) => {
  switch (status) {
    case "approved":
      return "bg-green-100 text-green-800 border-green-200";
    case "pending":
      return "bg-yellow-100 text-yellow-800 border-yellow-200";
    case "rejected":
      return "bg-red-100 text-red-800 border-red-200";
    case "flagged":
      return "bg-orange-100 text-orange-800 border-orange-200";
    default:
      return "bg-gray-100 text-gray-800 border-gray-200";
  }
};

const getStatusLabel = (status: string) => {
  switch (status) {
    case "approved":
      return "✅ Approved";
    case "pending":
      return "⏳ Pending";
    case "rejected":
      return "❌ Rejected";
    case "flagged":
      return "🚩 Flagged";
    default:
      return "❓ Unknown";
  }
};

function DirectoryPage() {
  // Search and filter state
  const [searchQuery, setSearchQuery] = useState("");
  const [cityFilter, setCityFilter] = useState("");
  const [provinceFilter, setProvinceFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [yearFilter, setYearFilter] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [limit] = useState(20);

  // Use the reports hook
  const {
    data: reportsData,
    isLoading: loading,
    error,
  } = useReports({
    query: searchQuery,
    city: cityFilter,
    province: provinceFilter,
    status: statusFilter,
    year: yearFilter,
    limit,
    offset: (currentPage - 1) * limit,
  });

  const reports = reportsData?.reports || [];
  const total = reportsData?.pagination?.total || 0;

  // Generate year options (current year back to 2015)
  const currentYear = new Date().getFullYear();
  const yearOptions = Array.from(
    { length: currentYear - 2014 },
    (_, i) => currentYear - i,
  );

  const handleSearch = () => {
    setCurrentPage(1);
  };

  const clearFilters = () => {
    setSearchQuery("");
    setCityFilter("");
    setProvinceFilter("");
    setStatusFilter("");
    setYearFilter("");
    setCurrentPage(1);
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <section>
      <AuthNav />
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Reports Directory</h1>
          <p className="text-gray-600">
            Community-submitted reports about business hiring practices in
            Canada
          </p>
        </div>

        {/* Search and Filters */}
        <div className="bg-white p-6 rounded-lg shadow-sm border mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-4">
            <div className="relative">
              <FontAwesomeIcon
                icon={faSearch}
                className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              />
              <Input
                placeholder="Search businesses..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>

            <Input
              placeholder="City"
              value={cityFilter}
              onChange={(e) => setCityFilter(e.target.value)}
            />

            <Input
              placeholder="Province"
              value={provinceFilter}
              onChange={(e) => setProvinceFilter(e.target.value)}
            />

            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="">All Statuses</option>
              <option value="approved">✅ Approved</option>
              <option value="pending">⏳ Pending</option>
              <option value="rejected">❌ Rejected</option>
              <option value="flagged">🚩 Flagged</option>
            </select>

            <select
              value={yearFilter}
              onChange={(e) => setYearFilter(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="">All Years</option>
              {yearOptions.map((year) => (
                <option key={year} value={year.toString()}>
                  {year}
                </option>
              ))}
            </select>
          </div>

          <div className="flex gap-2">
            <Button onClick={handleSearch} size="sm">
              <FontAwesomeIcon icon={faSearch} className="mr-2" />
              Search
            </Button>
            <Button onClick={clearFilters} variant="outline" size="sm">
              <FontAwesomeIcon icon={faFilter} className="mr-2" />
              Clear Filters
            </Button>
          </div>
        </div>

        {/* Results */}
        {loading && (
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="mt-2 text-gray-600">Loading reports...</p>
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <p className="text-red-800">
              Error:{" "}
              {error instanceof Error ? error.message : "An error occurred"}
            </p>
          </div>
        )}

        {!loading && !error && (
          <>
            <div className="mb-4 text-sm text-gray-600">
              Showing {reports.length} of {total} reports
            </div>

            <div className="bg-white rounded-lg shadow-sm border mb-8">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Business Name</TableHead>
                    <TableHead>Address</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Confidence</TableHead>
                    <TableHead>Source</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Notes</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {reports.map((report) => (
                    <TableRow key={report.id}>
                      <TableCell className="font-medium">
                        {report.business_name}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center text-sm">
                          <FontAwesomeIcon
                            icon={faMapMarkerAlt}
                            className="mr-1 text-gray-400"
                          />
                          <span className="truncate max-w-xs">
                            {report.business_address}
                          </span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge
                          className={getStatusColor(report.status)}
                          variant="outline"
                        >
                          {getStatusLabel(report.status)}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {report.confidence_level ? (
                          <span className="font-medium">
                            {report.confidence_level}/10
                          </span>
                        ) : (
                          <span className="text-gray-400">-</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <span className="capitalize">
                          {report.report_source.replace('_', ' ')}
                        </span>
                      </TableCell>
                      <TableCell>
                        <span className="text-sm text-gray-600">
                          {new Date(report.created_at).toLocaleDateString()}
                        </span>
                      </TableCell>
                      <TableCell>
                        {report.additional_notes ? (
                          <span className="text-sm text-gray-600 truncate max-w-xs block">
                            {report.additional_notes}
                          </span>
                        ) : (
                          <span className="text-gray-400">-</span>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center space-x-2">
                <Button
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                  variant="outline"
                  size="sm"
                >
                  Previous
                </Button>

                <span className="text-sm text-gray-600">
                  Page {currentPage} of {totalPages}
                </span>

                <Button
                  onClick={() =>
                    setCurrentPage(Math.min(totalPages, currentPage + 1))
                  }
                  disabled={currentPage === totalPages}
                  variant="outline"
                  size="sm"
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}

        {!loading && !error && reports.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">
              No reports found matching your criteria.
            </p>
            <Button onClick={clearFilters} className="mt-4" variant="outline">
              Clear filters to see all reports
            </Button>
          </div>
        )}
      </div>
    </section>
  );
}

export const Route = createFileRoute("/directory")({
  component: DirectoryPage,
});
