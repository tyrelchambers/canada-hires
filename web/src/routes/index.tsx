import { buttonVariants } from "@/components/ui/button";
import { createFileRoute, Link } from "@tanstack/react-router";
import { AuthNav } from "@/components/AuthNav";
import { StripedBackground } from "@/components/StripedBackground";
import { TodaysJobs } from "@/components/TodaysJobs";
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

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: stats, isLoading: isLMIAStatsLoading } = useLMIAStats();
  const { data: statsData, isLoading: isJobStatsLoading } = useJobStats();
  const { data: reportStats, isLoading: isReportStatsLoading } =
    useReportStats();
  console.log(reportStats);
  return (
    <div className="min-h-screen">
      <AuthNav />

      <section className="bg-secondary border-b border-border relative h-[550px]">
        <StripedBackground />

        <div className="max-w-5xl mx-auto border-x border-border bg-white z-10 relative p-4 lg:p-20 h-full flex flex-col items-center justify-center">
          <h1 className="text-3xl lg:text-5xl -tracking-[0.015em] font-medium mb-6 text-center">
            Selling Out Canadian Jobs
          </h1>
          <p className="md:text-xl text-gray-500 font-light  text-center mb-10 max-w-3xl">
            The Temporary Foreign Worker (TFW) program is meant to fill labour
            shortages, but some companies exploit it to hire cheaper foreign
            labour instead of Canadians. We track the data so you can see which
            companies are abusing the system and choose where you spend your
            money.
          </p>

          <div className="flex md:flex-row flex-col gap-6 justify-center">
            <Link
              to="/lmia"
              className={clsx(
                buttonVariants({ variant: "outline", size: "lg" }),
              )}
            >
              Search LMIA records
            </Link>
            <Link
              to="/jobs"
              className={clsx(
                buttonVariants({ variant: "default", size: "lg" }),
              )}
            >
              Browse Job Postings
            </Link>
          </div>
        </div>
      </section>

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

      <section className="max-w-5xl mx-auto w-full my-20 px-4 ">
        <TodaysJobs />
      </section>

      <section className="grid grid-cols-1 lg:grid-cols-2 border-y border-border">
        <div className="flex flex-col gap-2 p-6 lg:p-20 max-w-3xl ml-auto border-r border-border">
          <Badge className="w-fit">Exploitation</Badge>
          <h2 className="font-bold text-3xl mb-4 -tracking-wide">
            How Companies Exploit the LMIA Program
          </h2>
          <p className=" text-lg text-muted-foreground">
            The Labour Market Impact Assessment (LMIA) program was designed to
            protect Canadian jobs by requiring employers to prove no qualified
            Canadian workers are available before hiring temporary foreign
            workers. However, many companies have found ways to circumvent this
            system's intent.
          </p>
          <p className=" text-lg text-muted-foreground">
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

      <section className="max-w-5xl mx-auto w-full my-20">
        <p className="lg:text-2xl p-4 font-light leading-relaxed text-center">
          Some employers use the LMIA process as a pathway to secure cheaper
          labor, posting positions they never intend to fill with Canadian
          workers. Others exploit the system's bureaucratic delays, using the
          application period to justify their "inability" to find local talent
          while their preferred temporary foreign worker waits in the wings.
        </p>
      </section>

      <section className="grid grid-cols-1 lg:grid-cols-2 border-y border-border">
        <div className="relative hidden lg:block">
          <StripedBackground />
        </div>
        <div className="flex flex-col gap-2 p-6 lg:p-20 max-w-3xl mr-auto border-l border-border">
          <Badge className="w-fit">Harming Canadians</Badge>

          <h2 className="font-bold text-3xl mb-4 -tracking-wide">
            LMIA exploitation is harming Canadians
          </h2>
          <p className=" text-lg text-muted-foreground">
            This exploitation undermines wage standards for all workers,
            displaces qualified Canadians from employment opportunities, and
            defeats the program's core purpose of protecting the domestic labor
            market. The result is a system that often serves corporate interests
            rather than Canadian workers or genuine labor market needs.
          </p>
        </div>
      </section>

      <section className="max-w-5xl mx-auto w-full my-20 p-4">
        <h2 className="text-3xl lg:text-4xl -tracking-wide text-center mb-6 font-medium">
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
        <div className="flex flex-col items-start">
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
