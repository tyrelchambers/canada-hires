import { createFileRoute, Link } from "@tanstack/react-router";
import { AuthNav } from "@/components/AuthNav";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFileText,
  faArrowRight,
  faCalendar,
  faDollarSign,
  faBuilding,
  faUsers,
} from "@fortawesome/free-solid-svg-icons";

export const Route = createFileRoute("/research/")({
  component: ResearchPage,
  head: () => ({
    meta: [
      {
        title: "Research - JobWatch Canada",
        description:
          "In-depth research and analysis on Canadian employment, immigration policies, and business hiring practices. Data-driven insights for informed decision making.",
      },
    ],
  }),
});

function ResearchPage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <AuthNav />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Research & Analysis
          </h1>
          <p className="text-xl text-gray-600 mb-6">
            Data-driven insights into Canadian employment policies, immigration
            programs, and business hiring practices
          </p>
        </div>

        {/* Featured Research */}
        <div className="mb-12">
          <h2 className="text-2xl font-semibold text-gray-900 mb-6">
            Featured Research
          </h2>

          <Card className="mb-6 border-blue-200 bg-blue-50">
            <CardContent>
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <Badge
                      variant="secondary"
                      className="bg-blue-100 text-gray-800"
                    >
                      New
                    </Badge>
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    Canadian Wage Subsidy Programs for Hiring Immigrants and
                    Newcomers (2025)
                  </h3>
                  <p className="text-gray-700 mb-4">
                    Comprehensive analysis of federal, provincial, and
                    territorial programs providing wage subsidies to businesses
                    hiring immigrants. Includes subsidy amounts, eligibility
                    criteria, and 2025 policy changes.
                  </p>
                  <div className="flex items-center gap-4 text-sm text-gray-600 mb-4">
                    <div className="flex items-center gap-1">
                      <FontAwesomeIcon icon={faCalendar} className="w-4 h-4" />
                      <span>January 2025</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <FontAwesomeIcon icon={faFileText} className="w-4 h-4" />
                      <span>Policy Analysis</span>
                    </div>
                  </div>
                  <Button asChild variant="default">
                    <Link to="/research/wage-subsidies-immigrants">
                      Read Full Analysis
                      <FontAwesomeIcon
                        icon={faArrowRight}
                        className="w-4 h-4 ml-2"
                      />
                    </Link>
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Research Categories */}
        <div className="mb-12">
          <h2 className="text-2xl font-semibold text-gray-900 mb-6">
            Research Categories
          </h2>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Immigration Policy */}
            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faUsers}
                    className="w-5 h-5 text-gray-600"
                  />
                  Immigration Policy
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600 mb-4">
                  Analysis of Canadian immigration programs, policy changes, and
                  their impact on the labor market.
                </p>
                <div className="space-y-2">
                  <div className="text-sm text-gray-600">Featured Topics:</div>
                  <ul className="text-sm space-y-1">
                    <li>• Temporary Foreign Worker Program</li>
                    <li>• Provincial Nominee Programs</li>
                    <li>• Immigration level targets</li>
                    <li>• Policy impact assessments</li>
                  </ul>
                </div>
              </CardContent>
            </Card>

            {/* Business Hiring Practices */}
            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faBuilding}
                    className="w-5 h-5 text-gray-600"
                  />
                  Business Hiring Practices
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600 mb-4">
                  Research into Canadian business hiring patterns, industry
                  trends, and employment practices.
                </p>
                <div className="space-y-2">
                  <div className="text-sm text-gray-500">Featured Topics:</div>
                  <ul className="text-sm space-y-1">
                    <li>• LMIA usage patterns</li>
                    <li>• Industry hiring trends</li>
                    <li>• Regional employment data</li>
                    <li>• Employer compliance analysis</li>
                  </ul>
                </div>
              </CardContent>
            </Card>

            {/* Economic Impact */}
            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faDollarSign}
                    className="w-5 h-5 text-gray-600"
                  />
                  Economic Impact
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600 mb-4">
                  Economic analysis of immigration policies, wage trends, and
                  labor market dynamics.
                </p>
                <div className="space-y-2">
                  <div className="text-sm text-gray-500">Featured Topics:</div>
                  <ul className="text-sm space-y-1">
                    <li>• Wage subsidy effectiveness</li>
                    <li>• Labor market impacts</li>
                    <li>• Regional economic effects</li>
                    <li>• Cost-benefit analysis</li>
                  </ul>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Research Methods */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Our Research Methodology</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  Data Sources
                </h4>
                <ul className="space-y-1 text-sm text-gray-600">
                  <li>• Official government databases and publications</li>
                  <li>• Employment and Social Development Canada (ESDC)</li>
                  <li>• Immigration, Refugees and Citizenship Canada (IRCC)</li>
                  <li>• Provincial and territorial government sources</li>
                  <li>• Statistics Canada labor force data</li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  Research Standards
                </h4>
                <ul className="space-y-1 text-sm text-gray-600">
                  <li>• Fact-checking through multiple sources</li>
                  <li>• Regular updates as policies change</li>
                  <li>• Transparent methodology documentation</li>
                  <li>• Clear distinction between data and analysis</li>
                  <li>• Correction process for identified errors</li>
                </ul>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Coming Soon */}
        <Card className="mb-8 border-gray-200">
          <CardHeader>
            <CardTitle className="text-gray-700">Coming Soon</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-gray-400 rounded-full mt-2"></div>
                <div>
                  <h4 className="font-semibold text-gray-700">
                    LMIA Processing Times Analysis
                  </h4>
                  <p className="text-sm text-gray-600">
                    Trends in LMIA application processing times across different
                    industries and regions
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-gray-400 rounded-full mt-2"></div>
                <div>
                  <h4 className="font-semibold text-gray-700">
                    Provincial Immigration Program Comparison
                  </h4>
                  <p className="text-sm text-gray-600">
                    Detailed comparison of provincial nominee programs and their
                    effectiveness
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-gray-400 rounded-full mt-2"></div>
                <div>
                  <h4 className="font-semibold text-gray-700">
                    Industry-Specific TFW Usage Patterns
                  </h4>
                  <p className="text-sm text-gray-600">
                    Analysis of Temporary Foreign Worker usage by industry
                    sector and geographic region
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-gray-400 rounded-full mt-2"></div>
                <div>
                  <h4 className="font-semibold text-gray-700">
                    LMIA Application Trends Dashboard
                  </h4>
                  <p className="text-sm text-gray-600">
                    Interactive charts showing monthly and yearly trends in LMIA job applications, processing times, and regional patterns
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-gray-400 rounded-full mt-2"></div>
                <div>
                  <h4 className="font-semibold text-gray-700">
                    Data Visualization Suite
                  </h4>
                  <p className="text-sm text-gray-600">
                    Comprehensive charts and graphs for wage subsidy effectiveness, employment patterns, and policy impact analysis
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Contact */}
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                Research Requests & Feedback
              </h3>
              <p className="text-gray-600 mb-4">
                Have a specific research question or found an issue with our
                analysis? We'd love to hear from you.
              </p>
              <Button asChild variant="outline">
                <Link to="/feedback">Submit Research Request</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
