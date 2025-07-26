import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { useCurrentUser } from "../hooks/useAuth";

export const Route = createRootRoute({
  component: () => {
    useCurrentUser();
    return (
      <>
        <Outlet />
        <TanStackRouterDevtools />
      </>
    );
  },
});
