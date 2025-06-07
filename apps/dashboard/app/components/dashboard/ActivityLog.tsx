import React from 'react';
import {Card} from '~/components/ui/card';
import {cn} from '~/lib/utils';
import {formatDistanceToNow} from 'date-fns';
import {Clock, Download, HardDrive, Server, Shield, Upload, UserCircle2} from 'lucide-react';

interface Activity {
    id: string;
    type: 'server_start' | 'server_stop' | 'backup_created' | 'file_upload' | 'file_download' | 'user_login' | 'user_action' | string;
    description: string;
    timestamp: string;
    server?: string;
    user?: string;
}

interface ActivityLogProps {
    activities: Activity[];
    className?: string;
    maxItems?: number;
}

export function ActivityLog({activities, className, maxItems = 10}: ActivityLogProps) {
    // Get icon based on activity type
    const getActivityIcon = (type: string) => {
        const iconMap: Record<string, React.ReactNode> = {
            server_start: <Server className="h-4 w-4 text-emerald-500"/>,
            server_stop: <Server className="h-4 w-4 text-amber-500"/>,
            backup_created: <HardDrive className="h-4 w-4 text-indigo-500"/>,
            file_upload: <Upload className="h-4 w-4 text-blue-500"/>,
            file_download: <Download className="h-4 w-4 text-blue-500"/>,
            user_login: <UserCircle2 className="h-4 w-4 text-slate-500"/>,
            user_action: <Shield className="h-4 w-4 text-purple-500"/>,
        };

        return iconMap[type] || <Clock className="h-4 w-4 text-slate-500"/>;
    };

    // Format relative time
    const formatRelativeTime = (timestamp: string) => {
        try {
            return formatDistanceToNow(new Date(timestamp), {addSuffix: true});
        } catch (e) {
            return 'unknown time';
        }
    };

    return (
        <Card className={cn("overflow-hidden py-0", className)}>
            <div className="divide-y divide-slate-200 dark:divide-slate-700">
                {activities.slice(0, maxItems).map((activity) => (
                    <div
                        key={activity.id}
                        className="flex items-start gap-4 p-4 transition-colors hover:bg-slate-50 dark:hover:bg-slate-800/50"
                    >
                        <div
                            className="flex h-8 w-8 items-center justify-center rounded-full bg-slate-100 dark:bg-slate-800">
                            {getActivityIcon(activity.type)}
                        </div>

                        <div className="flex-1 space-y-1">
                            <p className="text-sm font-medium text-slate-900 dark:text-slate-100">
                                {activity.description}
                            </p>

                            <div className="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs">
                                {activity.server && (
                                    <span
                                        className="inline-flex items-center rounded-full border border-slate-200 dark:border-slate-700 px-2 py-0.5 text-xs font-medium text-slate-800 dark:text-slate-300">
                    <Server className="mr-1 h-3 w-3 text-slate-500"/>
                                        {activity.server}
                  </span>
                                )}

                                {activity.user && (
                                    <span className="text-slate-500 dark:text-slate-400">
                    by {activity.user}
                  </span>
                                )}

                                <span className="text-slate-500 dark:text-slate-400">
                  {formatRelativeTime(activity.timestamp)}
                </span>
                            </div>
                        </div>
                    </div>
                ))}

                {activities.length === 0 && (
                    <div className="flex flex-col items-center justify-center py-8 px-4 text-center">
                        <div className="rounded-full bg-slate-100 dark:bg-slate-800 p-3 mb-3">
                            <Clock className="h-6 w-6 text-slate-400"/>
                        </div>
                        <h3 className="text-sm font-medium text-slate-900 dark:text-slate-100">No recent activity</h3>
                        <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                            Activity will appear here as you use your servers
                        </p>
                    </div>
                )}
            </div>
        </Card>
    );
}
