import {index, layout, prefix, route, type RouteConfig} from "@react-router/dev/routes";
import {sampleAdminModule, sampleUserModule} from "./modules/sample-module"; // Import sample modules

// Helper to extract file routes from modules, ensuring they are in the correct format
const getModuleFileRoutes = (moduleRoutes: Array<{ path: string, file: string }> | undefined) => {
    if (!moduleRoutes) return [];
    return moduleRoutes.map(r => route(r.path, r.file));
};

const userModuleRoutes = getModuleFileRoutes(sampleUserModule.routes?.user);
const adminModuleRoutes = getModuleFileRoutes(sampleAdminModule.routes?.admin);

export default [
    ...prefix("auth", [
        layout("layouts/AuthLayout.tsx", [ // Added base path "auth"
            route("login", "routes/auth/login.tsx"),
        ])
    ]),

    layout("layouts/UserDashboardLayout.tsx", [
        index("routes/dashboard/overview.tsx"),
        route("server/*", "routes/dashboard/server/[id].tsx"),
        route("settings", "routes/dashboard/settings.tsx"),
        ...userModuleRoutes, // Spread user module routes here
    ]),

    ...prefix("admin", [
        layout("layouts/AdminDashboardLayout.tsx", [
            index("routes/admin/overview.tsx"),
            route("modules", "routes/admin/modules.tsx"),
            ...adminModuleRoutes, // Spread admin module routes here
        ]),
    ]),

    // route("*", "routes/not-found.tsx"),
] satisfies RouteConfig;