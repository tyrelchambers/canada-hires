import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFlag,
  faUsers,
  faArrowRight,
} from "@fortawesome/free-solid-svg-icons";
export function ReportingBanner() {
  return (
    <section className="py-16 bg-white border-b border-border">
      <div className="max-w-3xl mx-auto px-4">
        {/* Left side - Call to Action */}
        <div className="space-y-6">
          <div className="flex items-center gap-3">
            <FontAwesomeIcon
              icon={faFlag}
              className="text-3xl text-indigo-600"
            />
            <Badge
              variant="secondary"
              className="bg-indigo-100 text-indigo-800"
            >
              Community Reporting
            </Badge>
          </div>

          <div>
            <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 mb-4">
              Help Build Our Community Database
            </h2>
            <p className="text-lg text-gray-600 leading-relaxed">
              Share your knowledge about local businesses' hiring practices.
              Every report helps fellow Canadians make informed decisions about
              where to spend their money.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-4">
            <Link to="/reports/create">
              <Button size="lg" className="w-full sm:w-auto">
                <FontAwesomeIcon icon={faFlag} className="mr-2" />
                Submit a Report
              </Button>
            </Link>

            <Link to="/reports">
              <Button variant="outline" size="lg" className="w-full sm:w-auto">
                <FontAwesomeIcon icon={faUsers} className="mr-2" />
                Browse All Reports
                <FontAwesomeIcon icon={faArrowRight} className="ml-2" />
              </Button>
            </Link>
          </div>

          <div className="flex items-center gap-3 text-sm text-gray-500">
            <FontAwesomeIcon icon={faUsers} />
            <span>Join thousands of Canadians building transparency</span>
          </div>
        </div>
      </div>
    </section>
  );
}
