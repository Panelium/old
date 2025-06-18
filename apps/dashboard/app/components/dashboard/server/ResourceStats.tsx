import React from "react";
import { Cpu, Database, HardDrive } from "lucide-react";
import { cn, formatDisk, formatMemory } from "~/lib/utils";

import { Button } from "~/components/ui/button";
import type { Server } from "~/components/cards/server-card";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";

const ResourceStats: React.FC<{ server: Server }> = ({ server }) => {
  return (
    <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
      <Card className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
            <Cpu className="mr-2 h-4 w-4 text-indigo-500" />
            CPU Usage
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold text-slate-900 dark:text-slate-50">
              {server.cpuUsage}%
            </div>
            <div className="text-xs text-slate-500 dark:text-slate-400">
              Core Limit: 100%
            </div>
          </div>
          <div className="mt-3 h-1.5 w-full bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
            <div
              className={cn(
                "h-full rounded-full",
                server.cpuUsage > 90
                  ? "bg-red-500"
                  : server.cpuUsage > 75
                  ? "bg-amber-500"
                  : "bg-emerald-500"
              )}
              style={{ width: `${server.cpuUsage}%` }}
            />
          </div>
        </CardContent>
      </Card>

      {/* Memory Usage Card */}
      <Card className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
            <HardDrive className="mr-2 h-4 w-4 text-indigo-500" />
            Memory Usage
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold text-slate-900 dark:text-slate-50">
              {formatMemory(server.memoryUsage.used)}
            </div>
            <div className="text-xs text-slate-500 dark:text-slate-400">
              of {formatMemory(server.memoryUsage.total)}
            </div>
          </div>
          <div className="mt-3 h-1.5 w-full bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
            <div
              className={cn(
                "h-full rounded-full",
                server.memoryUsage.used / server.memoryUsage.total > 0.9
                  ? "bg-red-500"
                  : server.memoryUsage.used / server.memoryUsage.total > 0.75
                  ? "bg-amber-500"
                  : "bg-emerald-500"
              )}
              style={{
                width: `${
                  (server.memoryUsage.used / server.memoryUsage.total) * 100
                }%`,
              }}
            />
          </div>
        </CardContent>
      </Card>

      {/* Disk Usage Card */}
      <Card className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
            <HardDrive className="mr-2 h-4 w-4 text-indigo-500" />
            Disk Usage
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold text-slate-900 dark:text-slate-50">
              {formatDisk(server.diskUsage.used)}
            </div>
            <div className="text-xs text-slate-500 dark:text-slate-400">
              of {formatDisk(server.diskUsage.total)}
            </div>
          </div>
          <div className="mt-3 h-1.5 w-full bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
            <div
              className={cn(
                "h-full rounded-full",
                server.diskUsage.used / server.diskUsage.total > 0.9
                  ? "bg-red-500"
                  : server.diskUsage.used / server.diskUsage.total > 0.75
                  ? "bg-amber-500"
                  : "bg-emerald-500"
              )}
              style={{
                width: `${
                  (server.diskUsage.used / server.diskUsage.total) * 100
                }%`,
              }}
            />
          </div>
        </CardContent>
      </Card>

      {/* Connection Card */}
      <Card className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
            <Database className="mr-2 h-4 w-4 text-indigo-500" />
            Connection
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col space-y-2">
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
            <div className="mt-2 pt-2 border-t border-slate-200 dark:border-slate-700">
              <Button
                variant="outline"
                size="sm"
                className="w-full h-7 text-xs border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700"
              >
                Copy Details
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default ResourceStats;
