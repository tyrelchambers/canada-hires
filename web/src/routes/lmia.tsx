import { createFileRoute } from "@tanstack/react-router";
import { AuthNav } from "@/components/AuthNav";
import { LMIASearch } from "@/components/LMIASearch";
import { LMIAStats } from "@/components/LMIAStats";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileAlt, faSearch, faChartBar, faInfoCircle } from "@fortawesome/free-solid-svg-icons";

type LMIASearchParams = {
  tab?: "search" | "stats" | "about";
};

export const Route = createFileRoute("/lmia")({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>): LMIASearchParams => {
    return {
      tab: (search.tab as "search" | "stats" | "about") || "search",
    };
  },
});

function RouteComponent() {
  const { tab } = Route.useSearch();
  const navigate = Route.useNavigate();

  const handleTabChange = (value: string) => {
    navigate({
      search: { tab: value as "search" | "stats" | "about" },
    });
  };
  return (
    <div className="min-h-screen bg-gray-50">
      <AuthNav />

      <div className="max-w-7xl mx-auto py-8 px-4">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            LMIA Employer Directory
          </h1>
          <p className="text-lg text-gray-600">
            Search and explore employers who have received positive Labour Market Impact Assessments (LMIA) 
            in Canada's Temporary Foreign Worker Program.
          </p>
        </div>

        <Tabs value={tab} onValueChange={handleTabChange} className="space-y-6">
          <TabsList className="grid w-full grid-cols-3 lg:w-96">
            <TabsTrigger value="search" className="flex items-center gap-2">
              <FontAwesomeIcon icon={faSearch} className="w-4 h-4" />
              Search
            </TabsTrigger>
            <TabsTrigger value="stats" className="flex items-center gap-2">
              <FontAwesomeIcon icon={faChartBar} className="w-4 h-4" />
              Statistics
            </TabsTrigger>
            <TabsTrigger value="about" className="flex items-center gap-2">
              <FontAwesomeIcon icon={faInfoCircle} className="w-4 h-4" />
              About
            </TabsTrigger>
          </TabsList>

          <TabsContent value="search" className="space-y-6">
            <LMIASearch />
          </TabsContent>

          <TabsContent value="stats" className="space-y-6">
            <LMIAStats />
          </TabsContent>

          <TabsContent value="about" className="space-y-6">
            <div className="grid gap-6 md:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FontAwesomeIcon icon={faFileAlt} className="w-5 h-5" />
                    About LMIA Data
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <p className="text-sm text-gray-600">
                    A positive Labour Market Impact Assessment (LMIA) is issued by Service Canada 
                    when an assessment indicates that hiring a temporary foreign worker (TFW) will 
                    have a positive or neutral impact on the Canadian labour market.
                  </p>
                  
                  <div className="space-y-2">
                    <h4 className="font-medium text-sm">Data Source</h4>
                    <p className="text-sm text-gray-600">
                      This data is sourced from Employment and Social Development Canada (ESDC) 
                      through Canada's Open Data portal. The data is updated quarterly and excludes 
                      personal names and businesses that use personal names.
                    </p>
                  </div>

                  <div className="space-y-2">
                    <h4 className="font-medium text-sm">Important Notes</h4>
                    <ul className="text-sm text-gray-600 space-y-1 list-disc list-inside">
                      <li>This data tracks LMIA positions only, not actual work permits issued</li>
                      <li>Not all approved positions result in a TFW entering Canada</li>
                      <li>Some positions may be cancelled after approval</li>
                      <li>The list is not complete due to privacy exclusions</li>
                    </ul>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Data Fields Explained</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-3">
                    <div>
                      <h4 className="font-medium text-sm">Employer Name</h4>
                      <p className="text-sm text-gray-600">
                        The legal name of the company that received the LMIA
                      </p>
                    </div>
                    
                    <div>
                      <h4 className="font-medium text-sm">NOC Code & Title</h4>
                      <p className="text-sm text-gray-600">
                        National Occupational Classification - describes the job role
                      </p>
                    </div>
                    
                    <div>
                      <h4 className="font-medium text-sm">Program Stream</h4>
                      <p className="text-sm text-gray-600">
                        The specific TFW program stream (e.g., High-wage, Low-wage, Caregiver)
                      </p>
                    </div>
                    
                    <div>
                      <h4 className="font-medium text-sm">Positions Approved</h4>
                      <p className="text-sm text-gray-600">
                        Number of TFW positions approved for this employer
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="md:col-span-2">
                <CardHeader>
                  <CardTitle>Data Privacy & Limitations</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <h4 className="font-medium text-sm mb-2">Excluded from Data</h4>
                      <ul className="text-sm text-gray-600 space-y-1 list-disc list-inside">
                        <li>Employers of caregivers (personal names)</li>
                        <li>Business names that include personal names</li>
                        <li>Individual contractors and sole proprietors</li>
                        <li>Some seasonal/temporary employers</li>
                      </ul>
                    </div>
                    
                    <div>
                      <h4 className="font-medium text-sm mb-2">Data Accuracy</h4>
                      <ul className="text-sm text-gray-600 space-y-1 list-disc list-inside">
                        <li>Data is updated quarterly by ESDC</li>
                        <li>Historical data may not reflect current status</li>
                        <li>Processing delays may affect timeliness</li>
                        <li>Some records may contain formatting variations</li>
                      </ul>
                    </div>
                  </div>
                  
                  <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
                    <p className="text-sm text-blue-800">
                      <strong>Disclaimer:</strong> This data is provided for informational purposes only. 
                      For official LMIA information or to verify current status, please contact Employment 
                      and Social Development Canada directly.
                    </p>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}