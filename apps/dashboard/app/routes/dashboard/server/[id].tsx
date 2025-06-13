import React, {useEffect, useState} from 'react';
import {Link, useParams, useSearchParams} from 'react-router-dom';
import {Card, CardContent, CardHeader, CardTitle} from '~/components/ui/card';
import {Button} from '~/components/ui/button';
import {Tabs, TabsContent, TabsList, TabsTrigger} from '~/components/ui/tabs';
import {ScrollArea} from '~/components/ui/scroll-area';
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from '~/components/ui/select';
import {Input} from '~/components/ui/input';
import {
    Activity,
    ArrowLeft,
    Clock,
    Cpu,
    Database,
    HardDrive,
    Power,
    RefreshCw,
    Server as ServerIcon,
    Settings,
    Terminal,
} from 'lucide-react';
import {cn} from '~/lib/utils';
import {FileManager} from '~/components/dashboard/FileManager';
import {ActivityLog} from '~/components/dashboard/ActivityLog';
import StatusBadge from "~/components/dashboard/StatusBadge";
import {ServerStatusType} from 'proto-gen-ts/daemon_pb';

// Placeholder data - would normally come from an API
const server = {
    id: '1',
    name: 'My Awesome Server 1',
    status: ServerStatusType.SERVER_STATUS_ONLINE,
    node: 'Node Alpha',
    cpuUsage: 25,
    memoryUsage: {
        used: 2048,
        total: 8192,
    },
    diskUsage: {
        used: 51200,
        total: 102400,
    },
    ip: '192.168.1.100',
    port: 25565,
    location: 'US East',
    game: 'Minecraft',
    uptime: '2d 15h 30m',
    console: [
        {time: '10:15:32', content: 'Server started on port 25565'},
        {time: '10:15:33', content: 'Loading world...'},
        {time: '10:15:40', content: 'World loaded in 7s'},
        {time: '10:15:41', content: 'Server ready for connections'},
        {time: '10:16:20', content: 'Player1 joined the game'},
        {time: '10:20:15', content: 'Player1: Hello everyone!'},
        {time: '10:25:30', content: 'Player2 joined the game'},
        {time: '10:30:45', content: 'Autosaving world...'},
        {time: '10:30:47', content: 'World autosaved'},
    ],
};

export default function ServerDetailsPage() {
    const {id: serverId} = useParams();
    const [searchParams] = useSearchParams();
    const tabParam = searchParams.get('tab');
    const [activeTab, setActiveTab] = useState(tabParam || 'console');
    const [command, setCommand] = useState('');

    // Log server ID for debugging - this would normally be used to fetch the server
    console.log(`Loading server: ${serverId}`);

    // Update active tab when URL parameter changes
    useEffect(() => {
        if (tabParam && ['console', 'files', 'activity', 'settings'].includes(tabParam)) {
            setActiveTab(tabParam);
        }
    }, [tabParam]);

    // Format memory and disk usage
    const formatMemory = (bytes: number): string => {
        const gb = bytes / 1024 / 1024;
        return gb >= 1 ? `${gb.toFixed(1)} GB` : `${(bytes / 1024).toFixed(0)} MB`;
    };

    const formatDisk = (bytes: number): string => {
        const gb = bytes / 1024 / 1024;
        return `${gb.toFixed(1)} GB`;
    };

    // Handle command submit
    const handleCommandSubmit = (e: React.FormEvent): void => {
        e.preventDefault();
        if (!command.trim()) return;

        // In a real app, you would send this command to the server
        console.log(`Sending command: ${command}`);
        setCommand('');
    };

    return (
        <div className="p-6 bg-slate-50 dark:bg-slate-900 min-h-full">
            <div className="max-w-7xl mx-auto space-y-8">
                {/* Server header */}
                <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between no-select">
                    <div className="flex items-center gap-3">
                        <Button variant="outline" size="icon"
                                className="rounded-lg shadow-sm bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700"
                                asChild>
                            <Link to="/">
                                <ArrowLeft className="h-4 w-4 text-slate-500 dark:text-slate-400"/>
                            </Link>
                        </Button>
                        <div>
                            <h1 className="text-2xl font-bold text-slate-900 dark:text-slate-50">{server.name}</h1>
                            <div className="flex items-center gap-3 mt-1">
                                <StatusBadge status={server.status}/>
                                <div className="flex items-center text-sm text-slate-500 dark:text-slate-400">
                                    <ServerIcon className="h-3.5 w-3.5 mr-1.5"/>
                                    {server.game}
                                </div>
                                <div className="flex items-center text-sm text-slate-500 dark:text-slate-400">
                                    <Clock className="h-3.5 w-3.5 mr-1.5"/>
                                    {server.uptime}
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            className="shadow-sm bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700"
                        >
                            <RefreshCw className="mr-2 h-4 w-4"/>
                            Restart
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            className="shadow-sm bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-700"
                        >
                            <Power className="mr-2 h-4 w-4"/>
                            Stop
                        </Button>
                    </div>
                </div>

                {/* Resource stats grid */}
                <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
                    {/* CPU Usage Card */}
                    <Card
                        className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
                        <CardHeader className="pb-2">
                            <CardTitle
                                className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
                                <Cpu className="mr-2 h-4 w-4 text-indigo-500"/>
                                CPU Usage
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between">
                                <div
                                    className="text-2xl font-bold text-slate-900 dark:text-slate-50">{server.cpuUsage}%
                                </div>
                                <div className="text-xs text-slate-500 dark:text-slate-400">Core Limit: 100%</div>
                            </div>
                            <div
                                className="mt-3 h-1.5 w-full bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
                                <div
                                    className={cn(
                                        'h-full rounded-full',
                                        server.cpuUsage > 90 ? 'bg-red-500' : server.cpuUsage > 75 ? 'bg-amber-500' : 'bg-emerald-500'
                                    )}
                                    style={{width: `${server.cpuUsage}%`}}
                                />
                            </div>
                        </CardContent>
                    </Card>

                    {/* Memory Usage Card */}
                    <Card
                        className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
                        <CardHeader className="pb-2">
                            <CardTitle
                                className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
                                <HardDrive className="mr-2 h-4 w-4 text-indigo-500"/>
                                Memory Usage
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between">
                                <div
                                    className="text-2xl font-bold text-slate-900 dark:text-slate-50">{formatMemory(server.memoryUsage.used)}</div>
                                <div
                                    className="text-xs text-slate-500 dark:text-slate-400">of {formatMemory(server.memoryUsage.total)}</div>
                            </div>
                            <div
                                className="mt-3 h-1.5 w-full bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
                                <div
                                    className={cn(
                                        'h-full rounded-full',
                                        server.memoryUsage.used / server.memoryUsage.total > 0.9
                                            ? 'bg-red-500'
                                            : server.memoryUsage.used / server.memoryUsage.total > 0.75
                                                ? 'bg-amber-500'
                                                : 'bg-emerald-500'
                                    )}
                                    style={{width: `${(server.memoryUsage.used / server.memoryUsage.total) * 100}%`}}
                                />
                            </div>
                        </CardContent>
                    </Card>

                    {/* Disk Usage Card */}
                    <Card
                        className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
                        <CardHeader className="pb-2">
                            <CardTitle
                                className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
                                <HardDrive className="mr-2 h-4 w-4 text-indigo-500"/>
                                Disk Usage
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between">
                                <div
                                    className="text-2xl font-bold text-slate-900 dark:text-slate-50">{formatDisk(server.diskUsage.used)}</div>
                                <div
                                    className="text-xs text-slate-500 dark:text-slate-400">of {formatDisk(server.diskUsage.total)}</div>
                            </div>
                            <div
                                className="mt-3 h-1.5 w-full bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
                                <div
                                    className={cn(
                                        'h-full rounded-full',
                                        server.diskUsage.used / server.diskUsage.total > 0.9
                                            ? 'bg-red-500'
                                            : server.diskUsage.used / server.diskUsage.total > 0.75
                                                ? 'bg-amber-500'
                                                : 'bg-emerald-500'
                                    )}
                                    style={{width: `${(server.diskUsage.used / server.diskUsage.total) * 100}%`}}
                                />
                            </div>
                        </CardContent>
                    </Card>

                    {/* Connection Card */}
                    <Card
                        className="border-slate-200/40 dark:border-slate-700/30 bg-white dark:bg-slate-800 shadow-sm hover:shadow transition-shadow rounded-xl no-select">
                        <CardHeader className="pb-2">
                            <CardTitle
                                className="text-sm font-medium text-slate-500 dark:text-slate-400 flex items-center">
                                <Database className="mr-2 h-4 w-4 text-indigo-500"/>
                                Connection
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex flex-col space-y-2">
                                <div className="flex items-center justify-between">
                                    <span className="text-xs text-slate-500 dark:text-slate-400">IP Address</span>
                                    <span
                                        className="text-xs font-medium text-slate-900 dark:text-slate-50">{server.ip}</span>
                                </div>
                                <div className="flex items-center justify-between">
                                    <span className="text-xs text-slate-500 dark:text-slate-400">Port</span>
                                    <span
                                        className="text-xs font-medium text-slate-900 dark:text-slate-50">{server.port}</span>
                                </div>
                                <div className="flex items-center justify-between">
                                    <span className="text-xs text-slate-500 dark:text-slate-400">Location</span>
                                    <span
                                        className="text-xs font-medium text-slate-900 dark:text-slate-50">{server.location}</span>
                                </div>
                                <div className="mt-2 pt-2 border-t border-slate-200 dark:border-slate-700">
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        className="w-full h-7 text-xs border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700"
                                    >
                                        Copy Details
                                    </Button>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Tabs section */}
                <Card
                    className="border-slate-200/60 dark:border-slate-700/40 bg-white dark:bg-slate-800 shadow-sm overflow-hidden rounded-xl py-0">
                    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                        <TabsList
                            className="grid grid-cols-4 rounded-none border-b border-slate-200 dark:border-slate-700 bg-transparent p-0 h-12 no-select">
                            <TabsTrigger
                                value="console"
                                className="relative rounded-none data-[state=active]:bg-transparent data-[state=active]:text-indigo-600 dark:data-[state=active]:text-indigo-400 data-[state=active]:shadow-none h-12 border-transparent border-0"
                            >
                                <Terminal className="mr-2 h-4 w-4"/>
                                Console
                                <span
                                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600 dark:bg-indigo-400 data-[state=active]:animate-tab-border-in data-[state=inactive]:w-0 transition-all duration-300"
                                    data-state={activeTab === 'console' ? 'active' : 'inactive'}></span>
                            </TabsTrigger>
                            <TabsTrigger
                                value="files"
                                className="relative rounded-none data-[state=active]:bg-transparent data-[state=active]:text-indigo-600 dark:data-[state=active]:text-indigo-400 data-[state=active]:shadow-none h-12 border-transparent border-0"
                            >
                                <HardDrive className="mr-2 h-4 w-4"/>
                                Files
                                <span
                                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600 dark:bg-indigo-400 data-[state=active]:animate-tab-border-in data-[state=inactive]:w-0 transition-all duration-300"
                                    data-state={activeTab === 'files' ? 'active' : 'inactive'}></span>
                            </TabsTrigger>
                            <TabsTrigger
                                value="activity"
                                className="relative rounded-none data-[state=active]:bg-transparent data-[state=active]:text-indigo-600 dark:data-[state=active]:text-indigo-400 data-[state=active]:shadow-none h-12 border-transparent border-0"
                            >
                                <Activity className="mr-2 h-4 w-4"/>
                                Activity
                                <span
                                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600 dark:bg-indigo-400 data-[state=active]:animate-tab-border-in data-[state=inactive]:w-0 transition-all duration-300"
                                    data-state={activeTab === 'activity' ? 'active' : 'inactive'}></span>
                            </TabsTrigger>
                            <TabsTrigger
                                value="settings"
                                className="relative rounded-none data-[state=active]:bg-transparent data-[state=active]:text-indigo-600 dark:data-[state=active]:text-indigo-400 data-[state=active]:shadow-none h-12 border-transparent border-0"
                            >
                                <Settings className="mr-2 h-4 w-4"/>
                                Settings
                                <span
                                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600 dark:bg-indigo-400 data-[state=active]:animate-tab-border-in data-[state=inactive]:w-0 transition-all duration-300"
                                    data-state={activeTab === 'settings' ? 'active' : 'inactive'}></span>
                            </TabsTrigger>
                        </TabsList>

                        <TabsContent value="console" className="m-0 px-6 pb-4">
                            <div className="flex flex-col h-full">
                                <div className="flex items-center justify-between mb-4 no-select">
                                    <div>
                                        <h3 className="text-lg font-medium text-slate-900 dark:text-slate-50">Server
                                            Console</h3>
                                        <p className="text-sm text-slate-500 dark:text-slate-400">View live server
                                            output and send commands</p>
                                    </div>
                                </div>

                                <div
                                    className="relative flex-1 border border-slate-200/40 dark:border-slate-700/30 rounded-lg overflow-hidden bg-slate-950">
                                    <ScrollArea
                                        className="h-[352px] w-full text-emerald-400 p-4 font-mono text-sm"
                                        style={{
                                            userSelect: 'none',
                                            WebkitTouchCallout: 'none'
                                        }}
                                        onFocus={(e) => {
                                            const element = e.currentTarget;
                                            element.style.userSelect = 'text';
                                        }}
                                        onBlur={(e) => {
                                            const element = e.currentTarget;
                                            element.style.userSelect = 'none';
                                        }}
                                        tabIndex={0} // Make the ScrollArea focusable
                                    >
                                        <div
                                            className="console-content"
                                            onKeyDown={(e) => {
                                                // Prevent Ctrl+A selection in the console
                                                if ((e.ctrlKey || e.metaKey) && e.key === 'a') {
                                                    e.preventDefault();
                                                }
                                            }}
                                            onSelectCapture={(e) => {
                                                // Check if the selection was triggered by Ctrl+A
                                                if (document.getSelection()?.toString() === document.querySelector('.console-content')?.textContent) {
                                                    // If entire content is selected, it's likely Ctrl+A
                                                    e.preventDefault();
                                                    // Clear selection
                                                    window.getSelection()?.removeAllRanges();
                                                }
                                            }}
                                        >
                                            <div className="pb-1 text-xs text-slate-500 no-select">--- Server started
                                                on {new Date().toLocaleString()} ---
                                            </div>
                                            {server.console.map((line, i) => {
                                                // Apply styling based on content type
                                                const isError = line.content.toLowerCase().includes('error') || line.content.toLowerCase().includes('exception');
                                                const isWarning = line.content.toLowerCase().includes('warn');
                                                const isInfo = line.content.toLowerCase().includes('info');

                                                return (
                                                    <div
                                                        key={i}
                                                        className={cn(
                                                            'pb-1',
                                                            isError && 'text-red-400',
                                                            isWarning && 'text-amber-400',
                                                            isInfo && 'text-blue-400'
                                                        )}
                                                    >
                                                        <span
                                                            className="text-slate-500">[{line.time}]</span> {line.content}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </ScrollArea>

                                    <form onSubmit={handleCommandSubmit}
                                          className="absolute bottom-0 left-0 right-0 border-t border-slate-700/50 bg-slate-950 px-4 py-2 no-select">
                                        <div className="flex items-center">
                                            <span className="text-slate-500 mr-2">$</span>
                                            <input
                                                type="text"
                                                value={command}
                                                onChange={(e) => setCommand(e.target.value)}
                                                placeholder="Type a command and press Enter..."
                                                className="flex-1 bg-transparent border-none outline-none text-emerald-400 font-mono text-sm focus:outline-none focus:ring-0"
                                                autoFocus
                                                onKeyDown={(e) => {
                                                    if (e.key === 'Enter') {
                                                        e.preventDefault();
                                                        handleCommandSubmit(e);
                                                    }
                                                }}
                                            />
                                        </div>
                                    </form>
                                </div>

                                <div className="mt-2 text-xs text-slate-500 dark:text-slate-400 no-select">
                                    <span className="font-semibold">Tip:</span> Type{' '}
                                    <code
                                        className="px-1 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-indigo-700 dark:text-indigo-400">help</code> to
                                    see available commands
                                </div>
                            </div>
                        </TabsContent>

                        <TabsContent value="files" className="m-0 p-6 py-4">
                            <div className="no-select">
                                <FileManager serverId={serverId!}/>
                            </div>
                        </TabsContent>

                        <TabsContent value="activity" className="m-0 p-6 py-4">
                            <div className="no-select">
                                <ActivityLog activities={[]}/>
                            </div>
                        </TabsContent>

                        <TabsContent value="settings" className="m-0 px-6 pb-4">
                            <div className="space-y-6 no-select">
                                <div className="no-select">
                                    <h3 className="text-lg font-medium text-slate-900 dark:text-slate-50">Server
                                        Settings</h3>
                                    <p className="text-sm text-slate-500 dark:text-slate-400">Configure your server
                                        properties and behavior</p>
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    <Card
                                        className="border-slate-200/60 dark:border-slate-700/40 bg-white dark:bg-slate-800 shadow-sm rounded-xl py-0 no-select">
                                        <CardHeader className="p-4">
                                            <CardTitle className="text-base text-slate-900 dark:text-slate-50">General
                                                Settings</CardTitle>
                                        </CardHeader>
                                        <CardContent className="p-4 pt-0 space-y-4">
                                            <div className="grid gap-2">
                                                <label
                                                    className="text-sm font-medium text-slate-700 dark:text-slate-300">Server
                                                    Name</label>
                                                <Input
                                                    type="text"
                                                    defaultValue={server.name}
                                                    className="border-slate-200/60 dark:border-slate-700/40 focus-visible:ring-indigo-500/70"
                                                />
                                            </div>

                                            <div className="grid gap-2">
                                                <label
                                                    className="text-sm font-medium text-slate-700 dark:text-slate-300">Memory
                                                    Limit</label>
                                                <Select defaultValue="2 GB">
                                                    <SelectTrigger
                                                        className="w-full border-slate-200/60 dark:border-slate-700/40 focus-visible:ring-indigo-500/70">
                                                        <SelectValue placeholder="Select memory limit"/>
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value="2 GB"
                                                                    className="!bg-white dark:!bg-slate-800">2
                                                            GB</SelectItem>
                                                        <SelectItem value="4 GB"
                                                                    className="!bg-white dark:!bg-slate-800">4
                                                            GB</SelectItem>
                                                        <SelectItem value="8 GB"
                                                                    className="!bg-white dark:!bg-slate-800">8
                                                            GB</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </div>

                                            <div className="grid gap-2">
                                                <label
                                                    className="text-sm font-medium text-slate-700 dark:text-slate-300">CPU
                                                    Limit</label>
                                                <Select defaultValue="1 Core">
                                                    <SelectTrigger
                                                        className="w-full border-slate-200/60 dark:border-slate-700/40 focus-visible:ring-indigo-500/70">
                                                        <SelectValue placeholder="Select CPU limit"/>
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value="1 Core"
                                                                    className="!bg-white dark:!bg-slate-800">1
                                                            Core</SelectItem>
                                                        <SelectItem value="2 Cores"
                                                                    className="!bg-white dark:!bg-slate-800">2
                                                            Cores</SelectItem>
                                                        <SelectItem value="4 Cores"
                                                                    className="!bg-white dark:!bg-slate-800">4
                                                            Cores</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </div>
                                        </CardContent>
                                    </Card>

                                    <Card
                                        className="border-slate-200/60 dark:border-slate-700/40 bg-white dark:bg-slate-800 shadow-sm rounded-xl py-0 no-select">
                                        <CardHeader className="p-4">
                                            <CardTitle className="text-base text-slate-900 dark:text-slate-50">Network
                                                Settings</CardTitle>
                                        </CardHeader>
                                        <CardContent className="p-4 pt-0 space-y-4">
                                            <div className="grid gap-2">
                                                <label
                                                    className="text-sm font-medium text-slate-700 dark:text-slate-300">Server
                                                    Port</label>
                                                <Input
                                                    type="text"
                                                    defaultValue={server.port}
                                                    className="border-slate-200/60 dark:border-slate-700/40 focus-visible:ring-indigo-500/70"
                                                />
                                            </div>

                                            <div className="grid gap-2">
                                                <label
                                                    className="text-sm font-medium text-slate-700 dark:text-slate-300">Connection
                                                    Address</label>
                                                <Input
                                                    type="text"
                                                    defaultValue={server.ip}
                                                    className="bg-slate-50 dark:bg-slate-700 text-slate-500 dark:text-slate-400 cursor-not-allowed border-slate-200/40 dark:border-slate-700/30"
                                                    readOnly
                                                />
                                            </div>

                                            <div className="grid gap-2">
                                                <label
                                                    className="text-sm font-medium text-slate-700 dark:text-slate-300">Location</label>
                                                <Select defaultValue={server.location}>
                                                    <SelectTrigger
                                                        className="w-full border-slate-200/60 dark:border-slate-700/40 focus-visible:ring-indigo-500/70">
                                                        <SelectValue placeholder="Select location"/>
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value="US East"
                                                                    className="!bg-white dark:!bg-slate-800">US
                                                            East</SelectItem>
                                                        <SelectItem value="US West"
                                                                    className="!bg-white dark:!bg-slate-800">US
                                                            West</SelectItem>
                                                        <SelectItem value="EU Central"
                                                                    className="!bg-white dark:!bg-slate-800">EU
                                                            Central</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </div>
                                        </CardContent>
                                    </Card>
                                </div>

                                <div className="flex justify-end gap-2 mt-6">
                                    <Button
                                        variant="outline"
                                        className="shadow-sm bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700"
                                    >
                                        Cancel
                                    </Button>
                                    <Button className="bg-indigo-600 hover:bg-indigo-700 text-white shadow-sm">Save
                                        Changes</Button>
                                </div>
                            </div>
                        </TabsContent>
                    </Tabs>
                </Card>
            </div>
        </div>
    );
}
