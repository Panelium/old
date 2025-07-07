import React, { useEffect, useState } from "react";
import {
  Activity,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  HardDrive,
  LayoutGrid,
  LogOut,
  type LucideProps,
  Moon,
  PanelLeftIcon,
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
import EntityAvatar from "~/components/avatars/EntityAvatar";
import { Sidebar, SidebarProvider, useSidebar } from "~/components/ui/sidebar";
import { PagePressedEvent, pagesEventBus } from "~/components/dashboard/server/Pages";
import FilesPage from "~/components/dashboard/server/pages/FilesPage";
import { getClientClient } from "~/lib/api-clients";
import { ScrollArea } from "~/components/ui/scroll-area";

interface NavigationItemProps {
  title?: string;
  icon?: React.ForwardRefExoticComponent<Omit<LucideProps, "ref">>;
  href?: string;
  type?: string;
  text?: string;
}

const NAVIGATION_ITEMS: NavigationItemProps[] = [
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
];

// TODO: do this differently
// // Adds the mockup servers to NAVIGATION_ITEMS
// useDashboard().mockServers.forEach((server) => {
//   NAVIGATION_ITEMS.push({
//     title: server.name,
//     icon: ServerIcon,
//     href: "/server/" + server.id,
//     type: "server",
//   });
// });

const SidebarHeader: React.FC = () => {
  return (
    <div className={cn("flex items-center gap-0 p-6 pl-4 h-16 no-select", "border-b border-sidebar-border")}>
      <div className="bg-tag-purple text-transparent h-12 min-w-2 m-0"></div>
      <img src="/logo/full-logo.svg" className="dark:filter-[invert()] mr-4 ml-2" />
    </div>
  );
};

const SidebarNavigationItem: React.FC<{ item: NavigationItemProps }> = ({ item }) => {
  function isActive(href: string | undefined): boolean {
    return location.pathname === href || location.pathname.startsWith(href + "/");
  }

  const IconComponent = item.icon;

  return (
    <Link
      key={item.href}
      to={item.href ?? "#"}
      className={cn(
        "flex w-full items-center gap-3 rounded-lg px-4 py-2.5",
        "text-sm font-medium no-select",
        "hover:bg-slate-100 dark:hover:bg-slate-700",
        "transition-all duration-200 ease",
        isActive(item.href) && "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300"
      )}
      onClick={() => {
        if (item.href === undefined) {
          return;
        }
        const splitHref = item.href.split("/");
        if (item.type !== "server-tab") {
          pagesEventBus.dispatchEvent(new PagePressedEvent(FilesPage.id));
          return;
        }
        pagesEventBus.dispatchEvent(new PagePressedEvent(splitHref[splitHref.length - 1]));
      }}
    >
      {item.type === "server-tab" ? <div className="w-4" /> : <></>}
      {IconComponent && (
        <IconComponent
          className={cn(
            "h-5 w-5 text-slate-500 dark:text-slate-400",
            isActive(item.href) && "text-indigo-600 dark:text-indigo-400"
          )}
        />
      )}
      {item.title}
    </Link>
  );
};
const SidebarNavigation: React.FC = () => {
  return (
    <ScrollArea className="flex-1 px-2 py-4">
      <nav className="flex flex-col w-full gap-1">
        {NAVIGATION_ITEMS.map((item, idx) => {
          const key =
            (item.type ? item.type + "-" : "") +
            (item.title ? item.title + "-" : "") +
            (item.text ? item.text + "-" : "") +
            (item.href ? item.href : "") +
            "-" +
            idx;
          switch (item.type) {
            case "seperator":
              return <hr key={key} className="border-muted-foreground m-2" />;
            case "header":
              return (
                <h1 key={key} className="px-4 pt-1.5">
                  {item.text}
                </h1>
              );
            case "server": {
              const [isDroppedDown, setIsDroppedDown] = useState(false);
              return (
                <div key={key} className="group/navigation-family">
                  <div className="flex w-full items-center gap-0 rounded-lg">
                    <SidebarNavigationItem item={item} />
                    <button
                      className="ml-[calc(var(--spacing)*-7)] cursor-pointer hover:bg-sidebar-accent rounded-sm p-0.5 transition-all duration-300"
                      onClick={(event) => {
                        setIsDroppedDown(!isDroppedDown);
                      }}
                    >
                      {isDroppedDown ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
                    </button>
                  </div>
                  <div className={isDroppedDown ? "" : "hidden"}>
                    <SidebarNavigationItem
                      key={key + "-files"}
                      item={{
                        title: "Files",
                        icon: HardDrive,
                        href: item.href + "/files",
                        type: "server-tab",
                      }}
                    />
                    <SidebarNavigationItem
                      key={key + "-console"}
                      item={{
                        title: "Console",
                        icon: Terminal,
                        href: item.href + "/console",
                        type: "server-tab",
                      }}
                    />
                    <SidebarNavigationItem
                      key={key + "-activity"}
                      item={{
                        title: "Activity",
                        icon: Activity,
                        href: item.href + "/activity",
                        type: "server-tab",
                      }}
                    />
                    <SidebarNavigationItem
                      key={key + "-settings"}
                      item={{
                        title: "Settings",
                        icon: Settings,
                        href: item.href + "/settings",
                        type: "server-tab",
                      }}
                    />
                  </div>
                </div>
              );
            }
            default:
              return <SidebarNavigationItem key={key} item={item} />;
          }
        })}
      </nav>
    </ScrollArea>
  );
};

const SidebarDropdownMenu: React.FC = () => {
  const performLogout = useLogout();

  const [userInfo, setUserInfo] = useState<{
    uid: string;
    username: string;
    email: string;
    admin: boolean;
  } | null>(null);

  useEffect(() => {
    let mounted = true;
    getClientClient().then((client) =>
      client.getInfo({}).then((info) => {
        if (mounted) setUserInfo(info);
      })
    );
    return () => {
      mounted = false;
    };
  }, []);

  const { uid, username, email, admin } = userInfo || {
    uid: "",
    username: "Loading...",
    email: "Loading...",
    admin: false,
  };

  return (
    <div className="border-t border-sidebar-border p-4">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="group/button w-full justify-start p-2 h-14 no-select">
            <EntityAvatar src="" alt="Avatar" title={username} subTitle={email} />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56 no-select">
          {admin && (
            <DropdownMenuItem asChild>
              <Link to="/admin" className="w-full">
                <Button variant="ghost" className="w-full justify-start">
                  <Settings className="mr-2 h-4 w-4" />
                  <span>Admin Panel</span>
                </Button>
              </Link>
            </DropdownMenuItem>
          )}
          <DropdownMenuItem className="bg-transparent focus:bg-transparent text-red-500 focus:text-red-500">
            <Button variant="ghost" className="w-full justify-start" onClick={() => performLogout()}>
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
        <ChevronLeft className={cn("h-5 w-5", state === "collapsed" && "rotate-180")} />
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
      <Button variant="ghost" size="icon" onClick={toggleTheme} aria-label="Toggle theme">
        <ThemeIcon className="h-5 w-5" />
      </Button>
    </header>
  );
};

const UserDashboardLayout: React.FC = () => {
  return (
    <SidebarProvider defaultOpen={true}>
      <div className="background-gradient" />
      <div className={cn("flex min-h-screen w-full", "background text-foreground")}>
        <Sidebar className={cn("flex flex-col h-full max-h-screen border-r", "border-sidebar-border b-white")}>
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
