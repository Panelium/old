import { type LucideProps, Terminal } from "lucide-react";
import type React from "react";
import { cn } from "~/lib/utils";
import ConsolePage from "./pages/ConsolePage";
import { ServerInfo } from "proto-gen-ts/backend/Client_pb";
import { useParams } from "react-router";
import { useNavigate } from "react-router-dom";

const transition = "transition-all";

interface PageProps {
  icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref">>;
  serverId: string;
  id: string;
  currentId: string;
  color?: string;
}

const Page: React.FC<PageProps> = ({ icon, serverId, id, currentId, color }) => {
  const navigate = useNavigate();

  const IconComponent = icon;

  // Weird workaround to get tailwind generating
  let bg70 = "bg-tag-gray/70";
  let bg20 = "bg-tag-gray/20";
  let iconColor = "group-hover/page:text-tag-gray text-tag-gray/70";

  if (color === "purple") {
    bg70 = "bg-tag-purple/70";
    bg20 = "bg-tag-purple/20";
    iconColor = "group-hover/page:text-tag-purple text-tag-purple/70";
  } else if (color === "red") {
    bg70 = "bg-tag-red/70";
    bg20 = "bg-tag-red/20";
    iconColor = "group-hover/page:text-tag-red text-tag-red/70";
  } else if (color === "green") {
    bg70 = "bg-tag-green/70";
    bg20 = "bg-tag-green/20";
    iconColor = "group-hover/page:text-tag-green text-tag-green/70";
  } else if (color === "orange") {
    bg70 = "bg-tag-orange/70";
    bg20 = "bg-tag-orange/20";
    iconColor = "group-hover/page:text-tag-orange text-tag-orange/70";
  }

  return (
    <button
      className={cn("flex items-center h-14 w-14 cursor-pointer group/page", currentId === id ? bg20 : "")}
      onClick={() => {
        if (currentId !== id) {
          navigate(`/dashboard/server/${serverId}/${id}`);
        }
      }}
    >
      <div
        className={cn(
          "h-full w-1 mr-[calc(var(--spacing)*-1)]",
          bg70,
          "transition-all duration-100 ease-linear",
          currentId === id ? "scale-x-100" : "scale-x-0"
        )}
      />
      <IconComponent
        className={cn(
          "h-7 w-7 m-auto",
          transition,
          currentId === id ? iconColor : "group-hover/page:text-foreground text-foreground/70"
        )}
      />
    </button>
  );
};

export default function Pages({ server }: { server: ServerInfo }) {
  const params = useParams<{ id: string; page: string }>() ?? { page: "console" };
  const { id, page }: { id: string; page: string } = params as any;

  // const FilesPageComponent = FilesPage.component(page);
  const ConsolePageComponent = ConsolePage.component(page);
  // const ActivityPageComponent = ActivityPage.component(page);
  // const SettingsPageComponent = SettingsPage.component(page);

  return (
    <div className="flex felx-row bg-server-card border-border shadow-lg shadow-black/20 rounded-xl overflow-hidden">
      <div className="flex flex-col items-center w-14 group/holder bg-card overflow-hidden shadow-md shadow-black/10">
        {/*<Page icon={HardDrive} serverId={id} id={FilesPage.id} currentId={page} color="orange" />*/}
        <Page icon={Terminal} serverId={id} id={ConsolePage.id} currentId={page} color="purple" />
        {/*<Page icon={Activity} serverId={id} id={ActivityPage.id} currentId={page} color="green" />*/}
        {/*<div className="mt-auto">*/}
        {/*  <Page icon={Settings} serverId={id} id={SettingsPage.id} currentId={page} color="gray" />*/}
        {/*</div>*/}
      </div>
      {/*<FilesPageComponent />*/}
      <ConsolePageComponent />
      {/*<ActivityPageComponent />*/}
      {/*<SettingsPageComponent />*/}
    </div>
  );
}
