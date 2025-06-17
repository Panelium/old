import type { ServerStatusType } from "proto-gen-ts/daemon_pb";

import ServerCardGrid from "./grid";

interface Server {
  id: string;
  name: string;
  status: ServerStatusType;
  description?: string;
  icon?: string;
  cpuUsage: number;
  memoryUsage: {
    used: number;
    total: number;
  };
  game?: string;
  players?: {
    online: number;
    max: number;
  };
  ip?: string;
  port?: number;
}

export default ServerCardGrid;
export type { Server };
