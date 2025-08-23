import { createFileRoute } from '@tanstack/react-router'
import { LMIAMapHeatmap } from '@/components/LMIAMapHeatmap'

export const Route = createFileRoute('/lmia-map')({
  component: LMIAMapComponent,
})

function LMIAMapComponent() {
  return (
    <div className="h-screen overflow-hidden">
      <LMIAMapHeatmap />
    </div>
  )
}