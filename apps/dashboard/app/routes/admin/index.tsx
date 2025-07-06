import { ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import { Card, CardTitle } from "~/components/ui/card"

interface Column {
    label: string,
    id: string,
}

const COLUMNS: Column[] = [
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


export default function Admin() {
    const [data, setData] = useState(USERS_DATA);

    const handleSort = (id: string, ascending: boolean) => {
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
        setData(sorted);
    }

    return (
        <>
        <Card>
            <CardTitle className="text-3xl font-bold">
                Users:
            </CardTitle>
            <Button className="w-fit">
                Add New
            </Button>
            <table>
            <TableHead columns={COLUMNS} handleSort={handleSort}>
            </TableHead>
            <TableBody columns={COLUMNS} data={data}>
            </TableBody>
            </table>
        </Card>
        </>
    );
}