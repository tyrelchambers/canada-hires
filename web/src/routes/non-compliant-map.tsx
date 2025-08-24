import { createFileRoute } from "@tanstack/react-router";
import { NonCompliantMapHeatmap } from "@/components/NonCompliantMapHeatmap";
import { AuthNav } from "@/components/AuthNav";

export const Route = createFileRoute("/non-compliant-map")({
  component: NonCompliantMapComponent,
});

function NonCompliantMapComponent() {
  return (
    <div className="lg:h-screen overflow-hidden">
      <AuthNav />
      <NonCompliantMapHeatmap />
    </div>
  );
}
