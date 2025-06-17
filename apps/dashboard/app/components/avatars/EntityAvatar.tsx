import React from "react";
import { cn } from "~/lib/utils";
import type { Server } from "../cards/server-card";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

const EntityAvatar: React.FC<{
  src?: string;
  alt?: string;
  title: string;
  subTitle: string;
  className: string;
}> = ({ src, alt, title, subTitle, className }) => {
  return (
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarImage
          src={src}
          alt={alt}
          className="h-full w-full object-cover"
        />
        <AvatarFallback>{title.charAt(0).toUpperCase()}</AvatarFallback>
      </Avatar>
      <div className="flex flex-col flex-1">
        <h3
          className={cn(
            "font-medium",
            "text-slate-900 dark:text-slate-100",
            "group-hover:text-indigo-600 dark:group-hover:text-indigo-400",
            className
          )}
        >
          {title}
        </h3>
        <p className="text-xs text-slate-500 dark:text-slate-400">{subTitle}</p>
      </div>
    </div>
  );
};

export default EntityAvatar;
