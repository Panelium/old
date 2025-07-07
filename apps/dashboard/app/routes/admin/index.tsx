import { AnimatePresence, motion } from "framer-motion";
import { ChevronDown, ChevronUp } from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "~/components/ui/button";
import { Card } from "~/components/ui/card";
import { cn } from "~/lib/utils";
import {
  getAdminLocationManagerClient,
  getAdminNodeManagerClient,
  getAdminServerManagerClient,
  getAdminUserManagerClient,
} from "~/lib/api-clients";
import { User } from "proto-gen-ts/backend/admin/UserManager_pb";
import { Location } from "proto-gen-ts/backend/admin/LocationManager_pb";
import { Node } from "proto-gen-ts/backend/admin/NodeManager_pb";
import { Server } from "proto-gen-ts/backend/admin/ServerManager_pb";
import { Pagination } from "proto-gen-ts/common_pb";

interface Column<T> {
  label: string;
  id: keyof T;
  type?: "string" | "number" | "boolean" | "percent" | "memory";
  sortable?: boolean;
  sortFunction?: (a: T, b: T, ascending: boolean) => number;
  fetchFunction?: (a: T) => string | number | boolean;
}

const USERS_COLUMNS: Column<User>[] = [
  { label: "User ID", id: "uid", sortable: false },
  { label: "Username", id: "username" },
  { label: "Email", id: "email" },
  { label: "Admin", id: "admin", type: "boolean" },
  { label: "MFA Needed", id: "mfaNeeded", type: "boolean" },
];

const LOCATIONS_COLUMNS: Column<Location>[] = [
  { label: "Location ID", id: "lid", sortable: false },
  { label: "Name", id: "name" },
];

const NODES_COLUMNS: Column<Node>[] = [
  { label: "Node ID", id: "nid", sortable: false },
  { label: "Name", id: "name" },
  { label: "Location ID", id: "lid" },
  { label: "FQDN", id: "fqdn" },
  { label: "Daemon Port", id: "daemonPort" },
  { label: "HTTPS", id: "https" },
  { label: "Max CPU", id: "maxCpu", type: "number" },
  { label: "Max RAM", id: "maxRam", type: "number" },
  { label: "Max Storage", id: "maxStorage", type: "number" },
  { label: "Max Swap", id: "maxSwap", type: "number" },
];

const SERVERS_COLUMNS: Column<Server>[] = [
  { label: "Server ID", id: "sid", sortable: false },
  { label: "Name", id: "name" },
  { label: "Owner UID", id: "ownerUid" },
  { label: "Node ID", id: "nid" },
  { label: "Blueprint ID", id: "bid" },
  { label: "Docker Image", id: "dockerImage" },
];

function TableHead<T>({
  columns,
  handleSort,
}: {
  columns: Column<T>[];
  handleSort: (id: string, ascending: boolean) => void;
}) {
  const [sortField, setSortField] = useState("");
  const [ascending, setAscending] = useState(false);

  const handleSortChanged = (id: string, sortable: boolean = true) => {
    if (!sortable) return;
    if (sortField === id) {
      setAscending(!ascending);
      handleSort(id, !ascending);
    } else {
      setSortField(id);
      setAscending(true);
      handleSort(id, true);
    }
  };

  return (
    <thead>
      <tr>
        {columns.map(({ label, id, sortable = true }, index) => {
          return (
            <th
              onClick={() => handleSortChanged(id as string, sortable)}
              className={`p-2 text-left ${sortable ? "cursor-pointer" : "cursor-default"} no-select`}
              key={index}
            >
              <div className="flex">
                {label}
                {sortable && sortField === id ? (
                  ascending ? (
                    <ChevronUp className="w-4" />
                  ) : (
                    <ChevronDown className="w-4" />
                  )
                ) : (
                  <div className="w-4" />
                )}
              </div>
            </th>
          );
        })}
        <th className="p-2 text-left no-select">Actions</th>
      </tr>
    </thead>
  );
}

function TableBody<T>({ columns, data }: { columns: Column<T>[]; data: any }) {
  return (
    <tbody>
      {data.map((d: any) => {
        return (
          <tr className="nth-[odd]:bg-white/5" key={d.id}>
            {columns.map(({ id, type }, index) => {
              let tData = d[id];
              if (type === "boolean") {
                tData = tData === true ? "Yes" : tData === false ? "No" : "";
              } else if (type === "number") {
                tData = typeof tData === "number" ? tData : tData ? Number(tData) : "";
              } else {
                tData = tData ?? "";
              }
              return (
                <td className="p-2 text-left" key={index}>
                  {tData}
                </td>
              );
            })}
            <th className="flex gap-2 p-2">
              <Button>Edit</Button>
              <Button variant="destructive">Delete</Button>
            </th>
          </tr>
        );
      })}
    </tbody>
  );
}

function Tab<T>({ data, columns }: { data: T[]; columns: Column<T>[] }) {
  const [sortField, setSortField] = useState<keyof T | null>(null);
  const [ascending, setAscending] = useState(true);
  const [sortedData, setSortedData] = useState<T[]>(data);

  useEffect(() => {
    setSortedData(data);
  }, [data]);

  const handleSort = (id: string, asc: boolean) => {
    setSortField(id as keyof T);
    setAscending(asc);
    const newData = handleSortFn(id as keyof T, asc, sortedData, columns);
    setSortedData(newData);
  };

  // Rename handleSort to handleSortFn to avoid shadowing
  function handleSortFn(id: keyof T, ascending: boolean, data: T[], columns: Column<T>[]): T[] {
    const column = columns.find((col) => col.id === id);
    if (column && column.sortFunction) {
      return [...data].sort((a, b) => column.sortFunction!(a, b, ascending));
    }
    return [...data].sort((a, b) => {
      const aVal = a[id];
      const bVal = b[id];
      if (aVal == null) return 1;
      if (bVal == null) return -1;
      if (aVal == null && bVal == null) return 0;
      // Handle booleans
      if (typeof aVal === "boolean" && typeof bVal === "boolean") {
        if (aVal === bVal) return 0;
        return (aVal ? 1 : -1) * (ascending ? 1 : -1);
      }
      // Fallback to string compare
      return (
        aVal.toString().localeCompare(bVal.toString(), "en", {
          numeric: true,
        }) * (ascending ? 1 : -1)
      );
    });
  }

  return (
    <div>
      <Button className="w-fit">Add New</Button>
      <table>
        <TableHead columns={columns} handleSort={handleSort} />
        <TableBody columns={columns} data={sortedData} />
      </table>
    </div>
  );
}

function TabButton({ id, currentId, setTab }: any) {
  return (
    <button
      className={cn(
        "transition-all cursor-pointer hover:bg-white/5 bg-transparent border-transparent border-b-2 p-2 pr-4 pl-4",
        id === currentId && "border-tag-purple bg-tag-purple/10 hover:bg-tag-purple/20"
      )}
      onClick={() => {
        setTab(id);
      }}
    >
      <span className="capitalize">{id}</span>
    </button>
  );
}

export default function Admin() {
  const [currentTab, setCurrentTab] = useState("users");

  const [usersData, setUsersData] = useState<User[]>([]);
  const [locationsData, setLocationsData] = useState<Location[]>([]);
  const [nodesData, setNodesData] = useState<Node[]>([]);
  const [serversData, setServersData] = useState<Server[]>([]);

  useEffect(() => {
    (async () => {
      const pagination: Pagination = {
        $typeName: "common.Pagination",
        page: 1,
        pageSize: 1000, // TODO: make actual pagination and somehow handle sorting
      };

      const userManagerServiceClient = await getAdminUserManagerClient();
      const users = await userManagerServiceClient.getUsers({ pagination });

      setUsersData(users.users);

      const locationManagerServiceClient = await getAdminLocationManagerClient();
      const locations = await locationManagerServiceClient.getLocations({ pagination });

      setLocationsData(locations.locations);

      const nodeManagerServiceClient = await getAdminNodeManagerClient();
      const nodes = await nodeManagerServiceClient.getNodes({ pagination });

      setNodesData(nodes.nodes);

      const serverManagerServiceClient = await getAdminServerManagerClient();
      const servers = await serverManagerServiceClient.getServers({ pagination });

      setServersData(servers.servers);
    })();
  }, []);

  return (
    <div className="flex">
      <Card className="overflow-hidden">
        <div className="flex">
          <TabButton id="users" currentId={currentTab} setTab={setCurrentTab}></TabButton>
          <TabButton id="locations" currentId={currentTab} setTab={setCurrentTab}></TabButton>
          <TabButton id="nodes" currentId={currentTab} setTab={setCurrentTab}></TabButton>
          <TabButton id="servers" currentId={currentTab} setTab={setCurrentTab}></TabButton>
        </div>

        <AnimatePresence mode="wait" initial={false}>
          {currentTab === "users" && (
            <motion.div
              key="users"
              initial={{ opacity: 0, x: 40 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -40 }}
              transition={{ duration: 0.35, ease: "easeInOut" }}
              className="w-full h-full flex items-center justify-center"
            >
              <Tab data={usersData} columns={USERS_COLUMNS}></Tab>
            </motion.div>
          )}
          {currentTab === "locations" && (
            <motion.div
              key="locations"
              initial={{ opacity: 0, x: 40 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -40 }}
              transition={{ duration: 0.35, ease: "easeInOut" }}
              className="w-full h-full flex items-center justify-center"
            >
              <Tab data={locationsData} columns={LOCATIONS_COLUMNS}></Tab>
            </motion.div>
          )}
          {currentTab === "nodes" && (
            <motion.div
              key="nodes"
              initial={{ opacity: 0, x: 40 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -40 }}
              transition={{ duration: 0.35, ease: "easeInOut" }}
              className="w-full h-full flex items-center justify-center"
            >
              <Tab data={nodesData} columns={NODES_COLUMNS}></Tab>
            </motion.div>
          )}
          {currentTab === "servers" && (
            <motion.div
              key="servers"
              initial={{ opacity: 0, x: 40 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -40 }}
              transition={{ duration: 0.35, ease: "easeInOut" }}
              className="w-full h-full flex items-center justify-center"
            >
              <Tab data={serversData} columns={SERVERS_COLUMNS}></Tab>
            </motion.div>
          )}
        </AnimatePresence>
      </Card>
    </div>
  );
}
