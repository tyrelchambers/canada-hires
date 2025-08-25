import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFlag,
  faUsers,
  faArrowRight,
} from "@fortawesome/free-solid-svg-icons";
import img1 from "@/assets/reports-screen.png";

export function ReportingBanner() {
  return (
    <section className="py-16 bg-white border-b border-border">
      <div className="max-w-screen-2xl mx-auto px-4 flex flex-col gap-6 md:gap-20 items-center">
        {/* Left side - Call to Action */}
        <div className="space-y-6 flex-1">
          <div className="flex items-center gap-3">
            <Badge variant="secondary">Community Reporting</Badge>
          </div>

          <div>
            <h2 className="text-3xl lg:text-4xl font-bold text-space-cadet mb-4">
              Help Build Our{" "}
              <span className="text-pale-dogwood">Community</span> Database
            </h2>
            <p className="text-lg text-gray-600 leading-relaxed">
              Share your knowledge about local businesses' hiring practices.
              Every report helps fellow Canadians make informed decisions about
              where to spend their money.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-4">
            <Link
              to="/reports/create"
              search={{ businessName: undefined, businessAddress: undefined }}
            >
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

        <div className="relative hidden md:block md:h-[400px]  aspect-video overflow-hidden bg-gradient-to-tl from-raising-black to-ultra-violet rounded-3xl shadow-lg">
          <img
            src={img1}
            className="absolute w-full h-full  top-[15%] left-[9%] rounded-xl"
          />
        </div>
      </div>
    </section>
  );
}
