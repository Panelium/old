import React from "react";
import useServer from "./useServer";
import TabsSection from "~/components/dashboard/server/TabsSection";
import ServerHeader from "~/components/dashboard/server/ServerHeader";
import ResourceStats from "~/components/dashboard/server/ResourceStats";
import type { Server } from "~/components/cards/server-card/ServerCard";

export default function ServerDetailsPage() {
  const { server } = useServer();

  return (
    <div className="p-6 pb-16">
      <div className="max-w-7xl mx-auto space-y-8">
        <ServerHeader server={server as unknown as Server} />
        <ResourceStats server={server as unknown as Server} />
        <TabsSection />
      </div>
    </div>
  );
}
