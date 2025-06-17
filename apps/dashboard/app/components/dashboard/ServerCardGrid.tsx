import { ServerCard } from "~/components/dashboard/ServerCard";
import { Card } from "~/components/ui/card";
import { Plus } from "lucide-react";
import React from "react";
import { ServerStatusType } from "proto-gen-ts/daemon_pb";
import { cn } from "~/lib/utils";

export interface Server {
  id: string;
  name: string;
  status: ServerStatusType;
  description?: string;
  icon?: string;
  cpuUsage: number;
  memoryUsage: {
    used: number;
    total: number;
  };
  game?: string;
  players?: {
    online: number;
    max: number;
  };
  ip?: string;
  port?: number;
}

const AddServerCard: React.FC = () => {
  return (
    <Card
      className={cn(
        "flex flex-col items-center justify-center no-select",
        "border-slate-200 bg-transparent hover:bg-slate-50",
        "dark:border-slate-700 dark:hover:bg-slate-800/20",
        "transition-colors"
      )}
    >
      <div className="rounded-full bg-slate-100 dark:bg-slate-800 p-3 mb-3">
        <Plus className="h-6 w-6 text-slate-500 dark:text-slate-400" />
      </div>
      <h3 className="text-lg font-medium text-slate-900 dark:text-slate-100">
        New Server
      </h3>
      <p className="text-sm text-center text-slate-500 dark:text-slate-400 mt-1 mb-4">
        Deploy a new server with your preferred configuration
      </p>
    </Card>
  );
};
export default function ServerCardGrid({ servers }: { servers: Server[] }) {
  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-xl font-bold text-slate-900 dark:text-slate-50 no-select">
        Servers
      </h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
        {servers.map((server) => (
          <ServerCard key={server.id} server={server} className="no-select" />
        ))}
        <AddServerCard />
      </div>
    </div>
  );
}
