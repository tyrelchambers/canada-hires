import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faExclamationTriangle } from "@fortawesome/free-solid-svg-icons";
import { useNonCompliantReasons } from "@/hooks/useNonCompliant";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

interface ReasonCodeModalProps {
  reasonCodes: string[];
  businessName: string;
  isOpen: boolean;
  onClose: () => void;
}

export function ReasonCodeModal({
  reasonCodes,
  businessName,
  isOpen,
  onClose,
}: ReasonCodeModalProps) {
  const { data: reasonsData, isLoading } = useNonCompliantReasons();

  // Filter reasons to only those matching the provided codes
  const relevantReasons =
    reasonsData?.reasons.filter((reason) =>
      reasonCodes.includes(reason.reason_code),
    ) || [];

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh]" showCloseButton={false}>
        <DialogHeader className="bg-orange-600 text-white p-4 -m-6 mb-4 rounded-t-lg">
          <div className="flex items-center justify-between">
            <DialogTitle className="flex items-center gap-2 text-lg font-semibold text-white">
              <FontAwesomeIcon
                icon={faExclamationTriangle}
                className="text-white"
              />
              Violation Details
            </DialogTitle>
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="text-white hover:text-gray-200 hover:bg-white/20 h-auto w-auto p-1"
            >
              <span className="sr-only">Close</span>
              ✕
            </Button>
          </div>
        </DialogHeader>

        <div className="overflow-y-auto max-h-[calc(90vh-200px)]">
          <div className="mb-4">
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              {businessName}
            </h3>
            <DialogDescription>
              This employer has been found non-compliant with the following
              regulation(s):
            </DialogDescription>
          </div>

          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <div className="text-gray-500">Loading violation details...</div>
            </div>
          ) : relevantReasons.length === 0 ? (
            <div className="text-center py-8">
              <div className="text-gray-500 mb-2">
                No detailed descriptions available for reason codes:{" "}
                {reasonCodes.join(", ")}
              </div>
              <div className="text-xs text-gray-400">
                The system is still processing violation descriptions. Please
                try again later.
              </div>
            </div>
          ) : (
            <div className="space-y-6">
              {relevantReasons.map((reason, index) => (
                <div
                  key={reason.id}
                  className="border-l-4 border-orange-500 pl-4"
                >
                  <div className="flex items-start gap-2 mb-2">
                    <Badge variant="secondary" className="bg-orange-100 text-orange-800 hover:bg-orange-200">
                      Code {reason.reason_code}
                    </Badge>
                  </div>
                  <div className="text-gray-800 leading-relaxed">
                    {reason.description || "Description not available"}
                  </div>
                  {index < relevantReasons.length - 1 && (
                    <div className="border-b border-gray-200 mt-4 mb-2"></div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        <DialogFooter className="bg-gray-50 px-6 py-4 -m-6 mt-4 rounded-b-lg">
          <div className="flex justify-between items-center w-full">
            <div className="text-xs text-gray-500">
              Data source: Immigration, Refugees and Citizenship Canada
            </div>
            <Button variant="secondary" onClick={onClose}>
              Close
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
