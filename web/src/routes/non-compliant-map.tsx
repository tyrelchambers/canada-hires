import { createFileRoute } from "@tanstack/react-router";
import { NonCompliantMapHeatmap } from "@/components/NonCompliantMapHeatmap";
import { AuthNav } from "@/components/AuthNav";
import { PageLoader } from "@/components/shared/PageLoader";
import { useNonCompliantLocations } from "@/hooks/useNonCompliant";

export const Route = createFileRoute("/non-compliant-map")({
  component: NonCompliantMapComponent,
});

function NonCompliantMapComponent() {
  const { isLoading } = useNonCompliantLocations(2000);

  if (isLoading) {
    return (
      <>
        <AuthNav />
        <PageLoader text="Loading non-compliant employers map data..." />
      </>
    );
  }

  return (
    <div className="lg:h-[calc(100vh-72px)]">
      <AuthNav />
      <NonCompliantMapHeatmap />
    </div>
  );
}
