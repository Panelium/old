import React from "react";
import { Link } from "react-router";
import { ArrowLeft, FlagIcon, WifiIcon } from "lucide-react";

import { PowerAction } from "proto-gen-ts/daemon/Server_pb";

import StatusBadge from "../StatusBadge";
import { Button } from "~/components/ui/button";
import SoftwareText from "~/components/texts/SoftwareText";
import DurationText from "~/components/texts/DurationText";
import PowerButton from "~/components/buttons/PowerButton";
import type { ServerData } from "~/components/cards/server-card/ServerCard";
import IconText from "~/components/texts/IconText";

const ServerHeader: React.FC<{ server: ServerData }> = ({ server }) => {
  return (
    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between no-select">
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          size="icon"
          className="rounded-lg shadow-sm bg-white dark:bg-slate-800 border-border"
          asChild
        >
          <Link to="/">
            <ArrowLeft className="h-4 w-4 text-slate-500 dark:text-slate-400" />
          </Link>
        </Button>
        <div className="flex gap-0 flex-col">
          <h1 className="text-2xl font-bold text-foreground">{server.serverInfo.name}</h1>
          <div className="flex items-center gap-3 mt-1">
            <StatusBadge status={server.status} />
            <SoftwareText software={server.serverInfo.software} />
            <DurationText startDate={server.onlineSince} />
            {server.serverInfo.mainAllocation && (
              <IconText
                text={`${server.serverInfo.mainAllocation.ip}:${server.serverInfo.mainAllocation.port}`}
                icon={WifiIcon}
              />
            )}
            <IconText text={server.serverInfo.location} icon={FlagIcon} />
          </div>
          <span className="mt-2 text-muted-foreground">{server.serverInfo.description}</span>
        </div>
      </div>
      <div className="flex gap-2">
        <PowerButton action={PowerAction.RESTART} />
        <PowerButton action={PowerAction.STOP} />
        <PowerButton action={PowerAction.KILL} />
      </div>
    </div>
  );
};

export default ServerHeader;
