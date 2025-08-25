import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

interface MapLoadingOverlayProps {
  isLoading: boolean;
  loadingText?: string;
}

export function MapLoadingOverlay({
  isLoading,
  loadingText = "Loading data...",
}: MapLoadingOverlayProps) {
  if (!isLoading) return null;

  return (
    <div className="absolute inset-0 bg-black/20 bg-opacity-20 flex items-center justify-center z-10">
      <div className="bg-white rounded-lg p-4 shadow-lg flex items-center gap-3">
        <FontAwesomeIcon
          icon={faSpinner}
          className="animate-spin text-blue-600"
        />
        <span>{loadingText}</span>
      </div>
    </div>
  );
}
