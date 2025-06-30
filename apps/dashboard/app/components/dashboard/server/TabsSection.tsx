import React, { type ReactElement } from "react";
import {Activity, HardDrive, Settings, Terminal} from "lucide-react";
import {cn} from "~/lib/utils";

import useServer from "~/routes/dashboard/server/useServer";

import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue,} from "~/components/ui/select";
import {FileManager} from "../FileManager";
import {ActivityLog} from "../ActivityLog";
import {Input} from "~/components/ui/input";
import {Button} from "~/components/ui/button";
import {ScrollArea} from "~/components/ui/scroll-area";
import {Card, CardContent, CardHeader, CardTitle} from "~/components/ui/card";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "~/components/ui/tabs";

interface SectionTabsTriggerProps {
    value: string,
    activeTab: string,
    children: any,
}

function SectionTabsTrigger({value, activeTab, children}: SectionTabsTriggerProps){
    return (<>
        <TabsTrigger
            value={value}
            className="relative rounded-none data-[state=active]:bg-tag-purple/5 data-[state=active]:text-tag-purple data-[state=active]:shadow-none h-12 border-transparent border-0"
        >
            {children}
            <span className="capitalize">{value}</span>
            <span
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600 dark:bg-indigo-400 data-[state=active]:animate-tab-border-in data-[state=inactive]:w-0 transition-all duration-300"
                data-state={activeTab === value ? "active" : "inactive"}
            ></span>
        </TabsTrigger>
    </>);
}

const TabsSection: React.FC = () => {
    const {
        activeTab,
        setActiveTab,
        server,
        handleCommandSubmit,
        command,
        setCommand,
    } = useServer();
    return (
        <Card className="border-border bg-server-card shadow-sm overflow-hidden rounded-xl py-0">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                <TabsList
                    className="grid grid-cols-4 rounded-none border-b border-border bg-transparent p-0 h-12 no-select">
                    <SectionTabsTrigger value="console"  activeTab={activeTab}><Terminal  className="mr-2 h-4 w-4"/></SectionTabsTrigger>
                    <SectionTabsTrigger value="files"    activeTab={activeTab}><HardDrive className="mr-2 h-4 w-4"/></SectionTabsTrigger>
                    <SectionTabsTrigger value="activity" activeTab={activeTab}><Activity  className="mr-2 h-4 w-4"/></SectionTabsTrigger>
                    <SectionTabsTrigger value="settings" activeTab={activeTab}><Settings  className="mr-2 h-4 w-4"/></SectionTabsTrigger>
                </TabsList>

                <TabsContent value="console" className="m-0 px-6 pb-4">
                    <div className="flex flex-col h-full">
                        <div className="flex items-center justify-between mb-4 no-select">
                            <div>
                                <h3 className="text-lg font-medium text-slate-900 dark:text-slate-50">
                                    Server Console
                                </h3>
                                <p className="text-sm text-slate-500 dark:text-slate-400">
                                    View live server output and send commands
                                </p>
                            </div>
                        </div>

                        <div
                            className="relative flex-1 border border-slate-200/40 dark:border-slate-700/30 rounded-lg overflow-hidden bg-slate-950">
                            <ScrollArea
                                className="h-[352px] w-full text-emerald-400 p-4 font-mono text-sm"
                                style={{
                                    userSelect: "none",
                                    WebkitTouchCallout: "none",
                                }}
                                onFocus={(e) => {
                                    const element = e.currentTarget;
                                    element.style.userSelect = "text";
                                }}
                                onBlur={(e) => {
                                    const element = e.currentTarget;
                                    element.style.userSelect = "none";
                                }}
                                tabIndex={0} // Make the ScrollArea focusable
                            >
                                <div
                                    className="console-content"
                                    onKeyDown={(e) => {
                                        // Prevent Ctrl+A selection in the console
                                        if ((e.ctrlKey || e.metaKey) && e.key === "a") {
                                            e.preventDefault();
                                        }
                                    }}
                                    onSelectCapture={(e) => {
                                        // Check if the selection was triggered by Ctrl+A
                                        if (
                                            document.getSelection()?.toString() ===
                                            document.querySelector(".console-content")?.textContent
                                        ) {
                                            // If entire content is selected, it's likely Ctrl+A
                                            e.preventDefault();
                                            // Clear selection
                                            window.getSelection()?.removeAllRanges();
                                        }
                                    }}
                                >
                                    <div className="pb-1 text-xs text-slate-500 no-select">
                                        --- Server started on {new Date().toLocaleString()} ---
                                    </div>
                                    {server.console.map((line, i) => {
                                        // Apply styling based on content type
                                        const isError =
                                            line.content.toLowerCase().includes("error") ||
                                            line.content.toLowerCase().includes("exception");
                                        const isWarning = line.content
                                            .toLowerCase()
                                            .includes("warn");
                                        const isInfo = line.content.toLowerCase().includes("info");

                                        return (
                                            <div
                                                key={i}
                                                className={cn(
                                                    "pb-1",
                                                    isError && "text-red-400",
                                                    isWarning && "text-amber-400",
                                                    isInfo && "text-blue-400"
                                                )}
                                            >
                                                <span className="text-slate-500">[{line.time}]</span>{" "}
                                                {line.content}
                                            </div>
                                        );
                                    })}
                                </div>
                            </ScrollArea>

                            <form
                                onSubmit={handleCommandSubmit}
                                className="absolute bottom-0 left-0 right-0 border-t border-slate-700/50 bg-slate-950 px-4 py-2 no-select"
                            >
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
                                            if (e.key === "Enter") {
                                                e.preventDefault();
                                                handleCommandSubmit(e);
                                            }
                                        }}
                                    />
                                </div>
                            </form>
                        </div>

                        <div className="mt-2 text-xs text-slate-500 dark:text-slate-400 no-select">
                            <span className="font-semibold">Tip:</span> Type{" "}
                            <code
                                className="px-1 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-indigo-700 dark:text-indigo-400">
                                help
                            </code>{" "}
                            to see available commands
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="files" className="m-0 p-6 py-4">
                    <div className="no-select">
                        <FileManager serverId={server.id!}/>
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
                            <h3 className="text-lg font-medium text-foreground">
                                Server Settings
                            </h3>
                            <p className="text-sm text-card-muted-foreground">
                                Configure your server properties and behavior
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <Card className="border-border bg-card shadow-sm rounded-xl py-0 no-select">
                                <CardHeader className="p-4">
                                    <CardTitle className="text-base text-foreground">
                                        General Settings
                                    </CardTitle>
                                </CardHeader>
                                <CardContent className="p-4 pt-0 space-y-4">
                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium text-card-foreground">
                                            Server Name
                                        </label>
                                        <Input
                                            type="text"
                                            defaultValue={server.name}
                                            className="border-border focus-visible:ring-indigo-500/70"
                                        />
                                    </div>

                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium text-card-foreground">
                                            Memory Limit
                                        </label>
                                        <Select defaultValue="2 GB">
                                            <SelectTrigger
                                                className="w-full border-border focus-visible:ring-indigo-500/70">
                                                <SelectValue placeholder="Select memory limit"/>
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem
                                                    value="2 GB"
                                                >
                                                    2 GB
                                                </SelectItem>
                                                <SelectItem
                                                    value="4 GB"
                                                >
                                                    4 GB
                                                </SelectItem>
                                                <SelectItem
                                                    value="8 GB"
                                                >
                                                    8 GB
                                                </SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>

                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium text-card-foreground">
                                            CPU Limit
                                        </label>
                                        <Select defaultValue="1 Core">
                                            <SelectTrigger
                                                className="w-full border-border focus-visible:ring-indigo-500/70">
                                                <SelectValue placeholder="Select CPU limit"/>
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem
                                                    value="1 Core"
                                                >
                                                    1 Core
                                                </SelectItem>
                                                <SelectItem
                                                    value="2 Cores"
                                                >
                                                    2 Cores
                                                </SelectItem>
                                                <SelectItem
                                                    value="4 Cores"
                                                >
                                                    4 Cores
                                                </SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card className="border-border bg-card shadow-sm rounded-xl py-0 no-select">
                                <CardHeader className="p-4">
                                    <CardTitle className="text-base text-foreground">
                                        Network Settings
                                    </CardTitle>
                                </CardHeader>
                                <CardContent className="p-4 pt-0 space-y-4">
                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium text-card-foreground">
                                            Server Port
                                        </label>
                                        <Input
                                            type="text"
                                            defaultValue={server.port}
                                            className="border-border focus-visible:ring-indigo-500/70"
                                        />
                                    </div>

                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                                            Connection Address
                                        </label>
                                        <Input
                                            type="text"
                                            defaultValue={server.ip}
                                            className="bg-white/5 text-muted-foreground cursor-not-allowed"
                                            readOnly
                                        />
                                    </div>

                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                                            Location
                                        </label>
                                        <Select defaultValue={server.location}>
                                            <SelectTrigger
                                                className="w-full border-border focus-visible:ring-indigo-500/70">
                                                <SelectValue placeholder="Select location"/>
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem
                                                    value="US East"
                                                >
                                                    US East
                                                </SelectItem>
                                                <SelectItem
                                                    value="US West"
                                                >
                                                    US West
                                                </SelectItem>
                                                <SelectItem
                                                    value="EU Central"
                                                >
                                                    EU Central
                                                </SelectItem>
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
                            <Button className="bg-indigo-600 hover:bg-indigo-700 text-white shadow-sm">
                                Save Changes
                            </Button>
                        </div>
                    </div>
                </TabsContent>
            </Tabs>
        </Card>
    );
};

export default TabsSection;
