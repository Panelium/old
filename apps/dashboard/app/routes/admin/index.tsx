import { AnimatePresence, motion } from "framer-motion";
import { ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import { Card, CardHeader, CardTitle } from "~/components/ui/card"
import { cn } from "~/lib/utils";

interface Column {
    label: string,
    id: string,
}

const USERS_COLUMNS: Column[] = [
    { label: "ID", id: "id" },
    { label: "Username", id: "username" },
    { label: "Email", id: "email" },
    { label: "Password", id: "password" },
]

const USERS_DATA = [
    {
        "id": "0",
        "username": "test",
        "email": "test@test.test",
        "password": "test123",
    },
    {
        "id": 1,
        "username": "test",
        "email": "test@test.test",
        "password": "test123",
    },
    {
        "id": 2,
        "username": "test",
        "email": "test@test.test",
        "password": "test123",
    },
    {
        "id": 3,
        "username": "test",
        "email": "test@test.test",
        "password": "test123",
    },
    {
        "id": 4,
        "username": "test",
        "email": "test@test.test",
        "password": "test123",
    }
]

function TableHead({columns, handleSort}: {columns: Column[], handleSort: (id: string, ascending: boolean) => void}) {
    const [sortField, setSortField] = useState("");
    const [ascending, setAscending] = useState(false);

    const handleSortChanged = (id: string) => {
        setSortField(id);
        if (sortField === id) {
            setAscending(!ascending);
        }
        handleSort(sortField, ascending);
    }

    return (
        <thead>
            <tr>
                {columns.map(({label, id}) => {
                    return <th onClick={() => {handleSortChanged(id);}} className="p-2 text-left cursor-pointer" key={id}>
                        <div className="flex">
                            {label}
                            {
                                sortField === id ? ascending ? <ChevronUp className="w-4" /> : <ChevronDown className="w-4" /> : <div className="w-4"/>
                            }
                        </div>
                    </th>;
                })}
                <th className="p-2 text-left">
                    Actions
                </th>
            </tr>
        </thead>
    );
}

function TableBody({columns, data}: {columns: Column[], data: any}) {
    return (
        <tbody>
            {data.map((d: any) => {
                return (
                    <tr className="nth-[odd]:bg-white/5" key={d.id}>
                        {columns.map(({id}) => {
                            const tData = d[id] ? d[id] : "";
                            return <td className="p-2 text-left" key={id}>{tData}</td>;
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

function handleSort(id: string, ascending: boolean, data: any) {
    const sorted = [...data].sort((a: any, b: any) => {
        if (a[id] === null) return 1;
        if (b[id] === null) return -1;
        if (a[id] === null && b[id] === null) return 0;
        return (
            a[id].toString().localeCompare(b[id].toString(), "en", {
                numeric: true,
            }) * (ascending ? 1 : -1)
        );
    });
    return sorted;
}


function Tab({data, columns}: any) {
    return (
        <div>
            <Button className="w-fit">
                Add New
            </Button>
            <table>
            <TableHead columns={columns} handleSort={(id, ascending) => {
                handleSort(id, ascending, data)
            }}>
            </TableHead>
            <TableBody columns={columns} data={data}>
            </TableBody>
            </table>
        </div>
    );
}

function TabButton({id, currentId, setTab}: any) {
    return (
        <button className={cn(
            "transition-all cursor-pointer hover:bg-white/5 bg-transparent border-transparent border-b-2 p-2 pr-4 pl-4",
            id === currentId && "border-tag-purple bg-tag-purple/10 hover:bg-tag-purple/20"
        )} onClick={() => {
            setTab(id);
        }}>
            <span className="capitalize">{id}</span>
        </button>
    );
}

export default function Admin() {
    const [currentTab, setCurrentTab] = useState("users");

    const [usersData, setUsersData] = useState(USERS_DATA);
    const [locationsData, setLocationsData] = useState(USERS_DATA);
    const [nodesData, setNodesData] = useState(USERS_DATA);
    const [serversData, setServersData] = useState(USERS_DATA);

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
                {currentTab === "users" && <motion.div
                        key="users"
                        initial={{ opacity: 0, x: 40 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0, x: -40 }}
                        transition={{ duration: 0.35, ease: "easeInOut" }}
                        className="w-full h-full flex items-center justify-center"
                    >
                        <Tab data={usersData} columns={USERS_COLUMNS}></Tab>
                    </motion.div>
                }
                {currentTab === "locations" && <motion.div
                        key="locations"
                        initial={{ opacity: 0, x: 40 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0, x: -40 }}
                        transition={{ duration: 0.35, ease: "easeInOut" }}
                        className="w-full h-full flex items-center justify-center"
                    >
                        <Tab data={locationsData} columns={USERS_COLUMNS}></Tab>
                    </motion.div>
                }
                {currentTab === "nodes" && <motion.div
                        key="nodes"
                        initial={{ opacity: 0, x: 40 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0, x: -40 }}
                        transition={{ duration: 0.35, ease: "easeInOut" }}
                        className="w-full h-full flex items-center justify-center"
                    >
                        <Tab data={nodesData} columns={USERS_COLUMNS}></Tab>
                    </motion.div>
                }
                {currentTab === "servers" && <motion.div
                        key="servers"
                        initial={{ opacity: 0, x: 40 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0, x: -40 }}
                        transition={{ duration: 0.35, ease: "easeInOut" }}
                        className="w-full h-full flex items-center justify-center"
                    >
                        <Tab data={serversData} columns={USERS_COLUMNS}></Tab>
                    </motion.div>
                }
            </AnimatePresence>
        </Card>
        </div>
    );
}