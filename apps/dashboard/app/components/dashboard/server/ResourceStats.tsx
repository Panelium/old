import React from "react";
import { Cpu, Database, HardDrive, Server as ServerIcon } from "lucide-react";

import { Button } from "~/components/ui/button";
import type { Server } from "~/components/cards/server-card";
import OverviewCard from "~/components/cards/overview-card/Card";

const ResourceStats: React.FC<{ server: Server }> = ({ server }) => {
  const CPU_USAGE = {
    title: "25%",
    value: 30,
    max: 60,
    uiValue: "Core limit 100%",
  };

  const MEMORY_USAGE = {
    title: "1.5 GB",
    value: 1.5,
    max: 4,
    uiValue: "Used 1.5 GB of 4 GB",
  };

  const DISK_USAGE = {
    title: "3.2 GB",
    value: 3.2,
    max: 8,
    uiValue: "Used 3.2 GB of 8 GB",
  };

  const ConnectionMainComponent = () => {
    return (
      <div className="flex flex-1 flex-col space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-xs text-slate-500 dark:text-slate-400">
            IP Address
          </span>
          <span className="text-xs font-medium text-slate-900 dark:text-slate-50">
            {server.ip}
          </span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-xs text-slate-500 dark:text-slate-400">
            Port
          </span>
          <span className="text-xs font-medium text-slate-900 dark:text-slate-50">
            {server.port}
          </span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-xs text-slate-500 dark:text-slate-400">
            Location
          </span>
          <span className="text-xs font-medium text-slate-900 dark:text-slate-50">
            {server.location}
          </span>
        </div>
      </div>
    );
  };

  const ConnectionFooterComponent = () => {
    return (
      <Button
        variant="outline"
        size="sm"
        className="flex-1 w-full h-7 text-xs border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700"
      >
        Copy Details
      </Button>
    );
  };
  return (
    <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
      <OverviewCard bar={CPU_USAGE} title="CPU Usage" icon={Cpu} />
      <OverviewCard bar={MEMORY_USAGE} title="Memory Usage" icon={HardDrive} />
      <OverviewCard bar={DISK_USAGE} title="Disk Usage" icon={ServerIcon} />
      <OverviewCard
        title="Connection"
        icon={Database}
        children={<ConnectionMainComponent />}
        footerChildren={<ConnectionFooterComponent />}
      />
    </div>
  );
};

export default ResourceStats;
