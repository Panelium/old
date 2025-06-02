import React, {useEffect, useState} from 'react';
import {Link, Outlet, useLocation} from 'react-router-dom';
import {cn} from '~/lib/utils';
import {BarChart, LayoutGrid, LogOut, Moon, Server, Settings, Sun} from 'lucide-react';
import {Button} from '~/components/ui/button';
import {ScrollArea} from '~/components/ui/scroll-area';
import {Avatar, AvatarFallback, AvatarImage} from '~/components/ui/avatar';
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from '~/components/ui/dropdown-menu';

const navigationItems = [
    {
        title: 'Overview',
        icon: LayoutGrid,
        href: '/',
    },
    {
        title: 'Analytics',
        icon: BarChart,
        href: '/analytics',
    },
    {
        title: 'Settings',
        icon: Settings,
        href: '/settings',
    },
];

export default function UserDashboardLayout() {
    const location = useLocation();
    const [theme, setTheme] = useState<'light' | 'dark'>('light');
    const [sidebarOpen, setSidebarOpen] = useState(false);

    // Toggle theme
    const toggleTheme = () => {
        const newTheme = theme === 'light' ? 'dark' : 'light';
        setTheme(newTheme);
        document.documentElement.classList.toggle('dark', newTheme === 'dark');
        localStorage.setItem('theme', newTheme);
    };

    // Initialize theme from localStorage
    useEffect(() => {
        const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' || 'light';
        setTheme(savedTheme);
        document.documentElement.classList.toggle('dark', savedTheme === 'dark');
    }, []);

    return (
        <div className="flex h-screen bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-50">
            {/* Sidebar - Desktop */}
            <aside
                className="fixed inset-y-0 z-50 flex w-72 flex-col bg-white dark:bg-slate-800 shadow-lg transition-transform duration-300 lg:relative lg:translate-x-0">
                {/* Logo */}
                <div className="flex h-16 items-center gap-2 border-b border-slate-200 dark:border-slate-700 px-6">
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
                                        "flex items-center gap-3 rounded-lg px-4 py-2.5 text-sm font-medium transition-all hover:bg-slate-100 dark:hover:bg-slate-700",
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
                            <Button variant="ghost" className="w-full justify-start px-2">
                                <Avatar className="h-8 w-8 mr-2">
                                    <AvatarImage src="" alt="Avatar"/>
                                    <AvatarFallback
                                        className="bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300">JD</AvatarFallback>
                                </Avatar>
                                <div className="flex flex-col items-start text-sm">
                                    <span className="font-medium">John Doe</span>
                                    <span className="text-xs text-slate-500 dark:text-slate-400">john@example.com</span>
                                </div>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-56">
                            <DropdownMenuItem className="text-red-500 focus:text-red-500">
                                <LogOut className="mr-2 h-4 w-4"/>
                                <span>Logout</span>
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 overflow-hidden">
                {/* Top Bar */}
                <header
                    className="flex h-16 items-center justify-between border-b border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-6 shadow-sm">
                    <div className="flex items-center gap-4">
                        <Button variant="ghost" size="icon" onClick={toggleTheme}>
                            {theme === 'light' ? <Moon className="h-5 w-5"/> : <Sun className="h-5 w-5"/>}
                        </Button>
                    </div>
                </header>

                {/* Page Content */}
                <div className="overflow-auto h-[calc(100vh-4rem)]">
                    <Outlet/>
                </div>
            </main>
        </div>
    );
}
