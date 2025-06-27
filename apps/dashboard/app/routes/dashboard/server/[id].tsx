import React from "react";
import useServer from "./useServer";
import TabsSection from "~/components/dashboard/server/TabsSection";
import ServerHeader from "~/components/dashboard/server/ServerHeader";
import ResourceStats from "~/components/dashboard/server/ResourceStats";

export default function ServerDetailsPage() {
    const {server} = useServer();

    return (
        <div className="p-6 bg-background min-h-full">
            <div className="max-w-7xl mx-auto space-y-8">
                <ServerHeader server={server}/>
                <ResourceStats server={server}/>
                <TabsSection/>
            </div>
        </div>
    );
}
