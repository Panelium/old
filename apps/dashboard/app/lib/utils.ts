import {twMerge} from "tailwind-merge";
import {type ClassValue, clsx} from "clsx";

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export const formatMemory = (bytes: number) => {
    const mb = bytes / 1024;
    if (mb >= 1024) {
        return `${(mb / 1024).toFixed(1)} GB`;
    }
    return `${Math.round(mb)} MB`;
};

export const formatDisk = (bytes: number): string => {
    const gb = bytes / 1024 / 1024;
    return `${gb.toFixed(1)} GB`;
};

export const getPercentage = (value: number, max: number) =>
    clampNumber((value / max) * 100, {min: 0, max: 100});

export const clampNumber = (
    n: number,
    {min, max}: { min?: number; max?: number }
) => {
    let c = n;
    if (max) {
        c = Math.min(c, max);
    }
    if (min) {
        c = Math.max(c, min);
    }
    return c;
};
