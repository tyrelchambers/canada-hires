import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

interface PageLoaderProps {
  text?: string;
}

export function PageLoader({ 
  text = "Loading..." 
}: PageLoaderProps) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <FontAwesomeIcon
          icon={faSpinner}
          className="animate-spin text-4xl text-gray-400 mb-4"
        />
        <p className="text-gray-600 text-lg">{text}</p>
      </div>
    </div>
  );
}