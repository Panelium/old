import React, { useEffect } from "react";
import { Link } from "react-router";
import { ArrowLeft, FlagIcon, WifiIcon } from "lucide-react";

import { PowerAction, ServerService, ServerStatusType } from "proto-gen-ts/daemon/Server_pb";

import StatusBadge from "../StatusBadge";
import { Button } from "~/components/ui/button";
import SoftwareText from "~/components/texts/SoftwareText";
import DurationText from "~/components/texts/DurationText";
import PowerButton from "~/components/buttons/PowerButton";
import IconText from "~/components/texts/IconText";
import { ClientService, ServerInfo } from "proto-gen-ts/backend/Client_pb";
import { Client } from "@connectrpc/connect";
import { getClientClient, getDaemonServerClient } from "~/lib/api-clients";

const ServerHeader: React.FC<{ server: ServerInfo }> = ({ server }) => {
  const [clientClient, setClientClient] = React.useState<Client<typeof ClientService>>();
  const [serverClient, setServerClient] = React.useState<Client<typeof ServerService>>();
  const [status, setStatus] = React.useState<ServerStatusType>(ServerStatusType.UNKNOWN);
  const [onlineSince, setOnlineSince] = React.useState<Date | null>(null);

  useEffect(() => {
    (async () => {
      const clientClient = await getClientClient();
      setClientClient(clientClient);

      const serverClient = await getDaemonServerClient(server.daemonHost);
      setServerClient(serverClient);
    })();
  }, [server]);

  useEffect(() => {
    if (!serverClient) return;

    let interval: NodeJS.Timeout;

    const fetchStatus = async () => {
      const serverStatusResponse = await serverClient.status({ id: server.sid });
      setStatus(serverStatusResponse.status);
      serverStatusResponse.timestampStart &&
        serverStatusResponse.timestampStart.seconds > BigInt(0) &&
        setOnlineSince(new Date(Number(serverStatusResponse.timestampStart.seconds * BigInt(1000))));
    };

    fetchStatus();
    interval = setInterval(fetchStatus, 5000); // Poll every 5 seconds

    return () => clearInterval(interval);
  }, [serverClient, server.sid]);

  const handlePowerAction = async (e: PowerAction) => {
    if (!serverClient) return;
    await serverClient.powerAction({
      serverId: server.sid,
      action: e,
    });
  };

  return (
    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between no-select">
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          size="icon"
          className="rounded-lg shadow-sm bg-white dark:bg-slate-800 border-border"
          asChild
        >
          <Link to="/">
            <ArrowLeft className="h-4 w-4 text-slate-500 dark:text-slate-400" />
          </Link>
        </Button>
        <div className="flex gap-0 flex-col">
          <h1 className="text-2xl font-bold text-foreground">{server.name}</h1>
          <div className="flex items-center gap-3 mt-1">
            <StatusBadge status={status} />
            <SoftwareText software={server.software} />
            {onlineSince && <DurationText startDate={onlineSince} />}
            {server.mainAllocation && (
              <IconText text={`${server.mainAllocation.ip}:${server.mainAllocation.port}`} icon={WifiIcon} copy />
            )}
            <IconText text={server.location} icon={FlagIcon} />
          </div>
          <span className="mt-2 text-muted-foreground">{server.description}</span>
        </div>
      </div>
      <div className="flex gap-2">
        {status === ServerStatusType.ONLINE ? (
          <>
            <PowerButton
              action={PowerAction.RESTART}
              onClick={() => {
                handlePowerAction(PowerAction.RESTART);
              }}
            />
            <PowerButton
              action={PowerAction.STOP}
              onClick={() => {
                handlePowerAction(PowerAction.STOP);
              }}
            />
            <PowerButton
              action={PowerAction.KILL}
              onClick={() => {
                handlePowerAction(PowerAction.KILL);
              }}
            />
          </>
        ) : (
          <PowerButton
            action={PowerAction.START}
            onClick={() => {
              handlePowerAction(PowerAction.START);
            }}
          />
        )}
      </div>
    </div>
  );
};

export default ServerHeader;
