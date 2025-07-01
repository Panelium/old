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
    if (percentage > 90) return "bg-chart-1";
    if (percentage > 75) return "bg-chart-5";
    return "bg-chart-2";
  };

  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-xs">
        <span className="text-card-muted-foreground">{title}</span>
        <span className="font-medium text-card-muted-foreground">
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
