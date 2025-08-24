import { createFileRoute } from "@tanstack/react-router";
import { LMIAMapHeatmap } from "@/components/LMIAMapHeatmap";
import { AuthNav } from "@/components/AuthNav";

export const Route = createFileRoute("/lmia-map")({
  component: LMIAMapComponent,
});

function LMIAMapComponent() {
  return (
    <div className="h-screen overflow-hidden">
      <AuthNav />
      <LMIAMapHeatmap />
    </div>
  );
}
