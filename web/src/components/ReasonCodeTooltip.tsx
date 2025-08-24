import { useState } from "react";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useNonCompliantReasons } from "@/hooks/useNonCompliant";

interface ReasonCodeTooltipProps {
  reasonCode: string;
  children: React.ReactNode;
  className?: string;
}

export function ReasonCodeTooltip({ reasonCode, children, className }: ReasonCodeTooltipProps) {
  const [open, setOpen] = useState(false);
  const { data: reasonsData, isLoading } = useNonCompliantReasons();

  // Find the reason description for this code
  const reason = reasonsData?.reasons.find(r => r.reason_code === reasonCode);
  
  // Don't show tooltip if no description found
  if (!reason?.description || isLoading) {
    return <span className={className}>{children}</span>;
  }

  // Truncate long descriptions for tooltip display
  const shortDescription = reason.description.length > 200 
    ? reason.description.substring(0, 200) + "..." 
    : reason.description;

  return (
    <Tooltip open={open} onOpenChange={setOpen}>
      <TooltipTrigger asChild>
        <span 
          className={`cursor-help underline decoration-dotted underline-offset-2 hover:text-orange-700 ${className || ""}`}
          onMouseEnter={() => setOpen(true)}
          onMouseLeave={() => setOpen(false)}
        >
          {children}
        </span>
      </TooltipTrigger>
      <TooltipContent 
        side="top" 
        className="max-w-sm p-3 text-xs bg-gray-900 text-white border border-gray-700"
        sideOffset={8}
      >
        <div className="space-y-1">
          <div className="font-semibold text-orange-400">Reason {reasonCode}:</div>
          <div className="leading-relaxed">{shortDescription}</div>
          {reason.description.length > 200 && (
            <div className="text-gray-400 text-xs italic">
              Click for full details...
            </div>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
}