import { useActiveSubreddits } from "@/hooks/useSubreddits";
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
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faTimes } from "@fortawesome/free-solid-svg-icons";
import { faReddit } from "@fortawesome/free-brands-svg-icons";

interface ApprovalConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  jobCount: number;
  isLoading?: boolean;
}

export function ApprovalConfirmationModal({
  isOpen,
  onClose,
  onConfirm,
  jobCount,
  isLoading = false,
}: ApprovalConfirmationModalProps) {
  const { data: activeSubreddits, isLoading: subredditsLoading } =
    useActiveSubreddits();

  const activeSubredditsList = activeSubreddits?.subreddits || [];

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <FontAwesomeIcon icon={faReddit} className="text-orange-500" />
            <span>Confirm Reddit Posting</span>
          </DialogTitle>
          <DialogDescription>
            You are about to approve and post job listings to Reddit.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Job Count */}
          <div className="flex items-center justify-between p-3 bg-blue-50 rounded-lg">
            <span className="font-medium text-blue-900">Jobs to post:</span>
            <Badge
              variant="secondary"
              className="bg-blue-100 text-blue-800 text-lg px-3 py-1"
            >
              {jobCount}
            </Badge>
          </div>

          {/* Subreddits */}
          <div className="space-y-2">
            <h4 className="font-medium text-gray-900">Will be posted to:</h4>
            {subredditsLoading ? (
              <div className="flex items-center space-x-2 text-gray-600">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-400"></div>
                <span className="text-sm">Loading subreddits...</span>
              </div>
            ) : activeSubredditsList.length > 0 ? (
              <div className="space-y-2">
                {activeSubredditsList.map((subreddit) => (
                  <div
                    key={subreddit.id}
                    className="flex items-center justify-between p-2 bg-gray-50 rounded"
                  >
                    <span className="font-mono text-sm">
                      r/{subreddit.name}
                    </span>
                    <Badge variant="outline" className="text-xs">
                      {subreddit.post_count} posts
                    </Badge>
                  </div>
                ))}
              </div>
            ) : (
              <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                <p className="text-sm text-yellow-800">
                  ⚠️ No active subreddits configured. Jobs will be approved but
                  not posted to Reddit.
                </p>
              </div>
            )}
          </div>

          {activeSubredditsList.length > 0 && (
            <div className="text-xs text-gray-500 bg-gray-50 p-2 rounded">
              💡 Jobs will be posted automatically after approval
            </div>
          )}
        </div>

        <DialogFooter className="flex flex-row space-x-2">
          <Button
            variant="outline"
            onClick={onClose}
            disabled={isLoading}
            className="flex-1"
          >
            <FontAwesomeIcon icon={faTimes} className="mr-2" />
            Cancel
          </Button>
          <Button
            onClick={onConfirm}
            disabled={isLoading}
            className="flex-1 bg-green-600 hover:bg-green-700"
          >
            <FontAwesomeIcon icon={faCheck} className="mr-2" />
            {isLoading
              ? "Posting..."
              : `Post ${jobCount} Job${jobCount !== 1 ? "s" : ""}`}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
