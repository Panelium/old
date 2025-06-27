"use client";

import * as React from "react";
import * as AvatarPrimitive from "@radix-ui/react-avatar";

import {cn} from "~/lib/utils";

function Avatar({
                    className,
                    ...props
                }: React.ComponentProps<typeof AvatarPrimitive.Root>) {
    return (
        <AvatarPrimitive.Root
            data-slot="avatar"
            className={cn(
                "flex h-10 w-10 items-center justify-center rounded-md bg-tag-purple-background flex-shrink-0",
                className
            )}
            {...props}
        />
    );
}

function AvatarImage({
                         className,
                         ...props
                     }: React.ComponentProps<typeof AvatarPrimitive.Image>) {
    return (
        <AvatarPrimitive.Image
            data-slot="avatar-image"
            className={cn("aspect-square size-full", className)}
            {...props}
        />
    );
}

function AvatarFallback({
                            className,
                            ...props
                        }: React.ComponentProps<typeof AvatarPrimitive.Fallback>) {
    return (
        <AvatarPrimitive.Fallback
            data-slot="avatar-fallback"
            className={cn(
                "text-lg font-bold text-tag-purple",
                className
            )}
            {...props}
        />
    );
}

export {Avatar, AvatarImage, AvatarFallback};
