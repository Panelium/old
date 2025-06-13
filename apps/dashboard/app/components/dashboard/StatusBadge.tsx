import {ServerStatusType} from "proto-gen-ts/daemon_pb"
import {cn} from "~/lib/utils";
import {Badge} from "~/components/ui/badge";

const PULSE_STYLE = "animate-pulse";

const BASE_BADGE_STYLE = "inline-flex items-center gap-1.5 px-2 py-0.5 text-xs font-medium rounded-full";
const BadgeStyleMap: Record<ServerStatusType, string> = {
    [ServerStatusType.SERVER_STATUS_UNKNOWN]: cn(BASE_BADGE_STYLE, `bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400 border-gray-300 dark:border-gray-800/80`),
    [ServerStatusType.SERVER_STATUS_ONLINE]: cn(BASE_BADGE_STYLE, `bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/50`),
    [ServerStatusType.SERVER_STATUS_STARTING]: cn(BASE_BADGE_STYLE, `bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400 border-amber-200 dark:border-amber-900/50`),
    [ServerStatusType.SERVER_STATUS_OFFLINE]: cn(BASE_BADGE_STYLE, `bg-slate-100 text-slate-700 dark:bg-slate-900/30 dark:text-slate-400 border-slate-300 dark:border-slate-800/80`),
    [ServerStatusType.SERVER_STATUS_STOPPING]: cn(BASE_BADGE_STYLE, `bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400 border-rose-200 dark:border-rose-900/50`),
}

const BASE_BADGE_DOT_STYLE = "mr-1 h-1.5 w-1.5 rounded-full inline-block";
const BadgeDotStyleMap: Record<ServerStatusType, string> = {
    [ServerStatusType.SERVER_STATUS_UNKNOWN]: cn(BASE_BADGE_DOT_STYLE, `bg-gray-500 dark:bg-gray-400`),
    [ServerStatusType.SERVER_STATUS_ONLINE]: cn(BASE_BADGE_DOT_STYLE, `bg-emerald-500 dark:bg-emerald-400`, PULSE_STYLE),
    [ServerStatusType.SERVER_STATUS_STARTING]: cn(BASE_BADGE_DOT_STYLE, `bg-amber-500 dark:bg-amber-400`, PULSE_STYLE),
    [ServerStatusType.SERVER_STATUS_OFFLINE]: cn(BASE_BADGE_DOT_STYLE, `bg-slate-500 dark:bg-slate-400`),
    [ServerStatusType.SERVER_STATUS_STOPPING]: cn(BASE_BADGE_DOT_STYLE, `bg-rose-500 dark:bg-rose-400`, PULSE_STYLE),
}

const ServerStatusTypeNames: Record<ServerStatusType, string> = {
    [ServerStatusType.SERVER_STATUS_UNKNOWN]: "Unknown",
    [ServerStatusType.SERVER_STATUS_ONLINE]: "Online",
    [ServerStatusType.SERVER_STATUS_OFFLINE]: "Offline",
    [ServerStatusType.SERVER_STATUS_STARTING]: "Starting",
    [ServerStatusType.SERVER_STATUS_STOPPING]: "Stopping",
}

export default function StatusBadge({status, customStatusName}: {
    status: ServerStatusType,
    customStatusName?: string
}) {
    const badgeStyle = BadgeStyleMap[status] ?? BadgeStyleMap[ServerStatusType.SERVER_STATUS_UNKNOWN];
    const badgeDotStyle = BadgeDotStyleMap[status] ?? BadgeDotStyleMap[ServerStatusType.SERVER_STATUS_UNKNOWN];
    const statusName = customStatusName ?? ServerStatusTypeNames[status] ?? "Unknown";

    return (
        <Badge variant="outline" className={badgeStyle}>
            {(status != ServerStatusType.SERVER_STATUS_UNKNOWN) && (<span className={badgeDotStyle}></span>)}
            {statusName}
        </Badge>
    );

}