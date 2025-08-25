import { ReactNode } from "react";
import { MapContainer as LeafletMapContainer, TileLayer } from "react-leaflet";
import { LatLngExpression } from "leaflet";

interface MapContainerProps {
  children: ReactNode;
  mapRef: (map: L.Map | null) => void;
  center?: LatLngExpression;
  zoom?: number;
  className?: string;
}

export function MapContainer({
  children,
  mapRef,
  center = [61.0666922, -95.712891],
  zoom = 4,
  className = "h-[300px] md:h-[400px] lg:h-full lg:flex-1 z-0",
}: MapContainerProps) {
  return (
    <div className={className}>
      <LeafletMapContainer
        center={center}
        zoom={zoom}
        style={{ height: "100%", width: "100%", zIndex: 0 }}
        ref={mapRef}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {children}
      </LeafletMapContainer>
    </div>
  );
}
