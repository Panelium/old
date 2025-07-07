import React from "react";
import ServerHeader from "~/components/dashboard/server/ServerHeader";
import Pages from "~/components/dashboard/server/Pages";
import { useParams } from "react-router";
import { ServerData } from "~/components/cards/server-card/ServerCard";
import { Client } from "@connectrpc/connect";
import { getClientClient, getDaemonServerClient } from "~/lib/api-clients";
import { ServerService } from "proto-gen-ts/daemon/Server_pb";
import { ClientService, ServerInfo } from "proto-gen-ts/backend/Client_pb";

export default function ServerDetailsPage() {
  const { id } = useParams<{ id: string }>();

  const [clientClient, setClientClient] = React.useState<Client<typeof ClientService>>();
  const [serverInfo, setServerInfo] = React.useState<ServerInfo>();
  const [serverClient, setServerClient] = React.useState<Client<typeof ServerService>>();

  React.useEffect(() => {
    (async () => {
      const clientClient = await getClientClient();
      setClientClient(clientClient);
    })();
  }, []);

  React.useEffect(() => {
    (async () => {
      if (!clientClient || !id) return;

      const serverInfoResponse = await clientClient.getServer({ id });
      if (serverInfoResponse) {
        setServerInfo(serverInfoResponse);
      } else {
        console.error("Failed to fetch server data");
      }
    })();
  }, [clientClient, id]);

  const [serverData, setServerData] = React.useState<ServerData>();

  React.useEffect(() => {
    (async () => {
      if (!serverInfo) return;
      const client = await getDaemonServerClient(serverInfo?.daemonHost);
      setServerClient(client);
    })();
  }, [serverInfo]);

  React.useEffect(() => {
    (async () => {
      if (!serverClient || !id || !serverInfo) return;

      const serverResourceUsageResponse = await serverClient.resourceUsage({ id });
      const serverStatusResponse = await serverClient.status({ id });
    })();
  }, [serverClient, id, serverInfo]);

  return (
    <div className="p-6 pb-16">
      <div className="max-w-7xl mx-auto space-y-4">
        {serverData ? (
          <>
            <ServerHeader server={serverData} />
            <Pages server={serverData} />
          </>
        ) : (
          <div className="text-center text-gray-500">Loading server details...</div>
        )}
      </div>
    </div>
  );
}
