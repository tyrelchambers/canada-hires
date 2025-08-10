import { Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import canadaHires from "@/assets/canada hires.svg";
import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

const navigationLinks = [
  { to: "/lmia", label: "LMIA Search" },
  { to: "/jobs", label: "Job Postings" },
  { to: "/research", label: "Research" },
  { to: "/feedback", label: "Feedback" },
];

const adminNavigationLinks = [{ to: "/admin", label: "Admin Dashboard" }];

interface NavLinksProps {
  onLinkClick?: () => void;
  className?: string;
  isAdmin?: boolean;
}

function NavLinks({
  onLinkClick,
  className = "",
  isAdmin = false,
}: NavLinksProps) {
  const links = isAdmin
    ? [...navigationLinks, ...adminNavigationLinks]
    : navigationLinks;

  return (
    <>
      {links.map((link) => (
        <Link
          key={link.to}
          to={link.to}
          className={`text-sm hover:text-primary px-4 py-2 rounded-full transition-colors ${className} ${
            link.to === "/admin" ? "relative" : ""
          }`}
          activeProps={{
            className:
              "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
          }}
          onClick={onLinkClick}
        >
          {link.label}
          {link.to === "/admin" && (
            <Badge className="ml-2 text-xs bg-red-100 text-red-800 hover:bg-red-100">
              Admin
            </Badge>
          )}
        </Link>
      ))}
    </>
  );
}

export function AuthNav() {
  const { data: user } = useCurrentUser();
  const logoutMutation = useLogout();
  const navigate = useNavigate();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const isAuthenticated = !!user;
  const isAdmin = user?.role === "admin";

  const handleLogin = () => {
    void navigate({ to: "/auth/login" });
  };

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        void navigate({ to: "/" });
      },
    });
  };

  return (
    <div className="border-b bg-white">
      <div className="flex items-center justify-between p-4">
        <Link to="/" className="font-bold text-lg">
          <img src={canadaHires} className="h-10" />
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex gap-4">
          <NavLinks isAdmin={isAdmin} />
        </nav>

        {/* Desktop Auth */}
        <div className="hidden md:flex items-center gap-4">
          {isAuthenticated ? (
            <>
              <div className="text-right">
                <p className="text-sm font-medium">{user.email}</p>
                <div className="flex items-center gap-2">
                  <Badge
                    variant={isAdmin ? "destructive" : "secondary"}
                    className="text-xs"
                  >
                    {isAdmin ? "Admin" : "User"}
                  </Badge>
                  <span className="text-xs text-gray-600">Logged in</span>
                </div>
              </div>
              <Button variant="outline" onClick={handleLogout}>
                Sign Out
              </Button>
            </>
          ) : (
            <>
              <div className="text-right">
                <p className="text-sm font-medium">JobWatch Canada</p>
                <p className="text-xs text-gray-600">
                  Sign in to access admin features
                </p>
              </div>
              <Button onClick={handleLogin}>Sign In</Button>
            </>
          )}
        </div>

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
            <NavLinks
              onLinkClick={() => setIsMobileMenuOpen(false)}
              isAdmin={isAdmin}
            />
            <div className="border-t pt-3 mt-3 text-center">
              {isAuthenticated ? (
                <>
                  <p className="text-sm font-medium mb-1">{user.email}</p>
                  <div className="flex items-center justify-center gap-2 mb-3">
                    <Badge
                      variant={isAdmin ? "destructive" : "secondary"}
                      className="text-xs"
                    >
                      {isAdmin ? "Admin" : "User"}
                    </Badge>
                  </div>
                  <Button
                    variant="outline"
                    onClick={() => {
                      handleLogout();
                      setIsMobileMenuOpen(false);
                    }}
                    className="w-full"
                  >
                    Sign Out
                  </Button>
                </>
              ) : (
                <>
                  <p className="text-sm font-medium mb-1">JobWatch Canada</p>
                  <p className="text-xs text-gray-600 mb-3">
                    Sign in to access admin features
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
                </>
              )}
            </div>
          </nav>
        </div>
      )}
    </div>
  );
}
