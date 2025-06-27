import "./app.css";
import React from "react";
import { Meta, Links, Outlet, Scripts, ScrollRestoration } from "react-router";

import { meta } from "~/lib/root-meta";
import { links } from "~/lib/root-links";

import ThemeProvider from "~/providers/ThemeProvider";
import SessionProvider from "~/providers/SessionProvider";
import { ErrorBoundary } from "~/components/ui/ErrorBoundary";

export { meta, links, ErrorBoundary };

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" type="image/svg" href="/public/favicon.svg"/>
        <Meta />
        <Links />
      </head>
      <body>
        <ThemeProvider>
          <SessionProvider>{children}</SessionProvider>
        </ThemeProvider>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}
