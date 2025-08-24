import { ReactNode } from "react";

interface MapMobileSectionProps {
  title: string;
  description: string;
  children: ReactNode;
  footer?: ReactNode;
}

export function MapMobileSection({ title, description, children, footer }: MapMobileSectionProps) {
  return (
    <div className="lg:hidden bg-gray-50 border-t border-gray-200">
      <div className="p-4 bg-white border-b border-gray-200">
        <h1 className="text-lg font-bold text-gray-900">{title}</h1>
        <p className="text-sm text-gray-600 mt-1">{description}</p>
      </div>

      <div className="p-4 space-y-4">
        {children}
      </div>

      {footer && (
        <div className="px-4 pb-2 text-xs text-gray-400">
          {footer}
        </div>
      )}
    </div>
  );
}