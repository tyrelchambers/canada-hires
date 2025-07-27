import { Link } from "@tanstack/react-router";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faReddit } from "@fortawesome/free-brands-svg-icons";
export function Footer() {
  return (
    <footer className="border-t border-border bg-gray-50">
      <div className="max-w-5xl mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Logo and Mission */}
          <div className="md:col-span-2">
            <div className="flex items-center gap-3 mb-4">
              <img
                src="/canada hires - logo.svg"
                alt="JobWatch Canada Logo"
                className="h-8 w-8"
              />
              <h3 className="text-xl font-semibold">JobWatch Canada</h3>
            </div>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Providing transparency about Canadian businesses' use of the
              Temporary Foreign Worker program, helping you support companies
              that prioritize hiring Canadians.
            </p>
            <p className="text-sm text-muted-foreground">
              Contact:{" "}
              <a
                href="mailto:contact@jobwatchcanada.com"
                className="text-blue-600 hover:underline"
              >
                contact@jobwatchcanada.com
              </a>
            </p>
          </div>

          {/* Navigation Links */}
          <div>
            <h4 className="font-semibold mb-4">Explore</h4>
            <ul className="space-y-2">
              <li>
                <Link
                  to="/lmia"
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  LMIA Records
                </Link>
              </li>
              <li>
                <Link
                  to="/jobs"
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Job Postings
                </Link>
              </li>
              <li>
                <Link
                  to="/directory"
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Business Directory
                </Link>
              </li>
              <li>
                <Link
                  to="/research"
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Research
                </Link>
              </li>
              <li>
                <Link
                  to="/feedback"
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Feedback
                </Link>
              </li>
            </ul>
          </div>

          {/* Community */}
          <div>
            <h4 className="font-semibold mb-4">Community</h4>
            <div className="space-y-3">
              <a
                href="https://reddit.com/user/jobwatchcanada"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
              >
                <FontAwesomeIcon icon={faReddit} className="text-lg" />
                u/jobwatchcanada
              </a>
              <a
                href="https://reddit.com/r/lmiascams"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
              >
                <FontAwesomeIcon icon={faReddit} className="text-lg" />
                r/lmiascams
              </a>
            </div>
          </div>
        </div>

        {/* Copyright */}
        <div className="border-t border-border mt-8 pt-8 text-center">
          <p className="text-sm text-muted-foreground">
            © {new Date().getFullYear()} JobWatch Canada. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
