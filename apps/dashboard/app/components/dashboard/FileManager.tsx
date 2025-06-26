import React, {useState} from 'react';
import {Card, CardContent} from '~/components/ui/card';
import {Button} from '~/components/ui/button';
import {Input} from '~/components/ui/input';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from '~/components/ui/table';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from '~/components/ui/dialog';
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbSeparator
} from '~/components/ui/breadcrumb';
import {
    ArrowLeft,
    Check,
    ChevronDown,
    Download,
    Edit,
    File,
    FileArchive,
    FileAudio,
    FileCode,
    FileImage,
    FilePlus,
    FileText,
    FileVideo,
    Folder,
    FolderPlus,
    LayoutGrid,
    List,
    MoreHorizontal,
    RefreshCw,
    Search,
    Trash,
    Upload
} from 'lucide-react';
import {Tabs, TabsList, TabsTrigger} from '~/components/ui/tabs';
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from '~/components/ui/tooltip';
import {cn} from '~/lib/utils';
import {ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger} from '~/components/ui/context-menu';
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from '~/components/ui/dropdown-menu';

// Sample file data structure
const sampleFiles = [
    {name: 'server.properties', type: 'file', size: '2.5 KB', modified: '2025-05-30 14:22', icon: FileText},
    {name: 'plugins', type: 'folder', items: 32, modified: '2025-05-29 09:10', icon: Folder},
    {name: 'world', type: 'folder', items: 156, modified: '2025-06-01 18:45', icon: Folder},
    {name: 'server.jar', type: 'file', size: '35.8 MB', modified: '2025-05-15 11:30', icon: FileCode},
    {name: 'backup.zip', type: 'file', size: '120.2 MB', modified: '2025-05-28 22:15', icon: FileArchive},
    {name: 'server-icon.png', type: 'file', size: '16.4 KB', modified: '2025-05-10 13:50', icon: FileImage},
    {name: 'logs', type: 'folder', items: 8, modified: '2025-06-01 20:30', icon: Folder},
    {name: 'config.yml', type: 'file', size: '4.7 KB', modified: '2025-05-30 15:10', icon: FileText},
];

// Helper function to determine file icon based on extension
const getFileIcon = (fileName: string) => {
    const extension = fileName.split('.').pop()?.toLowerCase() || '';

    const iconMap: {
        [key: string]: React.ComponentType<{ className?: string }>;
    } = {
        txt: FileText,
        log: FileText,
        properties: FileText,
        yml: FileText,
        yaml: FileText,
        json: FileText,

        js: FileCode,
        ts: FileCode,
        jsx: FileCode,
        tsx: FileCode,
        java: FileCode,
        py: FileCode,
        html: FileCode,
        css: FileCode,

        zip: FileArchive,
        tar: FileArchive,
        gz: FileArchive,
        rar: FileArchive,

        png: FileImage,
        jpg: FileImage,
        jpeg: FileImage,
        gif: FileImage,
        svg: FileImage,

        mp3: FileAudio,
        wav: FileAudio,
        ogg: FileAudio,

        mp4: FileVideo,
        avi: FileVideo,
        mov: FileVideo,
        mkv: FileVideo,
    };

    return iconMap[extension] || File;
};

interface FileManagerProps {
    serverId: string;
    title?: string;
    className?: string;
}

export function FileManager({serverId, title = 'File Manager', className}: FileManagerProps) {
    const [currentPath, setCurrentPath] = useState<string[]>([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);
    const [sortBy, setSortBy] = useState<'name' | 'size' | 'modified'>('name');
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
    const [viewMode, setViewMode] = useState<'list' | 'grid'>('list');
    const [isUploading, setIsUploading] = useState(false);
    const [showNewFileDialog, setShowNewFileDialog] = useState(false);
    const [newFileName, setNewFileName] = useState('');
    const [newFileContent, setNewFileContent] = useState('');

    // Sort files based on current sort settings
    const sortedFiles = [...sampleFiles].sort((a, b) => {
        // Folders always come first
        if (a.type !== b.type) {
            return a.type === 'folder' ? -1 : 1;
        }

        // Then sort by the selected column
        if (sortBy === 'name') {
            return sortDirection === 'asc'
                ? a.name.localeCompare(b.name)
                : b.name.localeCompare(a.name);
        } else if (sortBy === 'modified') {
            return sortDirection === 'asc'
                ? a.modified.localeCompare(b.modified)
                : b.modified.localeCompare(a.modified);
        } else if (sortBy === 'size' && 'size' in a && 'size' in b) {
            // For simplicity - in a real app you'd parse the size values
            return sortDirection === 'asc'
                ? a.size!.localeCompare(b.size!)
                : b.size!.localeCompare(a.size!);
        }
        return 0;
    });

    // Filter files based on search query
    const filteredFiles = sortedFiles.filter(file =>
        file.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    // Handle navigation to a folder
    const navigateToFolder = (folderName: string) => {
        setCurrentPath([...currentPath, folderName]);
        setSelectedFiles([]);
    };

    // Handle navigation back
    const navigateBack = () => {
        if (currentPath.length > 0) {
            setCurrentPath(currentPath.slice(0, -1));
            setSelectedFiles([]);
        }
    };

    // Handle file selection
    const toggleFileSelection = (fileName: string) => {
        if (selectedFiles.includes(fileName)) {
            setSelectedFiles(selectedFiles.filter(name => name !== fileName));
        } else {
            setSelectedFiles([...selectedFiles, fileName]);
        }
    };

    // Handle select all files
    const toggleSelectAll = () => {
        if (selectedFiles.length === filteredFiles.length) {
            setSelectedFiles([]);
        } else {
            setSelectedFiles(filteredFiles.map(file => file.name));
        }
    };

    // Handle sort change
    const handleSortChange = (column: 'name' | 'size' | 'modified') => {
        if (sortBy === column) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
        } else {
            setSortBy(column);
            setSortDirection('asc');
        }
    };

    // Handle file upload
    const handleFileUpload = () => {
        setIsUploading(true);
        // Simulate file upload
        setTimeout(() => {
            setIsUploading(false);
        }, 2000);
    };

    // Handle new file creation
    const handleCreateNewFile = () => {
        console.log(`Creating new file: ${newFileName} with content length: ${newFileContent.length}`);
        setShowNewFileDialog(false);
        setNewFileName('');
        setNewFileContent('');
    };

    return (
        <Card className={cn("shadow-sm py-0", className)}>
            <CardContent className="p-0">
                {/* File Manager Header */}
                <div className="border-b border-border p-4">
                    <div className="flex items-center justify-between gap-4 mb-3">
                        <h3 className="text-lg font-medium">{title}</h3>
                        <div className="flex items-center gap-2">
                            <Tabs value={viewMode} onValueChange={(value) => setViewMode(value as 'list' | 'grid')}>
                                <TabsList className="h-9">
                                    <TabsTrigger value="list" className="px-3">
                                        <List className="w-4 h-4"/>
                                    </TabsTrigger>
                                    <TabsTrigger value="grid" className="px-3">
                                        <LayoutGrid className="w-4 h-4"/>
                                    </TabsTrigger>
                                </TabsList>
                            </Tabs>
                        </div>
                    </div>

                    {/* Breadcrumb */}
                    <div className="flex items-center gap-2 mb-3">
                        <Button
                            variant="outline"
                            size="icon"
                            onClick={navigateBack}
                            disabled={currentPath.length === 0}
                            className="h-8 w-8"
                        >
                            <ArrowLeft className="h-4 w-4"/>
                        </Button>
                        <Breadcrumb>
                            <BreadcrumbList>
                                <BreadcrumbItem>
                                    <BreadcrumbLink
                                        onClick={() => setCurrentPath([])}
                                    >
                                        home
                                    </BreadcrumbLink>
                                </BreadcrumbItem>
                                {currentPath.map((folder, index) => (
                                    <React.Fragment key={index}>
                                        <BreadcrumbSeparator/>
                                        <BreadcrumbItem>
                                            <BreadcrumbLink
                                                onClick={() => setCurrentPath(currentPath.slice(0, index + 1))}
                                            >
                                                {folder}
                                            </BreadcrumbLink>
                                        </BreadcrumbItem>
                                    </React.Fragment>
                                ))}
                            </BreadcrumbList>
                        </Breadcrumb>
                        <Button variant="ghost" size="icon" className="h-8 w-8 ml-1">
                            <RefreshCw className="h-4 w-4"/>
                        </Button>
                    </div>

                    {/* Actions Bar */}
                    <div className="flex flex-col sm:flex-row gap-3">
                        <div className="relative flex-1">
                            <Search
                                className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500 dark:text-slate-400"/>
                            <Input
                                type="text"
                                placeholder="Search files..."
                                className="pl-10 h-9"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                            />
                        </div>
                        <div className="flex gap-2">
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Button variant="outline" size="sm" className="h-9" onClick={handleFileUpload}
                                                disabled={isUploading}>
                                            {isUploading ? <RefreshCw className="h-4 w-4 animate-spin"/> :
                                                <Upload className="h-4 w-4 mr-1"/>}
                                            <span className="hidden sm:inline-block">Upload</span>
                                        </Button>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>Upload files</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>

                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="outline" size="sm" className="h-9">
                                        <FilePlus className="h-4 w-4 mr-1"/>
                                        <span className="hidden sm:inline-block">New</span>
                                        <ChevronDown className="h-4 w-4 ml-1"/>
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem onClick={() => setShowNewFileDialog(true)}>
                                        <FileText className="h-4 w-4 mr-2"/>
                                        New File
                                    </DropdownMenuItem>
                                    <DropdownMenuItem>
                                        <FolderPlus className="h-4 w-4 mr-2"/>
                                        New Folder
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    </div>
                </div>

                {/* File List */}
                {viewMode === 'list' ? (
                    <div className="overflow-x-auto">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-10">
                                        <input
                                            type="checkbox"
                                            checked={selectedFiles.length === filteredFiles.length && filteredFiles.length > 0}
                                            onChange={toggleSelectAll}
                                            className="h-4 w-4 rounded border-slate-300 dark:border-slate-700"
                                        />
                                    </TableHead>
                                    <TableHead
                                        onClick={() => handleSortChange('name')}
                                    >
                                        <div className="flex items-center">
                                            Name
                                            {sortBy === 'name' && (
                                                <ChevronDown
                                                    className={cn(
                                                        "h-4 w-4 ml-1",
                                                        sortDirection === 'desc' && "rotate-180"
                                                    )}
                                                />
                                            )}
                                        </div>
                                    </TableHead>
                                    <TableHead
                                        onClick={() => handleSortChange('size')}
                                    >
                                        <div className="flex items-center">
                                            Size
                                            {sortBy === 'size' && (
                                                <ChevronDown
                                                    className={cn(
                                                        "h-4 w-4 ml-1",
                                                        sortDirection === 'desc' && "rotate-180"
                                                    )}
                                                />
                                            )}
                                        </div>
                                    </TableHead>
                                    <TableHead
                                        onClick={() => handleSortChange('modified')}
                                    >
                                        <div className="flex items-center">
                                            Last Modified
                                            {sortBy === 'modified' && (
                                                <ChevronDown
                                                    className={cn(
                                                        "h-4 w-4 ml-1",
                                                        sortDirection === 'desc' && "rotate-180"
                                                    )}
                                                />
                                            )}
                                        </div>
                                    </TableHead>
                                    <TableHead className="w-10"></TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {filteredFiles.length > 0 ? (
                                    filteredFiles.map((file) => (
                                        <ContextMenu key={file.name}>
                                            <ContextMenuTrigger asChild>
                                                <TableRow
                                                    className={cn(
                                                        "cursor-pointer transition-colors",
                                                        selectedFiles.includes(file.name) && "bg-slate-100 dark:bg-slate-800",
                                                        file.type === 'folder' && "hover:bg-slate-50 dark:hover:bg-slate-800/60"
                                                    )}
                                                    onClick={() => file.type === 'folder' ? navigateToFolder(file.name) : toggleFileSelection(file.name)}
                                                    onDoubleClick={() => file.type !== 'folder' && console.log(`Opening file: ${file.name}`)}
                                                >
                                                    <TableCell className="p-2">
                                                        <input
                                                            type="checkbox"
                                                            checked={selectedFiles.includes(file.name)}
                                                            onChange={(e) => {
                                                                e.stopPropagation();
                                                                toggleFileSelection(file.name);
                                                            }}
                                                            className="h-4 w-4 rounded border-slate-300 dark:border-slate-700"
                                                        />
                                                    </TableCell>
                                                    <TableCell className="p-2">
                                                        <div className="flex items-center gap-2">
                                                            <file.icon
                                                                className="h-4 w-4 text-slate-500 dark:text-slate-400 flex-shrink-0"/>
                                                            <span>{file.name}</span>
                                                        </div>
                                                    </TableCell>
                                                    <TableCell
                                                        className="p-2 text-slate-500 dark:text-slate-400 text-sm">
                                                        {file.type === 'folder' ? `${file.items} items` : file.size}
                                                    </TableCell>
                                                    <TableCell
                                                        className="p-2 text-slate-500 dark:text-slate-400 text-sm">
                                                        {file.modified}
                                                    </TableCell>
                                                    <TableCell className="p-2">
                                                        <DropdownMenu>
                                                            <DropdownMenuTrigger asChild>
                                                                <Button variant="ghost" size="icon" className="h-8 w-8">
                                                                    <MoreHorizontal className="h-4 w-4"/>
                                                                </Button>
                                                            </DropdownMenuTrigger>
                                                            <DropdownMenuContent align="end">
                                                                {file.type !== 'folder' && (
                                                                    <DropdownMenuItem>
                                                                        <Edit className="h-4 w-4 mr-2"/>
                                                                        Edit
                                                                    </DropdownMenuItem>
                                                                )}
                                                                {file.type !== 'folder' && (
                                                                    <DropdownMenuItem>
                                                                        <Download className="h-4 w-4 mr-2"/>
                                                                        Download
                                                                    </DropdownMenuItem>
                                                                )}
                                                                <DropdownMenuItem
                                                                    className="text-red-500 focus:text-red-500">
                                                                    <Trash className="w-4 h-4 mr-2"/>
                                                                    Delete
                                                                </DropdownMenuItem>
                                                            </DropdownMenuContent>
                                                        </DropdownMenu>
                                                    </TableCell>
                                                </TableRow>
                                            </ContextMenuTrigger>
                                            <ContextMenuContent>
                                                {file.type !== 'folder' && (
                                                    <ContextMenuItem>
                                                        <Edit className="h-4 w-4 mr-2"/>
                                                        Edit
                                                    </ContextMenuItem>
                                                )}
                                                {file.type !== 'folder' && (
                                                    <ContextMenuItem>
                                                        <Download className="h-4 w-4 mr-2"/>
                                                        Download
                                                    </ContextMenuItem>
                                                )}
                                                <ContextMenuItem className="text-red-500 focus:text-red-500">
                                                    <Trash className="w-4 h-4 mr-2"/>
                                                    Delete
                                                </ContextMenuItem>
                                            </ContextMenuContent>
                                        </ContextMenu>
                                    ))
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={5} className="h-32 text-center">
                                            <div className="flex flex-col items-center justify-center">
                                                <div className="rounded-full bg-slate-100 dark:bg-slate-800 p-3 mb-3">
                                                    <File className="h-6 w-6 text-slate-400"/>
                                                </div>
                                                <p className="text-sm text-slate-500 dark:text-slate-400">No files
                                                    found</p>
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </div>
                ) : (
                    <div className="p-4">
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                            {filteredFiles.length > 0 ? (
                                filteredFiles.map((file) => (
                                    <ContextMenu key={file.name}>
                                        <ContextMenuTrigger asChild>
                                            <div
                                                className={cn(
                                                    "flex flex-col items-center justify-center p-4 rounded-lg border border-border transition-all",
                                                    selectedFiles.includes(file.name) && "bg-slate-100 dark:bg-slate-800 border-slate-300 dark:border-slate-600",
                                                    file.type === 'folder' ? "hover:bg-slate-50 dark:hover:bg-slate-800/60" : "hover:border-slate-300 dark:hover:border-slate-600",
                                                    "cursor-pointer"
                                                )}
                                                onClick={() => file.type === 'folder' ? navigateToFolder(file.name) : toggleFileSelection(file.name)}
                                                onDoubleClick={() => file.type !== 'folder' && console.log(`Opening file: ${file.name}`)}
                                            >
                                                <div className="relative mb-2">
                                                    <file.icon
                                                        className="h-10 w-10 text-slate-500 dark:text-slate-400"/>
                                                    {selectedFiles.includes(file.name) && (
                                                        <div
                                                            className="absolute -top-2 -right-2 h-5 w-5 bg-indigo-600 rounded-full flex items-center justify-center">
                                                            <Check className="w-3 h-3 text-white"/>
                                                        </div>
                                                    )}
                                                </div>
                                                <span
                                                    className="text-sm font-medium text-center line-clamp-1">{file.name}</span>
                                                <span className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                                                    {file.type === 'folder' ? `${file.items} items` : file.size}
                                                </span>
                                            </div>
                                        </ContextMenuTrigger>
                                        <ContextMenuContent>
                                            {file.type !== 'folder' && (
                                                <ContextMenuItem>
                                                    <Edit className="h-4 w-4 mr-2"/>
                                                    Edit
                                                </ContextMenuItem>
                                            )}
                                            {file.type !== 'folder' && (
                                                <ContextMenuItem>
                                                    <Download className="h-4 w-4 mr-2"/>
                                                    Download
                                                </ContextMenuItem>
                                            )}
                                            <ContextMenuItem className="text-red-500 focus:text-red-500">
                                                <Trash className="w-4 h-4 mr-2"/>
                                                Delete
                                            </ContextMenuItem>
                                        </ContextMenuContent>
                                    </ContextMenu>
                                ))
                            ) : (
                                <div className="col-span-full h-32 flex flex-col items-center justify-center">
                                    <div className="rounded-full bg-slate-100 dark:bg-slate-800 p-3 mb-3">
                                        <File className="h-6 w-6 text-slate-400"/>
                                    </div>
                                    <p className="text-sm text-slate-500 dark:text-slate-400">No files found</p>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </CardContent>

            {/* New File Dialog */}
            <Dialog open={showNewFileDialog} onOpenChange={setShowNewFileDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Create New File</DialogTitle>
                        <DialogDescription>
                            Enter a name and content for your new file.
                        </DialogDescription>
                    </DialogHeader>

                    <div className="space-y-4 py-2">
                        <div className="space-y-2">
                            <label htmlFor="filename" className="text-sm font-medium">
                                File Name
                            </label>
                            <Input
                                id="filename"
                                value={newFileName}
                                onChange={(e) => setNewFileName(e.target.value)}
                                placeholder="Enter file name"
                            />
                        </div>

                        <div className="space-y-2">
                            <label htmlFor="content" className="text-sm font-medium">
                                Content
                            </label>
                            <textarea
                                id="content"
                                value={newFileContent}
                                onChange={(e) => setNewFileContent(e.target.value)}
                                placeholder="Enter file content"
                                className="w-full min-h-[200px] rounded-md border border-border bg-white dark:bg-slate-800 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-indigo-500"
                            />
                        </div>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setShowNewFileDialog(false)}>
                            Cancel
                        </Button>
                        <Button onClick={handleCreateNewFile}>
                            Create File
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </Card>
    );
}
