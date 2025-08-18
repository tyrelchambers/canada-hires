import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSearch,
  faMapMarkerAlt,
  faFilter,
} from "@fortawesome/free-solid-svg-icons";
import { useReportsGroupedByAddress } from "@/hooks/useReports";
import { AuthNav } from "@/components/AuthNav";
import { StripedBackground } from "@/components/StripedBackground";
import { BusinessCard } from "@/components/BusinessCard";

function ReportsPage() {
  // Search and filter state
  const [searchQuery, setSearchQuery] = useState("");
  const [cityFilter, setCityFilter] = useState("");
  const [provinceFilter, setProvinceFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [yearFilter, setYearFilter] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [limit] = useState(20);

  // Always use grouped reports with filters
  const {
    data: groupedData,
    isLoading: loading,
    error,
  } = useReportsGroupedByAddress({
    query: searchQuery,
    city: cityFilter,
    province: provinceFilter,
    status: statusFilter,
    year: yearFilter,
    limit,
    offset: (currentPage - 1) * limit,
  });

  const businesses = groupedData?.data || [];
  const total = groupedData?.count || 0;

  // Generate year options (current year back to 2015)
  const currentYear = new Date().getFullYear();
  const yearOptions = Array.from(
    { length: currentYear - 2014 },
    (_, i) => currentYear - i,
  );

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
    <section className="bg-gray-50 min-h-screen">
      <AuthNav />

      {/* Hero Search Section */}
      <div className="bg-gradient-to-r from-gray-900 to-gray-700 text-white py-12 relative">
        <StripedBackground />
        <div className="container mx-auto px-4 relative z-10">
          <div className="max-w-4xl mx-auto text-center">
            <h1 className="text-4xl font-bold mb-4">
              Business Reports Directory
            </h1>
            <p className="text-xl mb-8 text-white">
              Discover hiring practices and make informed decisions about
              Canadian businesses
            </p>

            {/* Main Search Bar */}
            <div className="bg-white rounded-lg p-4 flex flex-col md:flex-row gap-4 shadow-lg">
              <div className="flex-1 relative">
                <FontAwesomeIcon
                  icon={faSearch}
                  className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                />
                <Input
                  placeholder="coffee, restaurants, retail..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 text-lg py-3 border-0 focus:ring-2 focus:ring-red-500"
                />
              </div>
              <div className="flex-1 relative">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                />
                <Input
                  placeholder="Toronto, ON"
                  value={cityFilter}
                  onChange={(e) => setCityFilter(e.target.value)}
                  className="pl-10 text-lg py-3 border-0 focus:ring-2 focus:ring-red-500"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col lg:flex-row gap-8">
          {/* Sidebar Filters */}
          <div className="lg:w-1/4">
            <Card className="sticky top-4">
              <CardContent className="p-6">
                <h3 className="font-bold text-lg mb-4 flex items-center">
                  <FontAwesomeIcon icon={faFilter} className="mr-2" />
                  Filters
                </h3>

                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium mb-2">
                      Province
                    </label>
                    <Input
                      placeholder="All Provinces"
                      value={provinceFilter}
                      onChange={(e) => setProvinceFilter(e.target.value)}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">
                      Verification Status
                    </label>
                    <select
                      value={statusFilter}
                      onChange={(e) => setStatusFilter(e.target.value)}
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                    >
                      <option value="">All Statuses</option>
                      <option value="approved">Verified</option>
                      <option value="pending">Under Review</option>
                      <option value="rejected">Unverified</option>
                      <option value="flagged">Flagged</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">
                      Report Year
                    </label>
                    <select
                      value={yearFilter}
                      onChange={(e) => setYearFilter(e.target.value)}
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                    >
                      <option value="">All Years</option>
                      {yearOptions.map((year) => (
                        <option key={year} value={year.toString()}>
                          {year}
                        </option>
                      ))}
                    </select>
                  </div>

                  <Button
                    onClick={clearFilters}
                    variant="outline"
                    className="w-full"
                  >
                    Clear All Filters
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Content */}
          <div className="lg:w-3/4">
            {/* Results Header */}
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold">Business Reports</h2>
              <Link to="/reports/create">
                <Button>Submit your report</Button>
              </Link>
            </div>

            {/* Loading State */}
            {loading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-red-600"></div>
                <p className="mt-4 text-gray-600 text-lg">
                  Finding businesses...
                </p>
              </div>
            )}

            {/* Error State */}
            {error && (
              <Card className="border-red-200 bg-red-50">
                <CardContent className="p-6">
                  <p className="text-red-800">
                    Error:{" "}
                    {error instanceof Error
                      ? error.message
                      : "An error occurred"}
                  </p>
                </CardContent>
              </Card>
            )}

            {/* Business Cards */}
            {!loading && !error && (
              <>
                <div className="space-y-4 mb-8">
                  {businesses.map((business, index) => (
                    <BusinessCard
                      key={`${business.business_address}-${index}`}
                      business={business}
                    />
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="flex justify-center items-center space-x-4 bg-white p-6 rounded-lg border">
                    <Button
                      onClick={() =>
                        setCurrentPage(Math.max(1, currentPage - 1))
                      }
                      disabled={currentPage === 1}
                      variant="outline"
                    >
                      Previous
                    </Button>

                    <div className="flex items-center space-x-2">
                      {Array.from(
                        { length: Math.min(5, totalPages) },
                        (_, i) => {
                          const pageNum = i + 1;
                          return (
                            <Button
                              key={pageNum}
                              onClick={() => setCurrentPage(pageNum)}
                              variant={
                                currentPage === pageNum ? "default" : "outline"
                              }
                              size="sm"
                            >
                              {pageNum}
                            </Button>
                          );
                        },
                      )}
                      {totalPages > 5 && (
                        <span className="text-gray-500">...</span>
                      )}
                    </div>

                    <Button
                      onClick={() =>
                        setCurrentPage(Math.min(totalPages, currentPage + 1))
                      }
                      disabled={currentPage === totalPages}
                      variant="outline"
                    >
                      Next
                    </Button>
                  </div>
                )}
              </>
            )}

            {/* Empty State */}
            {!loading && !error && businesses.length === 0 && (
              <Card className="text-center py-12">
                <CardContent>
                  <FontAwesomeIcon
                    icon={faSearch}
                    className="text-gray-300 text-6xl mb-4"
                  />
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    No businesses found
                  </h3>
                  <p className="text-gray-500 mb-4">
                    Try adjusting your search criteria or location
                  </p>
                  <Button onClick={clearFilters} variant="outline">
                    Clear all filters
                  </Button>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}

export const Route = createFileRoute("/reports/")({
  component: ReportsPage,
});
