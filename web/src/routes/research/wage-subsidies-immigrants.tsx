import { createFileRoute } from "@tanstack/react-router";
import { AuthNav } from "@/components/AuthNav";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faInfoCircle,
  faMapMarkerAlt,
  faCalendar,
  faDollarSign,
} from "@fortawesome/free-solid-svg-icons";

export const Route = createFileRoute("/research/wage-subsidies-immigrants")({
  component: WageSubsidiesPage,
  head: () => ({
    meta: [
      {
        title:
          "How Canadian Wage Subsidy Programs Work: What Canadians Need to Know (2025) - JobWatch Canada",
        description:
          "Understanding how your tax dollars fund wage subsidy programs for immigrant hiring across Canada. A comprehensive guide to federal, provincial, and territorial programs affecting Canadian workers and businesses.",
      },
    ],
  }),
});

function WageSubsidiesPage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <AuthNav />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Understanding Canadian Wage Subsidy Programs: How They Work and What
            They Mean for Canadians
          </h1>
          <p className="text-xl text-gray-600 mb-6">
            A comprehensive guide to how Canadian tax dollars fund wage
            subsidies for immigrant hiring and what this means for Canadian
            workers and communities (2025)
          </p>

          <div className="flex items-center gap-2 text-sm text-gray-500">
            <FontAwesomeIcon icon={faCalendar} className="w-4 h-4" />
            <span>Last updated: July 2025</span>
          </div>
        </div>

        {/* What This Means for Canadians */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon icon={faInfoCircle} className="w-5 h-5" />
              What This Means for Canadians
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p>
              Canadian taxpayers fund numerous wage subsidy programs that
              provide financial incentives to businesses hiring immigrants and
              newcomers. These publicly-funded programs cover 15% to 80% of
              wages paid by employers, with durations from 4 months to 1 year.
              Understanding these programs helps Canadians see how their tax
              dollars are being used and what impact they have on the Canadian
              job market.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-blue-50 p-4 rounded-lg">
                <h4 className="font-semibold text-blue-900">
                  Federal Programs
                </h4>
                <p className="text-sm text-blue-700">
                  Federal tax dollars fund summer jobs, student placements, and
                  skills training programs
                </p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <h4 className="font-semibold text-green-900">
                  Provincial Programs
                </h4>
                <p className="text-sm text-green-700">
                  Provincial tax dollars fund integration programs,
                  apprenticeships, and employment support
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg">
                <h4 className="font-semibold text-purple-900">
                  Territorial Programs
                </h4>
                <p className="text-sm text-purple-700">
                  Territorial tax dollars fund remote area support, skills
                  development, and workforce integration
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Understanding Wage Subsidies */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon icon={faDollarSign} className="w-5 h-5" />
              How Wage Subsidies Work and Impact Canadian Workers
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="bg-blue-50 p-6 rounded-lg">
              <h3 className="font-semibold text-blue-900 mb-3">
                What Are Wage Subsidies?
              </h3>
              <p className="text-blue-800">
                Wage subsidies are government payments that cover a portion of
                an employee's salary, paid directly to employers who hire from
                specific groups. The employer pays the full wage to the worker,
                then receives reimbursement from the government for the
                subsidized portion.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="bg-green-50 p-4 rounded-lg">
                <h4 className="font-semibold text-green-900 mb-2">
                  Potential Benefits for Canadians
                </h4>
                <ul className="text-sm text-green-800 space-y-1">
                  <li>• Skills integration helps fill labour shortages</li>
                  <li>
                    • Economic growth through increased workforce participation
                  </li>
                  <li>• Knowledge transfer and innovation</li>
                  <li>• Regional development in areas with worker shortages</li>
                </ul>
              </div>

              <div className="bg-yellow-50 p-4 rounded-lg">
                <h4 className="font-semibold text-yellow-900 mb-2">
                  Considerations for Canadian Workers
                </h4>
                <ul className="text-sm text-yellow-800 space-y-1">
                  <li>
                    • May create financial incentives to hire subsidized workers
                    over others
                  </li>
                  <li>• Programs are time-limited (typically 4-12 months)</li>
                  <li>• Focus on specific sectors or skill areas</li>
                  <li>• Transparency in program outcomes varies by region</li>
                </ul>
              </div>
            </div>

            <div className="bg-gray-100 p-4 rounded-lg">
              <h4 className="font-semibold text-gray-900 mb-2">
                Public Funding Transparency
              </h4>
              <p className="text-sm text-gray-700">
                These programs are funded through federal, provincial, and
                territorial budgets using Canadian tax dollars. While program
                details are publicly available, comprehensive data on outcomes,
                job displacement effects, and long-term employment retention is
                often limited or not consistently published across
                jurisdictions.
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Quick Reference Table */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>
              Quick Reference: Wage Subsidy Programs by Region
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Region/Program</TableHead>
                    <TableHead>Subsidy Amount</TableHead>
                    <TableHead>Duration</TableHead>
                    <TableHead>Target Group</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow>
                    <TableCell className="font-medium">Quebec PRIIME</TableCell>
                    <TableCell>60-70% (capped at min wage)</TableCell>
                    <TableCell>40-52 weeks</TableCell>
                    <TableCell>Immigrants/visible minorities</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">BC WorkBC</TableCell>
                    <TableCell>50% → 25% → 15%</TableCell>
                    <TableCell>6 months</TableCell>
                    <TableCell>Unemployed residents</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">
                      Yukon Staffing UP
                    </TableCell>
                    <TableCell>60% of base wage</TableCell>
                    <TableCell>Up to 1 year</TableCell>
                    <TableCell>Under-represented groups</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">Ontario TOP</TableCell>
                    <TableCell>50-70% (max $5K-$7K)</TableCell>
                    <TableCell>Work placement period</TableCell>
                    <TableCell>Students including newcomers</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell className="font-medium">
                      NL Training Program
                    </TableCell>
                    <TableCell>60-80% ($12/hour max)</TableCell>
                    <TableCell>28 weeks</TableCell>
                    <TableCell>Unemployed residents</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>

        {/* Federal Programs */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Federal Programs</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h3 className="text-lg font-semibold mb-3">
                Canada Summer Jobs Program
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p>
                  <strong>Eligibility:</strong> Canadian citizens, permanent
                  residents, and protected refugees aged 18-30
                </p>
                <p>
                  <strong>Subsidy:</strong> 50% to 100% of minimum hourly wage
                </p>
                <p>
                  <strong>Duration:</strong> Up to 4 months
                </p>
                <p>
                  <strong>Purpose:</strong> Create summer employment
                  opportunities for youth
                </p>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-3">
                Student Work Placement Program
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p>
                  <strong>Subsidy:</strong> Up to $5,000 for regular students,
                  up to $7,000 for under-represented groups (including
                  newcomers)
                </p>
                <p>
                  <strong>Eligibility:</strong> Students enrolled in Canadian
                  institutions where work placement is part of study plan
                </p>
                <p>
                  <strong>Requirements:</strong> Must be Canadian citizens,
                  permanent residents, or persons with refugee protection
                </p>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-3">
                Skills for Success Program
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p>
                  <strong>Purpose:</strong> Foundational and transferable skills
                  training
                </p>
                <p>
                  <strong>Streams:</strong> Training and Tools Stream, Research
                  and Innovation Stream
                </p>
                <p>
                  <strong>Focus:</strong> All skill levels, includes support for
                  newcomers
                </p>
                <p>
                  <strong>Funding:</strong> Variable based on project scope and
                  participants
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Provincial Programs */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Provincial Programs</CardTitle>
          </CardHeader>
          <CardContent className="space-y-8">
            {/* Quebec */}
            <div>
              <div className="flex items-center gap-2 mb-4">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 text-blue-600"
                />
                <h3 className="text-xl font-semibold">Quebec</h3>
              </div>

              <div className="ml-6">
                <h4 className="text-lg font-semibold mb-3">
                  Employment Integration Program for Immigrants and Visible
                  Minorities (PRIIME)
                </h4>
                <div className="bg-blue-50 p-4 rounded-lg space-y-2">
                  <p>
                    <strong>Subsidy:</strong> Up to 60% of gross salary (capped
                    at minimum wage), up to 70% for individuals with
                    disabilities
                  </p>
                  <p>
                    <strong>Duration:</strong> Up to 40-52 weeks depending on
                    integration needs
                  </p>
                  <p>
                    <strong>Additional Support:</strong>
                  </p>
                  <ul className="list-disc list-inside ml-4 space-y-1">
                    <li>
                      Accompaniment grant up to $2,000 (100% of accompanist
                      costs)
                    </li>
                    <li>HR adaptation reimbursement up to $5,000</li>
                  </ul>
                  <p>
                    <strong>Eligibility:</strong> Immigrants/visible minorities
                    with no North American work experience in their field
                  </p>
                </div>
              </div>
            </div>

            {/* British Columbia */}
            <div>
              <div className="flex items-center gap-2 mb-4">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 text-green-600"
                />
                <h3 className="text-xl font-semibold">British Columbia</h3>
              </div>

              <div className="ml-6">
                <h4 className="text-lg font-semibold mb-3">
                  WorkBC Wage Subsidy
                </h4>
                <div className="bg-green-50 p-4 rounded-lg space-y-2">
                  <p>
                    <strong>Structure:</strong> Three-tier system over 6 months
                  </p>
                  <ul className="list-disc list-inside ml-4 space-y-1">
                    <li>First third: 50% wage subsidy (up to $500 weekly)</li>
                    <li>Second third: 25% wage subsidy</li>
                    <li>Final third: 15% wage subsidy</li>
                  </ul>
                  <p>
                    <strong>Eligibility:</strong> Permanent residents and
                    Canadian citizens facing unemployment
                  </p>
                </div>
              </div>
            </div>

            {/* Ontario */}
            <div>
              <div className="flex items-center gap-2 mb-4">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 text-red-600"
                />
                <h3 className="text-xl font-semibold">Ontario</h3>
              </div>

              <div className="ml-6 space-y-4">
                <div>
                  <h4 className="text-lg font-semibold mb-3">
                    Talent Opportunities Program (TOP)
                  </h4>
                  <div className="bg-red-50 p-4 rounded-lg space-y-2">
                    <p>
                      <strong>Subsidy:</strong> Up to 50% of wages (max $5,000)
                      or 70% (max $7,000 for under-represented groups including
                      newcomers)
                    </p>
                    <p>
                      <strong>Focus:</strong> College and university student
                      work placements
                    </p>
                  </div>
                </div>

                <div>
                  <h4 className="text-lg font-semibold mb-3">
                    Graduated Apprenticeship Grant for Employers (GAGE)
                  </h4>
                  <div className="bg-red-50 p-4 rounded-lg space-y-2">
                    <p>
                      <strong>Subsidy:</strong> Up to $16,700 to train
                      apprentices
                    </p>
                    <p>
                      <strong>Bonus:</strong> Additional $2,500 for
                      under-represented groups including newcomers
                    </p>
                    <p>
                      <strong>Trades:</strong> More than 100 eligible trades
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Newfoundland and Labrador */}
            <div>
              <div className="flex items-center gap-2 mb-4">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 text-purple-600"
                />
                <h3 className="text-xl font-semibold">
                  Newfoundland and Labrador
                </h3>
              </div>

              <div className="ml-6">
                <h4 className="text-lg font-semibold mb-3">
                  Provincial Training Program
                </h4>
                <div className="bg-purple-50 p-4 rounded-lg space-y-2">
                  <p>
                    <strong>Subsidy:</strong> Up to $12/hour for 28 weeks within
                    42-week period
                  </p>
                  <p>
                    <strong>Structure:</strong> 60% subsidy first 14 weeks, 80%
                    subsidy last 14 weeks
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Territorial Programs */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Territorial Programs</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 text-yellow-600"
                />
                <h3 className="text-xl font-semibold">Yukon Territory</h3>
              </div>

              <div className="ml-6">
                <h4 className="text-lg font-semibold mb-3">
                  Staffing UP Program
                </h4>
                <div className="bg-yellow-50 p-4 rounded-lg space-y-2">
                  <p>
                    <strong>Subsidy:</strong> 60% of base wage for up to one
                    year
                  </p>
                  <p>
                    <strong>Additional Support:</strong>
                  </p>
                  <ul className="list-disc list-inside ml-4 space-y-1">
                    <li>On-the-job training costs up to 6 months</li>
                    <li>Skill assessments up to 8 weeks</li>
                  </ul>
                  <p>
                    <strong>Target:</strong> Under-represented groups including
                    newcomers
                  </p>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-center gap-2 mb-4">
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 text-indigo-600"
                />
                <h3 className="text-xl font-semibold">Nunavut</h3>
              </div>

              <div className="ml-6">
                <h4 className="text-lg font-semibold mb-3">
                  Training on the Job Program
                </h4>
                <div className="bg-indigo-50 p-4 rounded-lg space-y-2">
                  <p>
                    <strong>Purpose:</strong> Encourage hiring of unemployed
                    EI-eligible individuals
                  </p>
                  <p>
                    <strong>Structure:</strong> Partial salary subsidy for
                    full-time, part-time (min 20 hrs/week), and seasonal jobs
                  </p>
                  <p>
                    <strong>Note:</strong> Nunavut currently has no immigrant
                    nominee program
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Sector-Specific Programs */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Sector-Specific Programs</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h3 className="text-lg font-semibold mb-3">
                Electricity Human Resources Canada
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p>
                  <strong>Subsidy:</strong> Up to 50% wage subsidy or $10,000
                  maximum
                </p>
                <p>
                  <strong>Duration:</strong> First few months of onboarding
                </p>
                <p>
                  <strong>Target:</strong> Skilled newcomers in electricity
                  sector
                </p>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-3">
                Pathways to Employment for Newcomers
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p>
                  <strong>Purpose:</strong> Wage subsidies for hiring skilled
                  newcomers
                </p>
                <p>
                  <strong>Additional:</strong> Upskilling opportunities and
                  integration best practices
                </p>
                <p>
                  <strong>Focus:</strong> Internationally educated professionals
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* 2025 Context */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>2025 Policy Context and Challenges</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h3 className="text-lg font-semibold mb-3 text-red-600">
                Reduced Provincial Nominee Program Allocations
              </h3>
              <div className="bg-red-50 p-4 rounded-lg space-y-2">
                <p>
                  Most provinces experienced a 50% reduction in PNP allocations
                  for 2025:
                </p>
                <ul className="list-disc list-inside ml-4 space-y-1">
                  <li>
                    <strong>Saskatchewan:</strong> Reduced to 3,625 spots
                    (lowest since 2009)
                  </li>
                  <li>
                    <strong>Manitoba:</strong> Advocating for 12,000 allocations
                    to meet labour demands
                  </li>
                  <li>
                    <strong>New Brunswick:</strong> Paused Atlantic Immigration
                    Program after reaching quota
                  </li>
                </ul>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-3 text-green-600">
                Minimum Wage Increases (October 2025)
              </h3>
              <div className="bg-green-50 p-4 rounded-lg space-y-2">
                <p>Five provinces implementing minimum wage increases:</p>
                <ul className="list-disc list-inside ml-4 space-y-1">
                  <li>
                    <strong>Saskatchewan:</strong> $15.35/hour (from $15.00)
                  </li>
                  <li>
                    <strong>New Brunswick:</strong> Expected ~$15.77/hour (April
                    2025)
                  </li>
                  <li>
                    <strong>Other provinces:</strong> Ontario, Manitoba, Nova
                    Scotia, Prince Edward Island
                  </li>
                </ul>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-3">
                Immigration Targets
              </h3>
              <div className="bg-blue-50 p-4 rounded-lg space-y-2">
                <p>
                  <strong>2025 Target:</strong> 395,000 new permanent residents
                </p>
                <p>
                  <strong>Future Reductions:</strong> 380,000 (2026), 365,000
                  (2027)
                </p>
                <p>
                  <strong>Atlantic Provinces:</strong> ~80% of immigrants
                  welcomed through subnational programs in 2022
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Conclusion */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Key Information for Canadians</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-semibold text-green-700 mb-2">
                  How Your Tax Dollars Are Used
                </h4>
                <ul className="space-y-1 text-sm">
                  <li>
                    • Multiple publicly-funded programs across all regions
                  </li>
                  <li>
                    • Government covers 15% to 80% of wages for participating
                    employers
                  </li>
                  <li>
                    • Additional funding for training and integration support
                  </li>
                  <li>
                    • Programs aim to address regional labour shortages and
                    skill gaps
                  </li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold text-blue-700 mb-2">
                  What Canadians Should Know
                </h4>
                <ul className="space-y-1 text-sm">
                  <li>
                    • Programs create temporary hiring incentives (4-12 months
                    typically)
                  </li>
                  <li>
                    • Eligibility requirements and oversight vary by region
                  </li>
                  <li>
                    • Some programs have limited funding and may have waitlists
                  </li>
                  <li>
                    • 2025 sees reduced immigration targets affecting program
                    demand
                  </li>
                </ul>
              </div>
            </div>

            <div className="mt-6 p-4 bg-gray-100 rounded-lg">
              <p className="text-sm text-gray-700">
                <strong>Disclaimer:</strong> This information is based on
                publicly available sources as of July 2025. Program details,
                funding amounts, and eligibility criteria may change. Canadians
                interested in learning more should contact relevant government
                agencies directly to confirm current program availability,
                outcomes data, and transparency measures.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
