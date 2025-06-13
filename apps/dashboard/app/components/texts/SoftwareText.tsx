import {Server as ServerIcon} from "lucide-react";
import React from "react";
import IconText from "~/components/texts/IconText";

export default function SoftwareText({software}: { software: string }) {
    return (
        <IconText icon={ServerIcon} text={software}/>
    )
}