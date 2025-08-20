import { createFileRoute, Link } from "@tanstack/react-router";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AuthNav } from "@/components/AuthNav";
import { useUserBoycotts } from "@/hooks/useBoycotts";
import { useCurrentUser } from "@/hooks/useAuth";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBan,
  faMapMarkerAlt,
  faCalendarAlt,
  faArrowLeft,
} from "@fortawesome/free-solid-svg-icons";
import { BoycottButton } from "@/components/BoycottButton";

function MyBoycottsPage() {
  const { data: user, isLoading: authLoading } = useCurrentUser();
  const { data: boycottsData, isLoading, error } = useUserBoycotts();

  // Redirect to login if not authenticated
  if (!authLoading && !user) {
    return (
      <div className="min-h-screen bg-slate-50">
        <AuthNav />
        <div className="max-w-4xl mx-auto px-4 py-8">
          <Card className="border-yellow-200 bg-yellow-50">
            <CardContent className="p-6 text-center">
              <h2 className="text-xl font-semibold mb-4">Please Log In</h2>
              <p className="text-gray-600 mb-4">
                You need to be logged in to view your boycotted companies.
              </p>
              <Button asChild>
                <Link to="/auth/login">Login</Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  if (authLoading || isLoading) {
    return (
      <div className="min-h-screen bg-slate-50">
        <AuthNav />
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-red-600"></div>
            <p className="mt-4 text-gray-600 text-lg">
              Loading your boycotts...
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-slate-50">
        <AuthNav />
        <div className="max-w-4xl mx-auto px-4 py-8">
          <Card className="border-red-200 bg-red-50">
            <CardContent className="p-6">
              <p className="text-red-800">
                Error loading boycotts:{" "}
                {error instanceof Error ? error.message : "An error occurred"}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  const boycotts = boycottsData?.data || [];

  return (
    <div className="min-h-screen bg-slate-50">
      <AuthNav />

      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-4 mb-4">
            <Button variant="ghost" size="sm" asChild>
              <Link to="/">
                <FontAwesomeIcon icon={faArrowLeft} className="mr-2" />
                Back to Home
              </Link>
            </Button>
          </div>

          <div className="flex items-center gap-3 mb-2">
            <FontAwesomeIcon icon={faBan} className="text-2xl text-red-600" />
            <h1 className="text-3xl font-bold">My Boycotts</h1>
          </div>
          <p className="text-gray-600">
            Companies you have chosen to boycott based on their TFW practices
          </p>
        </div>

        {/* Boycotts List */}
        <div className="space-y-4">
          {boycotts.length === 0 ? (
            <Card>
              <CardContent className="p-8 text-center">
                <FontAwesomeIcon
                  icon={faBan}
                  className="text-6xl text-gray-300 mb-4"
                />
                <h3 className="text-xl font-semibold mb-2">No Boycotts Yet</h3>
                <p className="text-gray-600 mb-6">
                  You haven't boycotted any companies yet. Browse our business
                  directory to find companies and take action.
                </p>
                <Button asChild>
                  <Link to="/reports">Browse Business Directory</Link>
                </Button>
              </CardContent>
            </Card>
          ) : (
            boycotts.map((boycott) => (
              <Card
                key={boycott.id}
                className="hover:shadow-md transition-shadow"
              >
                <CardContent className="p-6">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="text-xl font-semibold mb-2">
                        {boycott.business_name}
                      </h3>

                      <div className="flex items-center text-gray-600 mb-3">
                        <FontAwesomeIcon
                          icon={faMapMarkerAlt}
                          className="w-4 h-4 mr-2"
                        />
                        <span>{boycott.business_address}</span>
                      </div>

                      <div className="flex items-center gap-4 mb-4">
                        <div className="flex items-center text-sm text-gray-500">
                          <FontAwesomeIcon
                            icon={faCalendarAlt}
                            className="w-4 h-4 mr-2"
                          />
                          Boycotting since{" "}
                          {new Date(boycott.created_at).toLocaleDateString()}
                        </div>
                      </div>

                      <div className="flex gap-3">
                        <Button variant="outline" size="sm" asChild>
                          <Link
                            to="/business/$address"
                            params={{
                              address: encodeURIComponent(
                                boycott.business_address,
                              ),
                            }}
                          >
                            View Details
                          </Link>
                        </Button>
                        <BoycottButton
                          businessName={boycott.business_name}
                          businessAddress={boycott.business_address}
                        />
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute("/boycotts")({
  component: MyBoycottsPage,
});
