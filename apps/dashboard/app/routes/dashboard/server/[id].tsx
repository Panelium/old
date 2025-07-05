import React from "react";
import useServer from "./useServer";
import TabsSection from "~/components/dashboard/server/TabsSection";
import ServerHeader from "~/components/dashboard/server/ServerHeader";
import Pages from "~/components/dashboard/server/Pages";
import { Card } from "~/components/ui/card";

export default function ServerDetailsPage() {
  const { server } = useServer();

  return (
    <div className="p-6 pb-16">
      <div className="max-w-7xl mx-auto space-y-4">
        <ServerHeader server={server} />
          <Pages/>
      </div>
    </div>
  );
}
