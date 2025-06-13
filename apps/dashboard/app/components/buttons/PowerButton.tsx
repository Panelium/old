import IconButton from "~/components/buttons/IconButton";
import {PowerAction} from "proto-gen-ts/daemon_pb"
import {BadgeHelp, type LucideIcon, Play, Power, RefreshCw, Square} from "lucide-react";

const PowerActionIconMap: Record<PowerAction, LucideIcon> = {
    [PowerAction.UNSPECIFIED]: BadgeHelp,
    [PowerAction.START]: Play,
    [PowerAction.RESTART]: RefreshCw,
    [PowerAction.STOP]: Square,
    [PowerAction.KILL]: Power,
}

const PowerActionTextMap: Record<PowerAction, string> = {
    [PowerAction.UNSPECIFIED]: "Unknown",
    [PowerAction.START]: "Start",
    [PowerAction.RESTART]: "Restart",
    [PowerAction.STOP]: "Stop",
    [PowerAction.KILL]: "Kill",
}

export interface PowerButtonProps {
    action: PowerAction;
    onClick?: () => void;
}

const PowerActionButtonStyleMap: Record<PowerAction, string> = {
    [PowerAction.UNSPECIFIED]: "",
    [PowerAction.START]: "",
    [PowerAction.RESTART]: "",
    [PowerAction.STOP]: "text-rose-700 dark:text-rose-300 hover:bg-rose-100 dark:hover:bg-rose-900/20 hover:text-rose-700 border-rose-200 dark:border-rose-900/20",
    [PowerAction.KILL]: "text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/20 hover:text-red-700 border-red-200 dark:border-red-900/20",
}

export default function PowerButton({action, onClick}: PowerButtonProps) {
    return (
        <IconButton
            icon={PowerActionIconMap[action]}
            text={PowerActionTextMap[action]}
            onClick={onClick}
            className={PowerActionButtonStyleMap[action]}
        />
    );
}