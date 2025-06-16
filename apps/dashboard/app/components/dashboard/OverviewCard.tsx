import {ResourceGauge, type ResourceGaugeProps} from "~/components/dashboard/ResourceGauge";
import type {LucideIcon} from "lucide-react";
import {Card, CardContent, CardHeader, CardTitle} from "~/components/ui/card";

export interface BasicOverviewCardContent {
    title: string;
    subtitle?: string;
    icon?: LucideIcon;
}

interface OverviewCardWithContent {
    title: string;
    content: BasicOverviewCardContent;
    gauge?: never;
    footer?: string;
}

interface OverviewCardWithGauge {
    title: string;
    content?: never;
    gauge: Omit<ResourceGaugeProps, 'className'> & {
        value: number;
        maxValue: number;
        size?: 'sm' | 'md' | 'lg';
        unit?: string;
        label?: string;
    } & {
        subtitle?: string;
    }
    footer?: string;
}

export type OverviewCardProps = OverviewCardWithContent | OverviewCardWithGauge;

export default function OverviewCard({title, content, gauge, footer}: OverviewCardProps) {
    return (
        <Card
            className="shadow-sm hover:shadow-md transition-shadow border-slate-200 dark:border-slate-700 no-select">
            <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400">
                    {title}
                </CardTitle>
            </CardHeader>
            <CardContent className="h-full">
                <div className="flex flex-col justify-between flex-grow h-full">
                    <div className="flex items-center justify-between">
                        {gauge ? (
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
                                    <span
                                        className="text-3xl font-bold">{gauge.unit != '%' ? (gauge.value.toFixed(1) + " ") : gauge.value}{gauge.unit}</span>
                                        {gauge.subtitle && (
                                            <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                                            {gauge.subtitle}
                                        </span>
                                        )}
                                    </div>
                                </div>
                            </>
                        ) : content && (
                            <div className="flex flex-col">
                                {content?.title && (
                                    <span className="text-3xl font-bold">
                                    {content.title}
                                </span>
                                )}
                                {content?.subtitle && (
                                    <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                                    {content.subtitle}
                                </span>
                                )}
                            </div>
                        )}
                        {content?.icon && (
                            <div
                                className="h-12 w-12 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center ml-4">
                                <content.icon className="h-6 w-6 text-emerald-600 dark:text-emerald-400"/>
                            </div>
                        )}
                    </div>
                    {footer && (
                        <div className="mt-2 pt-2 border-t border-slate-200 dark:border-slate-700">
                            <span className="text-xs text-slate-500 dark:text-slate-400">
                                {footer}
                            </span>
                        </div>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}
