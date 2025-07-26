import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { useRouter, Link } from "@tanstack/react-router";

export function AuthNav() {
  const { data: user } = useCurrentUser();
  const logoutMutation = useLogout();
  const router = useRouter();

  const isAuthenticated = !!user;

  const handleLogin = () => {
    router.navigate({ to: "/auth/login" });
  };

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        router.navigate({ to: "/" });
      },
    });
  };

  if (isAuthenticated && user) {
    return (
      <div className="flex items-center gap-4 p-4 bg-gray-100 border-b">
        <Link to="/" className="font-bold text-lg">
          Canada Hires
        </Link>
        <nav className="flex gap-4 ml-8">
          <Link to="/lmia" className="text-sm hover:text-blue-600">
            LMIA Search
          </Link>
          <Link to="/jobs" className="text-sm hover:text-blue-600">
            Job Postings
          </Link>
        </nav>
        <div className="flex-1">
          <p className="text-sm font-medium text-right">
            Welcome, {user.username}
          </p>
          <p className="text-xs text-gray-600 text-right capitalize">
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
    );
  }

  return (
    <div className="flex items-center justify-between p-4 bg-gray-100 border-b">
      <div className="flex items-center gap-8">
        <Link to="/" className="font-bold text-lg">
          Canada Hires
        </Link>
        <nav className="flex gap-4">
          <Link to="/lmia" className="text-sm hover:text-blue-600">
            LMIA Search
          </Link>
          <Link to="/jobs" className="text-sm hover:text-blue-600">
            Job Postings
          </Link>
        </nav>
      </div>
      <div className="flex items-center gap-4">
        <div className="text-right">
          <p className="text-sm font-medium">Canada Hires</p>
          <p className="text-xs text-gray-600">Sign in to submit reports</p>
        </div>
        <Button onClick={handleLogin}>Sign In</Button>
      </div>
    </div>
  );
}
