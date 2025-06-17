import React from "react";
import { cn } from "~/lib/utils";

interface BarProps {
  value: number;
  max: number;
  className?: string;
}

interface ResourceBarProps {
  title: string;
  value: number;
  uiValue?: string;
  max: number;
}

const Bar: React.FC<BarProps> = ({ value, max, className }) => {
  const percentage = Math.min((value / max) * 100, 100);

  const getBarColor = () => {
    if (percentage > 90) return "bg-rose-500";
    if (percentage > 75) return "bg-amber-500";
    return "bg-emerald-500";
  };

  return (
    <div
      className={cn(
        "h-1.5 w-full overflow-hidden rounded-full",
        "bg-slate-200 dark:bg-slate-800",
        className
      )}
    >
      <div
        className={cn("h-full rounded-full", getBarColor())}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
};

const ServerBar: React.FC<ResourceBarProps> = ({
  title,
  uiValue,
  value,
  max,
}) => {
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-xs">
        <span className="text-slate-600 dark:text-slate-400">{title}</span>
        <span className="font-medium text-slate-700 dark:text-slate-300">
          {uiValue}
        </span>
      </div>
      <Bar value={value} max={max} />
    </div>
  );
};

export default ServerBar;
