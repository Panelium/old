import React from "react";
import { motion } from "framer-motion";
import { cn } from "~/lib/utils";

export interface ResourceGaugeProps {
  value: number;
  maxValue: number;
  subtitle?: string;
  label?: string;
  unit?: string;
  size?: "sm" | "md" | "lg";
  showValue?: boolean;
  className?: string;
}

export function ResourceGauge({
  value,
  maxValue,
  label,
  unit = "",
  size = "md",
  showValue = true,
  className,
}: ResourceGaugeProps) {
  const percentage = Math.min((value / maxValue) * 100, 100);

  const getColor = () => {
    if (percentage > 90) return "text-red-500";
    if (percentage > 75) return "text-amber-500";
    return "text-emerald-500";
  };

  const getSize = () => {
    switch (size) {
      case "sm":
        return "w-24 h-24";
      case "lg":
        return "w-40 h-40";
      default:
        return "w-32 h-32";
    }
  };

  const getThickness = () => {
    switch (size) {
      case "sm":
        return 4;
      case "lg":
        return 8;
      default:
        return 6;
    }
  };

  const getFontSize = () => {
    switch (size) {
      case "sm":
        return "text-lg";
      case "lg":
        return "text-3xl";
      default:
        return "text-2xl";
    }
  };

  // Calculate stroke-dasharray and stroke-dashoffset
  const radius = 45; // SVG viewBox is 100x100, so radius is set to 45
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  return (
    <div
      className={cn(
        "relative flex items-center justify-center",
        getSize(),
        className
      )}
    >
      {/* Background circle */}
      <svg viewBox="0 0 100 100" className="w-full h-full -rotate-90">
        <circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={getThickness()}
          className="text-slate-200 dark:text-slate-800"
        />

        {/* Foreground circle with animation */}
        <motion.circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={getThickness()}
          strokeDasharray={circumference}
          initial={{ strokeDashoffset: circumference }}
          animate={{ strokeDashoffset }}
          transition={{ duration: 1, ease: "easeOut" }}
          className={getColor()}
          strokeLinecap="round"
        />
      </svg>

      {/* Center text */}
      {showValue && (
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className={cn("font-bold", getFontSize())}>
            {value}
            {unit}
          </span>
          {label && (
            <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
              {label}
            </span>
          )}
        </div>
      )}
    </div>
  );
}
