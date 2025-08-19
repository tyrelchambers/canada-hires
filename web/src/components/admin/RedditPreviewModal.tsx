import { JobPosting } from "@/types";
import { usePreviewRedditPost } from "@/hooks/useAdminJobs";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faRobot,
  faFileAlt,
  faExclamationTriangle,
} from "@fortawesome/free-solid-svg-icons";

interface RedditPreviewModalProps {
  job: JobPosting;
  isOpen: boolean;
  onClose: () => void;
}

export function RedditPreviewModal({
  job,
  isOpen,
  onClose,
}: RedditPreviewModalProps) {
  const previewMutation = usePreviewRedditPost();

  const handleGeneratePreview = () => {
    previewMutation.mutate(job.id);
  };

  const preview = previewMutation.data;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="!max-w-4xl w-full max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <span>Reddit Post Preview</span>
            <Badge variant="outline" className="text-xs">
              {job.title}
            </Badge>
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          {/* Job Info */}
          <div className="bg-gray-50 p-4 rounded-lg">
            <h3 className="font-semibold text-sm text-gray-700 mb-2">
              Job Details
            </h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="font-medium">Title:</span> {job.title}
              </div>
              <div>
                <span className="font-medium">Employer:</span> {job.employer}
              </div>
              <div>
                <span className="font-medium">Location:</span> {job.location}
              </div>
              <div className="col-span-2">
                <div className="mb-1">
                  <span className="font-medium">URL:</span>
                </div>
                <a
                  href={job.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:underline text-xs break-all block"
                  title={job.url}
                >
                  {job.url}
                </a>
              </div>
            </div>
          </div>

          {/* Preview Controls */}
          <div className="flex items-center justify-between">
            <h3 className="font-semibold">What will be posted to Reddit:</h3>
            {!preview && (
              <Button
                onClick={handleGeneratePreview}
                disabled={previewMutation.isPending}
                className="bg-blue-600 hover:bg-blue-700"
              >
                {previewMutation.isPending
                  ? "Generating..."
                  : "Generate Preview"}
              </Button>
            )}
          </div>

          {/* Loading State */}
          {previewMutation.isPending && (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-2 text-gray-600">Generating preview...</p>
            </div>
          )}

          {/* Error State */}
          {previewMutation.error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <div className="flex items-center space-x-2">
                <FontAwesomeIcon
                  icon={faExclamationTriangle}
                  className="text-red-500"
                />
                <span className="font-medium text-red-800">
                  Failed to generate preview
                </span>
              </div>
              <p className="text-red-600 text-sm mt-1">
                {previewMutation.error.message}
              </p>
            </div>
          )}

          {/* Preview Content */}
          {preview && (
            <div className="space-y-4">
              {/* Content Type Badge */}
              <div className="flex items-center space-x-2">
                <Badge
                  variant={
                    preview.content_type === "ai" ? "default" : "secondary"
                  }
                  className={
                    preview.content_type === "ai"
                      ? "bg-green-600"
                      : "bg-gray-600"
                  }
                >
                  <FontAwesomeIcon
                    icon={preview.content_type === "ai" ? faRobot : faFileAlt}
                    className="mr-1"
                  />
                  {preview.content_type === "ai" ? "AI Generated" : "Template"}
                </Badge>
                {preview.error && (
                  <span className="text-sm text-amber-600 flex items-center">
                    <FontAwesomeIcon
                      icon={faExclamationTriangle}
                      className="mr-1"
                    />
                    {preview.error}
                  </span>
                )}
              </div>

              {/* Reddit Post Preview */}
              <div className="border border-gray-300 rounded-lg overflow-hidden bg-white">
                {/* Reddit-style header */}
                <div className="bg-gray-100 px-4 py-2 border-b border-gray-200">
                  <div className="flex items-center space-x-2 text-sm text-gray-600">
                    <span className="font-medium">r/jobwatchcanada</span>
                    <span>•</span>
                    <span>Posted by u/JobWatchBot</span>
                    <span>•</span>
                    <span>just now</span>
                  </div>
                </div>

                {/* Post content */}
                <div className="p-4">
                  <h2 className="text-lg font-semibold text-gray-900 mb-3">
                    {preview.title}
                  </h2>
                  <div className="prose prose-sm max-w-none">
                    <pre className="whitespace-pre-wrap font-sans text-gray-800 leading-relaxed">
                      {preview.body}
                    </pre>
                  </div>
                </div>
              </div>

              {/* Actions */}
              <div className="flex items-center justify-between pt-4">
                <Button
                  variant="outline"
                  onClick={handleGeneratePreview}
                  disabled={previewMutation.isPending}
                >
                  Regenerate Preview
                </Button>
                <div className="flex space-x-2">
                  <Button variant="outline" onClick={onClose}>
                    Close
                  </Button>
                </div>
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
