import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useCurrentUser } from "@/hooks/useAuth";
import { JobApprovalDashboard } from "@/components/admin/JobApprovalDashboard";
import { ReportManagementDashboard } from "@/components/admin/ReportManagementDashboard";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

interface Search {
  tab?: string;
  jobTab?: string;
}

export const Route = createFileRoute("/admin")({
  component: AdminPage,
  validateSearch: (search: Search) => ({
    tab: search?.tab,
    jobTab: search?.jobTab,
  }),
});

function AdminPage() {
  const { data: user, isLoading, error } = useCurrentUser();
  const { tab, jobTab } = Route.useSearch();
  const navigate = useNavigate();

  const handleTabChange = async (value: string) => {
    await navigate({
      to: "/admin",
      search: { tab: value },
    });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Access Denied
          </h1>
          <p className="text-gray-600 mb-4">
            Please log in to access the admin dashboard.
          </p>
          <a
            href="/auth/login"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
          >
            Go to Login
          </a>
        </div>
      </div>
    );
  }

  if (user.role !== "admin") {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Access Denied
          </h1>
          <p className="text-gray-600 mb-4">
            You need administrator privileges to access this page.
          </p>
          <a
            href="/"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
          >
            Go to Home
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white shadow">
        <div className="px-4 sm:px-6 lg:px-8">
          <div className="py-6">
            <h1 className="text-3xl font-bold text-gray-900">
              Admin Dashboard
            </h1>
            <p className="mt-2 text-gray-600">
              Manage job postings and business reports
            </p>
          </div>
        </div>
      </div>

      <div className="py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <Tabs value={tab || "jobs"} onValueChange={handleTabChange} className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-6">
              <TabsTrigger value="jobs">Job Management</TabsTrigger>
              <TabsTrigger value="reports">Report Management</TabsTrigger>
            </TabsList>
            
            <TabsContent value="jobs">
              <JobApprovalDashboard user={user} activeTab={jobTab || "pending"} />
            </TabsContent>
            
            <TabsContent value="reports">
              <ReportManagementDashboard />
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
