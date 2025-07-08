import { index, layout, prefix, route, type RouteConfig } from "@react-router/dev/routes";

export default [
  ...prefix("auth", [layout("layouts/AuthLayout.tsx", [index("routes/auth/index.tsx")])]),

  layout("layouts/ProtectedLayout.tsx", [
    layout("layouts/UserDashboardLayout.tsx", [
      index("routes/dashboard/index.tsx"),
      route("create-server", "routes/dashboard/create-server.tsx"),
      route("server/:id", "routes/dashboard/server/:id/index.tsx"),
      route("server/:id/:page", "routes/dashboard/server/:id/:page/index.tsx"),
    ]),
    ...prefix("admin", [layout("layouts/AdminLayout.tsx", [index("routes/admin/index.tsx")])]),
  ]),

  // route("*", "routes/not-found.tsx"),
] satisfies RouteConfig;
