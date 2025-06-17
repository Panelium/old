import React from "react";
import { cn } from "~/lib/utils";
import type { Server } from "../cards/server-card";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

const EntityAvatar: React.FC<{ server: Server; className: string }> = ({
  server,
  className,
}) => {
  return (
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarImage
          src={server.icon}
          alt={server.name}
          className="h-full w-full object-cover"
        />
        <AvatarFallback>{server.name.charAt(0).toUpperCase()}</AvatarFallback>
      </Avatar>
      <div className="flex flex-col flex-1">
        <h3
          className={cn(
            "font-medium",
            "text-slate-900 dark:text-slate-100",
            "group-hover:text-indigo-600 dark:group-hover:text-indigo-400",
            className
          )}
        >
          {server.name}
        </h3>
        <p className="text-xs text-slate-500 dark:text-slate-400">
          {server.game ? server.game : server.name || " "}
        </p>
      </div>
    </div>
  );
};

export default EntityAvatar;
