import React from "react";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { ResourceGauge, type ResourceGaugeProps } from "./ResourceGauge";
import { type LucideIcon } from "lucide-react";

interface OverviewCardContent {
  title: string;
  subtitle?: string;
  icon?: LucideIcon;
}

interface BaseOverviewCard {
  title: string;
  footer?: string;
}

interface OverviewCardWithContent extends BaseOverviewCard {
  content: OverviewCardContent;
  gauge?: never;
}

interface OverviewCardWithGauge extends BaseOverviewCard {
  content?: never;
  gauge: Omit<ResourceGaugeProps, "className">;
}

export type OverviewCardProps = OverviewCardWithContent | OverviewCardWithGauge;

const GaugeCard: React.FC<ResourceGaugeProps> = (gauge) => {
  return (
    <>
      <div className="flex items-center justify-between w-full">
        <ResourceGauge
          value={gauge.value}
          maxValue={gauge.maxValue}
          size={gauge.size}
          unit={gauge.unit}
          label={gauge.label}
        />
        <div className="flex flex-col items-end ml-4">
          <span className="text-3xl font-bold">
            {gauge.unit != "%" ? gauge.value.toFixed(1) + " " : gauge.value}
            {gauge.unit}
          </span>
          {gauge.subtitle && (
            <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
              {gauge.subtitle}
            </span>
          )}
        </div>
      </div>
    </>
  );
};

const ContentCard: React.FC<OverviewCardContent> = (content) => {
  return (
    <div className="flex flex-col">
      {content?.title && (
        <span className="text-3xl font-bold">{content.title}</span>
      )}
      {content?.subtitle && (
        <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
          {content.subtitle}
        </span>
      )}
    </div>
  );
};

const CardIcon: React.FC<OverviewCardContent> = (content) => {
  if (!content.icon) return null;
  return (
    <div className="h-12 w-12 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center ml-4">
      <content.icon className="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
    </div>
  );
};

const OverviewCard: React.FC<OverviewCardProps> = ({
  title,
  content,
  gauge,
  footer,
}) => {
  return (
    <Card className="shadow-sm hover:shadow-md transition-shadow border-slate-200 dark:border-slate-700 no-select">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex h-full flex-row items-center justify-between flex-grow h-full gap-2">
        {gauge && !content && <GaugeCard {...gauge} />}
        {!gauge && content && <ContentCard {...content} />}
        {!gauge && content?.icon && <CardIcon {...content} />}
      </CardContent>
      <CardFooter className="border-t border-slate-200 dark:border-slate-700">
        <span className="text-xs text-slate-500 dark:text-slate-400 h-3">
          {footer ? footer : ""}
        </span>
      </CardFooter>
    </Card>
  );
};

export default OverviewCard;
