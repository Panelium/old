import * as React from "react";

import { cn } from "~/lib/utils";
import { Eye, EyeClosed } from "lucide-react";

function ToggleVisibility({updateVisible}: {updateVisible: (visible: boolean) => void}) {
  const [visible, setVisible] = React.useState(false);
  return (
    <div className="flex cursor-pointer w-9 h-9 items-center justify-center content-center ml-[calc(var(--spacing)*-9)]" onClick={() => {
      updateVisible(!visible);
      setVisible(!visible);
    }}>
      {visible ? <Eye className="w-5 h-5"/> : <EyeClosed className="w-5 h-5"/>}
    </div>
  );
}

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  const [visible, setVisible] = React.useState(false);

  return (
    <div className="flex">
    <input
      type={type === "password" ? (visible ? "text" : "password") : type}
      data-slot="input"
      className={cn(
        "file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground border-input flex h-9 w-full min-w-0 rounded-md border bg-white/5 px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
        "focus-visible:border-ring focus-visible:ring-indigo-500/70 focus-visible:ring-[3px]",
        "aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
        className
      )}
      {...props}
    />
    {type === "password" ? <ToggleVisibility updateVisible={(v) => {
      setVisible(v);
    }}/> : undefined}
    </div>
  );
}

export { Input };
