import {type LucideIcon} from "lucide-react";
import {Button} from "~/components/ui/button";
import React from "react";
import {cn} from "~/lib/utils";

interface IconTextProps {
    icon: LucideIcon;
    text: string;
    className?: string;
    onClick?: () => void;
}

export default function IconButton({icon, text, className, onClick}: IconTextProps) {
    const IconComponent = icon;
    return (
        <Button
            variant="outline"
            size="sm"
            className={cn("shadow-sm bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700", className)}
            onClick={onClick}
        >
            <IconComponent className="mr-2 h-4 w-4"/>
            {text}
        </Button>
    );
}