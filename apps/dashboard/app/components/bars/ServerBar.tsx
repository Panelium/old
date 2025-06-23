import React from "react";
import { getPercentage } from "~/lib/utils";
import Bar from "./Bar";

interface ResourceBarProps {
  title: string;
  value: number;
  uiValue?: string;
  max: number;
}

const ServerBar: React.FC<ResourceBarProps> = ({
  title,
  uiValue,
  value,
  max,
}) => {
  const percentage = getPercentage(value, max);

  const getBarColor = (percentage: number) => {
    if (percentage > 90) return "bg-rose-500";
    if (percentage > 75) return "bg-amber-500";
    return "bg-emerald-500";
  };

  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-xs">
        <span className="text-slate-600 dark:text-slate-400">{title}</span>
        <span className="font-medium text-slate-700 dark:text-slate-300">
          {uiValue}
        </span>
      </div>
      <Bar
        percentage={percentage}
        barColor={getBarColor(percentage)}
        size="md"
      />
    </div>
  );
};

export default ServerBar;
