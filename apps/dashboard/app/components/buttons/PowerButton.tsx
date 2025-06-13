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
    danger?: boolean;
}

export default function PowerButton({action, onClick, danger}: PowerButtonProps) {
    return (
        <IconButton icon={PowerActionIconMap[action]} text={PowerActionTextMap[action]}/>
    );
}