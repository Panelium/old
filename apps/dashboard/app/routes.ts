import {index, layout, prefix, route, type RouteConfig,} from "@react-router/dev/routes";

export default [
    ...prefix("auth", [
        layout("layouts/AuthLayout.tsx", [index("routes/auth/index.tsx")]),
    ]),

    layout("layouts/ProtectedLayout.tsx", [
        layout("layouts/UserDashboardLayout.tsx", [
            index("routes/dashboard/index.tsx"),
            route("server/*", "routes/dashboard/server/[id].tsx"),
        ]),
    ]),

    // route("*", "routes/not-found.tsx"),
] satisfies RouteConfig;
