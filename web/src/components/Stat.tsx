import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Skeleton } from "./ui/skeleton";

interface Props {
  label: string;
  value: string;
  icon?: IconDefinition;
}
const Stat = ({ label, value, icon }: Props) => {
  return (
    <div className="flex items-center flex-col md:px-4">
      {icon && <FontAwesomeIcon icon={icon} className="text-2xl mb-2" />}
      <p className="font-mono text-5xl mb-2 text-center">{value}</p>
      <p className="uppercase font-medium text-xs text-muted-foreground text-center">
        {label}
      </p>
    </div>
  );
};

export const StatSkeleton = () => {
  return (
    <div className="flex items-center flex-col md:px-4">
      <Skeleton className="h-12 w-32 mb-2" />
      <Skeleton className="h-4 w-24" />
    </div>
  );
};

export default Stat;
