import React from "react";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { type LucideIcon } from "lucide-react";
import { ResourceGauge, type ResourceGaugeProps } from "../ResourceGauge";
import OverviewBar from "~/components/bars/OverviewBar";
import { cn } from "~/lib/utils";

type GaugeProps = Omit<ResourceGaugeProps, "className">;

interface OverviewProps {
  title: string;
  subtitle?: string;
  icon?: LucideIcon;
}

interface BarProps {
  title: string;
  value: number;
  max: number;
  uiValue?: string;
}

interface OverviewCardBaseProps {
  title: string;
  icon?: LucideIcon;
  children?: React.ReactNode;
  content?: OverviewProps;
  gauge?: GaugeProps;
  bar?: BarProps;
}

type OverviewCardProps =
  | (OverviewCardBaseProps & { footer?: never; footerChildren?: never })
  | (OverviewCardBaseProps & { footer: string; footerChildren?: never })
  | (OverviewCardBaseProps & {
      footer?: never;
      footerChildren: React.ReactNode;
    });

const overviewCardTransition =
  "transition-all duration-300 ease-in-out truncate";

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
            <span className="text-xs card-foreground mt-1">
              {gauge.subtitle}
            </span>
          )}
        </div>
      </div>
    </>
  );
};

const ContentCard: React.FC<OverviewProps> = (content) => {
  return (
    <div className="flex flex-col">
      {content?.title && (
        <span className="text-3xl font-bold">{content.title}</span>
      )}
      {content?.subtitle && (
        <span className="text-xs card-foreground mt-1">{content.subtitle}</span>
      )}
    </div>
  );
};

const CardIcon: React.FC<OverviewProps> = (content) => {
  if (!content.icon) return null;
  return (
    <div className="h-12 w-12 rounded-full bg-tag-green-background/60 flex items-center justify-center ml-4">
      <content.icon className="h-6 w-6 text-tag-green" />
    </div>
  );
};

const CardBar: React.FC<BarProps> = (bar) => {
  return (
    <OverviewBar
      title={bar.title}
      value={bar.value}
      max={bar.max}
      uiValue={bar.uiValue}
    />
  );
};

const OverviewCard: React.FC<OverviewCardProps> = ({
  title,
  content,
  gauge,
  bar,
  footer,
  icon,
  children = null,
  footerChildren = null,
}) => {
  const IconComponent = icon ? icon : () => null;

  return (
    <Card
      className={cn(
        "shadow-sm hover:shadow-md transition-shadow border-border no-select",
        overviewCardTransition
      )}
    >
      <CardHeader className="flex flex-row pb-2 items-center gap-2">
        <IconComponent className="h-4 w-4 text-indigo-500" />
        <CardTitle className="text-sm font-medium text-card-foreground">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex h-full flex-row items-center justify-between flex-grow gap-2">
        {gauge && !content && <GaugeCard {...gauge} />}
        {!gauge && content && <ContentCard {...content} />}
        {!gauge && content?.icon && <CardIcon {...content} />}
        {bar && <CardBar {...bar} />}
        {!gauge && !content && !bar && children}
      </CardContent>
      {(footer || footerChildren) && (
        <CardFooter className="border-t border-border">
          <span className="flex-1 text-xs text-card-foreground h-3">
            {footer || footerChildren}
          </span>
        </CardFooter>
      )}
    </Card>
  );
};

export default OverviewCard;
export type { OverviewCardProps };
