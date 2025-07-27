interface Props {
  label: string;
  value: string;
}
const Stat = ({ label, value }: Props) => {
  return (
    <div className="flex items-center flex-col md:px-4">
      <p className="font-mono text-5xl mb-2 text-center">{value}</p>
      <p className="uppercase font-medium text-xs text-muted-foreground text-center">
        {label}
      </p>
    </div>
  );
};

export default Stat;
