export function StripedBackground({ className }: { className?: string }) {
  return (
    <div
      className={`absolute inset-0 opacity-10 ${className}`}
      style={{
        backgroundImage: `repeating-linear-gradient(
          45deg,
          #4a4e68 0px,
          #4a4e68 20px,
          transparent 20px,
          transparent 40px
        )`,
      }}
    />
  );
}
