import { ChevronDown, ChevronUp } from "lucide-react";
import { Link } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import { Button } from "~/components/ui/button";
import { Card } from "~/components/ui/card";
import { cn } from "~/lib/utils";
import {
  getAdminLocationManagerClient,
  getAdminNodeManagerClient,
  getAdminServerManagerClient,
  getAdminUserManagerClient,
} from "~/lib/api-clients";
import { User, UserManagerService } from "proto-gen-ts/backend/admin/UserManager_pb";
import { Location, LocationManagerService } from "proto-gen-ts/backend/admin/LocationManager_pb";
import { Node, NodeManagerService } from "proto-gen-ts/backend/admin/NodeManager_pb";
import { Server, ServerManagerService } from "proto-gen-ts/backend/admin/ServerManager_pb";
import { Pagination } from "proto-gen-ts/common_pb";
import { Client } from "@connectrpc/connect";
import { Dialog, DialogClose, DialogContent, DialogTrigger } from "~/components/ui/dialog";
import { Input } from "~/components/ui/input";
import { AnimatePresence, motion } from "framer-motion";

interface Column<T> {
  label: string;
  id: keyof T;
  type?: "string" | "number" | "boolean" | "percent" | "memory";
  sortable?: boolean;
  sortFunction?: (a: T, b: T, ascending: boolean) => number;
  fetchFunction?: (a: T) => string | number | boolean;
  hidden?: boolean; // only show in create form
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
  form,
  setForm,
  loading,
  setLoading,
  open,
  setOpen,
  formRef,
  handleInputChange,
  handleSubmit,
}: {
  columns: Column<T>[];
  handleSort: (id: string, ascending: boolean) => void;
  form: any;
  setForm: React.Dispatch<React.SetStateAction<any>>;
  loading: boolean;
  setLoading: React.Dispatch<React.SetStateAction<boolean>>;
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  formRef: React.RefObject<HTMLFormElement | null>;
  handleInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleSubmit: (e: React.FormEvent) => void;
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
        <th colSpan={columns.length + 1} className="p-0">
          <div className="flex w-full items-center justify-between">
            <div className="flex flex-1">
              {columns.map(({ label, id, sortable = true }, index) => (
                <div
                  key={index}
                  onClick={() => handleSortChanged(id as string, sortable)}
                  className={`text-left ${
                    sortable ? "cursor-pointer" : "cursor-default"
                  } no-select flex-1 min-w-0 px-2 py-1`}
                >
                  <div className="flex items-center">
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
                </div>
              ))}
            </div>
            <div className="flex items-center justify-end gap-2 p-2 min-w-[192px]">
              <Dialog open={open} onOpenChange={setOpen}>
                <DialogTrigger asChild>
                  <Button className="w-fit min-w-[80px]" onClick={() => setOpen(true)}>
                    Create
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <form ref={formRef} onSubmit={handleSubmit} className="flex flex-col gap-4">
                    {Object.keys(form)
                      .filter(
                        (key) =>
                          key !== "id" &&
                          key !== "uid" &&
                          key !== "lid" &&
                          key !== "nid" &&
                          key !== "sid" &&
                          key !== "bid" &&
                          key !== "ownerUid"
                      )
                      .map((key) => (
                        <div key={key}>
                          <label className="block mb-1 capitalize">{key}</label>
                          <Input name={key} value={form[key] || ""} onChange={handleInputChange} required />
                        </div>
                      ))}
                    <div className="flex gap-2 justify-end">
                      <DialogClose asChild>
                        <Button type="button" variant="secondary">
                          Cancel
                        </Button>
                      </DialogClose>
                      <Button type="submit" disabled={loading}>
                        {loading ? "Creating..." : "Create"}
                      </Button>
                    </div>
                  </form>
                </DialogContent>
              </Dialog>
            </div>
          </div>
        </th>
      </tr>
    </thead>
  );
}

function TableBody<T>({
  columns,
  data,
  onDelete,
}: {
  columns: Column<T>[];
  data: T[];
  onDelete: (id: string) => void;
}) {
  const idKey = columns[0]?.id;
  return (
    <tbody>
      {data.map((d: T, index) => {
        const idData = d[idKey];
        const id = typeof idData === "string" ? idData : String(idData);
        return (
          <tr className="nth-[odd]:bg-white/5" key={index}>
            <td colSpan={columns.length + 1} className="p-0">
              <div className="flex w-full items-center justify-between">
                <div className="flex flex-1">
                  {columns.map(({ id, type }, index) => {
                    let tData = d[id];
                    let td: string | number | boolean = "";
                    if (type === "boolean") {
                      td = tData === true ? "Yes" : tData === false ? "No" : "";
                    } else if (type === "number") {
                      td = typeof tData === "number" ? tData : tData ? Number(tData) : "";
                    } else if (type === "string") {
                      td = typeof tData === "string" ? tData : tData ? String(tData) : "";
                    } else if (type === "percent") {
                      td = typeof tData === "number" ? `${tData}%` : tData ? `${Number(tData)}%` : "";
                    } else if (type === "memory") {
                      if (typeof tData === "number") {
                        td = `${(tData / 1024).toFixed(2)} GB`;
                      } else if (typeof tData === "string") {
                        const num = parseFloat(tData);
                        td = isNaN(num) ? "" : `${(num / 1024).toFixed(2)} GB`;
                      } else {
                        td = "Error";
                      }
                    } else {
                      td = tData !== undefined && tData !== null ? String(tData) : "";
                    }
                    return (
                      <div className="text-left flex-1 min-w-0 px-2 py-1" key={index}>
                        {td}
                      </div>
                    );
                  })}
                </div>
                <div className="flex items-center justify-end gap-2 p-2 min-w-[192px]">
                  <Button className="min-w-[80px]">Edit</Button>
                  <Button
                    className="min-w-[80px]"
                    variant="destructive"
                    onClick={() => {
                      if (typeof d[idKey] === "string") onDelete(d[idKey]);
                      else console.error("ID is not a string", d[idKey]);
                    }}
                  >
                    Delete
                  </Button>
                </div>
              </div>
            </td>
          </tr>
        );
      })}
    </tbody>
  );
}

function Tab<T>({
  data,
  columns,
  onCreate,
  onDelete,
}: {
  data: T[];
  columns: Column<T>[];
  onCreate: (values: T) => Promise<void>;
  onDelete: (id: string) => void;
}) {
  const [sortField, setSortField] = useState<keyof T | null>(null);
  const [ascending, setAscending] = useState(true);
  const [sortedData, setSortedData] = useState<T[]>(data);
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<any>({});
  const [loading, setLoading] = useState(false);
  const formRef = useRef<HTMLFormElement>(null);

  useEffect(() => {
    setSortedData(data);
  }, [data]);

  const handleSort = (id: string, asc: boolean) => {
    setSortField(id as keyof T);
    setAscending(asc);
    const newData = handleSortFn(id as keyof T, asc, sortedData, columns);
    setSortedData(newData);
  };

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

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!onCreate) return;
    setLoading(true);
    await onCreate(form);
    setLoading(false);
    setForm({});
    setOpen(false);
  };

  return (
    <table className="w-full">
      <TableHead
        columns={columns}
        handleSort={handleSort}
        form={form}
        setForm={setForm}
        loading={loading}
        setLoading={setLoading}
        open={open}
        setOpen={setOpen}
        formRef={formRef}
        handleInputChange={handleInputChange}
        handleSubmit={handleSubmit}
      />
      <TableBody columns={columns} data={sortedData} onDelete={onDelete} />
    </table>
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

export default function AdminPage() {
  const [currentTab, setCurrentTab] = useState("users");

  const pagination: Pagination = {
    $typeName: "common.Pagination",
    page: 1,
    pageSize: 1000, // TODO: make actual pagination and somehow handle sorting
  };

  const [userManagerClient, setUserManagerClient] = useState<Client<typeof UserManagerService>>();
  const [locationManagerClient, setLocationManagerClient] = useState<Client<typeof LocationManagerService>>();
  const [nodeManagerClient, setNodeManagerClient] = useState<Client<typeof NodeManagerService>>();
  const [serverManagerClient, setServerManagerClient] = useState<Client<typeof ServerManagerService>>();

  const [usersData, setUsersData] = useState<User[]>([]);
  const [locationsData, setLocationsData] = useState<Location[]>([]);
  const [nodesData, setNodesData] = useState<Node[]>([]);
  const [serversData, setServersData] = useState<Server[]>([]);

  useEffect(() => {
    (async () => {
      const userManagerServiceClient = await getAdminUserManagerClient();
      setUserManagerClient(userManagerServiceClient);
      const users = await userManagerServiceClient.getUsers({ pagination });
      setUsersData(users.users);

      const locationManagerServiceClient = await getAdminLocationManagerClient();
      setLocationManagerClient(locationManagerServiceClient);
      const locations = await locationManagerServiceClient.getLocations({ pagination });
      setLocationsData(locations.locations);

      const nodeManagerServiceClient = await getAdminNodeManagerClient();
      setNodeManagerClient(nodeManagerServiceClient);
      const nodes = await nodeManagerServiceClient.getNodes({ pagination });
      setNodesData(nodes.nodes);

      const serverManagerServiceClient = await getAdminServerManagerClient();
      setServerManagerClient(serverManagerServiceClient);
      const servers = await serverManagerServiceClient.getServers({ pagination });
      setServersData(servers.servers);
    })();
  }, []);

  const tryRefreshUsers = async () => {
    const users = await userManagerClient?.getUsers({ pagination });

    if (!users) {
      console.error("Failed to fetch users");
      return;
    }

    setUsersData(users.users);
  };

  const tryRefreshLocations = async () => {
    const locations = await locationManagerClient?.getLocations({ pagination });

    if (!locations) {
      console.error("Failed to fetch locations");
      return;
    }

    setLocationsData(locations.locations);
  };

  const tryRefreshNodes = async () => {
    const nodes = await nodeManagerClient?.getNodes({ pagination });

    if (!nodes) {
      console.error("Failed to fetch nodes");
      return;
    }

    setNodesData(nodes.nodes);
  };

  const tryRefreshServers = async () => {
    const servers = await serverManagerClient?.getServers({ pagination });

    if (!servers) {
      console.error("Failed to fetch servers");
      return;
    }

    setServersData(servers.servers);
  };

  const handleCreateUser = async (user: User) => {
    const res = await userManagerClient?.createUser({ user });
    if (!res || !res.success) {
      console.error("Failed to create user");
      return;
    }
    await tryRefreshUsers();
  };

  const handleCreateLocation = async (location: Location) => {
    const res = await locationManagerClient?.createLocation({ location });
    if (!res || !res.success) {
      console.error("Failed to create location");
      return;
    }
    await tryRefreshLocations();
  };

  const handleCreateNode = async (node: Node) => {
    const res = await nodeManagerClient?.createNode({ node });
    if (!res || !res.success) {
      console.error("Failed to create node");
      return;
    }
    await tryRefreshNodes();
  };

  const handleCreateServer = async (server: Server) => {
    const res = await serverManagerClient?.createServer({ server });
    if (!res || !res.success) {
      console.error("Failed to create server");
      return;
    }
    await tryRefreshServers();
  };

  const handleDeleteUser = async (uid: string) => {
    if (!userManagerClient) return;
    await userManagerClient.deleteUser({ uid });
    await tryRefreshUsers();
  };
  const handleDeleteLocation = async (lid: string) => {
    if (!locationManagerClient) return;
    await locationManagerClient.deleteLocation({ lid });
    await tryRefreshLocations();
  };
  const handleDeleteNode = async (nid: string) => {
    if (!nodeManagerClient) return;
    await nodeManagerClient.deleteNode({ nid });
    await tryRefreshNodes();
  };
  const handleDeleteServer = async (sid: string) => {
    if (!serverManagerClient) return;
    await serverManagerClient.deleteServer({ sid });
    await tryRefreshServers();
  };

  return (
    <div className="min-h-screen min-w-screen flex items-center justify-center">
      <div className="w-[calc(100vw-4rem)] h-[calc(100vh-4rem)] max-w-[1800px] max-h-[1000px] m-8 flex flex-col">
        <div className="mb-4">
          <Link to="/" className="text-tag-purple hover:underline font-semibold">
            ← Back to Dashboard
          </Link>
        </div>
        <Card className="overflow-hidden w-full h-full flex flex-col bg-server-card">
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
                className="w-full h-full"
              >
                <Tab data={usersData} columns={USERS_COLUMNS} onCreate={handleCreateUser} onDelete={handleDeleteUser} />
              </motion.div>
            )}
            {currentTab === "locations" && (
              <motion.div
                key="locations"
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -40 }}
                transition={{ duration: 0.35, ease: "easeInOut" }}
                className="w-full h-full"
              >
                <Tab
                  data={locationsData}
                  columns={LOCATIONS_COLUMNS}
                  onCreate={handleCreateLocation}
                  onDelete={handleDeleteLocation}
                />
              </motion.div>
            )}
            {currentTab === "nodes" && (
              <motion.div
                key="nodes"
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -40 }}
                transition={{ duration: 0.35, ease: "easeInOut" }}
                className="w-full h-full"
              >
                <Tab data={nodesData} columns={NODES_COLUMNS} onCreate={handleCreateNode} onDelete={handleDeleteNode} />
              </motion.div>
            )}
            {currentTab === "servers" && (
              <motion.div
                key="servers"
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -40 }}
                transition={{ duration: 0.35, ease: "easeInOut" }}
                className="w-full h-full"
              >
                <Tab
                  data={serversData}
                  columns={SERVERS_COLUMNS}
                  onCreate={handleCreateServer}
                  onDelete={handleDeleteServer}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </Card>
      </div>
    </div>
  );
}
