import {ServerStatusType} from "proto-gen-ts/daemon_pb";
import {useEffect, useState} from "react";
import {useSearchParams} from "react-router";

const server = {
    id: "1",
    name: "My Awesome Server 1",
    status: ServerStatusType.ONLINE,
    node: "Node Alpha",
    cpuUsage: 25,
    memoryUsage: {
        used: 2048,
        total: 8192,
    },
    diskUsage: {
        used: 51200,
        total: 102400,
    },
    ip: "192.168.1.100",
    port: 25565,
    location: "US East",
    game: "Minecraft",
    uptime: "2d 15h 30m",
    console: [
        {time: "10:15:32", content: "Server started on port 25565"},
        {time: "10:15:33", content: "Loading world..."},
        {time: "10:15:40", content: "World loaded in 7s"},
        {time: "10:15:41", content: "Server ready for connections"},
        {time: "10:16:20", content: "Player1 joined the game"},
        {time: "10:20:15", content: "Player1: Hello everyone!"},
        {time: "10:25:30", content: "Player2 joined the game"},
        {time: "10:30:45", content: "Autosaving world..."},
        {time: "10:30:47", content: "World autosaved"},
    ],
};

const useServer = () => {
    const [searchParams] = useSearchParams();
    const tabParam = searchParams.get("tab");

    const [activeTab, setActiveTab] = useState(tabParam || "console");
    const [command, setCommand] = useState("");

    useEffect(() => {
        if (
            tabParam &&
            ["console", "files", "activity", "settings"].includes(tabParam)
        ) {
            setActiveTab(tabParam);
        }
    }, [tabParam]);

    // Handle command submit
    const handleCommandSubmit = (e: React.FormEvent): void => {
        e.preventDefault();
        if (!command.trim()) return;

        // In a real app, you would send this command to the server
        console.log(`Sending command: ${command}`);
        setCommand("");
    };

    return {
        activeTab,
        setActiveTab,
        command,
        setCommand,
        handleCommandSubmit,
        server,
    };
};

export default useServer;
