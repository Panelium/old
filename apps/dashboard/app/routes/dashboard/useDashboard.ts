import { type Server } from "./ServerCardGrid";
import { ServerStatusType } from "proto-gen-ts/daemon_pb";

interface UseDashboard {
  mockServers: Server[];
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

const useDashboard = (): UseDashboard => {
  const mockServers: Server[] = [
    {
      id: "1",
      name: "Minecraft SMP",
      status: ServerStatusType.ONLINE,
      description: "Survival multiplayer server",
      cpuUsage: 45,
      memoryUsage: {
        used: 2048 * 1024,
        total: 4096 * 1024,
      },
      game: "Minecraft",
      players: {
        online: 12,
        max: 50,
      },
      ip: "192.168.1.1",
      port: 25565,
    },
    {
      id: "2",
      name: "Web Server",
      status: ServerStatusType.STARTING,
      description: "NGINX web server",
      cpuUsage: 12,
      memoryUsage: {
        used: 512 * 1024,
        total: 1024 * 1024,
      },
      game: "Nginx",
      ip: "192.168.1.6",
      port: 80,
    },
    {
      id: "3",
      name: "Terraria Adventure",
      status: ServerStatusType.OFFLINE,
      description: "Terraria modded adventure server",
      cpuUsage: 0,
      memoryUsage: {
        used: 0,
        total: 2048 * 1024,
      },
      game: "Terraria",
      players: {
        online: 0,
        max: 16,
      },
      ip: "192.168.1.3",
      port: 7777,
    },
    {
      id: "4",
      name: "Rust Survival",
      status: ServerStatusType.STOPPING,
      description: "Rust survival server",
      cpuUsage: 30,
      memoryUsage: {
        used: 1024 * 1024,
        total: 2048 * 1024,
      },
      game: "Rust",
      players: {
        online: 5,
        max: 100,
      },
      ip: "192.168.1.23",
      port: 28015,
    },
    {
      id: "5",
      name: "ARK: Survival Evolved",
      status: ServerStatusType.UNKNOWN,
      description: "ARK survival server",
      cpuUsage: 60,
      memoryUsage: {
        used: 3072 * 1024,
        total: 8192 * 1024,
      },
      game: "ARK",
      players: {
        online: 20,
        max: 70,
      },
      ip: "192.168.1.13",
      port: 7778,
    },
  ];

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
    mockServers,
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
