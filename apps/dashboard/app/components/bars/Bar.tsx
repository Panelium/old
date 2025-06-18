import { cn } from "~/lib/utils";

interface BarProps {
  percentage: number;
  barColor: string;
  size: "sm" | "md" | "lg";
  className?: string;
}

const Bar: React.FC<BarProps> = ({
  percentage,
  barColor,
  className,
  size = "md",
}) => {
  const sizeClasses = {
    sm: "h-1",
    md: "h-1.5",
    lg: "h-3",
  };

  return (
    <div
      className={cn(
        sizeClasses[size],
        "w-full overflow-hidden rounded-full",
        "bg-slate-200 dark:bg-slate-800",
        className
      )}
    >
      <div
        className={cn("h-full rounded-full", barColor)}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
};

export default Bar;
