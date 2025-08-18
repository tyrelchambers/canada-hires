import { useState } from "react";
import { ReportStatsCards } from "./ReportStatsCards";
import { ReportTable } from "./ReportTable";
import { ReportDetailsModal } from "./ReportDetailsModal";
import { type Report } from "@/hooks/useReports";

export function ReportManagementDashboard() {
  const [selectedReport, setSelectedReport] = useState<Report | null>(null);
  const [modalMode, setModalMode] = useState<"view" | "edit">("view");
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleViewReport = (report: Report) => {
    setSelectedReport(report);
    setModalMode("view");
    setIsModalOpen(true);
  };

  const handleEditReport = (report: Report) => {
    setSelectedReport(report);
    setModalMode("edit");
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedReport(null);
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="border-b pb-4">
        <h1 className="text-3xl font-bold text-gray-900">Report Management</h1>
        <p className="mt-2 text-gray-600">
          Manage and review community-submitted business reports
        </p>
      </div>

      {/* Statistics Cards */}
      <ReportStatsCards />

      {/* Reports Table */}
      <ReportTable 
        onViewReport={handleViewReport}
        onEditReport={handleEditReport}
      />

      {/* Report Details Modal */}
      <ReportDetailsModal
        report={selectedReport}
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        mode={modalMode}
      />
    </div>
  );
}