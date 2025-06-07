import React from 'react';
import {Link, Outlet, useLocation} from 'react-router-dom';
import {cn} from '~/lib/utils';
import {ChevronLeft, LayoutGrid, LogOut, Moon, PanelLeftIcon, Server, Sun} from 'lucide-react';
import {Button} from '~/components/ui/button';
import {ScrollArea} from '~/components/ui/scroll-area';
import {Avatar, AvatarFallback, AvatarImage} from '~/components/ui/avatar';
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from '~/components/ui/dropdown-menu';
import {Sidebar, SidebarProvider, useSidebar} from '~/components/ui/sidebar';
import {useTheme} from '~/providers/ThemeProvider';
import {useLogout} from "~/providers/SessionProvider";

const navigationItems = [
    {
        title: 'Overview',
        icon: LayoutGrid,
        href: '/',
    },
    // {
    //     title: 'Analytics',
    //     icon: BarChart,
    //     href: '/analytics',
    // },
    // {
    //     title: 'Settings',
    //     icon: Settings,
    //     href: '/settings',
    // },
];

// Mobile Sidebar Trigger Component
function MobileSidebarTrigger() {
    const {toggleSidebar} = useSidebar();
    return (
        <Button
            variant="ghost"
            size="icon"
            className="md:hidden no-select"
            onClick={toggleSidebar}
            aria-label="Toggle navigation menu"
        >
            <PanelLeftIcon className="h-6 w-6"/>
        </Button>
    );
}

// Desktop Sidebar Toggle Component
function DesktopSidebarToggle() {
    const {toggleSidebar, state} = useSidebar();
    return (
        <Button
            variant="ghost"
            size="icon"
            className="hidden md:flex no-select"
            onClick={toggleSidebar}
            aria-label="Toggle sidebar"
        >
            <ChevronLeft className={cn(
                "h-5 w-5 transition-transform",
                state === "collapsed" && "rotate-180"
            )}/>
        </Button>
    );
}

export default function UserDashboardLayout() {
    const location = useLocation();
    const {theme, setTheme} = useTheme();
    const performLogout = useLogout();

    // Toggle theme
    const toggleTheme = () => {
        setTheme(theme === 'light' ? 'dark' : 'light');
    };

    async function handleLogout() {
        await performLogout();
    }

    return (
        <SidebarProvider defaultOpen={true}>
            <div className="flex min-h-screen w-full bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-50">
                <Sidebar
                    className="hidden border-r dark:border-slate-700 border-slate-200 bg-white dark:bg-slate-800 md:block">
                    <div className="flex h-full max-h-screen flex-col">
                        {/* Logo */}
                        <div
                            className="flex h-16 items-center gap-2 border-b border-slate-200 dark:border-slate-700 px-6 no-select">
                            <div className="flex h-8 w-8 items-center justify-center rounded-md bg-indigo-600">
                                <Server className="h-5 w-5 text-white"/>
                            </div>
                            <span className="text-lg font-semibold">Panelium</span>
                        </div>

                        {/* Navigation */}
                        <ScrollArea className="flex-1 py-4">
                            <nav className="grid gap-1 px-2">
                                {navigationItems.map((item) => {
                                    const isActive = location.pathname === item.href || location.pathname.startsWith(item.href + '/');

                                    return (
                                        <Link
                                            key={item.href}
                                            to={item.href}
                                            className={cn(
                                                "flex items-center gap-3 rounded-lg px-4 py-2.5 text-sm font-medium transition-all hover:bg-slate-100 dark:hover:bg-slate-700 no-select",
                                                isActive && "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300"
                                            )}
                                        >
                                            <item.icon className={cn(
                                                "h-5 w-5 text-slate-500 dark:text-slate-400",
                                                isActive && "text-indigo-600 dark:text-indigo-400"
                                            )}/>
                                            {item.title}
                                        </Link>
                                    );
                                })}
                            </nav>
                        </ScrollArea>

                        {/* User Menu */}
                        <div className="border-t border-slate-200 dark:border-slate-700 p-4">
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" className="w-full justify-start px-2 no-select">
                                        <Avatar className="h-8 w-8 mr-2">
                                            <AvatarImage src="" alt="Avatar"/>
                                            <AvatarFallback
                                                className="bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300">JD</AvatarFallback>
                                        </Avatar>
                                        <div className="flex flex-col items-start text-sm">
                                            <span className="font-medium">John Doe</span>
                                            <span
                                                className="text-xs text-slate-500 dark:text-slate-400">john@example.com</span>
                                        </div>
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end" className="w-56 no-select">
                                    <DropdownMenuItem className="text-red-500 focus:text-red-500">
                                        <Button
                                            variant="ghost"
                                            className="w-full justify-start"
                                            onClick={handleLogout}
                                        >
                                            <LogOut className="mr-2 h-4 w-4"/>
                                            <span>Logout</span>
                                        </Button>
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    </div>
                </Sidebar>

                {/* Main Content */}
                <div className="flex flex-1 flex-col">
                    {/* Top Bar */}
                    <header
                        className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-6 shadow-sm no-select">
                        <div className="flex items-center gap-4">
                            <MobileSidebarTrigger/>
                            <DesktopSidebarToggle/>
                            {/* Additional header content can go here */}
                        </div>
                        <Button variant="ghost" size="icon" onClick={toggleTheme} className="no-select"
                                aria-label="Toggle theme">
                            {theme === 'light' ? <Moon className="h-5 w-5"/> : <Sun className="h-5 w-5"/>}
                        </Button>
                    </header>

                    {/* Page Content */}
                    <div className="overflow-auto h-[calc(100vh-4rem)]">
                        <Outlet/>
                    </div>
                </div>
            </div>
        </SidebarProvider>
    );
}
