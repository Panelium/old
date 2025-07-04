import { Activity, HardDrive, Settings, Terminal, type LucideProps } from "lucide-react";
import type React from "react"
import { useState } from "react";
import { cn } from "~/lib/utils";
import ConsolePage from "./pages/ConsolePage";
import ActivityPage from "./pages/ActivityPage";
import FilesPage from "./pages/FilesPage";
import SettingsPage from "./pages/SettingsPage";

export const pagesEventBus = new Comment("event-bus");

export class PagePressedEvent extends Event {
    #id;

    constructor(id: string) {
        super("pagePressed");
        this.#id = id;
    }

    get id() {
        return this.#id;
    }
}

const transition = "transition-all"

interface PageProps {
    icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref">>,
    id: string,
    currentId: string,
    color?: string,
}

const Page: React.FC<PageProps> = ({icon, id, currentId, color}) => {
    const IconComponent = icon;

    // Weird workaround to get tailwind generating
    let bg70 = "bg-tag-gray/70";
    let bg20 = "bg-tag-gray/20";
    let iconColor = "group-hover/page:text-tag-gray text-tag-gray/70";
    
    if (color === "purple") { bg70 = "bg-tag-purple/70"; bg20 = "bg-tag-purple/20"; iconColor = "group-hover/page:text-tag-purple text-tag-purple/70"; }
    else if (color === "red") { bg70 = "bg-tag-red/70"; bg20 = "bg-tag-red/20"; iconColor = "group-hover/page:text-tag-red text-tag-red/70"; }
    else if (color === "green") { bg70 = "bg-tag-green/70"; bg20 = "bg-tag-green/20"; iconColor = "group-hover/page:text-tag-green text-tag-green/70"; }
    else if (color === "orange") { bg70 = "bg-tag-orange/70"; bg20 = "bg-tag-orange/20"; iconColor = "group-hover/page:text-tag-orange text-tag-orange/70"; }

    return (
        <button className={cn(
            "flex items-center h-14 w-14 cursor-pointer group/page",
            currentId === id ? bg20 : ""
        )} onClick={() => {
            pagesEventBus.dispatchEvent(new PagePressedEvent(id));
        }}>
            <div className={cn(
                "h-full w-1 mr-[calc(var(--spacing)*-1)]",
                bg70, 
                "transition-all duration-100 ease-linear", 
                currentId === id ? "scale-x-100" : "scale-x-0"
            )}/>
            <IconComponent className={cn(
                "h-7 w-7 m-auto",
                transition,
                currentId === id ? iconColor : "group-hover/page:text-foreground text-foreground/70"
            )}/>
        </button>
    );
}

const Pages: React.FC = () => {
    const [id, setId] = useState(location.pathname.split("/").length > 6 ? location.pathname.split("/")[3] : FilesPage.id);

    const FilesPageComponent = FilesPage.component(id);
    const ConsolePageComponent = ConsolePage.component(id);
    const ActivityPageComponent = ActivityPage.component(id);
    const SettingsPageComponent = SettingsPage.component(id);

    pagesEventBus.addEventListener("pagePressed", (e) => {
        if (!(e instanceof PagePressedEvent)) { return; }
        setId(e.id);
    });

    return (
        <div className="flex felx-row bg-card rounded-xl border-border shadow-lg shadow-black/20">
        <div className="flex flex-col items-center w-14 group/holder bg-accent/20 rounded-xl overflow-hidden shadow-md shadow-black/10">
            <Page icon={HardDrive} id={FilesPage.id}    currentId={id} color="orange"/>
            <Page icon={Terminal}  id={ConsolePage.id}  currentId={id} color="purple"/>
            <Page icon={Activity}  id={ActivityPage.id} currentId={id} color="green"/>
            <div className="mt-auto">
                <Page icon={Settings}  id={SettingsPage.id} currentId={id} color="gray"/>
            </div>
        </div>
        <FilesPageComponent/>
        <ConsolePageComponent/>
        <ActivityPageComponent/>
        <SettingsPageComponent/>
        </div>
    );
}

export default Pages