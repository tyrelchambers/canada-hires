import { useState } from "react";
import { Subreddit } from "@/types";
import {
  useSubreddits,
  useCreateSubreddit,
  useUpdateSubreddit,
  useDeleteSubreddit,
} from "@/hooks/useSubreddits";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faEdit,
  faTrash,
  faSave,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { faReddit } from "@fortawesome/free-brands-svg-icons";

interface EditingSubreddit {
  id: string;
  is_active: boolean;
}

export function SubredditManager() {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editingData, setEditingData] = useState<EditingSubreddit | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [newSubreddit, setNewSubreddit] = useState({
    name: "",
    is_active: true,
  });

  const { data: subredditsData, isLoading, error } = useSubreddits();

  // Sort subreddits by name to maintain consistent order
  const subreddits = subredditsData
    ? {
        ...subredditsData,
        subreddits: [...subredditsData.subreddits].sort((a, b) =>
          a.name.localeCompare(b.name),
        ),
      }
    : undefined;
  const createMutation = useCreateSubreddit();
  const updateMutation = useUpdateSubreddit();
  const deleteMutation = useDeleteSubreddit();

  const handleToggleActive = async (subreddit: Subreddit) => {
    try {
      await updateMutation.mutateAsync({
        id: subreddit.id,
        data: { is_active: !subreddit.is_active },
      });
    } catch (error) {
      console.error("Failed to update subreddit:", error);
    }
  };

  const handleStartEdit = (subreddit: Subreddit) => {
    setEditingId(subreddit.id);
    setEditingData({
      id: subreddit.id,
      is_active: subreddit.is_active,
    });
  };

  const handleSaveEdit = async () => {
    if (!editingData) return;

    try {
      await updateMutation.mutateAsync({
        id: editingData.id,
        data: {
          is_active: editingData.is_active,
        },
      });
      setEditingId(null);
      setEditingData(null);
    } catch (error) {
      console.error("Failed to save subreddit:", error);
    }
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditingData(null);
  };

  const handleCreate = async () => {
    if (!newSubreddit.name.trim()) return;

    try {
      await createMutation.mutateAsync({
        name: newSubreddit.name.trim(),
        is_active: newSubreddit.is_active,
      });

      // Reset form
      setNewSubreddit({
        name: "",
        is_active: true,
      });
      setShowAddForm(false);
    } catch (error) {
      console.error("Failed to create subreddit:", error);
    }
  };

  const handleDelete = async (id: string) => {
    if (
      !confirm(
        "Are you sure you want to delete this subreddit? This action cannot be undone.",
      )
    ) {
      return;
    }

    try {
      await deleteMutation.mutateAsync(id);
    } catch (error) {
      console.error("Failed to delete subreddit:", error);
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "Never";
    return new Date(dateString).toLocaleDateString();
  };

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-2 text-gray-600">Loading subreddits...</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-red-600">
            <p className="font-medium">Error loading subreddits</p>
            <p className="text-sm text-gray-600 mt-1">
              Please check your permissions and try again.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center space-x-2">
                <FontAwesomeIcon icon={faReddit} className="text-orange-500" />
                <span>Subreddit Management</span>
              </CardTitle>
              <p className="text-sm text-gray-600 mt-1">
                Track and manage which existing subreddits to post job listings to
              </p>
            </div>
            <Button onClick={() => setShowAddForm(true)} disabled={showAddForm}>
              <FontAwesomeIcon icon={faPlus} className="mr-2" />
              Track Subreddit
            </Button>
          </div>
        </CardHeader>
      </Card>

      {/* Add New Subreddit Form */}
      {showAddForm && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Track New Subreddit</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label htmlFor="name">Subreddit Name *</Label>
              <Input
                id="name"
                placeholder="e.g., jobwatchcanada (without r/)"
                value={newSubreddit.name}
                onChange={(e) =>
                  setNewSubreddit((prev) => ({
                    ...prev,
                    name: e.target.value,
                  }))
                }
              />
            </div>
            <div className="flex items-center space-x-2">
              <Switch
                id="is_active"
                checked={newSubreddit.is_active}
                onCheckedChange={(checked) =>
                  setNewSubreddit((prev) => ({ ...prev, is_active: checked }))
                }
              />
              <Label htmlFor="is_active">
                Enabled (job postings will be published to this subreddit)
              </Label>
            </div>
            <div className="flex space-x-2">
              <Button
                onClick={handleCreate}
                disabled={!newSubreddit.name.trim() || createMutation.isPending}
              >
                <FontAwesomeIcon icon={faSave} className="mr-2" />
                {createMutation.isPending ? "Adding..." : "Add to List"}
              </Button>
              <Button variant="outline" onClick={() => setShowAddForm(false)}>
                <FontAwesomeIcon icon={faTimes} className="mr-2" />
                Cancel
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Subreddits Table */}
      <Card>
        <CardHeader>
          <CardTitle>
            Tracked Subreddits ({subreddits?.subreddits.length || 0})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Subreddit</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Posts</TableHead>
                <TableHead>Last Posted</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {subreddits?.subreddits.map((subreddit) => (
                <TableRow key={subreddit.id}>
                  <TableCell className="font-medium">
                    r/{subreddit.name}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      <Switch
                        checked={
                          editingId === subreddit.id
                            ? editingData?.is_active || false
                            : subreddit.is_active
                        }
                        onCheckedChange={(checked) => {
                          if (editingId === subreddit.id) {
                            setEditingData((prev) =>
                              prev ? { ...prev, is_active: checked } : null,
                            );
                          } else {
                            void handleToggleActive(subreddit);
                          }
                        }}
                        disabled={updateMutation.isPending}
                      />
                      <Badge
                        variant={subreddit.is_active ? "default" : "secondary"}
                      >
                        {subreddit.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">{subreddit.post_count}</Badge>
                  </TableCell>
                  <TableCell className="text-sm">
                    {formatDate(subreddit.last_posted_at)}
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end space-x-2">
                      {editingId === subreddit.id ? (
                        <>
                          <Button
                            size="sm"
                            onClick={handleSaveEdit}
                            disabled={updateMutation.isPending}
                          >
                            <FontAwesomeIcon icon={faSave} />
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={handleCancelEdit}
                          >
                            <FontAwesomeIcon icon={faTimes} />
                          </Button>
                        </>
                      ) : (
                        <>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleStartEdit(subreddit)}
                          >
                            <FontAwesomeIcon icon={faEdit} />
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleDelete(subreddit.id)}
                            disabled={deleteMutation.isPending}
                            className="text-red-600 border-red-300 hover:bg-red-50"
                          >
                            <FontAwesomeIcon icon={faTrash} />
                          </Button>
                        </>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {subreddits?.subreddits.length === 0 && (
            <div className="text-center py-12">
              <FontAwesomeIcon
                icon={faReddit}
                className="mx-auto text-4xl text-gray-400"
              />
              <h3 className="mt-4 text-lg font-medium text-gray-900">
                No subreddits tracked
              </h3>
              <p className="mt-2 text-gray-600">
                Add existing subreddits to your posting list to start publishing
                job listings to Reddit.
              </p>
              <Button className="mt-4" onClick={() => setShowAddForm(true)}>
                <FontAwesomeIcon icon={faPlus} className="mr-2" />
                Track First Subreddit
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
