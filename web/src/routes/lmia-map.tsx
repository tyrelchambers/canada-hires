import { createFileRoute } from "@tanstack/react-router";
import { LMIAMapHeatmap } from "@/components/LMIAMapHeatmap";
import { AuthNav } from "@/components/AuthNav";
import { PageLoader } from "@/components/shared/PageLoader";
import { useLMIAPostalCodeLocations } from "@/hooks/useLMIA";

export const Route = createFileRoute("/lmia-map")({
  component: LMIAMapComponent,
});

function LMIAMapComponent() {
  const currentYear = new Date().getFullYear();
  const { isLoading } = useLMIAPostalCodeLocations(
    currentYear,
    undefined,
    1000,
  );

  if (isLoading) {
    return (
      <>
        <AuthNav />
        <PageLoader text="Loading LMIA map data..." />
      </>
    );
  }

  return (
    <div className="lg:h-[calc(100vh-72px)]">
      <AuthNav />
      <LMIAMapHeatmap />
    </div>
  );
}
