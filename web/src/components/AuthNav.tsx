import { Link } from "@tanstack/react-router";
import { useState } from "react";
import canadaHires from "@/assets/canada hires.svg";

const navigationLinks = [
  { to: "/lmia", label: "LMIA Search" },
  { to: "/jobs", label: "Job Postings" },
  { to: "/feedback", label: "Feedback" },
];

interface NavLinksProps {
  onLinkClick?: () => void;
  className?: string;
}

function NavLinks({ onLinkClick, className = "" }: NavLinksProps) {
  return (
    <>
      {navigationLinks.map((link) => (
        <Link
          key={link.to}
          to={link.to}
          className={`text-sm hover:text-primary px-4 py-2 rounded-full transition-colors ${className}`}
          activeProps={{
            className:
              "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
          }}
          onClick={onLinkClick}
        >
          {link.label}
        </Link>
      ))}
    </>
  );
}

export function AuthNav() {
  // const { data: user } = useCurrentUser();
  // const logoutMutation = useLogout();
  // const navigate = useNavigate();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // const isAuthenticated = !!user;

  // const handleLogin = () => {
  //   void navigate({ to: "/auth/login" });
  // };

  // const handleLogout = () => {
  //   logoutMutation.mutate(undefined, {
  //     onSuccess: () => {
  //       void navigate({ to: "/" });
  //     },
  //   });
  // };

  return (
    <div className="bg-gray-100 border-b">
      <div className="flex items-center justify-between p-4">
        <Link to="/" className="font-bold text-lg">
          <img src={canadaHires} className="h-10" />
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex gap-4">
          <NavLinks />
        </nav>

        {/* Desktop Sign In */}
        <div></div>
        {/* <div className="hidden md:flex items-center gap-4">
          <div className="text-right">
            <p className="text-sm font-medium">Canada Hires</p>
            <p className="text-xs text-gray-600">Sign in to submit reports</p>
          </div>
          <Button onClick={handleLogin}>Sign In</Button>
        </div> */}

        {/* Mobile Menu Button */}
        <button
          className="md:hidden p-2"
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            {isMobileMenuOpen ? (
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            ) : (
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 6h16M4 12h16M4 18h16"
              />
            )}
          </svg>
        </button>
      </div>

      {/* Mobile Menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t bg-white">
          <nav className="flex flex-col p-4 space-y-3">
            <NavLinks onLinkClick={() => setIsMobileMenuOpen(false)} />
            <div></div>
            {/* <div className="border-t pt-3 mt-3 text-center">
              <p className="text-sm font-medium mb-1">Canada Hires</p>
              <p className="text-xs text-gray-600 mb-3">
                Sign in to submit reports
              </p>
              <Button
                onClick={() => {
                  handleLogin();
                  setIsMobileMenuOpen(false);
                }}
                className="w-full"
              >
                Sign In
              </Button>
            </div> */}
          </nav>
        </div>
      )}
    </div>
  );
}
