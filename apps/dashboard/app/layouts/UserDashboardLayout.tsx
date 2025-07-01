import React, { useState } from "react";
import {
  Activity,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  HardDrive,
  LayoutGrid,
  LogOut,
  Moon,
  PanelLeftIcon,
  ServerIcon,
  Settings,
  Sun,
  Terminal,
} from "lucide-react";
import { Link, Outlet } from "react-router-dom";
import { cn } from "~/lib/utils";

import { useTheme } from "~/providers/ThemeProvider";
import { useLogout } from "~/providers/SessionProvider";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { Button } from "~/components/ui/button";
import { ScrollArea } from "~/components/ui/scroll-area";
import EntityAvatar from "~/components/avatars/EntityAvatar";
import { Sidebar, SidebarProvider, useSidebar } from "~/components/ui/sidebar";

const NAVIGATION_ITEMS = [
  {
    title: "Overview",
    icon: LayoutGrid,
    href: "/",
  },
  {
    type: "seperator",
  },
  {
    text: "Your Servers",
    type: "header",
  },
  {
    title: "My Awesome Server 1",
    icon: ServerIcon,
    href: "/server/1",
    type: "server",
  },
  {
    title: "Web Server",
    icon: ServerIcon,
    href: "/server/2",
    type: "server",
  },
  {
    title: "Terraria Adventure",
    icon: ServerIcon,
    href: "/server/3",
    type: "server",
  },
  {
    title: "Rust Survival",
    icon: ServerIcon,
    href: "/server/4",
    type: "server",
  },
  {
    title: "ARK: Survival Evolved",
    icon: ServerIcon,
    href: "/server/5",
    type: "server",
  },
];

const SidebarHeader: React.FC = () => {
  return (
    <div
      className={cn(
        "flex items-center gap-0 p-6 pl-4 h-16 no-select",
        "border-b border-sidebar-border"
      )}
    >
      <div className="bg-tag-purple text-transparent h-12 min-w-2 m-0"></div>
      <img
        src="/logo/full-logo.svg"
        className="dark:filter-[invert()] mr-4 ml-2"
      />
    </div>
  );
};

const SidebarNavigation: React.FC = () => {
  function isActive(href: string | undefined): boolean {
    return (
      location.pathname === href || location.pathname.startsWith(href + "/")
    );
  }
  return (
    <ScrollArea className="flex-1 px-2 py-4">
      <nav className="flex flex-col w-full gap-1">
        {NAVIGATION_ITEMS.map((item) => {
          const IconComponent = item.icon;

          if (item.type !== null) {
            if (item.type === "seperator") {
              return <hr className="border-muted-foreground m-2" />;
            }
            if (item.type === "header") {
              return <h1 className="px-4 pt-1.5">{item.text}</h1>;
            }
            if (item.type === "server") {
              function SubButton({ item, value, children }: any) {
                return (
                  <Link
                    key={item.href + "/" + value}
                    to={item.href + "/" + value}
                    className={cn(
                      "flex w-full items-center gap-2 rounded-lg px-4 py-1.5",
                      "text-sm font-medium no-select",
                      "hover:bg-accent",
                      "transition-all duration-200 ease",
                      isActive(item.href + "/" + value) &&
                        "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300"
                    )}
                  >
                    <div className="w-4" />
                    {children}
                    <span className="capitalize">{value}</span>
                  </Link>
                );
              }

              const [isDroppedDown, setIsDroppedDown] = useState(false);

              return (
                <div className="group/navigation-family">
                  <div
                    className={cn(
                      "flex w-full items-center gap-0 rounded-lg",
                      "text-sm font-medium no-select",
                      "hover:bg-slate-100 dark:hover:bg-slate-700",
                      "transition-all duration-200 ease",
                      isActive(item.href) &&
                        "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300"
                    )}
                  >
                    <Link
                      key={item.href}
                      to={item.href ?? "#"}
                      className={cn(
                        "flex min-w-full items-center gap-3 rounded-lg px-4 py-2.5",
                        "text-sm font-medium no-select"
                      )}
                    >
                      {IconComponent && (
                        <IconComponent
                          className={cn(
                            "h-5 w-5 text-slate-500 dark:text-slate-400",
                            isActive(item.href) &&
                              "text-indigo-600 dark:text-indigo-400"
                          )}
                        />
                      )}
                      <span className="w-full">{item.title}</span>
                    </Link>
                    <button
                      className="ml-[calc(var(--spacing)*-7)] cursor-pointer hover:bg-sidebar-accent rounded-sm p-0.5 transition-all duration-300"
                      onClick={(event) => {
                        setIsDroppedDown(!isDroppedDown);
                      }}
                    >
                      {isDroppedDown ? (
                        <ChevronDown className="w-4 h-4" />
                      ) : (
                        <ChevronRight className="w-4 h-4" />
                      )}
                    </button>
                  </div>
                  <div className={isDroppedDown ? "" : "hidden"}>
                    <SubButton item={item} value="console">
                      <Terminal
                        className={cn(
                          "h-5 w-5 text-sidebar-accent-foreground",
                          isActive(item.href + "/console") &&
                            "text-indigo-600 dark:text-indigo-400"
                        )}
                      />
                    </SubButton>
                    <SubButton item={item} value="files">
                      <HardDrive
                        className={cn(
                          "h-5 w-5 text-sidebar-accent-foreground",
                          isActive(item.href + "/files") &&
                            "text-indigo-600 dark:text-indigo-400"
                        )}
                      />
                    </SubButton>
                    <SubButton item={item} value="activity">
                      <Activity
                        className={cn(
                          "h-5 w-5 text-sidebar-accent-foreground",
                          isActive(item.href + "/activity") &&
                            "text-indigo-600 dark:text-indigo-400"
                        )}
                      />
                    </SubButton>
                    <SubButton item={item} value="settings">
                      <Settings
                        className={cn(
                          "h-5 w-5 text-sidebar-accent-foreground",
                          isActive(item.href + "/settings") &&
                            "text-indigo-600 dark:text-indigo-400"
                        )}
                      />
                    </SubButton>
                  </div>
                </div>
              );
            }
          }
          return (
            <Link
              key={item.href}
              to={item.href ?? "#"}
              className={cn(
                "flex w-full items-center gap-3 rounded-lg px-4 py-2.5",
                "text-sm font-medium no-select",
                "hover:bg-slate-100 dark:hover:bg-slate-700",
                "transition-all duration-200 ease",
                isActive(item.href) &&
                  "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300"
              )}
            >
              {IconComponent && (
                <IconComponent
                  className={cn(
                    "h-5 w-5 text-slate-500 dark:text-slate-400",
                    isActive(item.href) &&
                      "text-indigo-600 dark:text-indigo-400"
                  )}
                />
              )}
              {item.title}
            </Link>
          );
        })}
      </nav>
    </ScrollArea>
  );
};

const SidebarDropdownMenu: React.FC = () => {
  const performLogout = useLogout();
  return (
    <div className="border-t border-sidebar-border p-4">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            className="group/button w-full justify-start p-2 h-14 no-select"
          >
            <EntityAvatar
              src=""
              alt="Avatar"
              title="Freddy Fazbear"
              subTitle="freddy.fazbear@freddys.pizza"
            />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56 no-select">
          <DropdownMenuItem className="bg-transparent focus:bg-transparent text-red-500 focus:text-red-500">
            <Button
              variant="ghost"
              className="w-full justify-start"
              onClick={() => performLogout()}
            >
              <LogOut className="mr-2 h-4 w-4" />
              <span>Logout</span>
            </Button>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

const SidebarTrigger: React.FC = () => {
  const { toggleSidebar, state } = useSidebar();

  const buttonProps = {
    variant: "ghost" as const,
    size: "icon" as const,
    onClick: toggleSidebar as () => void,
    "aria-label": "Toggle sidebar" as const,
  };

  return (
    <>
      <Button {...buttonProps} className="md:hidden no-select">
        <PanelLeftIcon className="h-6 w-6" />
      </Button>
      <Button {...buttonProps} className="hidden md:flex no-select">
        <ChevronLeft
          className={cn("h-5 w-5", state === "collapsed" && "rotate-180")}
        />
      </Button>
    </>
  );
};

const TopBar: React.FC = () => {
  const { theme, setTheme } = useTheme();
  const toggleTheme = () => {
    setTheme(theme === "light" ? "dark" : "light");
  };
  const ThemeIcon = theme === "light" ? Moon : Sun;

  return (
    <header
      className={cn(
        "sticky top-0 z-30 no-select",
        "flex items-center justify-between h-16 px-6 shadow-sm",
        "border-b border-border bg-topbar text-topbar-foreground"
      )}
    >
      <SidebarTrigger />
      <Button
        variant="ghost"
        size="icon"
        onClick={toggleTheme}
        aria-label="Toggle theme"
      >
        <ThemeIcon className="h-5 w-5" />
      </Button>
    </header>
  );
};

const UserDashboardLayout: React.FC = () => {
  return (
    <SidebarProvider defaultOpen={true}>
      <div className="background-gradient" />
      <div
        className={cn("flex min-h-screen w-full", "background text-foreground")}
      >
        <Sidebar
          className={cn(
            "flex flex-col h-full max-h-screen border-r",
            "border-sidebar-border b-white"
          )}
        >
          <SidebarHeader />
          <SidebarNavigation />
          <SidebarDropdownMenu />
        </Sidebar>
        <div className="flex-1">
          <TopBar />
          <Outlet />
        </div>
      </div>
    </SidebarProvider>
  );
};

export default UserDashboardLayout;
