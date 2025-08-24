interface MapBackButtonProps {
  onBack: () => void;
}

export function MapBackButton({ onBack }: MapBackButtonProps) {
  return (
    <div className="mb-4">
      <button
        onClick={onBack}
        className="flex items-center gap-2 text-gray-600 hover:text-gray-900 text-sm font-medium"
      >
        ← Back to Map Overview
      </button>
    </div>
  );
}