import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";

interface QuarterSelectorProps {
  year: number;
  quarter?: string;
  onYearChange: (year: number) => void;
  onQuarterChange: (quarter?: string) => void;
}

export function QuarterSelector({
  year,
  quarter,
  onYearChange,
  onQuarterChange,
}: QuarterSelectorProps) {
  const currentYear = new Date().getFullYear();
  
  // Generate years from 2020 to current year
  const years = Array.from(
    { length: currentYear - 2019 },
    (_, i) => currentYear - i,
  );

  const quarters = [
    { value: "all", label: "All Quarters" },
    { value: "Q1", label: "Q1 (Jan-Mar)" },
    { value: "Q2", label: "Q2 (Apr-Jun)" },
    { value: "Q3", label: "Q3 (Jul-Sep)" },
    { value: "Q4", label: "Q4 (Oct-Dec)" },
  ];

  return (
    <div className="flex flex-col space-y-4 p-4 bg-white rounded-lg shadow">
      <div>
        <Label htmlFor="year-select" className="text-sm font-medium">
          Year
        </Label>
        <Select
          value={year.toString()}
          onValueChange={(value) => onYearChange(parseInt(value))}
        >
          <SelectTrigger id="year-select" className="w-full">
            <SelectValue placeholder="Select year" />
          </SelectTrigger>
          <SelectContent>
            {years.map((y) => (
              <SelectItem key={y} value={y.toString()}>
                {y}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div>
        <Label htmlFor="quarter-select" className="text-sm font-medium">
          Quarter
        </Label>
        <Select value={quarter || "all"} onValueChange={(value) => onQuarterChange(value === "all" ? undefined : value)}>
          <SelectTrigger id="quarter-select" className="w-full">
            <SelectValue placeholder="Select quarter" />
          </SelectTrigger>
          <SelectContent>
            {quarters.map((q) => (
              <SelectItem key={q.value} value={q.value}>
                {q.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="text-xs text-gray-500">
        <p>Filter LMIA data by year and quarter.</p>
        <p>Showing data for {year} {quarter || "all quarters"}</p>
      </div>
    </div>
  );
}