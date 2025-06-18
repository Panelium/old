import React from "react";
import { getPercentage } from "~/lib/utils";
import Bar from "./Bar";

interface OverviewBarProps {
  title: string;
  value: number;
  uiValue?: string;
  max: number;
}

const ServerBar: React.FC<OverviewBarProps> = ({
  title,
  uiValue,
  value,
  max,
}) => {
  const percentage = getPercentage(value, max);

  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-xs">
        <span className="text-slate-600 dark:text-slate-400">{title}</span>
        <span className="font-medium text-slate-700 dark:text-slate-300">
          {uiValue}
        </span>
      </div>
      <Bar percentage={percentage} size="lg" />
    </div>
  );
};

export default ServerBar;
