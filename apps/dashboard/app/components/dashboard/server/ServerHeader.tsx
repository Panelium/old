import React from "react";
import { Link } from "react-router";
import { ArrowLeft } from "lucide-react";

import { PowerAction } from "proto-gen-ts/daemon_pb";

import StatusBadge from "../StatusBadge";
import { Button } from "~/components/ui/button";
import SoftwareText from "~/components/texts/SoftwareText";
import DurationText from "~/components/texts/DurationText";
import PowerButton from "~/components/buttons/PowerButton";
import type { Server } from "~/components/cards/server-card";

const ServerHeader: React.FC<{ server: Server }> = ({ server }) => {
  return (
    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between no-select">
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          size="icon"
          className="rounded-lg shadow-sm bg-white dark:bg-slate-800 border-border"
          asChild
        >
          <Link to="/">
            <ArrowLeft className="h-4 w-4 text-slate-500 dark:text-slate-400" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-slate-900 dark:text-slate-50">
            {server.name}
          </h1>
          <div className="flex items-center gap-3 mt-1">
            <StatusBadge status={server.status} />
            <SoftwareText software={"Minecraft"} /> {/* TODO: CHANGE ME */}
            <DurationText startDate={new Date(1747096311000)} />{" "}
            {/* TODO: CHANGE ME */}
          </div>
        </div>
      </div>
      <div className="flex gap-2">
        <PowerButton action={PowerAction.RESTART} />
        <PowerButton action={PowerAction.STOP} />
        <PowerButton action={PowerAction.KILL} />
      </div>
    </div>
  );
};

export default ServerHeader;
