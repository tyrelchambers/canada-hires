import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

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

export default Stat;
