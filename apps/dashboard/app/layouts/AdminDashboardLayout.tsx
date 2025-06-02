import React from 'react';
import {NavLink, Outlet} from 'react-router-dom';
import {Sidebar, SidebarProvider, useSidebar} from '~/components/ui/sidebar';
import {Button} from '~/components/ui/button';
import {
    LayoutGrid,
    Moon,
    PanelLeftIcon,
    Puzzle,
    Server as ServerIcon,
    Settings,
    ShieldCheck,
    Sun,
    UserCircle,
    Users
} from 'lucide-react';
import {cn} from '~/lib/utils';
import {useTheme} from "next-themes";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from '~/components/ui/dropdown-menu';

// Mobile Sidebar Trigger Component
function MobileSidebarTrigger() {
    const {toggleSidebar} = useSidebar();
    return (
        <Button
            variant="ghost"
            size="icon"
            className="md:hidden"
            onClick={toggleSidebar}
            aria-label="Toggle navigation menu"
        >
            <PanelLeftIcon className="h-6 w-6"/>
        </Button>
    );
}

// Theme Toggle Component
function ThemeToggle() {
    const {setTheme, theme} = useTheme();

    return (
        <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === "light" ? "dark" : "light")}
            aria-label="Toggle theme"
        >
            <Sun className="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"/>
            <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"/>
            <span className="sr-only">Toggle theme</span>
        </Button>
    );
}

export default function AdminDashboardLayout() {
    const navigation = [
        {name: 'Overview', href: '/admin', icon: LayoutGrid},
        {name: 'Users', href: '/admin/users', icon: Users},
        {name: 'Servers', href: '/admin/servers', icon: ServerIcon},
        {name: 'Nodes', href: '/admin/nodes', icon: ServerIcon},
        {name: 'Blueprints', href: '/admin/blueprints', icon: Puzzle},
        {name: 'System Settings', href: '/admin/system-settings', icon: Settings},
        {name: 'Modules', href: '/admin/modules', icon: Puzzle},
        {name: 'Security', href: '/admin/security', icon: ShieldCheck},
    ];

    return (
        <SidebarProvider defaultOpen={true}>
            <div className="flex min-h-screen w-full bg-muted/40">
                <Sidebar className="hidden border-r bg-background md:block">
                    <div className="flex h-full max-h-screen flex-col gap-2">
                        <div className="flex h-16 items-center border-b px-4 lg:px-6 shrink-0">
                            <NavLink to="/admin/overview" className="flex items-center gap-2 font-semibold">
                                <span className="">Panelium Admin</span>
                            </NavLink>
                        </div>
                        <div className="flex-1 overflow-y-auto">
                            <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
                                {navigation.map((item) => (
                                    <NavLink
                                        key={item.name}
                                        to={item.href}
                                        end
                                        className={({isActive}) =>
                                            cn(
                                                "flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary",
                                                isActive && "bg-muted text-primary"
                                            )
                                        }
                                    >
                                        <item.icon className="h-5 w-5"/>
                                        {item.name}
                                    </NavLink>
                                ))}
                            </nav>
                        </div>
                    </div>
                </Sidebar>

                <div className="flex flex-1 flex-col">
                    <header
                        className="sticky top-0 z-30 flex h-16 items-center gap-4 border-b bg-background px-4 sm:px-6">
                        <MobileSidebarTrigger/>
                        <div className="flex-1">
                            {/* Breadcrumbs or dynamic page title can go here */}
                        </div>
                        <ThemeToggle/>
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon" className="rounded-full">
                                    <UserCircle className="h-6 w-6"/>
                                    <span className="sr-only">Toggle admin menu</span>
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuLabel>Admin Actions</DropdownMenuLabel>
                                <DropdownMenuSeparator/>
                                <DropdownMenuItem>View Profile</DropdownMenuItem>
                                <DropdownMenuItem>Audit Logs</DropdownMenuItem>
                                <DropdownMenuSeparator/>
                                <DropdownMenuItem>Logout</DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </header>
                    <main className="flex-1 overflow-y-auto p-4 sm:p-6">
                        <Outlet/>
                    </main>
                </div>
            </div>
        </SidebarProvider>
    );
}
