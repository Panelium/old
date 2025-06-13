import React from 'react';
import {Users} from 'lucide-react';
import ServerCardGrid, {type Server} from './ServerCardGrid';
import type {OverviewCardProps} from './OverviewCard';
import OverviewCardGrid from './OverviewCardGrid';
import {ServerStatusType} from 'proto-gen-ts/daemon_pb';

// Mock data for demonstration
const mockServers: Server[] = [
    {
        id: '1',
        name: 'Minecraft SMP',
        status: ServerStatusType.SERVER_STATUS_ONLINE,
        description: 'Survival multiplayer server',
        cpuUsage: 45,
        memoryUsage: {
            used: 2048 * 1024,
            total: 4096 * 1024,
        },
        game: 'Minecraft',
        players: {
            online: 12,
            max: 50,
        },
        ip: '192.168.1.1',
        port: 25565,
    },
    {
        id: '2',
        name: 'Web Server',
        status: ServerStatusType.SERVER_STATUS_STARTING,
        description: 'NGINX web server',
        cpuUsage: 12,
        memoryUsage: {
            used: 512 * 1024,
            total: 1024 * 1024,
        },
        game: 'Nginx',
        ip: '192.168.1.6',
        port: 80,
    },
    {
        id: '3',
        name: 'Terraria Adventure',
        status: ServerStatusType.SERVER_STATUS_OFFLINE,
        description: 'Terraria modded adventure server',
        cpuUsage: 0,
        memoryUsage: {
            used: 0,
            total: 2048 * 1024,
        },
        game: 'Terraria',
        players: {
            online: 0,
            max: 16,
        },
        ip: '192.168.1.3',
        port: 7777,
    },
    {
        id: '4',
        name: 'Rust Survival',
        status: ServerStatusType.SERVER_STATUS_STOPPING,
        description: 'Rust survival server',
        cpuUsage: 30,
        memoryUsage: {
            used: 1024 * 1024,
            total: 2048 * 1024,
        },
        game: 'Rust',
        players: {
            online: 5,
            max: 100,
        },
        ip: '192.168.1.23',
        port: 28015,
    },
    {
        id: '5',
        name: 'ARK: Survival Evolved',
        status: ServerStatusType.SERVER_STATUS_UNKNOWN,
        description: 'ARK survival server',
        cpuUsage: 60,
        memoryUsage: {
            used: 3072 * 1024,
            total: 8192 * 1024,
        },
        game: 'ARK',
        players: {
            online: 20,
            max: 70,
        },
        ip: '192.168.1.13',
        port: 7778,
    },
];

export default function DashboardOverviewPage() {
    // Calculate actual stats from server data
    const totalServers = mockServers.length;
    const onlineServers = mockServers.filter(s => s.status === ServerStatusType.SERVER_STATUS_ONLINE).length;

    // Calculate total players across all servers
    const totalPlayers = mockServers.reduce((acc, server) => {
        return acc + (server.players?.online || 0);
    }, 0);
    const maxPlayers = mockServers.reduce((acc, server) => {
        return acc + (server.players?.max || 0);
    }, 0);

    // Calculate average CPU usage from online servers
    const onlineServersArray = mockServers.filter(s => s.status === ServerStatusType.SERVER_STATUS_ONLINE);
    const avgCpuUsage = onlineServersArray.length
        ? Math.round(onlineServersArray.reduce((acc, server) => acc + server.cpuUsage, 0) / onlineServersArray.length)
        : 0;

    // Calculate total memory usage
    const totalMemoryUsed = mockServers.reduce((acc, server) => acc + server.memoryUsage.used, 0);
    const totalMemoryAllocated = mockServers.reduce((acc, server) => acc + server.memoryUsage.total, 0);

    // Convert bytes to GB for display
    const usedMemoryGB = (totalMemoryUsed / (1024 * 1024 * 1024)).toFixed(1);
    const totalMemoryGB = (totalMemoryAllocated / (1024 * 1024 * 1024)).toFixed(1);

    const overviewCardsData: OverviewCardProps[] = [
        {
            title: "Player Activity",
            content: {
                title: totalPlayers.toString(),
                subtitle: `across ${onlineServers} online servers`,
                icon: Users,
            },
            footer: `${totalPlayers}/${maxPlayers} total capacity`,
        },
        {
            title: "CPU Usage",
            gauge: {
                value: avgCpuUsage,
                maxValue: 100,
                size: "sm",
                unit: "%",
                label: "Average",
                subtitle: "Avg across 2 servers",
            },
            footer: `Peak: ${Math.max(...onlineServersArray.map(s => s.cpuUsage), 0)}% on ${
                onlineServersArray.length > 0
                    ? onlineServersArray.reduce((max, server) =>
                        server.cpuUsage > max.cpuUsage ? server : max
                    ).name
                    : "N/A"
            }`,
        },
        {
            title: "Memory Usage",
            gauge: {
                value: parseFloat(usedMemoryGB),
                maxValue: parseFloat(totalMemoryGB),
                size: "sm",
                unit: "GB",
                label: "Used",
                subtitle: `of 0.0 GB allocated`,
            }
        },
    ];

    return (
        <div className="container py-6 px-4 mx-auto max-w-7xl">
            <div className="mb-6 no-select">
                <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">Dashboard Overview</h1>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-2">
                    Welcome back! Here's what's happening with your servers.
                </p>
            </div>

            <OverviewCardGrid cards={overviewCardsData}/>

            <ServerCardGrid servers={mockServers}/>
        </div>
    );
}

