import {index, layout, prefix, route, type RouteConfig} from "@react-router/dev/routes";

import {sampleAdminModule, sampleUserModule} from "./modules/sample-module";

const getModuleFileRoutes = (moduleRoutes: Array<{ path: string, file: string }> | undefined) => {
    if (!moduleRoutes) return [];
    return moduleRoutes.map(r => route(r.path, r.file));
};

const userModuleRoutes = getModuleFileRoutes(sampleUserModule.routes?.user);
const adminModuleRoutes = getModuleFileRoutes(sampleAdminModule.routes?.admin);

export default [
    ...prefix("auth", [
        layout("layouts/AuthLayout.tsx", [
            index("routes/auth/index.tsx"),
        ])
    ]),

    layout("layouts/ProtectedLayout.tsx", [
        layout("layouts/UserDashboardLayout.tsx", [
            index("routes/dashboard/index.tsx"),
            route("server/*", "routes/dashboard/server/[id].tsx"),
            ...userModuleRoutes,
        ]),

        ...prefix("admin", [
            layout("layouts/AdminDashboardLayout.tsx", [
                index("routes/admin/index.tsx"),
                ...adminModuleRoutes,
            ]),
        ]),
    ]),

    // route("*", "routes/not-found.tsx"),
] satisfies RouteConfig;