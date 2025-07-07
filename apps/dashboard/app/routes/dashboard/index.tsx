import React, { useEffect, useState } from "react";
import { Users } from "lucide-react";

import ServerCardGrid from "~/components/cards/server-card/ServerCardGrid";
import type { OverviewCardProps } from "~/components/cards/overview-card/OverviewCard";
import { getClientClient } from "~/lib/api-clients";
import { ServerInfo } from "proto-gen-ts/backend/Client_pb";

export default function DashboardOverviewPage() {
  const overviewCardsData: OverviewCardProps[] = [
    {
      title: "Player Activity",
      content: {
        title: "x",
        subtitle: `across y online servers`,
        icon: Users,
      },
      footer: `x/z total capacity`,
    },
    {
      title: "CPU Usage",
      gauge: {
        value: 0,
        maxValue: 100,
        size: "sm",
        unit: "%",
        label: "Average",
        subtitle: "Avg across 2 servers",
      },
      footer: `Peak: x%`,
    },
    {
      title: "Memory Usage",
      gauge: {
        value: 0.5,
        maxValue: 1,
        size: "sm",
        unit: "GB",
        label: "Used",
        subtitle: `of 0.0 GB allocated`,
      },
      footer: "",
    },
  ];

  const [serverInfos, setServerInfos] = useState<ServerInfo[]>([]);

  useEffect(() => {
    const fetchServerData = async () => {
      try {
        const client = await getClientClient();
        const servers = await client.getServerList({});
        setServerInfos(servers.servers);
        console.log(servers);
      } catch (error) {
        console.error("Failed to fetch server data:", error);
      }
    };

    fetchServerData();
  }, []);

  return (
    <div className="container py-6 px-4 mx-auto max-w-7xl">
      <div className="mb-6 no-select">
        <h1 className="text-3xl font-bold text-foreground">Dashboard Overview</h1>
        <p className="text-sm text-muted-foreground">Welcome back! Here's what's happening with your servers.</p>
      </div>
      {/*TODO: readd this*/}
      {/*<OverviewCardGrid cards={overviewCardsData} />*/}
      <ServerCardGrid serverInfos={serverInfos} />
    </div>
  );
}
