import { useEffect, useRef, useState } from "react";
import Page from "./Page";
import { ScrollArea } from "~/components/ui/scroll-area";
import { cn } from "~/lib/utils";
import { ServerService } from "proto-gen-ts/daemon/Server_pb";
import { Client } from "@connectrpc/connect";
import { useParams } from "react-router";
import { getClientClient, getDaemonServerClient } from "~/lib/api-clients";
import { ClientService, ServerInfo } from "proto-gen-ts/backend/Client_pb";

const ConsolePage: Page = new Page("console", () => {
  const params = useParams<{ id: string }>();
  const { id }: { id: string } = params as any;

  const [command, setCommand] = useState("");
  const [clientClient, setClientClient] = useState<Client<typeof ClientService>>();
  const [serverInfo, setServerInfo] = useState<ServerInfo>();
  const [serverClient, setServerClient] = useState<Client<typeof ServerService>>();
  const [consoleLines, setConsoleLines] = useState<string[]>([]);

  const scrollAreaRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (consoleLines[consoleLines.length - 1].includes("DOWNLOAD FINISHED")) {
      setTimeout(() => {
        window.location.reload();
      }, 3000);
      return;
    }
    // Auto-scroll to bottom when new lines are added
    if (scrollAreaRef.current) {
      scrollAreaRef.current.scrollTop = scrollAreaRef.current.scrollHeight;
    }
  }, [consoleLines]);

  useEffect(() => {
    (async () => {
      const clientClient = await getClientClient();
      setClientClient(clientClient);
    })();
  }, []);

  useEffect(() => {
    if (!clientClient) return;

    (async () => {
      const serverInfoResponse = await clientClient.getServer({ id });
      setServerInfo(serverInfoResponse);

      const serverClient = await getDaemonServerClient(serverInfoResponse.daemonHost);
      setServerClient(serverClient);
    })();
  }, [clientClient, id]);

  useEffect(() => {
    if (!serverClient) return;

    (async () => {
      const stream = serverClient.console({ id });

      for await (const message of stream) {
        setConsoleLines((prev) => [...prev, message.text]);
      }
    })();
  }, [serverClient]);

  const handleCommandSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!serverClient || !command.trim()) {
      setCommand(""); // Clear input even if command is empty or client is missing
      return;
    }

    try {
      await serverClient.consoleCommand({ id, text: command.trim() });
      setCommand("");
    } catch (error) {
      setCommand(""); // Clear input on error as well
      console.error("Error sending command:", error);
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between mb-4 no-select">
        <div>
          <h3 className="text-lg font-medium text-card-foreground">Server Console</h3>
          <p className="text-sm text-card-muted-foreground">View live server output and send commands</p>
        </div>
      </div>

      <div className="relative flex-1 border border-slate-200/40 dark:border-slate-700/30 rounded-lg overflow-hidden bg-slate-950">
        <ScrollArea
          className="flex h-full w-full text-emerald-400 p-4 font-mono text-sm"
          style={{
            userSelect: "none",
            WebkitTouchCallout: "none",
          }}
          onFocus={(e) => {
            const element = e.currentTarget;
            element.style.userSelect = "text";
          }}
          onBlur={(e) => {
            const element = e.currentTarget;
            element.style.userSelect = "none";
          }}
          tabIndex={0} // Make the ScrollArea focusable
        >
          <div
            ref={scrollAreaRef}
            className="console-content flex flex-col h-full pb-8" // Added pb-16 for bottom padding
            style={{ userSelect: "none", overflowY: "auto" }}
            onKeyDown={(e) => {
              // Prevent Ctrl+A selection in the console
              if ((e.ctrlKey || e.metaKey) && e.key === "a") {
                e.preventDefault();
              }
            }}
            onSelectCapture={(e) => {
              // Check if the selection was triggered by Ctrl+A
              if (document.getSelection()?.toString() === document.querySelector(".console-content")?.textContent) {
                // If entire content is selected, it's likely Ctrl+A
                e.preventDefault();
                // Clear selection
                window.getSelection()?.removeAllRanges();
              }
            }}
          >
            <div className="pb-1 text-xs text-slate-500 no-select">
              --- Server started on {new Date().toLocaleString()} ---
            </div>
            {consoleLines.map((line, i) => {
              return (
                <div key={i} className={cn("pb-1")}>
                  {line}
                </div>
              );
            })}
          </div>
        </ScrollArea>

        <form
          onSubmit={handleCommandSubmit}
          className="absolute bottom-0 left-0 right-0 border-t border-slate-700/50 bg-slate-950 px-4 py-2 no-select"
        >
          <div className="flex items-center">
            <span className="text-slate-500 mr-2">$</span>
            <input
              type="text"
              value={command}
              onChange={(e) => setCommand(e.target.value)}
              placeholder="Type a command and press Enter..."
              className="flex-1 bg-transparent border-none outline-none text-emerald-400 font-mono text-sm focus:outline-none focus:ring-0"
              autoFocus
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  handleCommandSubmit(e);
                }
              }}
            />
          </div>
        </form>
      </div>

      <div className="mt-2 text-xs text-slate-500 dark:text-slate-400 no-select">
        <span className="font-semibold">Tip:</span> Type{" "}
        <code className="px-1 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-indigo-700 dark:text-indigo-400">
          help
        </code>{" "}
        to see available commands
      </div>
    </div>
  );
});

export default ConsolePage;
