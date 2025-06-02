import React from 'react';
import {Card, CardContent, CardHeader, CardTitle} from '~/components/ui/card';
import {ResourceGauge} from '~/components/dashboard/resource-gauge';
import {ServerCard} from '~/components/dashboard/server-card';
import {Plus, Users} from 'lucide-react';

// Mock data for demonstration
const mockServers = [
    {
        id: '1',
        name: 'Minecraft SMP',
        status: 'Online',
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
        status: 'Online',
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
        status: 'Offline',
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
];

export default function DashboardOverviewPage() {
    // Calculate actual stats from server data
    const totalServers = mockServers.length;
    const onlineServers = mockServers.filter(s => s.status === 'Online').length;

    // Calculate total players across all servers
    const totalPlayers = mockServers.reduce((acc, server) => {
        return acc + (server.players?.online || 0);
    }, 0);
    const maxPlayers = mockServers.reduce((acc, server) => {
        return acc + (server.players?.max || 0);
    }, 0);

    // Calculate average CPU usage from online servers
    const onlineServersArray = mockServers.filter(s => s.status === 'Online');
    const avgCpuUsage = onlineServersArray.length
        ? Math.round(onlineServersArray.reduce((acc, server) => acc + server.cpuUsage, 0) / onlineServersArray.length)
        : 0;

    // Calculate total memory usage
    const totalMemoryUsed = mockServers.reduce((acc, server) => acc + server.memoryUsage.used, 0);
    const totalMemoryAllocated = mockServers.reduce((acc, server) => acc + server.memoryUsage.total, 0);

    // Convert bytes to GB for display
    const usedMemoryGB = (totalMemoryUsed / (1024 * 1024 * 1024)).toFixed(1);
    const totalMemoryGB = (totalMemoryAllocated / (1024 * 1024 * 1024)).toFixed(1);

    return (
        <div className="container py-6 px-4 mx-auto max-w-7xl">
            {/* Header with welcome message */}
            <div className="mb-6 no-select">
                <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">Dashboard Overview</h1>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-2">
                    Welcome back! Here's what's happening with your servers.
                </p>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
                <Card className="shadow-sm hover:shadow-md transition-shadow border-slate-200 dark:border-slate-700 no-select">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400">
                            Player Activity
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-between">
                            <div className="flex flex-col">
                                <span className="text-3xl font-bold">{totalPlayers}</span>
                                <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                                    across {onlineServers} online servers
                                </span>
                            </div>
                            <div
                                className="h-12 w-12 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center">
                                <Users className="h-6 w-6 text-emerald-600 dark:text-emerald-400"/>
                            </div>
                        </div>
                        <div className="mt-2 pt-2 border-t border-slate-200 dark:border-slate-700">
                            <span className="text-xs text-slate-500 dark:text-slate-400">
                                {totalPlayers}/{maxPlayers} total capacity
                            </span>
                        </div>
                    </CardContent>
                </Card>

                <Card className="shadow-sm hover:shadow-md transition-shadow border-slate-200 dark:border-slate-700 no-select">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400">
                            CPU Usage
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-between">
                            <ResourceGauge value={avgCpuUsage} maxValue={100} size="sm" unit="%" label="Average"/>
                            <div className="flex flex-col items-end">
                                <span className="text-3xl font-bold">{avgCpuUsage}%</span>
                                <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                                    Avg across {onlineServers} servers
                                </span>
                            </div>
                        </div>
                        <div className="mt-2 pt-2 border-t border-slate-200 dark:border-slate-700">
                            <span className="text-xs text-slate-500 dark:text-slate-400">
                                Peak: {Math.max(...onlineServersArray.map(s => s.cpuUsage), 0)}% on {onlineServersArray.length > 0 ? onlineServersArray.reduce((max, server) => server.cpuUsage > max.cpuUsage ? server : max).name : 'N/A'}
                            </span>
                        </div>
                    </CardContent>
                </Card>

                <Card className="shadow-sm hover:shadow-md transition-shadow border-slate-200 dark:border-slate-700 no-select">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-sm font-medium text-slate-500 dark:text-slate-400">
                            Memory Usage
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-between">
                            <ResourceGauge
                                value={parseFloat(usedMemoryGB)}
                                maxValue={parseFloat(totalMemoryGB)}
                                size="sm"
                                unit="GB"
                                label="Used"
                            />
                            <div className="flex flex-col items-end">
                                <span className="text-3xl font-bold">{usedMemoryGB} GB</span>
                                <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                                    of {totalMemoryGB} GB allocated
                                </span>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            <div className="space-y-4">
                <div className="flex items-center justify-between no-select">
                    <h2 className="text-xl font-bold text-slate-900 dark:text-slate-50">Servers</h2>
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
                    {mockServers.map((server) => (
                        <ServerCard
                            key={server.id}
                            server={server}
                            className="no-select"
                        />
                    ))}
                    <Card
                        className="flex flex-col items-center justify-center p-6 border-dashed border-2 border-slate-200 dark:border-slate-700 bg-transparent hover:bg-slate-50 dark:hover:bg-slate-800/20 transition-colors no-select">
                        <div className="rounded-full bg-slate-100 dark:bg-slate-800 p-3 mb-3">
                            <Plus className="h-6 w-6 text-slate-500 dark:text-slate-400"/>
                        </div>
                        <h3 className="text-lg font-medium text-slate-900 dark:text-slate-100">New Server</h3>
                        <p className="text-sm text-center text-slate-500 dark:text-slate-400 mt-1 mb-4">
                            Deploy a new server with your preferred configuration
                        </p>
                    </Card>
                </div>
            </div>
        </div>
    );
}

