export function StripedBackground() {
  return (
    <div
      className="absolute inset-0 opacity-10"
      style={{
        backgroundImage: `repeating-linear-gradient(
          45deg,
          #aaa 0px,
          #aaa 20px,
          transparent 20px,
          transparent 40px
        )`,
      }}
    />
  );
}
