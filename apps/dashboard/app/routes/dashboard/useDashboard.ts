import { type Server } from "./ServerCardGrid";
import { ServerStatusType } from "proto-gen-ts/daemon_pb";

interface UseDashboard {
  mockServers: Server[];
}

interface UseDashboardReturn {
  totalServers: number;
  onlineServers: number;
  totalPlayers: number;
  maxPlayers: number;
  onlineServersArray: Server[];
  avgCpuUsage: number;
  totalMemoryUsed: number;
  totalMemoryAllocated: number;
  usedMemoryGB: string;
  totalMemoryGB: string;
}

const useDashboard = ({ mockServers }: UseDashboard): UseDashboardReturn => {
  const totalServers = mockServers.length;
  const onlineServers = mockServers.filter(
    (s) => s.status === ServerStatusType.ONLINE
  ).length;

  const totalPlayers = mockServers.reduce((acc, server) => {
    return acc + (server.players?.online || 0);
  }, 0);
  const maxPlayers = mockServers.reduce((acc, server) => {
    return acc + (server.players?.max || 0);
  }, 0);

  const onlineServersArray = mockServers.filter(
    (s) => s.status === ServerStatusType.ONLINE
  );
  const avgCpuUsage = onlineServersArray.length
    ? Math.round(
        onlineServersArray.reduce((acc, server) => acc + server.cpuUsage, 0) /
          onlineServersArray.length
      )
    : 0;

  const totalMemoryUsed = mockServers.reduce(
    (acc, server) => acc + server.memoryUsage.used,
    0
  );
  const totalMemoryAllocated = mockServers.reduce(
    (acc, server) => acc + server.memoryUsage.total,
    0
  );

  const usedMemoryGB = (totalMemoryUsed / (1024 * 1024 * 1024)).toFixed(1);
  const totalMemoryGB = (totalMemoryAllocated / (1024 * 1024 * 1024)).toFixed(
    1
  );

  return {
    totalServers,
    onlineServers,
    totalPlayers,
    maxPlayers,
    onlineServersArray,
    avgCpuUsage,
    totalMemoryUsed,
    totalMemoryAllocated,
    usedMemoryGB,
    totalMemoryGB,
  };
};

export default useDashboard;
