import type { LucideIcon } from "lucide-react";
import { cn } from "~/lib/utils";
import React from "react";

interface IconTextProps {
  icon: LucideIcon;
  text: string;
  copy?: boolean; // Optional prop to indicate if text should be copyable
}

export default function IconText({ icon, text, copy }: IconTextProps) {
  const IconComponent = icon;
  const [copied, setCopied] = React.useState(false);

  const handleCopy = () => {
    if (copy) {
      navigator.clipboard
        .writeText(text)
        .then(() => {
          setCopied(true);
          setTimeout(() => setCopied(false), 1200);
        })
        .catch((err) => {
          console.error("Failed to copy text: ", err);
        });
    }
  };

  return (
    <div
      className={cn(
        "flex items-center text-sm text-card-muted-foreground",
        copy ? "cursor-pointer hover:text-card-foreground" : ""
      )}
      onClick={handleCopy}
    >
      <IconComponent className="h-3.5 w-3.5 mr-1.5" />
      <span>{text}</span>
      {copy && (
        <span
          className={cn("ml-2 text-xs transition-opacity duration-200", copied ? "opacity-100" : "opacity-0 size-0")}
          aria-live="polite"
        >
          Copied!
        </span>
      )}
    </div>
  );
}
