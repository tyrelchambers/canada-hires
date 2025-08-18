import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { useUpdateReport, type Report } from "@/hooks/useReports";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { 
  faCalendar, 
  faUser, 
  faMapMarkerAlt,
  faSave,
  faTimes
} from "@fortawesome/free-solid-svg-icons";

interface ReportDetailsModalProps {
  report: Report | null;
  isOpen: boolean;
  onClose: () => void;
  mode: "view" | "edit";
}

export function ReportDetailsModal({ 
  report, 
  isOpen, 
  onClose, 
  mode 
}: ReportDetailsModalProps) {
  const [isEditing, setIsEditing] = useState(mode === "edit");
  const [formData, setFormData] = useState({
    business_name: "",
    business_address: "",
    report_source: "employment" as "employment" | "observation" | "public_record",
    confidence_level: undefined as number | undefined,
    additional_notes: "",
  });

  const updateReportMutation = useUpdateReport();

  // Initialize form data when report changes
  useEffect(() => {
    if (report) {
      setFormData({
        business_name: report.business_name,
        business_address: report.business_address,
        report_source: report.report_source as "employment" | "observation" | "public_record",
        confidence_level: report.confidence_level,
        additional_notes: report.additional_notes || "",
      });
    }
    setIsEditing(mode === "edit");
  }, [report, mode]);

  const handleSave = async () => {
    if (!report) return;

    try {
      await updateReportMutation.mutateAsync({
        id: report.id,
        data: {
          business_name: formData.business_name,
          business_address: formData.business_address,
          report_source: formData.report_source,
          confidence_level: formData.confidence_level,
          additional_notes: formData.additional_notes || undefined,
        },
      });
      setIsEditing(false);
      onClose();
    } catch (error) {
      console.error("Failed to update report:", error);
      alert("Failed to update report. Please try again.");
    }
  };

  const handleCancel = () => {
    if (report) {
      // Reset form data to original values
      setFormData({
        business_name: report.business_name,
        business_address: report.business_address,
        report_source: report.report_source as "employment" | "observation" | "public_record",
        confidence_level: report.confidence_level,
        additional_notes: report.additional_notes || "",
      });
    }
    setIsEditing(false);
    if (mode === "edit") {
      onClose();
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-CA', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getReportSourceBadge = (source: string) => {
    const variants = {
      employment: "bg-blue-100 text-blue-800",
      observation: "bg-green-100 text-green-800",
      public_record: "bg-purple-100 text-purple-800"
    };
    
    const labels = {
      employment: "Employment",
      observation: "Observation",
      public_record: "Public Record"
    };

    return (
      <Badge className={variants[source as keyof typeof variants] || "bg-gray-100 text-gray-800"}>
        {labels[source as keyof typeof labels] || source}
      </Badge>
    );
  };

  if (!report) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center justify-between">
            <span>
              {isEditing ? "Edit Report" : "Report Details"}
            </span>
            {!isEditing && mode === "view" && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsEditing(true)}
              >
                Edit
              </Button>
            )}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          {/* Report Metadata */}
          <div className="bg-gray-50 p-4 rounded-lg">
            <h3 className="font-semibold text-sm text-gray-700 mb-3">Report Information</h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div className="flex items-center text-gray-600">
                <FontAwesomeIcon icon={faUser} className="w-4 h-4 mr-2" />
                <span>User ID: {report.user_id.slice(0, 8)}...</span>
              </div>
              <div className="flex items-center text-gray-600">
                <FontAwesomeIcon icon={faCalendar} className="w-4 h-4 mr-2" />
                <span>Created: {formatDate(report.created_at)}</span>
              </div>
              {report.updated_at !== report.created_at && (
                <div className="flex items-center text-gray-600">
                  <FontAwesomeIcon icon={faCalendar} className="w-4 h-4 mr-2" />
                  <span>Updated: {formatDate(report.updated_at)}</span>
                </div>
              )}
            </div>
          </div>

          {/* Business Information */}
          <div className="space-y-4">
            <div>
              <Label htmlFor="business_name">Business Name</Label>
              {isEditing ? (
                <Input
                  id="business_name"
                  value={formData.business_name}
                  onChange={(e) => setFormData(prev => ({ ...prev, business_name: e.target.value }))}
                  className="mt-1"
                />
              ) : (
                <div className="mt-1 p-2 bg-gray-50 rounded-md font-medium">
                  {report.business_name}
                </div>
              )}
            </div>

            <div>
              <Label htmlFor="business_address">Business Address</Label>
              {isEditing ? (
                <Input
                  id="business_address"
                  value={formData.business_address}
                  onChange={(e) => setFormData(prev => ({ ...prev, business_address: e.target.value }))}
                  className="mt-1"
                />
              ) : (
                <div className="mt-1 p-2 bg-gray-50 rounded-md flex items-center">
                  <FontAwesomeIcon icon={faMapMarkerAlt} className="w-4 h-4 mr-2 text-gray-500" />
                  {report.business_address}
                </div>
              )}
            </div>

            <div>
              <Label htmlFor="report_source">Report Source</Label>
              {isEditing ? (
                <select
                  id="report_source"
                  value={formData.report_source}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    report_source: e.target.value as "employment" | "observation" | "public_record"
                  }))}
                  className="mt-1 flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="employment">Employment</option>
                  <option value="observation">Observation</option>
                  <option value="public_record">Public Record</option>
                </select>
              ) : (
                <div className="mt-1">
                  {getReportSourceBadge(report.report_source)}
                </div>
              )}
            </div>

            <div>
              <Label htmlFor="confidence_level">Confidence Level (1-10)</Label>
              {isEditing ? (
                <Input
                  id="confidence_level"
                  type="number"
                  min="1"
                  max="10"
                  value={formData.confidence_level || ""}
                  onChange={(e) => setFormData(prev => ({ 
                    ...prev, 
                    confidence_level: e.target.value ? parseInt(e.target.value) : undefined 
                  }))}
                  className="mt-1"
                  placeholder="Optional"
                />
              ) : (
                <div className="mt-1 p-2 bg-gray-50 rounded-md">
                  {report.confidence_level ? (
                    <span className="font-semibold">{report.confidence_level}/10</span>
                  ) : (
                    <span className="text-gray-500">Not specified</span>
                  )}
                </div>
              )}
            </div>

            <div>
              <Label htmlFor="additional_notes">Additional Notes</Label>
              {isEditing ? (
                <Textarea
                  id="additional_notes"
                  value={formData.additional_notes}
                  onChange={(e) => setFormData(prev => ({ ...prev, additional_notes: e.target.value }))}
                  className="mt-1"
                  rows={4}
                  placeholder="Optional additional information..."
                />
              ) : (
                <div className="mt-1 p-2 bg-gray-50 rounded-md min-h-[100px]">
                  {report.additional_notes || (
                    <span className="text-gray-500">No additional notes</span>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>

        <DialogFooter>
          {isEditing ? (
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={handleCancel}
                disabled={updateReportMutation.isPending}
              >
                <FontAwesomeIcon icon={faTimes} className="w-4 h-4 mr-2" />
                Cancel
              </Button>
              <Button
                onClick={handleSave}
                disabled={updateReportMutation.isPending}
              >
                <FontAwesomeIcon icon={faSave} className="w-4 h-4 mr-2" />
                {updateReportMutation.isPending ? "Saving..." : "Save Changes"}
              </Button>
            </div>
          ) : (
            <Button variant="outline" onClick={onClose}>
              Close
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}