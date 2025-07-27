import { createFileRoute } from "@tanstack/react-router";
import { AuthNav } from "@/components/AuthNav";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
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
  faExclamationTriangle,
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
          "Canadian Wage Subsidy Programs for Hiring Immigrants (2025) - JobWatch Canada",
        description:
          "Comprehensive guide to federal, provincial, and territorial wage subsidy programs available to Canadian businesses hiring immigrants and newcomers in 2025.",
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
            Canadian Wage Subsidy Programs for Hiring Immigrants and Newcomers
          </h1>
          <p className="text-xl text-gray-600 mb-6">
            Comprehensive analysis of federal, provincial, and territorial
            programs providing wage subsidies to businesses hiring immigrants
            (2025)
          </p>

          <div className="flex items-center gap-2 text-sm text-gray-500">
            <FontAwesomeIcon icon={faCalendar} className="w-4 h-4" />
            <span>Last updated: January 2025</span>
          </div>
        </div>

       {/* Executive Summary */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FontAwesomeIcon icon={faInfoCircle} className="w-5 h-5" />
              Executive Summary
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p>
              Canada offers numerous wage subsidy programs that can benefit
              businesses hiring immigrants and newcomers. These programs range
              from 15% to 80% wage subsidies, with durations from 4 months to 1
              year depending on the specific program.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-blue-50 p-4 rounded-lg">
                <h4 className="font-semibold text-blue-900">
                  Federal Programs
                </h4>
                <p className="text-sm text-blue-700">
                  Summer jobs, student placements, skills training
                </p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <h4 className="font-semibold text-green-900">
                  Provincial Programs
                </h4>
                <p className="text-sm text-green-700">
                  Integration programs, apprenticeships, employment support
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg">
                <h4 className="font-semibold text-purple-900">
                  Territorial Programs
                </h4>
                <p className="text-sm text-purple-700">
                  Remote area support, skills development, workforce integration
                </p>
              </div>
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
            <CardTitle>Key Takeaways for Employers</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-semibold text-green-700 mb-2">
                  Available Support
                </h4>
                <ul className="space-y-1 text-sm">
                  <li>
                    • Multiple funding streams available across all regions
                  </li>
                  <li>• Subsidies range from 15% to 80% of wages</li>
                  <li>• Additional support for training and integration</li>
                  <li>
                    • Programs target skills shortages and underrepresented
                    groups
                  </li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold text-blue-700 mb-2">
                  Important Considerations
                </h4>
                <ul className="space-y-1 text-sm">
                  <li>• Programs are not exclusive to immigrants</li>
                  <li>• Eligibility requirements vary by region</li>
                  <li>• Some programs have limited funding/waitlists</li>
                  <li>• 2025 sees reduced immigration allocations</li>
                </ul>
              </div>
            </div>

            <div className="mt-6 p-4 bg-gray-100 rounded-lg">
              <p className="text-sm text-gray-700">
                <strong>Disclaimer:</strong> This information is based on
                publicly available sources as of January 2025. Program details,
                funding amounts, and eligibility criteria may change. Employers
                should contact relevant government agencies directly to confirm
                current program availability and requirements.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
