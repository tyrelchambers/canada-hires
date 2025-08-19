import { useState, useEffect, useMemo } from "react";
import { useActiveSubreddits } from "@/hooks/useSubreddits";
import { useGenerateRedditPosts } from "@/hooks/useAdminJobs";
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
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCheck,
  faTimes,
  faChevronDown,
  faChevronRight,
  faSpinner,
  faEdit,
} from "@fortawesome/free-solid-svg-icons";
import { faReddit } from "@fortawesome/free-brands-svg-icons";
import { GeneratedRedditPost } from "@/types";

interface ApprovalConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (
    selectedSubredditIds: string[],
    generatedContent?: GeneratedRedditPost[],
  ) => void;
  jobCount: number;
  jobIds: string[];
  isLoading?: boolean;
}

export function ApprovalConfirmationModal({
  isOpen,
  onClose,
  onConfirm,
  jobCount,
  jobIds,
  isLoading = false,
}: ApprovalConfirmationModalProps) {
  const { data: activeSubreddits, isLoading: subredditsLoading } =
    useActiveSubreddits();
  const generateContentMutation = useGenerateRedditPosts();

  const activeSubredditsList = useMemo(
    () => activeSubreddits?.subreddits || [],
    [activeSubreddits?.subreddits],
  );

  // State for per-post subreddit selection
  const [selectedSubreddits, setSelectedSubreddits] = useState<string[]>([]);

  // State for content generation
  const [generationStep, setGenerationStep] = useState<
    "subreddit-selection" | "generating" | "review"
  >("subreddit-selection");
  const [editablePosts, setEditablePosts] = useState<GeneratedRedditPost[]>([]);
  const [expandedPosts, setExpandedPosts] = useState<Set<string>>(new Set());

  // Initialize state when modal opens
  useEffect(() => {
    if (isOpen) {
      setGenerationStep("subreddit-selection");
      setEditablePosts([]);
      setExpandedPosts(new Set());
    } else {
      // Reset state when modal closes
      setSelectedSubreddits([]);
    }
  }, [isOpen]);

  // Initialize selected subreddits when subreddits data loads
  useEffect(() => {
    if (
      isOpen &&
      activeSubredditsList.length > 0 &&
      selectedSubreddits.length === 0
    ) {
      setSelectedSubreddits(activeSubredditsList.map((s) => s.id));
    }
  }, [isOpen, activeSubredditsList.length, selectedSubreddits.length]);

  const handleSubredditToggle = (subredditId: string) => {
    setSelectedSubreddits((prev) =>
      prev.includes(subredditId)
        ? prev.filter((id) => id !== subredditId)
        : [...prev, subredditId],
    );
  };

  const handleGenerateContent = async () => {
    setGenerationStep("generating");

    try {
      const result = await generateContentMutation.mutateAsync(jobIds);
      setEditablePosts(result.posts.map((post) => ({ ...post })));
      setGenerationStep("review");
    } catch (error) {
      console.error("Failed to generate content:", error);
      // Fall back to approval without generated content
      setGenerationStep("subreddit-selection");
    }
  };

  const handleConfirm = () => {
    if (generationStep === "review") {
      onConfirm(selectedSubreddits, editablePosts);
    } else {
      // Generate content first
      void handleGenerateContent();
    }
  };

  const handleFinalConfirm = () => {
    onConfirm(selectedSubreddits, editablePosts);
  };

  const handlePostContentChange = (jobId: string, newContent: string) => {
    setEditablePosts((prev) =>
      prev.map((post) =>
        post.job_id === jobId ? { ...post, content: newContent } : post,
      ),
    );
  };

  const togglePostExpansion = (jobId: string) => {
    setExpandedPosts((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(jobId)) {
        newSet.delete(jobId);
      } else {
        newSet.add(jobId);
      }
      return newSet;
    });
  };

  const selectedSubredditCount = selectedSubreddits.length;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-5xl w-full">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <FontAwesomeIcon icon={faReddit} className="text-orange-500" />
            <span>
              {generationStep === "subreddit-selection" &&
                "Configure Reddit Posting"}
              {generationStep === "generating" && "Generating Content..."}
              {generationStep === "review" && "Review Generated Content"}
            </span>
          </DialogTitle>
          <DialogDescription>
            {generationStep === "subreddit-selection" &&
              "Select subreddits and generate content for your job postings."}
            {generationStep === "generating" &&
              "Creating sarcastic Reddit posts for your job listings..."}
            {generationStep === "review" &&
              "Review and edit the generated content before posting."}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Job Count - Always visible */}
          <div className="flex items-center justify-between p-3 bg-blue-50 rounded-lg">
            <span className="font-medium text-blue-900">Jobs to post:</span>
            <Badge
              variant="secondary"
              className="bg-blue-100 text-blue-800 text-lg px-3 py-1"
            >
              {jobCount}
            </Badge>
          </div>

          {/* Step 1: Subreddit Selection */}
          {generationStep === "subreddit-selection" && (
            <div className="space-y-3">
              <div>
                <h4 className="font-medium text-gray-900">
                  Select subreddits ({selectedSubredditCount} selected):
                </h4>
                <p className="text-sm text-gray-600 mt-1">
                  💡 Toggle subreddits on/off for this batch of jobs
                </p>
              </div>

              {subredditsLoading ? (
                <div className="flex items-center space-x-2 text-gray-600">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-400"></div>
                  <span className="text-sm">Loading subreddits...</span>
                </div>
              ) : activeSubredditsList.length > 0 ? (
                <div className="space-y-2 max-h-32 overflow-y-auto">
                  {activeSubredditsList.map((subreddit) => {
                    const isSelected = selectedSubreddits.includes(
                      subreddit.id,
                    );
                    return (
                      <div
                        key={subreddit.id}
                        className={`flex items-center justify-between p-3 rounded-lg border transition-colors ${
                          isSelected
                            ? "bg-green-50 border-green-200"
                            : "bg-gray-50 border-gray-200"
                        }`}
                      >
                        <div className="flex items-center space-x-3">
                          <span className="font-mono text-sm font-medium">
                            r/{subreddit.name}
                          </span>
                          <Badge variant="outline" className="text-xs">
                            {subreddit.post_count} posts
                          </Badge>
                        </div>
                        <Switch
                          checked={isSelected}
                          onCheckedChange={() =>
                            handleSubredditToggle(subreddit.id)
                          }
                        />
                      </div>
                    );
                  })}
                </div>
              ) : (
                <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                  <p className="text-sm text-yellow-800">
                    ⚠️ No active subreddits configured. Jobs will be approved
                    but not posted to Reddit.
                  </p>
                </div>
              )}

              {selectedSubredditCount === 0 &&
                activeSubredditsList.length > 0 && (
                  <div className="p-3 bg-orange-50 border border-orange-200 rounded-lg">
                    <p className="text-sm text-orange-800">
                      ⚠️ No subreddits selected. Jobs will be approved but not
                      posted to Reddit.
                    </p>
                  </div>
                )}
            </div>
          )}

          {/* Step 2: Generating Content */}
          {generationStep === "generating" && (
            <div className="space-y-4">
              <div className="flex flex-col items-center justify-center p-8">
                <FontAwesomeIcon
                  icon={faSpinner}
                  className="text-4xl text-blue-500 animate-spin mb-4"
                />
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  Generating Reddit Posts...
                </h3>
                <p className="text-sm text-gray-600 text-center">
                  Gemini AI is creating engaging, sarcastic content for your{" "}
                  {jobCount} job{jobCount !== 1 ? "s" : ""}.<br />
                  This may take a few seconds.
                </p>
              </div>
            </div>
          )}

          {/* Step 3: Review Generated Content */}
          {generationStep === "review" && (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h4 className="font-medium text-gray-900">
                  Generated Content ({editablePosts.length} posts)
                </h4>
                <p className="text-sm text-gray-600">
                  Click posts to expand and edit content
                </p>
              </div>

              <div className="max-h-96 overflow-y-auto space-y-3">
                {editablePosts.map((post, index) => {
                  const isExpanded = expandedPosts.has(post.job_id);
                  const hasError = !!post.error;

                  return (
                    <div
                      key={post.job_id}
                      className={`border rounded-lg p-3 transition-colors ${
                        hasError
                          ? "border-red-200 bg-red-50"
                          : "border-gray-200 hover:border-blue-300"
                      }`}
                    >
                      <Collapsible>
                        <CollapsibleTrigger
                          onClick={() => togglePostExpansion(post.job_id)}
                          className="flex items-center justify-between w-full text-left"
                        >
                          <div className="flex items-center space-x-2">
                            <FontAwesomeIcon
                              icon={isExpanded ? faChevronDown : faChevronRight}
                              className="text-gray-400"
                            />
                            <span className="font-medium">
                              Post {index + 1}
                              {hasError && (
                                <span className="text-red-600 ml-2 text-sm">
                                  (Generation failed)
                                </span>
                              )}
                            </span>
                          </div>
                          <FontAwesomeIcon
                            icon={faEdit}
                            className="text-gray-400"
                          />
                        </CollapsibleTrigger>

                        <CollapsibleContent>
                          {hasError ? (
                            <div className="mt-3 p-3 bg-red-100 border border-red-200 rounded">
                              <p className="text-red-800 text-sm">
                                <strong>Error:</strong> {post.error}
                              </p>
                              <p className="text-red-700 text-xs mt-1">
                                This job will be posted without generated
                                content.
                              </p>
                            </div>
                          ) : (
                            <div className="mt-3">
                              <Textarea
                                value={post.content}
                                onChange={(e) =>
                                  handlePostContentChange(
                                    post.job_id,
                                    e.target.value,
                                  )
                                }
                                className="w-full min-h-32 resize-none"
                                placeholder="Generated content will appear here..."
                              />
                            </div>
                          )}
                        </CollapsibleContent>
                      </Collapsible>
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>

        <DialogFooter className="flex flex-row space-x-2">
          <Button
            variant="outline"
            onClick={onClose}
            disabled={isLoading || generateContentMutation.isPending}
            className="flex-1"
          >
            <FontAwesomeIcon icon={faTimes} className="mr-2" />
            Cancel
          </Button>

          {generationStep === "subreddit-selection" && (
            <>
              <Button
                onClick={handleConfirm}
                disabled={isLoading || generateContentMutation.isPending}
                className="flex-1 bg-blue-600 hover:bg-blue-700"
              >
                <FontAwesomeIcon icon={faSpinner} className="mr-2" />
                Generate Content
              </Button>
            </>
          )}

          {generationStep === "review" && (
            <Button
              onClick={handleFinalConfirm}
              disabled={isLoading}
              className="flex-1 bg-green-600 hover:bg-green-700"
            >
              <FontAwesomeIcon icon={faCheck} className="mr-2" />
              {isLoading
                ? "Posting..."
                : `Post ${jobCount} Job${jobCount !== 1 ? "s" : ""}`}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
