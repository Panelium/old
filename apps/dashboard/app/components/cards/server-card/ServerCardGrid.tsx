import React from "react";
import {Plus} from "lucide-react";
import {cn} from "~/lib/utils";

import type {Server} from "./ServerCard";

import {Card} from "~/components/ui/card";
import ServerCard from "~/components/cards/server-card/ServerCard";

const AddServerCard: React.FC = () => {
    return (
        <Card
            className={cn(
                "flex flex-col items-center justify-center no-select",
                "border-border bg-server-card/50 hover:border-border-hover/50 hover:bg-server-card",
                "transition-all cursor-pointer backdrop-blur-[3px]",
            )}
        >
            <div className="rounded-full bg-tag-gray-background p-3 mb-3">
                <Plus className="h-6 w-6 text-tag-gray"/>
            </div>
            <h3 className="text-lg font-medium text-foreground">
                New Server
            </h3>
            <p className="text-sm text-center text-server-card-foreground mt-1 mb-4">
                Deploy a new server with your preferred configuration
            </p>
        </Card>
    );
};
export default function ServerCardGrid({servers}: { servers: Server[] }) {
    return (
        <div className="flex flex-col gap-4">
            <h2 className="text-xl font-bold text-foreground no-select">
                Servers
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
                {servers.map((server) => (
                    <ServerCard key={server.id} server={server} className="no-select"/>
                ))}
                <AddServerCard/>
            </div>
        </div>
    );
}
