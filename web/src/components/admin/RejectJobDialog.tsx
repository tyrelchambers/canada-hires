import { useState } from "react";
import { JobPosting } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTimes, faExternalLinkAlt } from "@fortawesome/free-solid-svg-icons";

interface RejectJobDialogProps {
  job: JobPosting;
  onReject: (jobId: string, reason?: string) => void;
  onCancel: () => void;
  isLoading: boolean;
}

export function RejectJobDialog({ job, onReject, onCancel, isLoading }: RejectJobDialogProps) {
  const [reason, setReason] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onReject(job.id, reason.trim() || undefined);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg">Reject Job for Reddit</CardTitle>
            <Button
              variant="ghost"
              size="sm"
              onClick={onCancel}
              disabled={isLoading}
            >
              <FontAwesomeIcon icon={faTimes} />
            </Button>
          </div>
        </CardHeader>
        
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {/* Job Details */}
            <div className="p-3 bg-gray-50 rounded-lg">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="font-medium text-gray-900">{job.title}</h3>
                  <p className="text-sm text-gray-600">{job.employer}</p>
                  <p className="text-sm text-gray-500">{job.location}</p>
                </div>
                <a
                  href={job.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:text-blue-800 ml-2"
                >
                  <FontAwesomeIcon icon={faExternalLinkAlt} />
                </a>
              </div>
            </div>

            {/* Rejection Reason */}
            <div className="space-y-2">
              <Label htmlFor="reason">
                Reason for rejection <span className="text-gray-500">(optional)</span>
              </Label>
              <Input
                id="reason"
                placeholder="e.g., Low quality posting, duplicate job, inappropriate content..."
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                disabled={isLoading}
              />
              <p className="text-xs text-gray-500">
                This reason will be logged for audit purposes.
              </p>
            </div>
          </CardContent>

          <CardFooter className="flex justify-end space-x-2">
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="destructive"
              disabled={isLoading}
            >
              {isLoading ? "Rejecting..." : "Reject Job"}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}