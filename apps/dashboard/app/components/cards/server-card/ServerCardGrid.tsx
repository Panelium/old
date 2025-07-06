import React from "react";
import { Plus } from "lucide-react";
import { cn } from "~/lib/utils";

import { Card } from "~/components/ui/card";
import ServerCard from "~/components/cards/server-card/ServerCard";
import { ServerInfo } from "proto-gen-ts/backend/Client_pb";

const AddServerCard: React.FC = () => {
  return (
    <Card
      className={cn(
        "hover:rounded-h-xl flex flex-col items-center justify-center no-select flex-1/2",
        "border-border bg-server-card/50 hover:border-border-hover/50 hover:bg-server-card",
        "transition-all cursor-pointer backdrop-blur-[3px]"
      )}
    >
      <div className="rounded-full bg-tag-gray-background p-3">
        <Plus className="h-6 w-6 text-tag-gray" />
      </div>
      <h3 className="text-lg font-medium text-foreground">New Server</h3>
      <p className="text-sm text-center text-server-card-foreground mt-1 mb-2 min-h-10">
        Deploy a new server with your preferred configuration
      </p>
    </Card>
  );
};
export default function ServerCardGrid({ serverInfos }: { serverInfos: ServerInfo[] }) {
  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-xl font-bold text-foreground no-select">Servers</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
        <AddServerCard />
        {serverInfos.map((serverInfo) => (
          <ServerCard key={serverInfo.sid} serverInfo={serverInfo} className="no-select" />
        ))}
      </div>
    </div>
  );
}
