import React, { useState } from "react";
import { cn } from "~/lib/utils";
import { Link, useNavigate } from "react-router-dom";
import {
  FolderOpen,
  Play,
  Settings,
  Square,
  Terminal,
  Users,
} from "lucide-react";
import type { Server } from "~/components/dashboard/ServerCardGrid";
import StatusBadge from "~/components/dashboard/StatusBadge";
import { ServerStatusType } from "proto-gen-ts/daemon_pb";
import { Card, CardContent, CardFooter, CardHeader } from "../ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

interface BarProps {
  value: number;
  max: number;
  className?: string;
}

interface ResourceBarProps {
  server: Server;
  title: string;
  value: number;
  uiValue?: string;
  max: number;
  className?: string;
}

interface ServerCardProps {
  server: Server;
  className?: string;
}

const Bar: React.FC<BarProps> = ({ value, max, className }) => {
  const percentage = Math.min((value / max) * 100, 100);

  const getBarColor = () => {
    if (percentage > 90) return "bg-rose-500";
    if (percentage > 75) return "bg-amber-500";
    return "bg-emerald-500";
  };

  return (
    <div
      className={cn(
        "h-1.5 w-full rounded-full bg-slate-200 dark:bg-slate-800",
        className
      )}
    >
      <div
        className={cn("h-full rounded-full", getBarColor())}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
};

const ResourceBar: React.FC<ResourceBarProps> = (props) => {
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-xs">
        <span className="text-slate-600 dark:text-slate-400">
          {props.title}
        </span>
        <span className="font-medium text-slate-700 dark:text-slate-300">
          {props.uiValue}
        </span>
      </div>
      <Bar value={props.value} max={props.max} />
    </div>
  );
};

export function ServerCard({ server, className }: ServerCardProps) {
  const navigate = useNavigate();
  const [actionsOpen, setActionsOpen] = useState(false);

  // Format memory values
  const formatMemory = (bytes: number) => {
    const mb = bytes / 1024;
    if (mb >= 1024) {
      return `${(mb / 1024).toFixed(1)} GB`;
    }
    return `${Math.round(mb)} MB`;
  };

  // Handle specific action button clicks
  const handleActionClick = (e: React.MouseEvent, action: string) => {
    e.preventDefault(); // Prevent the link navigation
    e.stopPropagation(); // Prevent triggering the card's onClick

    switch (action) {
      case "power":
        navigate(`/server/${server.id}/power`);
        break;
      case "console":
        navigate(`/server/${server.id}/console`);
        break;
      case "files":
        navigate(`/server/${server.id}/files`);
        break;
      case "settings":
        navigate(`/server/${server.id}/settings`);
        break;
      default:
        navigate(`/server/${server.id}`);
    }
  };

  return (
    <>
      <Link to={`/server/${server.id}`}>
        <Card
          className={cn(
            "group relative overflow-hidden rounded-lg border bg-white shadow-sm border-slate-200 dark:border-slate-700",
            "hover:shadow-md hover:border-indigo-200 dark:hover:border-indigo-900/50",
            "dark:bg-slate-900 dark:border-slate-800",
            "flex h-full",
            className
          )}
        >
          <div className="flex flex-col gap-3">
            <CardHeader className="gap-3">
              <div className="flex items-start gap-3">
                <Avatar>
                  <AvatarImage
                    src={server.icon}
                    alt={server.name}
                    className="h-full w-full object-cover"
                  />
                  <AvatarFallback>
                    {server.name.charAt(0).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1 transition-colors truncate">
                  <h3 className="font-medium text-slate-900 dark:text-slate-100 group-hover:text-indigo-600 dark:group-hover:text-indigo-400">
                    {server.name}
                  </h3>
                  <p className="text-xs text-slate-500 dark:text-slate-400">
                    {server.game ? server.game : server.name || " "}
                  </p>
                </div>
                <StatusBadge status={server.status} />
              </div>
              <p className="text-sm text-slate-600 dark:text-slate-300 line-clamp-2">
                {server.description ? server.description : " "}
              </p>
            </CardHeader>
            <CardContent className="space-y-3">
              <ResourceBar
                server={server}
                title="CPU"
                uiValue={`${server.cpuUsage.toFixed(1)}%`}
                value={server.cpuUsage}
                max={100}
              />
              <ResourceBar
                server={server}
                title="Memory"
                uiValue={`${formatMemory(
                  server.memoryUsage.used
                )} / ${formatMemory(server.memoryUsage.total)}`}
                value={server.memoryUsage.used}
                max={server.memoryUsage.total}
              />
            </CardContent>
            <CardFooter>
              <div className="flex flex-1 items-center justify-between text-xs text-slate-600 dark:text-slate-400">
                <div className="font-mono">
                  {`${server.ip ? server.ip : "N/A"}:${
                    server.port ? server.port : "N/A"
                  }`}
                </div>

                <div className="flex flex-row gap-1">
                  {server.players && (
                    <>
                      <Users className="h-4 w-4 text-slate-500 dark:text-slate-400" />
                      <span>{`${server.players.online} / ${server.players.max} Players`}</span>
                    </>
                  )}
                </div>
              </div>
            </CardFooter>
          </div>
        </Card>
      </Link>
      <Link
        to={`/server/${server.id}`}
        className={cn(
          "group relative overflow-hidden rounded-lg border bg-white shadow-sm border-slate-200 dark:border-slate-700",
          "hover:shadow-md hover:border-indigo-200 dark:hover:border-indigo-900/50",
          "dark:bg-slate-900 dark:border-slate-800",
          "flex h-full",
          className
        )}
        onMouseEnter={() => setActionsOpen(true)}
        onMouseLeave={() => setActionsOpen(false)}
      >
        {/* Main content */}
        <div
          className={cn(
            "flex flex-col p-5 overflow-hidden transition-all duration-200 ease-in-out w-full",
            actionsOpen && "w-[calc(100%-60px)]"
          )}
        >
          {/* Status indicator in top right corner */}
          <div
            className={cn(
              "absolute top-3 transition-all duration-200 ease-in-out",
              actionsOpen ? "right-[calc(60px+0.75rem)]" : "right-3"
            )}
          >
            <StatusBadge status={server.status} />
          </div>

          {/* Server icon and name */}
          <div className="flex items-start gap-3 mb-3">
            {server.icon ? (
              <div className="h-10 w-10 rounded-md bg-slate-100 dark:bg-slate-800 overflow-hidden flex-shrink-0">
                <img
                  src={server.icon}
                  alt={server.name}
                  className="h-full w-full object-cover"
                />
              </div>
            ) : (
              <div className="flex h-10 w-10 items-center justify-center rounded-md bg-indigo-100 dark:bg-indigo-900/30 flex-shrink-0">
                <span className="text-lg font-bold text-indigo-700 dark:text-indigo-300">
                  {server.name.charAt(0).toUpperCase()}
                </span>
              </div>
            )}
            <div className="flex-1 min-w-0 pr-12">
              <h3 className="font-medium text-slate-900 dark:text-slate-100 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors truncate">
                {server.name}
              </h3>
              <div className="flex items-center justify-between">
                {server.game && (
                  <p className="text-xs text-slate-500 dark:text-slate-400">
                    {server.game}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* Server description if available */}
          {server.description && (
            <p className="text-sm text-slate-600 dark:text-slate-300 mb-4 line-clamp-2">
              {server.description}
            </p>
          )}

          {/* Resource usage */}
          <div className="space-y-3">
            <ResourceBar
              server={server}
              title="CPU"
              uiValue={`${server.cpuUsage.toFixed(1)}%`}
              value={server.cpuUsage}
              max={100}
            />
            <ResourceBar
              server={server}
              title="Memory"
              uiValue={`${formatMemory(
                server.memoryUsage.used
              )} / ${formatMemory(server.memoryUsage.total)}`}
              value={server.memoryUsage.used}
              max={server.memoryUsage.total}
            />
          </div>

          {/* Players and IP info */}
          <div className="mt-3 flex items-center justify-between text-xs text-slate-600 dark:text-slate-400">
            <div className="flex items-center gap-1">
              {server.ip && server.port && (
                <div className="ml-auto font-mono">
                  {server.ip}:{server.port}
                </div>
              )}
            </div>
            <div className="flex items-center gap-1">
              {server.players && (
                <>
                  <Users className="h-4 w-4 text-slate-500 dark:text-slate-400" />
                  <span>
                    {server.players.online} / {server.players.max} Players
                  </span>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Right-side action panel */}
        <div
          className={cn(
            "h-full flex flex-col items-center justify-center gap-3 overflow-hidden",
            "transition-all duration-200 ease-in-out",
            actionsOpen
              ? "w-[60px] border-l border-indigo-100 dark:border-indigo-800 bg-indigo-50 dark:bg-indigo-950/30 px-2"
              : "w-0 border-l-0 bg-transparent p-0"
          )}
          style={{
            minWidth: actionsOpen ? "60px" : "0",
            transition:
              "width 0.2s ease-in-out, min-width 0.2s ease-in-out, background-color 0.2s ease-in-out, border-color 0.2s ease-in-out, padding 0.2s ease-in-out",
          }}
          onClick={(e) => e.preventDefault()}
        >
          {/* Action buttons with visibility toggling */}
          <div
            className={cn(
              "flex flex-col items-center gap-3 w-full",
              actionsOpen ? "opacity-100" : "opacity-0 invisible",
              "transition-all duration-100"
            )}
          >
            {/* Start/Stop button */}
            <button
              className={cn(
                "h-10 w-10 rounded-full bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-400 dark:hover:bg-emerald-800/50 flex items-center justify-center shadow-sm",
                "transition-all duration-200",
                actionsOpen ? "scale-100" : "scale-75"
              )}
              style={{ transitionDelay: actionsOpen ? "25ms" : "0ms" }}
              title={
                server.status === ServerStatusType.ONLINE
                  ? "Stop Server"
                  : "Start Server"
              }
              onClick={(e) => handleActionClick(e, "power")}
            >
              {server.status === ServerStatusType.ONLINE ? (
                <Square className="w-5 h-5" />
              ) : (
                <Play className="w-5 h-5" />
              )}
            </button>

            {/* Console/Terminal button */}
            <button
              className={cn(
                "h-10 w-10 rounded-full bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400 dark:hover:bg-indigo-800/50 flex items-center justify-center shadow-sm",
                "transition-all duration-200",
                actionsOpen ? "scale-100" : "scale-75"
              )}
              style={{ transitionDelay: actionsOpen ? "25ms" : "0ms" }}
              title="Open Console"
              onClick={(e) => handleActionClick(e, "console")}
            >
              <Terminal className="w-5 h-5" />
            </button>

            {/* File Manager button */}
            <button
              className={cn(
                "h-10 w-10 rounded-full bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-400 dark:hover:bg-amber-800/50 flex items-center justify-center shadow-sm",
                "transition-all duration-200",
                actionsOpen ? "scale-100" : "scale-75"
              )}
              style={{ transitionDelay: actionsOpen ? "25ms" : "0ms" }}
              title="File Manager"
              onClick={(e) => handleActionClick(e, "files")}
            >
              <FolderOpen className="w-5 h-5" />
            </button>

            {/* Settings button */}
            <button
              className={cn(
                "h-10 w-10 rounded-full bg-slate-100 text-slate-700 hover:bg-slate-200 dark:bg-slate-800 dark:text-slate-400 dark:hover:bg-slate-700 flex items-center justify-center shadow-sm",
                "transition-all duration-200",
                actionsOpen ? "scale-100" : "scale-75"
              )}
              style={{ transitionDelay: actionsOpen ? "25ms" : "0ms" }}
              title="Server Settings"
              onClick={(e) => handleActionClick(e, "settings")}
            >
              <Settings className="w-5 h-5" />
            </button>
          </div>
        </div>
      </Link>
    </>
  );
}
