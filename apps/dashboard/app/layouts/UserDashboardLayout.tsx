import React from "react";
import {
  ChevronLeft,
  LayoutGrid,
  LogOut,
  Moon,
  PanelLeftIcon,
  Server,
  Sun,
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
];

const SidebarHeader: React.FC = () => {
  return (
    <div
      className={cn(
        "flex items-center gap-2 p-6 h-16 no-select",
        "border-b border-sidebar-border"
      )}
    >
      <div className="flex h-8 w-8 items-center justify-center rounded-md bg-indigo-600">
        <Server className="h-5 w-5 text-white" />
      </div>
      <span className="text-lg font-semibold">Panelium</span>
    </div>
  );
};

const SidebarNavigation: React.FC = () => {
  return (
    <ScrollArea className="flex-1 px-2 py-4">
      <nav className="flex flex-col w-full gap-1">
        {NAVIGATION_ITEMS.map((item) => {
          const isActive =
            location.pathname === item.href ||
            location.pathname.startsWith(item.href + "/");

          return (
            <Link
              key={item.href}
              to={item.href}
              className={cn(
                "flex w-full items-center gap-3 rounded-lg px-4 py-2.5",
                "text-sm font-medium no-select",
                "hover:bg-slate-100 dark:hover:bg-slate-700",
                "transition-all duration-200 ease",
                isActive &&
                  "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300"
              )}
            >
              <item.icon
                className={cn(
                  "h-5 w-5 text-slate-500 dark:text-slate-400",
                  isActive && "text-indigo-600 dark:text-indigo-400"
                )}
              />
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
            className="w-full justify-start p-2 h-14 no-select"
          >
            <EntityAvatar
              src=""
              alt="Avatar"
              title="John Doe"
              subTitle="john@example.com"
            />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56 no-select">
          <DropdownMenuItem className="text-red-500 focus:text-red-500">
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
      <div
        className={cn(
          "flex min-h-screen w-full",
          "bg-background text-foreground"
        )}
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
