import { ReactNode } from "react";

interface MapSidebarProps {
  title: string;
  description: string;
  children: ReactNode;
  footer?: ReactNode;
}

export function MapSidebar({ title, description, children, footer }: MapSidebarProps) {
  return (
    <div className="hidden lg:flex w-[400px] bg-gray-50 border-r border-gray-200 flex-col overflow-hidden">
      <div className="p-4 border-b border-gray-200 bg-white">
        <h1 className="text-xl font-bold text-gray-900">{title}</h1>
        <p className="text-sm text-gray-600 mt-1">{description}</p>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
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