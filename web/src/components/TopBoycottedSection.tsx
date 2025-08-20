import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useTopBoycotted } from "@/hooks/useBoycotts";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBan, faMapMarkerAlt } from "@fortawesome/free-solid-svg-icons";
import { Link } from "@tanstack/react-router";
import { BoycottButton } from "./BoycottButton";

export function TopBoycottedSection() {
  const { data: topBoycotted, isLoading, error } = useTopBoycotted(3);

  if (isLoading) {
    return (
      <section className="my-20 mx-4 lg:max-w-5xl lg:w-full lg:mx-auto">
        <h2 className="text-3xl lg:text-4xl -tracking-wide text-center mb-6 font-medium">
          Most Boycotted Companies
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[...Array.from({ length: 3 })].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardContent className="p-6">
                <div className="h-4 bg-gray-200 rounded mb-2"></div>
                <div className="h-3 bg-gray-200 rounded mb-4"></div>
                <div className="h-6 bg-gray-200 rounded"></div>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>
    );
  }

  if (error || !topBoycotted || topBoycotted.length === 0) {
    return null; // Don't show the section if there's an error or no data
  }

  return (
    <section className="my-20 mx-4 lg:max-w-5xl lg:w-full lg:mx-auto">
      <div className="text-center mb-12">
        <div className="flex justify-center mb-4">
          <FontAwesomeIcon icon={faBan} className="text-4xl text-red-600" />
        </div>
        <h2 className="text-3xl lg:text-4xl -tracking-wide font-medium mb-4">
          Most Boycotted Companies
        </h2>
        <p className="text-lg lg:text-xl text-gray-600 font-light max-w-3xl mx-auto">
          These companies have received the most boycotts from our community
          members based on their TFW practices.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {topBoycotted.map((company, index) => (
          <Card
            key={`${company.business_name}-${company.business_address}`}
            className="hover:shadow-lg transition-shadow relative"
          >
            {index === 0 && (
              <div className="absolute -top-3 -right-3 bg-red-600 text-white text-xs font-bold px-2 py-1 rounded-full">
                #1 Most Boycotted
              </div>
            )}
            <CardHeader className="pb-4">
              <CardTitle className="text-lg line-clamp-2">
                {company.business_name}
              </CardTitle>
              <div className="flex items-start text-sm text-gray-600">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 mr-2 mt-0.5 flex-shrink-0"
                />
                <span className="line-clamp-2">{company.business_address}</span>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Badge variant="secondary" className="text-sm">
                  {company.boycott_count} boycotting
                </Badge>
                <span className="text-xs text-gray-500">#{index + 1}</span>
              </div>

              <div className="space-y-2">
                <Button variant="outline" size="sm" className="w-full" asChild>
                  <Link
                    to="/business/$address"
                    params={{
                      address: encodeURIComponent(company.business_address),
                    }}
                  >
                    View Details
                  </Link>
                </Button>
                <BoycottButton
                  businessName={company.business_name}
                  businessAddress={company.business_address}
                  className="w-full"
                />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="text-center mt-8">
        <Button variant="outline" asChild>
          <Link to="/reports">Browse All Companies</Link>
        </Button>
      </div>
    </section>
  );
}
