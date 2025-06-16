import "./app.css";
import React from "react";
import { Meta, Links, Outlet, Scripts, ScrollRestoration } from "react-router";

import ThemeProvider from "~/providers/ThemeProvider";
import SessionProvider from "~/providers/SessionProvider";

import { meta } from "~/lib/root-meta";
import { links } from "~/lib/root-links";
import { ErrorBoundary } from "~/components/ui/ErrorBoundary";

export { meta, links, ErrorBoundary };

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
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
