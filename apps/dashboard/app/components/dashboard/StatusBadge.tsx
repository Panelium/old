import {ServerStatusType} from "proto-gen-ts/daemon_pb"
import {cn} from "~/lib/utils";
import {Badge} from "~/components/ui/badge";

const PULSE_STYLE = "animate-pulse";

const BASE_BADGE_STYLE = "inline-flex items-center gap-1.5 px-2 py-0.5 text-xs font-medium rounded-full";
const BadgeStyleMap: Record<ServerStatusType, string> = {
    [ServerStatusType.UNKNOWN]: cn(BASE_BADGE_STYLE,  `bg-tag-gray-background/30   text-tag-gray   border-tag-gray-border/80`),
    [ServerStatusType.ONLINE]: cn(BASE_BADGE_STYLE,   `bg-tag-green-background/30  text-tag-green  border-tag-green-border/50`),
    [ServerStatusType.STARTING]: cn(BASE_BADGE_STYLE, `bg-tag-yellow-background/30 text-tag-yellow border-tag-yellow-border/50`),
    [ServerStatusType.OFFLINE]: cn(BASE_BADGE_STYLE,  `bg-tag-gray-background/30   text-tag-gray   border-tag-gray-border/80`),
    [ServerStatusType.STOPPING]: cn(BASE_BADGE_STYLE, `bg-tag-red-background/30    text-tag-red    border-tag-red-border/50`),
}

const BASE_BADGE_DOT_STYLE = "mr-1 h-1.5 w-1.5 rounded-full inline-block";
const BadgeDotStyleMap: Record<ServerStatusType, string> = {
    [ServerStatusType.UNKNOWN]: cn(BASE_BADGE_DOT_STYLE,  `bg-tag-gray`),
    [ServerStatusType.ONLINE]: cn(BASE_BADGE_DOT_STYLE,   `bg-tag-green`, PULSE_STYLE),
    [ServerStatusType.STARTING]: cn(BASE_BADGE_DOT_STYLE, `bg-tag-yellow`, PULSE_STYLE),
    [ServerStatusType.OFFLINE]: cn(BASE_BADGE_DOT_STYLE,  `bg-tag-gray`),
    [ServerStatusType.STOPPING]: cn(BASE_BADGE_DOT_STYLE, `bg-tag-red`, PULSE_STYLE),
}

const ServerStatusTypeNames: Record<ServerStatusType, string> = {
    [ServerStatusType.UNKNOWN]: "Unknown",
    [ServerStatusType.ONLINE]: "Online",
    [ServerStatusType.OFFLINE]: "Offline",
    [ServerStatusType.STARTING]: "Starting",
    [ServerStatusType.STOPPING]: "Stopping",
}

export default function StatusBadge({status, customStatusName}: {
    status: ServerStatusType,
    customStatusName?: string
}) {
    const badgeStyle = BadgeStyleMap[status] ?? BadgeStyleMap[ServerStatusType.UNKNOWN];
    const badgeDotStyle = BadgeDotStyleMap[status] ?? BadgeDotStyleMap[ServerStatusType.UNKNOWN];
    const statusName = customStatusName ?? ServerStatusTypeNames[status] ?? "Unknown";

    return (
        <Badge variant="outline" className={badgeStyle}>
            {(status != ServerStatusType.UNKNOWN) && (<span className={badgeDotStyle}></span>)}
            {statusName}
        </Badge>
    );

}