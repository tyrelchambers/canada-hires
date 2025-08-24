import { createFileRoute } from "@tanstack/react-router";
import { LMIAMapHeatmap } from "@/components/LMIAMapHeatmap";
import { AuthNav } from "@/components/AuthNav";

export const Route = createFileRoute("/lmia-map")({
  component: LMIAMapComponent,
});

function LMIAMapComponent() {
  return (
    <div className="lg:h-[calc(100vh-72px)]">
      <AuthNav />
      <LMIAMapHeatmap />
    </div>
  );
}
