import { createFileRoute } from "@tanstack/react-router";
import { JobPostings } from "@/components/JobPostings";
import { AuthNav } from "@/components/AuthNav";
import { DataDisclaimer } from "@/components/DataDisclaimer";

export const Route = createFileRoute("/jobs")({
  component: JobsPage,
  head: () => ({
    meta: [
      {
        title: "LMIA Job Postings - Canada Hires",
        description:
          "Browse current LMIA job postings from Canadian employers. Find opportunities and apply for positions that require Labour Market Impact Assessment approval.",
      },
    ],
  }),
});

function JobsPage() {
  return (
    <section>
      <AuthNav />
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900">
              LMIA Job Postings
            </h1>
            <p className="mt-2 text-lg text-gray-600">
              Current job postings from Canadian employers that have received
              Labour Market Impact Assessment (LMIA) approval for hiring foreign
              workers. Canadian citizens and permanent residents are encouraged
              to apply.
            </p>
          </div>
          <JobPostings />
          <DataDisclaimer className="mt-4" />
        </div>
      </div>
    </section>
  );
}
