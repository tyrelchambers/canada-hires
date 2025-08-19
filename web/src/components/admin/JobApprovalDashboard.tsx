import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { User } from "@/types";
import {
  usePendingJobs,
  usePostedJobs,
} from "@/hooks/useAdminJobs";
import { PendingJobsList } from "./PendingJobsList";
import { PostedJobsList } from "./PostedJobsList";
import { SubredditManager } from "./SubredditManager";
import { ScraperManager } from "./ScraperManager";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";

interface JobApprovalDashboardProps {
  user: User;
  activeTab: string;
}

export function JobApprovalDashboard({
  user,
  activeTab,
}: JobApprovalDashboardProps) {
  const navigate = useNavigate();
  const [selectedJobs, setSelectedJobs] = useState<string[]>([]);
  const [currentPage, setCurrentPage] = useState(0);
  const [pageSize, setPageSize] = useState(100); // Larger page size to show more jobs
  const [postedCurrentPage, setPostedCurrentPage] = useState(0);
  const [postedPageSize, setPostedPageSize] = useState(100);

  const handleTabChange = async (value: string) => {
    await navigate({
      to: "/admin",
      search: (prev) => ({ ...prev, jobTab: value }),
    });
  };

  const {
    data: pendingJobs,
    isLoading: pendingLoading,
    error: pendingError,
  } = usePendingJobs(pageSize, currentPage * pageSize);

  const {
    data: postedJobs,
    isLoading: postedLoading,
    error: postedError,
  } = usePostedJobs(postedPageSize, postedCurrentPage * postedPageSize);


  const handleJobSelect = (jobId: string, selected: boolean) => {
    if (selected) {
      setSelectedJobs((prev) => [...prev, jobId]);
    } else {
      setSelectedJobs((prev) => prev.filter((id) => id !== jobId));
    }
  };

  const handleSelectAll = () => {
    if (!pendingJobs?.jobs) return;
    setSelectedJobs(pendingJobs.jobs.map((job) => job.id));
  };

  const handleClearSelection = () => {
    setSelectedJobs([]);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    setSelectedJobs([]); // Clear selection when changing pages
  };

  const handlePageSizeChange = (newPageSize: number) => {
    setPageSize(newPageSize);
    setCurrentPage(0); // Reset to first page
    setSelectedJobs([]); // Clear selection
  };

  const handlePostedPageChange = (page: number) => {
    setPostedCurrentPage(page);
  };

  const handlePostedPageSizeChange = (newPageSize: number) => {
    setPostedPageSize(newPageSize);
    setPostedCurrentPage(0); // Reset to first page
  };

  if (pendingLoading) {
    return (
      <div className="space-y-6">
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading pending jobs...</p>
        </div>
      </div>
    );
  }

  if (pendingError) {
    return (
      <div className="space-y-6">
        <Card>
          <CardContent className="pt-6">
            <div className="text-center text-red-600">
              <p className="font-medium">Error loading pending jobs</p>
              <p className="text-sm text-gray-600 mt-1">
                Please check if you have admin permissions and try again.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">
              Pending Approval
            </CardTitle>
            <Badge
              variant="secondary"
              className="bg-yellow-100 text-yellow-800"
            >
              {pendingJobs?.total || 0}
            </Badge>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{pendingJobs?.total || 0}</div>
            <p className="text-xs text-muted-foreground">
              Jobs waiting for Reddit approval
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Selected</CardTitle>
            <Badge variant="secondary" className="bg-blue-100 text-blue-800">
              {selectedJobs.length}
            </Badge>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{selectedJobs.length}</div>
            <p className="text-xs text-muted-foreground">
              Jobs selected for bulk action
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Admin User</CardTitle>
            <Badge variant="default" className="bg-green-100 text-green-800">
              {user.role}
            </Badge>
          </CardHeader>
          <CardContent>
            <div className="text-lg font-medium">{user.email}</div>
            <p className="text-xs text-muted-foreground">
              Logged in administrator
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Tabs
        value={activeTab}
        onValueChange={handleTabChange}
        className="space-y-4"
      >
        <TabsList>
          <TabsTrigger value="pending">Pending Jobs</TabsTrigger>
          <TabsTrigger value="posted">Posted Jobs</TabsTrigger>
          <TabsTrigger value="scraper">Scraper</TabsTrigger>
          <TabsTrigger value="subreddits">Subreddits</TabsTrigger>
        </TabsList>

        <TabsContent value="pending" className="space-y-4">
          <section>
            <header className="mb-8">
              <h3 className="font-bold text-xl mb-2">
                Reddit Posting Approval
              </h3>
              <p className="text-muted-foreground">
                Review and approve job postings for Reddit. Approved jobs will
                be automatically posted to Reddit.
              </p>
            </header>
            <PendingJobsList
              jobs={pendingJobs?.jobs || []}
              selectedJobs={selectedJobs}
              onJobSelect={handleJobSelect}
              onSelectAll={handleSelectAll}
              onClearSelection={handleClearSelection}
              user={user}
              pagination={{
                currentPage,
                pageSize,
                total: pendingJobs?.total || 0,
                hasMore: pendingJobs?.has_more || false,
              }}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
            />
          </section>
        </TabsContent>

        <TabsContent value="posted" className="space-y-4">
          <section>
            <header className="mb-8">
              <h3 className="font-bold text-xl mb-2">Posted Jobs</h3>
              <p className="text-muted-foreground">
                View jobs that have been approved and posted to Reddit.
              </p>
            </header>
            {postedError ? (
              <Card>
                <CardContent className="pt-6">
                  <div className="text-center text-red-600">
                    <p className="font-medium">Error loading posted jobs</p>
                    <p className="text-sm text-gray-600 mt-1">
                      Please check your permissions and try again.
                    </p>
                  </div>
                </CardContent>
              </Card>
            ) : (
              <PostedJobsList
                jobs={postedJobs?.jobs || []}
                pagination={{
                  currentPage: postedCurrentPage,
                  pageSize: postedPageSize,
                  total: postedJobs?.total || 0,
                  hasMore: postedJobs?.has_more || false,
                }}
                onPageChange={handlePostedPageChange}
                onPageSizeChange={handlePostedPageSizeChange}
                isLoading={postedLoading}
              />
            )}
          </section>
        </TabsContent>

        <TabsContent value="scraper" className="space-y-4">
          <ScraperManager />
        </TabsContent>

        <TabsContent value="subreddits" className="space-y-4">
          <SubredditManager />
        </TabsContent>

      </Tabs>
    </div>
  );
}
