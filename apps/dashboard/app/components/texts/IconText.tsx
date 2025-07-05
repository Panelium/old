import type { LucideIcon } from "lucide-react";

interface IconTextProps {
  icon: LucideIcon;
  text: string;
}

export default function IconText({ icon, text }: IconTextProps) {
  const IconComponent = icon;
  return (
    <div className="flex items-center text-sm text-card-muted-foreground">
      <IconComponent className="h-3.5 w-3.5 mr-1.5" />
      {text}
    </div>
  );
}
