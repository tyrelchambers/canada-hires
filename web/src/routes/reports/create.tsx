import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { AuthNav } from "@/components/AuthNav";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { AddressSearch } from "@/components/AddressSearch";
import { useCreateReport } from "@/hooks/useReports";
import { useCurrentUser } from "@/hooks/useAuth";
import { CreateReportRequest } from "@/types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faShieldAlt,
  faHeart,
  faUsers,
} from "@fortawesome/free-solid-svg-icons";

function CreateReportPage() {
  const navigate = useNavigate();
  const { data: user, isLoading, error } = useCurrentUser();
  const createReportMutation = useCreateReport();

  const [formData, setFormData] = useState<CreateReportRequest>({
    business_name: "",
    business_address: "",
    report_source: "employment",
    confidence_level: 5,
    additional_notes: "",
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await createReportMutation.mutateAsync(formData);
      // Navigate back to directory or reports list
      void navigate({ to: "/directory" });
    } catch (error) {
      console.error("Failed to create report:", error);
    }
  };

  const handleInputChange = (
    field: keyof CreateReportRequest,
    value: string | number,
  ) => {
    console.log(value);
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  // Show loading state while checking authentication
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

  // Redirect to login if not authenticated
  if (error || !user) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">
            Sign In Required
          </h1>
          <p className="text-gray-600 mb-4">
            Please sign in to submit a business report.
          </p>
          <a
            href="/auth/login"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
          >
            Go to Sign In
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <AuthNav />

      <div className="container mx-auto px-4 py-8 max-w-6xl">
        <div className="mb-6">
          <Button
            variant="ghost"
            onClick={() => navigate({ to: "/directory" })}
            className="mb-4"
          >
            <FontAwesomeIcon icon={faArrowLeft} className="mr-2" />
            Back to Directory
          </Button>

          <h1 className="text-3xl font-bold mb-2 flex items-center">
            Create Business Report
          </h1>
          <p className="text-gray-600">
            Help build our community database by reporting on business hiring
            practices
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Form */}
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Report Details</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-6">
                  <div className="space-y-2">
                    <Label htmlFor="business_name">Business Name *</Label>
                    <Input
                      id="business_name"
                      required
                      value={formData.business_name}
                      onChange={(e) =>
                        handleInputChange("business_name", e.target.value)
                      }
                      placeholder="Enter the business name"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="business_address">Business Address *</Label>
                    <AddressSearch
                      id="business_address"
                      required
                      value={formData.business_address}
                      onChange={(value) =>
                        handleInputChange("business_address", value)
                      }
                      placeholder="Search for business address..."
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="report_source">Report Source *</Label>
                    <select
                      id="report_source"
                      required
                      value={formData.report_source}
                      onChange={(e) =>
                        handleInputChange(
                          "report_source",
                          e.target.value as
                            | "employment"
                            | "observation"
                            | "public_record",
                        )
                      }
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      <option value="employment">Employment Experience</option>
                      <option value="observation">Personal Observation</option>
                      <option value="public_record">Public Record</option>
                    </select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="confidence_level">
                      Confidence Level: {formData.confidence_level}/10
                    </Label>
                    <input
                      type="range"
                      id="confidence_level"
                      min="1"
                      max="10"
                      value={formData.confidence_level || 5}
                      onChange={(e) =>
                        handleInputChange(
                          "confidence_level",
                          parseInt(e.target.value),
                        )
                      }
                      className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                    />
                    <div className="flex justify-between text-xs text-gray-500">
                      <span>1 - Low</span>
                      <span>5 - Medium</span>
                      <span>10 - High</span>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="additional_notes">Additional Notes</Label>
                    <Textarea
                      id="additional_notes"
                      value={formData.additional_notes}
                      onChange={(e) =>
                        handleInputChange(
                          "additional_notes",
                          e.currentTarget.value,
                        )
                      }
                      placeholder="Provide any additional details about your observations..."
                      rows={4}
                    />
                  </div>

                  <div className="flex gap-4 pt-4">
                    <Button
                      type="submit"
                      disabled={createReportMutation.isPending}
                      className="flex-1"
                    >
                      {createReportMutation.isPending
                        ? "Submitting..."
                        : "Submit Report"}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => navigate({ to: "/directory" })}
                    >
                      Cancel
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar Cards */}
          <div className="space-y-6">
            {/* Integrity Card */}
            <Card>
              <CardContent className="p-6">
                <div className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <FontAwesomeIcon icon={faShieldAlt} className="text-2xl" />
                  </div>
                  <div>
                    <h3 className="font-semibold mb-2">
                      Your Integrity Matters
                    </h3>
                    <p className="text-sm">
                      Accurate reporting helps fellow Canadians make informed
                      decisions. Please only submit information you can verify
                      and trust.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Community Impact Card */}
            <Card>
              <CardContent className="p-6">
                <div className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <FontAwesomeIcon icon={faUsers} className="text-2xl " />
                  </div>
                  <div>
                    <h3 className="font-semibold  mb-2">
                      Building Our Community
                    </h3>
                    <p className="text-sm ">
                      Every truthful report strengthens our collective knowledge
                      and helps Canadian workers find businesses that value
                      local hiring.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Value Card */}
            <Card className="border-purple-200 bg-purple-50">
              <CardContent className="p-6">
                <div className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <FontAwesomeIcon
                      icon={faHeart}
                      className="text-2xl text-purple-600"
                    />
                  </div>
                  <div>
                    <h3 className="font-semibold text-purple-900 mb-2">
                      Your Voice Has Power
                    </h3>
                    <p className="text-sm text-purple-800">
                      Your honest experiences and observations contribute to
                      transparency in Canadian hiring practices. Thank you for
                      taking the time to share.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute("/reports/create")({
  component: CreateReportPage,
});
