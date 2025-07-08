import React from "react";
import ServerHeader from "~/components/dashboard/server/ServerHeader";
import Pages from "~/components/dashboard/server/Pages";
import { useParams } from "react-router";
import { Client } from "@connectrpc/connect";
import { getClientClient } from "~/lib/api-clients";
import { ClientService, ServerInfo } from "proto-gen-ts/backend/Client_pb";

export default function ServerDetailsPage() {
  const { id } = useParams<{ id: string }>();

  const [clientClient, setClientClient] = React.useState<Client<typeof ClientService>>();
  const [serverInfo, setServerInfo] = React.useState<ServerInfo>();

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

  return (
    <div className="p-6 pb-16">
      <div className="max-w-7xl mx-auto space-y-4">
        {serverInfo ? (
          <>
            <ServerHeader server={serverInfo} />
            <Pages server={serverInfo} />
          </>
        ) : (
          <div className="text-center text-gray-500">Loading server details...</div>
        )}
      </div>
    </div>
  );
}
