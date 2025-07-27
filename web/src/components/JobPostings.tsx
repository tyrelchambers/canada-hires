import { useState, useMemo, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Pagination } from "@/components/ui/pagination";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSearch,
  faMapMarkerAlt,
  faExclamationCircle,
  faBuilding,
  faCalendar,
  faExternalLinkAlt,
  faFlag,
  faInfoCircle,
  faExclamationTriangle,
} from "@fortawesome/free-solid-svg-icons";
import { useJobPostings, useJobStats } from "@/hooks/useJobPostings";
import { JobPostingFilters } from "@/types";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import Stat from "./Stat";

export function JobPostings() {
  const [searchQuery, setSearchQuery] = useState("");
  const [employer, setEmployer] = useState("");
  const [city, setCity] = useState("");
  const [province, setProvince] = useState("");
  const [title, setTitle] = useState("");
  const [salaryMin, setSalaryMin] = useState("");
  const [sortBy, setSortBy] = useState("posting_date");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 25;

  // Build filters object
  const filters: JobPostingFilters = useMemo(() => {
    const f: JobPostingFilters = {
      limit: itemsPerPage,
      offset: (currentPage - 1) * itemsPerPage,
      sort_by: sortBy,
      sort_order: sortOrder,
      days: 30, // Default to last 30 days
    };

    if (searchQuery.trim()) f.search = searchQuery.trim();
    if (employer.trim()) f.employer = employer.trim();
    if (city.trim()) f.city = city.trim();
    if (province.trim()) f.province = province.trim();
    if (title.trim()) f.title = title.trim();
    if (salaryMin.trim()) {
      const parsed = parseFloat(salaryMin);
      if (!isNaN(parsed)) f.salary_min = parsed;
    }

    return f;
  }, [
    searchQuery,
    employer,
    city,
    province,
    title,
    salaryMin,
    sortBy,
    sortOrder,
    currentPage,
  ]);

  const { data: jobData, isLoading, error } = useJobPostings(filters);

  const { data: statsData } = useJobStats();

  const totalPages = jobData ? Math.ceil(jobData.total / itemsPerPage) : 0;

  // Reset to first page when search parameters change
  useEffect(() => {
    setCurrentPage(1);
  }, [
    searchQuery,
    employer,
    city,
    province,
    title,
    salaryMin,
    sortBy,
    sortOrder,
  ]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    // The query will automatically trigger due to dependency changes
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleClearFilters = () => {
    setSearchQuery("");
    setEmployer("");
    setCity("");
    setProvince("");
    setTitle("");
    setSalaryMin("");
    setSortBy("posting_date");
    setSortOrder("desc");
    setCurrentPage(1);
  };

  const formatSalary = (min?: number, max?: number, type?: string) => {
    if (!min && !max) return null;

    const formatAmount = (amount: number) => {
      return new Intl.NumberFormat("en-CA", {
        style: "currency",
        currency: "CAD",
        minimumFractionDigits: 0,
        maximumFractionDigits: 2,
      }).format(amount);
    };

    let salaryText = "";
    if (min && max && min !== max) {
      salaryText = `${formatAmount(min)} - ${formatAmount(max)}`;
    } else if (min) {
      salaryText = formatAmount(min);
    } else if (max) {
      salaryText = formatAmount(max);
    }

    if (type && salaryText) {
      salaryText += ` ${type}`;
    }

    return salaryText;
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "N/A";
    return new Date(dateString).toLocaleDateString("en-CA");
  };

  const canadianProvinces = [
    "AB",
    "BC",
    "MB",
    "NB",
    "NL",
    "NT",
    "NS",
    "NU",
    "ON",
    "PE",
    "QC",
    "SK",
    "YT",
  ];

  return (
    <div className="space-y-6">
      {/* Educational Content About LMIA */}
      <div className="border border-border p-6 rounded-xl ">
        <p>
          <FontAwesomeIcon icon={faInfoCircle} /> About Labour Market Impact
          Assessment (LMIA) Job Postings
        </p>
        <div className="space-y-4 mt-6">
          <p className="text-sm mb-3">
            <strong>What is LMIA?</strong> A Labour Market Impact Assessment
            (LMIA) is a document that an employer in Canada may need to get
            before hiring a foreign worker. It shows that there is a need for a
            foreign worker to fill the job and that no Canadian worker or
            permanent resident is available to do the job.
          </p>
          <p className="text-sm mb-3">
            <strong>Canadian workers are encouraged to apply!</strong> Before
            hiring through the Temporary Foreign Worker (TFW) Program, employers
            must demonstrate they couldn't find qualified Canadian workers. If
            you're qualified for these positions, you should still apply even if
            they show LMIA approval.
          </p>
          <p className="text-sm">
            <strong>No response from employers?</strong> If you apply for these
            jobs but don't receive interviews or responses despite being
            qualified, this could indicate potential misuse of the TFW program.
            Report such cases to{" "}
            <a
              href="https://www.canada.ca/en/employment-social-development/services/foreign-workers/report-abuse.html"
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 underline"
            >
              the TFW program tip line
            </a>
            .
          </p>
        </div>
      </div>

      {/* Stats Card */}
      {statsData && (
        <div className="grid grid-cols-1 md:grid-cols-3  py-10">
          <Stat
            label="Total LMIA Job Postings"
            value={statsData.total_jobs.toLocaleString()}
          />
          <Stat
            label="Unique Employers"
            value={statsData.total_employers.toLocaleString()}
          />
          <Stat label="Default View Period" value="Last 30 Days" />
        </div>
      )}

      {/* Search and Filter Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FontAwesomeIcon icon={faSearch} className="w-5 h-5" />
            Search LMIA Job Postings
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSearch} className="space-y-4">
            {/* General Search */}
            <div>
              <Input
                placeholder="Search across job titles, employers, and locations..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full"
              />
            </div>

            {/* Specific Filters */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div>
                <Input
                  placeholder="Job Title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                />
              </div>
              <div>
                <Input
                  placeholder="Employer/Company"
                  value={employer}
                  onChange={(e) => setEmployer(e.target.value)}
                />
              </div>
              <div>
                <Input
                  placeholder="City"
                  value={city}
                  onChange={(e) => setCity(e.target.value)}
                />
              </div>
              <div>
                <select
                  value={province}
                  onChange={(e) => setProvince(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value="">All Provinces</option>
                  {canadianProvinces.map((prov) => (
                    <option key={prov} value={prov}>
                      {prov}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <Input
                  placeholder="Min Salary (CAD)"
                  type="number"
                  value={salaryMin}
                  onChange={(e) => setSalaryMin(e.target.value)}
                />
              </div>
              <div>
                <select
                  value={`${sortBy}-${sortOrder}`}
                  onChange={(e) => {
                    const [field, order] = e.target.value.split("-");
                    setSortBy(field);
                    setSortOrder(order as "asc" | "desc");
                  }}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value="posting_date-desc">Newest First</option>
                  <option value="posting_date-asc">Oldest First</option>
                  <option value="salary_min-desc">Highest Salary</option>
                  <option value="salary_min-asc">Lowest Salary</option>
                  <option value="employer-asc">Employer A-Z</option>
                  <option value="title-asc">Job Title A-Z</option>
                </select>
              </div>
            </div>

            <div className="flex gap-2">
              <Button type="submit">
                <FontAwesomeIcon icon={faSearch} className="w-4 h-4 mr-2" />
                Search
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={handleClearFilters}
              >
                Clear Filters
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Results */}
      {isLoading && (
        <Card>
          <CardContent className="py-8 text-center">
            <div className="flex items-center justify-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-900"></div>
              <span>Loading job postings...</span>
            </div>
          </CardContent>
        </Card>
      )}

      {error && (
        <Card>
          <CardContent className="py-8 text-center">
            <FontAwesomeIcon
              icon={faExclamationCircle}
              className="w-8 h-8 text-red-500 mx-auto mb-2"
            />
            <p className="text-red-600">
              Error loading job postings. Please try again.
            </p>
          </CardContent>
        </Card>
      )}

      {jobData && jobData.jobs && jobData.jobs.length > 0 && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">
              Job Postings ({jobData.total.toLocaleString()} found)
            </h3>
            <div className="text-sm text-gray-600">
              Showing {(currentPage - 1) * itemsPerPage + 1} to{" "}
              {Math.min(currentPage * itemsPerPage, jobData.total)} of{" "}
              {jobData.total.toLocaleString()}
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-sm border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Job Title</TableHead>
                  <TableHead>Employer</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead>Salary</TableHead>
                  <TableHead>Posted Date</TableHead>
                  <TableHead>LMIA Status</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {jobData.jobs.map((job) => (
                  <TableRow key={job.id}>
                    <TableCell className="font-medium">
                      <div className="max-w-xs">
                        <div
                          className="font-semibold text-gray-900 truncate"
                          title={job.title}
                        >
                          {job.title}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="max-w-xs">
                        <div className="flex items-center">
                          <FontAwesomeIcon
                            icon={faBuilding}
                            className="mr-2 text-gray-400 w-3 h-3"
                          />
                          <span className="truncate" title={job.employer}>
                            {job.employer}
                          </span>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center">
                        <FontAwesomeIcon
                          icon={faMapMarkerAlt}
                          className="mr-2 text-gray-400 w-3 h-3"
                        />
                        <span
                          className="text-sm truncate max-w-32"
                          title={job.location}
                        >
                          {job.location}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      {formatSalary(
                        job.salary_min,
                        job.salary_max,
                        job.salary_type,
                      ) ? (
                        <div className="flex items-center">
                          <span className="text-sm font-medium">
                            {formatSalary(
                              job.salary_min,
                              job.salary_max,
                              job.salary_type,
                            )}
                          </span>
                        </div>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center">
                        <FontAwesomeIcon
                          icon={faCalendar}
                          className="mr-1 text-gray-400 w-3 h-3"
                        />
                        <span className="text-sm">
                          {formatDate(job.posting_date)}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Badge
                          variant="outline"
                          className="bg-yellow-50 text-yellow-700 border-yellow-200"
                        >
                          <FontAwesomeIcon
                            icon={faFlag}
                            className="w-3 h-3 mr-1"
                          />
                          Pending LMIA
                        </Badge>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button variant="outline" size="sm" asChild>
                          <a
                            href={job.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-1"
                          >
                            <FontAwesomeIcon
                              icon={faExternalLinkAlt}
                              className="w-3 h-3"
                            />
                            Apply
                          </a>
                        </Button>
                        {job.job_bank_id && (
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="outline" size="sm" asChild>
                                <a
                                  href={`https://www.jobbank.gc.ca/surveyreportmisuse/${job.job_bank_id}`}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="flex items-center gap-1"
                                >
                                  <FontAwesomeIcon
                                    icon={faExclamationTriangle}
                                    className="w-3 h-3"
                                  />
                                  Report
                                </a>
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              Report potential misuse of this job posting to the
                              Government of Canada
                            </TooltipContent>
                          </Tooltip>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              totalItems={jobData.total}
              itemsPerPage={itemsPerPage}
              onPageChange={handlePageChange}
              className="mt-6"
            />
          )}
        </div>
      )}

      {jobData && jobData.jobs && jobData.jobs.length === 0 && !isLoading && (
        <Card>
          <CardContent className="py-8 text-center">
            <FontAwesomeIcon
              icon={faSearch}
              className="w-8 h-8 text-gray-400 mx-auto mb-2"
            />
            <p className="text-gray-600">
              No job postings found matching your criteria.
            </p>
            <p className="text-sm text-gray-500 mt-1">
              Try adjusting your search filters or check back later for new
              postings.
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
