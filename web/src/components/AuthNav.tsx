import { Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import canadaHires from "@/assets/canada hires.png";
import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFileAlt,
  faSignOutAlt,
  faChevronDown,
  faBan,
} from "@fortawesome/free-solid-svg-icons";

const navigationLinks = [
  { to: "/lmia", label: "LMIA Search" },
  { to: "/lmia-map", label: "LMIA Heatmap" },
  { to: "/trends", label: "Trends" },
  { to: "/research", label: "Research" },
  { to: "/feedback", label: "Feedback" },
];

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

  const handleCreateReport = () => {
    void navigate({ to: "/reports/create", search: { businessName: undefined, businessAddress: undefined } });
  };

  // Generate initials from email for avatar fallback
  const getInitials = (email: string) => {
    return email.slice(0, 2).toUpperCase();
  };

  return (
    <div className="border-b bg-white">
      <div className="flex items-center justify-between p-4">
        <Link to="/" className="font-bold text-lg">
          <img src={canadaHires} className="h-10" />
        </Link>

        <div className="flex gap-4">
          {/* Desktop Navigation */}
          <nav className="hidden md:flex gap-4">
            <Link
              to="/jobs"
              className="text-sm hover:text-primary px-4 py-2 rounded-full transition-colors"
              activeProps={{
                className:
                  "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
              }}
            >
              Job Postings
            </Link>
            <Link
              to="/reports"
              className="text-sm hover:text-primary px-4 py-2 rounded-full transition-colors"
              activeProps={{
                className:
                  "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
              }}
            >
              Reports
            </Link>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="flex items-center gap-1">
                  Explore
                  <FontAwesomeIcon icon={faChevronDown} className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                {navigationLinks.map((link) => (
                  <DropdownMenuItem key={link.to} asChild>
                    <Link to={link.to} className="flex items-center">
                      {link.label}
                    </Link>
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </nav>
          {/* Desktop Auth */}
          <div className="hidden md:flex items-center gap-4">
            {isAuthenticated ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="flex items-center gap-2 p-2"
                  >
                    <Avatar className="h-8 w-8">
                      <AvatarFallback className="bg-primary text-primary-foreground text-xs">
                        {getInitials(user.email)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="hidden sm:block text-left">
                      <p className="text-sm font-medium">{user.email}</p>
                    </div>
                    <FontAwesomeIcon icon={faChevronDown} className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56">
                  <DropdownMenuItem onClick={handleCreateReport}>
                    <FontAwesomeIcon
                      icon={faFileAlt}
                      className="mr-2 h-4 w-4"
                    />
                    Create Report
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link to="/boycotts" className="flex items-center">
                      <FontAwesomeIcon icon={faBan} className="mr-2 h-4 w-4" />
                      My Boycotts
                    </Link>
                  </DropdownMenuItem>
                  {isAdmin && (
                    <>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem asChild>
                        <Link
                          to="/admin"
                          search={{ tab: undefined, jobTab: undefined }}
                          className="flex items-center"
                        >
                          Admin Dashboard
                        </Link>
                      </DropdownMenuItem>
                    </>
                  )}
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={handleLogout}
                    disabled={logoutMutation.isPending}
                  >
                    <FontAwesomeIcon
                      icon={faSignOutAlt}
                      className="mr-2 h-4 w-4"
                    />
                    {logoutMutation.isPending ? "Logging out..." : "Log out"}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <Button onClick={handleLogin}>Sign In</Button>
            )}
          </div>
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
            <Link
              to="/jobs"
              className="text-sm hover:text-primary px-4 py-2 rounded-full transition-colors"
              activeProps={{
                className:
                  "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
              }}
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Job Postings
            </Link>
            <Link
              to="/reports"
              className="text-sm hover:text-primary px-4 py-2 rounded-full transition-colors"
              activeProps={{
                className:
                  "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
              }}
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Reports
            </Link>
            {navigationLinks.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                className="text-sm hover:text-primary px-4 py-2 rounded-full transition-colors"
                activeProps={{
                  className:
                    "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground",
                }}
                onClick={() => setIsMobileMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
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
                  <div className="space-y-2">
                    <Button
                      variant="outline"
                      onClick={() => {
                        handleCreateReport();
                        setIsMobileMenuOpen(false);
                      }}
                      className="w-full"
                    >
                      <FontAwesomeIcon icon={faFileAlt} className="mr-2" />
                      Create Report
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() => {
                        void navigate({ to: "/boycotts" });
                        setIsMobileMenuOpen(false);
                      }}
                      className="w-full"
                    >
                      <FontAwesomeIcon icon={faBan} className="mr-2" />
                      My Boycotts
                    </Button>
                    {isAdmin && (
                      <Button
                        variant="outline"
                        onClick={() => {
                          void navigate({
                            to: "/admin",
                            search: { tab: undefined, jobTab: undefined },
                          });
                          setIsMobileMenuOpen(false);
                        }}
                        className="w-full"
                      >
                        Admin Dashboard
                        <Badge className="ml-2 text-xs bg-red-100 text-red-800">
                          Admin
                        </Badge>
                      </Button>
                    )}
                    <Button
                      variant="outline"
                      onClick={() => {
                        handleLogout();
                        setIsMobileMenuOpen(false);
                      }}
                      className="w-full"
                    >
                      <FontAwesomeIcon icon={faSignOutAlt} className="mr-2" />
                      {logoutMutation.isPending ? "Logging out..." : "Log out"}
                    </Button>
                  </div>
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
