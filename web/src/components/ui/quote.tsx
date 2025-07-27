import * as React from "react";

import { cn } from "@/lib/utils";

interface QuoteProps extends React.ComponentProps<"blockquote"> {
  source?: string;
  url?: string;
}

function Quote({ className, source, url, children, ...props }: QuoteProps) {
  return (
    <blockquote
      className={cn(
        "border-l-4 border-primary/60 bg-white px-6 py-4 italic text-muted-foreground shadow-sm my-6",
        className,
      )}
      {...props}
    >
      <div className="text-base leading-relaxed">{children}</div>
      {source && (
        <footer className="mt-3 text-sm font-medium text-foreground">
          {url ? (
            <a
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:underline"
            >
              — {source}
            </a>
          ) : (
            <span>— {source}</span>
          )}
        </footer>
      )}
    </blockquote>
  );
}

export { Quote };
