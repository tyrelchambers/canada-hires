import { ReactElement } from "react";

interface Props {
  label: string;
  value: number;
  icon: () => ReactElement;
}
const Stat = ({ label, value, icon }: Props) => {
  return (
    <div className="p-6 w-full bg-card rounded-xl ">
      <header className="flex items-center gap-8 justify-between  ">
        <div className="flex flex-col">
          <h3 className="text-muted-foreground">{label}</h3>
          <p className="text-3xl font-bold">{value}</p>
        </div>
        <div className="bg-background h-12 w-12 justify-center flex items-center rounded-full ">
          {icon()}
        </div>
      </header>
    </div>
  );
};

export default Stat;
