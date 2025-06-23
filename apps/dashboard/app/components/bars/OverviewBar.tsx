import React from "react";
import { getPercentage } from "~/lib/utils";
import Bar from "./Bar";

interface OverviewBarProps {
  title: string;
  value: number;
  uiValue?: string;
  max: number;
}

const OverviewBar: React.FC<OverviewBarProps> = ({
  title,
  uiValue,
  value,
  max,
}) => {
  const percentage = getPercentage(value, max);

  return (
    <div className="space-y-3 w-full">
      <div className="flex items-center justify-between text-xs">
        <span className="text-2xl font-bold text-slate-900 dark:text-white">
          {title}
        </span>
        <span className="font-normal text-slate-500 dark:text-slate-400">
          {uiValue}
        </span>
      </div>
      <Bar percentage={percentage} size="md" />
    </div>
  );
};

export default OverviewBar;
