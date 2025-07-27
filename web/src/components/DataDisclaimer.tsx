import { Link } from "@tanstack/react-router";

interface DataDisclaimerProps {
  variant?: "default" | "compact";
  className?: string;
}

export function DataDisclaimer({
  variant = "default",
  className = "",
}: DataDisclaimerProps) {
  if (variant === "compact") {
    return (
      <div
        className={`bg-blue-50 p-3 rounded-lg border border-blue-200 text-sm ${className}`}
      >
        <p className="text-blue-800">
          <strong>Data Disclaimer:</strong> All government and business data is
          sourced from official databases. While we do our best to ensure
          accuracy, occasional errors may occur due to system complexity.{" "}
          <Link to="/feedback" className="underline hover:text-blue-900">
            Report issues here
          </Link>
          .
        </p>
      </div>
    );
  }

  return (
    <div
      className={`bg-blue-50 p-4 rounded-lg border border-blue-200 ${className}`}
    >
      <h3 className="font-semibold text-blue-900 mb-2">Data Disclaimer</h3>
      <p className="text-blue-800 text-sm">
        All government and business data is sourced from official databases and
        public records. While we make every effort to ensure accuracy, the
        complexity of these data systems may result in occasional errors. Our
        goal is to make this information more accessible and transparent. If you
        notice any inaccuracies or have feedback, please{" "}
        <Link to="/feedback" className="underline hover:text-blue-900">
          contact us
        </Link>{" "}
        - we value your input and will work promptly to address any issues.
      </p>
    </div>
  );
}
