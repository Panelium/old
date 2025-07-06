import React from "react";
import { FolderOpen, type LucideIcon, Play, Settings, Square, Terminal } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import { ServerStatusType } from "proto-gen-ts/daemon/Server_pb";
import { clampNumber, cn, formatMemory } from "~/lib/utils";

import { Card, CardContent, CardFooter, CardHeader } from "~/components/ui/card";
import ServerBar from "~/components/bars/ServerBar";
import EntityAvatar from "~/components/avatars/EntityAvatar";
import StatusBadge from "~/components/dashboard/StatusBadge";
import { ServerInfo } from "proto-gen-ts/backend/Client_pb";
import { ResourceUsage, ResourceUsageSchema } from "proto-gen-ts/common_pb";

export interface ServerData {
  serverInfo: ServerInfo;
  status: ServerStatusType;
  onlineSince: Date;
  resourceUsage: ResourceUsage;
}

interface ServerCardProps {
  serverInfo: ServerInfo;
  className?: string;
}

interface ServerCardLinkProps {
  icon: LucideIcon;
  color?: string;
  bgColor?: string;
  className?: string;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
}

interface ServerCardActionProps {
  serverData: ServerData;
  onAction: (e: React.MouseEvent, action: string) => void;
}

const serverCardTransition = "transition-all duration-300 ease-in-out truncate";

const ServerCardButton: React.FC<ServerCardLinkProps> = ({ icon, color, bgColor, className, onClick }) => {
  const IconComponent = icon;
  return (
    <div>
      <div
        className={cn("scale-20 group-hover:scale-100 w-10 h-10 bg-card rounded-full absolute", serverCardTransition)}
      />
      <button
        className={cn(
          "flex items-center justify-center h-10 w-10 rounded-full shadow-sm",
          "scale-20 group-hover:scale-100",
          "cursor-pointer",
          serverCardTransition,
          className
        )}
        style={{ color: color, backgroundColor: bgColor }}
        onClick={onClick}
      >
        <IconComponent className="w-5 h-5" />
      </button>
    </div>
  );
};

const ServerCardHeader: React.FC<{ serverData: ServerData }> = ({ serverData }) => {
  return (
    <CardHeader className="gap-3">
      <div className="flex items-start justify-between gap-3 truncate">
        <EntityAvatar
          src={serverData.serverInfo.softwareIcon}
          alt={serverData.serverInfo.name}
          title={serverData.serverInfo.name}
          subTitle={serverData.serverInfo.software}
          className={serverCardTransition}
        />
        <StatusBadge status={serverData.status} />
      </div>
      <p className="text-sm text-server-card-foreground line-clamp-2">{serverData.serverInfo.description}</p>
    </CardHeader>
  );
};

const ServerCardContent: React.FC<{ serverData: ServerData }> = ({ serverData }) => {
  const cpuUsage = serverData.resourceUsage.cpu;
  const cpuLimit = serverData.serverInfo.resourceLimit!.cpu;

  const cpuUsageText = `${cpuUsage.toFixed(1)}% / ${cpuLimit.toFixed(0)}%`;

  const cpuUsagePercentage = clampNumber((cpuUsage / cpuLimit) * 100, {
    min: 0,
    max: 100,
  });

  const ramUsage = serverData.resourceUsage.ram;
  const ramLimit = serverData.serverInfo.resourceLimit!.ram;

  const ramUsageText = `${formatMemory(ramUsage)} / ${formatMemory(ramLimit)}`;

  const storageUsage = serverData.resourceUsage.storage;
  const storageLimit = serverData.serverInfo.resourceLimit!.storage;

  const storageUsageText = `${formatMemory(storageUsage)} / ${formatMemory(storageLimit)}`;

  return (
    <CardContent className="space-y-3">
      <div className="flex flex-row gap-3">
        <div className="flex-1">
          <ServerBar title="CPU" uiValue={cpuUsageText} value={cpuUsagePercentage} max={100} />
        </div>
        <div className="flex-1">
          <ServerBar title="Memory" uiValue={ramUsageText} value={ramUsage} max={ramLimit} />
        </div>
      </div>
      <div>
        <ServerBar title="Storage" uiValue={storageUsageText} value={storageUsage} max={storageLimit} />
      </div>
    </CardContent>
  );
};

const ServerCardFooter: React.FC<{ serverData: ServerData }> = ({ serverData }) => {
  return (
    <CardFooter>
      <div className="flex flex-1 items-center justify-between text-xs text-card-muted-foreground">
        {serverData.serverInfo.mainAllocation && (
          <div className="font-mono">
            {serverData.serverInfo.mainAllocation.ip}
            {serverData.serverInfo.mainAllocation.port}
          </div>
        )}

        <div className="flex flex-row gap-1">
          {/*TODO: readd this*/}
          {/*{server.players && (*/}
          {/*  <>*/}
          {/*    <Users className="h-4 w-4 text-card-muted-foreground" />*/}
          {/*    <span>*/}
          {/*      {server.players.online} / {server.players.max} Players*/}
          {/*    </span>*/}
          {/*  </>*/}
          {/*)}*/}
        </div>
      </div>
    </CardFooter>
  );
};

const ServerCardAction: React.FC<ServerCardActionProps> = ({ serverData, onAction }) => {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center gap-3",
        "border-l border-tag-purple/30 bg-tag-purple-background/20",
        "absolute top-0 bottom-0 right-0 w-[60px] translate-x-[60px] ",
        "group-hover:w-[60px] group-hover:translate-x-[0px]",
        serverCardTransition
      )}
    >
      <ServerCardButton
        icon={serverData.status === ServerStatusType.ONLINE ? Square : Play}
        className={cn("bg-tag-green-background/30 text-tag-green hover:bg-tag-green-background/70")}
        onClick={(e) => onAction(e, "power")}
      />
      <ServerCardButton
        icon={Terminal}
        className={cn("bg-tag-purple-background/30 text-tag-purple hover:bg-tag-purple-background/70")}
        onClick={(e) => onAction(e, "console")}
      />
      <ServerCardButton
        icon={FolderOpen}
        className={cn("bg-tag-orange-background/30 text-tag-orange hover:bg-tag-orange-background/70")}
        onClick={(e) => onAction(e, "files")}
      />
      <ServerCardButton
        icon={Settings}
        className={cn("bg-tag-gray-background/30 text-tag-gray hover:bg-tag-gray-background/70")}
        onClick={(e) => onAction(e, "settings")}
      />
    </div>
  );
};

const ServerCard: React.FC<ServerCardProps> = ({ serverInfo, className }) => {
  const navigate = useNavigate();

  const handleServerCardButtonClick = (e: React.MouseEvent, action: string) => {
    e.preventDefault();
    e.stopPropagation();

    switch (action) {
      case "power":
        navigate(`/server/${serverInfo.sid}/power`);
        break;
      case "console":
        navigate(`/server/${serverInfo.sid}/console`);
        break;
      case "files":
        navigate(`/server/${serverInfo.sid}/files`);
        break;
      case "settings":
        navigate(`/server/${serverInfo.sid}/settings`);
        break;
      default:
        navigate(`/server/${serverInfo.sid}`);
    }
  };

  const serverData: ServerData = {
    serverInfo,
    status: ServerStatusType.UNKNOWN,
    onlineSince: new Date(),
    resourceUsage: {
      $typeName: ResourceUsageSchema.typeName,
      cpu: 0,
      ram: 0,
      storage: 0,
    },
  };

  // TODO: fetch actual data

  return (
    <Link to={`/server/${serverInfo.sid}`} className="group">
      <Card
        className={cn(
          "hover:rounded-h-xl relative flex flex-row overflow-hidden shadow-md shadow-black/20",
          "bg-server-card border-border hover:border-border-hover",
          "hover:scale-102 hover:shadow-xl hover:z-2 hover:shadow-black/80",
          serverCardTransition,
          className
        )}
      >
        <div className={cn("flex flex-col w-full gap-3", "group-hover:w-[calc(100%-60px)]", serverCardTransition)}>
          <ServerCardHeader serverData={serverData} />
          <ServerCardContent serverData={serverData} />
          <ServerCardFooter serverData={serverData} />
        </div>
        <ServerCardAction serverData={serverData} onAction={(e, action) => handleServerCardButtonClick(e, action)} />
      </Card>
    </Link>
  );
};

export default ServerCard;
