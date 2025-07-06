import React from "react";
import { Users } from "lucide-react";
import useDashboard from "./useDashboard";

import ServerCardGrid from "~/components/cards/server-card/ServerCardGrid";
import type { OverviewCardProps } from "~/components/cards/overview-card/OverviewCard";

export default function DashboardOverviewPage() {
  const dashboardData = useDashboard();

  const overviewCardsData: OverviewCardProps[] = [
    {
      title: "Player Activity",
      content: {
        title: dashboardData.totalPlayers.toString(),
        subtitle: `across ${dashboardData.onlineServers} online servers`,
        icon: Users,
      },
      footer: `${dashboardData.totalPlayers}/${dashboardData.maxPlayers} total capacity`,
    },
    {
      title: "CPU Usage",
      gauge: {
        value: dashboardData.avgCpuUsage,
        maxValue: 100,
        size: "sm",
        unit: "%",
        label: "Average",
        subtitle: "Avg across 2 servers",
      },
      footer: `Peak: ${Math.max(
        ...dashboardData.onlineServersArray.map((s) => s.cpuUsage),
        0
      )}% on ${
        dashboardData.onlineServersArray.length > 0
          ? dashboardData.onlineServersArray.reduce((max, server) =>
              server.cpuUsage > max.cpuUsage ? server : max
            ).name
          : "N/A"
      }`,
    },
    {
      title: "Memory Usage",
      gauge: {
        value: parseFloat(dashboardData.usedMemoryGB),
        maxValue: parseFloat(dashboardData.totalMemoryGB),
        size: "sm",
        unit: "GB",
        label: "Used",
        subtitle: `of 0.0 GB allocated`,
      },
      footer: "",
    },
  ];

  return (
    <div className="container py-6 px-4 mx-auto max-w-7xl">
      <div className="mb-6 no-select">
        <h1 className="text-3xl font-bold text-foreground">
          Dashboard Overview
        </h1>
        <p className="text-sm text-muted-foreground">
          Welcome back! Here's what's happening with your servers.
        </p>
      </div>
      {/*TODO: readd this*/}
      {/*<OverviewCardGrid cards={overviewCardsData} />*/}
      <ServerCardGrid servers={dashboardData.mockServers} />
    </div>
  );
}
