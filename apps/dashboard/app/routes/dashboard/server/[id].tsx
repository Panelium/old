import React from "react";
import ServerHeader from "~/components/dashboard/server/ServerHeader";
import Pages from "~/components/dashboard/server/Pages";

export default function ServerDetailsPage() {
  const { server } = {} as any; // TODO: replace with actual data fetching logic

  return (
    <div className="p-6 pb-16">
      <div className="max-w-7xl mx-auto space-y-4">
        <ServerHeader server={server} />
        <Pages />
      </div>
    </div>
  );
}
