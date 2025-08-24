import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

interface MapLoadingSpinnerProps {
  text?: string;
  size?: "sm" | "md";
}

export function MapLoadingSpinner({ 
  text = "Loading...", 
  size = "md" 
}: MapLoadingSpinnerProps) {
  return (
    <div className={`flex items-center gap-2 text-gray-500 ${size === "sm" ? "text-sm" : ""}`}>
      <FontAwesomeIcon
        icon={faSpinner}
        className="animate-spin"
      />
      <span>{text}</span>
    </div>
  );
}