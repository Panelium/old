import { ChevronDown, ChevronUp } from "lucide-react";
import { Link } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import { Button } from "~/components/ui/button";
import { Card } from "~/components/ui/card";
import { cn } from "~/lib/utils";
import {
  getAdminBlueprintManagerClient,
  getAdminLocationManagerClient,
  getAdminNodeAllocationManagerClient,
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
import { Dialog, DialogClose, DialogContent, DialogTitle, DialogTrigger } from "~/components/ui/dialog";
import { Input } from "~/components/ui/input";
import { AnimatePresence, motion } from "framer-motion";
import { NodeAllocation, NodeAllocationManagerService } from "proto-gen-ts/backend/admin/NodeAllocationManager_pb";
import { Blueprint, BlueprintManagerService } from "proto-gen-ts/backend/admin/BlueprintManager_pb";
import { Label } from "~/components/ui/label";

type ColumnType = User | Location | Node | NodeAllocation | Server | Blueprint;

interface Column<T extends ColumnType> {
  label: string;
  id: keyof T;
  type?: "string" | "number" | "boolean" | "percent" | "memory";
  optional?: boolean; // if true, the field is optional in the create form
  sortable?: boolean;
  sortFunction?: (a: T, b: T, ascending: boolean) => number;
  fetchFunction?: (a: T) => string | number | boolean;
  setFunction?: (a: T, value: any) => T;
  hidden?: boolean; // only show in create form
  linksTo?: () => { value: string | number; label: string }[];
}

function TableHead<T extends ColumnType>({
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

  const [fileName, setFileName] = useState("");

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
            <div className="flex items-center justify-end gap-2 p-2 min-w-[220px]">
              <Dialog open={open} onOpenChange={setOpen}>
                <DialogTrigger asChild>
                  <Button className="w-fit min-w-[80px]" onClick={() => setOpen(true)}>
                    Create
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogTitle className="text-lg font-semibold mb-4">Create</DialogTitle>
                  <form ref={formRef} onSubmit={handleSubmit} className="flex flex-col gap-4">
                    {/* Blueprint JSON upload (only for blueprints) */}
                    {columns[0].label === "Blueprint ID" ? (
                      <div>
                        <Label className="block mb-1">Upload Blueprint (.bp)</Label>
                        <div className="flex items-center gap-2">
                          <Button asChild type="button" variant="outline">
                            <label htmlFor="bp-upload" className="cursor-pointer">
                              Choose File
                              <Input
                                id="bp-upload"
                                type="file"
                                accept=".bp,application/json"
                                className="hidden"
                                onChange={async (e) => {
                                  const file = e.target.files?.[0];
                                  if (!file) return;
                                  if (!file.name.endsWith(".bp")) {
                                    alert("Please select a .bp file.");
                                    return;
                                  }

                                  const text = await file.text();
                                  let validJson = true;
                                  try {
                                    JSON.parse(text);
                                  } catch (err) {
                                    alert("Invalid JSON in .bp file.");
                                    validJson = false;
                                  }
                                  if (!validJson) return;

                                  setFileName(file.name);
                                  setForm(text);
                                }}
                              />
                            </label>
                          </Button>
                          <span className="text-gray-500 text-sm" id="bp-upload-filename">
                            {fileName || "No file chosen"}
                          </span>
                        </div>
                      </div>
                    ) : (
                      columns
                        .filter(
                          (col, index) => index > 0 // Skip the first column (ID)
                        )
                        .map((column, index) => (
                          <div key={index}>
                            <label className="block mb-1 capitalize">{column.label}</label>
                            {column.linksTo ? (
                              <select
                                name={column.id.toString()}
                                value={form[column.id] ?? ""}
                                onChange={(e) => setForm({ ...form, [column.id]: e.target.value })}
                                required
                                className="block w-full border rounded px-2 py-1"
                              >
                                <option value="" disabled>
                                  Select {column.label}
                                </option>
                                {column.linksTo().map((opt) => (
                                  <option key={opt.value} value={opt.value}>
                                    {opt.label}
                                  </option>
                                ))}
                              </select>
                            ) : column.type === "boolean" ? (
                              <Input
                                name={column.id.toString()}
                                type="checkbox"
                                checked={!!form[column.id]}
                                onChange={handleInputChange}
                              />
                            ) : (
                              <Input
                                name={column.id.toString()}
                                type={column.type === "number" ? "number" : "text"}
                                value={
                                  column.type === "number"
                                    ? form[column.id] !== undefined && form[column.id] !== ""
                                      ? form[column.id].toString()
                                      : ""
                                    : form[column.id] || ""
                                }
                                onChange={handleInputChange}
                                required={!column.optional}
                              />
                            )}
                          </div>
                        ))
                    )}
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

function TableBody<T extends ColumnType>({
  columns,
  data,
  onDelete,
}: {
  columns: Column<T>[];
  data: T[];
  onDelete: ((id: string) => void) | ((id: number) => void);
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
                  {columns.map(({ id, type, linksTo, fetchFunction }, index) => {
                    let tData = d[id];
                    let td: string | number | boolean = "";
                    if (fetchFunction) {
                      td = fetchFunction(d);
                    } else if (linksTo) {
                      const opts = linksTo();
                      const found = opts.find((opt) => String(opt.value) === String(tData));
                      td = found ? found.label : String(tData ?? "");
                    } else if (type === "boolean") {
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
                      <div
                        className="text-left flex-1 min-w-0 px-2 py-1 truncate"
                        key={index}
                        title={typeof td === "string" ? td : undefined}
                      >
                        {td}
                      </div>
                    );
                  })}
                </div>
                <div className="flex items-center justify-end gap-2 p-2 min-w-[220px]">
                  {columns[0].label === "Node ID" && (
                    <Button
                      variant="secondary"
                      onClick={async () => {
                        const nodeManagerClient = await getAdminNodeManagerClient();
                        if (!nodeManagerClient) {
                          alert("Failed to send request to generate token");
                          return;
                        }
                        let res = await nodeManagerClient.generateBackendToken({
                          nid: id,
                        });
                        if (!res) {
                          alert("Failed to generate token");
                          return;
                        }
                        if (!res.success) {
                          const regenerate = confirm(
                            "Backend token for this node seems to already exist. Do you want to regenerate it?"
                          );
                          if (!regenerate) {
                            return;
                          }
                          res = await nodeManagerClient.generateBackendToken({
                            nid: id,
                            regenerate: true,
                          });
                          if (!res || !res.success) {
                            alert("Failed to regenerate token");
                            return;
                          }
                        }

                        const token = res.backendToken;
                        if (!token) {
                          alert("Failed to generate token");
                          return;
                        }

                        await navigator.clipboard.writeText(token);
                        alert("Backend token successfully generated and copied to clipboard.");
                      }}
                    >
                      T
                    </Button>
                  )}
                  <Button>Edit</Button>
                  <Button
                    variant="destructive"
                    onClick={() => {
                      onDelete(d[idKey] as never); // fuck typescript
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

function Tab<T extends ColumnType>({
  data,
  columns,
  onCreate,
  onDelete,
}: {
  data: T[];
  columns: Column<T>[];
  onCreate: ((values: T) => Promise<void>) | ((json: string) => Promise<void>);
  onDelete: ((id: string) => void) | ((id: number) => void);
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
      if (typeof aVal === "boolean" && typeof bVal === "boolean") {
        if (aVal === bVal) return 0;
        return (aVal ? 1 : -1) * (ascending ? 1 : -1);
      }
      return (
        aVal.toString().localeCompare(bVal.toString(), "en", {
          numeric: true,
        }) * (ascending ? 1 : -1)
      );
    });
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.type === "file") return;

    const { name, value, type, checked } = e.target;
    const column = columns.find((col) => col.id === name);
    let parsedValue: any = value;
    if (column?.type === "number") {
      parsedValue = value === "" ? "" : Number(value);
    } else if (column?.type === "boolean") {
      parsedValue = type === "checkbox" ? checked : value === "true";
    }
    // setFunction applies on submit, not on change
    if (column?.setFunction) {
      setForm({ ...form, [name]: value });
    } else {
      setForm({ ...form, [name]: parsedValue });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!onCreate) return;
    setLoading(true);
    let f = form;
    if (typeof form !== "string") {
      let submitForm = { ...form };
      columns.forEach((column) => {
        if (column.setFunction && submitForm[column.id] !== undefined) {
          submitForm = column.setFunction(submitForm, submitForm[column.id]);
        }
      });
      f = submitForm;
    }
    await onCreate(f);
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
      <span className="capitalize">{id.replaceAll("_", " ")}</span>
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
  const [nodeAllocationManagerClient, setNodeAllocationManagerClient] =
    useState<Client<typeof NodeAllocationManagerService>>();
  const [serverManagerClient, setServerManagerClient] = useState<Client<typeof ServerManagerService>>();
  const [blueprintManagerClient, setBlueprintManagerClient] = useState<Client<typeof BlueprintManagerService>>();

  const [usersData, setUsersData] = useState<User[]>([]);
  const [locationsData, setLocationsData] = useState<Location[]>([]);
  const [nodesData, setNodesData] = useState<Node[]>([]);
  const [nodeAllocationsData, setNodeAllocationsData] = useState<NodeAllocation[]>([]);
  const [serversData, setServersData] = useState<Server[]>([]);
  const [blueprintsData, setBlueprintsData] = useState<Blueprint[]>([]);

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

      const nodeAllocationManagerServiceClient = await getAdminNodeAllocationManagerClient();
      setNodeAllocationManagerClient(nodeAllocationManagerServiceClient);
      const nodeAllocations = await nodeAllocationManagerServiceClient.getNodeAllocations({ pagination });
      setNodeAllocationsData(nodeAllocations.nodeAllocations);

      const serverManagerServiceClient = await getAdminServerManagerClient();
      setServerManagerClient(serverManagerServiceClient);
      const servers = await serverManagerServiceClient.getServers({ pagination });
      setServersData(servers.servers);

      const blueprintManagerServiceClient = await getAdminBlueprintManagerClient();
      setBlueprintManagerClient(blueprintManagerServiceClient);
      const blueprints = await blueprintManagerServiceClient.getBlueprints({ pagination });
      setBlueprintsData(blueprints.blueprints);
    })();
  }, []);

  function USERS_COLUMNS(): Column<User>[] {
    return [
      { label: "User ID", id: "uid", sortable: false },
      { label: "Username", id: "username" },
      { label: "Email", id: "email" },
      { label: "Admin", id: "admin", type: "boolean" },
      { label: "MFA Needed", id: "mfaNeeded", type: "boolean" },
    ];
  }

  function LOCATIONS_COLUMNS(): Column<Location>[] {
    return [
      { label: "Location ID", id: "lid", sortable: false },
      { label: "Name", id: "name" },
    ];
  }

  function NODES_COLUMNS(locationsData: Location[]): Column<Node>[] {
    return [
      { label: "Node ID", id: "nid", sortable: false },
      { label: "Name", id: "name" },
      {
        label: "Location",
        id: "lid",
        linksTo: () => (locationsData ? locationsData.map((l) => ({ value: l.lid, label: l.name })) : []),
      },
      { label: "FQDN", id: "fqdn" },
      { label: "Daemon Port", id: "daemonPort", type: "number" },
      { label: "HTTPS", id: "https", type: "boolean" },
      { label: "Max CPU", id: "maxCpu", type: "number" },
      { label: "Max RAM", id: "maxRam", type: "number" },
      { label: "Max Storage", id: "maxStorage", type: "number" },
      { label: "Max Swap", id: "maxSwap", type: "number" },
    ];
  }

  function NODE_ALLOCATIONS_COLUMNS(nodesData: Node[], serversData: Server[]): Column<NodeAllocation>[] {
    return [
      { label: "ID", id: "id", sortable: true, type: "number" },
      {
        label: "Node",
        id: "nid",
        linksTo: () => (nodesData ? nodesData.map((n) => ({ value: n.nid, label: n.name })) : []),
      },
      {
        label: "Allocation",
        id: "ipAllocation",
        sortable: false,
        fetchFunction: ({ ipAllocation: a }) => (a ? `${a.ip}:${a.port}` : "Error"),
        setFunction: (na, value): NodeAllocation => {
          const [ip, port] = value.split(":");
          return {
            ...na,
            ipAllocation: {
              ...na.ipAllocation,
              $typeName: "common.IPAllocation",
              ip,
              port: parseInt(port, 10),
            },
          };
        },
      },
      {
        label: "Server ID",
        id: "sid",
        sortable: false,
        optional: true,
        linksTo: () => (serversData ? serversData.map((s) => ({ value: s.sid, label: s.name })) : []),
      },
    ];
  }

  function SERVERS_COLUMNS(nodesData: Node[], blueprintsData: Blueprint[], usersData: User[]): Column<Server>[] {
    return [
      { label: "Server ID", id: "sid", sortable: false },
      { label: "Name", id: "name" },
      {
        label: "Owner",
        id: "ownerUid",
        linksTo: () => (usersData ? usersData.map((u) => ({ value: u.uid, label: u.username })) : []),
      },
      {
        label: "Node",
        id: "nid",
        linksTo: () => (nodesData ? nodesData.map((n) => ({ value: n.nid, label: n.name })) : []),
      },
      {
        label: "Blueprint",
        id: "bid",
        linksTo: () => (blueprintsData ? blueprintsData.map((b) => ({ value: b.bid, label: b.name })) : []),
      },
      { label: "Docker Image", id: "dockerImage" },
    ];
  }

  function BLUEPRINTS_COLUMNS(): Column<Blueprint>[] {
    return [
      { label: "Blueprint ID", id: "bid", sortable: false },
      { label: "Name", id: "name" },
      { label: "Category", id: "category" },
      { label: "Description", id: "description" },
      { label: "Version", id: "version" },
    ];
  }

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

  const tryRefreshNodeAllocations = async () => {
    const nodeAllocations = await nodeAllocationManagerClient?.getNodeAllocations({ pagination });

    if (!nodeAllocations) {
      console.error("Failed to fetch node allocations");
      return;
    }

    setNodeAllocationsData(nodeAllocations.nodeAllocations);
  };

  const tryRefreshServers = async () => {
    const servers = await serverManagerClient?.getServers({ pagination });

    if (!servers) {
      console.error("Failed to fetch servers");
      return;
    }

    setServersData(servers.servers);
  };

  const tryRefreshBlueprints = async () => {
    const blueprints = await blueprintManagerClient?.getBlueprints({ pagination });

    if (!blueprints) {
      console.error("Failed to fetch blueprints");
      return;
    }

    setBlueprintsData(blueprints.blueprints);
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

  const handleCreateNodeAllocation = async (nodeAllocation: NodeAllocation) => {
    const res = await nodeAllocationManagerClient?.createNodeAllocation({ nodeAllocation });
    if (!res || !res.success) {
      console.error("Failed to create node allocation");
      return;
    }
    await tryRefreshNodeAllocations();
  };

  const handleCreateServer = async (server: Server) => {
    const res = await serverManagerClient?.createServer({ server });
    if (!res || !res.success) {
      console.error("Failed to create server");
      return;
    }
    await tryRefreshServers();
  };

  const handleCreateBlueprint = async (json: string) => {
    const res = await blueprintManagerClient?.createBlueprint({
      blueprintOrJson: { value: json, case: "blueprintJson" },
    });
    if (!res || !res.success) {
      console.error("Failed to create blueprint");
      return;
    }
    await tryRefreshBlueprints();
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
  const handleDeleteNodeAllocation = async (id: number) => {
    if (!nodeAllocationManagerClient) return;
    await nodeAllocationManagerClient.deleteNodeAllocation({ id });
    await tryRefreshNodeAllocations();
  };
  const handleDeleteServer = async (sid: string) => {
    if (!serverManagerClient) return;
    await serverManagerClient.deleteServer({ sid });
    await tryRefreshServers();
  };
  const handleDeleteBlueprint = async (bid: string) => {
    if (!blueprintManagerClient) return;
    await blueprintManagerClient.deleteBlueprint({ bid });
    await tryRefreshBlueprints();
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
            <TabButton id="node_allocations" currentId={currentTab} setTab={setCurrentTab}></TabButton>
            <TabButton id="servers" currentId={currentTab} setTab={setCurrentTab}></TabButton>
            <TabButton id="blueprints" currentId={currentTab} setTab={setCurrentTab}></TabButton>
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
                <Tab
                  data={usersData}
                  columns={USERS_COLUMNS()}
                  onCreate={handleCreateUser}
                  onDelete={handleDeleteUser}
                />
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
                  columns={LOCATIONS_COLUMNS()}
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
                transition={{ duration: 0.35, ease: "easeInOut" }}
                className="w-full h-full"
              >
                <Tab
                  data={nodesData}
                  columns={NODES_COLUMNS(locationsData)}
                  onCreate={handleCreateNode}
                  onDelete={handleDeleteNode}
                />
              </motion.div>
            )}
            {currentTab === "node_allocations" && (
              <motion.div
                key="node_allocations"
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -40 }}
                transition={{ duration: 0.35, ease: "easeInOut" }}
                className="w-full h-full"
              >
                <Tab
                  data={nodeAllocationsData}
                  columns={NODE_ALLOCATIONS_COLUMNS(nodesData, serversData)}
                  onCreate={handleCreateNodeAllocation}
                  onDelete={handleDeleteNodeAllocation}
                />
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
                  columns={SERVERS_COLUMNS(nodesData, blueprintsData, usersData)}
                  onCreate={handleCreateServer}
                  onDelete={handleDeleteServer}
                />
              </motion.div>
            )}
            {currentTab === "blueprints" && (
              <motion.div
                key="blueprints"
                initial={{ opacity: 0, x: 40 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -40 }}
                transition={{ duration: 0.35, ease: "easeInOut" }}
                className="w-full h-full"
              >
                <Tab
                  data={blueprintsData}
                  columns={BLUEPRINTS_COLUMNS()}
                  onCreate={handleCreateBlueprint}
                  onDelete={handleDeleteBlueprint}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </Card>
      </div>
    </div>
  );
}
