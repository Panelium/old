import React from "react";
import { cn } from "~/lib/utils";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

interface EntityAvatarProps {
  title: string;
  subTitle: string;
  src?: string;
  alt?: string;
  className?: string;
}

const EntityAvatar: React.FC<EntityAvatarProps> = ({
  src,
  alt,
  title,
  subTitle,
  className,
}) => {
  return (
    <div className="flex items-center text-left gap-3">
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
            "text-foreground",
            "group-hover:text-indigo-600 dark:group-hover:text-indigo-400",
            className
          )}
        >
          {title}
        </h3>
        <p className="text-xs text-faded-foreground">{subTitle}</p>
      </div>
    </div>
  );
};

export default EntityAvatar;
