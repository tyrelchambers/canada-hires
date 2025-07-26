import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";

export function AuthNav() {
  const { data: user } = useCurrentUser();
  const logoutMutation = useLogout();
  const navigate = useNavigate();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const isAuthenticated = !!user;

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

  if (isAuthenticated && user) {
    return (
      <div className="bg-gray-100 border-b">
        <div className="flex items-center justify-between p-4">
          <Link to="/" className="font-bold text-lg">
            Canada Hires
          </Link>
          
          {/* Desktop Navigation */}
          <nav className="hidden md:flex gap-4">
            <Link to="/lmia" className="text-sm hover:text-blue-600">
              LMIA Search
            </Link>
            <Link to="/jobs" className="text-sm hover:text-blue-600">
              Job Postings
            </Link>
          </nav>
          
          {/* Desktop User Info & Logout */}
          <div className="hidden md:flex items-center gap-4">
            <div className="text-right">
              <p className="text-sm font-medium">
                Welcome, {user.username}
              </p>
              <p className="text-xs text-gray-600 capitalize">
                {user.verification_tier} account •{" "}
                {user.email_verified ? "Verified" : "Unverified"}
              </p>
            </div>
            <Button
              onClick={handleLogout}
              variant="outline"
              size="sm"
              disabled={logoutMutation.isPending}
            >
              {logoutMutation.isPending ? "Signing Out..." : "Sign Out"}
            </Button>
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
                to="/lmia" 
                className="text-sm hover:text-blue-600 py-2"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                LMIA Search
              </Link>
              <Link 
                to="/jobs" 
                className="text-sm hover:text-blue-600 py-2"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Job Postings
              </Link>
              <div className="border-t pt-3 mt-3">
                <p className="text-sm font-medium mb-1">
                  Welcome, {user.username}
                </p>
                <p className="text-xs text-gray-600 capitalize mb-3">
                  {user.verification_tier} account •{" "}
                  {user.email_verified ? "Verified" : "Unverified"}
                </p>
                <Button
                  onClick={() => {
                    handleLogout();
                    setIsMobileMenuOpen(false);
                  }}
                  variant="outline"
                  size="sm"
                  className="w-full"
                  disabled={logoutMutation.isPending}
                >
                  {logoutMutation.isPending ? "Signing Out..." : "Sign Out"}
                </Button>
              </div>
            </nav>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="bg-gray-100 border-b">
      <div className="flex items-center justify-between p-4">
        <Link to="/" className="font-bold text-lg">
          Canada Hires
        </Link>
        
        {/* Desktop Navigation */}
        <nav className="hidden md:flex gap-4">
          <Link to="/lmia" className="text-sm hover:text-blue-600">
            LMIA Search
          </Link>
          <Link to="/jobs" className="text-sm hover:text-blue-600">
            Job Postings
          </Link>
        </nav>
        
        {/* Desktop Sign In */}
        <div className="hidden md:flex items-center gap-4">
          <div className="text-right">
            <p className="text-sm font-medium">Canada Hires</p>
            <p className="text-xs text-gray-600">Sign in to submit reports</p>
          </div>
          <Button onClick={handleLogin}>Sign In</Button>
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
              to="/lmia" 
              className="text-sm hover:text-blue-600 py-2"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              LMIA Search
            </Link>
            <Link 
              to="/jobs" 
              className="text-sm hover:text-blue-600 py-2"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Job Postings
            </Link>
            <div className="border-t pt-3 mt-3 text-center">
              <p className="text-sm font-medium mb-1">Canada Hires</p>
              <p className="text-xs text-gray-600 mb-3">Sign in to submit reports</p>
              <Button 
                onClick={() => {
                  handleLogin();
                  setIsMobileMenuOpen(false);
                }}
                className="w-full"
              >
                Sign In
              </Button>
            </div>
          </nav>
        </div>
      )}
    </div>
  );
}
