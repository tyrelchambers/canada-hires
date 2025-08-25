import { Button, buttonVariants } from "@/components/ui/button";
import { createFileRoute, Link } from "@tanstack/react-router";
import { AuthNav } from "@/components/AuthNav";
import { StripedBackground } from "@/components/StripedBackground";
import { Footer } from "@/components/Footer";
import { ReportingBanner } from "@/components/ReportingBanner";
import clsx from "clsx";
import { useLMIAStats } from "@/hooks/useLMIA";
import { Badge } from "@/components/ui/badge";
import Stat, { StatSkeleton } from "@/components/Stat";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowRight, faFlask } from "@fortawesome/free-solid-svg-icons";
import { useJobStats } from "@/hooks/useJobPostings";
import { useReportStats } from "@/hooks/useReports";
import { TopBoycottedSection } from "@/components/TopBoycottedSection";
import Hero from "@/components/index/Hero";
import lmiaMap from "@/assets/lmia-map.png";

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: stats, isLoading: isLMIAStatsLoading } = useLMIAStats();
  const { data: statsData, isLoading: isJobStatsLoading } = useJobStats();
  const { data: reportStats, isLoading: isReportStatsLoading } =
    useReportStats();
  return (
    <div className="min-h-screen h-full">
      <AuthNav />

      <Hero />

      <ReportingBanner />

      <section className="min-h-[300px] items-center flex">
        <div className="max-w-5xl mx-auto flex flex-col flex-wrap md:flex-row py-10 gap-10 w-full justify-evenly">
          {isLMIAStatsLoading ? (
            <>
              <StatSkeleton />
              <StatSkeleton />
              <StatSkeleton />
            </>
          ) : (
            <>
              <Stat
                label="Distinct employers"
                value={String(
                  stats?.distinct_employers.toLocaleString() ?? "0",
                )}
              />

              <Stat
                label="Years"
                value={`${stats?.year_range.min_year}-${stats?.year_range.max_year}`}
              />

              <Stat
                label="Total records"
                value={String(stats?.total_records.toLocaleString() ?? "0")}
              />
            </>
          )}

          {isJobStatsLoading ? (
            <>
              <StatSkeleton />
              <StatSkeleton />
            </>
          ) : (
            <>
              <Stat
                label="Total Jobs"
                value={String(statsData?.total_jobs.toLocaleString() ?? "0")}
              />

              <Stat
                label="Total Employers"
                value={String(
                  statsData?.total_employers.toLocaleString() ?? "0",
                )}
              />
            </>
          )}

          {isReportStatsLoading ? (
            <StatSkeleton />
          ) : (
            <Stat
              label="Reports Submitted"
              value={String(reportStats?.total_reports.toLocaleString() ?? "0")}
            />
          )}
        </div>
      </section>

      <section className="grid grid-cols-1 lg:grid-cols-2 border-y border-border">
        <div className="flex flex-col gap-2 p-6 lg:p-20 max-w-3xl ml-auto border-r border-border">
          <Badge className="w-fit">Exploitation</Badge>
          <h2 className="font-bold text-3xl mb-4 -tracking-wide text-space-cadet">
            How Companies Exploit the LMIA Program
          </h2>
          <p className=" text-lg text-space-cadet">
            The Labour Market Impact Assessment (LMIA) program was designed to
            protect Canadian jobs by requiring employers to prove no qualified
            Canadian workers are available before hiring temporary foreign
            workers. However, many companies have found ways to circumvent this
            system's intent.
          </p>
          <p className=" text-lg text-space-cadet">
            Common exploitation tactics include posting job requirements with
            unrealistic qualifications, offering below-market wages that
            discourage Canadian applicants, advertising positions in obscure
            locations or with minimal visibility, and timing job postings
            strategically to fulfill technical requirements while ensuring few
            Canadians will apply.
          </p>
        </div>
        <div className="relative hidden lg:block">
          <StripedBackground />
        </div>
      </section>

      <section className="py-20">
        <section className="max-w-screen-2xl mx-auto w-full  flex-col flex gap-6 md:gap-20 items-center">
          <div className="flex flex-col gap-2 max-w-2xl p-4">
            <Badge>LMIA Map</Badge>
            <h2 className="font-bold text-3xl mb-4 -tracking-wide text-space-cadet mt-2">
              See where LMIA exploitation is happening
            </h2>
            <p className=" text-lg text-space-cadet">
              Our LMIA visualization platform transforms complex government
              employment data into an accessible interactive map, showing
              Temporary Foreign Worker program usage patterns across Canada.
              Users can explore real-time statistics on employer hiring
              practices, approved positions, and historical trends dating back
              to 2015.
            </p>
            <p className="text-lg text-space-cadet">
              Perfect for job seekers, researchers, and community advocates, the
              tool provides geographic insights and detailed employer
              information to help Canadians make informed decisions about local
              job markets and understand how immigration policies affect
              employment opportunities in their communities.
            </p>
            <Button asChild className="mt-8 w-fit">
              <Link to="/lmia-map">Go to full map</Link>
            </Button>
          </div>

          <div className="aspect-video md:h-[400px]  m-4 rounded-3xl shadow-xl overflow-hidden border-4 border-rose-quartz ">
            <img src={lmiaMap} className="object-cover w-full h-full" />
          </div>
        </section>
      </section>

      <section className="grid grid-cols-1 lg:grid-cols-2 border-y border-border">
        <div className="relative hidden lg:block">
          <StripedBackground />
        </div>
        <div className="flex flex-col gap-2 p-6 lg:p-20 max-w-3xl mr-auto border-l border-border">
          <Badge className="w-fit">Harming Canadians</Badge>

          <h2 className="font-bold text-3xl mb-4 -tracking-wide text-space-cadet">
            LMIA exploitation is harming Canadians
          </h2>
          <p className=" text-lg text-space-cadet">
            This exploitation undermines wage standards for all workers,
            displaces qualified Canadians from employment opportunities, and
            defeats the program's core purpose of protecting the domestic labor
            market. The result is a system that often serves corporate interests
            rather than Canadian workers or genuine labor market needs.
          </p>
        </div>
      </section>

      <section className="max-w-5xl mx-auto w-full my-20 p-4">
        <h2 className="text-3xl lg:text-4xl -tracking-wide text-center mb-6 font-medium text-space-cadet">
          The goal of JobWatch Canada
        </h2>
        <p className="lg:text-2xl font-light leading-relaxed text-center">
          JobWatch Canada aims to bring transparency to these practices by
          tracking which employers frequently rely on the TFW program and
          providing Canadians with the information needed to make informed
          decisions about where to spend their money.
        </p>
      </section>

      <section className="my-20 mx-4 lg:max-w-5xl lg:w-full lg:mx-auto p-8 lg:p-20 bg-gradient-to-r from-gray-900 to-gray-800 rounded-xl shadow-xl shadow-gray-900/60 relative">
        <StripedBackground className="mask-l-from-0%" />
        <div className="flex flex-col items-start z-20 relative">
          <FontAwesomeIcon icon={faFlask} className="text-3xl text-white" />
          <h3 className="lg:text-4xl font-medium -tracking-wide max-w-lg text-white my-4">
            See how JobWatch Canada is tracking LMIA exploitation
          </h3>
          <p className="text-sm lg:text-lg text-white/70 mb-8 ">
            JobWatch Canada is tracking LMIA exploitation by analyzing data on
            temporary foreign worker applications and comparing them to the
            domestic labor market. This information is then made available to
            Canadians through our website and social media channels.
          </p>
          <Link
            to="/research"
            className={clsx(buttonVariants({ variant: "default" }))}
          >
            <FontAwesomeIcon icon={faArrowRight} />
            See research
          </Link>
        </div>
      </section>

      <TopBoycottedSection />

      <Footer />
    </div>
  );
}
